package sports

import (
	"fmt"
	"html"
	"net/http"
	"strings"
	"time"
	_ "time/tzdata"
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

// NewHandler constructs and returns the HTTP handler with all registered routes.
func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleSchedule)
	mux.HandleFunc("GET /game/{gamePk}", handleGame)
	mux.HandleFunc("GET /api/games", handleAPIGames)
	mux.HandleFunc("GET /api/game/{gamePk}", handleAPIGameDetail)

	// NBA Routes
	mux.HandleFunc("GET /nba", handleNBASchedule)
	mux.HandleFunc("GET /nba/game/{gamePk}", handleNBAGame)
	mux.HandleFunc("GET /api/nba/games", handleAPINBAGames)
	mux.HandleFunc("GET /api/nba/game/{gamePk}", handleAPINBAGameDetail)

	return mux
}

const htmlPage = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>sportstxt - MLB Scoreboard</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="color-scheme" content="light dark">
    <link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;700&display=swap" rel="stylesheet">
    <style>
        :root {
            /* Light theme variables by default */
            --term-bg: #f4f5f8;
            --term-container-bg: #ffffff;
            --term-border: #e1e3ec;
            --term-green: #0d6b38;
            --term-yellow: #a27b00;
            --term-red: #c62828;
            --term-cyan: #007791;
            --term-blue: #1565c0;
            --term-magenta: #8e24aa;
            --term-gray: #6b7280;
            --term-link-hover: #000000;
            --color-primary: var(--term-green);
        }

        @media (prefers-color-scheme: dark) {
            :root {
                /* Dark theme overrides */
                --term-bg: #07080c;
                --term-container-bg: #07080c;
                --term-border: #1f222e;
                --term-green: #39ff14;
                --term-yellow: #ffeb3b;
                --term-red: #ff3b30;
                --term-cyan: #00f0ff;
                --term-blue: #3b82f6;
                --term-magenta: #d946ef;
                --term-gray: #555866;
                --term-link-hover: #ffffff;
            }
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
            transition: background-color 0.3s ease, color 0.3s ease;
        }

        .term-container {
            width: 100%;
            max-width: 900px;
            background: var(--term-container-bg);
            border: 1px solid var(--term-border);
            border-radius: 8px;
            padding: 30px;
            box-sizing: border-box;
            position: relative;
            transition: background-color 0.3s ease, border-color 0.3s ease;
        }

        pre {
            margin: 0;
            font-family: inherit;
            font-size: 15px;
            line-height: 1.6;
            white-space: pre;
            position: relative;
            overflow-x: auto;
            --scrollbar-thumb: var(--term-border);
            --scrollbar-track: transparent;
            scrollbar-color: var(--scrollbar-thumb) var(--scrollbar-track);
            scrollbar-width: thin;
        }

        /* Legacy fallback for WebKit/Blink browsers */
        @supports not (scrollbar-color: auto) {
            pre::-webkit-scrollbar {
                height: 6px;
            }
            pre::-webkit-scrollbar-thumb {
                background: var(--scrollbar-thumb);
                border-radius: 3px;
            }
            pre::-webkit-scrollbar-track {
                background: var(--scrollbar-track);
            }
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
            transition: color 0.2s ease;
        }
        .term-link:hover {
            color: var(--term-link-hover) !important;
        }

        /* Status bar */
        .status-bar {
            display: flex;
            justify-content: space-between;
            margin-bottom: 15px;
            border-bottom: 1px solid var(--term-border);
            padding-bottom: 10px;
            font-size: 12px;
            color: var(--term-gray);
            transition: border-color 0.3s ease, color 0.3s ease;
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
            transition: background-color 0.3s ease;
        }

        @keyframes pulse {
            0% { transform: scale(0.95); opacity: 0.5; }
            70% { transform: scale(1); opacity: 1; }
            100% { transform: scale(0.95); opacity: 0.5; }
        }

        /* Responsive design queries for mobile / smaller viewports */
        @media (max-width: 768px) {
            body {
                padding: 12px;
                justify-content: flex-start;
            }
            .term-container {
                padding: 16px;
                border-radius: 6px;
            }
            pre {
                font-size: 13px;
            }
            .status-bar {
                font-size: 10px;
                margin-bottom: 12px;
                padding-bottom: 8px;
            }
        }

        @media (max-width: 600px) {
            body {
                padding: 8px;
            }
            .term-container {
                padding: 12px;
                border-radius: 4px;
            }
            pre {
                font-size: 11px;
            }
        }

        @media (max-width: 480px) {
            body {
                padding: 6px;
            }
            .term-container {
                padding: 10px;
            }
            pre {
                font-size: 9.5px;
            }
            .status-bar {
                font-size: 9px;
            }
        }

        @media (max-width: 380px) {
            pre {
                font-size: 8.5px;
            }
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
            const isNBAGamePage = window.location.pathname.startsWith('/nba/game/');

            const leftEl = document.getElementById('status-left');
            const rightEl = document.getElementById('status-right');

            if (isGamePage || isNBAGamePage) {
                if (!leftEl.querySelector('a')) {
                    const backUrl = isNBAGamePage ? '/nba' : '/';
                    leftEl.innerHTML = '<a href="' + backUrl + '" class="term-link">&lt;&lt; BACK TO SCOREBOARD</a>';
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
                if (!url.searchParams.has('tz')) {
                    url.searchParams.set('tz', Intl.DateTimeFormat().resolvedOptions().timeZone);
                }
                
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
