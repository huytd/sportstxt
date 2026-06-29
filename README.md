# sportstxt

Get sports updates in plain text.

"Screenshot"

```
========================================================================
 [FINAL]  ATH 1  @  SF 2 - Bottom 9th
 Athletics @ San Francisco Giants
========================================================================
        2nd                [COUNT & OUTS]
        [◇]                Balls:   ○ ○ ○
   3rd       1st           Strikes: ● ○
   [◇]       [◇]           Outs:    ● ●
                           Batter:  Victor Bericoto
                                    Today:  2-4 | HR, K, RBI
                                    Season: .227 AVG, 2 HR, 2 RBI
                           Pitcher: Erik Miller
                                    Today:  1.0 IP, 0 ER, K, BB
                                    Season: 4.29 ERA, 1-0, 1.67 WHIP
 Pitches:
  P1: 0-0, 100mph four-seam fastball, foul
  P2: 0-1, 87mph slider, in play, run(s)

------------------------------------------------------------------------
 INNINGS     1  2  3  4  5  6  7  8  9 |  R  H  E
------------------------------------------------------------------------
 ATH         0  0  0  0  0  0  0  1  0 |  1  3  0
 SF          0  0  0  0  0  0  0  0  2 |  2  6  1
------------------------------------------------------------------------
 RECENT PLAYS:

 --- Bottom 9 ---
 [B9] Victor Bericoto homers (2) on a fly ball to center field. (FE) (Current Play)
 Score: 1-2
 [B9] Jung Hoo Lee flies out to right fielder Lawrence Butler. (CFFBFX)
 Score: 1-1
 [B9] Willy Adames flies out to center fielder Henry Bolte. (X)
 Score: 1-1
 [B9] Rafael Devers homers (12) on a fly ball to center field. (BSE)
 Score: 1-1

 --- Top 9 ---
 [T9] Jonah Heim flies out to right fielder Jung Hoo Lee. (X)
 Score: 1-0

 ...

```

## Betting Odds Feature

The game page can display betting odds from major sportsbooks including moneyline, spread, and over/under totals.

### Setup

1. Get a free API key from [The Odds API](https://the-odds-api.com/) (500 calls/day free tier)
2. Set the environment variable:

```bash
export THEODDS_API_KEY="your-api-key-here"
./sportstxt
```

Or with docker:

```bash
docker run -e THEODDS_API_KEY="your-api-key-here" -p 9090:9090 your-image
```

### What You'll See

When odds are available, the game page displays:

- **Moneyline**: Win probability odds for each team (e.g., LAD -150 | ATL +130)
- **Spread**: Run line with associated odds
- **Total**: Over/under run totals

Odds are aggregated from multiple bookmakers (DraftKings, FanDuel, BetMGM, Caesars, etc.) and updated every 15 minutes.

### Example Output

```
================================================================================
                        BETTING ODDS
================================================================================

 MONEYLINE:
  DraftKings: CWS +145 | BAL -165
  FanDuel: CWS +150 | BAL -170
  BetMGM: CWS +148 | BAL -168

 SPREAD:
  DraftKings: Chicago White Sox +1.5 (-135) | Baltimore Orioles -1.5 (+115)
  FanDuel: Chicago White Sox +1.5 (-140) | Baltimore Orioles -1.5 (+120)

 TOTAL (OVER/UNDER):
  DraftKings: Over 8.5 (-110) | Under 8.5 (-110)
  FanDuel: Over 8.5 (-108) | Under 8.5 (-112)

  Last updated: Jun 29, 2026 1:00 PM PST
================================================================================
```

### Notes

- If no API key is set, the odds section simply won't appear
- Odds are only shown for games that have betting lines available
- The feature works in both ANSI (terminal) and HTML modes
