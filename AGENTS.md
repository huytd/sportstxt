# sportstxt - Agent Knowledge Base

## Project Overview
A Go-based scoreboard application that displays MLB/NBA game data in both terminal (ANSI) and web (HTML) formats. Serves as a live feed with auto-refresh for active games.

## Directory Structure
```
/home/king/code/sportstxt/
├── main.go              # Entry point, HTTP server setup
├── api/
│   └── index.go         # API routes
├── sports/
│   ├── sports.go        # Shared handlers, HTML template, common functions
│   ├── mlb.go           # MLB-specific rendering and handlers
│   └── nba.go           # NBA-specific rendering and handlers
├── sportstxt            # Compiled binary
└── vercel.json          # Deployment config
```

## Key Files & Their Purposes

### sports/sports.go
- **HTML Template (`htmlPage`)**: Main web UI with navigation, status bar, terminal content area
- **Shared Handlers**: `handleSchedule`, `handleStandings`, `handleTeamPage`, `handleCompareTeams`
- **Common Functions**:
  - `style(text, styleCode, format)` - Applies ANSI or HTML styling
  - `txt(text, format)` - Escapes text for HTML mode
  - `getFormat(r)` - Returns "html" or "ansi" based on request
  - `serveHTMLWrapper(w, r)` - Serves HTML shell for browser requests
  - `writeResponse(w, format, text)` - Sets content type and writes output
- **ANSI Color Codes**: `ansiReset`, `ansiBold`, `ansiGreen`, `ansiYellow`, `ansiRed`, `ansiCyan`, `ansiBlue`, `ansiMagenta`, `ansiGray`

### sports/mlb.go
- **MLB Game Rendering**: `renderGame()` - Main game detail page
- **MLB Handlers**: `handleGame()`, `handleSchedule()`
- **Helper Functions**:
  - `renderDiamondAndMatchup()` - Pitcher/batter display
  - `renderCurrentPitches()` - Ball/strike/out dots
  - `renderBoxscore()` - Player stats
  - `renderStandings()` - League standings table
  - `renderTeamPage()` - Individual team page
  - `renderCompareTeams()` - Team comparison view

### sports/nba.go
- **NBA Game Rendering**: Similar structure to MLB but for basketball
- **Handlers**: `handleNBASchedule()`, `handleNBAGame()`

## Routes
| Route | Handler | Description |
|-------|---------|-------------|
| `/` | `handleSchedule` | Main scoreboard (today's games) |
| `/game/{gamePk}` | `handleGame` | MLB game detail page |
| `/mlb/standings` | `handleStandings` | MLB standings table |
| `/mlb/team/{teamId}` | `handleTeamPage` | Individual team info |
| `/mlb/compare` | `handleCompareTeams` | Compare two teams (query: `?team1=X&team2=Y`) |
| `/nba` | `handleNBASchedule` | NBA schedule |
| `/nba/game/{gamePk}` | `handleNBAGame` | NBA game detail |

## API Endpoints
- `GET /api/games` - MLB games JSON
- `GET /api/game/{gamePk}` - Single MLB game detail
- `GET /api/nba/games` - NBA games JSON
- `GET /api/nba/game/{gamePk}` - Single NBA game detail

## External APIs
- **MLB Stats API**: `https://statsapi.mlb.com/api/v1/`
  - `/teams?sportId=1` - All teams
  - `/schedule` - Game schedule
  - `/game/{gamePk}/feed` - Live game data
  - `/standings` - League standings

## Rendering Format Logic
- **HTML Mode**: Browser requests (no `?raw=1`, not curl) → HTML with styled spans/links
- **ANSI Mode**: Curl requests or `?raw=1` param → ANSI escape codes
- Links in HTML mode use `<a href="..." class="term-link">...</a>`

## Common Patterns for Modifications

### Adding a Link to Game Page
Modify `renderGame()` in `sports/mlb.go`:
```go
// Around line 840, in the HTML format block
link := fmt.Sprintf(`<a href="/path?param=value" class="term-link">LINK TEXT</a>`)
line := fmt.Sprintf(`<span class="term-gray"> ... %s ...</span>\n`, link)
```

### Adding Navigation Item
Edit `htmlPage` constant in `sports/sports.go`:
```go
<nav style="display: flex; gap: 12px; margin-bottom: 15px;">
    <a href="/path" class="term-link" style="padding: 6px 14px; border-radius: 4px;">Label</a>
</nav>
```

### Adding a New Route
1. Add handler function (e.g., `func handleNewRoute(w http.ResponseWriter, r *http.Request)`)
2. Register in `NewHandler()`: `mux.HandleFunc("GET /path", handleNewRoute)`
3. Optionally add navigation link

## Styling Classes (HTML Mode)
- `.term-green`, `.term-yellow`, `.term-red`, `.term-cyan`, `.term-blue`, `.term-magenta` - Colors
- `.term-gray` - Gray text
- `.term-bold` - Bold text
- `.term-link` - Clickable links with hover effect
- `.term-highlight-better` - Highlighted stat (better value)

## Team IDs (Common)
| Team | ID | Abbreviation |
|------|-----|--------------|
| LA Dodgers | 119 | LAD |
| NY Mets | 121 | NYM |
| Chicago Cubs | 112 | CHC |
| St. Louis Cardinals | 138 | STL |
| San Francisco Giants | 137 | SF |

## Quick Reference Commands
```bash
# Run locally
go run .

# Build binary
go build -o sportstxt .

# View routes
grep "mux.HandleFunc" sports/sports.go

# Find function definition
grep -n "func renderGame" sports/*.go
```

## Key Data Structures
- `GameFeedResponse` - Full game data (MLB)
- `TeamInfo` - Team details (id, name, abbreviation, city, league, division)
- `LeagueStandingsResponse` - Standings with records array

## Auto-Polling Behavior
- Scoreboard (`/`) and game pages (`/game/{id}`): Poll every 10 seconds
- Static pages (team info, compare, standings): No polling
- Controlled by `updatePolling()` in HTML template
