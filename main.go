package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainModel struct {
  Tabs []string
  TabContent []Model
  activeTab int
  width int
  height int
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
    case "left", "n":
      m.activeTab = max(m.activeTab-1, 0)
    }
   case tea.WindowSizeMsg:
     m.width = msg.Width
     m.height = msg.Height
  }
  model, cmd := m.TabContent[m.activeTab].Update(msg)
  m.TabContent[m.activeTab] = model.(Model)
  return m, cmd 
}

var (
	docStyle          = lipgloss.NewStyle().Padding(1, 2, 1, 2)
	highlightColor    = lipgloss.Color("#E95420")
  inactiveTabStyle  = lipgloss.NewStyle()
	activeTabStyle    = inactiveTabStyle.Background(highlightColor)
	windowStyle       = lipgloss.NewStyle().BorderForeground(highlightColor).Padding(2, 0).Align(lipgloss.Center).Border(lipgloss.NormalBorder()).UnsetBorderTop()

)

func (m MainModel) View() string {
  outerWidth := int(float64(m.width) * 0.8)
  outerHeight := int(float64(m.height) * 0.8)
  innerWidth := outerWidth - 4
  innerHeight := 15

  var renderedTabs []string

  for i, t := range m.Tabs {
    var style lipgloss.Style
    isActive := i == m.activeTab

    style = inactiveTabStyle
    if isActive {
     style = activeTabStyle 
    }

    style = style.Border(lipgloss.RoundedBorder()).
    Width((innerWidth - 1)/3).
    Height(2).
    Align(lipgloss.Center, lipgloss.Center).
    BorderForeground(highlightColor)
    
    renderedTabs = append(renderedTabs, style.Render(t))
  }

  row := lipgloss.JoinHorizontal(lipgloss.Left, renderedTabs...)
  activeTab := m.TabContent[m.activeTab].View()

  mainSquare := CenterSquareWithText(m.width, m.height, outerWidth, outerHeight, innerWidth, innerHeight, activeTab)

  squareWithTab := CenterTabsAboveSquare(m.width, m.height, outerWidth, outerHeight, row, mainSquare)

  return squareWithTab
}


func main() {
  tabs := []string{"Pomodoro", "Short Break", "Long Break"}
  tabContent := []Model {
    New(Pomodoro()),
    New(ShortBreak()),
    New(LongBreak()),
  }
  m := MainModel{Tabs: tabs, TabContent: tabContent}
  if _, err := tea.NewProgram(m).Run(); err != nil {
    fmt.Println("Error running program:", err)
    os.Exit(1)
  }
}
  
