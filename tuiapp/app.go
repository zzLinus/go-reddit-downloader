package tuiapp

// this file is basicly working on bubble team tui framework

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zzLinus/GoRedditDownloader/downloader"
	_ "github.com/zzLinus/GoRedditDownloader/extractor/reddit"
)

type errMsg error
type tickMsg time.Time
type respMsg int

const (
	padding  = 2
	maxWidth = 80
)

var (
	videoDownloader *downloader.Downloader
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy().PaddingTop(4).PaddingLeft(4)
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	focusedButton       = focusedStyle.Copy().Render("[ Submit ]")
	style               = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("63"))
	blurredButton       = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	spinnerStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).PaddingLeft(4).PaddingTop(4)
	anotherStyle        = lipgloss.NewStyle().PaddingLeft(4).BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("204"))
)

type model struct {
	spinner    spinner.Model
	loading    bool
	progress   progress.Model
	textInput  textinput.Model
	cursorMode textinput.CursorMode
	quitting   bool
	err        error
}

func New() *tea.Program {
	p := tea.NewProgram(initialModel())
	return p
}

// Add a purple, rectangular border

func initialModel() model {

	ti := textinput.New()
	ti.Placeholder = "pase an url that support by use"
	ti.Focus()
	ti.CharLimit = 80
	ti.Width = 80
	ti.CursorStyle = focusedStyle

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Italic(true)

	pro := progress.New(progress.WithDefaultGradient())

	return model{spinner: s, textInput: ti, loading: false, progress: pro}
}

func (m model) Init() tea.Cmd {
	videoDownloader = downloader.New()
	return tea.Batch(textinput.Blink)

}

// TODO:this section is just broken.... i started this project just to play around with bubbletea tui framework
// but end up with this reddit downloader stuff,so many code are useless the only reason thy are still here is
// because i'm lazy as fuck and after i play around the bubbletea stuff i forget to clean them...
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Batch(tickCmd())
		case "enter":
			m.loading = true
			return m, tea.Batch(m.spinner.Tick,
				func() tea.Msg {
					rowURL := m.textInput.Value()
					statusCode, err := videoDownloader.Download(rowURL)
					if err != nil {
						panic(err)
					}
					return respMsg(statusCode)
				})
		case "ctrl+r":
			m.cursorMode++
			if m.cursorMode > textinput.CursorHide {
				m.cursorMode = textinput.CursorBlink
			}
			cmd := m.textInput.SetCursorMode(m.cursorMode)
			return m, cmd

		default:
			var cmd tea.Cmd
			if !m.loading {
				m.textInput, cmd = m.textInput.Update(msg)
			}
			return m, tea.Batch(cmd, textinput.Blink)
		}

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case errMsg:
		m.err = msg
		return m, nil

	case respMsg:
		if msg == 200 {
			m.loading = false
		}
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}
		cmd := m.progress.IncrPercent(0.25)
		return m, tea.Batch(tickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		var cmd tea.Cmd
		if m.loading && !m.quitting {
			m.spinner, cmd = m.spinner.Update(msg)
		}
		return m, cmd
	}
}

func (m model) View() string {
	var str string = ""

	if m.err != nil {
		return m.err.Error()
	}

	if m.loading {
		str += spinnerStyle.Render(fmt.Sprintf("\n%s Downloading content\n", m.spinner.View()))
	}

	if !m.loading {
		str += cursorStyle.Render(fmt.Sprintf("%s", m.textInput.View()))
	}

	if m.quitting {
		pad := strings.Repeat(" ", padding)
		return "\n" +
			pad + m.progress.View() + "\n\n"
	}

	return str
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg { return tickMsg(t) })
}
