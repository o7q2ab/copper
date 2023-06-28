package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/o7q2ab/copper/internal/pearls/finder"
	"github.com/o7q2ab/copper/internal/pearls/menu"
	"github.com/o7q2ab/copper/internal/pearls/picker"
)

const version = "copper day-11"

const (
	color1 = "#454D66"
	color2 = "#309975"
	color3 = "#58B368"
	color4 = "#DAD873"
	color5 = "#EFEEB4"
)

func main() {
	pickerModel, err := picker.New()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	model := &copper{
		stage: pickerModel,

		menu:   menu.New(),
		picker: pickerModel,
		finder: finder.New(),
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
}

type copper struct {
	windowWidth int

	stage tea.Model

	menu   *menu.Model
	picker *picker.Model
	finder *finder.Model
}

func (c *copper) Init() tea.Cmd { return nil }

func (c *copper) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		c.windowWidth = msg.Width

	case tea.KeyMsg:
		switch msg.String() {

		case "ctrl+c", "q":
			return c, tea.Quit

		case "esc":
			switch c.stage.(type) {
			case *menu.Model:
				c.stage = c.picker
			case *picker.Model:
				c.stage = c.menu
			case *finder.Model:
				c.stage = c.picker
			}

		case "f":
			if _, ok := c.stage.(*picker.Model); ok {
				c.finder.SetCurrent(c.picker.GetCurrent())
				c.stage = c.finder
			}

		case "enter":
			if _, ok := c.stage.(*finder.Model); ok {
				if selected := c.finder.GetSelected(); selected != "" {
					c.picker.SetCurrent(filepath.Dir(selected))
				}
				c.stage = c.picker
				return c, nil
			}
		}
	}

	_, cmd := c.stage.Update(msg)
	return c, cmd
}

func (c *copper) View() string {
	var style = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(color5)).
		Background(lipgloss.Color(color1)).
		Width(c.windowWidth).
		Align(lipgloss.Center)

	s := style.Render(version) + "\n"
	s += c.stage.View()
	s += "\nPress q to quit.\n"
	return s
}
