package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func init() {
	_ = os.Mkdir("./tmp", os.ModePerm)
	emptyTmp()
}

// Page data structure
type PageData struct {
	Algorithm     string
	ImageData     string
	AnimationData string
	IsGenerated   bool
	SolutionSteps int
	NodesExplored int
	TimeTaken     string
	Width         int
	Height        int
	HasAnimation  bool
	MazeType      string
}

// HTML template with animation support
const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>🌐 Maze Visualizer</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Courier New', monospace;
            background: linear-gradient(135deg, #0a0a19 0%, #1a0a2e 100%);
            color: #00ffff;
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
        }

        header {
            text-align: center;
            margin-bottom: 40px;
            padding: 30px;
            background: rgba(25, 15, 45, 0.6);
            border: 2px solid #00c8ff;
            border-radius: 15px;
            box-shadow: 0 0 30px rgba(0, 200, 255, 0.3);
        }

        h1 {
            font-size: 3em;
            color: #39ff14;
            text-shadow: 0 0 20px #39ff14;
            margin-bottom: 10px;
        }

        .subtitle {
            color: #ff0080;
            font-size: 1.2em;
            text-shadow: 0 0 10px #ff0080;
        }

        .controls {
            background: rgba(20, 25, 45, 0.8);
            padding: 30px;
            border-radius: 15px;
            border: 2px solid #bf40bf;
            box-shadow: 0 0 20px rgba(191, 64, 191, 0.3);
            margin-bottom: 30px;
        }

        .form-group {
            margin-bottom: 20px;
        }

        label {
            display: block;
            margin-bottom: 10px;
            color: #00ffff;
            font-size: 1.1em;
            text-shadow: 0 0 5px #00ffff;
        }

        select, input[type="number"], input[type="text"] {
            width: 100%;
            padding: 12px;
            background: #0a0a19;
            border: 2px solid #00c8ff;
            border-radius: 8px;
            color: #00ffff;
            font-size: 1em;
            font-family: 'Courier New', monospace;
            transition: all 0.3s;
        }

        select:hover, input:hover {
            border-color: #39ff14;
            box-shadow: 0 0 15px rgba(57, 255, 20, 0.3);
        }

        select:focus, input:focus {
            outline: none;
            border-color: #ff0080;
            box-shadow: 0 0 20px rgba(255, 0, 128, 0.5);
        }

        button {
            width: 100%;
            padding: 15px 30px;
            background: linear-gradient(135deg, #bf40bf 0%, #ff0080 100%);
            border: none;
            border-radius: 8px;
            color: white;
            font-size: 1.1em;
            font-weight: bold;
            cursor: pointer;
            transition: all 0.3s;
            text-shadow: 0 0 10px rgba(0, 0, 0, 0.5);
            box-shadow: 0 0 20px rgba(191, 64, 191, 0.5);
            margin-top: 20px;
        }

        button:hover {
            transform: translateY(-2px);
            box-shadow: 0 0 30px rgba(255, 0, 128, 0.8);
        }

        .result-container {
            background: rgba(20, 25, 45, 0.8);
            padding: 30px;
            border-radius: 15px;
            border: 2px solid #00ffff;
            box-shadow: 0 0 20px rgba(0, 255, 255, 0.3);
            text-align: center;
        }

        .result-container h2 {
            color: #39ff14;
            margin-bottom: 20px;
            text-shadow: 0 0 15px #39ff14;
        }

        .view-toggle {
            display: flex;
            gap: 10px;
            justify-content: center;
            margin-bottom: 20px;
        }

        .toggle-btn {
            padding: 10px 20px;
            background: rgba(75, 50, 150, 0.5);
            border: 2px solid #4b3296;
            border-radius: 8px;
            color: #00ffff;
            cursor: pointer;
            transition: all 0.3s;
            font-family: 'Courier New', monospace;
            font-size: 1em;
        }

        .toggle-btn.active {
            background: linear-gradient(135deg, #00c8ff 0%, #39ff14 100%);
            color: #0a0a19;
            border-color: #39ff14;
            font-weight: bold;
        }

        .toggle-btn:hover {
            border-color: #00ffff;
            box-shadow: 0 0 15px rgba(0, 255, 255, 0.3);
        }

        .maze-display {
            position: relative;
        }

        .maze-image {
            max-width: 100%;
            border: 3px solid #00c8ff;
            border-radius: 10px;
            box-shadow: 0 0 40px rgba(0, 200, 255, 0.5);
            display: none;
        }

        .maze-image.active {
            display: block;
        }

        .loading {
            display: none;
            color: #ff0080;
            font-size: 1.5em;
            margin: 20px 0;
            animation: pulse 1.5s ease-in-out infinite;
        }

        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.3; }
        }

        .stats-box {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin: 20px 0;
        }

        .stat-item {
            background: rgba(75, 50, 150, 0.3);
            padding: 15px;
            border-radius: 10px;
            border: 1px solid #4b3296;
        }

        .stat-label {
            color: #00c8ff;
            font-size: 0.9em;
            margin-bottom: 5px;
        }

        .stat-value {
            color: #39ff14;
            font-size: 1.5em;
            font-weight: bold;
            text-shadow: 0 0 10px #39ff14;
        }

        .info-box {
            background: rgba(75, 50, 150, 0.3);
            padding: 20px;
            border-radius: 10px;
            border: 1px solid #4b3296;
            margin-top: 20px;
            text-align: left;
        }

        .info-box h3 {
            color: #00ffff;
            margin-bottom: 10px;
        }

        .info-box ul {
            list-style: none;
            padding-left: 0;
        }

        .info-box li {
            padding: 5px 0;
            color: #00c8ff;
        }

        .info-box li:before {
            content: "▸ ";
            color: #39ff14;
        }
    </style>
</head>
<body>
    <div class="container">
        <header>
            <h1>🌐 MAZE VISUALIZER</h1>
            <p class="subtitle">Pathfinding Algorithms</p>
        </header>

        <form method="POST" action="/" id="mazeForm">
            <div class="controls">
                <div class="form-group">
                    <label for="algorithm">🔮 Select Algorithm:</label>
                    <select name="algorithm" id="algorithm" required>
                        <option value="">-- Choose Algorithm --</option>
                        <option value="DFS" {{if eq .Algorithm "DFS"}}selected{{end}}>DFS (Depth-First Search)</option>
                        <option value="BFS" {{if eq .Algorithm "BFS"}}selected{{end}}>BFS (Breadth-First Search)</option>
                        <option value="Dijkstra" {{if eq .Algorithm "Dijkstra"}}selected{{end}}>Dijkstra's Algorithm</option>
                        <option value="AStar" {{if eq .Algorithm "AStar"}}selected{{end}}>A* Algorithm</option>
                    </select>
                </div>

                <div class="form-group">
                    <label for="maze">📁 Maze File:</label>
                    <select name="maze" id="maze" required>
                        <option value="">-- Choose MazeFile --</option>
                        <option value="maze.txt" {{if eq .MazeType "maze.txt"}}selected{{end}}>maze1.txt</option>
                        <option value="maze2.txt" {{if eq .MazeType "maze2.txt"}}selected{{end}}>maze2.txt</option>
                    </select>
                </div>

                <button type="submit">⚡ Generate Maze Animation</button>
            </div>
        </form>

        <div class="loading" id="loading">
            🎬 Generating maze animation... Please wait...
        </div>

        {{if .IsGenerated}}
        <div class="result-container">
            <h2>✨ Generated Maze - {{.Algorithm}} Algorithm</h2>
            
            <div class="stats-box">
                <div class="stat-item">
                    <div class="stat-label">📊 Solution Steps</div>
                    <div class="stat-value">{{.SolutionSteps}}</div>
                </div>
                <div class="stat-item">
                    <div class="stat-label">🔍 Nodes Explored</div>
                    <div class="stat-value">{{.NodesExplored}}</div>
                </div>
                <div class="stat-item">
                    <div class="stat-label">⏱️ Time Taken</div>
                    <div class="stat-value">{{.TimeTaken}}</div>
                </div>
                <div class="stat-item">
                    <div class="stat-label">📐 Maze Size</div>
                    <div class="stat-value">{{.Width}}×{{.Height}}</div>
                </div>
            </div>

            {{if .HasAnimation}}
            <div class="view-toggle">
                <button type="button" class="toggle-btn active" onclick="showAnimation()">🎬 Animation</button>
                <button type="button" class="toggle-btn" onclick="showStatic()">🖼️ Final Result</button>
            </div>
            {{end}}

            <div class="maze-display">
                {{if .HasAnimation}}
                <img id="animationImg" src="data:image/png;base64,{{.AnimationData}}" 
                     alt="Maze Animation" class="maze-image active">
                {{end}}
                <img id="staticImg" src="data:image/png;base64,{{.ImageData}}" 
                     alt="Final Maze" class="maze-image {{if not .HasAnimation}}active{{end}}">
            </div>
            
            <div class="info-box">
                <h3>Algorithm Info:</h3>
                <ul>
                    {{if eq .Algorithm "DFS"}}
                    <li>Type: Uninformed Search</li>
                    <li>Strategy: Explores deep into paths before backtracking</li>
                    <li>Complete: Yes (for finite graphs)</li>
                    <li>Optimal: No</li>
                    {{else if eq .Algorithm "BFS"}}
                    <li>Type: Uninformed Search</li>
                    <li>Strategy: Explores level by level</li>
                    <li>Complete: Yes</li>
                    <li>Optimal: Yes (for unweighted graphs)</li>
                    {{else if eq .Algorithm "Dijkstra"}}
                    <li>Type: Uninformed Search</li>
                    <li>Strategy: Explores by shortest distance from start</li>
                    <li>Complete: Yes</li>
                    <li>Optimal: Yes</li>
                    {{else if eq .Algorithm "AStar"}}
                    <li>Type: Informed Search (Heuristic)</li>
                    <li>Strategy: Uses heuristic to guide search toward goal</li>
                    <li>Complete: Yes</li>
                    <li>Optimal: Yes (with admissible heuristic)</li>
                    {{end}}
                </ul>
            </div>
        </div>
        {{end}}
    </div>

    <script>
        document.getElementById('mazeForm').addEventListener('submit', function() {
            document.getElementById('loading').style.display = 'block';
        });

        function showAnimation() {
            document.getElementById('animationImg').classList.add('active');
            document.getElementById('staticImg').classList.remove('active');
            document.querySelectorAll('.toggle-btn')[0].classList.add('active');
            document.querySelectorAll('.toggle-btn')[1].classList.remove('active');
            
            // Force reload animation to play again
            const animImg = document.getElementById('animationImg');
            const src = animImg.src;
            animImg.src = '';
            animImg.src = src;
        }

        function showStatic() {
            document.getElementById('animationImg').classList.remove('active');
            document.getElementById('staticImg').classList.add('active');
            document.querySelectorAll('.toggle-btn')[0].classList.remove('active');
            document.querySelectorAll('.toggle-btn')[1].classList.add('active');
        }
    </script>
</body>
</html>
`

func main() {
	tmpl := template.Must(template.New("index").Parse(htmlTemplate))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			data := PageData{
				IsGenerated: false,
			}
			tmpl.Execute(w, data)
			return
		}

		if r.Method == "POST" {
			algorithm := r.FormValue("algorithm")
			mazeFile := r.FormValue("maze")

			fmt.Println("mazefile : ", mazeFile)
			if mazeFile == "" {
				mazeFile = "maze.txt"
			}

			fmt.Printf("🎯 Generating maze with %s algorithm from %s...\n", algorithm, mazeFile)

			var m Maze
			m.Animate = true
			err := m.loadMaze(mazeFile)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error loading maze: %v", err), http.StatusInternalServerError)
				return
			}

			fmt.Printf("📐 Maze dimensions: %d×%d\n", m.Height, m.Width-1)

			startTime := time.Now()

			// Solve based on algorithm
			switch algorithm {
			case "DFS":
				m.SearchType = DFS
				solveDFS(&m)
			case "BFS":
				m.SearchType = BFS
				solveBFS(&m)
			case "Dijkstra":
				m.SearchType = DIJKSTRA
				solveDijkstra(&m)
			case "AStar":
				m.SearchType = ASTAR
				solveAStar(&m)
			default:
				http.Error(w, "Invalid search type", http.StatusBadRequest)
				return
			}

			timeTaken := time.Since(startTime)

			if len(m.Solution.Actions) == 0 {
				http.Error(w, "No solution found for this maze", http.StatusInternalServerError)
				return
			}

			fmt.Println("✅ Solution found!")
			fmt.Printf("📊 Solution steps: %d\n", len(m.Solution.Cells))
			fmt.Printf("🔍 Nodes explored: %d\n", len(m.Explored))
			fmt.Printf("⏱️  Time taken: %v\n", timeTaken)

			// Generate static image
			m.OutputImage("image.png")

			// Generate animation
			m.OutputAnimatedImage()

			// Read both static and animation images
			staticData := ""
			animationData := ""
			hasAnimation := false

			// Read static image (always available)
			if imgBytes, err := ioutil.ReadFile("./image.png"); err == nil {
				staticData = base64.StdEncoding.EncodeToString(imgBytes)
			} else {
				http.Error(w, "Error reading static image", http.StatusInternalServerError)
				return
			}

			// Read animation if it exists
			if _, err := os.Stat("./animation.png"); err == nil {
				if imgBytes, err := ioutil.ReadFile("./animation.png"); err == nil {
					animationData = base64.StdEncoding.EncodeToString(imgBytes)
					hasAnimation = true
					fmt.Println("✅ Animation loaded successfully")
				}
			} else {
				fmt.Println("⚠️  Animation file not found, showing static image only")
			}

			// Render result
			data := PageData{
				Algorithm:     algorithm,
				ImageData:     staticData,
				AnimationData: animationData,
				IsGenerated:   true,
				SolutionSteps: len(m.Solution.Cells),
				NodesExplored: len(m.Explored),
				TimeTaken:     timeTaken.String(),
				Width:         m.Width,
				Height:        m.Height,
				HasAnimation:  hasAnimation,
				MazeType:      mazeFile,
			}
			tmpl.Execute(w, data)
			return
		}
	})

	port := ":8080"
	fmt.Printf("🚀 Maze Visualizer starting on http://localhost%s\n", port)
	fmt.Println("✨ theme activated!")
	fmt.Println("📁 Make sure your maze.txt file is in the same directory")
	log.Fatal(http.ListenAndServe(port, nil))
}

func solveDFS(m *Maze) {
	var s DepthFirstSearch
	s.Game = m
	fmt.Println("🎯 Goal is:", s.Game.Goal)
	s.Solve()
}

func solveBFS(m *Maze) {
	var s BreadthFirstSearch
	s.Game = m
	s.Solve()
	fmt.Println("🔄 BFS solver called")
}

func solveDijkstra(m *Maze) {
	var s DijkstraSearch
	s.Game = m
	s.Solve()
	fmt.Println("📊 Dijkstra solver called")
}

func solveAStar(m *Maze) {
	var s AstrSearch
	s.Game = m
	s.Solve()
	fmt.Println("⭐ A* solver called")
}

func atoi(s string, defaultValue int) int {
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	return defaultValue
}
