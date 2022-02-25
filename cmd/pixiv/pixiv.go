package main

import (
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
	currentIndex    int
	currentFileName string
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

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
		pad + helpStyle(currentFileName)
}

func DownloadCmd() tea.Cmd {
	return (func() tea.Msg {
		currentIllustration := illustrations[currentIndex]
		currentFileName = "temp-" + currentIllustration.Title + ".jpg"
		pixivapi.DownloadIllustration(currentIllustration, currentFileName)
		currentIndex = currentIndex + 1
		return nil
	})
}
