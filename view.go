package main

import (
  "strings"
)

func (m model) View() string {
  var b strings.Builder

	for i := range m.input {
		b.WriteString(m.input[i].View())
		if i < len(m.input)-1 {
			b.WriteRune('\n')
		}
	}

  return baseStyle.Render(m.table.View()) + "\n" + b.String() + "\n"
}
