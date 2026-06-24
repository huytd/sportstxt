package main

import (
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// ANSI terminal style escape codes
const (
	ansiReset  = "\033[0m"
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

// renderSchedule creates the plain-text scoreboard view
func renderSchedule(sched ScheduleResponse, dateStr string, format string) string {
	var sb strings.Builder

	title := fmt.Sprintf("MLB LIVE SCOREBOARD (%s)", dateStr)
	padding := (72 - len(title)) / 2
	if padding < 0 {
		padding = 0
	}
	sb.WriteString(style("========================================================================\n", ansiCyan, format))
	sb.WriteString(txt(strings.Repeat(" ", padding), format))
	sb.WriteString(style(title+"\n", ansiBold+ansiCyan, format))
	sb.WriteString(style("========================================================================\n", ansiCyan, format))
	sb.WriteString(style(fmt.Sprintf(" %-8s %-22s %2s  @  %-22s %2s  %-11s\n", "ID", "AWAY TEAM", "R", "HOME TEAM", "R", "STATUS"), ansiBold, format))
	sb.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))

	if len(sched.Dates) == 0 || len(sched.Dates[0].Games) == 0 {
		sb.WriteString(txt(fmt.Sprintf(" %s\n", "No games scheduled for today."), format))
		sb.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
		sb.WriteString(txt(" Use ?date=YYYY-MM-DD to query another date.\n", format))
		sb.WriteString(style("========================================================================\n", ansiCyan, format))
		return sb.String()
	}

	for _, game := range sched.Dates[0].Games {
		idStr := strconv.Itoa(game.GamePk)
		awayName := game.Teams.Away.Team.Name
		homeName := game.Teams.Home.Team.Name

		if len(awayName) > 22 {
			awayName = awayName[:21] + "."
		}
		if len(homeName) > 22 {
			homeName = homeName[:21] + "."
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

		row := fmt.Sprintf(" %-8s %-22s %2s  @  %-22s %2s  %-11s\n",
			idStr,
			awayName,
			awayScoreStr,
			homeName,
			homeScoreStr,
			statusStr,
		)
		sb.WriteString(style(row, rowStyle, format))
	}

	sb.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))
	if format == "ansi" {
		sb.WriteString(txt(" Run 'curl http://localhost:8080/game/<ID>' to view a game in real-time.\n", format))
	} else {
		sb.WriteString(txt(" Click on a game ID to view the game in real-time.\n", format))
	}
	sb.WriteString(style("========================================================================\n", ansiCyan, format))

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
		style("[COUNT & OUTS]", ansiBold+ansiCyan, format),
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
	sb.WriteString(style(" Pitches:\n", ansiBold+ansiCyan, format))
	for _, p := range pitches {
		sb.WriteString(txt(p+"\n", format))
	}
	return sb.String()
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

	sb.WriteString(style("------------------------------------------------------------------------\n", ansiCyan, format))

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
		sb.WriteString(txt(" Run 'curl http://localhost:8080/' to return to the scoreboard.\n", format))
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

	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
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
		text := renderSchedule(sched, dateStr, "html")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(text))
	} else if isCurl {
		text := renderSchedule(sched, dateStr, "ansi")
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
			sb.WriteString(txt(" Run 'curl http://localhost:8080/' to return to the scoreboard.\n", format))
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
	dateStr := r.URL.Query().Get("date")
	if dateStr == "" {
		dateStr = time.Now().Format("2006-01-02")
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

func main() {
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	// Go 1.22 path pattern routing
	http.HandleFunc("GET /", handleSchedule)
	http.HandleFunc("GET /game/{gamePk}", handleGame)
	http.HandleFunc("GET /api/games", handleAPIGames)
	http.HandleFunc("GET /api/game/{gamePk}", handleAPIGameDetail)

	log.Printf("Starting sportstxt MLB tracker on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}

const htmlPage = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>sportstxt - MLB Scoreboard</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;700&display=swap" rel="stylesheet">
    <style>
        :root {
            --term-bg: #07080c;
            --term-green: #39ff14;
            --term-yellow: #ffeb3b;
            --term-red: #ff3b30;
            --term-cyan: #00f0ff;
            --term-blue: #3b82f6;
            --term-magenta: #d946ef;
            --term-gray: #555866;
            --color-primary: var(--term-green);
        }

        body {
            background-color: var(--term-bg);
            color: var(--color-primary);
            font-family: 'JetBrains Mono', monospace;
            margin: 0;
            padding: 20px;
            display: flex;
            flex-direction: column;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            box-sizing: border-box;
            overflow-x: hidden;
        }

        .term-container {
            width: 100%;
            max-width: 900px;
            background: #07080c;
            border: 1px solid #1f222e;
            border-radius: 8px;
            padding: 30px;
            box-sizing: border-box;
            position: relative;
        }

        pre {
            margin: 0;
            font-family: inherit;
            font-size: 15px;
            line-height: 1.6;
            white-space: pre-wrap;
            position: relative;
            overflow-x: auto;
        }

        /* Colored spans rendered by Go backend */
        .term-green { color: var(--term-green); }
        .term-yellow { color: var(--term-yellow); }
        .term-red { color: var(--term-red); }
        .term-cyan { color: var(--term-cyan); }
        .term-blue { color: var(--term-blue); }
        .term-magenta { color: var(--term-magenta); }
        .term-gray { color: var(--term-gray); }
        .term-bold { font-weight: bold; }

        /* Links styling */
        .term-link {
            color: var(--color-primary);
            text-decoration: underline;
            cursor: pointer;
            font-weight: bold;
        }
        .term-link:hover {
            color: #ffffff !important;
        }

        /* Status bar */
        .status-bar {
            display: flex;
            justify-content: space-between;
            margin-bottom: 15px;
            border-bottom: 1px solid #1f222e;
            padding-bottom: 10px;
            font-size: 12px;
            color: #8a8f98;
        }

        .status-bar span {
            display: flex;
            align-items: center;
            gap: 6px;
        }

        .dot-pulse {
            width: 8px;
            height: 8px;
            background-color: var(--term-green);
            border-radius: 50%;
            display: inline-block;
            animation: pulse 1.5s infinite;
        }

        @keyframes pulse {
            0% { transform: scale(0.95); opacity: 0.5; }
            70% { transform: scale(1); opacity: 1; }
            100% { transform: scale(0.95); opacity: 0.5; }
        }
    </style>
</head>
<body>
    <div class="term-container">
        <div class="status-bar">
            <span id="status-left"><span class="dot-pulse"></span>LIVE FEED</span>
            <span id="status-right">SYS TIME: --:--:--</span>
        </div>
        <pre id="terminal-content">LOADING SCOREBOARD...</pre>
    </div>

    <script>
        // Update clock and status layout
        function updateStatus() {
            const now = new Date();
            const timeStr = now.toTimeString().split(' ')[0];
            const isGamePage = window.location.pathname.startsWith('/game/');

            const leftEl = document.getElementById('status-left');
            const rightEl = document.getElementById('status-right');

            if (isGamePage) {
                if (!leftEl.querySelector('a')) {
                    leftEl.innerHTML = '<a href="/" class="term-link">&lt;&lt; BACK TO SCOREBOARD</a>';
                }
                rightEl.innerHTML = '<span class="dot-pulse"></span>LIVE FEED &bull; SYS TIME: ' + timeStr;
            } else {
                if (leftEl.querySelector('a') || leftEl.innerHTML === '') {
                    leftEl.innerHTML = '<span class="dot-pulse"></span>LIVE FEED';
                }
                rightEl.innerHTML = 'SYS TIME: ' + timeStr;
            }
        }
        updateStatus();
        setInterval(updateStatus, 1000);

        // Fetch live terminal data
        async function fetchTerminalData() {
            try {
                const url = new URL(window.location.href);
                url.searchParams.set('raw', '1');
                
                const res = await fetch(url.toString());
                if (!res.ok) throw new Error('HTTP status ' + res.status);
                const htmlText = await res.text();
                
                document.getElementById('terminal-content').innerHTML = htmlText;
            } catch (err) {
                console.error(err);
                document.getElementById('terminal-content').innerHTML = 
                    '<span class="term-red">ERROR CONNECTING TO STREAM: ' + err.message + '</span>\nRetrying in 10s...';
            }
        }

        // Initial fetch
        fetchTerminalData();
        // Poll every 10 seconds
        setInterval(fetchTerminalData, 10000);

        // Handle navigation inside terminal via AJAX
        document.addEventListener('click', async (e) => {
            if (e.target.classList.contains('term-link')) {
                e.preventDefault();
                const href = e.target.getAttribute('href');
                history.pushState(null, '', href);
                updateStatus(); // Immediate layout update
                document.getElementById('terminal-content').innerHTML = 'RETRIEVING FEED...';
                await fetchTerminalData();
            }
        });

        // Listen for back button
        window.addEventListener('popstate', () => {
            updateStatus();
            fetchTerminalData();
        });
    </script>
</body>
</html>
`
