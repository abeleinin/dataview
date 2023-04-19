package main
  
import (
  "bufio"
  "strings"
	"fmt"
	"os"
  "golang.org/x/term"

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

  // Detech terminal height
  _, height, err := term.GetSize(0)
  // columnWidth := 0
  
  if err != nil {
    panic(err)
  }

  for scanner.Scan() {
    myString := scanner.Text()
    values := strings.Split(myString, ",")
    
    // set columns 
    if inc == 0 { 
      for _, value := range values {
        columns = append(columns, table.Column{Title: value, Width: 5})
      }
    // set rows
    } else { 
      rows = append(rows, values)
    }
    inc++
  }

  if err := scanner.Err(); err != nil {
    fmt.Println("Error reading file:", err)
    return
  }

  // 80% of the terminal height
  scale := 0.8 
  tableHeight := int(float64(height) * scale)
  // tableWidth := 0
  // if columnWidth < tableWidth {
  //   tableWidth = columnWidth
  // } else {
  //   tableWidth = int(float64(width) * scale)
  // }

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(tableHeight),
    // table.WithWidth(tableWidth),
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

