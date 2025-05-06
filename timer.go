package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

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
}

func New() Model {
	return Model{
		Interval: time.Second,
    max: 1 * time.Minute,
    title: "Pomodoro",
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
  s += fmt.Sprintf("Ongoing Time -> %s / %s", m.d, m.max) 
  if m.finished {
    s += "\n\nTime finished"
  }
  return s
}

func tick(d time.Duration) tea.Cmd {
  return tea.Tick(d, func(_ time.Time) tea.Msg {
    return TickMsg{}
  })
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
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
    switch msg.Type {
    case tea.KeyCtrlQ:
      return m, tea.Quit
    case tea.KeyCtrlP:
      return m, m.Toggle()
    case tea.KeyCtrlR:
      return m, m.Reset()
    }
  }

  return m, nil
}
