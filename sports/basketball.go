package sports

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// basketballLeague identifies a basketball tour supported by ESPN.
type basketballLeague struct {
	Slug string // ESPN sport slug, e.g. "nba", "wnba"
	Name string // Display label, e.g. "NBA"
}

// basketballLeagues lists every basketball tour we aggregate.
var basketballLeagues = []basketballLeague{
	{"nba", "NBA"},
	{"wnba", "WNBA"},
	{"mens-college-basketball", "NCAA M"},
	{"womens-college-basketball", "NCAA W"},
}

// basketballScheduleEntry is one league's scoreboard for the combined view.
type basketballScheduleEntry struct {
	League basketballLeague
	Sched  ESPNBAScoreboard
}

// ESPNBAEvent represents a single game within a basketball scoreboard.
type ESPNBAEvent struct {
	ID        string `json:"id"`
	Date      string `json:"date"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Status    struct {
		Clock        float64 `json:"clock"`
		DisplayClock string  `json:"displayClock"`
		Period       int     `json:"period"`
		Type         struct {
			State       string `json:"state"`
			Detail      string `json:"detail"`
			Description string `json:"description"`
		} `json:"type"`
	} `json:"status"`
	Competitions []struct {
		Competitors []struct {
			HomeAway string `json:"homeAway"`
			Score    string `json:"score"`
			Team     struct {
				ID           string `json:"id"`
				DisplayName  string `json:"displayName"`
				Abbreviation string `json:"abbreviation"`
			} `json:"team"`
			Linescores []struct {
				DisplayValue string `json:"displayValue"`
			} `json:"linescores"`
		} `json:"competitors"`
	} `json:"competitions"`
}

// ESPNBAScoreboard represents the response from ESPN basketball scoreboard API
type ESPNBAScoreboard struct {
	Events []ESPNBAEvent `json:"events"`
}

// ESPNBAPlayerStats represents individual player statistics in boxscore
type ESPNBAPlayerStats struct {
	Active     bool `json:"active"`
	Starter    bool `json:"starter"`
	DidNotPlay bool `json:"didNotPlay"`
	Reason     string `json:"reason"`
	Athlete    struct {
		ID          string `json:"id"`
		DisplayName string `json:"displayName"`
		ShortName   string `json:"shortName"`
		Jersey      string `json:"jersey"`
		Position    struct {
			Abbreviation string `json:"abbreviation"`
		} `json:"position"`
	} `json:"athlete"`
	Stats []string `json:"stats"`
}

// ESPNBASeries represents series info
type ESPNBASeries struct {
	Summary string `json:"summary"`
}

// ESPNBABoxscoreTeamStats represents team stats in the boxscore
type ESPNBABoxscoreTeamStats struct {
	Team struct {
		ID           string `json:"id"`
		DisplayName  string `json:"displayName"`
		Abbreviation string `json:"abbreviation"`
	} `json:"team"`
	Statistics []struct {
		Name         string `json:"name"`
		DisplayValue string `json:"displayValue"`
		Label        string `json:"label"`
	} `json:"statistics"`
	HomeAway string `json:"homeAway"`
}

// ESPNBABoxscoreTeam represents team data in the boxscore
type ESPNBABoxscoreTeam struct {
	Team struct {
		ID           string `json:"id"`
		DisplayName  string `json:"displayName"`
		Abbreviation string `json:"abbreviation"`
	} `json:"team"`
	Statistics []struct {
		Names   []string            `json:"names"`
		Keys    []string            `json:"keys"`
		Athletes []ESPNBAPlayerStats `json:"athletes"`
		Totals  []string            `json:"totals"`
	} `json:"statistics"`
}

// ESPNBAGameSummary represents the response from ESPN NBA summary API
type ESPNBAGameSummary struct {
	Header struct {
		ID           string `json:"id"`
		Competitions []struct {
			Status struct {
				Type struct {
					State       string `json:"state"`
					Detail      string `json:"detail"`
					Description string `json:"description"`
				} `json:"type"`
			} `json:"status"`
			Competitors []struct {
				HomeAway   string `json:"homeAway"`
				Score      string `json:"score"`
				Possession bool   `json:"possession"`
				Team       struct {
					ID           string `json:"id"`
					DisplayName  string `json:"displayName"`
					Abbreviation string `json:"abbreviation"`
				} `json:"team"`
				Linescores []struct {
					DisplayValue string `json:"displayValue"`
				} `json:"linescores"`
			} `json:"competitors"`
			Series []ESPNBASeries `json:"series"`
		} `json:"competitions"`
	} `json:"header"`
	Leaders []struct {
		Team struct {
			ID           string `json:"id"`
			Abbreviation string `json:"abbreviation"`
		} `json:"team"`
		Leaders []struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
			Leaders     []struct {
				Value        float64 `json:"value"`
				DisplayValue string  `json:"displayValue"`
				Athlete      struct {
					DisplayName string `json:"displayName"`
				} `json:"athlete"`
			} `json:"leaders"`
		} `json:"leaders"`
	} `json:"leaders"`
	Plays []struct {
		Text       string `json:"text"`
		AwayScore  int    `json:"awayScore"`
		HomeScore  int    `json:"homeScore"`
		Period     struct {
			Number int `json:"number"`
		} `json:"period"`
		Clock struct {
			DisplayValue string `json:"displayValue"`
		} `json:"clock"`
	} `json:"plays"`
	Boxscore struct {
		Teams   []ESPNBABoxscoreTeamStats `json:"teams"`
		Players []ESPNBABoxscoreTeam      `json:"players"`
	} `json:"boxscore"`
}

// renderBasketballSchedule creates the plain-text scoreboard view for one or
// more basketball leagues. When more than one entry is supplied the view is a
// combined "all basketball" board with a LEAGUE column.
func renderBasketballSchedule(entries []basketballScheduleEntry, dateStr string, format string, loc *time.Location) string {
	var sb strings.Builder

	combined := len(entries) > 1

	zoneName, _ := time.Now().In(loc).Zone()
	var title string
	if combined {
		title = fmt.Sprintf("BASKETBALL LIVE SCOREBOARD (%s %s)", dateStr, zoneName)
	} else {
		title = fmt.Sprintf("%s LIVE SCOREBOARD (%s %s)", entries[0].League.Name, dateStr, zoneName)
	}
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

	var basePath string
	if combined {
		basePath = "/basketball"
	} else {
		basePath = "/nba"
	}

	var banner strings.Builder
	if format != "html" {
		banner.WriteString(style("==============================================================================\n", ansiCyan, format))
		if combined {
			banner.WriteString(txt("           ", format) + style("[MLB]", ansiGray, format) + txt("             ", format) + style("[BASKETBALL]", ansiBold+ansiGreen, format) + txt("             ", format) + style("[TENNIS]", ansiGray, format) + "\n")
		} else {
			banner.WriteString(txt("           ", format) + style("[MLB]", ansiGray, format) + txt("             ", format) + style("[NBA]", ansiBold+ansiGreen, format) + txt("             ", format) + style("[TENNIS]", ansiGray, format) + "\n")
		}
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
		prevLink := fmt.Sprintf(`<a href="%s?date=%s" class="term-link">%s</a>`, basePath, prevDateStr, prevLinkText)
		nextLink := fmt.Sprintf(`<a href="%s?date=%s" class="term-link">%s</a>`, basePath, nextDateStr, nextLinkText)
		banner.WriteString(prevLink + strings.Repeat(" ", spacerSize) + nextLink + "\n")
	} else {
		banner.WriteString(style(prevLinkText, ansiGreen, format) + strings.Repeat(" ", spacerSize) + style(nextLinkText, ansiGreen, format) + "\n")
	}
	banner.WriteString(style("==============================================================================\n", ansiCyan, format))
	sb.WriteString(termPre(format, banner.String()))

	totalEvents := 0
	for _, e := range entries {
		totalEvents += len(e.Sched.Events)
	}
	if totalEvents == 0 {
		sb.WriteString(termPre(format, txt(" No games scheduled for this date.\n", format)+
			style("==============================================================================\n", ansiCyan, format)))
		return sb.String()
	}

	// Flatten events across leagues for sorting/rendering.
	var events []struct {
		league basketballLeague
		event  ESPNBAEvent
	}
	for _, e := range entries {
		for i := range e.Sched.Events {
			events = append(events, struct {
				league basketballLeague
				event  ESPNBAEvent
			}{e.League, e.Sched.Events[i]})
		}
	}

	sort.SliceStable(events, func(i, j int) bool {
		parseTime := func(s string) time.Time {
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				return t
			}
			if t, err := time.Parse("2006-01-02T15:04Z", s); err == nil {
				return t
			}
			return time.Time{}
		}
		t1 := parseTime(events[i].event.Date)
		t2 := parseTime(events[j].event.Date)
		if !t1.IsZero() && !t2.IsZero() {
			return t1.Before(t2)
		}
		return events[i].event.Date < events[j].event.Date
	})

	cols := []TableCol{}
	if combined {
		cols = append(cols, TableCol{Title: "LEAGUE", Width: 9})
	}
	cols = append(cols,
		TableCol{Title: "ID", Width: 9},
		TableCol{Title: "TIME", Width: 8},
		TableCol{Title: "TEAM", Width: 18},
		TableCol{Title: "PTS", Width: 3, Align: alignRight},
		TableCol{Title: "STATUS", Width: 10},
	)
	t := NewTable(format, cols...)

	classFor := func(code string) string {
		switch code {
		case ansiGreen:
			return "term-green"
		case ansiGray:
			return "term-gray"
		case ansiBold:
			return "term-bold"
		case ansiBold + ansiGreen:
			return "term-bold term-green"
		}
		return ""
	}

	for _, ev := range events {
		event := ev.event
		idStr := event.ID
		var homeComp, awayComp struct {
			HomeAway string `json:"homeAway"`
			Score    string `json:"score"`
			Team     struct {
				ID           string `json:"id"`
				DisplayName  string `json:"displayName"`
				Abbreviation string `json:"abbreviation"`
			} `json:"team"`
			Linescores []struct {
				DisplayValue string `json:"displayValue"`
			} `json:"linescores"`
		}

		if len(event.Competitions) > 0 {
			for _, competitor := range event.Competitions[0].Competitors {
				if competitor.HomeAway == "home" {
					homeComp = competitor
				} else {
					awayComp = competitor
				}
			}
		}

		awayName := awayComp.Team.DisplayName
		homeName := homeComp.Team.DisplayName

		if utf8.RuneCountInString(awayName) > 17 {
			awayName = string([]rune(awayName)[:16]) + "."
		}
		if utf8.RuneCountInString(homeName) > 17 {
			homeName = string([]rune(homeName)[:16]) + "."
		}

		state := event.Status.Type.State
		isLive := state == "in"
		isFinal := state == "post"

		awayScoreStr := "-"
		homeScoreStr := "-"

		if isLive || isFinal {
			awayScoreStr = awayComp.Score
			homeScoreStr = homeComp.Score
		}

		statusStr := event.Status.Type.Detail
		if utf8.RuneCountInString(statusStr) > 10 {
			statusStr = string([]rune(statusStr)[:9]) + "."
		}

		var awayStyle, homeStyle string
		var baseStyle string

		if isLive {
			baseStyle = ansiGreen
			awayStyle = ansiGreen
			homeStyle = ansiGreen
			awayScoreVal, _ := strconv.Atoi(awayComp.Score)
			homeScoreVal, _ := strconv.Atoi(homeComp.Score)
			if awayScoreVal > homeScoreVal {
				awayStyle = ansiBold + ansiGreen
			} else if homeScoreVal > awayScoreVal {
				homeStyle = ansiBold + ansiGreen
			}
		} else if isFinal {
			baseStyle = ""
			awayStyle = ""
			homeStyle = ""
			awayScoreVal, _ := strconv.Atoi(awayComp.Score)
			homeScoreVal, _ := strconv.Atoi(homeComp.Score)
			if awayScoreVal > homeScoreVal {
				awayStyle = ansiBold
			} else if homeScoreVal > awayScoreVal {
				homeStyle = ansiBold
			}
		} else {
			baseStyle = ansiGray
			awayStyle = ansiGray
			homeStyle = ansiGray
		}

		gameTime := "--:--"
		if t, err := time.Parse(time.RFC3339, event.Date); err == nil {
			gameTime = t.In(loc).Format("03:04 PM")
		} else if t, err := time.Parse("2006-01-02T15:04Z", event.Date); err == nil {
			gameTime = t.In(loc).Format("03:04 PM")
		}

		gameLink := "/nba/game/" + idStr
		if combined {
			gameLink = "/basketball/game/" + idStr + "?league=" + ev.league.Slug
		}

		awayRow := []TableCell{}
		homeRow := []TableCell{}
		if combined {
			awayRow = append(awayRow, TableCell{Text: ev.league.Name, Class: classFor(baseStyle), ANSI: baseStyle})
			homeRow = append(homeRow, TableCell{Text: ""})
		}
		awayRow = append(awayRow,
			TableCell{Text: idStr, Link: gameLink},
			TableCell{Text: gameTime, Class: classFor(baseStyle), ANSI: baseStyle},
			TableCell{Text: awayName, Class: classFor(awayStyle), ANSI: awayStyle},
			TableCell{Text: awayScoreStr, Align: alignRight, Class: classFor(awayStyle), ANSI: awayStyle},
			TableCell{Text: statusStr, Class: classFor(baseStyle), ANSI: baseStyle},
		)
		homeRow = append(homeRow,
			TableCell{Text: ""},
			TableCell{Text: ""},
			TableCell{Text: homeName, Class: classFor(homeStyle), ANSI: homeStyle},
			TableCell{Text: homeScoreStr, Align: alignRight, Class: classFor(homeStyle), ANSI: homeStyle},
			TableCell{Text: ""},
		)
		t.AddRow(awayRow...)
		t.AddRow(homeRow...)
	}

	sb.WriteString(t.Render())

	var footer strings.Builder
	footer.WriteString(style("------------------------------------------------------------------------------\n", ansiCyan, format))
	if format == "ansi" {
		linkPath := "/nba/game/<ID>"
		if combined {
			linkPath = "/basketball/game/<ID>?league=<LEAGUE>"
		}
		footer.WriteString(txt(fmt.Sprintf(" Run 'curl ...%s' for live game detail.\n", linkPath), format))
	} else {
		footer.WriteString(txt(" Click on a game ID to view the game in real-time.\n", format))
	}
	footer.WriteString(style("==============================================================================\n", ansiCyan, format))
	sb.WriteString(termPre(format, footer.String()))

	return sb.String()
}

// renderNBABoxscore generates a team's boxscore table
func renderNBABoxscore(teamEntry ESPNBABoxscoreTeam, format string) string {
	var sb strings.Builder

	sb.WriteString(termPre(format, style(fmt.Sprintf("\n %s STATISTICS\n", strings.ToUpper(teamEntry.Team.DisplayName)), ansiBold+ansiCyan, format)))

	if len(teamEntry.Statistics) == 0 {
		sb.WriteString(termPre(format, txt(" No statistics available.\n", format)))
		return sb.String()
	}

	statIndexes := make(map[string]int)
	for idx, name := range teamEntry.Statistics[0].Names {
		statIndexes[name] = idx
	}

	getStat := func(stats []string, name string) string {
		idx, ok := statIndexes[name]
		if !ok || idx >= len(stats) {
			return "-"
		}
		return stats[idx]
	}

	t := NewTable(format,
		TableCol{Title: "PLAYER", Width: 24},
		TableCol{Title: "MIN", Width: 3, Align: alignRight},
		TableCol{Title: "PTS", Width: 3, Align: alignRight},
		TableCol{Title: "FG", Width: 5, Align: alignRight},
		TableCol{Title: "3PT", Width: 5, Align: alignRight},
		TableCol{Title: "FT", Width: 5, Align: alignRight},
		TableCol{Title: "REB", Width: 3, Align: alignRight},
		TableCol{Title: "AST", Width: 3, Align: alignRight},
		TableCol{Title: "TO", Width: 2, Align: alignRight},
		TableCol{Title: "STL", Width: 3, Align: alignRight},
		TableCol{Title: "BLK", Width: 3, Align: alignRight},
		TableCol{Title: "+/-", Width: 4, Align: alignRight},
	)

	for _, athlete := range teamEntry.Statistics[0].Athletes {
		name := athlete.Athlete.ShortName
		pos := athlete.Athlete.Position.Abbreviation
		displayName := name
		if pos != "" {
			displayName = name + " " + pos
		}
		if utf8.RuneCountInString(displayName) > 24 {
			displayName = string([]rune(displayName)[:23]) + "."
		}

		if athlete.DidNotPlay {
			reason := athlete.Reason
			if reason == "" {
				reason = "Did Not Play"
			}
			sb.WriteString(termPre(format, style(fmt.Sprintf(" %-24s DNP - %s\n", displayName, reason), ansiGray, format)))
			continue
		}

		min := getStat(athlete.Stats, "MIN")
		pts := getStat(athlete.Stats, "PTS")
		fg := getStat(athlete.Stats, "FG")
		tpt := getStat(athlete.Stats, "3PT")
		ft := getStat(athlete.Stats, "FT")
		reb := getStat(athlete.Stats, "REB")
		ast := getStat(athlete.Stats, "AST")
		to := getStat(athlete.Stats, "TO")
		stl := getStat(athlete.Stats, "STL")
		blk := getStat(athlete.Stats, "BLK")
		pm := getStat(athlete.Stats, "+/-")

		t.AddRow(
			TableCell{Text: displayName},
			TableCell{Text: min, Align: alignRight},
			TableCell{Text: pts, Align: alignRight},
			TableCell{Text: fg, Align: alignRight},
			TableCell{Text: tpt, Align: alignRight},
			TableCell{Text: ft, Align: alignRight},
			TableCell{Text: reb, Align: alignRight},
			TableCell{Text: ast, Align: alignRight},
			TableCell{Text: to, Align: alignRight},
			TableCell{Text: stl, Align: alignRight},
			TableCell{Text: blk, Align: alignRight},
			TableCell{Text: pm, Align: alignRight},
		)
	}

	sb.WriteString(t.Render())

	totals := teamEntry.Statistics[0].Totals
	minTot := getStat(totals, "MIN")
	ptsTot := getStat(totals, "PTS")
	fgTot := getStat(totals, "FG")
	tptTot := getStat(totals, "3PT")
	ftTot := getStat(totals, "FT")
	rebTot := getStat(totals, "REB")
	astTot := getStat(totals, "AST")
	toTot := getStat(totals, "TO")
	stlTot := getStat(totals, "STL")
	blkTot := getStat(totals, "BLK")
	pmTot := getStat(totals, "+/-")

	totalsRow := fmt.Sprintf(" %-24s %3s %3s %5s %5s %5s %3s %3s %2s %3s %3s %4s\n",
		"TOTALS", minTot, ptsTot, fgTot, tptTot, ftTot, rebTot, astTot, toTot, stlTot, blkTot, pmTot,
	)
	sb.WriteString(termPre(format, style(totalsRow, ansiBold, format)))

	return sb.String()
}

// renderNBAGame creates the detailed view of a game
func renderNBAGame(summary ESPNBAGameSummary, leagueLabel string, format string) string {
	var sb strings.Builder

	if len(summary.Header.Competitions) == 0 {
		return "Error: No competition details found for this game."
	}

	comp := summary.Header.Competitions[0]

	var homeComp, awayComp struct {
		HomeAway   string `json:"homeAway"`
		Score      string `json:"score"`
		Possession bool   `json:"possession"`
		Team       struct {
			ID           string `json:"id"`
			DisplayName  string `json:"displayName"`
			Abbreviation string `json:"abbreviation"`
		} `json:"team"`
		Linescores []struct {
			DisplayValue string `json:"displayValue"`
		} `json:"linescores"`
	}

	for _, c := range comp.Competitors {
		if c.HomeAway == "home" {
			homeComp = c
		} else {
			awayComp = c
		}
	}

	awayName := awayComp.Team.DisplayName
	homeName := homeComp.Team.DisplayName
	awayAbb := awayComp.Team.Abbreviation
	homeAbb := homeComp.Team.Abbreviation

	awayScoreVal := awayComp.Score
	homeScoreVal := homeComp.Score

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

	titleLine := fmt.Sprintf(" %s  %s %s  @  %s %s\n",
		badgeStyled,
		style(awayAbb, ansiBold, format),
		awayScoreVal,
		style(homeAbb, ansiBold, format),
		homeScoreVal,
	)

	seriesStr := ""
	if len(comp.Series) > 0 {
		seriesStr = comp.Series[0].Summary
	}

	subTitleLine := fmt.Sprintf(" %s @ %s\n", awayName, homeName)
	if seriesStr != "" {
		subTitleLine = fmt.Sprintf(" %s @ %s (%s)\n", awayName, homeName, seriesStr)
	}
	if leagueLabel != "" {
		subTitleLine = subTitleLine + fmt.Sprintf(" League: %s\n", leagueLabel)
	}

	var header strings.Builder
	header.WriteString(style("========================================================================\n", ansiCyan, format))
	header.WriteString(titleLine)
	header.WriteString(style(subTitleLine, ansiGray, format))
	header.WriteString(style("========================================================================\n", ansiCyan, format))

	hasPoss := ""
	if awayComp.Possession {
		hasPoss = fmt.Sprintf("🏀 POSSESSION: %s (%s)", awayName, awayAbb)
	} else if homeComp.Possession {
		hasPoss = fmt.Sprintf("🏀 POSSESSION: %s (%s)", homeName, homeAbb)
	}

	header.WriteString("\n")
	if hasPoss != "" {
		header.WriteString("               " + style(hasPoss, ansiYellow, format) + "\n")
	}
	header.WriteString("               Status: " + style(detail, ansiYellow, format) + "\n\n")
	sb.WriteString(termPre(format, header.String()))

	// Period linescores
	numPeriods := 4
	if len(awayComp.Linescores) > 4 {
		numPeriods = len(awayComp.Linescores)
	}

	periodCols := []TableCol{{Title: "TEAM", Width: 10}}
	for i := 1; i <= numPeriods; i++ {
		colName := strconv.Itoa(i)
		if i > 4 {
			if numPeriods == 5 {
				colName = "OT"
			} else {
				colName = fmt.Sprintf("O%d", i-4)
			}
		}
		periodCols = append(periodCols, TableCol{Title: colName, Width: 3, Align: alignRight})
	}
	periodCols = append(periodCols, TableCol{Title: "T", Width: 3, Align: alignRight})

	pt := NewTable(format, periodCols...)

	buildPeriodRow := func(abbr, score string, lines []struct {
		DisplayValue string `json:"displayValue"`
	}) []TableCell {
		row := []TableCell{{Text: abbr}}
		for i := 1; i <= numPeriods; i++ {
			val := "-"
			if i-1 < len(lines) {
				val = lines[i-1].DisplayValue
			}
			row = append(row, TableCell{Text: val, Align: alignRight})
		}
		row = append(row, TableCell{Text: score, Align: alignRight})
		return row
	}
	pt.AddRow(buildPeriodRow(awayAbb, awayScoreVal, awayComp.Linescores)...)
	pt.AddRow(buildPeriodRow(homeAbb, homeScoreVal, homeComp.Linescores)...)
	sb.WriteString(pt.Render())

	// Team stats section
	var awayTeamStats, homeTeamStats ESPNBABoxscoreTeamStats
	hasTeamStats := false
	for _, ts := range summary.Boxscore.Teams {
		if ts.Team.ID == awayComp.Team.ID {
			awayTeamStats = ts
			hasTeamStats = true
		} else if ts.Team.ID == homeComp.Team.ID {
			homeTeamStats = ts
			hasTeamStats = true
		}
	}

	if hasTeamStats {
		getStatVal := func(ts ESPNBABoxscoreTeamStats, name string) string {
			for _, s := range ts.Statistics {
				if s.Name == name || s.Label == name {
					return s.DisplayValue
				}
			}
			return "-"
		}

		// Team stats counts table
		ct := NewTable(format,
			TableCol{Title: "TEAM", Width: 10},
			TableCol{Title: "PTS", Width: 3, Align: alignRight},
			TableCol{Title: "FG", Width: 7, Align: alignRight},
			TableCol{Title: "3PT", Width: 6, Align: alignRight},
			TableCol{Title: "FT", Width: 6, Align: alignRight},
			TableCol{Title: "ORB", Width: 3, Align: alignRight},
			TableCol{Title: "DRB", Width: 3, Align: alignRight},
			TableCol{Title: "TRB", Width: 3, Align: alignRight},
			TableCol{Title: "AST", Width: 3, Align: alignRight},
			TableCol{Title: "STL", Width: 3, Align: alignRight},
			TableCol{Title: "BLK", Width: 3, Align: alignRight},
		)
		ct.SetCaption("TEAM STATS")
		ct.AddRow(
			TableCell{Text: awayAbb},
			TableCell{Text: awayScoreVal, Align: alignRight},
			TableCell{Text: getStatVal(awayTeamStats, "fieldGoalsMade-fieldGoalsAttempted"), Align: alignRight},
			TableCell{Text: getStatVal(awayTeamStats, "threePointFieldGoalsMade-threePointFieldGoalsAttempted"), Align: alignRight},
			TableCell{Text: getStatVal(awayTeamStats, "freeThrowsMade-freeThrowsAttempted"), Align: alignRight},
			TableCell{Text: getStatVal(awayTeamStats, "offensiveRebounds"), Align: alignRight},
			TableCell{Text: getStatVal(awayTeamStats, "defensiveRebounds"), Align: alignRight},
			TableCell{Text: getStatVal(awayTeamStats, "totalRebounds"), Align: alignRight},
			TableCell{Text: getStatVal(awayTeamStats, "assists"), Align: alignRight},
			TableCell{Text: getStatVal(awayTeamStats, "steals"), Align: alignRight},
			TableCell{Text: getStatVal(awayTeamStats, "blocks"), Align: alignRight},
		)
		ct.AddRow(
			TableCell{Text: homeAbb},
			TableCell{Text: homeScoreVal, Align: alignRight},
			TableCell{Text: getStatVal(homeTeamStats, "fieldGoalsMade-fieldGoalsAttempted"), Align: alignRight},
			TableCell{Text: getStatVal(homeTeamStats, "threePointFieldGoalsMade-threePointFieldGoalsAttempted"), Align: alignRight},
			TableCell{Text: getStatVal(homeTeamStats, "freeThrowsMade-freeThrowsAttempted"), Align: alignRight},
			TableCell{Text: getStatVal(homeTeamStats, "offensiveRebounds"), Align: alignRight},
			TableCell{Text: getStatVal(homeTeamStats, "defensiveRebounds"), Align: alignRight},
			TableCell{Text: getStatVal(homeTeamStats, "totalRebounds"), Align: alignRight},
			TableCell{Text: getStatVal(homeTeamStats, "assists"), Align: alignRight},
			TableCell{Text: getStatVal(homeTeamStats, "steals"), Align: alignRight},
			TableCell{Text: getStatVal(homeTeamStats, "blocks"), Align: alignRight},
		)
		sb.WriteString(ct.Render())
		sb.WriteString("\n")

		// Shooting percentages table
		pt := NewTable(format,
			TableCol{Title: "TEAM", Width: 10},
			TableCol{Title: "FG%", Width: 6, Align: alignRight},
			TableCol{Title: "3P%", Width: 6, Align: alignRight},
			TableCol{Title: "FT%", Width: 6, Align: alignRight},
		)
		pt.SetCaption("SHOOTING %")
		pt.AddRow(
			TableCell{Text: awayAbb},
			TableCell{Text: getStatVal(awayTeamStats, "fieldGoalPct") + "%", Align: alignRight},
			TableCell{Text: getStatVal(awayTeamStats, "threePointFieldGoalPct") + "%", Align: alignRight},
			TableCell{Text: getStatVal(awayTeamStats, "freeThrowPct") + "%", Align: alignRight},
		)
		pt.AddRow(
			TableCell{Text: homeAbb},
			TableCell{Text: getStatVal(homeTeamStats, "fieldGoalPct") + "%", Align: alignRight},
			TableCell{Text: getStatVal(homeTeamStats, "threePointFieldGoalPct") + "%", Align: alignRight},
			TableCell{Text: getStatVal(homeTeamStats, "freeThrowPct") + "%", Align: alignRight},
		)
		sb.WriteString(pt.Render())
		sb.WriteString("\n")

		leadChanges := getStatVal(awayTeamStats, "leadChanges")
		awayLL := getStatVal(awayTeamStats, "largestLead")
		homeLL := getStatVal(homeTeamStats, "largestLead")

		var extra strings.Builder
		extra.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
		extra.WriteString(txt(fmt.Sprintf(" Lead Changes: %s\n", leadChanges), format))
		extra.WriteString(txt(fmt.Sprintf(" Biggest Lead: %s %s, %s %s\n", awayAbb, awayLL, homeAbb, homeLL), format))
		sb.WriteString(termPre(format, extra.String()))
	}

	// Leaders section
	type teamLeader struct {
		ptsName, ptsVal string
		astName, astVal string
		rebName, rebVal string
	}
	leaderMap := map[string]*teamLeader{}
	for _, tl := range summary.Leaders {
		l := &teamLeader{}
		for _, cat := range tl.Leaders {
			if len(cat.Leaders) > 0 {
				leaderAthlete := cat.Leaders[0]
				name := leaderAthlete.Athlete.DisplayName
				val := leaderAthlete.DisplayValue

				switch cat.Name {
				case "points":
					l.ptsName, l.ptsVal = name, val
				case "assists":
					l.astName, l.astVal = name, val
				case "rebounds":
					l.rebName, l.rebVal = name, val
				}
			}
		}
		leaderMap[tl.Team.ID] = l
	}

	awayL := leaderMap[awayComp.Team.ID]
	homeL := leaderMap[homeComp.Team.ID]

	if awayL != nil || homeL != nil {
		lt := NewTable(format,
			TableCol{Title: "TEAM", Width: 10},
			TableCol{Title: "POINTS", Width: 24},
			TableCol{Title: "ASSISTS", Width: 24},
			TableCol{Title: "REBOUNDS", Width: 24},
		)
		lt.SetCaption("TEAM LEADERS")

		leaderCell := func(l *teamLeader, name, val string) TableCell {
			if l == nil || name == "" {
				return TableCell{Text: "-"}
			}
			return TableCell{Text: fmt.Sprintf("%s (%s)", name, val)}
		}

		lt.AddRow(
			TableCell{Text: awayAbb},
			leaderCell(awayL, awayL.ptsName, awayL.ptsVal),
			leaderCell(awayL, awayL.astName, awayL.astVal),
			leaderCell(awayL, awayL.rebName, awayL.rebVal),
		)
		lt.AddRow(
			TableCell{Text: homeAbb},
			leaderCell(homeL, homeL.ptsName, homeL.ptsVal),
			leaderCell(homeL, homeL.astName, homeL.astVal),
			leaderCell(homeL, homeL.rebName, homeL.rebVal),
		)
		sb.WriteString(lt.Render())
		sb.WriteString("\n")
	}

	// Recent plays section
	var playBlock strings.Builder
	playBlock.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
	playBlock.WriteString(style(" RECENT PLAYS:\n", ansiBold+ansiCyan, format))

	plays := summary.Plays

	var lastPeriod int = -1
	for i := len(plays) - 1; i >= 0; i-- {
		play := plays[i]
		desc := play.Text
		if desc != "" {
			period := play.Period.Number
			clockVal := play.Clock.DisplayValue

			if lastPeriod != period {
				header := fmt.Sprintf("\n --- Quarter %d ---\n", period)
				playBlock.WriteString(style(header, ansiBold+ansiCyan, format))
				lastPeriod = period
			}

			isLast := i == len(plays)-1
			prefix := fmt.Sprintf(" [Q%d] ", period)

			var playLine string
			if isLast && state == "in" {
				playLine = style(prefix+desc, ansiGreen, format) + style(fmt.Sprintf(" (%s)\n", clockVal), ansiGreen, format) + style(" (Current Play)\n", ansiGreen, format)
			} else {
				playLine = txt(prefix+desc, format) + style(fmt.Sprintf(" (%s)\n", clockVal), ansiGray, format)
			}
			playBlock.WriteString(playLine)

			// Show score after play indented under the play
			scoreLine := fmt.Sprintf("      Score: %d-%d\n", play.AwayScore, play.HomeScore)
			playBlock.WriteString(style(scoreLine, ansiGray, format))
		}
	}

	if len(plays) == 0 {
		playBlock.WriteString(txt(" No plays recorded yet.\n", format))
	}
	sb.WriteString(termPre(format, playBlock.String()))

	// Boxscore Statistics Section
	sb.WriteString("\n")
	sb.WriteString(style("========================================================================\n", ansiCyan, format))
	sb.WriteString(style("                          BOXSCORE STATISTICS\n", ansiBold+ansiCyan, format))
	sb.WriteString(style("========================================================================\n", ansiCyan, format))

	// Determine team display order (away first, then home)
	var awayBox, homeBox *ESPNBABoxscoreTeam
	for i := range summary.Boxscore.Players {
		teamBox := &summary.Boxscore.Players[i]
		if teamBox.Team.ID == awayComp.Team.ID {
			awayBox = teamBox
		} else if teamBox.Team.ID == homeComp.Team.ID {
			homeBox = teamBox
		}
	}

	if awayBox != nil {
		sb.WriteString(renderNBABoxscore(*awayBox, format))
	}
	if homeBox != nil {
		sb.WriteString(renderNBABoxscore(*homeBox, format))
	}

	sb.WriteString(style("========================================================================\n", ansiCyan, format))

	if format == "ansi" {
		sb.WriteString(txt(" Run 'curl http://localhost:9090/nba' to return to the scoreboard.\n", format))
		sb.WriteString(style("========================================================================\n", ansiCyan, format))
	}

	return sb.String()
}

func fetchBasketballScoreboard(slug string, dateStr string) (*ESPNBAScoreboard, error) {
	espnDate := strings.ReplaceAll(dateStr, "-", "")
	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/basketball/%s/scoreboard?dates=%s", slug, espnDate)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESPN API returned status code %d", resp.StatusCode)
	}

	var sched ESPNBAScoreboard
	if err := json.NewDecoder(resp.Body).Decode(&sched); err != nil {
		return nil, err
	}
	return &sched, nil
}

func resolveBasketballLeague(dateStr string, gamePk string) (basketballLeague, bool) {
	for _, lg := range basketballLeagues {
		sched, err := fetchBasketballScoreboard(lg.Slug, dateStr)
		if err != nil {
			continue
		}
		for _, ev := range sched.Events {
			if ev.ID == gamePk {
				return lg, true
			}
		}
	}
	return basketballLeague{}, false
}

func handleNBASchedule(w http.ResponseWriter, r *http.Request) {
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

	espnDate := strings.ReplaceAll(dateStr, "-", "")
	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard?dates=%s", espnDate)
	resp, err := client.Get(url)
	if err != nil {
		http.Error(w, "Failed to connect to ESPN NBA API: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("ESPN NBA API returned status code %d", resp.StatusCode), http.StatusBadGateway)
		return
	}

	var sched ESPNBAScoreboard
	if err := json.NewDecoder(resp.Body).Decode(&sched); err != nil {
		http.Error(w, "Failed to decode schedule JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	format := getFormat(r)
	text := renderBasketballSchedule([]basketballScheduleEntry{{basketballLeagues[0], sched}}, dateStr, format, loc)
	writeResponse(w, format, text)
}

// handleBasketballSchedule shows every basketball league on one board.
func handleBasketballSchedule(w http.ResponseWriter, r *http.Request) {
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
	results := make([]basketballScheduleEntry, len(basketballLeagues))
	errs := make([]error, len(basketballLeagues))

	for i, lg := range basketballLeagues {
		wg.Add(1)
		go func(i int, lg basketballLeague) {
			defer wg.Done()
			sched, e := fetchBasketballScoreboard(lg.Slug, dateStr)
			if e != nil {
				errs[i] = e
				results[i] = basketballScheduleEntry{League: lg, Sched: ESPNBAScoreboard{}}
				return
			}
			results[i] = basketballScheduleEntry{League: lg, Sched: *sched}
		}(i, lg)
	}
	wg.Wait()

	format := getFormat(r)
	text := renderBasketballSchedule(results, dateStr, format, loc)
	writeResponse(w, format, text)
}

func handleNBAGame(w http.ResponseWriter, r *http.Request) {
	gamePk := r.PathValue("gamePk")
	if gamePk == "" {
		http.Error(w, "Missing game ID", http.StatusBadRequest)
		return
	}

	if serveHTMLWrapper(w, r) {
		return
	}

	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/basketball/nba/summary?event=%s", gamePk)
	resp, err := client.Get(url)
	if err != nil {
		http.Error(w, "Failed to connect to ESPN NBA API: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	format := getFormat(r)

	if resp.StatusCode != http.StatusOK {
		var sb strings.Builder
		sb.WriteString(style("========================================================================\n", ansiRed, format))
		sb.WriteString(style(fmt.Sprintf("                       ERROR: GAME NOT FOUND (%s)\n", gamePk), ansiBold+ansiRed, format))
		sb.WriteString(style("========================================================================\n", ansiRed, format))
		sb.WriteString(txt(" Details: The requested game ID is invalid or not yet available.\n", format))
		sb.WriteString(style("========================================================================\n", ansiRed, format))
		if format == "ansi" {
			sb.WriteString(txt(" Run 'curl http://localhost:9090/nba' to return to the scoreboard.\n", format))
			sb.WriteString(style("========================================================================\n", ansiRed, format))
		}

		writeResponse(w, format, sb.String())
		return
	}

	var summary ESPNBAGameSummary
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		http.Error(w, "Failed to decode game summary JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	text := renderNBAGame(summary, "NBA", format)
	writeResponse(w, format, text)
}

// handleBasketballGame shows a single game from any league. The league is taken
// from the ?league= query param when present, otherwise it is resolved by
// searching every league's scoreboard for the event ID.
func handleBasketballGame(w http.ResponseWriter, r *http.Request) {
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

	leagueSlug := r.URL.Query().Get("league")
	format := getFormat(r)

	var lg basketballLeague
	found := false
	if leagueSlug != "" {
		for _, l := range basketballLeagues {
			if l.Slug == leagueSlug {
				lg = l
				found = true
				break
			}
		}
	}
	if !found {
		lg, found = resolveBasketballLeague(dateStr, gamePk)
	}

	if !found {
		var sb strings.Builder
		sb.WriteString(style("========================================================================\n", ansiRed, format))
		sb.WriteString(style(fmt.Sprintf("                       ERROR: GAME NOT FOUND (%s)\n", gamePk), ansiBold+ansiRed, format))
		sb.WriteString(style("========================================================================\n", ansiRed, format))
		sb.WriteString(txt(" Details: The requested game ID was not found in any basketball league.\n", format))
		sb.WriteString(style("========================================================================\n", ansiRed, format))
		if format == "ansi" {
			sb.WriteString(txt(" Run 'curl http://localhost:9090/basketball' to return to the scoreboard.\n", format))
			sb.WriteString(style("========================================================================\n", ansiRed, format))
		}
		writeResponse(w, format, sb.String())
		return
	}

	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/basketball/%s/summary?event=%s", lg.Slug, gamePk)
	resp, err := client.Get(url)
	if err != nil {
		http.Error(w, "Failed to connect to ESPN API: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var sb strings.Builder
		sb.WriteString(style("========================================================================\n", ansiRed, format))
		sb.WriteString(style(fmt.Sprintf("                       ERROR: GAME NOT FOUND (%s)\n", gamePk), ansiBold+ansiRed, format))
		sb.WriteString(style("========================================================================\n", ansiRed, format))
		sb.WriteString(txt(" Details: The requested game ID is invalid or not yet available.\n", format))
		sb.WriteString(style("========================================================================\n", ansiRed, format))
		if format == "ansi" {
			sb.WriteString(txt(" Run 'curl http://localhost:9090/basketball' to return to the scoreboard.\n", format))
			sb.WriteString(style("========================================================================\n", ansiRed, format))
		}
		writeResponse(w, format, sb.String())
		return
	}

	var summary ESPNBAGameSummary
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		http.Error(w, "Failed to decode game summary JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	text := renderNBAGame(summary, lg.Name, format)
	writeResponse(w, format, text)
}

func handleAPINBAGames(w http.ResponseWriter, r *http.Request) {
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

	espnDate := strings.ReplaceAll(dateStr, "-", "")
	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard?dates=%s", espnDate)
	resp, err := client.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func handleAPINBAGameDetail(w http.ResponseWriter, r *http.Request) {
	gamePk := r.PathValue("gamePk")
	if gamePk == "" {
		http.Error(w, "Missing game ID", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/basketball/nba/summary?event=%s", gamePk)
	resp, err := client.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

// handleAPIBasketballGames returns every league's scoreboard merged into one JSON object.
func handleAPIBasketballGames(w http.ResponseWriter, r *http.Request) {
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
	results := make(map[string]*ESPNBAScoreboard, len(basketballLeagues))
	for _, lg := range basketballLeagues {
		wg.Add(1)
		go func(lg basketballLeague) {
			defer wg.Done()
			sched, e := fetchBasketballScoreboard(lg.Slug, dateStr)
			if e == nil {
				results[lg.Slug] = sched
			}
		}(lg)
	}
	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(results)
}

func handleAPIBasketballGameDetail(w http.ResponseWriter, r *http.Request) {
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

	leagueSlug := r.URL.Query().Get("league")
	var lg basketballLeague
	found := false
	if leagueSlug != "" {
		for _, l := range basketballLeagues {
			if l.Slug == leagueSlug {
				lg = l
				found = true
				break
			}
		}
	}
	if !found {
		lg, found = resolveBasketballLeague(dateStr, gamePk)
	}
	if !found {
		http.Error(w, "Game not found in any basketball league", http.StatusNotFound)
		return
	}

	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/basketball/%s/summary?event=%s", lg.Slug, gamePk)
	resp, err := client.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
