package sports

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// TennisScoreboard represents the response from ESPN tennis scoreboard API

func padRight(s string, n int) string {
	r := []rune(s)
	if len(r) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(r))
}

func classFor(code string) string {
	switch code {
	case ansiGreen:
		return "term-green"
	case ansiGray:
		return "term-gray"
	case ansiBold:
		return "term-bold"
	case ansiBold + ansiGreen:
		return "term-bold term-green"
	case ansiBold + ansiRed:
		return "term-bold term-red"
	}
	return ""
}

type TennisScoreboard struct {
	Events []TennisEvent `json:"events"`
}

type TennisEvent struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	ShortName string           `json:"shortName"`
	Groupings []TennisGrouping `json:"groupings"`
}

type TennisGrouping struct {
	Grouping struct {
		ID          string `json:"id"`
		Slug        string `json:"slug"`
		DisplayName string `json:"displayName"`
	} `json:"grouping"`
	Competitions []TennisCompetition `json:"competitions"`
}

type TennisCompetition struct {
	ID          string             `json:"id"`
	Date        string             `json:"date"`
	Status      TennisStatus       `json:"status"`
	Venue       TennisVenue        `json:"venue"`
	Competitors []TennisCompetitor `json:"competitors"`
	Round       TennisRound        `json:"round"`
	Notes       []TennisNote       `json:"notes"`
	Type        TennisType         `json:"type"`
}

type TennisStatus struct {
	Period int `json:"period"`
	Type   struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		State       string `json:"state"`
		Completed   bool   `json:"completed"`
		Description string `json:"description"`
		Detail      string `json:"detail"`
		ShortDetail string `json:"shortDetail"`
	} `json:"type"`
}

type TennisVenue struct {
	FullName string `json:"fullName"`
	Court    string `json:"court"`
}

type TennisRound struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
}

type TennisNote struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type TennisType struct {
	ID   string `json:"id"`
	Text string `json:"text"`
	Slug string `json:"slug"`
}

type TennisCompetitor struct {
	ID         string `json:"id"`
	HomeAway   string `json:"homeAway"`
	Winner     bool   `json:"winner"`
	Possession bool   `json:"possession"`
	Linescores []struct {
		Value    float64 `json:"value"`
		Winner   bool    `json:"winner"`
		Tiebreak float64 `json:"tiebreak,omitempty"`
	} `json:"linescores"`
	Athlete struct {
		DisplayName string `json:"displayName"`
		ShortName   string `json:"shortName"`
		FullName    string `json:"fullName"`
		Flag        struct {
			Href string   `json:"href"`
			Alt  string   `json:"alt"`
			Rel  []string `json:"rel"`
		} `json:"flag"`
	} `json:"athlete"`
	Roster struct {
		DisplayName      string `json:"displayName"`
		ShortDisplayName string `json:"shortDisplayName"`
	} `json:"roster"`
}

type MergedTournament struct {
	Tour  string
	Event TennisEvent
}

func getCompetitorName(c TennisCompetitor) string {
	if c.Roster.ShortDisplayName != "" {
		return c.Roster.ShortDisplayName
	}
	if c.Roster.DisplayName != "" {
		return c.Roster.DisplayName
	}
	if c.Athlete.ShortName != "" {
		return c.Athlete.ShortName
	}
	if c.Athlete.DisplayName != "" {
		return c.Athlete.DisplayName
	}
	return "Unknown Player"
}

func getCompetitorFullName(c TennisCompetitor) string {
	if c.Roster.DisplayName != "" {
		return c.Roster.DisplayName
	}
	if c.Athlete.FullName != "" {
		return c.Athlete.FullName
	}
	if c.Athlete.DisplayName != "" {
		return c.Athlete.DisplayName
	}
	return "Unknown Player"
}

func abbreviateRound(r string) string {
	r = strings.ToLower(r)
	if strings.Contains(r, "round of 16") {
		return "R16"
	}
	if strings.Contains(r, "quarterfinal") {
		return "QF"
	}
	if strings.Contains(r, "semifinal") {
		return "SF"
	}
	if strings.Contains(r, "final") {
		return "F"
	}
	if strings.Contains(r, "round 1") {
		return "R1"
	}
	if strings.Contains(r, "round 2") {
		return "R2"
	}
	if strings.Contains(r, "round 3") {
		return "R3"
	}
	if strings.Contains(r, "qualifying round 1") {
		return "Q1"
	}
	if strings.Contains(r, "qualifying round 2") {
		return "Q2"
	}
	if len(r) > 5 {
		return r[:5]
	}
	return r
}

func fetchTennisScoreboard(tour string, dateStr string) (*TennisScoreboard, error) {
	espnDate := strings.ReplaceAll(dateStr, "-", "")
	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/tennis/%s/scoreboard?dates=%s", tour, espnDate)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESPN API returned status code %d", resp.StatusCode)
	}

	var sched TennisScoreboard
	if err := json.NewDecoder(resp.Body).Decode(&sched); err != nil {
		return nil, err
	}
	return &sched, nil
}

func renderTennisSchedule(tournaments []MergedTournament, dateStr string, format string, loc *time.Location) string {
	var sb strings.Builder

	zoneName, _ := time.Now().In(loc).Zone()
	title := fmt.Sprintf("TENNIS LIVE SCOREBOARD (%s %s)", dateStr, zoneName)
	padding := (layoutWidth - len(title)) / 2
	if padding < 0 {
		padding = 0
	}

	currentDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		currentDate = time.Now().In(loc)
	}
	prevDateStr := currentDate.AddDate(0, 0, -1).Format("2006-01-02")
	nextDateStr := currentDate.AddDate(0, 0, 1).Format("2006-01-02")

	// Sports Selector row
	var banner strings.Builder
	if format != "html" {
		banner.WriteString(style("==============================================================================\n", ansiCyan, format))
		banner.WriteString(txt("           ", format) + style("[MLB]", ansiGray, format) + txt("             ", format) + style("[BASKETBALL]", ansiGray, format) + txt("             ", format) + style("[TENNIS]", ansiBold+ansiGreen, format) + "\n")
		banner.WriteString(style("==============================================================================\n", ansiCyan, format))
	}
	banner.WriteString(txt(strings.Repeat(" ", padding), format))
	banner.WriteString(style(title+"\n", ansiBold+ansiCyan, format))

	// Date Navigation Row
	banner.WriteString(style("==============================================================================\n", ansiCyan, format))
	banner.WriteString(txt(" ", format))
	prevLinkText := fmt.Sprintf("<< PREV DAY (%s)", prevDateStr)
	nextLinkText := fmt.Sprintf("NEXT DAY (%s) >>", nextDateStr)
	spacerSize := layoutWidth - 1 - len(prevLinkText) - len(nextLinkText)
	if spacerSize < 1 {
		spacerSize = 1
	}
	if format == "html" {
		prevLink := fmt.Sprintf(`<a href="/tennis?date=%s" class="term-link">%s</a>`, prevDateStr, prevLinkText)
		nextLink := fmt.Sprintf(`<a href="/tennis?date=%s" class="term-link">%s</a>`, nextDateStr, nextLinkText)
		banner.WriteString(prevLink + strings.Repeat(" ", spacerSize) + nextLink + "\n")
	} else {
		banner.WriteString(style(prevLinkText, ansiGreen, format) + strings.Repeat(" ", spacerSize) + style(nextLinkText, ansiGreen, format) + "\n")
	}
	banner.WriteString(style("==============================================================================\n", ansiCyan, format))
	sb.WriteString(termPre(format, banner.String()))

	if len(tournaments) == 0 {
		sb.WriteString(txt(" No matches scheduled for this date.\n", format))
		sb.WriteString(style("==============================================================================\n", ansiCyan, format))
		return termPre(format, sb.String())
	}

	// buildSetCell returns a cell whose HTML carries per-set winner coloring.
	buildSetCell := func(comp TennisCompetitor) TableCell {
		var plain, htmlB strings.Builder
		for sIdx := 0; sIdx < 5; sIdx++ {
			var valStr string
			win, tb := false, false
			if sIdx < len(comp.Linescores) {
				valStr = fmt.Sprintf("%.0f", comp.Linescores[sIdx].Value)
				win = comp.Linescores[sIdx].Winner
				tb = comp.Linescores[sIdx].Tiebreak > 0
			} else {
				valStr = "-"
			}
			tbChar := ""
			if tb {
				tbChar = "ₜ"
			}
			plain.WriteString(valStr + tbChar + " ")
			if format == "html" {
				cls := ""
				if win {
					cls = "term-bold term-red"
				}
				if cls != "" {
					htmlB.WriteString(`<span class="` + cls + `">` + valStr + tbChar + `</span> `)
				} else {
					htmlB.WriteString(`<span>` + valStr + tbChar + `</span> `)
				}
			}
		}
		if format == "html" {
			return TableCell{Text: strings.TrimRight(plain.String(), " "), HTML: strings.TrimRight(htmlB.String(), " ")}
		}
		return TableCell{Text: strings.TrimRight(plain.String(), " ")}
	}

	playerCell := func(comp TennisCompetitor, isLive, isFinal bool) TableCell {
		name := getCompetitorName(comp)
		if utf8.RuneCountInString(name) > 24 {
			name = string([]rune(name)[:23]) + "."
		}
		if comp.Possession {
			name += "ₛ"
		}
		var code string
		switch {
		case isLive:
			code = ansiGreen
			if comp.Winner {
				code = ansiBold + ansiGreen
			}
		case isFinal:
			if comp.Winner {
				code = ansiBold
			}
		default:
			code = ansiGray
		}
		return TableCell{Text: name, ANSI: code, Class: classFor(code)}
	}

	addMatchRow := func(t *Table, match TennisCompetition, dateStr string) {
		if len(match.Competitors) < 2 {
			return
		}
		var homeComp, awayComp TennisCompetitor
		for _, c := range match.Competitors {
			if c.HomeAway == "home" {
				homeComp = c
			} else {
				awayComp = c
			}
		}

		state := match.Status.Type.State
		isLive := state == "in"
		isFinal := state == "post"

		var idStyle, roundStyle string
		if isLive {
			idStyle = ansiGreen
			roundStyle = ansiGreen
		} else {
			idStyle = ansiGray
			roundStyle = ansiGray
		}

		statusStr := strings.TrimSuffix(match.Status.Type.Detail, " Set")
		if utf8.RuneCountInString(statusStr) > 12 {
			statusStr = string([]rune(statusStr)[:11]) + "."
		}

		t.AddRow(
			TableCell{Text: match.ID, Link: "/tennis/game/" + match.ID + "?date=" + dateStr},
			TableCell{Text: abbreviateRound(match.Round.DisplayName), ANSI: roundStyle, Class: classFor(roundStyle)},
			playerCell(awayComp, isLive, isFinal),
			buildSetCell(awayComp),
			TableCell{Text: statusStr, ANSI: idStyle, Class: classFor(idStyle)},
		)
		t.AddRow(
			TableCell{Text: ""},
			TableCell{Text: ""},
			playerCell(homeComp, isLive, isFinal),
			buildSetCell(homeComp),
			TableCell{Text: ""},
		)
	}

	renderTable := func(matches []TennisCompetition, header string, headerStyle string) {
		if len(matches) == 0 {
			return
		}
		var h strings.Builder
		h.WriteString(style(fmt.Sprintf("\n %s\n", header), ansiBold+headerStyle, format))
		sb.WriteString(termPre(format, h.String()))

		t := NewTable(format,
			TableCol{Title: "ID", Width: 9},
			TableCol{Title: "RND", Width: 5},
			TableCol{Title: "PLAYER", Width: 24},
			TableCol{Title: "SETS", Width: 14},
			TableCol{Title: "STATUS", Width: 12},
		)
		for _, m := range matches {
			addMatchRow(t, m, dateStr)
		}
		sb.WriteString(t.Render())
	}

	// First pass: collect all live matches (flat, across tournaments) and non-live per tournament
	type liveEntry struct {
		match TennisCompetition
		label string
	}
	var allLive []liveEntry

	type nonLiveGroup struct {
		tournament MergedTournament
		matches    []TennisCompetition
	}
	var allNonLive []nonLiveGroup

	for _, mt := range tournaments {
		var matches []TennisCompetition
		for _, grp := range mt.Event.Groupings {
			matches = append(matches, grp.Competitions...)
		}
		if len(matches) == 0 {
			continue
		}

		var nonLive []TennisCompetition
		for _, m := range matches {
			if m.Status.Type.State == "in" {
				allLive = append(allLive, liveEntry{
					match: m,
					label: fmt.Sprintf("%s (%s)", mt.Event.Name, mt.Tour),
				})
			} else {
				nonLive = append(nonLive, m)
			}
		}
		if len(nonLive) > 0 {
			sort.SliceStable(nonLive, func(i, j int) bool {
				t1, _ := time.Parse(time.RFC3339, nonLive[i].Date)
				t2, _ := time.Parse(time.RFC3339, nonLive[j].Date)
				return t1.Before(t2)
			})
			allNonLive = append(allNonLive, nonLiveGroup{mt, nonLive})
		}
	}

	// Sort all live entries by time
	sort.SliceStable(allLive, func(i, j int) bool {
		t1, _ := time.Parse(time.RFC3339, allLive[i].match.Date)
		t2, _ := time.Parse(time.RFC3339, allLive[j].match.Date)
		return t1.Before(t2)
	})

	// Render live matches section at top (flat across all tournaments)
	if len(allLive) > 0 {
		var live []TennisCompetition
		for _, e := range allLive {
			live = append(live, e.match)
		}
		renderTable(live, ">> LIVE MATCHES", ansiGreen)
	}

	// Render per-tournament non-live sections
	for _, ng := range allNonLive {
		renderTable(ng.matches, fmt.Sprintf("TOURNAMENT: %s (%s)", ng.tournament.Event.Name, ng.tournament.Tour), ansiCyan)
	}

	var footer strings.Builder
	if format == "ansi" {
		footer.WriteString(txt(" Run 'curl http://localhost:9090/tennis/game/<ID>' to view a match in real-time.\n", format))
	} else {
		footer.WriteString(txt(" Click on a game ID to view the match in real-time.\n", format))
	}
	footer.WriteString(style("==============================================================================\n", ansiCyan, format))
	sb.WriteString(termPre(format, footer.String()))

	return sb.String()
}

func findTennisCompetition(dateStr string, gamePk string) (*TennisCompetition, *TennisEvent, string, error) {
	atpSched, atpErr := fetchTennisScoreboard("atp", dateStr)
	if atpErr == nil && atpSched != nil {
		for _, ev := range atpSched.Events {
			for _, grp := range ev.Groupings {
				for _, comp := range grp.Competitions {
					if comp.ID == gamePk {
						return &comp, &ev, "ATP", nil
					}
				}
			}
		}
	}

	wtaSched, wtaErr := fetchTennisScoreboard("wta", dateStr)
	if wtaErr == nil && wtaSched != nil {
		for _, ev := range wtaSched.Events {
			for _, grp := range ev.Groupings {
				for _, comp := range grp.Competitions {
					if comp.ID == gamePk {
						return &comp, &ev, "WTA", nil
					}
				}
			}
		}
	}

	return nil, nil, "", fmt.Errorf("match %s not found on date %s", gamePk, dateStr)
}

// ---------------------------------------------------------------------------
// TennisLive.net — point-by-point live scores
//
// ESPN's tennis feed only exposes set-level data (games won per set). To show
// the live point score within the current game (0/15/30/40, server, break
// points) we supplement it with TennisLive.net, which server-renders the full
// point-by-point history on each match detail page. This is best-effort: the
// ESPN set view is always shown even if this lookup fails or is unavailable.
// ---------------------------------------------------------------------------

type tennisLivePoint struct {
	Found        bool
	CurrentGame  string // e.g. "4-5"
	CurrentPoint string // e.g. "40-40", "A-40"
	Server       string
	BreakPoint   bool
}

type tennisLiveDetailCacheEntry struct {
	html      string
	timestamp time.Time
}

var (
	tennisLiveListingMu    sync.RWMutex
	tennisLiveListingData  map[string]string
	tennisLiveListingStamp time.Time

	tennisLiveDetailMu    sync.Mutex
	tennisLiveDetailCache = map[string]tennisLiveDetailCacheEntry{}
)

// tennisLiveSurnames returns the surname token(s) used for fuzzy matching.
// Only the last word of each name (or slash-separated doubles partner) is
// used so that "J. Kym" matches "Jerome Kym".
func tennisLiveSurnames(name string) []string {
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '/' || r == ' ' || r == '.'
	})
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.ToLower(p)
		var b strings.Builder
		for _, r := range p {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
				b.WriteRune(r)
			}
		}
		if b.Len() > 0 {
			out = append(out, b.String())
		}
	}
	if len(out) == 0 {
		return out
	}
	return []string{out[len(out)-1]}
}

// tennisLiveMatchKey builds an order-independent key from the two player names
// plus the tour, so we can locate a match in the TennisLive listing.
func tennisLiveMatchKey(tour, p1, p2 string) string {
	s := append(tennisLiveSurnames(p1), tennisLiveSurnames(p2)...)
	sort.Strings(s)
	return strings.ToLower(tour) + ":" + strings.Join(s, "|")
}

// tennisLiveFetch fetches a URL, optionally caching the body for cacheFor.
func tennisLiveFetch(url string, cacheFor time.Duration) (string, error) {
	if cacheFor > 0 {
		tennisLiveDetailMu.Lock()
		if e, ok := tennisLiveDetailCache[url]; ok && time.Since(e.timestamp) < cacheFor {
			h := e.html
			tennisLiveDetailMu.Unlock()
			return h, nil
		}
		tennisLiveDetailMu.Unlock()
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0 Safari/537.36")
	req.Header.Set("Cookie", "verified=1")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	h := string(body)
	if cacheFor > 0 {
		tennisLiveDetailMu.Lock()
		tennisLiveDetailCache[url] = tennisLiveDetailCacheEntry{html: h, timestamp: time.Now()}
		tennisLiveDetailMu.Unlock()
	}
	return h, nil
}

// fetchTennisLiveListing parses the TennisLive homepage (which lists all live
// matches with player names and detail links) into a matchKey -> detailURL map.
func fetchTennisLiveListing() (map[string]string, error) {
	tennisLiveListingMu.RLock()
	if tennisLiveListingData != nil && time.Since(tennisLiveListingStamp) < 60*time.Second {
		d := tennisLiveListingData
		tennisLiveListingMu.RUnlock()
		return d, nil
	}
	tennisLiveListingMu.RUnlock()

	htmlStr, err := tennisLiveFetch("https://www.tennislive.net/", 0)
	if err != nil {
		return nil, err
	}

	linkRe := regexp.MustCompile(`<a href="(https://www\.tennislive\.net/(atp|wta)/match/[^"]+)"[^>]*class="nu-reward"`)
	matchRe := regexp.MustCompile(`class="match "><a href="https://www\.tennislive\.net/(?:atp|wta)/[^"]+"[^>]*>([^<]+)</a>`)

	data := map[string]string{}
	for _, l := range linkRe.FindAllStringSubmatch(htmlStr, -1) {
		url := l[1]
		tour := l[2]
		pos := strings.Index(htmlStr, url)
		if pos < 0 {
			continue
		}
		before := matchRe.FindAllStringSubmatch(htmlStr[:pos], -1)
		p1 := ""
		if len(before) > 0 {
			p1 = html.UnescapeString(strings.TrimSpace(before[len(before)-1][1]))
		}
		after := matchRe.FindStringSubmatch(htmlStr[pos:])
		p2 := ""
		if after != nil {
			p2 = html.UnescapeString(strings.TrimSpace(after[1]))
		}
		if p1 == "" || p2 == "" {
			continue
		}
		data[tennisLiveMatchKey(tour, p1, p2)] = url
	}

	tennisLiveListingMu.Lock()
	tennisLiveListingData = data
	tennisLiveListingStamp = time.Now()
	tennisLiveListingMu.Unlock()
	return data, nil
}

// parseTennisLivePoint extracts the current (most recent) game's point score,
// server and break-point flag from a TennisLive match detail page. It never
// panics: any unexpected input yields a zero (Found=false) value so a bad
// scrape can never crash the server.
func parseTennisLivePoint(htmlStr string) (res tennisLivePoint) {
	defer func() {
		if r := recover(); r != nil {
			res = tennisLivePoint{}
		}
	}()

	re := regexp.MustCompile(`class="mp_info_txt">([^<]*)</td>.*?class="mp_15">([^<]*)</td>`)
	locs := re.FindAllStringSubmatchIndex(htmlStr, -1)
	if len(locs) == 0 {
		return tennisLivePoint{}
	}
	last := locs[len(locs)-1]
	gameScore := strings.TrimSpace(html.UnescapeString(htmlStr[last[2]:last[3]]))
	prog := strings.TrimSpace(html.UnescapeString(htmlStr[last[4]:last[5]]))
	tokens := strings.Split(prog, ",")
	if len(tokens) == 0 {
		return tennisLivePoint{}
	}
	cur := strings.TrimSpace(tokens[len(tokens)-1])
	bp := strings.Contains(cur, "[BP]")
	cur = strings.ReplaceAll(cur, "[BP]", "")
	cur = strings.TrimSpace(cur)

	infoStart := last[2]
	trStart := strings.LastIndex(htmlStr[:infoStart], "<tr")
	rowEnd := strings.Index(htmlStr[infoStart:], "</tr>")
	server := ""
	if trStart >= 0 && rowEnd >= 0 && trStart < infoStart+rowEnd {
		row := htmlStr[trStart : infoStart+rowEnd]
		serveRe := regexp.MustCompile(`class="mp_serve">([^<]*?)(<img[^>]*>)?</td>`)
		for _, s := range serveRe.FindAllStringSubmatch(row, -1) {
			name := strings.TrimSpace(html.UnescapeString(s[1]))
			if s[2] != "" {
				server = name
			}
		}
	}

	return tennisLivePoint{
		Found:        true,
		CurrentGame:  gameScore,
		CurrentPoint: cur,
		Server:       server,
		BreakPoint:   bp,
	}
}

// getTennisLivePoint looks up a match by player names and returns its live
// point score, or a zero (Found=false) value if it cannot be resolved. It never
// panics, so a failed TennisLive lookup degrades gracefully to the ESPN view.
func getTennisLivePoint(tour, p1, p2 string) (res tennisLivePoint) {
	defer func() {
		if r := recover(); r != nil {
			res = tennisLivePoint{}
		}
	}()
	listing, err := fetchTennisLiveListing()
	if err != nil {
		return tennisLivePoint{}
	}
	url, ok := listing[tennisLiveMatchKey(tour, p1, p2)]
	if !ok {
		return tennisLivePoint{}
	}
	h, err := tennisLiveFetch(url, 12*time.Second)
	if err != nil {
		return tennisLivePoint{}
	}
	return parseTennisLivePoint(h)
}

func renderTennisGame(comp TennisCompetition, event TennisEvent, tour string, dateStr string, format string) string {
	var sb strings.Builder

	if len(comp.Competitors) < 2 {
		return "Error: Incomplete competitor data for this match."
	}

	var homeComp, awayComp TennisCompetitor
	for _, c := range comp.Competitors {
		if c.HomeAway == "home" {
			homeComp = c
		} else {
			awayComp = c
		}
	}

	awayName := getCompetitorFullName(awayComp)
	homeName := getCompetitorFullName(homeComp)

	state := comp.Status.Type.State
	detail := comp.Status.Type.Detail

	badge := fmt.Sprintf("[%s]", strings.ToUpper(state))
	var badgeColor string
	switch state {
	case "in":
		badgeColor = ansiGreen
	case "post":
		badgeColor = ansiBold
	default:
		badgeColor = ansiGray
	}

	badgeStyled := style(badge, badgeColor, format)

	titleLine := fmt.Sprintf(" %s  %s vs %s (%s)\n",
		badgeStyled,
		style(awayName, ansiBold, format),
		style(homeName, ansiBold, format),
		tour,
	)

	subTitleLine := fmt.Sprintf(" Tournament: %s\n Venue: %s - %s\n", event.Name, comp.Venue.FullName, comp.Venue.Court)

	servingPlayer := ""
	if awayComp.Possession {
		servingPlayer = awayName
	} else if homeComp.Possession {
		servingPlayer = homeName
	}

	maxSets := 5
	if len(awayComp.Linescores) > maxSets {
		maxSets = len(awayComp.Linescores)
	}
	if len(homeComp.Linescores) > maxSets {
		maxSets = len(homeComp.Linescores)
	}
	if maxSets < 3 {
		maxSets = 3
	}

	var header strings.Builder
	header.WriteString(style("========================================================================\n", ansiCyan, format))
	header.WriteString(titleLine)
	header.WriteString(style(subTitleLine, ansiGray, format))
	header.WriteString(style("========================================================================\n", ansiCyan, format))
	header.WriteString("\n")
	header.WriteString(" Status: " + style(detail, ansiYellow, format) + "\n")
	if servingPlayer != "" {
		header.WriteString(" " + style("* SERVING: "+servingPlayer, ansiYellow, format) + "\n")
	}
	sb.WriteString(termPre(format, header.String()))

	if len(comp.Notes) > 0 {
		var notes strings.Builder
		for _, note := range comp.Notes {
			notes.WriteString(" " + style(">> "+note.Text, ansiBold+ansiGreen, format) + "\n")
		}
		sb.WriteString(termPre(format, notes.String()))
	}

	isLive := state == "in"
	isFinal := state == "post"

	playerCell := func(comp TennisCompetitor) TableCell {
		name := getCompetitorName(comp)
		if utf8.RuneCountInString(name) > 24 {
			name = string([]rune(name)[:23]) + "."
		}
		if comp.Possession {
			name += "ₛ"
		}
		var code string
		switch {
		case isLive:
			code = ansiGreen
			if comp.Winner {
				code = ansiBold + ansiGreen
			}
		case isFinal:
			if comp.Winner {
				code = ansiBold
			}
		default:
			code = ansiGray
		}
		return TableCell{Text: name, ANSI: code, Class: classFor(code)}
	}

	setCell := func(comp TennisCompetitor, sIdx int) TableCell {
		var valStr string
		win, tb := false, false
		if sIdx < len(comp.Linescores) {
			valStr = fmt.Sprintf("%.0f", comp.Linescores[sIdx].Value)
			win = comp.Linescores[sIdx].Winner
			tb = comp.Linescores[sIdx].Tiebreak > 0
		} else {
			valStr = "--"
		}
		tbChar := ""
		if tb {
			tbChar = "ₜ"
		}
		if format == "html" {
			cls := ""
			if win {
				cls = "term-bold term-red"
			}
			return TableCell{Text: valStr + tbChar, HTML: `<span class="` + cls + `">` + valStr + tbChar + `</span>`, Align: alignRight}
		}
		code := ""
		if win {
			code = ansiBold + ansiRed
		}
		return TableCell{Text: valStr + tbChar, ANSI: code, Align: alignRight}
	}

	awaySetsWon := 0
	homeSetsWon := 0
	for _, ls := range awayComp.Linescores {
		if ls.Winner {
			awaySetsWon++
		}
	}
	for _, ls := range homeComp.Linescores {
		if ls.Winner {
			homeSetsWon++
		}
	}

	cols := []TableCol{{Title: "PLAYER", Width: 26}}
	for i := 1; i <= maxSets; i++ {
		cols = append(cols, TableCol{Title: fmt.Sprintf("S%d", i), Width: 3, Align: alignRight})
	}
	cols = append(cols, TableCol{Title: "SETS", Width: 3, Align: alignRight})

	st := NewTable(format, cols...)

	awayRow := []TableCell{playerCell(awayComp)}
	for sIdx := 0; sIdx < maxSets; sIdx++ {
		awayRow = append(awayRow, setCell(awayComp, sIdx))
	}
	awayRow = append(awayRow, TableCell{Text: fmt.Sprintf("%d", awaySetsWon), Align: alignRight})
	st.AddRow(awayRow...)

	homeRow := []TableCell{playerCell(homeComp)}
	for sIdx := 0; sIdx < maxSets; sIdx++ {
		homeRow = append(homeRow, setCell(homeComp, sIdx))
	}
	homeRow = append(homeRow, TableCell{Text: fmt.Sprintf("%d", homeSetsWon), Align: alignRight})
	st.AddRow(homeRow...)

	sb.WriteString(st.Render())
	sb.WriteString("\n")

	// Live point-by-point score (supplements ESPN's set-only data)
	if state == "in" {
		pt := getTennisLivePoint(tour, awayName, homeName)
		if pt.Found {
			pointDisplay := pt.CurrentPoint
			switch pt.CurrentPoint {
			case "40-40":
				pointDisplay = "40-40 (Deuce)"
			case "A-40":
				pointDisplay = "Ad (Server)"
			case "40-A":
				pointDisplay = "Ad (Returner)"
			}
			bp := ""
			if pt.BreakPoint {
				bp = "  [BREAK POINT]"
			}

			var b strings.Builder
			b.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
			b.WriteString(style(" LIVE GAME SCORE  (point-by-point · TennisLive)\n", ansiBold+ansiYellow, format))
			b.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
			b.WriteString(" " + style("CURRENT GAME", ansiYellow, format) + ": " + style(pt.CurrentGame, ansiGreen, format) + "\n")
			b.WriteString(" " + style("POINTS", ansiYellow, format) + ": " + style(pointDisplay+bp, ansiGreen, format) + "\n")
			if pt.Server != "" {
				b.WriteString(" " + style("SERVER", ansiYellow, format) + ": " + style("★ "+pt.Server, ansiGreen, format) + "\n")
			}
			b.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
			sb.WriteString(termPre(format, b.String()))
		}
	}

	if format == "ansi" {
		sb.WriteString(txt(fmt.Sprintf(" Run 'curl http://localhost:9090/tennis?date=%s' to return to the scoreboard.\n", dateStr), format))
		sb.WriteString(style("========================================================================\n", ansiCyan, format))
	}

	return termPre(format, sb.String())
}

func handleTennisSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if serveHTMLWrapper(w, r) {
		return
	}

	tzStr := r.URL.Query().Get("tz")
	if tzStr == "" {
		tzStr = "America/Los_Angeles"
	}
	loc, err := time.LoadLocation(tzStr)
	if err != nil {
		loc, _ = time.LoadLocation("America/Los_Angeles")
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = time.Now().In(loc).Format("2006-01-02")
	}

	var wg sync.WaitGroup
	var atpSched, wtaSched *TennisScoreboard
	var atpErr, wtaErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		atpSched, atpErr = fetchTennisScoreboard("atp", dateStr)
	}()
	go func() {
		defer wg.Done()
		wtaSched, wtaErr = fetchTennisScoreboard("wta", dateStr)
	}()
	wg.Wait()

	if atpErr != nil && wtaErr != nil {
		http.Error(w, "Failed to connect to ESPN Tennis API: "+atpErr.Error(), http.StatusBadGateway)
		return
	}

	var tournaments []MergedTournament
	if atpErr == nil && atpSched != nil {
		for _, ev := range atpSched.Events {
			tournaments = append(tournaments, MergedTournament{Tour: "ATP", Event: ev})
		}
	}
	if wtaErr == nil && wtaSched != nil {
		for _, ev := range wtaSched.Events {
			tournaments = append(tournaments, MergedTournament{Tour: "WTA", Event: ev})
		}
	}

	format := getFormat(r)
	text := renderTennisSchedule(tournaments, dateStr, format, loc)
	writeResponse(w, format, text)
}

func handleTennisGame(w http.ResponseWriter, r *http.Request) {
	gamePk := r.PathValue("gamePk")
	if gamePk == "" {
		http.Error(w, "Missing game ID", http.StatusBadRequest)
		return
	}

	if serveHTMLWrapper(w, r) {
		return
	}

	tzStr := r.URL.Query().Get("tz")
	if tzStr == "" {
		tzStr = "America/Los_Angeles"
	}
	loc, err := time.LoadLocation(tzStr)
	if err != nil {
		loc, _ = time.LoadLocation("America/Los_Angeles")
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = time.Now().In(loc).Format("2006-01-02")
	}

	comp, event, tour, err := findTennisCompetition(dateStr, gamePk)
	format := getFormat(r)

	if err != nil {
		var sb strings.Builder
		sb.WriteString(style("========================================================================\n", ansiRed, format))
		sb.WriteString(style(fmt.Sprintf("                       ERROR: GAME NOT FOUND (%s)\n", gamePk), ansiBold+ansiRed, format))
		sb.WriteString(style("========================================================================\n", ansiRed, format))
		sb.WriteString(txt(fmt.Sprintf(" Details: %s\n", err.Error()), format))
		sb.WriteString(style("========================================================================\n", ansiRed, format))
		if format == "ansi" {
			sb.WriteString(txt(fmt.Sprintf(" Run 'curl http://localhost:9090/tennis?date=%s' to return to the scoreboard.\n", dateStr), format))
			sb.WriteString(style("========================================================================\n", ansiRed, format))
		}

		writeResponse(w, format, termPre(format, sb.String()))
		return
	}

	text := renderTennisGame(*comp, *event, tour, dateStr, format)
	writeResponse(w, format, text)
}

func handleAPITennisGames(w http.ResponseWriter, r *http.Request) {
	tzStr := r.URL.Query().Get("tz")
	if tzStr == "" {
		tzStr = "America/Los_Angeles"
	}
	loc, err := time.LoadLocation(tzStr)
	if err != nil {
		loc, _ = time.LoadLocation("America/Los_Angeles")
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = time.Now().In(loc).Format("2006-01-02")
	}

	var wg sync.WaitGroup
	var atpSched, wtaSched *TennisScoreboard
	var atpErr, wtaErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		atpSched, atpErr = fetchTennisScoreboard("atp", dateStr)
	}()
	go func() {
		defer wg.Done()
		wtaSched, wtaErr = fetchTennisScoreboard("wta", dateStr)
	}()
	wg.Wait()

	if atpErr != nil && wtaErr != nil {
		http.Error(w, "Failed to fetch tennis data", http.StatusBadGateway)
		return
	}

	type CombinedResponse struct {
		ATP *TennisScoreboard `json:"atp,omitempty"`
		WTA *TennisScoreboard `json:"wta,omitempty"`
	}

	res := CombinedResponse{
		ATP: atpSched,
		WTA: wtaSched,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func handleAPITennisGameDetail(w http.ResponseWriter, r *http.Request) {
	gamePk := r.PathValue("gamePk")
	if gamePk == "" {
		http.Error(w, "Missing game ID", http.StatusBadRequest)
		return
	}

	tzStr := r.URL.Query().Get("tz")
	if tzStr == "" {
		tzStr = "America/Los_Angeles"
	}
	loc, err := time.LoadLocation(tzStr)
	if err != nil {
		loc, _ = time.LoadLocation("America/Los_Angeles")
	}

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = time.Now().In(loc).Format("2006-01-02")
	}

	comp, _, _, err := findTennisCompetition(dateStr, gamePk)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(comp)
}
