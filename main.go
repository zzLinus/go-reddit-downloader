package main

// A simple program demonstrating the spinner component from the Bubbles
// component library.

import (
	"fmt"

	// "os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/zzLinus/GoTUITODOList/downloader"
	// "github.com/zzLinus/GoTUITODOList/extractor"
)

type errMsg error
type tickMsg time.Time

const (
	padding  = 2
	maxWidth = 80
)

// Add a purple, rectangular border
var style = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("63"))

// Set a rounded, yellow-on-purple border to the top and left
var anotherStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("204"))

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type model struct {
	spinner   spinner.Model
	border    lipgloss.Border
	loading   bool
	progress  progress.Model
	textInput textinput.Model
	quitting  bool
	err       error
}

func initialModel() model {
	s := spinner.New()
	ti := textinput.New()
	ti.Placeholder = "plz input some text"
	ti.Focus()
	ti.CharLimit = 80
	ti.Width = 20
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Italic(true)
	pro := progress.New(progress.WithDefaultGradient())
	return model{spinner: s, textInput: ti, loading: true, progress: pro}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(spinner.Tick)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Batch(tickCmd())
		case "enter":
			m.loading = false
			return m, nil
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
	pad := strings.Repeat(" ", padding)
	if m.err != nil {
		return m.err.Error()
	}
	if m.loading && !m.quitting {
		str += anotherStyle.Render(fmt.Sprintf("\n%s Loading forever...press q to quit\n",
			m.spinner.View()))
		str += "\n"
		str += anotherStyle.Render(lipgloss.NewStyle().Italic(true).Render("Hello, kitty."))
	}
	if !m.loading && !m.quitting {
		str += fmt.Sprintf("Show me what you got\n\n%s", m.textInput.View())
	}

	if m.quitting {
		return "\n" +
			pad + m.progress.View() + "\n\n"
	}
	return str
}

func main() {
	rowURL := "https://v.redd.it/8akffrc6fqx81/DASH_720.mp4"
	// url := extractor.Extractor(rowURL)
	err := downloader.Download(rowURL)
	if err != nil {
		panic(err)
	}
	p := tea.NewProgram(initialModel())
	// start CLI gorutian
	p.Start()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg { return tickMsg(t) })
}
