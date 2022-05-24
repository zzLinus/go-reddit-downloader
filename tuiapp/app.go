package tuiapp

// this file is basicly working on bubble team tui framework

import (
	"fmt"
	"math/rand"
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

type processFinishedMsg time.Duration

const (
	padding  = 2
	maxWidth = 80
)

var (
	videoDownloader *downloader.Downloader
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).PaddingLeft(4)
	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	spinnerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).PaddingTop(2)
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).PaddingLeft(4)
	blurredText  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	focusedText  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle      = lipgloss.NewStyle()

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = blurredStyle.Copy().Render("[ Submit ]")
)

type model struct {
	spinner     spinner.Model
	progress    progress.Model
	textInput   textinput.Model
	cursorMode  textinput.CursorMode
	focusButton bool
	loading     bool
	quitting    bool
	err         error
	results     []result
}

type result struct {
	duration time.Duration
	emoji    string
}

func New() *tea.Program {
	p := tea.NewProgram(initialModel())
	return p
}

func initialModel() model {

	ti := textinput.New()
	ti.Placeholder = "pase an url that support by use"
	ti.Focus()
	ti.CharLimit = 80
	ti.TextStyle = focusedText
	ti.CursorStyle = focusedStyle
	ti.PromptStyle = focusedStyle

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Italic(true)

	pro := progress.New(progress.WithDefaultGradient())

	return model{
		spinner:     s,
		textInput:   ti,
		loading:     false,
		progress:    pro,
		focusButton: false,
		results:     make([]result, 6),
	}
}

func (m model) Init() tea.Cmd {
	rand.Seed(time.Now().UTC().UnixNano())
	videoDownloader = downloader.New()
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)

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
		case "tab", "down":
			m.focusButton = true
			m.textInput.TextStyle = blurredText
			m.textInput.CursorStyle = blurredStyle
			m.textInput.PromptStyle = blurredStyle
			m.textInput.Blur()
			return m, nil
		case "shift+tab", "up":
			m.focusButton = false
			m.textInput.TextStyle = focusedText
			m.textInput.CursorStyle = focusedStyle
			m.textInput.PromptStyle = focusedStyle
			m.textInput.Focus()
			return m, nil
		case "enter":
			if m.focusButton {
				m.loading = true
				return m, tea.Batch(
					m.spinner.Tick,
					func() tea.Msg {
						rowURL := m.textInput.Value()
						statusCode, err := videoDownloader.Download(rowURL)
						if err != nil {
							panic(err)
						}
						return respMsg(statusCode)
					},
					runPretendProcess)
			}
			return m, nil
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

	case processFinishedMsg:
		d := time.Duration(msg)
		res := result{emoji: randomEmoji(), duration: d}
		m.results = append(m.results[1:], res)
		return m, runPretendProcess

	default:
		var cmd tea.Cmd
		if m.loading && !m.quitting {
			m.spinner, cmd = m.spinner.Update(msg)
		}
		return m, cmd
	}
}

func (m model) View() string {
	var b strings.Builder
	if m.err != nil {
		return m.err.Error()
	}

	if m.loading {
		b.WriteString(spinnerStyle.Render(fmt.Sprintf("%s Downloading content", m.spinner.View())))
		b.WriteString("\n\n")
		for _, res := range m.results {
			if res.duration == 0 {
				b.WriteString(".............................\n")
			} else {
				b.WriteString(fmt.Sprintf("%s Fake Job finished in %s\n", res.emoji, res.duration))
			}
		}
	}

	if !m.loading {
		// str += cursorStyle.Render(fmt.Sprintf("%s", m.textInput.View()))
		b.WriteString(cursorStyle.Render(fmt.Sprintf("%s", m.textInput.View())))
		button := &blurredButton
		if m.focusButton == true {
			button = &focusedButton
		}
		fmt.Fprintf(&b, "\n\n%s\n\n", *button)
	}

	if m.quitting {
		pad := strings.Repeat(" ", padding)
		return "\n\n" + pad + m.progress.View() + "\n\n"
	}

	return b.String()
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg { return tickMsg(t) })
}

func randomEmoji() string {
	emojis := []rune("🍦🧋🍡🤠👾😭🦊🐯🦆🥨🎏🍔🍒🍥🎮📦🦁🐶🐸🍕🥐🧲🚒🥇🏆🌽")
	return string(emojis[rand.Intn(len(emojis))])
}

func runPretendProcess() tea.Msg {
	pause := time.Duration(rand.Int63n(899)+100) * time.Millisecond
	time.Sleep(pause)
	return processFinishedMsg(pause)
}
