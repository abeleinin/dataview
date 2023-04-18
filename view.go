package main

import (
  "fmt"
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

	button := &blurredButton
	if m.focusIndex == len(m.input) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(m.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style)"))

  return baseStyle.Render(m.table.View()) + "\n" + b.String() + "\n"
}
