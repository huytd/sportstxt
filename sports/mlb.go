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
						Name string `json:"name"`
					} `json:"team"`
					Score int `json:"score"`
				} `json:"away"`
				Home struct {
					Team struct {
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

	sb.WriteString(style("================================================================================\n", ansiCyan, format))
	// Sports Selector row
	if format == "html" {
		sb.WriteString(txt("                         ", format) + style("[MLB]", ansiBold+ansiGreen, format) + txt("             ", format) + fmt.Sprintf(`<a href="/nba?date=%s" class="term-link">[NBA]</a>`, dateStr) + "\n")
	} else {
		sb.WriteString(txt("                         ", format) + style("[MLB]", ansiBold+ansiGreen, format) + txt("             ", format) + style("[NBA]", ansiGray, format) + "\n")
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
	Date        string `json:"date"`
	Opponent    string `json:"opponent"`
	IsHome      bool   `json:"isHome"`
	Result      string `json:"result"` // W, L, T
	Score       string `json:"score"` // e.g., "5-3"
	GamePk      int    `json:"gamePk"`
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

	isRaw := r.URL.Query().Get("raw") == "1"
	isCurl := strings.Contains(strings.ToLower(r.UserAgent()), "curl")

	if isRaw {
		text := renderSchedule(sched, dateStr, "html", loc)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(text))
	} else if isCurl {
		text := renderSchedule(sched, dateStr, "ansi", loc)
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(text))
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlPage))
	}
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	gamePk := r.PathValue("gamePk")
	if gamePk == "" {
		http.Error(w, "Missing game ID", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("https://statsapi.mlb.com/api/v1.1/game/%s/feed/live", gamePk)
	resp, err := client.Get(url)
	if err != nil {
		http.Error(w, "Failed to connect to MLB Stats API: "+err.Error(), http.StatusBadGateway)
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
			sb.WriteString(txt(" Run 'curl http://localhost:9090/' to return to the scoreboard.\n", format))
			sb.WriteString(style("========================================================================\n", ansiRed, format))
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if format == "ansi" {
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		}
		w.Write([]byte(sb.String()))
		return
	}

	var game GameFeedResponse
	if err := json.NewDecoder(resp.Body).Decode(&game); err != nil {
		http.Error(w, "Failed to decode game feed JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if isRaw {
		text := renderGame(game, "html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(text))
	} else if isCurl {
		text := renderGame(game, "ansi")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte(text))
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(htmlPage))
	}
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
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	Ties         int    `json:"ties"`
	PCT          string `json:"pct"`
	GB           string `json:"gb"`
	Venue        struct {
		City     string `json:"city"`
		Name     string `json:"name"`
	} `json:"venue"`
}

// LeagueStandingsResponse represents the response from the league-based standings API
type LeagueStandingsResponse struct {
	Records []struct {
		Division struct {
			Id int `json:"id"`
		} `json:"division"`
		TeamRecords []struct {
			Team struct {
				Id       int    `json:"id"`
				Name     string `json:"name"`
				Abbreviation string `json:"abbreviation"`
			} `json:"team"`
			LeagueRecord struct {
				Wins   int  `json:"wins"`
				Losses int  `json:"losses"`
				Ties   int  `json:"ties"`
				Pct    string `json:"pct"`
			} `json:"leagueRecord"`
			DivisionRank    string `json:"divisionRank"`
			LeagueRank      string `json:"leagueRank"`
			GamesBack       string `json:"gamesBack"`
		} `json:"teamRecords"`
	} `json:"records"`
}

// renderStandings creates the plain-text standings view
func renderStandings(standings LeagueStandingsResponse, teamMap map[int]struct{ Name, Abbreviation string }, format string) string {
	var sb strings.Builder

	loc := time.Now().Location()
	zoneName, _ := time.Now().In(loc).Zone()
	now := time.Now().Format("2006-01-02")
	title := fmt.Sprintf("MLB STANDINGS (%s %s)", zoneName, now)
	padding := (80 - len(title)) / 2
	if padding < 0 {
		padding = 0
	}

	sb.WriteString(style("================================================================================\n", ansiCyan, format))
	// Sports Selector row
	if format == "html" {
		sb.WriteString(txt("                         ", format) + style("[MLB]", ansiBold+ansiGreen, format) + txt("             ", format) + fmt.Sprintf(`<a href="/nba?date=%s" class="term-link">[NBA]</a>`, now) + "\n")
	} else {
		sb.WriteString(txt("                         ", format) + style("[MLB]", ansiBold+ansiGreen, format) + txt("             ", format) + style("[NBA]", ansiGray, format) + "\n")
	}
	sb.WriteString(style("================================================================================\n", ansiCyan, format))
	sb.WriteString(txt(strings.Repeat(" ", padding), format))
	sb.WriteString(style(title+"\n", ansiBold+ansiCyan, format))

	// Group standings by division: AL East, AL Central, AL West, NL East, NL Central, NL West
	alEast := []DivisionTeam{}
	alCentral := []DivisionTeam{}
	alWest := []DivisionTeam{}
	nlEast := []DivisionTeam{}
	nlCentral := []DivisionTeam{}
	nlWest := []DivisionTeam{}

	for _, rec := range standings.Records {
		divId := rec.Division.Id
		for _, tr := range rec.TeamRecords {
			tm, ok := teamMap[tr.Team.Id]
			if !ok {
				tm = struct{ Name, Abbreviation string }{tr.Team.Name, ""}
			}
			t := DivisionTeam{
				Abbreviation: tm.Abbreviation,
				Id:           tr.Team.Id,
				Name:         tm.Name,
				Wins:         tr.LeagueRecord.Wins,
				Losses:       tr.LeagueRecord.Losses,
				Ties:         tr.LeagueRecord.Ties,
				PCT:          tr.LeagueRecord.Pct,
				GB:           tr.GamesBack,
			}
			switch divId {
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

	renderDivision := func(teams []DivisionTeam, divAbbr string) {
		if len(teams) == 0 {
			return
		}
		sb.WriteString(style(fmt.Sprintf("\n%s DIVISION STANDINGS", strings.ToUpper(divAbbr)), ansiBold+ansiCyan, format))
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
		sb.WriteString(style(fmt.Sprintf(" %-10s %-28s %4s %4s %3s  %6s  %s  %s\n", "TEAM", "VENUE", "W", "L", "T", "PCT", "GB", "LINK"), ansiBold, format))
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))

		for _, t := range teams {
			link := fmt.Sprintf(`<a href="/mlb/team/%d" class="term-link">%s</a>`, t.Id, t.Abbreviation)
			if format == "ansi" {
				link = t.Abbreviation
			}
			row := fmt.Sprintf(" %-10s %-28s %4d %4d %3d  %6s  %s\n",
				link,
				t.Venue.Name,
				t.Wins,
				t.Losses,
				t.Ties,
				t.PCT,
				t.GB,
			)
			sb.WriteString(style(row, ansiBold, format))
		}
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
	}

	renderDivision(alEast, "AL EAST")
	renderDivision(alCentral, "AL CENTRAL")
	renderDivision(alWest, "AL WEST")
	renderDivision(nlEast, "NL EAST")
	renderDivision(nlCentral, "NL CENTRAL")
	renderDivision(nlWest, "NL WEST")

	sb.WriteString(style("================================================================================\n", ansiCyan, format))
	if format == "ansi" {
		sb.WriteString(txt(" Run 'curl http://localhost:9090/' to return to the scoreboard.\n", format))
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
	}

	return sb.String()
}

// fetchTeamGames fetches a team's schedule for a given date range
func fetchTeamGames(teamId int, teamName string, startDate string, endDate string) ([]TeamGame, error) {
	url := fmt.Sprintf("https://statsapi.mlb.com/api/v1/schedule?sportId=1&startDate=%s&endDate=%s", startDate, endDate)
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
			awayName := game.Teams.Away.Team.Name
			homeName := game.Teams.Home.Team.Name
			isHome := homeName == teamName
			gameDate := date.Date
			state := game.Status.DetailedState
			awayScore := game.Teams.Away.Score
			homeScore := game.Teams.Home.Score
			result := "-"
			if state == "Final" || state == "Game Over" {
				if isHome {
					result = fmt.Sprintf("%d-%d", awayScore, homeScore)
				} else {
					result = fmt.Sprintf("%d-%d", homeScore, awayScore)
				}
			}
			games = append(games, TeamGame{
				Date:    gameDate,
				Opponent: awayName + " @ " + homeName,
				IsHome: isHome,
				Result: result,
				Score:   result,
				GamePk:  game.GamePk,
			})
		}
	}
	return games, nil
}

// renderTeamPage creates the plain-text team detail view
func renderTeamPage(teamId int, teamName string, teamAbb string, format string) string {
	var sb strings.Builder

	title := fmt.Sprintf("%s - %s", teamAbb, teamName)
	padding := (80 - len(title)) / 2
	if padding < 0 {
		padding = 0
	}

	sb.WriteString(style("================================================================================\n", ansiCyan, format))
	if format == "html" {
		sb.WriteString(txt("                         ", format) + style("[MLB]", ansiBold+ansiGreen, format) + txt("             ", format) + fmt.Sprintf(`<a href="/nba?date=`+`%s`+`" class="term-link">[NBA]</a>`, time.Now().Format("2006-01-02")) + "\n")
	} else {
		sb.WriteString(txt("                         ", format) + style("[MLB]", ansiGray, format) + txt("             ", format) + style("[NBA]", ansiGray, format) + "\n")
	}
	sb.WriteString(style("================================================================================\n", ansiCyan, format))
	sb.WriteString(txt(strings.Repeat(" ", padding), format))
	sb.WriteString(style(title+"\n", ansiBold+ansiCyan, format))

	// Fetch team season stats from MLB Stats API
	statsUrl := fmt.Sprintf("https://statsapi.mlb.com/api/v1/teams/%d/stats?season=2025&seasonStage=regularSeason&granularity=full", teamId)
	resp, err := client.Get(statsUrl)
	if err != nil {
		sb.WriteString(style("\n ERROR: Could not fetch team stats.\n", ansiRed, format))
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
		return sb.String()
	}
	defer resp.Body.Close()

	var teamData struct {
		Team struct {
			Id       int    `json:"id"`
			Name     string `json:"name"`
			Abbreviation string `json:"abbreviation"`
		} `json:"team"`
		Stats []struct {
			Type struct {
				DisplayName string `json:"displayName"`
			} `json:"type"`
			Splits []struct {
				Stat struct {
					DisplayValue string `json:"displayValue"`
				} `json:"stat"`
			} `json:"splits"`
		} `json:"stats"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&teamData); err != nil {
		sb.WriteString(style("\n ERROR: Could not decode team stats.\n", ansiRed, format))
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
		return sb.String()
	}

	// Extract season stats from the API response
	seasonStats := ""
	for _, stat := range teamData.Stats {
		if stat.Type.DisplayName == "Season" {
			for _, split := range stat.Splits {
				statName := strings.ToUpper(strings.ReplaceAll(split.Stat.DisplayValue, " ", "_"))
				seasonStats += fmt.Sprintf("  %-30s %s\n", statName, split.Stat.DisplayValue)
			}
			break
		}
	}

	sb.WriteString(style("\n TEAM INFORMATION\n", ansiBold+ansiCyan, format))
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
	sb.WriteString(txt(fmt.Sprintf("  Team: %s\n", teamName), format))
	sb.WriteString(txt(fmt.Sprintf("  Abbreviation: %s\n", teamAbb), format))
	sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))

	if seasonStats != "" {
		sb.WriteString(style("\n SEASON STATISTICS\n", ansiBold+ansiCyan, format))
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
		sb.WriteString(txt(seasonStats, format))
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
	}

	// Fetch and display team's schedule (all games)
	today := time.Now()
	startDate := today.AddDate(0, 0, -365).Format("2006-01-02") // Last 365 days
	endDate := today.Format("2006-01-02")

	games, err := fetchTeamGames(teamId, teamName, startDate, endDate)
	if err != nil {
		sb.WriteString(style("\n WARNING: Could not fetch team schedule.\n", ansiYellow, format))
	} else if len(games) == 0 {
		sb.WriteString(txt("\n No games found for this team.\n", format))
	} else {
		sb.WriteString(style("\n ALL GAMES (Last 365 Days)\n", ansiBold+ansiCyan, format))
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
		sb.WriteString(style(fmt.Sprintf(" %-12s %-40s %8s %4s %s\n", "DATE", "OPPONENT", "HOME/AWAY", "SCORE", "RESULT"), ansiBold, format))
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))

		for _, g := range games {
			homeStr := "H"
			if !g.IsHome {
				homeStr = "A"
			}
			resultStr := g.Result
			if resultStr == "-" {
				resultStr = "UPCOMING"
			}
			gameLink := fmt.Sprintf(`<a href="/game/%d" class="term-link">%d</a>`, g.GamePk, g.GamePk)
			if format == "ansi" {
				gameLink = "-"
			}
			row := fmt.Sprintf(" %-12s %-40s %8s %4s  %s\n",
				g.Date,
				g.Opponent,
				homeStr,
				resultStr,
				gameLink,
			)
			sb.WriteString(style(row, ansiBold, format))
		}
		sb.WriteString(style("--------------------------------------------------------------------------------\n", ansiCyan, format))
		sb.WriteString(txt(fmt.Sprintf(" Total: %d games in last 365 days\n", len(games)), format))
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
	}

	if format == "ansi" {
		sb.WriteString(txt(" Run 'curl http://localhost:9090/' to return to the scoreboard.\n", format))
		sb.WriteString(style("================================================================================\n", ansiCyan, format))
	}

	return sb.String()
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
