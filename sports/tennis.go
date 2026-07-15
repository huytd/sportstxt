package sports

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// TennisScoreboard represents the response from ESPN tennis scoreboard API

func padRight(s string, n int) string {
	r := []rune(s)
	if len(r) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(r))
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
		Value     float64 `json:"value"`
		Winner    bool    `json:"winner"`
		Tiebreak  float64 `json:"tiebreak,omitempty"`
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
	if format != "html" {
		sb.WriteString(style("==============================================================================\n", ansiCyan, format))
		sb.WriteString(txt("           ", format) + style("[MLB]", ansiGray, format) + txt("             ", format) + style("[NBA]", ansiGray, format) + txt("             ", format) + style("[TENNIS]", ansiBold+ansiGreen, format) + "\n")
		sb.WriteString(style("==============================================================================\n", ansiCyan, format))
	}
	sb.WriteString(txt(strings.Repeat(" ", padding), format))
	sb.WriteString(style(title+"\n", ansiBold+ansiCyan, format))

	// Date Navigation Row
	sb.WriteString(style("==============================================================================\n", ansiCyan, format))
	sb.WriteString(txt(" ", format))
	prevLinkText := fmt.Sprintf("<< PREV DAY (%s)", prevDateStr)
	nextLinkText := fmt.Sprintf("NEXT DAY (%s) >>", nextDateStr)
	spacerSize := layoutWidth - 1 - len(prevLinkText) - len(nextLinkText)
	if spacerSize < 1 {
		spacerSize = 1
	}
	if format == "html" {
		prevLink := fmt.Sprintf(`<a href="/tennis?date=%s" class="term-link">%s</a>`, prevDateStr, prevLinkText)
		nextLink := fmt.Sprintf(`<a href="/tennis?date=%s" class="term-link">%s</a>`, nextDateStr, nextLinkText)
		sb.WriteString(prevLink + strings.Repeat(" ", spacerSize) + nextLink + "\n")
	} else {
		sb.WriteString(style(prevLinkText, ansiGreen, format) + strings.Repeat(" ", spacerSize) + style(nextLinkText, ansiGreen, format) + "\n")
	}
	sb.WriteString(style("==============================================================================\n", ansiCyan, format))

	if len(tournaments) == 0 {
		sb.WriteString(txt(" No matches scheduled for this date.\n", format))
		sb.WriteString(style("==============================================================================\n", ansiCyan, format))
		return sb.String()
	}

	renderMatchRow := func(match TennisCompetition) {
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

		awayName := getCompetitorName(awayComp)
		homeName := getCompetitorName(homeComp)
		if len(awayName) > 23 {
			awayName = awayName[:22] + "."
		}
		if len(homeName) > 23 {
			homeName = homeName[:22] + "."
		}

		awayServes := ""
		homeServes := ""
		if awayComp.Possession {
			awayServes = "ₛ"
		}
		if homeComp.Possession {
			homeServes = "ₛ"
		}

		awayNameWithServe := awayName + awayServes
		homeNameWithServe := homeName + homeServes

		state := match.Status.Type.State
		isLive := state == "in"
		isFinal := state == "post"

		awaySets := ""
		homeSets := ""

		for sIdx := 0; sIdx < 5; sIdx++ {
			var aScore, hScore string
			aScoreVal := -1.0
			hScoreVal := -1.0
			aWin := false
			hWin := false
			var aTB, hTB float64

			if sIdx < len(awayComp.Linescores) {
				aScoreVal = awayComp.Linescores[sIdx].Value
				aWin = awayComp.Linescores[sIdx].Winner
				aTB = awayComp.Linescores[sIdx].Tiebreak
			}
			if sIdx < len(homeComp.Linescores) {
				hScoreVal = homeComp.Linescores[sIdx].Value
				hWin = homeComp.Linescores[sIdx].Winner
				hTB = homeComp.Linescores[sIdx].Tiebreak
			}

			if aScoreVal >= 0 {
				valStr := fmt.Sprintf("%.0f", aScoreVal)
				switch {
				case aWin && aTB > 0:
					aScore = style(fmt.Sprintf(" %sₜ ", valStr), ansiBold+ansiRed, format)
				case aWin:
					aScore = style(fmt.Sprintf(" %s  ", valStr), ansiBold+ansiRed, format)
				case aTB > 0:
					aScore = style(fmt.Sprintf(" %sₜ ", valStr), "", format)
				default:
					aScore = style(fmt.Sprintf(" %s  ", valStr), "", format)
				}
			} else {
				aScore = "    "
			}

			if hScoreVal >= 0 {
				valStr := fmt.Sprintf("%.0f", hScoreVal)
				switch {
				case hWin && hTB > 0:
					hScore = style(fmt.Sprintf(" %sₜ ", valStr), ansiBold+ansiRed, format)
				case hWin:
					hScore = style(fmt.Sprintf(" %s  ", valStr), ansiBold+ansiRed, format)
				case hTB > 0:
					hScore = style(fmt.Sprintf(" %sₜ ", valStr), "", format)
				default:
					hScore = style(fmt.Sprintf(" %s  ", valStr), "", format)
				}
			} else {
				hScore = "    "
			}

			awaySets += aScore
			homeSets += hScore
		}

		setColWidth := 20
		awaySets = padRight(awaySets, setColWidth)
		homeSets = padRight(homeSets, setColWidth)

		statusStr := strings.TrimSuffix(match.Status.Type.Detail, " Set")
		if len(statusStr) > 12 {
			statusStr = statusStr[:11] + "."
		}

		var idStyle, roundStyle, awayPlayerStyle, homePlayerStyle string
		if isLive {
			idStyle = ansiGreen
			roundStyle = ansiGreen
			awayPlayerStyle = ansiGreen
			homePlayerStyle = ansiGreen
			if awayComp.Winner {
				awayPlayerStyle = ansiBold + ansiGreen
			} else if homeComp.Winner {
				homePlayerStyle = ansiBold + ansiGreen
			}
		} else if isFinal {
			idStyle = ansiGray
			roundStyle = ansiGray
			awayPlayerStyle = ""
			homePlayerStyle = ""
			if awayComp.Winner {
				awayPlayerStyle = ansiBold
			} else if homeComp.Winner {
				homePlayerStyle = ansiBold
			}
		} else {
			idStyle = ansiGray
			roundStyle = ansiGray
			awayPlayerStyle = ansiGray
			homePlayerStyle = ansiGray
		}

		idVal := match.ID

		var idPart string
		if format == "html" {
			idPart = fmt.Sprintf(`<a href="/tennis/game/%s?date=%s" style="text-decoration:none">%s</a>`, idVal, dateStr, style(padRight(idVal, 7), idStyle, format))
		} else {
			idPart = style(padRight(idVal, 7), idStyle, format)
		}
		roundPart := style(padRight(abbreviateRound(match.Round.DisplayName), 5), roundStyle, format)
		awayPlayerPart := style(padRight(awayNameWithServe, 25), awayPlayerStyle, format)
		homePlayerPart := style(padRight(homeNameWithServe, 25), homePlayerStyle, format)
		statusPart := style(padRight(statusStr, 12), idStyle, format)

		sb.WriteString(fmt.Sprintf(" %s %s %s %s %s\n", idPart, roundPart, awayPlayerPart, awaySets, statusPart))
		sb.WriteString(fmt.Sprintf(" %s %s %s\n", strings.Repeat(" ", 13), homePlayerPart, homeSets))
		sb.WriteString(txt("\n", format))
	}

	renderSection := func(matches []TennisCompetition, header string, headerStyle string) {
		if len(matches) == 0 {
			return
		}
		sb.WriteString(style(" -----------------------------------------------------------------------------\n", ansiCyan, format))
		sb.WriteString(style(fmt.Sprintf(" %-7s %-5s %-25s %-20s %-12s\n", "ID", "ROUND", "PLAYERS", "SETS", "STATUS"), ansiBold, format))
		sb.WriteString(style(" -----------------------------------------------------------------------------\n", ansiCyan, format))
		sb.WriteString(style(fmt.Sprintf(" %s\n", header), ansiBold+headerStyle, format))
		for _, m := range matches {
			renderMatchRow(m)
		}
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
		sb.WriteString(style("\n >> LIVE MATCHES\n", ansiBold+ansiGreen, format))
		sb.WriteString(style(" -----------------------------------------------------------------------------\n", ansiCyan, format))
		sb.WriteString(style(fmt.Sprintf(" %-7s %-5s %-25s %-20s %-12s\n", "ID", "ROUND", "PLAYERS", "SETS", "STATUS"), ansiBold, format))
		sb.WriteString(style(" -----------------------------------------------------------------------------\n", ansiCyan, format))
		for _, e := range allLive {
			renderMatchRow(e.match)
		}
	}

	// Render per-tournament non-live sections
	for _, ng := range allNonLive {
		sb.WriteString(style(fmt.Sprintf("\n TOURNAMENT: %s (%s)\n", ng.tournament.Event.Name, ng.tournament.Tour), ansiBold+ansiCyan, format))
		renderSection(ng.matches, ">> FINAL / UPCOMING", ansiGray)
	}

	if format == "ansi" {
		sb.WriteString(txt(" Run 'curl http://localhost:9090/tennis/game/<ID>' to view a match in real-time.\n", format))
	} else {
		sb.WriteString(txt(" Click on a game ID to view the match in real-time.\n", format))
	}
	sb.WriteString(style("==============================================================================\n", ansiCyan, format))

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

	sb.WriteString(style("========================================================================\n", ansiCyan, format))
	sb.WriteString(titleLine)
	sb.WriteString(style(subTitleLine, ansiGray, format))
	sb.WriteString(style("========================================================================\n", ansiCyan, format))

	sb.WriteString("\n")
	sb.WriteString(" Status: " + style(detail, ansiYellow, format) + "\n")

	servingPlayer := ""
	if awayComp.Possession {
		servingPlayer = awayName
	} else if homeComp.Possession {
		servingPlayer = homeName
	}
	if servingPlayer != "" {
		sb.WriteString(" " + style("* SERVING: "+servingPlayer, ansiYellow, format) + "\n")
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

	sb.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
	sb.WriteString(style(" PLAYERS                    ", ansiBold, format))
	for i := 1; i <= maxSets; i++ {
		sb.WriteString(style(fmt.Sprintf("S%-4d", i), ansiBold, format))
	}
	sb.WriteString(style("| SETS\n", ansiBold, format))
	if len(comp.Notes) > 0 {
		for _, note := range comp.Notes {
			sb.WriteString(" " + style(">> "+note.Text, ansiBold+ansiGreen, format) + "\n")
		}
	}
	sb.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))

	awayNameShort := getCompetitorName(awayComp)
	homeNameShort := getCompetitorName(homeComp)
	if len(awayNameShort) > 25 {
		awayNameShort = awayNameShort[:24] + "."
	}
	if len(homeNameShort) > 25 {
		homeNameShort = homeNameShort[:24] + "."
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

	sb.WriteString(txt(" "+padRight(awayNameShort, 26)+" ", format))
	for sIdx := 0; sIdx < maxSets; sIdx++ {
		val := " --  "
		if sIdx < len(awayComp.Linescores) {
			scoreVal := awayComp.Linescores[sIdx].Value
			win := awayComp.Linescores[sIdx].Winner
			tb := awayComp.Linescores[sIdx].Tiebreak
			valStr := fmt.Sprintf("%.0f", scoreVal)
			switch {
			case win && tb > 0:
				val = style(fmt.Sprintf(" %sₜ  ", valStr), ansiBold+ansiRed, format)
			case win:
				val = style(fmt.Sprintf(" %s   ", valStr), ansiBold+ansiRed, format)
			case tb > 0:
				val = style(fmt.Sprintf(" %sₜ  ", valStr), "", format)
			default:
				val = style(fmt.Sprintf(" %s   ", valStr), "", format)
			}
		}
		sb.WriteString(val)
	}
	sb.WriteString(txt(fmt.Sprintf("|  %d\n", awaySetsWon), format))

	sb.WriteString(txt(" "+padRight(homeNameShort, 26)+" ", format))
	for sIdx := 0; sIdx < maxSets; sIdx++ {
		val := " --  "
		if sIdx < len(homeComp.Linescores) {
			scoreVal := homeComp.Linescores[sIdx].Value
			win := homeComp.Linescores[sIdx].Winner
			tb := homeComp.Linescores[sIdx].Tiebreak
			valStr := fmt.Sprintf("%.0f", scoreVal)
			switch {
			case win && tb > 0:
				val = style(fmt.Sprintf(" %sₜ  ", valStr), ansiBold+ansiRed, format)
			case win:
				val = style(fmt.Sprintf(" %s   ", valStr), ansiBold+ansiRed, format)
			case tb > 0:
				val = style(fmt.Sprintf(" %sₜ  ", valStr), "", format)
			default:
				val = style(fmt.Sprintf(" %s   ", valStr), "", format)
			}
		}
		sb.WriteString(val)
	}
	sb.WriteString(txt(fmt.Sprintf("|  %d\n", homeSetsWon), format))

	sb.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
	sb.WriteString("\n")
	sb.WriteString(style("========================================================================\n", ansiCyan, format))

	if format == "ansi" {
		sb.WriteString(txt(fmt.Sprintf(" Run 'curl http://localhost:9090/tennis?date=%s' to return to the scoreboard.\n", dateStr), format))
		sb.WriteString(style("========================================================================\n", ansiCyan, format))
	}

	return sb.String()
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

		writeResponse(w, format, sb.String())
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
