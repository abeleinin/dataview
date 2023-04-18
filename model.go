package main

import (
	"./table"
	// "github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

  cursorY = 0
  editing = false
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
  focusIndex int
	input      []textinput.Model
	cursorMode textinput.CursorMode
	table      table.Model
}

func (m model) Init() tea.Cmd { 
  	return tea.Batch(nil, tea.EnterAltScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
      if !editing {
			  return m, tea.Quit
      }
    case "j":
      if !editing {
        cursorY += 1
      }
    case "k":
      if !editing {
        cursorY -= 1
      }
    case "l":
      if !editing {
        currRows := m.table.Rows()
        currColumns := m.table.Columns()
        newRows := []table.Row{}
        for _, row := range currRows {
          row = append(row[1:], row[0])
          newRows = append(newRows, row)
        }
        newColumns := append(currColumns[1:], currColumns[0])
        m.table.SetRows(newRows)
        m.table.SetColumns(newColumns)
      }
    case "h":
      if !editing {
        currRows := m.table.Rows()
        currColumns := m.table.Columns()
        newRows := []table.Row{}
        for _, row := range currRows {
          row = append([]string{row[len(row)-1]}, row[:len(row)-1]...)
          newRows = append(newRows, row)
        }
        newColumns := append([]table.Column{currColumns[len(currColumns)-1]}, currColumns[:len(currColumns)-1]...)
        m.table.SetRows(newRows)
        m.table.SetColumns(newColumns)
      }
    case "i":
      if !editing {
        editing = true
        m.table.Blur()
        m.input[0].Focus()
      }
    case "esc":
      editing = false
      m.table.Focus()
      m.input[0].Blur()
      m.input[0].Reset()
    case "enter":
      editing = false
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.input) {
				return m, tea.Quit
			}

			if m.focusIndex > len(m.input) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.input)
			}

			cmds := make([]tea.Cmd, 1)
			cmds[0] = m.input[0].Focus()
			m.input[0].PromptStyle = focusedStyle
			m.input[0].TextStyle = focusedStyle
			m.input[0].Blur()
			m.input[0].PromptStyle = noStyle
			m.input[0].TextStyle = noStyle

      // Set new value at currRows[cursorY][0]
      currRows := m.table.Rows()
      currRows[cursorY][0] = m.input[0].Value()
      m.table.SetRows(currRows)
      m.table.Focus()
      m.input[0].Reset()

			return m, tea.Batch(cmds...)
    }
  }

  cmd = m.updateInputs(msg)

	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.input))

	for i := range m.input {
		m.input[i], cmds[i] = m.input[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func DataviewStyles() table.Styles {
  return table.Styles{
    Selected: lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")),
		Header:   lipgloss.NewStyle().Bold(true).Padding(0, 1),
		Cell:     lipgloss.NewStyle().Padding(0, 1),
  }
} 

