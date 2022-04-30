// package main
//
// // A simple program demonstrating the spinner component from the Bubbles
// // component library.
//
// import (
// 	"fmt"
// 	"os"
// 	"strings"
// 	"time"
//
// 	"github.com/charmbracelet/bubbles/progress"
// 	"github.com/charmbracelet/bubbles/spinner"
// 	"github.com/charmbracelet/bubbles/textinput"
// 	tea "github.com/charmbracelet/bubbletea"
// 	"github.com/charmbracelet/lipgloss"
// )
//
// type errMsg error
// type tickMsg time.Time
//
// const (
// 	padding  = 2
// 	maxWidth = 80
// )
//
// var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
//
// type model struct {
// 	spinner   spinner.Model
// 	border    lipgloss.Border
// 	loading   bool
// 	progress  progress.Model
// 	textInput textinput.Model
// 	quitting  bool
// 	err       error
// }
//
// func initialModel() model {
// 	s := spinner.New()
// 	ti := textinput.New()
// 	ti.Placeholder = "plz input some text"
// 	ti.Focus()
// 	ti.CharLimit = 80
// 	ti.Width = 20
// 	s.Spinner = spinner.Dot
// 	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Italic(true)
// 	pro := progress.New(progress.WithDefaultGradient())
// 	return model{spinner: s, textInput: ti, loading: true, progress: pro}
// }
//
// func (m model) Init() tea.Cmd {
// 	return m.spinner.Tick
// }
//
// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
//
// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "q", "esc", "ctrl+c":
// 			m.quitting = true
// 			return m, tea.Batch(tickCmd())
// 		case "enter":
// 			m.loading = false
// 			return m, nil
// 		default:
// 			var cmd tea.Cmd
// 			if !m.loading {
// 				m.textInput, cmd = m.textInput.Update(msg)
// 			}
// 			return m, cmd
// 		}
//
// 	case tea.WindowSizeMsg:
// 		m.progress.Width = msg.Width - padding*2 - 4
// 		if m.progress.Width > maxWidth {
// 			m.progress.Width = maxWidth
// 		}
// 		return m, nil
//
// 	case errMsg:
// 		m.err = msg
// 		return m, nil
//
// 	case tickMsg:
// 		if m.progress.Percent() == 1.0 {
// 			return m, tea.Quit
// 		}
// 		cmd := m.progress.IncrPercent(0.25)
// 		return m, cmd
//
// 	default:
// 		var cmd tea.Cmd
// 		if m.loading && !m.quitting{
// 			m.spinner, cmd = m.spinner.Update(msg)
// 		}
// 		if m.quitting {
// 			tickCmd()
// 		}
// 		return m, cmd
// 	}
// }
//
// func (m model) View() string {
// 	var str string = ""
// 	pad := strings.Repeat(" ", padding)
// 	if m.err != nil {
// 		return m.err.Error()
// 	}
// 	if m.loading && !m.quitting {
// 		str += fmt.Sprintf("\n\n   %s Loading forever...press q to quit\n\n",
// 			m.spinner.View())
// 	}
// 	if !m.loading && !m.quitting {
// 		str += fmt.Sprintf("Show me what you got\n\n%s", m.textInput.View())
// 	}
//
// 	if m.quitting {
// 		return "\n" +
// 			pad + m.progress.View() + "\n\n"
// 	}
//
// 	return str
// }
//
// func main() {
// 	p := tea.NewProgram(initialModel())
// 	if err := p.Start(); err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}
// }
//
// func tickCmd() tea.Cmd {
// 	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg { return tickMsg(t) })
// }


package main

// A simple example that shows how to render an animated progress bar. In this
// example we bump the progress by 25% every two seconds, animating our
// progress bar to its new target state.
//
// It's also possible to render a progress bar in a more static fashion without
// transitions. For details on that approach see the progress-static example.

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

func main() {
	m := model{
		progress: progress.New(progress.WithDefaultGradient()),
	}

	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

type tickMsg time.Time

type model struct {
	progress progress.Model
}

func (_ model) Init() tea.Cmd {
	return tickCmd()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}

		// Note that you can also use progress.Model.SetPercent to set the
		// percentage value explicitly, too.
		cmd := m.progress.IncrPercent(0.25)
		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (e model) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + e.progress.View() + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
