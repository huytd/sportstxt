package sports

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// ScheduleGame represents a single game in the schedule response
type ScheduleGame struct {
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
	Linescore *struct {
		CurrentInning        int    `json:"currentInning"`
		CurrentInningOrdinal string `json:"currentInningOrdinal"`
		InningState          string `json:"inningState"`
		IsTopInning          bool   `json:"isTopInning"`
		Balls                int    `json:"balls"`
		Strikes              int    `json:"strikes"`
		Outs                 int    `json:"outs"`
	} `json:"linescore"`
}

// ScheduleResponse represents the response from statsapi.mlb.com schedule endpoint
type ScheduleResponse struct {
	Dates []struct {
		Date  string         `json:"date"`
		Games []ScheduleGame `json:"games"`
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

// renderSingleGameBox renders a single game box as a slice of 5 styled text lines
func renderSingleGameBox(game ScheduleGame, loc *time.Location, format string, width int) []string {
	var lines []string

	awayName := game.Teams.Away.Team.Name
	homeName := game.Teams.Home.Team.Name

	awayScoreStr := "-"
	homeScoreStr := "-"

	state := game.Status.DetailedState
	isLive := state == "In Progress" || state == "Live" || state == "In Progress - Warmup" || state == "Warmup"
	isFinal := state == "Final" || state == "Game Over"

	if isLive || isFinal {
		awayScoreStr = strconv.Itoa(game.Teams.Away.Score)
		homeScoreStr = strconv.Itoa(game.Teams.Home.Score)
	}

	gameTime := "--:--"
	if t, err := time.Parse(time.RFC3339, game.GameDate); err == nil {
		gameTime = t.In(loc).Format("03:04 PM")
	}

	var header string
	var headerStyle string
	var awayStyle, homeStyle string
	var awayScoreStyle, homeScoreStyle string

	if isLive {
		headerStyle = ansiGreen
		if game.Linescore != nil {
			half := "BOT"
			if game.Linescore.IsTopInning {
				half = "TOP"
			}
			header = fmt.Sprintf("%s %s", half, game.Linescore.CurrentInningOrdinal)
		} else {
			header = "LIVE"
		}

		awayStyle = ansiGreen
		homeStyle = ansiGreen
		awayScoreStyle = ansiGreen
		homeScoreStyle = ansiGreen

		if game.Teams.Away.Score > game.Teams.Home.Score {
			awayStyle = ansiBold + ansiGreen
			awayScoreStyle = ansiBold + ansiGreen
		} else if game.Teams.Home.Score > game.Teams.Away.Score {
			homeStyle = ansiBold + ansiGreen
			homeScoreStyle = ansiBold + ansiGreen
		}
	} else if isFinal {
		header = "FINAL"
		headerStyle = ansiBold

		awayStyle = ""
		homeStyle = ""
		awayScoreStyle = ""
		homeScoreStyle = ""

		if game.Teams.Away.Score > game.Teams.Home.Score {
			awayStyle = ansiBold
			awayScoreStyle = ansiBold
		} else if game.Teams.Home.Score > game.Teams.Away.Score {
			homeStyle = ansiBold
			homeScoreStyle = ansiBold
		}
	} else {
		header = gameTime
		headerStyle = ansiGray
		awayStyle = ansiGray
		homeStyle = ansiGray
		awayScoreStyle = ansiGray
		homeScoreStyle = ansiGray
	}

	// Compute right-side content
	var awayRight, homeRight string
	awayRightVisible, homeRightVisible := 0, 0
	if isLive && game.Linescore != nil {
		ballsDots := dots(game.Linescore.Balls, 3)
		strikesDots := dots(game.Linescore.Strikes, 2)
		dotsStr := ballsDots + " " + strikesDots
		awayRight = style(ballsDots, ansiGreen, format) + " " + style(strikesDots, ansiRed, format) + " - " + style(awayScoreStr, awayScoreStyle, format)
		awayRightVisible = utf8.RuneCountInString(dotsStr) + 3 + len(awayScoreStr)

		outsStr := fmt.Sprintf("%d OUT", game.Linescore.Outs)
		homeRight = style(outsStr, ansiBold, format) + " - " + style(homeScoreStr, homeScoreStyle, format)
		homeRightVisible = len(outsStr) + 3 + len(homeScoreStr)
	} else {
		awayRight = style(awayScoreStr, awayScoreStyle, format)
		awayRightVisible = len(awayScoreStr)
		homeRight = style(homeScoreStr, homeScoreStyle, format)
		homeRightVisible = len(homeScoreStr)
	}

	innerWidth := width - 2

	// Truncate team names if they are too long to fit
	truncLimitAway := innerWidth - 4 - awayRightVisible
	if truncLimitAway < 12 {
		truncLimitAway = 12
	}
	truncLimitHome := innerWidth - 4 - homeRightVisible
	if truncLimitHome < 12 {
		truncLimitHome = 12
	}

	trunc := func(name string, maxLen int) string {
		if utf8.RuneCountInString(name) > maxLen {
			runes := []rune(name)
			return string(runes[:maxLen-2]) + ".."
		}
		return name
	}

	awayName = trunc(awayName, truncLimitAway)
	homeName = trunc(homeName, truncLimitHome)

	// Top border
	border := "+" + strings.Repeat("-", innerWidth) + "+"
	lines = append(lines, style(border, ansiCyan, format))

	// Header line (centered)
	headerPad := innerWidth - len(header)
	headerLeft := headerPad / 2
	headerRight := headerPad - headerLeft
	headerLine := fmt.Sprintf("|%s%s%s|", strings.Repeat(" ", headerLeft), header, strings.Repeat(" ", headerRight))
	lines = append(lines, style(headerLine, headerStyle, format))

	// Away team line
	{
		gap := innerWidth - 3 - utf8.RuneCountInString(awayName) - awayRightVisible
		if gap < 1 {
			gap = 1
		}
		awayLine := fmt.Sprintf("| %s%s %s |", style(html.EscapeString(awayName), awayStyle, format), strings.Repeat(" ", gap), awayRight)
		lines = append(lines, awayLine)
	}

	// Home team line
	{
		gap := innerWidth - 3 - utf8.RuneCountInString(homeName) - homeRightVisible
		if gap < 1 {
			gap = 1
		}
		homeLine := fmt.Sprintf("| %s%s %s |", style(html.EscapeString(homeName), homeStyle, format), strings.Repeat(" ", gap), homeRight)
		lines = append(lines, homeLine)
	}

	// Bottom border
	lines = append(lines, style(border, ansiCyan, format))

	return lines
}

// renderSchedule creates the plain-text scoreboard view
func renderSchedule(sched ScheduleResponse, dateStr string, format string, loc *time.Location) string {
	var sb strings.Builder

	zoneName, _ := time.Now().In(loc).Zone()
	title := fmt.Sprintf("MLB LIVE SCOREBOARD (%s %s)", dateStr, zoneName)
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
		banner.WriteString(txt("           ", format) + style("[MLB]", ansiBold+ansiGreen, format) + txt("             ", format) + style("[NBA]", ansiGray, format) + txt("             ", format) + style("[TENNIS]", ansiGray, format) + "\n")
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
		prevLink := fmt.Sprintf(`<a href="/?date=%s" class="term-link">%s</a>`, prevDateStr, prevLinkText)
		nextLink := fmt.Sprintf(`<a href="/?date=%s" class="term-link">%s</a>`, nextDateStr, nextLinkText)
		banner.WriteString(prevLink + strings.Repeat(" ", spacerSize) + nextLink + "\n")
	} else {
		banner.WriteString(style(prevLinkText, ansiGreen, format) + strings.Repeat(" ", spacerSize) + style(nextLinkText, ansiGreen, format) + "\n")
	}
	banner.WriteString(style("==============================================================================\n", ansiCyan, format))

	sb.WriteString(termPre(format, banner.String()))

	if len(sched.Dates) == 0 || len(sched.Dates[0].Games) == 0 {
		sb.WriteString(termPre(format, txt(" No games scheduled for this date.\n", format)+
			style("==============================================================================\n", ansiCyan, format)))
		return sb.String()
	}

	sort.SliceStable(sched.Dates[0].Games, func(i, j int) bool {
		t1, err1 := time.Parse(time.RFC3339, sched.Dates[0].Games[i].GameDate)
		t2, err2 := time.Parse(time.RFC3339, sched.Dates[0].Games[j].GameDate)
		if err1 == nil && err2 == nil {
			return t1.Before(t2)
		}
		return sched.Dates[0].Games[i].GameDate < sched.Dates[0].Games[j].GameDate
	})

	if format == "html" {
		sb.WriteString(`<div class="games-grid">`)
		for _, game := range sched.Dates[0].Games {
			sb.WriteString(fmt.Sprintf(`<a href="/game/%d" class="term-box-link">`, game.GamePk))
			sb.WriteString(`<pre class="term-art">`)
			boxLines := renderSingleGameBox(game, loc, format, 38)
			for _, line := range boxLines {
				sb.WriteString(line + "\n")
			}
			sb.WriteString(`</pre></a>` + "\n")
		}
		sb.WriteString(`</div>` + "\n")
	} else {
		games := sched.Dates[0].Games
		for i := 0; i < len(games); i += 2 {
			lines1 := renderSingleGameBox(games[i], loc, format, 38)
			var lines2 []string
			if i+1 < len(games) {
				lines2 = renderSingleGameBox(games[i+1], loc, format, 38)
			}

			for j := 0; j < 5; j++ {
				l1 := lines1[j]
				l2 := ""
				if j < len(lines2) {
					l2 = lines2[j]
				}

				if l2 != "" {
					sb.WriteString(l1 + "  " + l2 + "\n")
				} else {
					sb.WriteString(l1 + "\n")
				}
			}
			sb.WriteString("\n")
		}
	}

	if format == "ansi" {
		sb.WriteString(txt(" Run 'curl http://localhost:9090/game/<ID>' to view a game in real-time.\n", format))
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

// renderBattingBoxscore creates the batting boxscore table
func renderBattingBoxscore(players map[string]BoxscorePlayer, batterIDs []int, teamAbb string, format string) string {
	t := NewTable(format,
		TableCol{Title: "PLAYER", Width: 24},
		TableCol{Title: "AB", Align: alignRight},
		TableCol{Title: "R", Align: alignRight},
		TableCol{Title: "H", Align: alignRight},
		TableCol{Title: "RBI", Align: alignRight},
		TableCol{Title: "BB", Align: alignRight},
		TableCol{Title: "SO", Align: alignRight},
		TableCol{Title: "AVG", Align: alignRight},
	)
	t.SetCaption(fmt.Sprintf("%s BATTING", teamAbb))

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
		if utf8.RuneCountInString(displayName) > 24 {
			displayName = string([]rune(displayName)[:23]) + "."
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

		t.AddRow(
			TableCell{Text: displayName},
			TableCell{Text: strconv.Itoa(ab), Align: alignRight},
			TableCell{Text: strconv.Itoa(r), Align: alignRight},
			TableCell{Text: strconv.Itoa(h), Align: alignRight},
			TableCell{Text: strconv.Itoa(rbi), Align: alignRight},
			TableCell{Text: strconv.Itoa(bb), Align: alignRight},
			TableCell{Text: strconv.Itoa(so), Align: alignRight},
			TableCell{Text: avg, Align: alignRight},
		)
	}

	t.AddRow(
		TableCell{Text: "TOTALS", Class: "term-bold", ANSI: ansiBold},
		TableCell{Text: strconv.Itoa(totAB), Align: alignRight, Class: "term-bold", ANSI: ansiBold},
		TableCell{Text: strconv.Itoa(totR), Align: alignRight, Class: "term-bold", ANSI: ansiBold},
		TableCell{Text: strconv.Itoa(totH), Align: alignRight, Class: "term-bold", ANSI: ansiBold},
		TableCell{Text: strconv.Itoa(totRBI), Align: alignRight, Class: "term-bold", ANSI: ansiBold},
		TableCell{Text: strconv.Itoa(totBB), Align: alignRight, Class: "term-bold", ANSI: ansiBold},
		TableCell{Text: strconv.Itoa(totSO), Align: alignRight, Class: "term-bold", ANSI: ansiBold},
		TableCell{Text: "", Align: alignRight, Class: "term-bold", ANSI: ansiBold},
	)

	return t.Render()
}

// renderPitchingBoxscore creates the pitching boxscore table
func renderPitchingBoxscore(players map[string]BoxscorePlayer, pitcherIDs []int, teamAbb string, format string) string {
	t := NewTable(format,
		TableCol{Title: "PLAYER", Width: 24},
		TableCol{Title: "IP", Align: alignRight},
		TableCol{Title: "H", Align: alignRight},
		TableCol{Title: "R", Align: alignRight},
		TableCol{Title: "ER", Align: alignRight},
		TableCol{Title: "BB", Align: alignRight},
		TableCol{Title: "SO", Align: alignRight},
		TableCol{Title: "HR", Align: alignRight},
		TableCol{Title: "ERA", Align: alignRight},
	)
	t.SetCaption(fmt.Sprintf("%s PITCHING", teamAbb))

	for _, id := range pitcherIDs {
		playerKey := "ID" + strconv.Itoa(id)
		player, ok := players[playerKey]
		if !ok {
			continue
		}

		name := player.Person.FullName
		if utf8.RuneCountInString(name) > 24 {
			name = string([]rune(name)[:23]) + "."
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

		t.AddRow(
			TableCell{Text: name},
			TableCell{Text: ip, Align: alignRight},
			TableCell{Text: strconv.Itoa(h), Align: alignRight},
			TableCell{Text: strconv.Itoa(r), Align: alignRight},
			TableCell{Text: strconv.Itoa(er), Align: alignRight},
			TableCell{Text: strconv.Itoa(bb), Align: alignRight},
			TableCell{Text: strconv.Itoa(so), Align: alignRight},
			TableCell{Text: strconv.Itoa(hr), Align: alignRight},
			TableCell{Text: era, Align: alignRight},
		)
	}

	return t.Render()
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
			if desc == "" {
				desc = "unknown result"
			}

			// Select color code matching colorizePitchSequence logic
			code := e.Details.Call.Code
			if code == "" {
				code = e.Details.Type.Code
			}

			var ansiColor string
			switch code {
			case "B":
				ansiColor = ansiGreen
			case "C":
				ansiColor = ansiRed
			case "S":
				ansiColor = ansiMagenta
			case "F":
				ansiColor = ansiYellow
			case "X":
				ansiColor = ansiBlue
			case "D":
				ansiColor = ansiCyan
			case "E":
				ansiColor = ansiBold + ansiCyan
			case "*":
				ansiColor = ansiGray
			case "W":
				ansiColor = ansiBold + ansiRed
			case "H":
				ansiColor = ansiBold + ansiYellow
			case "I":
				ansiColor = ansiBold + ansiGreen
			case "L":
				ansiColor = ansiBold + ansiMagenta
			default:
				ansiColor = ansiReset
			}

			// Format prefix and calculate padding to align outcome column at index 38
			prefix := fmt.Sprintf("  P%d: %d-%d, %s%s", pitchNum, balls, strikes, speedStr, pitchType)
			runeCount := utf8.RuneCountInString(prefix)
			padSize := 38 - runeCount
			if padSize < 2 {
				padSize = 2
			}
			padding := strings.Repeat(" ", padSize)

			line := fmt.Sprintf("%s%s%s", prefix, padding, style(desc, ansiColor, format))
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
		sb.WriteString(p + "\n")
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
func renderGame(game GameFeedResponse, gamePk int, format string) string {
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

	var banner strings.Builder
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
		banner.WriteString(style("========================================================================\n", ansiCyan, format))
		banner.WriteString(titleLine)
		banner.WriteString(subTitleLine)
		banner.WriteString(style("========================================================================\n", ansiCyan, format))
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
		banner.WriteString(style("========================================================================\n", ansiCyan, format))
		banner.WriteString(titleLine)
		banner.WriteString(style(subTitleLine, ansiGray, format))
		banner.WriteString(style("========================================================================\n", ansiCyan, format))
	}
	sb.WriteString(termPre(format, banner.String()))

	var tb strings.Builder
	tb.WriteString(renderDiamondAndMatchup(game, format))
	if pStr := renderCurrentPitches(game, format); pStr != "" {
		tb.WriteString(pStr)
	}
	tb.WriteString("\n")
	sb.WriteString(termPre(format, tb.String()))

	numInnings := 9
	if len(game.LiveData.Linescore.Innings) > 9 {
		numInnings = len(game.LiveData.Linescore.Innings)
	}

	currentInning := game.LiveData.Linescore.CurrentInning
	isLiveGame := state == "In Progress" || state == "Live" || state == "In Progress - Warmup" || state == "Warmup"

	cols := []TableCol{{Title: "TEAM", Width: 10}}
	for i := 1; i <= numInnings; i++ {
		cols = append(cols, TableCol{Title: strconv.Itoa(i), Align: alignRight})
	}
	cols = append(cols,
		TableCol{Title: "R", Align: alignRight},
		TableCol{Title: "H", Align: alignRight},
		TableCol{Title: "E", Align: alignRight},
	)
	lineT := NewTable(format, cols...)
	lineT.SetCaption("LINE SCORE")

	// Away row
	awayCells := []TableCell{{Text: awayAbb}}
	for i := 1; i <= numInnings; i++ {
		val := "-"
		if i-1 < len(game.LiveData.Linescore.Innings) {
			inn := game.LiveData.Linescore.Innings[i-1]
			if inn.Away.Runs != nil {
				val = strconv.Itoa(*inn.Away.Runs)
			}
		}
		cell := TableCell{Text: val, Align: alignRight}
		if isLiveGame && i == currentInning && game.LiveData.Linescore.IsTopInning {
			cell.ANSI = ansiBold + ansiYellow
		}
		awayCells = append(awayCells, cell)
	}
	awayR, awayH, awayE := "-", "-", "-"
	if game.LiveData.Linescore.Teams.Away.Runs != nil {
		awayR = strconv.Itoa(*game.LiveData.Linescore.Teams.Away.Runs)
	}
	if game.LiveData.Linescore.Teams.Away.Hits != nil {
		awayH = strconv.Itoa(*game.LiveData.Linescore.Teams.Away.Hits)
	}
	if game.LiveData.Linescore.Teams.Away.Errors != nil {
		awayE = strconv.Itoa(*game.LiveData.Linescore.Teams.Away.Errors)
	}
	awayCells = append(awayCells,
		TableCell{Text: awayR, Align: alignRight},
		TableCell{Text: awayH, Align: alignRight},
		TableCell{Text: awayE, Align: alignRight},
	)
	lineT.AddRow(awayCells...)

	// Home row
	homeCells := []TableCell{{Text: homeAbb}}
	for i := 1; i <= numInnings; i++ {
		val := "-"
		if i-1 < len(game.LiveData.Linescore.Innings) {
			inn := game.LiveData.Linescore.Innings[i-1]
			if inn.Home.Runs != nil {
				val = strconv.Itoa(*inn.Home.Runs)
			}
		}
		cell := TableCell{Text: val, Align: alignRight}
		if isLiveGame && i == currentInning && !game.LiveData.Linescore.IsTopInning {
			cell.ANSI = ansiBold + ansiYellow
		}
		homeCells = append(homeCells, cell)
	}
	homeR, homeH, homeE := "-", "-", "-"
	if game.LiveData.Linescore.Teams.Home.Runs != nil {
		homeR = strconv.Itoa(*game.LiveData.Linescore.Teams.Home.Runs)
	}
	if game.LiveData.Linescore.Teams.Home.Hits != nil {
		homeH = strconv.Itoa(*game.LiveData.Linescore.Teams.Home.Hits)
	}
	if game.LiveData.Linescore.Teams.Home.Errors != nil {
		homeE = strconv.Itoa(*game.LiveData.Linescore.Teams.Home.Errors)
	}
	homeCells = append(homeCells,
		TableCell{Text: homeR, Align: alignRight},
		TableCell{Text: homeH, Align: alignRight},
		TableCell{Text: homeE, Align: alignRight},
	)
	lineT.AddRow(homeCells...)
	sb.WriteString(lineT.Render())

	// Recent plays
	var plays strings.Builder
	plays.WriteString(style(" RECENT PLAYS:\n", ansiBold+ansiCyan, format))
	allPlays := game.LiveData.Plays.AllPlays

	var lastInning int = -1
	var lastIsTop bool = false
	var hasLastSeen bool = false

	for i := len(allPlays) - 1; i >= 0; i-- {
		play := allPlays[i]
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
				plays.WriteString(style(header, ansiBold+ansiCyan, format))
				lastInning = inning
				lastIsTop = isTop
				hasLastSeen = true
			}

			hasRun := false
			if i > 0 {
				prevPlay := allPlays[i-1]
				if play.Result.AwayScore != prevPlay.Result.AwayScore || play.Result.HomeScore != prevPlay.Result.HomeScore {
					hasRun = true
				}
			}

			isOut := false
			lowerDesc := strings.ToLower(desc)
			if strings.Contains(lowerDesc, "out") ||
				strings.Contains(lowerDesc, "strike out") ||
				strings.Contains(lowerDesc, "grounded into") {
				isOut = true
			}

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

			isLast := i == len(allPlays)-1
			halfCode := "B"
			if isTop {
				halfCode = "T"
			}
			halfCode = fmt.Sprintf("%s%d", halfCode, inning)

			prefix := fmt.Sprintf(" [%s] ", halfCode)

			var playLine string
			if isOut && isLast {
				playLine = style(prefix+desc, ansiRed, format) + seqStyled + style(" (Current Play)\n", ansiRed, format)
			} else if isOut {
				playLine = style(prefix+desc, ansiRed, format) + seqStyled + txt("\n", format)
			} else if isLast {
				playLine = style(prefix+desc, ansiGreen, format) + seqStyled + style(" (Current Play)\n", ansiGreen, format)
			} else if hasRun {
				playLine = style(prefix+desc, ansiBold, format) + seqStyled + txt("\n", format)
			} else {
				playLine = txt(prefix+desc, format) + seqStyled + txt("\n", format)
			}
			plays.WriteString(playLine)

			scoreLine := fmt.Sprintf(" Score: %d-%d\n", play.Result.AwayScore, play.Result.HomeScore)
			plays.WriteString(style(scoreLine, ansiGray, format))
		}
	}

	if len(allPlays) == 0 {
		plays.WriteString(txt(" No plays recorded yet.\n", format))
	}
	sb.WriteString(termPre(format, plays.String()))

	// Boxscore Statistics Section
	var boxHead strings.Builder
	boxHead.WriteString("\n")
	boxHead.WriteString(style("========================================================================\n", ansiCyan, format))
	boxHead.WriteString(style("                          BOXSCORE STATISTICS\n", ansiBold+ansiCyan, format))
	boxHead.WriteString(style("========================================================================\n", ansiCyan, format))
	sb.WriteString(termPre(format, boxHead.String()))

	// Away batting
	sb.WriteString(renderBattingBoxscore(game.LiveData.Boxscore.Teams.Away.Players, game.LiveData.Boxscore.Teams.Away.Batters, awayAbb, format))
	// Home batting
	sb.WriteString(renderBattingBoxscore(game.LiveData.Boxscore.Teams.Home.Players, game.LiveData.Boxscore.Teams.Home.Batters, homeAbb, format))

	sb.WriteString(termPre(format, style("\n------------------------------------------------------------------------\n", ansiCyan, format)))

	// Away pitching
	sb.WriteString(renderPitchingBoxscore(game.LiveData.Boxscore.Teams.Away.Players, game.LiveData.Boxscore.Teams.Away.Pitchers, awayAbb, format))
	// Home pitching
	sb.WriteString(renderPitchingBoxscore(game.LiveData.Boxscore.Teams.Home.Players, game.LiveData.Boxscore.Teams.Home.Pitchers, homeAbb, format))

	// Pitch legend
	var legend strings.Builder
	legend.WriteString(style("========================================================================\n", ansiCyan, format))
	legend.WriteString(style(" PITCH LEGEND:\n", ansiBold+ansiCyan, format))
	legend.WriteString(txt("  ", format) + style("B", ansiGreen, format) + txt(": Ball               ", format) + style("C", ansiRed, format) + txt(": Called Strike      ", format) + style("S", ansiMagenta, format) + txt(": Swinging Strike\n", format))
	legend.WriteString(txt("  ", format) + style("F", ansiYellow, format) + txt(": Foul               ", format) + style("X", ansiBlue, format) + txt(": In Play, Out       ", format) + style("D", ansiCyan, format) + txt(": In Play, No Out (Hit)\n", format))
	legend.WriteString(txt("  ", format) + style("E", ansiBold+ansiCyan, format) + txt(": In Play, Run(s)    ", format) + style("*", ansiGray, format) + txt(": Ball in Dirt       ", format) + style("W", ansiBold+ansiRed, format) + txt(": Swinging Strike (Pitchout)\n", format))
	legend.WriteString(txt("  ", format) + style("H", ansiBold+ansiYellow, format) + txt(": Hit By Pitch       ", format) + style("I", ansiBold+ansiGreen, format) + txt(": Intentional Ball   ", format) + style("L", ansiBold+ansiMagenta, format) + txt(": Foul Tip\n", format))
	legend.WriteString(style("========================================================================\n", ansiCyan, format))
	sb.WriteString(termPre(format, legend.String()))

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

	url := fmt.Sprintf("https://statsapi.mlb.com/api/v1/schedule?sportId=1&date=%s&hydrate=linescore", dateStr)
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

	// Parse gamePk to int for betting odds
	gamePkInt, _ := strconv.Atoi(gamePk)
	text := renderGame(game, gamePkInt, format)
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
	DivisionId   int    `json:"-"`
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
	padding := (layoutWidth - len(title)) / 2
	if padding < 0 {
		padding = 0
	}

	var banner strings.Builder
	if format != "html" {
		banner.WriteString(style("==============================================================================\n", ansiCyan, format))
		banner.WriteString(txt(" [SCOREBOARD]   [COMPARE]   ", format) + style("[MLB]", ansiBold+ansiGreen, format) + txt("             ", format) + style("[NBA]", ansiGray, format) + "\n")
		banner.WriteString(style("==============================================================================\n", ansiCyan, format))
	}
	banner.WriteString(txt(strings.Repeat(" ", padding), format))
	banner.WriteString(style(title+"\n", ansiBold+ansiCyan, format))
	banner.WriteString(style("==============================================================================\n", ansiCyan, format))
	sb.WriteString(termPre(format, banner.String()))

	// Collect all teams for unified ranking
	var allTeams []DivisionTeam

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
				DivisionId:   rec.Division.Id,
			}
			allTeams = append(allTeams, t)
		}
	}

	// Sort all teams by winning percentage (string comparison works for MLB PCT format like ".567")
	sort.SliceStable(allTeams, func(i, j int) bool {
		return allTeams[i].PCT > allTeams[j].PCT
	})

	trunc20 := func(s string) string {
		if utf8.RuneCountInString(s) > 20 {
			return string([]rune(s)[:19]) + "."
		}
		return s
	}
	abbrCell := func(t DivisionTeam) TableCell {
		return TableCell{Text: t.Abbreviation, Link: fmt.Sprintf("/mlb/team/%d", t.Id)}
	}

	// Render unified all-teams table
	allT := NewTable(format,
		TableCol{Title: "RANK", Width: 4, Align: alignRight},
		TableCol{Title: "TEAM", Width: 4},
		TableCol{Title: "NAME", Width: 20},
		TableCol{Title: "W", Align: alignRight},
		TableCol{Title: "L", Align: alignRight},
		TableCol{Title: "PCT", Align: alignRight},
		TableCol{Title: "GB", Align: alignRight},
		TableCol{Title: "LEAGUE"},
	)
	allT.SetCaption("ALL TEAMS RANKINGS")
	for i, t := range allTeams {
		allT.AddRow(
			TableCell{Text: strconv.Itoa(i + 1), Align: alignRight},
			abbrCell(t),
			TableCell{Text: trunc20(t.Name)},
			TableCell{Text: strconv.Itoa(t.Wins), Align: alignRight},
			TableCell{Text: strconv.Itoa(t.Losses), Align: alignRight},
			TableCell{Text: t.PCT, Align: alignRight},
			TableCell{Text: t.GB, Align: alignRight},
			TableCell{Text: t.League},
		)
	}
	sb.WriteString(allT.Render())

	// Group standings by division + wild card
	alEast := []DivisionTeam{}
	alCentral := []DivisionTeam{}
	alWest := []DivisionTeam{}
	nlEast := []DivisionTeam{}
	nlCentral := []DivisionTeam{}
	nlWest := []DivisionTeam{}
	alWildCard := []DivisionTeam{}
	nlWildCard := []DivisionTeam{}

	for _, t := range allTeams {
		if t.WildCard && t.DivisionRank != "1" {
			if t.League == "AL" {
				alWildCard = append(alWildCard, t)
			} else {
				nlWildCard = append(nlWildCard, t)
			}
			continue
		}
		switch t.DivisionId {
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

	renderDivision := func(teams []DivisionTeam, divName string) {
		if len(teams) == 0 {
			return
		}
		divT := NewTable(format,
			TableCol{Title: "TEAM", Width: 4},
			TableCol{Title: "NAME", Width: 20},
			TableCol{Title: "W", Align: alignRight},
			TableCol{Title: "L", Align: alignRight},
			TableCol{Title: "PCT", Align: alignRight},
			TableCol{Title: "GB", Align: alignRight},
		)
		divT.SetCaption(strings.ToUpper(divName))
		for i, t := range teams {
			rowClass := ""
			rowANSI := ""
			if i == 0 || t.GB == "-" || t.GB == "0" {
				rowClass = "term-bold"
				rowANSI = ansiBold
			}
			cells := []TableCell{
				func() TableCell { c := abbrCell(t); c.Class, c.ANSI = rowClass, rowANSI; return c }(),
				func() TableCell { c := TableCell{Text: trunc20(t.Name), Class: rowClass, ANSI: rowANSI}; return c }(),
				TableCell{Text: strconv.Itoa(t.Wins), Align: alignRight, Class: rowClass, ANSI: rowANSI},
				TableCell{Text: strconv.Itoa(t.Losses), Align: alignRight, Class: rowClass, ANSI: rowANSI},
				TableCell{Text: t.PCT, Align: alignRight, Class: rowClass, ANSI: rowANSI},
				TableCell{Text: t.GB, Align: alignRight, Class: rowClass, ANSI: rowANSI},
			}
			divT.AddRow(cells...)
		}
		sb.WriteString(divT.Render())
	}

	renderWildCard := func(teams []DivisionTeam, label string) {
		if len(teams) == 0 {
			return
		}
		wcT := NewTable(format,
			TableCol{Title: "TEAM", Width: 4},
			TableCol{Title: "NAME", Width: 20},
			TableCol{Title: "W", Align: alignRight},
			TableCol{Title: "L", Align: alignRight},
			TableCol{Title: "PCT", Align: alignRight},
			TableCol{Title: "GB", Align: alignRight},
		)
		wcT.SetCaption(strings.ToUpper(label) + " WILD CARD")
		for _, t := range teams {
			wcT.AddRow(
				func() TableCell { c := abbrCell(t); c.Class, c.ANSI = "term-green", ansiGreen; return c }(),
				TableCell{Text: trunc20(t.Name), Class: "term-green", ANSI: ansiGreen},
				TableCell{Text: strconv.Itoa(t.Wins), Align: alignRight, Class: "term-green", ANSI: ansiGreen},
				TableCell{Text: strconv.Itoa(t.Losses), Align: alignRight, Class: "term-green", ANSI: ansiGreen},
				TableCell{Text: t.PCT, Align: alignRight, Class: "term-green", ANSI: ansiGreen},
				TableCell{Text: t.GB, Align: alignRight, Class: "term-green", ANSI: ansiGreen},
			)
		}
		sb.WriteString(wcT.Render())
	}

	renderDivision(alEast, "AL EAST")
	renderDivision(alCentral, "AL CENTRAL")
	renderDivision(alWest, "AL WEST")
	renderWildCard(alWildCard, "AL")
	renderDivision(nlEast, "NL EAST")
	renderDivision(nlCentral, "NL CENTRAL")
	renderDivision(nlWest, "NL WEST")
	renderWildCard(nlWildCard, "NL")

	sb.WriteString(termPre(format, style("==============================================================================\n", ansiCyan, format)))
	if format == "ansi" {
		sb.WriteString(txt(" Run 'curl http://localhost:9090/' to return to the scoreboard.\n", format))
		sb.WriteString(style("==============================================================================\n", ansiCyan, format))
	}

	return sb.String()
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
	padding := (layoutWidth - len(title)) / 2
	if padding < 0 {
		padding = 0
	}

	var banner strings.Builder
	if format != "html" {
		banner.WriteString(style("==============================================================================\n", ansiCyan, format))
		banner.WriteString(txt(" [SCOREBOARD]   [STANDINGS]   [COMPARE]   ", format) + style("[MLB]", ansiBold+ansiGreen, format) + txt("             ", format) + style("[NBA]", ansiGray, format) + "\n")
		banner.WriteString(style("==============================================================================\n", ansiCyan, format))
	}
	banner.WriteString(txt(strings.Repeat(" ", padding), format))
	banner.WriteString(style(title+"\n", ansiBold+ansiCyan, format))
	banner.WriteString(style("==============================================================================\n", ansiCyan, format))
	sb.WriteString(termPre(format, banner.String()))

	// Team information section
	infoT := NewTable(format,
		TableCol{Title: "FIELD", Width: 16},
		TableCol{Title: "VALUE"},
	)
	infoT.SetCaption("TEAM INFORMATION")
	infoT.AddRow(TableCell{Text: "Name:"}, TableCell{Text: teamName})
	infoT.AddRow(TableCell{Text: "City:"}, TableCell{Text: teamCity})
	infoT.AddRow(TableCell{Text: "Abbreviation:"}, TableCell{Text: teamAbb})
	infoT.AddRow(TableCell{Text: "League:"}, TableCell{Text: teamLeague})
	infoT.AddRow(TableCell{Text: "Division:"}, TableCell{Text: teamDivision})
	sb.WriteString(infoT.Render())

	// Fetch team season stats from MLB Stats API (dynamic year)
	statsUrl := fmt.Sprintf("https://statsapi.mlb.com/api/v1/teams/%d/stats?season=%d&gameType=R&stats=season&group=hitting,pitching,fielding", teamId, seasonYear)
	resp, err := client.Get(statsUrl)
	if err != nil {
		sb.WriteString(termPre(format, style("\n WARNING: Could not fetch team stats.\n", ansiYellow, format)))
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
			sb.WriteString(termPre(format, style("\n WARNING: Could not decode team stats.\n", ansiYellow, format)))
		} else {
			getStr := func(m map[string]interface{}, key string) string {
				if v, ok := m[key]; ok {
					return fmt.Sprintf("%v", v)
				}
				return "-"
			}

			valCell := func(val string, highlight bool) TableCell {
				c := TableCell{Text: val, Align: alignRight}
				if highlight && val != "-" {
					c.Class = "term-bold term-green term-highlight-better"
					c.ANSI = ansiBold + ansiGreen
				}
				return c
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

					t := NewTable(format,
						TableCol{Title: "GP", Width: 5, Align: alignRight},
						TableCol{Title: "AVG", Width: 5, Align: alignRight},
						TableCol{Title: "OBP", Width: 5, Align: alignRight},
						TableCol{Title: "SLG", Width: 5, Align: alignRight},
						TableCol{Title: "OPS", Width: 5, Align: alignRight},
						TableCol{Title: "R", Width: 5, Align: alignRight},
						TableCol{Title: "H", Width: 5, Align: alignRight},
						TableCol{Title: "HR", Width: 4, Align: alignRight},
						TableCol{Title: "RBI", Width: 4, Align: alignRight},
						TableCol{Title: "BB", Width: 4, Align: alignRight},
						TableCol{Title: "SO", Width: 4, Align: alignRight},
						TableCol{Title: "SB", Width: 4, Align: alignRight},
					)
					t.SetCaption("BATTING")
					t.AddRow(
						valCell(gp, false),
						valCell(avg, isAboveAvg(avg, 0.243)),
						valCell(obp, isAboveAvg(obp, 0.319)),
						valCell(slg, isAboveAvg(slg, 0.400)),
						valCell(ops, isAboveAvg(ops, 0.719)),
						valCell(r, false),
						valCell(h, false),
						valCell(hr, false),
						valCell(rbi, false),
						valCell(bb, false),
						valCell(so, false),
						valCell(stl, false),
					)
					sb.WriteString(t.Render())
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

					t := NewTable(format,
						TableCol{Title: "W", Width: 4, Align: alignRight},
						TableCol{Title: "L", Width: 4, Align: alignRight},
						TableCol{Title: "ERA", Width: 5, Align: alignRight},
						TableCol{Title: "WHIP", Width: 5, Align: alignRight},
						TableCol{Title: "IP", Width: 6, Align: alignRight},
						TableCol{Title: "SO", Width: 4, Align: alignRight},
						TableCol{Title: "BB", Width: 4, Align: alignRight},
						TableCol{Title: "HR", Width: 4, Align: alignRight},
						TableCol{Title: "SV", Width: 4, Align: alignRight},
						TableCol{Title: "HLD", Width: 4, Align: alignRight},
						TableCol{Title: "BS", Width: 4, Align: alignRight},
						TableCol{Title: "AVG", Width: 5, Align: alignRight},
					)
					t.SetCaption("PITCHING")
					t.AddRow(
						valCell(w, false),
						valCell(l, false),
						valCell(era, isBelowAvg(era, 4.18)),
						valCell(whip, isBelowAvg(whip, 1.308)),
						valCell(ip, false),
						valCell(so, false),
						valCell(bb, false),
						valCell(hr, false),
						valCell(sv, false),
						valCell(hld, false),
						valCell(bs, false),
						valCell(avg, isBelowAvg(avg, 0.243)),
					)
					sb.WriteString(t.Render())
				} else if gn == "fielding" {
					fpct := getStr(s, "fielding")
					e := getStr(s, "errors")
					dp := getStr(s, "doublePlays")
					pb := getStr(s, "passedBall")

					t := NewTable(format,
						TableCol{Title: "FPCT", Width: 5, Align: alignRight},
						TableCol{Title: "E", Width: 5, Align: alignRight},
						TableCol{Title: "DP", Width: 5, Align: alignRight},
						TableCol{Title: "PB", Width: 5, Align: alignRight},
					)
					t.SetCaption("FIELDING")
					t.AddRow(
						valCell(fpct, isAboveAvg(fpct, 0.985)),
						valCell(e, false),
						valCell(dp, false),
						valCell(pb, false),
					)
					sb.WriteString(t.Render())
				}
			}
		}
	}

	// Fetch and display recent games (last 30 days)
	today := time.Now()
	startDate := today.AddDate(0, 0, -30).Format("2006-01-02")
	endDate := today.Format("2006-01-02")

	games, err := fetchTeamGames(teamId, startDate, endDate)
	if err != nil {
		sb.WriteString(termPre(format, style("\n WARNING: Could not fetch recent games.\n", ansiYellow, format)))
	} else if len(games) == 0 {
		sb.WriteString(termPre(format, txt("\n No games found in the last 30 days (off-season).\n", format)))
	} else {
		// Reverse to show most recent first
		for i, j := 0, len(games)-1; i < j; i, j = i+1, j-1 {
			games[i], games[j] = games[j], games[i]
		}

		rgT := NewTable(format,
			TableCol{Title: "DATE", Width: 10},
			TableCol{Title: "OPPONENT", Width: 24},
			TableCol{Title: "W/L", Width: 4, Align: alignRight},
			TableCol{Title: "SCORE", Width: 5, Align: alignRight},
		)
		rgT.SetCaption(fmt.Sprintf("RECENT GAMES (Last 30 Days) - %d Games", len(games)))
		for _, g := range games {
			homeStr := "(H) "
			if !g.IsHome {
				homeStr = "(A) "
			}
			opponent := g.Opponent
			if utf8.RuneCountInString(opponent) > 22 {
				opponent = string([]rune(opponent)[:19]) + "..."
			}
			oppDisplay := homeStr + opponent

			var wlText, wlClass, wlANSI string
			switch g.Result {
			case "W":
				wlText, wlClass, wlANSI = "W", "term-green", ansiGreen
			case "L":
				wlText, wlClass, wlANSI = "L", "term-red", ansiRed
			case "T":
				wlText, wlClass, wlANSI = "T", "term-yellow", ansiYellow
			default:
				if g.State == "In Progress" || g.State == "Live" {
					wlText, wlClass, wlANSI = "LIVE", "term-green", ansiGreen
				} else {
					wlText, wlClass, wlANSI = "", "term-gray", ansiGray
				}
			}

			scoreText := g.Score
			if scoreText == "-" {
				scoreText = "UPCOMING"
			}
			rgT.AddRow(
				TableCell{Text: g.Date},
				TableCell{Text: oppDisplay},
				TableCell{Text: wlText, Align: alignRight, Class: wlClass, ANSI: wlANSI},
				TableCell{Text: scoreText, Align: alignRight, Link: fmt.Sprintf("/game/%d", g.GamePk)},
			)
		}
		sb.WriteString(rgT.Render())
	}

	sb.WriteString(termPre(format, style("==============================================================================\n", ansiCyan, format)))
	if format == "ansi" {
		sb.WriteString(txt(" Run 'curl http://localhost:9090/' to return to the scoreboard.\n", format))
		sb.WriteString(style("==============================================================================\n", ansiCyan, format))
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

func statRow(label, key, group string, stats1, stats2 map[string]map[string]string, lowerIsBetter bool, format string) []TableCell {
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

	mk := func(v string, better bool) TableCell {
		c := TableCell{Text: v, Align: alignRight}
		if better {
			c.Class = "term-bold term-green term-highlight-better"
			c.ANSI = ansiBold + ansiGreen
		}
		return c
	}

	return []TableCell{
		TableCell{Text: label},
		mk(val1, better1),
		mk(val2, better2),
	}
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

	padding := (layoutWidth - len(title)) / 2
	if padding < 0 {
		padding = 0
	}

	var banner strings.Builder
	if format != "html" {
		banner.WriteString(style("==============================================================================\n", ansiCyan, format))
		banner.WriteString(txt(" [SCOREBOARD]   [STANDINGS]   [COMPARE]   ", format) + style("[MLB]", ansiBold+ansiGreen, format) + txt("             ", format) + style("[NBA]", ansiGray, format) + "\n")
		banner.WriteString(style("==============================================================================\n", ansiCyan, format))
	}
	banner.WriteString(txt(strings.Repeat(" ", padding), format))
	banner.WriteString(style(title+"\n", ansiBold+ansiCyan, format))
	banner.WriteString(style("==============================================================================\n", ansiCyan, format))
	sb.WriteString(termPre(format, banner.String()))

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
		sb.WriteString(style("==============================================================================\n", ansiCyan, format))
	} else {
		// In curl/ANSI format, list options or show command help
		if teamA == nil || teamB == nil {
			var help strings.Builder
			help.WriteString(style("\n HOW TO COMPARE TEAMS:\n", ansiBold+ansiCyan, format))
			help.WriteString(txt(" Specify team IDs in query parameters: ?team1=<id1>&team2=<id2>\n", format))
			help.WriteString(txt(" Example: curl \"http://localhost:9090/mlb/compare?team1=130&team2=147\"\n\n", format))
			help.WriteString(style(" AVAILABLE MLB TEAMS:\n", ansiBold+ansiCyan, format))
			help.WriteString(style("------------------------------------------------------------------------------\n", ansiCyan, format))
			for i, t := range allTeams {
				help.WriteString(txt(fmt.Sprintf("  %-4d: %-30s (%s)", t.Id, t.Name, t.Abbreviation), format))
				if (i+1)%2 == 0 {
					help.WriteString("\n")
				} else {
					help.WriteString("   | ")
				}
			}
			if len(allTeams)%2 != 0 {
				help.WriteString("\n")
			}
			help.WriteString(style("------------------------------------------------------------------------------\n", ansiCyan, format))
			sb.WriteString(termPre(format, help.String()))
			return sb.String()
		}
	}

	if teamA == nil || teamB == nil {
		if format == "html" {
			sb.WriteString(style("\n Select two teams above to compare their season statistics.\n\n", ansiYellow, format))
			sb.WriteString(style("==============================================================================\n", ansiCyan, format))
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
	infoT := NewTable(format,
		TableCol{Title: "STAT", Width: 21},
		TableCol{Title: teamA.Abbreviation, Width: 27},
		TableCol{Title: teamB.Abbreviation, Width: 26},
	)
	infoT.SetCaption("TEAM INFORMATION")

	leagueA := "AL"
	if strings.Contains(strings.ToLower(teamA.League.Name), "national") {
		leagueA = "NL"
	}
	leagueB := "AL"
	if strings.Contains(strings.ToLower(teamB.League.Name), "national") {
		leagueB = "NL"
	}
	infoT.AddRow(TableCell{Text: "Team Name:"}, TableCell{Text: teamA.Name}, TableCell{Text: teamB.Name})
	infoT.AddRow(TableCell{Text: "City:"}, TableCell{Text: teamA.LocationName}, TableCell{Text: teamB.LocationName})
	infoT.AddRow(TableCell{Text: "League:"}, TableCell{Text: leagueA}, TableCell{Text: leagueB})
	infoT.AddRow(TableCell{Text: "Division:"}, TableCell{Text: teamA.Division.Name}, TableCell{Text: teamB.Division.Name})
	sb.WriteString(infoT.Render())

	// 3. Batting Stats
	batT := NewTable(format,
		TableCol{Title: "STAT", Width: 25},
		TableCol{Title: teamA.Abbreviation, Width: 27, Align: alignRight},
		TableCol{Title: teamB.Abbreviation, Width: 26, Align: alignRight},
	)
	batT.SetCaption("BATTING STATS")
	for _, f := range battingFields {
		batT.AddRow(statRow(f.label, f.key, "hitting", stats1, stats2, f.lowerIsBetter, format)...)
	}
	sb.WriteString(batT.Render())

	// 4. Pitching Stats
	pitT := NewTable(format,
		TableCol{Title: "STAT", Width: 25},
		TableCol{Title: teamA.Abbreviation, Width: 27, Align: alignRight},
		TableCol{Title: teamB.Abbreviation, Width: 26, Align: alignRight},
	)
	pitT.SetCaption("PITCHING STATS")
	for _, f := range pitchingFields {
		pitT.AddRow(statRow(f.label, f.key, "pitching", stats1, stats2, f.lowerIsBetter, format)...)
	}
	sb.WriteString(pitT.Render())

	// 5. Fielding Stats
	fldT := NewTable(format,
		TableCol{Title: "STAT", Width: 25},
		TableCol{Title: teamA.Abbreviation, Width: 27, Align: alignRight},
		TableCol{Title: teamB.Abbreviation, Width: 26, Align: alignRight},
	)
	fldT.SetCaption("FIELDING STATS")
	for _, f := range fieldingFields {
		fldT.AddRow(statRow(f.label, f.key, "fielding", stats1, stats2, f.lowerIsBetter, format)...)
	}
	sb.WriteString(fldT.Render())

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

	var rg strings.Builder
	rg.WriteString(style("\n RECENT GAMES (Last 30 Days)\n", ansiBold+ansiCyan, format))
	rg.WriteString(style("------------------------------------------------------------------------------\n", ansiCyan, format))
	headerA := fmt.Sprintf(" %s RECENT GAMES", teamA.Abbreviation)
	headerB := fmt.Sprintf(" %s RECENT GAMES", teamB.Abbreviation)
	rg.WriteString(style(pad(headerA, 38, true, format)+"|"+pad(headerB, 39, true, format)+"\n", ansiBold, format))
	rg.WriteString(style("------------------------------------------------------------------------------\n", ansiCyan, format))

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

		rg.WriteString(pad(colA, 38, true, format) + "|" + pad(colB, 39, true, format) + "\n")
	}
	rg.WriteString(style("------------------------------------------------------------------------------\n", ansiCyan, format))
	sb.WriteString(termPre(format, rg.String()))

	sb.WriteString(termPre(format, style("==============================================================================\n", ansiCyan, format)))
	if format == "ansi" {
		sb.WriteString(txt(" Run 'curl http://localhost:9090/' to return to the scoreboard.\n", format))
		sb.WriteString(style("==============================================================================\n", ansiCyan, format))
	}

	return sb.String()
}


