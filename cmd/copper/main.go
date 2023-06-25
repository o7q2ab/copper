package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/o7q2ab/copper/internal/pearls/picker"
)

const version = "day-3"

func main() {
	fmt.Printf("copper %s\n", version)

	pickerModel, err := picker.New()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	model := &copper{
		stage: pickerModel,
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

type copper struct {
	stage tea.Model
}

func (c *copper) Init() tea.Cmd { return nil }

func (c *copper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return c, tea.Quit
		}
	}

	_, cmd := c.stage.Update(msg)

	return c, cmd
}

func (c *copper) View() string {
	s := c.stage.View()
	s += "\nPress q to quit.\n"
	return s
}
