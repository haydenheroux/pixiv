package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"haydenheroux.xyz/pixivapi"
)

const (
	padding  = 2
	maxWidth = 80
)

var (
	illustrations   []pixivapi.PixivIllustration
	downloadLog     map[string]bool
	currentIndex    int
	currentFileName string
)

var dimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
var okStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Render
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Render

func main() {
	m := model{
		progress: progress.New(progress.WithGradient("#8BE9FD", "#FF79C6")),
	}

	if err := tea.NewProgram(m).Start(); err != nil {
		os.Exit(1)
	}
}

type model struct {
	progress progress.Model
}

func (_ model) Init() tea.Cmd {
	illustrations, _ = pixivapi.GetTopIllustrations()
	downloadLog = make(map[string]bool, 0)
	return DownloadCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}

		if currentIndex >= len(illustrations) {
			return m, tea.Quit
		}

		cmd := m.progress.IncrPercent(1 / float64(len(illustrations)))
		return m, tea.Batch(DownloadCmd(), cmd)
	}
}

func (e model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + e.progress.View() + "\n\n" +
		pad + dimStyle(currentFileName) + "\n\n" +
		BuildDownloadLog()
}

func DownloadCmd() tea.Cmd {
	return (func() tea.Msg {
		currentIllustration := illustrations[currentIndex]
		currentFileName = "temp-" + strings.Join(currentIllustration.Tags, "+") + ".jpg"
		_, err := pixivapi.DownloadIllustration(currentIllustration, currentFileName)
		if err != nil {
			downloadLog[currentFileName] = true
		} else {
			downloadLog[currentFileName] = false
		}
		currentIndex = currentIndex + 1
		return nil
	})
}

func BuildDownloadLog() string {
	pad := strings.Repeat(" ", padding)
	var str string
	for fileName, status := range downloadLog {
		var indicator string
		if status == true {
			// There was an error
			indicator = errorStyle("✗")
		} else {
			// No error
			indicator = okStyle("✓")
		}
		str += fmt.Sprintf("%s[%s] %s\n", pad, indicator, dimStyle(fileName))
	}
	return str
}
