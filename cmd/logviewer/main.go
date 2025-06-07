package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image/color"
	"net"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type LogViewer struct {
	logContainer *fyne.Container
	searchBar    *widget.Entry
	consoleBar   *widget.Entry
	logBinding   binding.String
	allLogs      []string
	cmd          *exec.Cmd
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("LLM NPC Backend - Log Viewer")
	myWindow.Resize(fyne.NewSize(1000, 700))

	logViewer := NewLogViewer()
	content := logViewer.CreateUI()

	myWindow.SetContent(content)

	// Start the backend process and log streaming
	go logViewer.StartBackendAndStream()

	myWindow.ShowAndRun()

	if logViewer.cmd != nil && logViewer.cmd.Process != nil {
		logViewer.cmd.Process.Kill()
	}
}

func NewLogViewer() *LogViewer {
	logBinding := binding.NewString()
	logContainer := container.NewVBox()

	return &LogViewer{
		logContainer: logContainer,
		logBinding:   logBinding,
		allLogs:      make([]string, 0),
	}
}

func (lv *LogViewer) CreateUI() *fyne.Container {
	lv.searchBar = widget.NewEntry()
	lv.searchBar.SetPlaceHolder("Search logs (Cmd+F)...")
	lv.searchBar.OnChanged = lv.searchLogs
	lv.searchBar.Resize(fyne.NewSize(600, 40))

	lv.consoleBar = widget.NewEntry()
	lv.consoleBar.SetPlaceHolder("Console command (e.g., 'read_scratchpads')...")
	lv.consoleBar.OnSubmitted = lv.executeConsoleCommand
	lv.consoleBar.Resize(fyne.NewSize(600, 40))

	clearBtn := widget.NewButton("Clear Logs", lv.clearLogs)

	toolbar := container.NewBorder(
		nil, nil, nil, clearBtn,
		lv.searchBar,
	)

	consolebar := container.NewBorder(
		nil, nil, widget.NewLabel("Console:"), nil,
		lv.consoleBar,
	)

	logScroll := container.NewScroll(lv.logContainer)
	logScroll.SetMinSize(fyne.NewSize(900, 550))

	return container.NewBorder(
		container.NewVBox(toolbar, consolebar),
		nil,
		nil,
		nil,
		logScroll,
	)
}

func (lv *LogViewer) StartBackendAndStream() {
	// Change to the project directory
	projectDir := "/Users/piercegovernale/Documents/llm-npc-backend"

	// Build the backend first
	buildCmd := exec.Command("go", "build", "./cmd/backend/...")
	buildCmd.Dir = projectDir
	if err := buildCmd.Run(); err != nil {
		lv.appendLog(fmt.Sprintf("Failed to build backend: %v\n", err))
		return
	}

	// Start the backend
	lv.cmd = exec.Command("./backend")
	lv.cmd.Dir = projectDir

	stdout, err := lv.cmd.StdoutPipe()
	if err != nil {
		lv.appendLog(fmt.Sprintf("Failed to create stdout pipe: %v\n", err))
		return
	}

	stderr, err := lv.cmd.StderrPipe()
	if err != nil {
		lv.appendLog(fmt.Sprintf("Failed to create stderr pipe: %v\n", err))
		return
	}

	if err := lv.cmd.Start(); err != nil {
		lv.appendLog(fmt.Sprintf("Failed to start backend: %v\n", err))
		return
	}

	lv.appendLog("Backend started successfully!\n")

	// Read stdout
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			lv.appendLog(line + "\n")
		}
	}()

	// Read stderr
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			lv.appendLog("[ERROR] " + line + "\n")
		}
	}()

	// Wait for process to finish
	go func() {
		if err := lv.cmd.Wait(); err != nil {
			lv.appendLog(fmt.Sprintf("Backend process ended with error: %v\n", err))
		} else {
			lv.appendLog("Backend process ended normally\n")
		}
	}()
}

func (lv *LogViewer) appendLog(text string) {
	timestamp := time.Now().Format("15:04:05")
	logLine := fmt.Sprintf("[%s] %s", timestamp, strings.TrimSpace(text))

	lv.allLogs = append(lv.allLogs, logLine)

	fyne.Do(func() {
		lv.updateDisplayInternal()
	})
}

func (lv *LogViewer) createLogEntry(text string, isAlternate bool) fyne.CanvasObject {
	entry := widget.NewEntry()
	entry.MultiLine = true
	entry.Wrapping = fyne.TextWrapWord
	entry.SetText(text)
	entry.OnChanged = func(string) {}
	entry.Disable()
	entry.Enable()

	var bgColor color.Color
	if isAlternate {
		bgColor = color.RGBA{40, 40, 40, 255}
	} else {
		bgColor = color.RGBA{30, 30, 30, 255}
	}

	bg := canvas.NewRectangle(bgColor)
	return container.NewStack(bg, container.NewPadded(entry))
}

func (lv *LogViewer) updateDisplay() {
	lv.updateDisplayInternal()
}

func (lv *LogViewer) updateDisplayInternal() {
	lv.logContainer.Objects = nil

	var logsToShow []string
	searchTerm := lv.searchBar.Text
	if searchTerm != "" {
		logsToShow = lv.filterLogs(searchTerm)
	} else {
		logsToShow = lv.allLogs
	}

	for i, logLine := range logsToShow {
		entry := lv.createLogEntry(logLine, i%2 == 1)
		lv.logContainer.Add(entry)
	}

	lv.logContainer.Refresh()
}

func (lv *LogViewer) searchLogs(searchTerm string) {
	lv.updateDisplay()
}

func (lv *LogViewer) filterLogs(searchTerm string) []string {
	if searchTerm == "" {
		return lv.allLogs
	}

	var filteredLines []string

	for _, line := range lv.allLogs {
		if strings.Contains(strings.ToLower(line), strings.ToLower(searchTerm)) {
			filteredLines = append(filteredLines, line)
		}
	}

	return filteredLines
}

func (lv *LogViewer) clearLogs() {
	lv.allLogs = make([]string, 0)
	lv.updateDisplay()
}

func (lv *LogViewer) executeConsoleCommand(command string) {
	lv.appendLog(fmt.Sprintf("> %s", command))

	switch command {
	case "read_scratchpads":
		lv.readScratchpads()
	default:
		lv.appendLog(fmt.Sprintf("Unknown command: %s", command))
	}

	lv.consoleBar.SetText("")
}

func (lv *LogViewer) readScratchpads() {
	socketPath := "/tmp/llm-npc-backend.sock"

	// Create HTTP client that uses Unix domain socket
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}

	// Make request to console endpoint
	resp, err := client.Get("http://unix/console/read_scratchpads")
	if err != nil {
		lv.appendLog(fmt.Sprintf("Failed to read scratchpads: %v", err))
		return
	}
	defer resp.Body.Close()

	// Parse response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		lv.appendLog(fmt.Sprintf("Failed to parse response: %v", err))
		return
	}

	// Display results
	if success, ok := result["success"].(bool); ok && success {
		if data, ok := result["data"].(map[string]interface{}); ok {
			if len(data) == 0 {
				lv.appendLog("No scratchpads found")
			} else {
				lv.appendLog(fmt.Sprintf("Found %d NPCs with scratchpads:", len(data)))
				for npcID, npcData := range data {
					if npcInfo, ok := npcData.(map[string]interface{}); ok {
						count := npcInfo["count"].(float64)
						lv.appendLog(fmt.Sprintf("  %s: %d entries", npcID, int(count)))

						if entries, ok := npcInfo["entries"].([]interface{}); ok {
							for _, entry := range entries {
								if entryMap, ok := entry.(map[string]interface{}); ok {
									key := entryMap["key"].(string)
									value := entryMap["value"].(string)
									timestamp := entryMap["timestamp"].(string)
									lv.appendLog(fmt.Sprintf("    %s: %s (at %s)", key, value, timestamp))
								}
							}
						}
					}
				}
			}
		}
	} else {
		lv.appendLog("Failed to read scratchpads")
	}
}
