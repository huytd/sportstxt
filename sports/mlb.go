package sports

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ScheduleResponse represents the response from statsapi.mlb.com schedule endpoint
type ScheduleResponse struct {
	Dates []struct {
		Date  string `json:"date"`
		Games []struct {
			GamePk int `json:"gamePk"`
			Status struct {
				DetailedState string `json:"detailedState"`
			} `json:"status"`
			Teams struct {
				Away struct {
					Team struct {
						Id   int    `json:"id"`
						Name string `json:"name"`
					} `json:"team"`
					Score int `json:"score"`
				} `json:"away"`
				Home struct {
					Team struct {
						Id   int    `json:"id"`
						Name string `json:"name"`
					} `json:"team"`
					Score int `json:"score"`
				} `json:"home"`
			} `json:"teams"`
			GameDate string `json:"gameDate"`
		} `json:"games"`
	} `json:"dates"`
}

// BoxscorePlayer represents individual player's game and season stats from boxscore
type BoxscorePlayer struct {
	Person struct {
		FullName string `json:"fullName"`
	} `json:"person"`
	Position struct {
		Abbreviation string `json:"abbreviation"`
	} `json:"position"`
	Stats struct {
		Batting struct {
			Summary     string `json:"summary"`
			AtBats      int    `json:"atBats"`
			Runs        int    `json:"runs"`
			Hits        int    `json:"hits"`
			Rbi         int    `json:"rbi"`
			BaseOnBalls int    `json:"baseOnBalls"`
			StrikeOuts  int    `json:"strikeOuts"`
		} `json:"batting"`
		Pitching struct {
			Summary        string `json:"summary"`
			InningsPitched string `json:"inningsPitched"`
			Hits           int    `json:"hits"`
			Runs           int    `json:"runs"`
			EarnedRuns     int    `json:"earnedRuns"`
			BaseOnBalls    int    `json:"baseOnBalls"`
			StrikeOuts     int    `json:"strikeOuts"`
			HomeRuns       int    `json:"homeRuns"`
		} `json:"pitching"`
	} `json:"stats"`
	SeasonStats struct {
		Batting struct {
			Avg string `json:"avg"`
			Obp string `json:"obp"`
			Slg string `json:"slg"`
			Ops string `json:"ops"`
			HR  int    `json:"homeRuns"`
			RBI int    `json:"rbi"`
		} `json:"batting"`
		Pitching struct {
			Era  string `json:"era"`
			Wins int    `json:"wins"`
			Loss int    `json:"losses"`
			SO   int    `json:"strikeOuts"`
			Whip string `json:"whip"`
		} `json:"pitching"`
	} `json:"seasonStats"`
}

// PlayEvent represents a pitch or action event within an at-bat
type PlayEvent struct {
	IsPitch bool `json:"isPitch"`
	Details struct {
		Description string `json:"description"`
		Event       string `json:"event"`
		Type        struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"type"`
		Call struct {
			Code        string `json:"code"`
			Description string `json:"description"`
		} `json:"call"`
	} `json:"details"`
	PitchData struct {
		StartSpeed float64 `json:"startSpeed"`
		EndSpeed   float64 `json:"endSpeed"`
	} `json:"pitchData"`
	Count struct {
		Balls   int `json:"balls"`
		Strikes int `json:"strikes"`
		Outs    int `json:"outs"`
	} `json:"count"`
}

// GameFeedResponse represents the response from statsapi.mlb.com feed/live endpoint
type GameFeedResponse struct {
	GameData struct {
		Status struct {
			DetailedState string `json:"detailedState"`
		} `json:"status"`
		Teams struct {
			Away struct {
				ID           int    `json:"id"`
				Name         string `json:"name"`
				Abbreviation string `json:"abbreviation"`
			} `json:"away"`
			Home struct {
				ID           int    `json:"id"`
				Name         string `json:"name"`
				Abbreviation string `json:"abbreviation"`
			} `json:"home"`
		} `json:"teams"`
	} `json:"gameData"`
	LiveData struct {
		Boxscore struct {
			Teams struct {
				Home struct {
					Batters  []int                     `json:"batters"`
					Pitchers []int                     `json:"pitchers"`
					Players  map[string]BoxscorePlayer `json:"players"`
				} `json:"home"`
				Away struct {
					Batters  []int                     `json:"batters"`
					Pitchers []int                     `json:"pitchers"`
					Players  map[string]BoxscorePlayer `json:"players"`
				} `json:"away"`
			} `json:"teams"`
		} `json:"boxscore"`
		Linescore struct {
			CurrentInning        int    `json:"currentInning"`
			CurrentInningOrdinal string `json:"currentInningOrdinal"`
			InningState          string `json:"inningState"`
			InningHalf           string `json:"inningHalf"`
			IsTopInning          bool   `json:"isTopInning"`
			Balls                int    `json:"balls"`
			Strikes              int    `json:"strikes"`
			Outs                 int    `json:"outs"`
			Innings              []struct {
				Num  int `json:"num"`
				Home struct {
					Runs   *int `json:"runs"`
					Hits   *int `json:"hits"`
					Errors *int `json:"errors"`
				} `json:"home"`
				Away struct {
					Runs   *int `json:"runs"`
					Hits   *int `json:"hits"`
					Errors *int `json:"errors"`
				} `json:"away"`
			} `json:"innings"`
			Teams struct {
				Home struct {
					Runs   *int `json:"runs"`
					Hits   *int `json:"hits"`
					Errors *int `json:"errors"`
				} `json:"home"`
				Away struct {
					Runs   *int `json:"runs"`
					Hits   *int `json:"hits"`
					Errors *int `json:"errors"`
				} `json:"away"`
			} `json:"teams"`
			Offense struct {
				Batter struct {
					ID       int    `json:"id"`
					FullName string `json:"fullName"`
				} `json:"batter"`
				Pitcher struct {
					ID       int    `json:"id"`
					FullName string `json:"fullName"`
				} `json:"pitcher"`
				First struct {
					FullName string `json:"fullName"`
				} `json:"first"`
				Second struct {
					FullName string `json:"fullName"`
				} `json:"second"`
				Third struct {
					FullName string `json:"fullName"`
				} `json:"third"`
			} `json:"offense"`
		} `json:"linescore"`
		Plays struct {
			CurrentPlay struct {
				PlayEvents []PlayEvent `json:"playEvents"`
				Result     struct {
					Description string `json:"description"`
				} `json:"result"`
				Count struct {
					Balls   int `json:"balls"`
					Strikes int `json:"strikes"`
					Outs    int `json:"outs"`
				} `json:"count"`
				Matchup struct {
					Batter struct {
						ID       int    `json:"id"`
						FullName string `json:"fullName"`
					} `json:"batter"`
					Pitcher struct {
						ID       int    `json:"id"`
						FullName string `json:"fullName"`
					} `json:"pitcher"`
					PostOnFirst struct {
						FullName string `json:"fullName"`
					} `json:"postOnFirst"`
					PostOnSecond struct {
						FullName string `json:"fullName"`
					} `json:"postOnSecond"`
					PostOnThird struct {
						FullName string `json:"fullName"`
					} `json:"postOnThird"`
				} `json:"matchup"`
			} `json:"currentPlay"`
			AllPlays []struct {
				Result struct {
					Description string `json:"description"`
					AwayScore   int    `json:"awayScore"`
					HomeScore   int    `json:"homeScore"`
				} `json:"result"`
				About struct {
					Inning      int  `json:"inning"`
					IsTopInning bool `json:"isTopInning"`
				} `json:"about"`
				PlayEvents []PlayEvent `json:"playEvents"`
			} `json:"allPlays"`
		} `json:"plays"`
	} `json:"liveData"`
}

// renderSchedule creates the plain-text scoreboard view
func renderSchedule(sched ScheduleResponse, dateStr string, format string, loc *time.Location) string {
	var sb strings.Builder

	zoneName, _ := time.Now().In(loc).Zone()
	title := fmt.Sprintf("MLB LIVE SCOREBOARD (%s %s)", dateStr, zoneName)
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

	if format != "html" {
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
		sb.WriteString(txt("                         ", format) + style("[MLB]", ansiBold+ansiGreen, format) + txt("             ", format) + style("[NBA]", ansiGray, format) + "\n")
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
	}
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
		prevLink := fmt.Sprintf(`<a href="/?date=%s" class="term-link">%s</a>`, prevDateStr, prevLinkText)
		nextLink := fmt.Sprintf(`<a href="/?date=%s" class="term-link">%s</a>`, nextDateStr, nextLinkText)
		sb.WriteString(prevLink + strings.Repeat(" ", spacerSize) + nextLink + "\n")
	} else {
		sb.WriteString(style(prevLinkText, ansiGreen, format) + strings.Repeat(" ", spacerSize) + style(nextLinkText, ansiGreen, format) + "\n")
	}
	sb.WriteString(style("================================================================================\n", ansiCyan, format))

	if len(sched.Dates) == 0 || len(sched.Dates[0].Games) == 0 {
		sb.WriteString(txt(" No games scheduled for this date.\n", format))
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
		return sb.String()
	}

	sb.WriteString(style(fmt.Sprintf(" %-8s %-8s %-19s %2s  @  %2s %-19s %-11s\n", "ID", "TIME", "AWAY TEAM", "R", "R", "HOME TEAM", "STATUS"), ansiBold, format))
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))

	for _, game := range sched.Dates[0].Games {
		idStr := strconv.Itoa(game.GamePk)
		awayName := game.Teams.Away.Team.Name
		homeName := game.Teams.Home.Team.Name

		if len(awayName) > 19 {
			awayName = awayName[:18] + "."
		}
		if len(homeName) > 19 {
			homeName = homeName[:18] + "."
		}

		awayScoreStr := "-"
		homeScoreStr := "-"

		state := game.Status.DetailedState
		isLive := state == "In Progress" || state == "Live" || state == "In Progress - Warmup" || state == "Warmup"
		isFinal := state == "Final" || state == "Game Over"

		if isLive || isFinal {
			awayScoreStr = strconv.Itoa(game.Teams.Away.Score)
			homeScoreStr = strconv.Itoa(game.Teams.Home.Score)
		}

		statusStr := state
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
		if t, err := time.Parse(time.RFC3339, game.GameDate); err == nil {
			gameTime = t.In(loc).Format("03:04 PM")
		}

		row := fmt.Sprintf(" %-8s %-8s %-19s %2s  @  %2s %-19s %-11s\n",
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
		sb.WriteString(txt(" Run 'curl http://localhost:9090/game/<ID>' to view a game in real-time.\n", format))
	} else {
		sb.WriteString(txt(" Click on a game ID to view the game in real-time.\n", format))
	}
	sb.WriteString(style("================================================================================\n", ansiCyan, format))

	if format == "html" {
		res := sb.String()
		for _, game := range sched.Dates[0].Games {
			idStr := strconv.Itoa(game.GamePk)
			link := fmt.Sprintf(`<a href="/game/%s" class="term-link">%s</a>`, idStr, idStr)
			res = strings.Replace(res, idStr, link, 1)

			awayTeam := game.Teams.Away.Team
			homeTeam := game.Teams.Home.Team
			awayName := awayTeam.Name
			homeName := homeTeam.Name
			if len(awayName) > 19 {
				awayName = awayName[:18] + "."
			}
			if len(homeName) > 19 {
				homeName = homeName[:18] + "."
			}
			awayLink := fmt.Sprintf(`<a href="/mlb/team/%d" class="term-link">%s</a>`, awayTeam.Id, awayName)
			homeLink := fmt.Sprintf(`<a href="/mlb/team/%d" class="term-link">%s</a>`, homeTeam.Id, homeName)
			res = strings.Replace(res, awayName, awayLink, 1)
			res = strings.Replace(res, homeName, homeLink, 1)
		}
		return res
	}

	return sb.String()
}

// findBoxscorePlayer searches for a player's stats by ID in both Home and Away boxscore maps
func findBoxscorePlayer(game GameFeedResponse, personId int) (BoxscorePlayer, bool) {
	idStr := "ID" + strconv.Itoa(personId)
	if player, ok := game.LiveData.Boxscore.Teams.Home.Players[idStr]; ok {
		return player, true
	}
	if player, ok := game.LiveData.Boxscore.Teams.Away.Players[idStr]; ok {
		return player, true
	}
	return BoxscorePlayer{}, false
}

// renderDiamondAndMatchup generates the text representation of bases and current count
func renderDiamondAndMatchup(game GameFeedResponse, format string) string {
	ls := game.LiveData.Linescore
	cp := game.LiveData.Plays.CurrentPlay

	state := game.GameData.Status.DetailedState
	isLive := state == "In Progress" || state == "Live" || state == "In Progress - Warmup" || state == "Warmup"

	var has1st, has2b, has3b bool
	var balls, strikes, outs int
	var batterName, pitcherName string
	var batterID, pitcherID int

	if isLive && (cp.Matchup.Batter.FullName != "" || cp.Matchup.Pitcher.FullName != "") {
		has1st = cp.Matchup.PostOnFirst.FullName != ""
		has2b = cp.Matchup.PostOnSecond.FullName != ""
		has3b = cp.Matchup.PostOnThird.FullName != ""
		balls = cp.Count.Balls
		strikes = cp.Count.Strikes
		outs = cp.Count.Outs
		batterName = cp.Matchup.Batter.FullName
		pitcherName = cp.Matchup.Pitcher.FullName
		batterID = cp.Matchup.Batter.ID
		pitcherID = cp.Matchup.Pitcher.ID
	} else {
		has1st = ls.Offense.First.FullName != ""
		has2b = ls.Offense.Second.FullName != ""
		has3b = ls.Offense.Third.FullName != ""
		balls = ls.Balls
		strikes = ls.Strikes
		outs = ls.Outs
		batterName = ls.Offense.Batter.FullName
		pitcherName = ls.Offense.Pitcher.FullName
		batterID = ls.Offense.Batter.ID
		pitcherID = ls.Offense.Pitcher.ID
	}

	b1 := style("◇", ansiGray, format)
	b2 := style("◇", ansiGray, format)
	b3 := style("◇", ansiGray, format)
	if has1st {
		b1 = style("◆", ansiYellow, format)
	}
	if has2b {
		b2 = style("◆", ansiYellow, format)
	}
	if has3b {
		b3 = style("◆", ansiYellow, format)
	}

	diamondLines := []string{
		"       2nd",
		fmt.Sprintf("       [%s]", b2),
		"  3rd       1st",
		fmt.Sprintf("  [%s]       [%s]", b3, b1),
	}

	ballsDots := dots(balls, 3)
	strikesDots := dots(strikes, 2)
	outsDots := dots(outs, 2)

	ballsDotsStyled := style(ballsDots, ansiGreen, format)
	strikesDotsStyled := style(strikesDots, ansiRed, format)
	outsDotsStyled := style(outsDots, ansiBold, format)

	if batterName == "" {
		batterName = "-"
	}
	if pitcherName == "" {
		pitcherName = "-"
	}

	// Fetch batter stats
	var batterToday, batterSeason string
	if batterID > 0 {
		if batterPlayer, ok := findBoxscorePlayer(game, batterID); ok {
			if batterPlayer.Stats.Batting.Summary != "" {
				batterToday = "Today:  " + batterPlayer.Stats.Batting.Summary
			}
			if batterPlayer.SeasonStats.Batting.Avg != "" {
				batterSeason = fmt.Sprintf("Season: %s AVG, %d HR, %d RBI",
					batterPlayer.SeasonStats.Batting.Avg,
					batterPlayer.SeasonStats.Batting.HR,
					batterPlayer.SeasonStats.Batting.RBI,
				)
			}
		}
	}

	// Fetch pitcher stats
	var pitcherToday, pitcherSeason string
	if pitcherID > 0 {
		if pitcherPlayer, ok := findBoxscorePlayer(game, pitcherID); ok {
			if pitcherPlayer.Stats.Pitching.Summary != "" {
				pitcherToday = "Today:  " + pitcherPlayer.Stats.Pitching.Summary
			}
			if pitcherPlayer.SeasonStats.Pitching.Era != "" {
				pitcherSeason = fmt.Sprintf("Season: %s ERA, %d-%d, %s WHIP",
					pitcherPlayer.SeasonStats.Pitching.Era,
					pitcherPlayer.SeasonStats.Pitching.Wins,
					pitcherPlayer.SeasonStats.Pitching.Loss,
					pitcherPlayer.SeasonStats.Pitching.Whip,
				)
			}
		}
	}

	matchupLines := []string{
		fmt.Sprintf("Balls:   %s", ballsDotsStyled),
		fmt.Sprintf("Strikes: %s", strikesDotsStyled),
		fmt.Sprintf("Outs:    %s", outsDotsStyled),
		fmt.Sprintf("Batter:  %s", style(batterName, ansiBold, format)),
	}
	if batterToday != "" {
		matchupLines = append(matchupLines, style("         "+batterToday, ansiGray, format))
	}
	if batterSeason != "" {
		matchupLines = append(matchupLines, style("         "+batterSeason, ansiGray, format))
	}

	matchupLines = append(matchupLines, fmt.Sprintf("Pitcher: %s", style(pitcherName, ansiBold, format)))
	if pitcherToday != "" {
		matchupLines = append(matchupLines, style("         "+pitcherToday, ansiGray, format))
	}
	if pitcherSeason != "" {
		matchupLines = append(matchupLines, style("         "+pitcherSeason, ansiGray, format))
	}

	maxLines := len(diamondLines)
	if len(matchupLines) > maxLines {
		maxLines = len(matchupLines)
	}

	var sb strings.Builder
	for i := 0; i < maxLines; i++ {
		diamondLine := ""
		if i < len(diamondLines) {
			diamondLine = diamondLines[i]
		}
		matchupLine := ""
		if i < len(matchupLines) {
			matchupLine = matchupLines[i]
		}

		rawDiamondLen := len([]rune(stripANSIAndHTML(diamondLine)))
		padSize := 26 - rawDiamondLen
		if padSize < 0 {
			padSize = 0
		}
		sb.WriteString(" " + diamondLine + strings.Repeat(" ", padSize) + matchupLine + "\n")
	}

	return sb.String()
}

// renderBattingBoxscore creates the plain-text batting boxscore table
func renderBattingBoxscore(players map[string]BoxscorePlayer, batterIDs []int, teamAbb string, format string) string {
	var sb strings.Builder

	sb.WriteString(style(fmt.Sprintf("\n %s BATTING\n", teamAbb), ansiBold+ansiCyan, format))
	sb.WriteString(style(" PLAYER                    AB  R  H RBI BB SO   AVG\n", ansiBold, format))
	sb.WriteString(style(" ---------------------------------------------------\n", ansiCyan, format))

	var totAB, totR, totH, totRBI, totBB, totSO int

	for _, id := range batterIDs {
		playerKey := "ID" + strconv.Itoa(id)
		player, ok := players[playerKey]
		if !ok {
			continue
		}

		name := player.Person.FullName
		pos := player.Position.Abbreviation
		displayName := name
		if pos != "" && pos != "P" {
			displayName = name + " " + pos
		}
		if len(displayName) > 24 {
			displayName = displayName[:23] + "."
		}

		ab := player.Stats.Batting.AtBats
		r := player.Stats.Batting.Runs
		h := player.Stats.Batting.Hits
		rbi := player.Stats.Batting.Rbi
		bb := player.Stats.Batting.BaseOnBalls
		so := player.Stats.Batting.StrikeOuts
		avg := player.SeasonStats.Batting.Avg
		if avg == "" {
			avg = ".000"
		}

		totAB += ab
		totR += r
		totH += h
		totRBI += rbi
		totBB += bb
		totSO += so

		row := fmt.Sprintf(" %-24s %2d %2d %2d %3d %2d %2d  %-5s\n",
			displayName, ab, r, h, rbi, bb, so, avg,
		)
		sb.WriteString(txt(row, format))
	}

	sb.WriteString(style(" ---------------------------------------------------\n", ansiCyan, format))
	totalsRow := fmt.Sprintf(" %-24s %2d %2d %2d %3d %2d %2d\n",
		"TOTALS", totAB, totR, totH, totRBI, totBB, totSO,
	)
	sb.WriteString(style(totalsRow, ansiBold, format))

	return sb.String()
}

// renderPitchingBoxscore creates the plain-text pitching boxscore table
func renderPitchingBoxscore(players map[string]BoxscorePlayer, pitcherIDs []int, teamAbb string, format string) string {
	var sb strings.Builder

	sb.WriteString(style(fmt.Sprintf("\n %s PITCHING\n", teamAbb), ansiBold+ansiCyan, format))
	sb.WriteString(style(" PLAYER                    IP  H  R ER BB SO HR   ERA\n", ansiBold, format))
	sb.WriteString(style(" ---------------------------------------------------\n", ansiCyan, format))

	for _, id := range pitcherIDs {
		playerKey := "ID" + strconv.Itoa(id)
		player, ok := players[playerKey]
		if !ok {
			continue
		}

		name := player.Person.FullName
		if len(name) > 24 {
			name = name[:23] + "."
		}

		ip := player.Stats.Pitching.InningsPitched
		h := player.Stats.Pitching.Hits
		r := player.Stats.Pitching.Runs
		er := player.Stats.Pitching.EarnedRuns
		bb := player.Stats.Pitching.BaseOnBalls
		so := player.Stats.Pitching.StrikeOuts
		hr := player.Stats.Pitching.HomeRuns
		era := player.SeasonStats.Pitching.Era
		if era == "" {
			era = "-.--"
		}

		row := fmt.Sprintf(" %-24s %4s %2d %2d %2d %2d %2d %2d  %-5s\n",
			name, ip, h, r, er, bb, so, hr, era,
		)
		sb.WriteString(txt(row, format))
	}

	sb.WriteString(style(" ---------------------------------------------------\n", ansiCyan, format))

	return sb.String()
}

// renderCurrentPitches creates the pitch list for the current at-bat
func renderCurrentPitches(game GameFeedResponse, format string) string {
	cp := game.LiveData.Plays.CurrentPlay
	var pitches []string

	balls := 0
	strikes := 0
	pitchNum := 1

	for _, e := range cp.PlayEvents {
		if e.IsPitch {
			speedStr := ""
			if e.PitchData.StartSpeed > 0 {
				speedStr = fmt.Sprintf("%dmph ", int(e.PitchData.StartSpeed+0.5))
			}

			pitchType := "pitch"
			if e.Details.Type.Description != "" {
				pitchType = strings.ToLower(e.Details.Type.Description)
			}

			desc := strings.ToLower(e.Details.Description)
			if desc == "" {
				desc = strings.ToLower(e.Details.Call.Description)
			}

			line := fmt.Sprintf("  P%d: %d-%d, %s%s, %s", pitchNum, balls, strikes, speedStr, pitchType, desc)
			pitches = append(pitches, line)
			pitchNum++
		}

		balls = e.Count.Balls
		strikes = e.Count.Strikes
	}

	if len(pitches) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(style(" PITCHES:\n\n", ansiBold+ansiCyan, format))
	for _, p := range pitches {
		sb.WriteString(txt(p+"\n", format))
	}
	return sb.String()
}

// TeamSeasonStats represents season statistics for a team
type TeamSeasonStats struct {
	Batting struct {
		Runs         int    `json:"runs"`
		Hits         int    `json:"hits"`
		HomeRuns     int    `json:"homeRuns"`
		RBI          int    `json:"rbi"`
		StolenBases  int    `json:"stolenBases"`
		BattingAvg   string `json:"battingAverage"`
		OBP          string `json:"onBasePercentage"`
		SLG          string `json:"slugPercentage"`
		OPS          string `json:"ops"`
		RunsAllowed  int    `json:"runsAllowed"`
	} `json:"batting"`
	Pitching struct {
		EarnedRunAverage string `json:"earnedRunAverage"`
		Wins             int    `json:"wins"`
		Losses           int    `json:"losses"`
		Saves            int    `json:"saves"`
		Strikeouts       int    `json:"strikeouts"`
		WHIP             string `json:"whip"`
	} `json:"pitching"`
}

// TeamGame represents a game in a team's season
type TeamGame struct {
	Date     string `json:"date"`
	Opponent string `json:"opponent"`
	IsHome   bool   `json:"isHome"`
	Result   string `json:"result"` // W, L, T, -
	Score    string `json:"score"` // e.g., "5-3"
	GamePk   int    `json:"gamePk"`
	State    string `json:"state"` // Final, In Progress, etc.
}

// colorizePitchSequence colorizes each letter of the sequence according to its pitch type color
func colorizePitchSequence(seq string, format string) string {
	var sb strings.Builder
	sb.WriteString(txt("(", format))
	for _, char := range seq {
		cStr := string(char)
		switch char {
		case 'B':
			sb.WriteString(style(cStr, ansiGreen, format))
		case 'C':
			sb.WriteString(style(cStr, ansiRed, format))
		case 'S':
			sb.WriteString(style(cStr, ansiMagenta, format))
		case 'F':
			sb.WriteString(style(cStr, ansiYellow, format))
		case 'X':
			sb.WriteString(style(cStr, ansiBlue, format))
		case 'D':
			sb.WriteString(style(cStr, ansiCyan, format))
		case 'E':
			sb.WriteString(style(cStr, ansiBold+ansiCyan, format))
		case '*':
			sb.WriteString(style(cStr, ansiGray, format))
		case 'W':
			sb.WriteString(style(cStr, ansiBold+ansiRed, format))
		case 'H':
			sb.WriteString(style(cStr, ansiBold+ansiYellow, format))
		case 'I':
			sb.WriteString(style(cStr, ansiBold+ansiGreen, format))
		case 'L':
			sb.WriteString(style(cStr, ansiBold+ansiMagenta, format))
		default:
			sb.WriteString(txt(cStr, format))
		}
	}
	sb.WriteString(txt(")", format))
	return sb.String()
}

// renderGame creates the detailed view of a game
func renderGame(game GameFeedResponse, format string) string {
	var sb strings.Builder

	awayName := game.GameData.Teams.Away.Name
	homeName := game.GameData.Teams.Home.Name
	awayAbb := game.GameData.Teams.Away.Abbreviation
	homeAbb := game.GameData.Teams.Home.Abbreviation

	awayScore := 0
	if game.LiveData.Linescore.Teams.Away.Runs != nil {
		awayScore = *game.LiveData.Linescore.Teams.Away.Runs
	}
	homeScore := 0
	if game.LiveData.Linescore.Teams.Home.Runs != nil {
		homeScore = *game.LiveData.Linescore.Teams.Home.Runs
	}

	state := game.GameData.Status.DetailedState
	badge := fmt.Sprintf("[%s]", strings.ToUpper(state))
	var badgeColor string
	switch state {
	case "In Progress", "Live", "In Progress - Warmup", "Warmup":
		badgeColor = ansiGreen
	case "Final", "Game Over":
		badgeColor = ansiBold
	case "Postponed", "Delayed", "Suspended":
		badgeColor = ansiRed
	default:
		badgeColor = ansiGray
	}

	badgeStyled := style(badge, badgeColor, format)

	inningInfo := ""
	if game.LiveData.Linescore.CurrentInningOrdinal != "" {
		inningInfo = fmt.Sprintf(" - %s %s", game.LiveData.Linescore.InningState, game.LiveData.Linescore.CurrentInningOrdinal)
	}

	if format == "html" {
		awayAbbLink := fmt.Sprintf(`<a href="/mlb/team/%d" class="term-link term-bold">%s</a>`, game.GameData.Teams.Away.ID, awayAbb)
		homeAbbLink := fmt.Sprintf(`<a href="/mlb/team/%d" class="term-link term-bold">%s</a>`, game.GameData.Teams.Home.ID, homeAbb)
		titleLine := fmt.Sprintf(" %s  %s %d  @  %s %d%s\n",
			badgeStyled,
			awayAbbLink,
			awayScore,
			homeAbbLink,
			homeScore,
			style(inningInfo, ansiYellow, format),
		)
		awayNameLink := fmt.Sprintf(`<a href="/mlb/team/%d" class="term-link">%s</a>`, game.GameData.Teams.Away.ID, awayName)
		homeNameLink := fmt.Sprintf(`<a href="/mlb/team/%d" class="term-link">%s</a>`, game.GameData.Teams.Home.ID, homeName)
		compareLink := fmt.Sprintf(`<a href="/mlb/compare?team1=%d&team2=%d" class="term-link">COMPARE</a>`, game.GameData.Teams.Away.ID, game.GameData.Teams.Home.ID)
		subTitleLine := fmt.Sprintf(`<span class="term-gray"> %s @ %s  |  %s</span>`+"\n", awayNameLink, homeNameLink, compareLink)
		sb.WriteString(style("========================================================================\n", ansiCyan, format))
		sb.WriteString(titleLine)
		sb.WriteString(subTitleLine)
		sb.WriteString(style("========================================================================\n", ansiCyan, format))
	} else {
		titleLine := fmt.Sprintf(" %s  %s %d  @  %s %d%s\n",
			badgeStyled,
			style(awayAbb, ansiBold, format),
			awayScore,
			style(homeAbb, ansiBold, format),
			homeScore,
			style(inningInfo, ansiYellow, format),
		)
		subTitleLine := fmt.Sprintf(" %s @ %s\n", awayName, homeName)
		sb.WriteString(style("========================================================================\n", ansiCyan, format))
		sb.WriteString(titleLine)
		sb.WriteString(style(subTitleLine, ansiGray, format))
		sb.WriteString(style("========================================================================\n", ansiCyan, format))
	}

	sb.WriteString(renderDiamondAndMatchup(game, format))
	if pStr := renderCurrentPitches(game, format); pStr != "" {
		sb.WriteString(pStr)
	}
	sb.WriteString("\n")

	numInnings := 9
	if len(game.LiveData.Linescore.Innings) > 9 {
		numInnings = len(game.LiveData.Linescore.Innings)
	}

	sb.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))

	currentInning := game.LiveData.Linescore.CurrentInning
	isLiveGame := state == "In Progress" || state == "Live" || state == "In Progress - Warmup" || state == "Warmup"

	sb.WriteString(style(" INNINGS    ", ansiBold, format))
	for i := 1; i <= numInnings; i++ {
		colStr := fmt.Sprintf("%2d ", i)
		if isLiveGame && i == currentInning {
			sb.WriteString(style(colStr, ansiBold+ansiYellow, format))
		} else {
			sb.WriteString(style(colStr, ansiBold, format))
		}
	}
	sb.WriteString(style("|  R  H  E\n", ansiBold, format))
	sb.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))

	// Away row
	sb.WriteString(txt(fmt.Sprintf(" %-10s ", awayAbb), format))
	for i := 1; i <= numInnings; i++ {
		val := "-"
		if i-1 < len(game.LiveData.Linescore.Innings) {
			inn := game.LiveData.Linescore.Innings[i-1]
			if inn.Away.Runs != nil {
				val = strconv.Itoa(*inn.Away.Runs)
			}
		}
		colStr := fmt.Sprintf("%2s ", val)
		isAwayActive := isLiveGame && i == currentInning && game.LiveData.Linescore.IsTopInning
		if isAwayActive {
			sb.WriteString(style(colStr, ansiBold+ansiYellow, format))
		} else {
			sb.WriteString(txt(colStr, format))
		}
	}

	awayR := "-"
	awayH := "-"
	awayE := "-"
	if game.LiveData.Linescore.Teams.Away.Runs != nil {
		awayR = strconv.Itoa(*game.LiveData.Linescore.Teams.Away.Runs)
	}
	if game.LiveData.Linescore.Teams.Away.Hits != nil {
		awayH = strconv.Itoa(*game.LiveData.Linescore.Teams.Away.Hits)
	}
	if game.LiveData.Linescore.Teams.Away.Errors != nil {
		awayE = strconv.Itoa(*game.LiveData.Linescore.Teams.Away.Errors)
	}
	sb.WriteString(txt(fmt.Sprintf("| %2s %2s %2s\n", awayR, awayH, awayE), format))

	// Home row
	sb.WriteString(txt(fmt.Sprintf(" %-10s ", homeAbb), format))
	for i := 1; i <= numInnings; i++ {
		val := "-"
		if i-1 < len(game.LiveData.Linescore.Innings) {
			inn := game.LiveData.Linescore.Innings[i-1]
			if inn.Home.Runs != nil {
				val = strconv.Itoa(*inn.Home.Runs)
			}
		}
		colStr := fmt.Sprintf("%2s ", val)
		isHomeActive := isLiveGame && i == currentInning && !game.LiveData.Linescore.IsTopInning
		if isHomeActive {
			sb.WriteString(style(colStr, ansiBold+ansiYellow, format))
		} else {
			sb.WriteString(txt(colStr, format))
		}
	}

	homeR := "-"
	homeH := "-"
	homeE := "-"
	if game.LiveData.Linescore.Teams.Home.Runs != nil {
		homeR = strconv.Itoa(*game.LiveData.Linescore.Teams.Home.Runs)
	}
	if game.LiveData.Linescore.Teams.Home.Hits != nil {
		homeH = strconv.Itoa(*game.LiveData.Linescore.Teams.Home.Hits)
	}
	if game.LiveData.Linescore.Teams.Home.Errors != nil {
		homeE = strconv.Itoa(*game.LiveData.Linescore.Teams.Home.Errors)
	}
	sb.WriteString(txt(fmt.Sprintf("| %2s %2s %2s\n", homeR, homeH, homeE), format))

	sb.WriteString(style("------------------------------------------------------------------------\n\n", ansiCyan, format))

	sb.WriteString(style(" RECENT PLAYS:\n", ansiBold+ansiCyan, format))
	plays := game.LiveData.Plays.AllPlays

	var lastInning int = -1
	var lastIsTop bool = false
	var hasLastSeen bool = false

	for i := len(plays) - 1; i >= 0; i-- {
		play := plays[i]
		desc := play.Result.Description
		if desc != "" {
			inning := play.About.Inning
			isTop := play.About.IsTopInning

			if !hasLastSeen || lastInning != inning || lastIsTop != isTop {
				halfStr := "Bottom"
				if isTop {
					halfStr = "Top"
				}
				header := fmt.Sprintf("\n --- %s %d ---\n", halfStr, inning)
				sb.WriteString(style(header, ansiBold+ansiCyan, format))
				lastInning = inning
				lastIsTop = isTop
				hasLastSeen = true
			}

			// Get the pitch sequence for this play
			var codes []string
			for _, e := range play.PlayEvents {
				if e.IsPitch {
					code := e.Details.Call.Code
					if code == "" {
						code = e.Details.Type.Code
					}
					if code != "" {
						codes = append(codes, code)
					}
				}
			}
			seqStyled := ""
			if len(codes) > 0 {
				seqStyled = " " + colorizePitchSequence(strings.Join(codes, ""), format)
			}

			isLast := i == len(plays)-1
			halfCode := "B"
			if isTop {
				halfCode = "T"
			}
			halfCode = fmt.Sprintf("%s%d", halfCode, inning)

			prefix := fmt.Sprintf(" [%s] ", halfCode)

			var playLine string
			if isLast {
				playLine = style(prefix+desc, ansiGreen, format) + seqStyled + style(" (Current Play)\n", ansiGreen, format)
			} else {
				playLine = txt(prefix+desc, format) + seqStyled + txt("\n", format)
			}
			sb.WriteString(playLine)

			// Show score after play indented under the play
			scoreLine := fmt.Sprintf(" Score: %d-%d\n", play.Result.AwayScore, play.Result.HomeScore)
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

	// Away batting
	sb.WriteString(renderBattingBoxscore(game.LiveData.Boxscore.Teams.Away.Players, game.LiveData.Boxscore.Teams.Away.Batters, awayAbb, format))
	// Home batting
	sb.WriteString(renderBattingBoxscore(game.LiveData.Boxscore.Teams.Home.Players, game.LiveData.Boxscore.Teams.Home.Batters, homeAbb, format))

	sb.WriteString(style("\n------------------------------------------------------------------------\n", ansiCyan, format))

	// Away pitching
	sb.WriteString(renderPitchingBoxscore(game.LiveData.Boxscore.Teams.Away.Players, game.LiveData.Boxscore.Teams.Away.Pitchers, awayAbb, format))
	// Home pitching
	sb.WriteString(renderPitchingBoxscore(game.LiveData.Boxscore.Teams.Home.Players, game.LiveData.Boxscore.Teams.Home.Pitchers, homeAbb, format))

	sb.WriteString(style("========================================================================\n", ansiCyan, format))
	sb.WriteString(style(" PITCH LEGEND:\n", ansiBold+ansiCyan, format))
	sb.WriteString(txt("  ", format) + style("B", ansiGreen, format) + txt(": Ball               ", format) + style("C", ansiRed, format) + txt(": Called Strike      ", format) + style("S", ansiMagenta, format) + txt(": Swinging Strike\n", format))
	sb.WriteString(txt("  ", format) + style("F", ansiYellow, format) + txt(": Foul               ", format) + style("X", ansiBlue, format) + txt(": In Play, Out       ", format) + style("D", ansiCyan, format) + txt(": In Play, No Out (Hit)\n", format))
	sb.WriteString(txt("  ", format) + style("E", ansiBold+ansiCyan, format) + txt(": In Play, Run(s)    ", format) + style("*", ansiGray, format) + txt(": Ball in Dirt       ", format) + style("W", ansiBold+ansiRed, format) + txt(": Swinging Strike (Pitchout)\n", format))
	sb.WriteString(txt("  ", format) + style("H", ansiBold+ansiYellow, format) + txt(": Hit By Pitch       ", format) + style("I", ansiBold+ansiGreen, format) + txt(": Intentional Ball   ", format) + style("L", ansiBold+ansiMagenta, format) + txt(": Foul Tip\n", format))
	sb.WriteString(style("========================================================================\n", ansiCyan, format))

	if format == "ansi" {
		sb.WriteString(txt(" Run 'curl http://localhost:9090/' to return to the scoreboard.\n", format))
		sb.WriteString(style("========================================================================\n", ansiCyan, format))
	}

	return sb.String()
}

func handleSchedule(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path != "/" {
		http.NotFound(w, r)
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

	url := fmt.Sprintf("https://statsapi.mlb.com/api/v1/schedule?sportId=1&date=%s", dateStr)
	resp, err := client.Get(url)
	if err != nil {
		http.Error(w, "Failed to connect to MLB Stats API: "+err.Error(), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, fmt.Sprintf("MLB Stats API returned status code %d", resp.StatusCode), http.StatusBadGateway)
		return
	}

	var sched ScheduleResponse
	if err := json.NewDecoder(resp.Body).Decode(&sched); err != nil {
		http.Error(w, "Failed to decode schedule JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	format := getFormat(r)
	text := renderSchedule(sched, dateStr, format, loc)
	writeResponse(w, format, text)
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	gamePk := r.PathValue("gamePk")
	if gamePk == "" {
		http.Error(w, "Missing game ID", http.StatusBadRequest)
		return
	}

	if serveHTMLWrapper(w, r) {
		return
	}

	url := fmt.Sprintf("https://statsapi.mlb.com/api/v1.1/game/%s/feed/live", gamePk)
	resp, err := client.Get(url)
	if err != nil {
		http.Error(w, "Failed to connect to MLB Stats API: "+err.Error(), http.StatusBadGateway)
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
			sb.WriteString(txt(" Run 'curl http://localhost:9090/' to return to the scoreboard.\n", format))
			sb.WriteString(style("========================================================================\n", ansiRed, format))
		}

		writeResponse(w, format, sb.String())
		return
	}

	var game GameFeedResponse
	if err := json.NewDecoder(resp.Body).Decode(&game); err != nil {
		http.Error(w, "Failed to decode game feed JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	text := renderGame(game, format)
	writeResponse(w, format, text)
}

func handleAPIGames(w http.ResponseWriter, r *http.Request) {
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

	url := fmt.Sprintf("https://statsapi.mlb.com/api/v1/schedule?sportId=1&date=%s", dateStr)
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

// DivisionTeam represents a team in standings grouped by division
type DivisionTeam struct {
	Abbreviation string `json:"abbreviation"`
	Id           int    `json:"id"`
	Name         string `json:"name"`
	City         string `json:"city"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	Ties         int    `json:"ties"`
	PCT          string `json:"pct"`
	GB           string `json:"gb"`
	League       string `json:"league"` // "AL" or "NL"
	DivisionRank string `json:"divisionRank"`
	WildCard     bool   `json:"wildCard"`
}

// TeamInfo represents basic team info from the teams API
type TeamInfo struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	LocationName string `json:"locationName"`
	League       struct {
		Name string `json:"name"`
	} `json:"league"`
	Division     struct {
		Name string `json:"name"`
	} `json:"division"`
}

// LeagueStandingsResponse represents the response from the league-based standings API
type LeagueStandingsResponse struct {
	Records []struct {
		League struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"league"`
		Division struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"division"`
		TeamRecords []struct {
			Team struct {
				Id           int    `json:"id"`
				Name         string `json:"name"`
				Abbreviation string `json:"abbreviation"`
			} `json:"team"`
			LeagueRecord struct {
				Wins   int    `json:"wins"`
				Losses int    `json:"losses"`
				Ties   int    `json:"ties"`
				Pct    string `json:"pct"`
			} `json:"leagueRecord"`
			DivisionRank    string `json:"divisionRank"`
			LeagueRank      string `json:"leagueRank"`
			GamesBack       string `json:"gamesBack"`
			SpringTraining  string `json:"springTraining"`
			DivWildcardSeed int    `json:"divWildcardSeed"`
		} `json:"teamRecords"`
	} `json:"records"`
}

// renderStandings creates the plain-text standings view
func renderStandings(standings LeagueStandingsResponse, teamMap map[int]TeamInfo, format string) string {
	var sb strings.Builder

	loc := time.Now().Location()
	zoneName, _ := time.Now().In(loc).Zone()
	seasonYear := time.Now().Year()
	title := fmt.Sprintf("MLB STANDINGS (%d) - %s", seasonYear, zoneName)
	padding := (80 - len(title)) / 2
	if padding < 0 {
		padding = 0
	}

	if format != "html" {
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
		sb.WriteString(txt(" [SCOREBOARD]   [COMPARE]   ", format) + style("[MLB]", ansiBold+ansiGreen, format) + txt("             ", format) + style("[NBA]", ansiGray, format) + "\n")
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
	}
	sb.WriteString(txt(strings.Repeat(" ", padding), format))
	sb.WriteString(style(title+"\n", ansiBold+ansiCyan, format))
	sb.WriteString(style("================================================================================\n", ansiCyan, format))

	// Group standings by division + wild card
	alEast := []DivisionTeam{}
	alCentral := []DivisionTeam{}
	alWest := []DivisionTeam{}
	nlEast := []DivisionTeam{}
	nlCentral := []DivisionTeam{}
	nlWest := []DivisionTeam{}
	alWildCard := []DivisionTeam{}
	nlWildCard := []DivisionTeam{}

	for _, rec := range standings.Records {
		leagueName := "AL"
		if strings.Contains(strings.ToLower(rec.League.Name), "national") || rec.League.Id == 104 {
			leagueName = "NL"
		}
		for _, tr := range rec.TeamRecords {
			tm, ok := teamMap[tr.Team.Id]
			if !ok {
				tm = TeamInfo{Name: tr.Team.Name, Abbreviation: tr.Team.Abbreviation}
			}
			t := DivisionTeam{
				Abbreviation: tm.Abbreviation,
				Id:           tr.Team.Id,
				Name:         tm.Name,
				City:         tm.LocationName,
				Wins:         tr.LeagueRecord.Wins,
				Losses:       tr.LeagueRecord.Losses,
				Ties:         tr.LeagueRecord.Ties,
				PCT:          tr.LeagueRecord.Pct,
				GB:           tr.GamesBack,
				League:       leagueName,
				DivisionRank: tr.DivisionRank,
				WildCard:     tr.DivWildcardSeed > 0,
			}
			// Wild card teams (not division leaders with a wildcard seed)
			if t.WildCard && t.DivisionRank != "1" {
				if leagueName == "AL" {
					alWildCard = append(alWildCard, t)
				} else {
					nlWildCard = append(nlWildCard, t)
				}
				continue
			}
			switch rec.Division.Id {
			case 200: // AL West
				alWest = append(alWest, t)
			case 201: // AL East
				alEast = append(alEast, t)
			case 202: // AL Central
				alCentral = append(alCentral, t)
			case 203: // NL West
				nlWest = append(nlWest, t)
			case 204: // NL East
				nlEast = append(nlEast, t)
			case 205: // NL Central
				nlCentral = append(nlCentral, t)
			}
		}
	}

	renderDivision := func(teams []DivisionTeam, divName string) {
		if len(teams) == 0 {
			return
		}
		sb.WriteString(style(fmt.Sprintf("\n %s", strings.ToUpper(divName)), ansiBold+ansiCyan, format))
		sb.WriteString("\n")
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
		sb.WriteString(style(fmt.Sprintf(" %-4s %-20s %4s %4s  %6s  %4s\n", "TEAM", "NAME", "W", "L", "PCT", "GB"), ansiBold, format))
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))

		for i, t := range teams {
			teamName := t.Name
			if len(teamName) > 20 {
				teamName = teamName[:19] + "."
			}

			// Use placeholder for team link (replaced after formatting to avoid HTML escaping)
			displayAbbr := fmt.Sprintf("__TL_%d__", t.Id)

			// Highlight division leader
			var rowStyle string
			if i == 0 {
				rowStyle = ansiBold
			} else if t.GB == "-" || t.GB == "0" {
				rowStyle = ansiBold
			} else {
				rowStyle = ""
			}

			gb := t.GB
			if gb == "-" {
				gb = "-"
			}

			abbrLen := len(t.Abbreviation)
			paddingSpaces := ""
			if abbrLen < 4 {
				paddingSpaces = strings.Repeat(" ", 4-abbrLen)
			}

			row := fmt.Sprintf(" %s%s %-20s %4d %4d  %6s  %4s\n",
				displayAbbr, paddingSpaces, teamName, t.Wins, t.Losses, t.PCT, gb,
			)
			if rowStyle != "" {
				sb.WriteString(style(row, rowStyle, format))
			} else {
				sb.WriteString(txt(row, format))
			}
		}
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
	}

	renderWildCard := func(teams []DivisionTeam, label string) {
		if len(teams) == 0 {
			return
		}
		sb.WriteString(style(fmt.Sprintf("\n %s WILD CARD", strings.ToUpper(label)), ansiBold+ansiCyan, format))
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
		sb.WriteString(style(fmt.Sprintf(" %-4s %-20s %4s %4s  %6s  %4s\n", "TEAM", "NAME", "W", "L", "PCT", "GB"), ansiBold, format))
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))

		for _, t := range teams {
			teamName := t.Name
			if len(teamName) > 20 {
				teamName = teamName[:19] + "."
			}
			// Use placeholder for team link (replaced after formatting to avoid HTML escaping)
			displayAbbr := fmt.Sprintf("__TL_%d__", t.Id)
			gb := t.GB
			if gb == "-" {
				gb = "-"
			}
			abbrLen := len(t.Abbreviation)
			paddingSpaces := ""
			if abbrLen < 4 {
				paddingSpaces = strings.Repeat(" ", 4-abbrLen)
			}
			row := fmt.Sprintf(" %s%s %-20s %4d %4d  %6s  %4s\n",
				displayAbbr, paddingSpaces, teamName, t.Wins, t.Losses, t.PCT, gb,
			)
			sb.WriteString(style(row, ansiGreen, format))
		}
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
	}

	renderDivision(alEast, "AL EAST")
	renderDivision(alCentral, "AL CENTRAL")
	renderDivision(alWest, "AL WEST")
	renderWildCard(alWildCard, "AL")
	renderDivision(nlEast, "NL EAST")
	renderDivision(nlCentral, "NL CENTRAL")
	renderDivision(nlWest, "NL WEST")
	renderWildCard(nlWildCard, "NL")

	sb.WriteString(style("================================================================================\n", ansiCyan, format))
	if format == "ansi" {
		sb.WriteString(txt(" Run 'curl http://localhost:9090/' to return to the scoreboard.\n", format))
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
	}

	result := sb.String()

	// Post-process: replace team link placeholders with actual links/abbreviations
	for _, rec := range standings.Records {
		for _, tr := range rec.TeamRecords {
			tm, ok := teamMap[tr.Team.Id]
			if !ok {
				tm = TeamInfo{Abbreviation: tr.Team.Abbreviation}
			}
			placeholder := fmt.Sprintf("__TL_%d__", tr.Team.Id)
			if format == "html" {
				link := fmt.Sprintf(`<a href="/mlb/team/%d" class="term-link">%s</a>`, tr.Team.Id, tm.Abbreviation)
				result = strings.ReplaceAll(result, placeholder, link)
			} else {
				result = strings.ReplaceAll(result, placeholder, tm.Abbreviation)
			}
		}
	}

	return result
}

// fetchTeamGames fetches a team's schedule for a given date range
func fetchTeamGames(teamId int, startDate string, endDate string) ([]TeamGame, error) {
	url := fmt.Sprintf("https://statsapi.mlb.com/api/v1/schedule?sportId=1&startDate=%s&endDate=%s&teamId=%d", startDate, endDate, teamId)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var sched ScheduleResponse
	if err := json.NewDecoder(resp.Body).Decode(&sched); err != nil {
		return nil, err
	}

	var games []TeamGame
	for _, date := range sched.Dates {
		for _, game := range date.Games {
			homeId := game.Teams.Home.Team.Id
			awayName := game.Teams.Away.Team.Name
			homeName := game.Teams.Home.Team.Name

			isHome := homeId == teamId

			opponent := awayName
			if isHome {
				opponent = awayName
			} else {
				opponent = homeName
			}

			gameDate := date.Date
			state := game.Status.DetailedState
			awayScore := game.Teams.Away.Score
			homeScore := game.Teams.Home.Score

			wl := "-"
			if state == "Final" || state == "Game Over" {
				if isHome {
					if homeScore > awayScore {
						wl = "W"
					} else if homeScore < awayScore {
						wl = "L"
					} else {
						wl = "T"
					}
				} else {
					if awayScore > homeScore {
						wl = "W"
					} else if awayScore < homeScore {
						wl = "L"
					} else {
						wl = "T"
					}
				}
			}

			score := fmt.Sprintf("%d-%d", awayScore, homeScore)
			if state != "Final" && state != "Game Over" && awayScore == 0 && homeScore == 0 {
				score = "-"
			}

			games = append(games, TeamGame{
				Date:     gameDate,
				Opponent: opponent,
				IsHome:   isHome,
				Result:   wl,
				Score:    score,
				GamePk:   game.GamePk,
				State:    state,
			})
		}
	}
	return games, nil
}

// renderTeamPage creates the plain-text team detail view
func renderTeamPage(teamId int, teamName string, teamAbb string, teamCity string, teamLeague string, teamDivision string, format string) string {
	var sb strings.Builder

	seasonYear := time.Now().Year()
	title := fmt.Sprintf("%s %s (%s)", teamCity, teamName, teamAbb)
	padding := (80 - len(title)) / 2
	if padding < 0 {
		padding = 0
	}

	if format != "html" {
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
		sb.WriteString(txt(" [SCOREBOARD]   [STANDINGS]   [COMPARE]   ", format) + style("[MLB]", ansiBold+ansiGreen, format) + txt("             ", format) + style("[NBA]", ansiGray, format) + "\n")
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
	}
	sb.WriteString(txt(strings.Repeat(" ", padding), format))
	sb.WriteString(style(title+"\n", ansiBold+ansiCyan, format))
	sb.WriteString(style("================================================================================\n", ansiCyan, format))

	// Team information section
	sb.WriteString(style("\n TEAM INFORMATION\n", ansiBold+ansiCyan, format))
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
	sb.WriteString(txt(fmt.Sprintf("  %-16s %s\n", "Name:", teamName), format))
	sb.WriteString(txt(fmt.Sprintf("  %-16s %s\n", "City:", teamCity), format))
	sb.WriteString(txt(fmt.Sprintf("  %-16s %s\n", "Abbreviation:", teamAbb), format))
	sb.WriteString(txt(fmt.Sprintf("  %-16s %s\n", "League:", teamLeague), format))
	sb.WriteString(txt(fmt.Sprintf("  %-16s %s\n", "Division:", teamDivision), format))
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))

	// Fetch team season stats from MLB Stats API (dynamic year)
	statsUrl := fmt.Sprintf("https://statsapi.mlb.com/api/v1/teams/%d/stats?season=%d&gameType=R&stats=season&group=hitting,pitching,fielding", teamId, seasonYear)
	resp, err := client.Get(statsUrl)
	if err != nil {
		sb.WriteString(style("\n WARNING: Could not fetch team stats.\n", ansiYellow, format))
	} else {
		defer resp.Body.Close()

		var teamData struct {
				Team struct {
					Id           int    `json:"id"`
					Name         string `json:"name"`
					Abbreviation string `json:"abbreviation"`
				} `json:"team"`
				Stats []struct {
					Type struct {
						DisplayName string `json:"displayName"`
					} `json:"type"`
					Group struct {
						DisplayName string `json:"displayName"`
					} `json:"group"`
					Splits []struct {
						Stat map[string]interface{} `json:"stat"`
					} `json:"splits"`
				} `json:"stats"`
			}

			if err := json.NewDecoder(resp.Body).Decode(&teamData); err != nil {
				sb.WriteString(style("\n WARNING: Could not decode team stats.\n", ansiYellow, format))
			} else {
				getStr := func(m map[string]interface{}, key string) string {
					if v, ok := m[key]; ok {
						return fmt.Sprintf("%v", v)
					}
					return "-"
				}

				formatVal := func(val string, width int, highlight bool) string {
					disp := val
					if highlight && val != "-" {
						if format == "html" {
							disp = fmt.Sprintf(`<span class="term-bold term-green term-highlight-better">%s</span>`, html.EscapeString(val))
						} else {
							disp = style(val, ansiBold+ansiGreen, format)
						}
					} else {
						if format == "html" {
							disp = html.EscapeString(val)
						}
					}
					return pad(disp, width, false, format)
				}

				isAboveAvg := func(valStr string, threshold float64) bool {
					if valStr == "-" {
						return false
					}
					v, err := strconv.ParseFloat(valStr, 64)
					if err != nil {
						return false
					}
					return v > threshold
				}

				isBelowAvg := func(valStr string, threshold float64) bool {
					if valStr == "-" {
						return false
					}
					v, err := strconv.ParseFloat(valStr, 64)
					if err != nil {
						return false
					}
					return v < threshold
				}

				for _, stat := range teamData.Stats {
					gn := stat.Group.DisplayName
					if gn == "" || len(stat.Splits) == 0 {
						continue
					}
					s := stat.Splits[0].Stat
					if gn == "hitting" {
						gp := getStr(s, "gamesPlayed")
						avg := getStr(s, "avg")
						obp := getStr(s, "obp")
						slg := getStr(s, "slg")
						ops := getStr(s, "ops")
						r := getStr(s, "runs")
						h := getStr(s, "hits")
						hr := getStr(s, "homeRuns")
						rbi := getStr(s, "rbi")
						bb := getStr(s, "baseOnBalls")
						so := getStr(s, "strikeOuts")
						stl := getStr(s, "stolenBases")

						sb.WriteString(style("\n BATTING\n", ansiBold+ansiCyan, format))
						sb.WriteString(style(fmt.Sprintf(" %5s %5s %5s %5s %5s %5s %5s %4s %4s %4s %4s %4s\n", "GP", "AVG", "OBP", "SLG", "OPS", "R", "H", "HR", "RBI", "BB", "SO", "SB"), ansiBold, format))
						sb.WriteString(style(" ------------------------------------------------------------------------\n", ansiCyan, format))
						row := fmt.Sprintf(" %s %s %s %s %s %s %s %s %s %s %s %s\n",
							formatVal(gp, 5, false),
							formatVal(avg, 5, isAboveAvg(avg, 0.243)),
							formatVal(obp, 5, isAboveAvg(obp, 0.319)),
							formatVal(slg, 5, isAboveAvg(slg, 0.400)),
							formatVal(ops, 5, isAboveAvg(ops, 0.719)),
							formatVal(r, 5, false),
							formatVal(h, 5, false),
							formatVal(hr, 4, false),
							formatVal(rbi, 4, false),
							formatVal(bb, 4, false),
							formatVal(so, 4, false),
							formatVal(stl, 4, false),
						)
						sb.WriteString(row)
						sb.WriteString(style(" ------------------------------------------------------------------------\n", ansiCyan, format))
					} else if gn == "pitching" {
						w := getStr(s, "wins")
						l := getStr(s, "losses")
						era := getStr(s, "era")
						whip := getStr(s, "whip")
						ip := getStr(s, "inningsPitched")
						so := getStr(s, "strikeOuts")
						bb := getStr(s, "baseOnBalls")
						hr := getStr(s, "homeRuns")
						sv := getStr(s, "saves")
						hld := getStr(s, "holds")
						bs := getStr(s, "blownSaves")
						avg := getStr(s, "avg")

						sb.WriteString(style("\n PITCHING\n", ansiBold+ansiCyan, format))
						sb.WriteString(style(fmt.Sprintf(" %4s %4s %5s %5s %6s %4s %4s %4s %4s %4s %4s %5s\n", "W", "L", "ERA", "WHIP", "IP", "SO", "BB", "HR", "SV", "HLD", "BS", "AVG"), ansiBold, format))
						sb.WriteString(style(" ------------------------------------------------------------------------\n", ansiCyan, format))
						row := fmt.Sprintf(" %s %s %s %s %s %s %s %s %s %s %s %s\n",
							formatVal(w, 4, false),
							formatVal(l, 4, false),
							formatVal(era, 5, isBelowAvg(era, 4.18)),
							formatVal(whip, 5, isBelowAvg(whip, 1.308)),
							formatVal(ip, 6, false),
							formatVal(so, 4, false),
							formatVal(bb, 4, false),
							formatVal(hr, 4, false),
							formatVal(sv, 4, false),
							formatVal(hld, 4, false),
							formatVal(bs, 4, false),
							formatVal(avg, 5, isBelowAvg(avg, 0.243)),
						)
						sb.WriteString(row)
						sb.WriteString(style(" ------------------------------------------------------------------------\n", ansiCyan, format))
					} else if gn == "fielding" {
						fpct := getStr(s, "fielding")
						e := getStr(s, "errors")
						dp := getStr(s, "doublePlays")
						pb := getStr(s, "passedBall")

						sb.WriteString(style("\n FIELDING\n", ansiBold+ansiCyan, format))
						sb.WriteString(style(fmt.Sprintf(" %5s %5s %5s %5s\n", "FPCT", "E", "DP", "PB"), ansiBold, format))
						sb.WriteString(style(" -----------------------------\n", ansiCyan, format))
						row := fmt.Sprintf(" %s %s %s %s\n",
							formatVal(fpct, 5, isAboveAvg(fpct, 0.985)),
							formatVal(e, 5, false),
							formatVal(dp, 5, false),
							formatVal(pb, 5, false),
						)
						sb.WriteString(row)
						sb.WriteString(style(" -----------------------------\n", ansiCyan, format))
					}
				}
			}
	}

	// Fetch and display recent games (last 30 days)
	today := time.Now()
	startDate := today.AddDate(0, 0, -30).Format("2006-01-02")
	endDate := today.Format("2006-01-02")

	games, err := fetchTeamGames(teamId, startDate, endDate)
	var allGames []TeamGame
	if err != nil {
		sb.WriteString(style("\n WARNING: Could not fetch recent games.\n", ansiYellow, format))
	} else if len(games) == 0 {
		sb.WriteString(txt("\n No games found in the last 30 days (off-season).\n", format))
	} else {
		allGames = games
		// Reverse to show most recent first
		for i, j := 0, len(games)-1; i < j; i, j = i+1, j-1 {
			games[i], games[j] = games[j], games[i]
		}

		sb.WriteString(style(fmt.Sprintf("\n RECENT GAMES (Last 30 Days) - %d Games\n", len(games)), ansiBold+ansiCyan, format))
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
		sb.WriteString(style(fmt.Sprintf(" %-10s %s %-22s %4s  %-5s\n", "DATE", "", "OPPONENT", "W/L", "SCORE"), ansiBold, format))
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))

		for _, g := range games {
			homeStr := "(H) "
			if !g.IsHome {
				homeStr = "(A) "
			}
			opponent := g.Opponent
			if len(opponent) > 22 {
				opponent = opponent[:19] + "..."
			}

			var wlPlaceholder string
			var rowStyle string
			switch g.Result {
			case "W":
				wlPlaceholder = "__WL_W__"
				rowStyle = ansiGreen
			case "L":
				wlPlaceholder = "__WL_L__"
				rowStyle = ""
			case "T":
				wlPlaceholder = "__WL_T__"
				rowStyle = ansiYellow
			default:
				if g.State == "In Progress" || g.State == "Live" {
					wlPlaceholder = "__WL_LIVE__"
					rowStyle = ansiGreen
				} else {
					wlPlaceholder = "__WL_EMPTY__"
					rowStyle = ansiGray
				}
			}

			// Build opponent with home/away indicator
			oppDisplay := homeStr + opponent

			// Use placeholder for game link (replaced after formatting to avoid HTML escaping)
			var scoreDisplay string
			if g.Score != "-" {
				scoreDisplay = fmt.Sprintf("__GL_%d__", g.GamePk)
			} else {
				scoreDisplay = fmt.Sprintf("__UL_%d__", g.GamePk)
			}

			row := fmt.Sprintf(" %-10s %-24s%4s  %-5s\n",
				g.Date, oppDisplay, wlPlaceholder, scoreDisplay,
			)
			if rowStyle != "" {
				sb.WriteString(style(row, rowStyle, format))
			} else {
				sb.WriteString(txt(row, format))
			}
		}
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
	}

	sb.WriteString(style("================================================================================\n", ansiCyan, format))
	if format == "ansi" {
		sb.WriteString(txt(" Run 'curl http://localhost:9090/' to return to the scoreboard.\n", format))
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
	}

	result := sb.String()

	// Post-process: replace W/L and game link placeholders with actual links or plain text
	wlReplacements := map[string]string{
		"__WL_W__":     style(" W ", ansiGreen, format),
		"__WL_L__":     style(" L ", ansiRed, format),
		"__WL_T__":     " T ",
		"__WL_LIVE__":  style(" LIVE", ansiGreen, format),
		"__WL_EMPTY__": "    ",
	}
	for placeholder, replacement := range wlReplacements {
		result = strings.ReplaceAll(result, placeholder, replacement)
	}

	for _, g := range allGames {
		gameLink := fmt.Sprintf("__GL_%d__", g.GamePk)
		if format == "html" {
			link := fmt.Sprintf(`<a href="/game/%d" class="term-link">%s</a>`, g.GamePk, g.Score)
			result = strings.ReplaceAll(result, gameLink, link)
		} else {
			result = strings.ReplaceAll(result, gameLink, g.Score)
		}

		upcomingLink := fmt.Sprintf("__UL_%d__", g.GamePk)
		if format == "html" {
			upcomingA := fmt.Sprintf(`<a href="/game/%d" class="term-link">UPCOMING</a>`, g.GamePk)
			result = strings.ReplaceAll(result, upcomingLink, upcomingA)
		} else {
			result = strings.ReplaceAll(result, upcomingLink, "UPCOMING")
		}
	}

	return result
}

func handleAPIGameDetail(w http.ResponseWriter, r *http.Request) {
	gamePk := r.PathValue("gamePk")
	if gamePk == "" {
		http.Error(w, "Missing game ID", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("https://statsapi.mlb.com/api/v1.1/game/%s/feed/live", gamePk)
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

// pad adjusts string length excluding HTML and ANSI codes
func pad(text string, width int, left bool, format string) string {
	raw := stripANSIAndHTML(text)
	diff := len(text) - len(raw)
	targetWidth := width + diff
	if left {
		return fmt.Sprintf("%-*s", targetWidth, text)
	}
	return fmt.Sprintf("%*s", targetWidth, text)
}

func compareStats(val1, val2 string, lowerIsBetter bool) (better1, better2 bool) {
	if val1 == "-" || val2 == "-" {
		return false, false
	}
	f1, err1 := strconv.ParseFloat(strings.TrimSpace(val1), 64)
	f2, err2 := strconv.ParseFloat(strings.TrimSpace(val2), 64)
	if err1 != nil || err2 != nil {
		return false, false
	}
	if f1 == f2 {
		return false, false
	}
	if lowerIsBetter {
		return f1 < f2, f2 < f1
	}
	return f1 > f2, f2 > f1
}

type statField struct {
	label         string
	key           string
	lowerIsBetter bool
}

var battingFields = []statField{
	{"Games Played", "gamesPlayed", false},
	{"Batting Average", "avg", false},
	{"On-Base Pct (OBP)", "obp", false},
	{"Slugging Pct (SLG)", "slg", false},
	{"On-Base+Slugging (OPS)", "ops", false},
	{"Runs", "runs", false},
	{"Hits", "hits", false},
	{"Home Runs", "homeRuns", false},
	{"Runs Batted In (RBI)", "rbi", false},
	{"Walks (BB)", "baseOnBalls", false},
	{"Strikeouts (SO)", "strikeOuts", true},
	{"Stolen Bases", "stolenBases", false},
}

var pitchingFields = []statField{
	{"Wins", "wins", false},
	{"Losses", "losses", true},
	{"Earned Run Avg (ERA)", "era", true},
	{"WHIP", "whip", true},
	{"Innings Pitched", "inningsPitched", false},
	{"Strikeouts (SO)", "strikeOuts", false},
	{"Walks (BB)", "baseOnBalls", true},
	{"Home Runs Allowed", "homeRuns", true},
	{"Saves", "saves", false},
	{"Holds", "holds", false},
	{"Blown Saves", "blownSaves", true},
	{"Opponent Avg", "avg", true},
}

var fieldingFields = []statField{
	{"Fielding Pct", "fielding", false},
	{"Errors", "errors", true},
	{"Double Plays", "doublePlays", false},
	{"Passed Balls", "passedBall", true},
}

func formatStatLine(label, key, group string, stats1, stats2 map[string]map[string]string, lowerIsBetter bool, format string) string {
	val1 := "-"
	if m, ok := stats1[group]; ok {
		if v, ok := m[key]; ok {
			val1 = v
		}
	}
	val2 := "-"
	if m, ok := stats2[group]; ok {
		if v, ok := m[key]; ok {
			val2 = v
		}
	}

	better1, better2 := compareStats(val1, val2, lowerIsBetter)

	disp1 := val1
	disp2 := val2

	if better1 {
		if format == "html" {
			disp1 = fmt.Sprintf(`<span class="term-bold term-green term-highlight-better">%s</span>`, html.EscapeString(val1))
		} else {
			disp1 = style(val1, ansiBold+ansiGreen, format)
		}
	} else if better2 {
		if format == "html" {
			disp2 = fmt.Sprintf(`<span class="term-bold term-green term-highlight-better">%s</span>`, html.EscapeString(val2))
		} else {
			disp2 = style(val2, ansiBold+ansiGreen, format)
		}
	}

	lbl := "  " + label
	col1 := pad(lbl, 26, true, format)
	col2 := pad(disp1, 27, true, format)
	col3 := pad(disp2, 27, true, format)

	return col1 + col2 + col3 + "\n"
}

func fetchStatsForTeam(teamId int, seasonYear int) (*TeamStatsResult, error) {
	statsUrl := fmt.Sprintf("https://statsapi.mlb.com/api/v1/teams/%d/stats?season=%d&gameType=R&stats=season&group=hitting,pitching,fielding", teamId, seasonYear)
	resp, err := client.Get(statsUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var res TeamStatsResult
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	return &res, nil
}

type TeamStatsResult struct {
	Stats []struct {
		Type struct {
			DisplayName string `json:"displayName"`
		} `json:"type"`
		Group struct {
			DisplayName string `json:"displayName"`
		} `json:"group"`
		Splits []struct {
			Stat map[string]interface{} `json:"stat"`
		} `json:"splits"`
	} `json:"stats"`
}

func extractStatsMap(res *TeamStatsResult) map[string]map[string]string {
	m := make(map[string]map[string]string)
	if res == nil {
		return m
	}
	for _, stat := range res.Stats {
		gn := stat.Group.DisplayName
		if gn == "" || len(stat.Splits) == 0 {
			continue
		}
		s := stat.Splits[0].Stat
		m[gn] = make(map[string]string)
		for k, v := range s {
			m[gn][k] = fmt.Sprintf("%v", v)
		}
	}
	return m
}

func dateShort(d string) string {
	if len(d) >= 10 {
		return d[5:10]
	}
	return d
}

func renderCompareTeams(team1Id, team2Id int, allTeams []TeamInfo, format string) string {
	var sb strings.Builder

	nowTime := time.Now()
	seasonYear := nowTime.Year()

	title := "MLB TEAM COMPARISON"
	if team1Id > 0 && team2Id > 0 {
		var abb1, abb2 string
		for _, t := range allTeams {
			if t.Id == team1Id {
				abb1 = t.Abbreviation
			}
			if t.Id == team2Id {
				abb2 = t.Abbreviation
			}
		}
		if abb1 != "" && abb2 != "" {
			title = fmt.Sprintf("MLB TEAM COMPARISON: %s VS %s", abb1, abb2)
		}
	}

	padding := (80 - len(title)) / 2
	if padding < 0 {
		padding = 0
	}

	if format != "html" {
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
		sb.WriteString(txt(" [SCOREBOARD]   [STANDINGS]   [COMPARE]   ", format) + style("[MLB]", ansiBold+ansiGreen, format) + txt("             ", format) + style("[NBA]", ansiGray, format) + "\n")
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
	}
	sb.WriteString(txt(strings.Repeat(" ", padding), format))
	sb.WriteString(style(title+"\n", ansiBold+ansiCyan, format))
	sb.WriteString(style("================================================================================\n", ansiCyan, format))

	// Find the compared teams in the list
	var teamA, teamB *TeamInfo
	for i := range allTeams {
		if allTeams[i].Id == team1Id {
			teamA = &allTeams[i]
		}
		if allTeams[i].Id == team2Id {
			teamB = &allTeams[i]
		}
	}

	// 1. Selector Form (Only for HTML format)
	if format == "html" {
		var opt1, opt2 strings.Builder
		opt1.WriteString(`<option value="">-- Select Team 1 --</option>`)
		opt2.WriteString(`<option value="">-- Select Team 2 --</option>`)
		for _, t := range allTeams {
			sel1 := ""
			if t.Id == team1Id {
				sel1 = "selected"
			}
			sel2 := ""
			if t.Id == team2Id {
				sel2 = "selected"
			}
			opt1.WriteString(fmt.Sprintf(`<option value="%d" %s>%s (%s)</option>`, t.Id, sel1, html.EscapeString(t.Name), html.EscapeString(t.Abbreviation)))
			opt2.WriteString(fmt.Sprintf(`<option value="%d" %s>%s (%s)</option>`, t.Id, sel2, html.EscapeString(t.Name), html.EscapeString(t.Abbreviation)))
		}

		sb.WriteString(`<div>
  <form id="compare-form" onsubmit="event.preventDefault(); const t1=document.getElementById('team1-select').value; const t2=document.getElementById('team2-select').value; if(t1&amp;&amp;t2){const href='/mlb/compare?team1='+t1+'&amp;team2='+t2; history.pushState(null,'',href); updateStatus(); document.getElementById('terminal-content').innerHTML='RETRIEVING FEED...'; fetchTerminalData();}" style="margin: 15px 0; display: flex; gap: 15px; align-items: center; flex-wrap: wrap;">
    <div>
      <label for="team1-select" style="font-weight: bold; margin-right: 8px;">Team 1:</label>
      <select id="team1-select" style="background: var(--term-container-bg); color: var(--color-primary); border: 1px solid var(--term-border); font-family: inherit; font-size: 14px; padding: 4px 8px; border-radius: 4px;">
` + opt1.String() + `
      </select>
    </div>
    <div>
      <label for="team2-select" style="font-weight: bold; margin-right: 8px;">Team 2:</label>
      <select id="team2-select" style="background: var(--term-container-bg); color: var(--color-primary); border: 1px solid var(--term-border); font-family: inherit; font-size: 14px; padding: 4px 8px; border-radius: 4px;">
` + opt2.String() + `
      </select>
    </div>
    <div>
      <label style="font-weight: bold; visibility: hidden;">Compare</label>
      <button type="submit" style="background: var(--color-primary); color: var(--term-container-bg); border: 1px solid var(--term-border); font-family: inherit; font-size: 14px; padding: 4px 12px; border-radius: 4px; cursor: pointer;">Compare</button>
    </div>
  </form>
</div>`)
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
	} else {
		// In curl/ANSI format, list options or show command help
		if teamA == nil || teamB == nil {
			sb.WriteString(style("\n HOW TO COMPARE TEAMS:\n", ansiBold+ansiCyan, format))
			sb.WriteString(txt(" Specify team IDs in query parameters: ?team1=<id1>&team2=<id2>\n", format))
			sb.WriteString(txt(" Example: curl \"http://localhost:9090/mlb/compare?team1=130&team2=147\"\n\n", format))
			sb.WriteString(style(" AVAILABLE MLB TEAMS:\n", ansiBold+ansiCyan, format))
			sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
			for i, t := range allTeams {
				sb.WriteString(txt(fmt.Sprintf("  %-4d: %-30s (%s)", t.Id, t.Name, t.Abbreviation), format))
				if (i+1)%2 == 0 {
					sb.WriteString("\n")
				} else {
					sb.WriteString("   | ")
				}
			}
			if len(allTeams)%2 != 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
			return sb.String()
		}
	}

	if teamA == nil || teamB == nil {
		if format == "html" {
			sb.WriteString(style("\n Select two teams above to compare their season statistics.\n\n", ansiYellow, format))
			sb.WriteString(style("================================================================================\n", ansiCyan, format))
		}
		return sb.String()
	}

	// Fetch Stats and Recent Games concurrently
	type fetchRes struct {
		stats *TeamStatsResult
		err   error
	}
	type gamesRes struct {
		games []TeamGame
		err   error
	}

	chStats1 := make(chan fetchRes, 1)
	chStats2 := make(chan fetchRes, 1)
	chGames1 := make(chan gamesRes, 1)
	chGames2 := make(chan gamesRes, 1)

	go func() {
		s, err := fetchStatsForTeam(team1Id, seasonYear)
		chStats1 <- fetchRes{s, err}
	}()
	go func() {
		s, err := fetchStatsForTeam(team2Id, seasonYear)
		chStats2 <- fetchRes{s, err}
	}()

	today := time.Now()
	startDate := today.AddDate(0, 0, -30).Format("2006-01-02")
	endDate := today.Format("2006-01-02")

	go func() {
		g, err := fetchTeamGames(team1Id, startDate, endDate)
		chGames1 <- gamesRes{g, err}
	}()
	go func() {
		g, err := fetchTeamGames(team2Id, startDate, endDate)
		chGames2 <- gamesRes{g, err}
	}()

	resStats1 := <-chStats1
	resStats2 := <-chStats2
	resGames1 := <-chGames1
	resGames2 := <-chGames2

	var stats1, stats2 map[string]map[string]string
	if resStats1.err == nil {
		stats1 = extractStatsMap(resStats1.stats)
	}
	if resStats2.err == nil {
		stats2 = extractStatsMap(resStats2.stats)
	}

	// 2. Team Info Comparison
	sb.WriteString(style("\n TEAM INFORMATION\n", ansiBold+ansiCyan, format))
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
	
	headerText := fmt.Sprintf("  %-24s %-27s %-27s\n", "STAT", teamA.Abbreviation, teamB.Abbreviation)
	sb.WriteString(style(headerText, ansiBold, format))
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))

	sb.WriteString(txt(fmt.Sprintf("  %-24s %-27s %-27s\n", "Team Name:", teamA.Name, teamB.Name), format))
	sb.WriteString(txt(fmt.Sprintf("  %-24s %-27s %-27s\n", "City:", teamA.LocationName, teamB.LocationName), format))
	
	leagueA := "AL"
	if strings.Contains(strings.ToLower(teamA.League.Name), "national") {
		leagueA = "NL"
	}
	leagueB := "AL"
	if strings.Contains(strings.ToLower(teamB.League.Name), "national") {
		leagueB = "NL"
	}
	sb.WriteString(txt(fmt.Sprintf("  %-24s %-27s %-27s\n", "League:", leagueA, leagueB), format))
	sb.WriteString(txt(fmt.Sprintf("  %-24s %-27s %-27s\n", "Division:", teamA.Division.Name, teamB.Division.Name), format))
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))

	// 3. Batting Stats
	sb.WriteString(style("\n BATTING STATS\n", ansiBold+ansiCyan, format))
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
	for _, f := range battingFields {
		sb.WriteString(formatStatLine(f.label, f.key, "hitting", stats1, stats2, f.lowerIsBetter, format))
	}
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))

	// 4. Pitching Stats
	sb.WriteString(style("\n PITCHING STATS\n", ansiBold+ansiCyan, format))
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
	for _, f := range pitchingFields {
		sb.WriteString(formatStatLine(f.label, f.key, "pitching", stats1, stats2, f.lowerIsBetter, format))
	}
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))

	// 5. Fielding Stats
	sb.WriteString(style("\n FIELDING STATS\n", ansiBold+ansiCyan, format))
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
	for _, f := range fieldingFields {
		sb.WriteString(formatStatLine(f.label, f.key, "fielding", stats1, stats2, f.lowerIsBetter, format))
	}
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))

	// 6. Recent Games (Side by Side)
	games1 := resGames1.games
	games2 := resGames2.games

	// Reverse to show most recent first
	for i, j := 0, len(games1)-1; i < j; i, j = i+1, j-1 {
		games1[i], games1[j] = games1[j], games1[i]
	}
	for i, j := 0, len(games2)-1; i < j; i, j = i+1, j-1 {
		games2[i], games2[j] = games2[j], games2[i]
	}

	maxGames := len(games1)
	if len(games2) > maxGames {
		maxGames = len(games2)
	}
	if maxGames > 5 {
		maxGames = 5
	}

	sb.WriteString(style("\n RECENT GAMES (Last 30 Days)\n", ansiBold+ansiCyan, format))
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
	headerA := fmt.Sprintf(" %s RECENT GAMES", teamA.Abbreviation)
	headerB := fmt.Sprintf(" %s RECENT GAMES", teamB.Abbreviation)
	sb.WriteString(style(pad(headerA, 39, true, format)+"|"+pad(headerB, 40, true, format)+"\n", ansiBold, format))
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))

	formatGameColumn := func(g *TeamGame, format string) string {
		if g == nil {
			return "  -"
		}
		homeStr := "(H) "
		if !g.IsHome {
			homeStr = "(A) "
		}
		opponent := g.Opponent
		if len(opponent) > 12 {
			opponent = opponent[:9] + "..."
		}
		var wlDisplay string
		switch g.Result {
		case "W":
			wlDisplay = style("W", ansiGreen, format)
		case "L":
			wlDisplay = style("L", ansiRed, format)
		case "T":
			wlDisplay = "T"
		default:
			if g.State == "In Progress" || g.State == "Live" {
				wlDisplay = style("LIVE", ansiGreen, format)
			} else {
				wlDisplay = style("UPC", ansiGray, format)
			}
		}

		scoreStr := g.Score
		if scoreStr != "-" {
			if format == "html" {
				scoreStr = fmt.Sprintf(`<a href="/game/%d" class="term-link">%s</a>`, g.GamePk, scoreStr)
			}
		} else {
			if format == "html" {
				scoreStr = fmt.Sprintf(`<a href="/game/%d" class="term-link">UPCOMING</a>`, g.GamePk)
			} else {
				scoreStr = "UPCOMING"
			}
		}

		colText := fmt.Sprintf(" %-5s %s%-12s %s  %s", dateShort(g.Date), homeStr, opponent, wlDisplay, scoreStr)
		return colText
	}

	for i := 0; i < maxGames; i++ {
		var g1, g2 *TeamGame
		if i < len(games1) {
			g1 = &games1[i]
		}
		if i < len(games2) {
			g2 = &games2[i]
		}

		colA := formatGameColumn(g1, format)
		colB := formatGameColumn(g2, format)

		sb.WriteString(pad(colA, 39, true, format) + "|" + pad(colB, 40, true, format) + "\n")
	}
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))

	sb.WriteString(style("================================================================================\n", ansiCyan, format))
	if format == "ansi" {
		sb.WriteString(txt(" Run 'curl http://localhost:9090/' to return to the scoreboard.\n", format))
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
	}

	return sb.String()
}
