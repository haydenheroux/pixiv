package main

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
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

var dimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262"))
var okStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#50FA7B"))
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5555"))

type State int

const (
	stateInput = iota
	stateRetrieving
	stateDownloading
)

type model struct {
	spinner  spinner.Model
	progress progress.Model
	input    textinput.Model

	state         State
	illustrations []pixivapi.PixivIllustration

	errors          map[string]error
	queue           []string
	directory       string
	currentIndex    int
	currentFileName string
}

func createInput() textinput.Model {
	input := textinput.New()
	input.Placeholder = "Search..."
	input.PlaceholderStyle = dimStyle
	input.Focus()
	input.CharLimit = 80
	input.Width = maxWidth - 2*padding - 4
	return input
}

func createSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return s
}

func main() {
	m := model{
		spinner:  createSpinner(),
		progress: progress.New(progress.WithGradient("#8BE9FD", "#FF79C6")),
		input:    createInput(),
		state:    stateInput,
	}

	if err := tea.NewProgram(m).Start(); err != nil {
		os.Exit(1)
	}
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

		// NOTE(hayden): dateTimeFromString causes a crash if the UpdateDate is empty
		illustrations = slices.DeleteFunc(illustrations, func(i pixivapi.PixivIllustration) bool {
			return len(i.UpdateDate) == 0
		})

		return illustrationsMsg(illustrations)
	}
}

// TODO(hayden): Test this
func (m model) Resize(screenWidth int) model {
	width := screenWidth - padding*2
	width = min(width, maxWidth)
	m.progress.Width = width
	m.input.Width = width
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
		return m.Resize(msg.Width), nil

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.state = stateRetrieving
			input := m.input.Value()
			if len(input) == 0 {
				m.directory = time.Now().Format(time.DateTime)
			} else {
				m.directory = input
			}
			_ = os.Mkdir(m.directory, 0755)
			return m, tea.Batch(
				m.spinner.Tick,
				GetResults(input),
			)
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case illustrationsMsg:
		m.illustrations = msg
		if len(m.illustrations) == 0 {
			return m, tea.Quit
		}
		m.errors = make(map[string]error)
		m.state = stateDownloading
		return m, DownloadNext(m.directory, m.illustrations[0])

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
			DownloadNext(m.directory, m.illustrations[m.currentIndex]),
			incr,
		)
	}

	var inputCmd, spinCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	m.spinner, spinCmd = m.spinner.Update(msg)
	return m, tea.Batch(inputCmd, spinCmd)
}

func (m model) InputView() string {
	pad := strings.Repeat(" ", padding)
	return pad + m.input.View()
}

func (m model) SearchView() string {
	pad := strings.Repeat(" ", padding)
	return pad + m.spinner.View() + " " + dimStyle.Render(m.input.Value())
}

func (m model) DownloadingView() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.progress.View() + "\n" +
		pad + dimStyle.Render(m.currentFileName) + "\n\n" +
		buildDownloadLog(m.queue, m.errors) + "\n"
}

func (m model) View() string {
	switch m.state {
	case stateInput:
		return m.InputView()
	case stateRetrieving:
		return m.SearchView()
	case stateDownloading:
		return m.InputView() + "\n" + m.DownloadingView()
	default:
		return "\n"
	}
}

type downloadedMsg struct {
	filename string
	err      error
}

func DownloadNext(directory string, illustration pixivapi.PixivIllustration) tea.Cmd {
	return func() tea.Msg {
		filename := strings.Join(illustration.Tags, "+") + ".jpg"
		destination := path.Join(directory, filename)
		_, err := pixivapi.DownloadIllustration(illustration, destination)
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
			indicator = errorStyle.Render("✗")
		} else {
			indicator = okStyle.Render("✓")
		}
		str += fmt.Sprintf("%s[%s] %s\n", pad, indicator, dimStyle.Render(name))
	}
	return str
}
