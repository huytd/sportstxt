package sports

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"
)

// ANSI terminal style escape codes
const (
	layoutWidth = 78
	ansiReset   = "\033[0m"
	ansiBold   = "\033[1m"
	ansiGreen  = "\033[32m"
	ansiYellow = "\033[33m"
	ansiRed    = "\033[31m"
	ansiCyan   = "\033[36m"
	ansiBlue   = "\033[34m"
	ansiMagenta = "\033[35m"
	ansiGray   = "\033[90m"
)

var client = &http.Client{Timeout: 10 * time.Second}

// style colorizes text for ANSI terminal output or wraps in styled HTML tags
func style(text string, styleCode string, format string) string {
	if format == "ansi" {
		return styleCode + text + ansiReset
	} else if format == "html" {
		var class string
		switch styleCode {
		case ansiGreen:
			class = "term-green"
		case ansiYellow:
			class = "term-yellow"
		case ansiRed:
			class = "term-red"
		case ansiCyan:
			class = "term-cyan"
		case ansiBlue:
			class = "term-blue"
		case ansiMagenta:
			class = "term-magenta"
		case ansiGray:
			class = "term-gray"
		case ansiBold:
			class = "term-bold"
		case ansiBold + ansiCyan:
			return fmt.Sprintf(`<span class="term-bold term-cyan">%s</span>`, html.EscapeString(text))
		case ansiBold + ansiRed:
			return fmt.Sprintf(`<span class="term-bold term-red">%s</span>`, html.EscapeString(text))
		case ansiBold + ansiGreen:
			return fmt.Sprintf(`<span class="term-bold term-green">%s</span>`, html.EscapeString(text))
		case ansiBold + ansiYellow:
			return fmt.Sprintf(`<span class="term-bold term-yellow">%s</span>`, html.EscapeString(text))
		case ansiBold + ansiMagenta:
			return fmt.Sprintf(`<span class="term-bold term-magenta">%s</span>`, html.EscapeString(text))
		default:
			return html.EscapeString(text)
		}
		return fmt.Sprintf(`<span class="%s">%s</span>`, class, html.EscapeString(text))
	}
	return text
}

// txt handles simple text escaping in HTML mode
func txt(text string, format string) string {
	if format == "html" {
		return html.EscapeString(text)
	}
	return text
}

// dots creates dot-indicators for balls/strikes/outs
func dots(count, max int) string {
	var parts []string
	for i := 0; i < max; i++ {
		if i < count {
			parts = append(parts, "●")
		} else {
			parts = append(parts, "○")
		}
	}
	return strings.Join(parts, " ")
}

// stripANSIAndHTML strips ANSI color escape sequences and HTML tags to calculate raw text length
func stripANSIAndHTML(s string) string {
	var sb strings.Builder
	inANSI := false
	inHTML := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if inANSI {
			if ch == 'm' {
				inANSI = false
			}
			continue
		}
		if inHTML {
			if ch == '>' {
				inHTML = false
			}
			continue
		}
		if ch == '\033' {
			inANSI = true
			continue
		}
		if ch == '<' {
			inHTML = true
			continue
		}
		sb.WriteByte(ch)
	}
	return sb.String()
}

// isBrowserRequest returns true if the request is from a web browser (not curl and not requesting raw content).
func isBrowserRequest(r *http.Request) bool {
	isRaw := r.URL.Query().Get("raw") == "1"
	isCurl := strings.Contains(strings.ToLower(r.UserAgent()), "curl")
	return !isRaw && !isCurl
}

// serveHTMLWrapper serves the main htmlPage shell if it's a browser request, and returns true.
func serveHTMLWrapper(w http.ResponseWriter, r *http.Request) bool {
	if isBrowserRequest(r) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlPage))
		return true
	}
	return false
}

// getFormat returns the rendering format ("html" or "ansi") based on user agent and raw query param.
func getFormat(r *http.Request) string {
	isRaw := r.URL.Query().Get("raw") == "1"
	isCurl := strings.Contains(strings.ToLower(r.UserAgent()), "curl")
	if isCurl && !isRaw {
		return "ansi"
	}
	return "html"
}

// writeResponse sets the appropriate content type and writes the output text.
func writeResponse(w http.ResponseWriter, format string, text string) {
	if format == "ansi" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	w.Write([]byte(text))
}

// NewHandler constructs and returns the HTTP handler with all registered routes.
func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleSchedule)
	mux.HandleFunc("GET /game/{gamePk}", handleGame)
	mux.HandleFunc("GET /api/games", handleAPIGames)
	mux.HandleFunc("GET /api/game/{gamePk}", handleAPIGameDetail)

	// MLB Standings and Team Routes
	mux.HandleFunc("GET /mlb/standings", handleStandings)
	mux.HandleFunc("GET /mlb/team/{teamId}", handleTeamPage)
	mux.HandleFunc("GET /mlb/compare", handleCompareTeams)

	// NBA Routes
	mux.HandleFunc("GET /nba", handleNBASchedule)
	mux.HandleFunc("GET /nba/game/{gamePk}", handleNBAGame)
	mux.HandleFunc("GET /api/nba/games", handleAPINBAGames)
	mux.HandleFunc("GET /api/nba/game/{gamePk}", handleAPINBAGameDetail)

	return mux
}

// handleStandings handles the /mlb/standings route
func handleStandings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if serveHTMLWrapper(w, r) {
		return
	}

	// Fetch all team info to get abbreviations and city names
	teamsUrl := "https://statsapi.mlb.com/api/v1/teams?sportId=1"
	teamsResp, err := client.Get(teamsUrl)
	if err != nil {
		http.Error(w, "Failed to connect to MLB Stats API: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer teamsResp.Body.Close()

	var allTeams struct {
		Teams []TeamInfo `json:"teams"`
	}
	if err := json.NewDecoder(teamsResp.Body).Decode(&allTeams); err != nil {
		http.Error(w, "Failed to decode teams JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Build team lookup map
	teamMap := make(map[int]TeamInfo)
	for _, t := range allTeams.Teams {
		teamMap[t.Id] = t
	}

	seasonYear := time.Now().Year()

	// Fetch standings for AL and NL separately using league-based endpoints
	alUrl := fmt.Sprintf("https://statsapi.mlb.com/api/v1/standings?sportId=1&leagueId=103&season=%d&seasonStage=regularSeason&granularity=team", seasonYear)
	nlUrl := fmt.Sprintf("https://statsapi.mlb.com/api/v1/standings?sportId=1&leagueId=104&season=%d&seasonStage=regularSeason&granularity=team", seasonYear)

	alResp, err := client.Get(alUrl)
	if err != nil {
		http.Error(w, "Failed to connect to MLB Stats API: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer alResp.Body.Close()

	nlResp, err := client.Get(nlUrl)
	if err != nil {
		http.Error(w, "Failed to connect to MLB Stats API: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer nlResp.Body.Close()

	if alResp.StatusCode != http.StatusOK || nlResp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("MLB Stats API returned status code %d", alResp.StatusCode), http.StatusBadGateway)
		return
	}

	var alStandings, nlStandings LeagueStandingsResponse
	if err := json.NewDecoder(alResp.Body).Decode(&alStandings); err != nil {
		http.Error(w, "Failed to decode AL standings JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewDecoder(nlResp.Body).Decode(&nlStandings); err != nil {
		http.Error(w, "Failed to decode NL standings JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Combine AL and NL standings
	combined := LeagueStandingsResponse{}
	combined.Records = append(combined.Records, alStandings.Records...)
	combined.Records = append(combined.Records, nlStandings.Records...)

	format := getFormat(r)
	text := renderStandings(combined, teamMap, format)
	writeResponse(w, format, text)
}

// handleTeamPage handles the /mlb/team/{teamId} route
func handleTeamPage(w http.ResponseWriter, r *http.Request) {
	teamId := r.PathValue("teamId")
	teamIdInt, err := strconv.Atoi(teamId)
	if err != nil || teamIdInt <= 0 {
		http.Error(w, "Invalid team ID", http.StatusBadRequest)
		return
	}

	if serveHTMLWrapper(w, r) {
		return
	}

	// Fetch team info first
	infoUrl := fmt.Sprintf("https://statsapi.mlb.com/api/v1/teams/%d?sportId=1", teamIdInt)
	resp, err := client.Get(infoUrl)
	if err != nil {
		http.Error(w, "Failed to connect to MLB Stats API: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	var teamInfo struct {
		Teams []TeamInfo `json:"teams"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&teamInfo); err != nil {
		http.Error(w, "Failed to decode team info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if len(teamInfo.Teams) == 0 {
		http.Error(w, "Team not found", http.StatusNotFound)
		return
	}

	team := teamInfo.Teams[0]

	// Determine league display name
	leagueName := "AL"
	if strings.Contains(strings.ToLower(team.League.Name), "national") {
		leagueName = "NL"
	}
	divisionName := team.Division.Name
	if divisionName == "" {
		divisionName = "Unknown"
	}
	cityName := team.LocationName
	if cityName == "" {
		cityName = "Unknown"
	}

	format := getFormat(r)
	text := renderTeamPage(team.Id, team.Name, team.Abbreviation, cityName, leagueName, divisionName, format)
	writeResponse(w, format, text)
}

// handleCompareTeams handles the /mlb/compare route
func handleCompareTeams(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if serveHTMLWrapper(w, r) {
		return
	}

	// Fetch all team info to get list of teams
	teamsUrl := "https://statsapi.mlb.com/api/v1/teams?sportId=1"
	teamsResp, err := client.Get(teamsUrl)
	if err != nil {
		http.Error(w, "Failed to connect to MLB Stats API: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer teamsResp.Body.Close()

	var allTeams struct {
		Teams []TeamInfo `json:"teams"`
	}
	if err := json.NewDecoder(teamsResp.Body).Decode(&allTeams); err != nil {
		http.Error(w, "Failed to decode teams JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Parse query parameters
	team1 := r.URL.Query().Get("team1")
	team2 := r.URL.Query().Get("team2")

	team1Id, _ := strconv.Atoi(team1)
	team2Id, _ := strconv.Atoi(team2)

	// In raw HTML or ANSI modes, generate comparison output
	format := getFormat(r)
	text := renderCompareTeams(team1Id, team2Id, allTeams.Teams, format)
	writeResponse(w, format, text)
}

//go:embed index.html
var htmlPage string

