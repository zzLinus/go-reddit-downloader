package tuiapp

// this file is basicly hacking on bubble tea tui framework

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
	"github.com/muesli/termenv"
	"github.com/zzLinus/GoRedditDownloader/downloader"
	"github.com/zzLinus/GoRedditDownloader/extractor"
)

type errMsg error
type tickMsg time.Time
type processFinishedMsg time.Duration
type respMsg int

const (
	padding  = 2
	maxWidth = 80
)

var (
	videoDownloader *downloader.Downloader
	duration        time.Duration
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).PaddingLeft(4)
	cursorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	spinnerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).PaddingTop(2).PaddingLeft(4)
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).PaddingLeft(4)
	blurredText  = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	focusedText  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle      = lipgloss.NewStyle()
	helpStyple   = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).PaddingLeft(4).PaddingTop(2)
	color        = termenv.EnvColorProfile().Color
	// Title        = termenv.Style{}.Foreground(color("204")).Background(color("235")).Styled
	Title = termenv.Style{}.Background(color("#cc3700")).Foreground(color("#e8e6e3")).Styled

	focusedSubmit = focusedStyle.Copy().PaddingLeft(4).Render("[ Submit ]")
	blurredSubmit = blurredStyle.Copy().PaddingLeft(4).Render("[ Submit ]")
	focusedQuit   = focusedStyle.Copy().PaddingLeft(4).Render("[ Quit ]")
	blurredQuit   = blurredStyle.Copy().PaddingLeft(4).Render("[ Quit ]")
)

type model struct {
	spinner    spinner.Model
	progress   progress.Model
	textInput  textinput.Model
	cursorMode textinput.CursorMode
	focusIndex int8
	loading    bool
	quitting   bool
	err        error
	results    []result
	sub        chan extractor.SubscriptMsg
}

type result struct {
	duration time.Duration
	emoji    string
	msg      string
}

func New() *tea.Program {
	p := tea.NewProgram(initialModel())
	return p
}

func initialModel() model {

	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 80
	ti.TextStyle = focusedText
	ti.Placeholder = "pase an url that support by use"
	ti.CursorStyle = focusedStyle
	ti.PromptStyle = focusedStyle

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Italic(true)

	pro := progress.New(progress.WithDefaultGradient())

	return model{
		spinner:    s,
		textInput:  ti,
		loading:    false,
		progress:   pro,
		focusIndex: 0,
		results:    make([]result, 9),
		sub:        make(chan extractor.SubscriptMsg, 10),
	}
}

func activityLisenter(sub chan extractor.SubscriptMsg) tea.Cmd {
	return func() tea.Msg {
		return extractor.SubscriptMsg(<-sub)
	}
}

func (m model) Init() tea.Cmd {
	rand.Seed(time.Now().UTC().UnixNano())
	videoDownloader = downloader.New()
	return tea.Batch(textinput.Blink, tea.EnterAltScreen, activityLisenter(m.sub))
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
			return m, tea.Quit
		case "tab", "down":
			if m.focusIndex == 0 {
				m.focusIndex++
				m.textInput.CursorStyle = blurredStyle
				m.textInput.PromptStyle = blurredStyle
				m.textInput.TextStyle = blurredText
				m.textInput.Blur()
			} else if m.focusIndex == 1 {
				m.focusIndex++
			} else if m.focusIndex == 2 {
				m.focusIndex = 0
				m.textInput.CursorStyle = focusedStyle
				m.textInput.PromptStyle = focusedStyle
				m.textInput.TextStyle = focusedText
				m.textInput.Focus()
			}
			return m, nil
		case "shift+tab", "up":
			if m.focusIndex == 0 {
				m.focusIndex = 2
				m.textInput.CursorStyle = blurredStyle
				m.textInput.PromptStyle = blurredStyle
				m.textInput.TextStyle = blurredText
				m.textInput.Blur()
			} else if m.focusIndex == 1 {
				m.focusIndex--
				m.textInput.CursorStyle = focusedStyle
				m.textInput.PromptStyle = focusedStyle
				m.textInput.TextStyle = focusedText
				m.textInput.Focus()
			} else if m.focusIndex == 2 {
				m.focusIndex--
			}
			return m, nil
		case "enter":
			if m.focusIndex == 1 {
				m.loading = true
				return m, tea.Batch(
					m.spinner.Tick,
					func() tea.Msg {
						rowURL := m.textInput.Value()
						statusCode, err := videoDownloader.Download(rowURL, m.sub)
						if err != nil {
							panic(err)
						}
						return respMsg(statusCode)
					})
			}
			if m.focusIndex == 2 {
				return m, tea.Quit
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

	case extractor.SubscriptMsg:
		d := time.Duration(10)
		res := result{emoji: randomEmoji(), duration: d, msg: msg.Msg}
		m.results = append(m.results[1:], res)
		return m, activityLisenter(m.sub)

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
	var b strings.Builder

	if m.err != nil {
		return m.err.Error()
	}

	b.WriteString(fmt.Sprintf("\n\n           %s", Title("RedditDownloader")))
	b.WriteString("\n\n\n")
	if m.loading {
		b.WriteString(spinnerStyle.Render(fmt.Sprintf("%s Downloading content", m.spinner.View())))
		b.WriteString("\n\n")
		for _, res := range m.results {
			if res.duration == 0 {
				b.WriteString("    .................................\n")
			} else {
				b.WriteString(fmt.Sprintf("    %s %-25s %s\n", res.emoji, res.msg, res.duration))
			}
		}
	}

	if !m.loading {
		b.WriteString(cursorStyle.Render(fmt.Sprintf("%s", m.textInput.View())))
		b.WriteString("\n\n")
		sub := &blurredSubmit
		qui := &blurredQuit
		if m.focusIndex == 1 {
			sub = &focusedSubmit
		} else if m.focusIndex == 2 {
			qui = &focusedQuit
		}
		h := helpStyple.Render("use tab shift+tab or ↓ ↑ to control,enter to choose")
		b.WriteString(fmt.Sprintf("%s%s", *sub, *qui))
		b.WriteString(fmt.Sprintf("%s", h))
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
	emojis := []rune("🍤🧋🍡👾🦊🐯🤖🎏🍔🍥🎮🍕🥐🎄🕹️")
	return string(emojis[rand.Intn(len(emojis))])
}
