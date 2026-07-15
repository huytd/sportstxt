package sports

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// ESPNBAScoreboard represents the response from ESPN NBA scoreboard API
type ESPNBAScoreboard struct {
	Events []struct {
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
	} `json:"events"`
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

// renderNBASchedule creates the plain-text scoreboard view for NBA
func renderNBASchedule(sched ESPNBAScoreboard, dateStr string, format string, loc *time.Location) string {
	var sb strings.Builder

	zoneName, _ := time.Now().In(loc).Zone()
	title := fmt.Sprintf("NBA LIVE SCOREBOARD (%s %s)", dateStr, zoneName)
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

	var banner strings.Builder
	if format != "html" {
		banner.WriteString(style("==============================================================================\n", ansiCyan, format))
		banner.WriteString(txt("           ", format) + style("[MLB]", ansiGray, format) + txt("             ", format) + style("[NBA]", ansiBold+ansiGreen, format) + txt("             ", format) + style("[TENNIS]", ansiGray, format) + "\n")
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
		prevLink := fmt.Sprintf(`<a href="/nba?date=%s" class="term-link">%s</a>`, prevDateStr, prevLinkText)
		nextLink := fmt.Sprintf(`<a href="/nba?date=%s" class="term-link">%s</a>`, nextDateStr, nextLinkText)
		banner.WriteString(prevLink + strings.Repeat(" ", spacerSize) + nextLink + "\n")
	} else {
		banner.WriteString(style(prevLinkText, ansiGreen, format) + strings.Repeat(" ", spacerSize) + style(nextLinkText, ansiGreen, format) + "\n")
	}
	banner.WriteString(style("==============================================================================\n", ansiCyan, format))
	sb.WriteString(termPre(format, banner.String()))

	if len(sched.Events) == 0 {
		sb.WriteString(termPre(format, txt(" No games scheduled for this date.\n", format)+
			style("==============================================================================\n", ansiCyan, format)))
		return sb.String()
	}

	sort.SliceStable(sched.Events, func(i, j int) bool {
		parseTime := func(s string) time.Time {
			if t, err := time.Parse(time.RFC3339, s); err == nil {
				return t
			}
			if t, err := time.Parse("2006-01-02T15:04Z", s); err == nil {
				return t
			}
			return time.Time{}
		}
		t1 := parseTime(sched.Events[i].Date)
		t2 := parseTime(sched.Events[j].Date)
		if !t1.IsZero() && !t2.IsZero() {
			return t1.Before(t2)
		}
		return sched.Events[i].Date < sched.Events[j].Date
	})

	t := NewTable(format,
		TableCol{Title: "ID", Width: 9},
		TableCol{Title: "TIME", Width: 8},
		TableCol{Title: "AWAY TEAM", Width: 17},
		TableCol{Title: "PTS", Width: 3, Align: alignRight},
		TableCol{Title: "PTS", Width: 3, Align: alignRight},
		TableCol{Title: "HOME TEAM", Width: 17},
		TableCol{Title: "STATUS", Width: 10},
	)

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

	for _, event := range sched.Events {
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

		t.AddRow(
			TableCell{Text: idStr, Link: "/nba/game/" + idStr},
			TableCell{Text: gameTime, Class: classFor(baseStyle), ANSI: baseStyle},
			TableCell{Text: awayName, Class: classFor(awayStyle), ANSI: awayStyle},
			TableCell{Text: awayScoreStr, Align: alignRight, Class: classFor(awayStyle), ANSI: awayStyle},
			TableCell{Text: homeScoreStr, Align: alignRight, Class: classFor(homeStyle), ANSI: homeStyle},
			TableCell{Text: homeName, Class: classFor(homeStyle), ANSI: homeStyle},
			TableCell{Text: statusStr, Class: classFor(baseStyle), ANSI: baseStyle},
		)
	}

	sb.WriteString(t.Render())

	var footer strings.Builder
	footer.WriteString(style("------------------------------------------------------------------------------\n", ansiCyan, format))
	if format == "ansi" {
		footer.WriteString(txt(" Run 'curl http://localhost:9090/nba/game/<ID>' to view a game in real-time.\n", format))
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
func renderNBAGame(summary ESPNBAGameSummary, format string) string {
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
		var ts strings.Builder
		ts.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
		ts.WriteString(style(" TEAM STATS:\n", ansiBold+ansiCyan, format))
		ts.WriteString(style("      PTS      FG    3PT     FT OR/DR/TR   A   S   B\n", ansiBold, format))

		getStatVal := func(ts ESPNBABoxscoreTeamStats, name string) string {
			for _, s := range ts.Statistics {
				if s.Name == name || s.Label == name {
					return s.DisplayValue
				}
			}
			return "-"
		}

		awayFG := getStatVal(awayTeamStats, "fieldGoalsMade-fieldGoalsAttempted")
		awayFGP := getStatVal(awayTeamStats, "fieldGoalPct")
		away3PT := getStatVal(awayTeamStats, "threePointFieldGoalsMade-threePointFieldGoalsAttempted")
		away3PTP := getStatVal(awayTeamStats, "threePointFieldGoalPct")
		awayFT := getStatVal(awayTeamStats, "freeThrowsMade-freeThrowsAttempted")
		awayFTP := getStatVal(awayTeamStats, "freeThrowPct")
		awayOR := getStatVal(awayTeamStats, "offensiveRebounds")
		awayDR := getStatVal(awayTeamStats, "defensiveRebounds")
		awayTR := getStatVal(awayTeamStats, "totalRebounds")
		awayA := getStatVal(awayTeamStats, "assists")
		awayS := getStatVal(awayTeamStats, "steals")
		awayB := getStatVal(awayTeamStats, "blocks")

		ts.WriteString(txt(fmt.Sprintf(" %-4s %3s %7s %6s %6s %2s/%2s/%2s %3s %3s %3s\n",
			awayAbb, awayScoreVal, awayFG, away3PT, awayFT, awayOR, awayDR, awayTR, awayA, awayS, awayB,
		), format))
		ts.WriteString(txt(fmt.Sprintf("          %-7s %-6s %-6s\n",
			awayFGP+"%", away3PTP+"%", awayFTP+"%",
		), format))

		homeFG := getStatVal(homeTeamStats, "fieldGoalsMade-fieldGoalsAttempted")
		homeFGP := getStatVal(homeTeamStats, "fieldGoalPct")
		home3PT := getStatVal(homeTeamStats, "threePointFieldGoalsMade-threePointFieldGoalsAttempted")
		home3PTP := getStatVal(homeTeamStats, "threePointFieldGoalPct")
		homeFT := getStatVal(homeTeamStats, "freeThrowsMade-freeThrowsAttempted")
		homeFTP := getStatVal(homeTeamStats, "freeThrowPct")
		homeOR := getStatVal(homeTeamStats, "offensiveRebounds")
		homeDR := getStatVal(homeTeamStats, "defensiveRebounds")
		homeTR := getStatVal(homeTeamStats, "totalRebounds")
		homeA := getStatVal(homeTeamStats, "assists")
		homeS := getStatVal(homeTeamStats, "steals")
		homeB := getStatVal(homeTeamStats, "blocks")

		ts.WriteString(txt(fmt.Sprintf(" %-4s %3s %7s %6s %6s %2s/%2s/%2s %3s %3s %3s\n",
			homeAbb, homeScoreVal, homeFG, home3PT, homeFT, homeOR, homeDR, homeTR, homeA, homeS, homeB,
		), format))
		ts.WriteString(txt(fmt.Sprintf("          %-7s %-6s %-6s\n",
			homeFGP+"%", home3PTP+"%", homeFTP+"%",
		), format))

		leadChanges := getStatVal(awayTeamStats, "leadChanges")
		awayLL := getStatVal(awayTeamStats, "largestLead")
		homeLL := getStatVal(homeTeamStats, "largestLead")

		ts.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
		ts.WriteString(txt(fmt.Sprintf(" Lead Changes: %s\n", leadChanges), format))
		ts.WriteString(txt(fmt.Sprintf(" Biggest Lead: %s %s, %s %s\n", awayAbb, awayLL, homeAbb, homeLL), format))
		sb.WriteString(termPre(format, ts.String()))
	}

	// Leaders section
	var awayLeadersStr, homeLeadersStr string

	for _, tl := range summary.Leaders {
		var ptsName, ptsVal string
		var astName, astVal string
		var rebName, rebVal string

		for _, cat := range tl.Leaders {
			if len(cat.Leaders) > 0 {
				leaderAthlete := cat.Leaders[0]
				name := leaderAthlete.Athlete.DisplayName
				val := leaderAthlete.DisplayValue

				switch cat.Name {
				case "points":
					ptsName = name
					ptsVal = val
				case "assists":
					astName = name
					astVal = val
				case "rebounds":
					rebName = name
					rebVal = val
				}
			}
		}

		leaderStr := fmt.Sprintf("%s (%s PTS), %s (%s AST), %s (%s REB)",
			ptsName, ptsVal, astName, astVal, rebName, rebVal,
		)

		if tl.Team.ID == awayComp.Team.ID {
			awayLeadersStr = leaderStr
		} else if tl.Team.ID == homeComp.Team.ID {
			homeLeadersStr = leaderStr
		}
	}

	var leaders strings.Builder
	if awayLeadersStr != "" || homeLeadersStr != "" {
		leaders.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
		leaders.WriteString(style(" TEAM LEADERS:\n", ansiBold+ansiCyan, format))
		if awayLeadersStr != "" {
			leaders.WriteString(txt(fmt.Sprintf("  %s: %s\n", awayAbb, awayLeadersStr), format))
		}
		if homeLeadersStr != "" {
			leaders.WriteString(txt(fmt.Sprintf("  %s: %s\n", homeAbb, homeLeadersStr), format))
		}
	}
	sb.WriteString(termPre(format, leaders.String()))

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
	text := renderNBASchedule(sched, dateStr, format, loc)
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

	text := renderNBAGame(summary, format)
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
