package main

import (
	"github.com/charmbracelet/lipgloss"
)

func CenterSquareWithText(appWidth, appHeight, outerWidth, outerHeight, innerWidth, innerHeight int, content string) string {
  outerStyle := lipgloss.NewStyle().
  Width(outerWidth).
  Height(outerHeight).
  Border(lipgloss.RoundedBorder()).
  BorderForeground(lipgloss.Color("#E95420"))

  innerRendered := lipgloss.Place(outerWidth, outerHeight, lipgloss.Center, lipgloss.Center, content)
  return outerStyle.Render(innerRendered)
}

func CenterTabsAboveSquare(appWidth, appHeight, squareWidth, squareHeight int, tabs, square string) string {
  squareLeft := (appWidth - squareWidth) / 2

  tabsStyle := lipgloss.NewStyle().
  MarginLeft(squareLeft)

  squareStyle := lipgloss.NewStyle().
  MarginLeft(squareLeft)

  combined := lipgloss.JoinVertical(
    lipgloss.Left,
    tabsStyle.Render(tabs),
    squareStyle.Render(square),
  )

  return combined
}
