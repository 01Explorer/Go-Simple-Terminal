package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainModel struct {
  Tabs []string
  TabContent []string
  activeTab int
}

func (m MainModel) Init() tea.Cmd {
 return nil 
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd){
  switch msg := msg.(type) {
  case tea.KeyMsg:
    switch keypress := msg.String(); keypress {
    case "q", "ctrl+c":
      return m, tea.Quit
    case "right", "l":
      m.activeTab = min(m.activeTab+1, len(m.Tabs)-1)
      return m, nil
    case "left", "n":
      m.activeTab = max(m.activeTab-1, 0)
      return m, nil
    }
  }
  return m, nil
}

func tabBorderWithBottom(left, middle, right string) lipgloss.Border {
  border := lipgloss.RoundedBorder()
  border.BottomLeft = left
  border.Bottom = middle
  border.BottomRight = right
  return border
}

var (
  inactiveTabBorder = tabBorderWithBottom("┴", "─", "┴")
	activeTabBorder   = tabBorderWithBottom("┘", " ", "└")
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	inactiveTabStyle  = lipgloss.NewStyle().Border(inactiveTabBorder, true).BorderForeground(highlightColor).Padding(0, 1)
	activeTabStyle    = inactiveTabStyle.Border(activeTabBorder, true)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()

)

func (m MainModel) View() string {
  doc := strings.Builder{}

  var renderedTabs []string

  for i, t := range m.Tabs {
    var style lipgloss.Style
    isFirst, isLast, isActive := i == 0, i == len(m.Tabs) - 1, i == m.activeTab

    style = inactiveTabStyle
    if isActive {
     style = activeTabStyle 
    }

    border, _, _, _, _ := style.GetBorder()
    if (isFirst || isLast) && isActive {
      border.BottomLeft = "|"
    } else if isFirst && !isActive {
    	border.BottomLeft = "├"
    } else if isLast && !isActive {
      border.BottomRight = "┤"
    }
    style = style.Border(border)
    renderedTabs = append(renderedTabs, style.Render(t))
  }

  row := lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
  doc.WriteString(row)
  doc.WriteString("\n")
  doc.WriteString(windowStyle.Width((lipgloss.Width(row)-windowStyle.GetHorizontalFrameSize())).Render(m.TabContent[m.activeTab]))
  return docStyle.Render(doc.String())
}


func main() {
  tabs := []string{"Pomodoro", "Short Rest", "Long Rest"}
  tabContent := []string{"Pomodoro configured Timer", "Shor Rest Timer", "Long Rest Timer"}
  m := MainModel{Tabs: tabs, TabContent: tabContent}
  if _, err := tea.NewProgram(m).Run(); err != nil {
    fmt.Println("Error running program:", err)
    os.Exit(1)
  }
}
  
