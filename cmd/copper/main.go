package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/o7q2ab/copper/internal/pearls/finder"
	"github.com/o7q2ab/copper/internal/pearls/picker"
)

const version = "day-4"

func main() {
	fmt.Printf("copper %s\n", version)

	pickerModel, err := picker.New()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	finderModel := finder.New()

	model := &copper{
		stage: pickerModel,

		picker: pickerModel,
		finder: finderModel,
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

type copper struct {
	stage tea.Model

	picker *picker.Model
	finder *finder.Model
}

func (c *copper) Init() tea.Cmd { return nil }

func (c *copper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return c, tea.Quit

		case "t":
			switch c.stage.(type) {
			case *picker.Model:
				c.stage = c.finder

			case *finder.Model:
				c.stage = c.picker
			}
		}
	}

	switch c.stage.(type) {
	case *picker.Model:
		if selected := c.finder.GetSelected(); selected != "" {
			c.picker.SetCurrent(filepath.Dir(selected))
		}

	case *finder.Model:
		c.finder.SetCurrent(c.picker.GetCurrent())

	}
	_, cmd := c.stage.Update(msg)

	return c, cmd
}

func (c *copper) View() string {
	s := c.stage.View()
	s += "\nPress q to quit.\n"
	return s
}
