package main
  
import (
  "bufio"
  "strings"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240")).
  Align(lipgloss.Center)

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { 
  	return tea.Batch(nil, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Batch(
				tea.Printf("Let's go to %s!", m.table.SelectedRow()[1]),
			)
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}

func main() {

  if len(os.Args) != 2 {
    fmt.Println("Usage: myprogram filename.txt")
    return
  }

  filename := os.Args[1]
  file, err := os.Open(filename)

  if err != nil {
    fmt.Println("Error opening file:", err)
    return
  }

  defer file.Close()

  // Inialize rows and columns
	rows := []table.Row{}
	columns := []table.Column{}

  scanner := bufio.NewScanner(file)
  inc := 1
  for scanner.Scan() {
    myString := scanner.Text()
    values := strings.Split(myString, ",")

    if inc == 1 { // set columns 
      for _, value := range values {
        columns = append(columns, table.Column{Title: value, Width: 5})
      }
    } else { // set rows
      rows = append(rows, values)
    }

    inc++
  }

  if err := scanner.Err(); err != nil {
    fmt.Println("Error reading file:", err)
    return
  }

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(30),
	)

  s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := model{t}
	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

