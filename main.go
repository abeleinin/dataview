package main
  
import (
  "bufio"
  "strings"
	"fmt"
	"os"

	"./table"
	// "github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

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
  inc := 0
  for scanner.Scan() {
    myString := scanner.Text()
    values := strings.Split(myString, ",")

    if inc == 0 { // set columns 
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

  s := DataviewStyles()

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

	m := model{
		input: make([]textinput.Model, 1),
    table: t,
	}

	var txt textinput.Model
	for i := range m.input {
		txt = textinput.New()
		txt.CursorStyle = cursorStyle
		txt.CharLimit = 32

		switch i {
		case 0:
			txt.Placeholder = "Enter Here"
			txt.PromptStyle = focusedStyle
			txt.TextStyle = focusedStyle
		}

		m.input[0] = txt
	}

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

