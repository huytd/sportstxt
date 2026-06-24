package sports

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	padding := (80 - len(title)) / 2
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
	sb.WriteString(style("================================================================================\n", ansiCyan, format))
	if format == "html" {
		sb.WriteString(txt("                         ", format) + fmt.Sprintf(`<a href="/?date=%s" class="term-link">[MLB]</a>`, dateStr) + txt("             ", format) + style("[NBA]", ansiBold+ansiGreen, format) + "\n")
	} else {
		sb.WriteString(txt("                         ", format) + style("[MLB]", ansiGray, format) + txt("             ", format) + style("[NBA]", ansiBold+ansiGreen, format) + "\n")
	}

	sb.WriteString(style("================================================================================\n", ansiCyan, format))
	sb.WriteString(txt(strings.Repeat(" ", padding), format))
	sb.WriteString(style(title+"\n", ansiBold+ansiCyan, format))
	
	// Date Navigation Row
	sb.WriteString(style("================================================================================\n", ansiCyan, format))
	sb.WriteString(txt(" ", format))
	prevLinkText := fmt.Sprintf("<< PREV DAY (%s)", prevDateStr)
	nextLinkText := fmt.Sprintf("NEXT DAY (%s) >>", nextDateStr)
	spacerSize := 79 - len(prevLinkText) - len(nextLinkText)
	if spacerSize < 1 {
		spacerSize = 1
	}
	if format == "html" {
		prevLink := fmt.Sprintf(`<a href="/nba?date=%s" class="term-link">%s</a>`, prevDateStr, prevLinkText)
		nextLink := fmt.Sprintf(`<a href="/nba?date=%s" class="term-link">%s</a>`, nextDateStr, nextLinkText)
		sb.WriteString(prevLink + strings.Repeat(" ", spacerSize) + nextLink + "\n")
	} else {
		sb.WriteString(style(prevLinkText, ansiGreen, format) + strings.Repeat(" ", spacerSize) + style(nextLinkText, ansiGreen, format) + "\n")
	}
	sb.WriteString(style("================================================================================\n", ansiCyan, format))

	if len(sched.Events) == 0 {
		sb.WriteString(txt(" No games scheduled for this date.\n", format))
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
		return sb.String()
	}

	sb.WriteString(style(fmt.Sprintf(" %-9s %-8s %-17s %3s  @  %3s %-17s %-11s\n", "ID", "TIME", "AWAY TEAM", "PTS", "PTS", "HOME TEAM", "STATUS"), ansiBold, format))
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))

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

		if len(awayName) > 17 {
			awayName = awayName[:16] + "."
		}
		if len(homeName) > 17 {
			homeName = homeName[:16] + "."
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
		if len(statusStr) > 11 {
			statusStr = statusStr[:10] + "."
		}

		var rowStyle string
		if isLive {
			rowStyle = ansiGreen
		} else if isFinal {
			rowStyle = ansiBold
		} else {
			rowStyle = ansiGray
		}

		gameTime := "--:--"
		if t, err := time.Parse(time.RFC3339, event.Date); err == nil {
			gameTime = t.In(loc).Format("03:04 PM")
		} else if t, err := time.Parse("2006-01-02T15:04Z", event.Date); err == nil {
			gameTime = t.In(loc).Format("03:04 PM")
		}

		row := fmt.Sprintf(" %-9s %-8s %-17s %3s  @  %3s %-17s %-11s\n",
			idStr,
			gameTime,
			awayName,
			awayScoreStr,
			homeScoreStr,
			homeName,
			statusStr,
		)
		sb.WriteString(style(row, rowStyle, format))
	}

	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
	if format == "ansi" {
		sb.WriteString(txt(" Run 'curl http://localhost:8080/nba/game/<ID>' to view a game in real-time.\n", format))
	} else {
		sb.WriteString(txt(" Click on a game ID to view the game in real-time.\n", format))
	}
	sb.WriteString(style("================================================================================\n", ansiCyan, format))

	if format == "html" {
		res := sb.String()
		for _, event := range sched.Events {
			idStr := event.ID
			link := fmt.Sprintf(`<a href="/nba/game/%s" class="term-link">%s</a>`, idStr, idStr)
			res = strings.Replace(res, idStr, link, 1)
		}
		return res
	}

	return sb.String()
}

// renderNBABoxscore generates a team's boxscore table
func renderNBABoxscore(teamEntry ESPNBABoxscoreTeam, format string) string {
	var sb strings.Builder

	sb.WriteString(style(fmt.Sprintf("\n %s STATISTICS\n", strings.ToUpper(teamEntry.Team.DisplayName)), ansiBold+ansiCyan, format))
	sb.WriteString(style(" PLAYER                    MIN PTS   FG   3PT    FT REB AST  TO STL BLK  +/-\n", ansiBold, format))
	sb.WriteString(style(" ---------------------------------------------------------------------------\n", ansiCyan, format))

	if len(teamEntry.Statistics) == 0 {
		sb.WriteString(txt(" No statistics available.\n", format))
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

	for _, athlete := range teamEntry.Statistics[0].Athletes {
		name := athlete.Athlete.ShortName
		pos := athlete.Athlete.Position.Abbreviation
		displayName := name
		if pos != "" {
			displayName = name + " " + pos
		}
		if len(displayName) > 24 {
			displayName = displayName[:23] + "."
		}

		var row string
		if athlete.DidNotPlay {
			reason := athlete.Reason
			if reason == "" {
				reason = "Did Not Play"
			}
			row = fmt.Sprintf(" %-24s DNP - %s\n", displayName, reason)
		} else {
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

			row = fmt.Sprintf(" %-24s %3s %3s %5s %5s %5s %3s %3s %2s %3s %3s %4s\n",
				displayName, min, pts, fg, tpt, ft, reb, ast, to, stl, blk, pm,
			)
		}
		sb.WriteString(txt(row, format))
	}

	sb.WriteString(style(" ---------------------------------------------------------------------------\n", ansiCyan, format))

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
	sb.WriteString(style(totalsRow, ansiBold, format))

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

	sb.WriteString(style("========================================================================\n", ansiCyan, format))
	sb.WriteString(titleLine)
	sb.WriteString(style(subTitleLine, ansiGray, format))
	sb.WriteString(style("========================================================================\n", ansiCyan, format))

	hasPoss := ""
	if awayComp.Possession {
		hasPoss = fmt.Sprintf("🏀 POSSESSION: %s (%s)", awayName, awayAbb)
	} else if homeComp.Possession {
		hasPoss = fmt.Sprintf("🏀 POSSESSION: %s (%s)", homeName, homeAbb)
	}

	sb.WriteString("\n")
	if hasPoss != "" {
		sb.WriteString("               " + style(hasPoss, ansiYellow, format) + "\n")
	}
	sb.WriteString("               Status: " + style(detail, ansiYellow, format) + "\n\n")

	// Period linescores
	numPeriods := 4
	if len(awayComp.Linescores) > 4 {
		numPeriods = len(awayComp.Linescores)
	}

	sb.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
	sb.WriteString(style(" PERIODS    ", ansiBold, format))
	for i := 1; i <= numPeriods; i++ {
		colName := strconv.Itoa(i)
		if i > 4 {
			if numPeriods == 5 {
				colName = "OT"
			} else {
				colName = fmt.Sprintf("O%d", i-4)
			}
		}
		sb.WriteString(style(fmt.Sprintf("%-3s", colName), ansiBold, format))
	}
	sb.WriteString(style("|   T\n", ansiBold, format))
	sb.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))

	// Away row
	sb.WriteString(txt(fmt.Sprintf(" %-10s ", awayAbb), format))
	for i := 1; i <= numPeriods; i++ {
		val := "-"
		if i-1 < len(awayComp.Linescores) {
			val = awayComp.Linescores[i-1].DisplayValue
		}
		sb.WriteString(txt(fmt.Sprintf("%-3s", val), format))
	}
	sb.WriteString(txt(fmt.Sprintf("| %3s\n", awayScoreVal), format))

	// Home row
	sb.WriteString(txt(fmt.Sprintf(" %-10s ", homeAbb), format))
	for i := 1; i <= numPeriods; i++ {
		val := "-"
		if i-1 < len(homeComp.Linescores) {
			val = homeComp.Linescores[i-1].DisplayValue
		}
		sb.WriteString(txt(fmt.Sprintf("%-3s", val), format))
	}
	sb.WriteString(txt(fmt.Sprintf("| %3s\n", homeScoreVal), format))

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
		sb.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
		sb.WriteString(style(" TEAM STATS:\n", ansiBold+ansiCyan, format))
		sb.WriteString(style("      PTS      FG    3PT     FT OR/DR/TR   A   S   B\n", ansiBold, format))

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

		sb.WriteString(txt(fmt.Sprintf(" %-4s %3s %7s %6s %6s %2s/%2s/%2s %3s %3s %3s\n",
			awayAbb, awayScoreVal, awayFG, away3PT, awayFT, awayOR, awayDR, awayTR, awayA, awayS, awayB,
		), format))
		sb.WriteString(txt(fmt.Sprintf("          %-7s %-6s %-6s\n",
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

		sb.WriteString(txt(fmt.Sprintf(" %-4s %3s %7s %6s %6s %2s/%2s/%2s %3s %3s %3s\n",
			homeAbb, homeScoreVal, homeFG, home3PT, homeFT, homeOR, homeDR, homeTR, homeA, homeS, homeB,
		), format))
		sb.WriteString(txt(fmt.Sprintf("          %-7s %-6s %-6s\n",
			homeFGP+"%", home3PTP+"%", homeFTP+"%",
		), format))

		leadChanges := getStatVal(awayTeamStats, "leadChanges")
		awayLL := getStatVal(awayTeamStats, "largestLead")
		homeLL := getStatVal(homeTeamStats, "largestLead")

		sb.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
		sb.WriteString(txt(fmt.Sprintf(" Lead Changes: %s\n", leadChanges), format))
		sb.WriteString(txt(fmt.Sprintf(" Biggest Lead: %s %s, %s %s\n", awayAbb, awayLL, homeAbb, homeLL), format))
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

	if awayLeadersStr != "" || homeLeadersStr != "" {
		sb.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
		sb.WriteString(style(" TEAM LEADERS:\n", ansiBold+ansiCyan, format))
		if awayLeadersStr != "" {
			sb.WriteString(txt(fmt.Sprintf("  %s: %s\n", awayAbb, awayLeadersStr), format))
		}
		if homeLeadersStr != "" {
			sb.WriteString(txt(fmt.Sprintf("  %s: %s\n", homeAbb, homeLeadersStr), format))
		}
	}

	// Recent plays section
	sb.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
	sb.WriteString(style(" RECENT PLAYS:\n", ansiBold+ansiCyan, format))

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
				sb.WriteString(style(header, ansiBold+ansiCyan, format))
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
			sb.WriteString(playLine)

			// Show score after play indented under the play
			scoreLine := fmt.Sprintf("      Score: %d-%d\n", play.AwayScore, play.HomeScore)
			sb.WriteString(style(scoreLine, ansiGray, format))
		}
	}

	if len(plays) == 0 {
		sb.WriteString(txt(" No plays recorded yet.\n", format))
	}

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
		sb.WriteString(txt(" Run 'curl http://localhost:8080/nba' to return to the scoreboard.\n", format))
		sb.WriteString(style("========================================================================\n", ansiCyan, format))
	}

	return sb.String()
}

func handleNBASchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

	isRaw := r.URL.Query().Get("raw") == "1"
	isCurl := strings.Contains(strings.ToLower(r.UserAgent()), "curl")

	if isRaw {
		text := renderNBASchedule(sched, dateStr, "html", loc)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(text))
	} else if isCurl {
		text := renderNBASchedule(sched, dateStr, "ansi", loc)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(text))
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlPage))
	}
}

func handleNBAGame(w http.ResponseWriter, r *http.Request) {
	gamePk := r.PathValue("gamePk")
	if gamePk == "" {
		http.Error(w, "Missing game ID", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("https://site.api.espn.com/apis/site/v2/sports/basketball/nba/summary?event=%s", gamePk)
	resp, err := client.Get(url)
	if err != nil {
		http.Error(w, "Failed to connect to ESPN NBA API: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	isCurl := strings.Contains(strings.ToLower(r.UserAgent()), "curl")
	isRaw := r.URL.Query().Get("raw") == "1"
	format := "html"
	if isCurl && !isRaw {
		format = "ansi"
	}

	if resp.StatusCode != http.StatusOK {
		var sb strings.Builder
		sb.WriteString(style("========================================================================\n", ansiRed, format))
		sb.WriteString(style(fmt.Sprintf("                       ERROR: GAME NOT FOUND (%s)\n", gamePk), ansiBold+ansiRed, format))
		sb.WriteString(style("========================================================================\n", ansiRed, format))
		sb.WriteString(txt(" Details: The requested game ID is invalid or not yet available.\n", format))
		sb.WriteString(style("========================================================================\n", ansiRed, format))
		if format == "ansi" {
			sb.WriteString(txt(" Run 'curl http://localhost:8080/nba' to return to the scoreboard.\n", format))
			sb.WriteString(style("========================================================================\n", ansiRed, format))
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if format == "ansi" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		}
		w.Write([]byte(sb.String()))
		return
	}

	var summary ESPNBAGameSummary
	if err := json.NewDecoder(resp.Body).Decode(&summary); err != nil {
		http.Error(w, "Failed to decode game summary JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if isRaw {
		text := renderNBAGame(summary, "html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(text))
	} else if isCurl {
		text := renderNBAGame(summary, "ansi")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(text))
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlPage))
	}
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
