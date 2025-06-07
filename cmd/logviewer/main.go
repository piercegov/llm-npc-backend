package main

import (
	"bufio"
	"fmt"
	"image/color"
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

	clearBtn := widget.NewButton("Clear Logs", lv.clearLogs)

	toolbar := container.NewBorder(
		nil, nil, nil, clearBtn,
		lv.searchBar,
	)

	logScroll := container.NewScroll(lv.logContainer)
	logScroll.SetMinSize(fyne.NewSize(900, 600))

	return container.NewBorder(
		toolbar,
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
