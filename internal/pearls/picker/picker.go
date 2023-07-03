package picker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/o7q2ab/copper/internal/gomod"
)

const height = 20

func New() (*Model, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd(): %w", err)
	}

	choices, err := readDir(cwd)
	if err != nil {
		return nil, fmt.Errorf("readDir(%s): %v", cwd, err)
	}

	delegate := list.DefaultDelegate{
		ShowDescription: false,
		Styles:          list.NewDefaultItemStyles(),
	}
	delegate.SetHeight(1)

	filelist := list.New(toBubble(choices), delegate, 0, min(height, len(choices)+5))
	filelist.SetShowHelp(false)
	filelist.Title = cwd

	model := &Model{
		list: filelist,
	}

	return model, nil
}

type Model struct {
	list list.Model
	err  error
}

func (m *Model) GetCurrent() string    { return m.list.SelectedItem().(*fileinfo).parent }
func (m *Model) SetCurrent(cwd string) { m.refresh(cwd) }

func (m *Model) refresh(path string) {
	choices, err := readDir(path)
	m.err = err
	m.list.Title = path
	_ = m.list.SetItems(toBubble(choices))
	m.list.ResetSelected()
	m.list.SetHeight(min(height, len(choices)+5))
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)

	case tea.KeyMsg:
		switch msg.String() {

		case "up", "k":
			m.list.CursorUp()

		case "down", "j":
			m.list.CursorDown()

		case "enter", "l":
			choice := m.list.SelectedItem().(*fileinfo)
			if !choice.isDir {
				break
			}
			m.refresh(filepath.Join(choice.parent, choice.name))

		case "backspace", "h":
			choice := m.list.SelectedItem().(*fileinfo)
			m.refresh(filepath.Dir(choice.parent))

		case "c":
			choice := m.list.SelectedItem().(*fileinfo)
			_ = exec.Command("code", choice.parent).Run()

		}
	}

	return m, nil
}

func (m *Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("error: %v\n\nPress backspace to go back.\nPress q to quit.\n", m.err)
	}

	var s string
	for _, itm := range m.list.Items() {
		if choice := itm.(*fileinfo); choice.name == "go.mod" {
			gomodInfo, err := gomod.Read(filepath.Join(choice.parent, choice.name))
			if err != nil {
				s = fmt.Sprintf("Module: (error) %s\n\n", err)
			}
			s = fmt.Sprintf("Module: %s\n\tdirect dependencies: %d\n\tindirect dependencies: %d\n\n",
				gomodInfo.Path, gomodInfo.DirectDepsCnt, gomodInfo.IndirectDepsCnt) + s
		}
	}

	return s + "\n" + m.list.View()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
