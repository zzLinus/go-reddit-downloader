package tuiapp

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zzLinus/GoTUITODOList/downloader"
	"github.com/zzLinus/GoTUITODOList/extractor"
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
	rowURLExtractor *extractor.Extractor
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	focusedButton       = focusedStyle.Copy().Render("[ Submit ]")
	style               = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("63"))
	blurredButton       = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
	spinnerStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).PaddingLeft(4)
	anotherStyle        = lipgloss.NewStyle().PaddingLeft(4).BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("204"))
)

type model struct {
	spinner    spinner.Model
	border     lipgloss.Border
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
	ti.Width = 20
	ti.CursorStyle = focusedStyle

	cm := textinput.CursorBlink

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Italic(true)

	pro := progress.New(progress.WithDefaultGradient())

	return model{spinner: s, textInput: ti, loading: false, progress: pro, cursorMode: cm}
}

func (m model) Init() tea.Cmd {
	videoDownloader = downloader.New()
	rowURLExtractor = extractor.New()
	return tea.Batch(textinput.Blink, spinner.Tick)

}

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
					// "https://v.redd.it/8akffrc6fqx81/DASH_720.mp4"
					rowURL := m.textInput.Value()
					statusCode, err := videoDownloader.Download(rowURL)
					if err != nil {
						panic(err)
					}
					return respMsg(statusCode)
				})
		default:
			var cmd tea.Cmd
			if !m.loading {
				m.textInput, cmd = m.textInput.Update(msg)
			}
			return m, cmd
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
		str += fmt.Sprintf("input or paset url here\n%s", m.textInput.View())
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
