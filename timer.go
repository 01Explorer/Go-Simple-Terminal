package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
  Pause key.Binding
  Reset key.Binding
  Help key.Binding
  Quit key.Binding
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

type TickMsg struct{}

type StartStopMsg struct {
	running bool
  finished bool
}

type ResetMsg struct{}

type Model struct {
	d        time.Duration
  max time.Duration
	running  bool
  finished bool 
	Interval time.Duration
  err string
  title string
  keys keyMap
  help help.Model
}

func New() Model {
	return Model{
		Interval: time.Second,
    max: 1 * time.Minute,
    title: "Pomodoro",
    keys: keys,
    help: help.New(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.Start()
}

func (m Model) Start() tea.Cmd {
	return tea.Sequence(func() tea.Msg {
		return StartStopMsg{running: true}
	}, tick(m.Interval))
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
  s = fmt.Sprintf("%s\n\n", m.title)
  s += fmt.Sprintf("Ongoing Time -> %s / %s\n\n", m.d, m.max) 
  if m.finished {
    s += "\n\nTime finished"
  }

  helpView := m.help.View(m.keys)
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
  case TickMsg:
    if !m.Running() {
     break
    }

    m.d += m.Interval

    if m.d >= m.max {
     return m, m.Finish() 
    }

    return m, tick(m.Interval)
  case tea.KeyMsg:
    switch {
    case key.Matches(msg, m.keys.Pause):
      return m, m.Toggle()
    case key.Matches(msg, m.keys.Reset):
      return m, m.Reset()
    case key.Matches(msg, m.keys.Help):
      m.help.ShowAll = !m.help.ShowAll
    case key.Matches(msg, m.keys.Quit):
      return m, tea.Quit
    }
  }

  return m, nil
}
