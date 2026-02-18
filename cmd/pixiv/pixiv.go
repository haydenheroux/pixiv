package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"haydenheroux.xyz/pixivapi"
)

const (
	padding        = 2
	maxWidth       = 80
	maxQueueLength = 5
)

var ()

var dimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
var okStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B")).Render
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555")).Render

func createInput() textinput.Model {
	input := textinput.New()
	input.Placeholder = "Search..."
	input.Focus()
	input.CharLimit = 80
	input.Width = maxWidth - 2*padding - 4
	return input
}

func main() {
	m := model{
		progress: progress.New(progress.WithGradient("#8BE9FD", "#FF79C6")),
		input:    createInput(),
	}

	if err := tea.NewProgram(m).Start(); err != nil {
		os.Exit(1)
	}
}

type model struct {
	progress progress.Model
	input    textinput.Model

	downloading   bool
	illustrations []pixivapi.PixivIllustration

	errors          map[string]error
	queue           []string
	currentIndex    int
	currentFileName string
}

func (_ model) Init() tea.Cmd {
	return textinput.Blink
}

type illustrationsMsg []pixivapi.PixivIllustration

func GetResults(search string) tea.Cmd {
	return func() tea.Msg {
		illustrations := make([]pixivapi.PixivIllustration, 0)
		if len(search) == 0 {
			illustrations, _ = pixivapi.GetTopIllustrations()
		} else {
			illustrations, _ = pixivapi.GetSearchIllustrations(search)
		}
		return illustrationsMsg(illustrations)
	}
}

func Resize(m model, width int) model {
	m.progress.Width = width - padding*2 - 4
	if m.progress.Width > maxWidth {
		m.progress.Width = maxWidth
	}
	return m
}

func PushAndCycleQueue(m model, filename string) model {
	m.queue = append(m.queue, filename)
	for len(m.queue) > maxQueueLength {
		removed := m.queue[0]
		delete(m.errors, removed)
		m.queue = m.queue[1:]
	}
	m.currentIndex = m.currentIndex + 1
	return m
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return Resize(m, msg.Width), nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.downloading = true
			return m, GetResults(m.input.Value())
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case illustrationsMsg:
		m.illustrations = msg
		m.errors = make(map[string]error)
		return m, DownloadNext(m.illustrations[0])

	case downloadedMsg:
		m = PushAndCycleQueue(m, msg.filename)
		m.errors[msg.filename] = msg.err

		progressDone := m.progress.Percent() >= 1.0
		allDownloaded := m.currentIndex >= len(m.illustrations)

		if progressDone || allDownloaded {
			return m, tea.Quit
		}

		delta := 1 / float64(len(m.illustrations))
		incr := m.progress.IncrPercent(delta)

		return m, tea.Batch(
			DownloadNext(m.illustrations[m.currentIndex]),
			incr,
		)
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.downloading {
		pad := strings.Repeat(" ", padding)
		return "\n" +
			pad + m.progress.View() + "\n\n" +
			pad + dimStyle(m.currentFileName) + "\n\n" +
			buildDownloadLog(m.queue, m.errors) + "\n"
	} else {
		return m.input.View()
	}
}

type downloadedMsg struct {
	filename string
	err      error
}

func DownloadNext(illustration pixivapi.PixivIllustration) tea.Cmd {
	return func() tea.Msg {
		filename := "temp-" + strings.Join(illustration.Tags, "+") + ".jpg"
		_, err := pixivapi.DownloadIllustration(illustration, filename)
		return downloadedMsg{filename, err}
	}
}

func buildDownloadLog(filenames []string, statuses map[string]error) string {
	pad := strings.Repeat(" ", padding)
	var str string
	for _, name := range filenames {
		var indicator string
		err := statuses[name]
		if err != nil {
			indicator = errorStyle("✗")
		} else {
			indicator = okStyle("✓")
		}
		str += fmt.Sprintf("%s[%s] %s\n", pad, indicator, dimStyle(name))
	}
	return str
}
