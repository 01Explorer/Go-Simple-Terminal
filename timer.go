package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/common-nighthawk/go-figure"
)

// Help

type keyMap struct {
  Pause key.Binding
  Reset key.Binding
  Help key.Binding
  Quit key.Binding
  Start key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
 return []key.Binding{k.Help, k.Quit} 
}

func (k keyMap) FullHelp() [][]key.Binding {
  return [][]key.Binding{
    {k.Pause, k.Reset},
    {k.Help, k.Quit},
  }
}

var keys = keyMap {
  Start: key.NewBinding(
    key.WithKeys("s"),
    key.WithHelp("S", "Start"),
  ),
  Pause: key.NewBinding(
    key.WithKeys("p"),
    key.WithHelp("P", "Pause"),
  ),
  Reset: key.NewBinding(
    key.WithKeys("r"),
    key.WithHelp("R", "Reset timer"),
  ),
  Help: key.NewBinding(
    key.WithKeys("?"),
    key.WithHelp("?", "Toggle help"),
  ),
  Quit: key.NewBinding(
    key.WithKeys("q", "esc", "ctrl+c"),
    key.WithHelp("q/esc/ctrl+c", "quit"),
  ),
}

var (
  titleStyle = lipgloss.NewStyle().
  PaddingBottom(10).
  Bold(true).
  Margin(1, 0)
)

type TickMsg struct{}

type StartStopMsg struct {
	running bool
  finished bool
}

type ResetMsg struct{}

type Model struct {
	d        time.Duration
	running  bool
  finished bool 
  err string
  keys keyMap
  help help.Model
  timerState TimerState
}


func New(state TimerState) Model {
	return Model{
    keys: keys,
    help: help.New(),
    timerState: state,
	}
}

func (m Model) Init() tea.Cmd {
  return nil
}

func (m Model) Start() tea.Cmd {
  if m.Running() {
   return nil 
  }

	return tea.Sequence(func() tea.Msg {
		return StartStopMsg{running: true, finished: false}
	}, tick(time.Second))
}

func (m Model) Stop() tea.Cmd {
	return func() tea.Msg {
		return StartStopMsg{running: false}
	}
}

func (m Model) Finish() tea.Cmd {
  return func() tea.Msg {
    return StartStopMsg{
      running: false,
      finished: true,
    }
  }
}

func (m Model) Toggle() tea.Cmd {
	if m.Running() {
    return m.Stop()
	}

  return m.Start()
}

func (m Model) Reset() tea.Cmd {
 return func () tea.Msg {
   return ResetMsg{}
 } 
}

func (m Model) Running() bool {
  return m.running
}

func (m Model) Elapsed() time.Duration {
  return m.d
}

func (m Model) View() string {
  var s string


  s = titleStyle.Render(m.timerState.title)
  timeToShow := m.timerState.max - m.Elapsed()

  timer := figure.NewColorFigure(fmt.Sprintf("%s",timeToShow), "slant", "gray", true)
  innerStyle := lipgloss.NewStyle().
  Border(lipgloss.HiddenBorder())
  
  figureString := innerStyle.Render(timer.ColorString())


  s += figureString
  if m.finished {
    s += "\n\nTime finished"
  }

  helpStyle := lipgloss.NewStyle().
  Border(lipgloss.HiddenBorder())
  helpView := helpStyle.Render(m.help.View(m.keys))
  return s + helpView
}

func tick(d time.Duration) tea.Cmd {
  return tea.Tick(d, func(_ time.Time) tea.Msg {
    return TickMsg{}
  })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.WindowSizeMsg:
    m.help.Width = msg.Width
  case StartStopMsg:
    m.running = msg.running
    m.finished = msg.finished
  case ResetMsg:
    m.d = 0
    if m.Running() {
      return m, m.Stop() 
    }
  case TickMsg:
    if !m.Running() {
     break
    }

    m.d += time.Second

    if m.d >= m.timerState.max {
     return m, m.Finish() 
    }

    return m, tick(time.Second)
  case tea.KeyMsg:
    switch {
    case key.Matches(msg, m.keys.Start):
      return m, m.Start()
    case key.Matches(msg, m.keys.Pause):
      return m, m.Toggle()
    case key.Matches(msg, m.keys.Reset):
      return m, m.Reset()
    case key.Matches(msg, m.keys.Help):
      m.help.ShowAll = !m.help.ShowAll
    }
  }

  return m, nil
}
