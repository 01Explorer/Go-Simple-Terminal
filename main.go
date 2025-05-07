package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/creack/pty"
)

type model struct {
  choices []string
  cursor int
  selected map[int]struct{}
}

func initialModel()  model {
  return model{
    choices: []string{"Buy carrots", "Buy celery", "Buy kohlraby"},
    selected: make(map[int]struct{}),
  }
}

func (m model) Init() tea.Cmd {
  return nil 
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd){
  switch msg := msg.(type) {
 case tea.KeyMsg:
   switch msg.String() {
   case "ctrl+c", "q":
     return m, tea.Quit
  case "up", "k":
   if m.cursor > 0  {
     m.cursor--
   } 
 case "down", "j":
   if m.cursor < len(m.choices)-1 {
     m.cursor++
    
   }
 case "enter", " ":
   _, ok := m.selected[m.cursor]
   if ok {
     delete(m.selected, m.cursor)
   } else {
     m.selected[m.cursor] = struct{}{}
   }
   }
 }

return m, nil
 }

 func (m model) View()  string {
   s := "What should we buy at the market??\n\n"

   for i, choice := range m.choices {
     cursor := " "
     if m.cursor == i {
       cursor = ">"
     }

     checked := " " 
     if _, ok := m.selected[i]; ok {
       checked = "x"
     }
     s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
   }

   s += "\nPress q to quit.\n"

   return s
 }

func main() {
  docStyle := lipgloss.NewStyle().Padding(1, 2, 1, 2)
  physicalWidth, _, _ := term.GetSize(os.Stdout.Fd())
  doc := strings.Builder{}

  row := lipgloss.JoinHorizontal(
    lipgloss.Top,
    activeTab.Render("Lip Gloss"),
    tab.Render("Blush"),
    tab.Render("Eye Shadow"),
    tab.Render("Mascara"),
    tab.Render("Foundation"),
  )
  gap := tabGap.Render(strings.Repeat(" ", max(0, 80 - lipgloss.Width(row) - 2)))
  row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
  doc.WriteString(row + "\n\n")

  if physicalWidth > 0 {
    docStyle = docStyle.MaxWidth(physicalWidth)
  }

  fmt.Println(docStyle.Render(doc.String()))

  p := tea.NewProgram(New())
  if _, err := p.Run(); err != nil {
    fmt.Printf("Alas, there's been an error: %v", err)
    os.Exit(1)
  }
}
  
func oldMain() {
  cmd := exec.Command("/bin/bash")
  pty, err := pty.Start(cmd)
  if err != nil {
    panic(err)
  }
  defer pty.Close()

  reader := bufio.NewReader(pty)
  go func() {
   for {
     buf := make([]byte, 1024)
     n, err := reader.Read(buf)
     if err != nil {
       if err == io.EOF {
         return
       }
       panic(err)
     }
     os.Stdout.Write(buf[:n])
   } 
  }()

  go func() {
    scanner := bufio.NewScanner(os.Stdin) 
    for scanner.Scan(){
      input := scanner.Text()
      pty.Write([]byte(input + "\n"))
    }
  }()

  time.Sleep(10 * time.Second)
}
