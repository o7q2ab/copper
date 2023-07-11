package picker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/o7q2ab/copper/internal/gomod"
)

const (
	listh = 20

	gomodw = 40
)

func New() (*Model, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd(): %w", err)
	}

	delegate := list.DefaultDelegate{
		ShowDescription: false,
		Styles:          list.NewDefaultItemStyles(),
	}
	delegate.SetHeight(1)

	filelist := list.New(nil, delegate, 0, listh)
	filelist.SetShowHelp(false)

	model := &Model{
		list: filelist,
	}
	model.readDir(cwd)

	return model, nil
}

type Model struct {
	cwd   string
	gomod string
	list  list.Model
	err   error
}

func (m *Model) GetCurrent() string    { return m.cwd }
func (m *Model) SetCurrent(cwd string) { m.readDir(cwd) }

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width - gomodw)

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
			m.readDir(filepath.Join(m.cwd, choice.name))

		case "backspace", "h":
			m.readDir(filepath.Dir(m.cwd))

		case "c":
			_ = exec.Command("code", m.cwd).Run()

		}
	}

	return m, nil
}

func (m *Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("error: %v\n\nPress backspace to go back.\nPress q to quit.\n", m.err)
	}

	if m.gomod == "" {
		return m.list.View()
	}

	var s string
	gomodInfo, err := gomod.Read(m.gomod)
	if err != nil {
		s = fmt.Sprintf("Module: (error) %s\n\n", err)
	}
	s = fmt.Sprintf("%s\n\ndirect dependencies: %d\nindirect dependencies: %d",
		gomodInfo.Path, gomodInfo.DirectDepsCnt, gomodInfo.IndirectDepsCnt) + s

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.list.View(),
		lipgloss.NewStyle().Width(gomodw).Padding(1, 1, 1, 1).Foreground(lipgloss.Color("#EFEEB4")).Background(lipgloss.Color("#454D66")).Render(s),
	)
}

func (m *Model) readDir(path string) {
	m.cwd = path
	m.gomod = ""

	files, err := os.ReadDir(path)
	if err != nil {
		m.err = err
		return
	}

	res := make([]*fileinfo, len(files))

	shift := 0

	for _, f := range files {
		if !f.Type().IsDir() {
			continue
		}

		item := &fileinfo{
			parent:   path,
			name:     f.Name(),
			isDir:    true,
			isHidden: strings.HasPrefix(f.Name(), "."),
			isGoCode: false,
			size:     -1,
		}

		if info, err := f.Info(); err == nil {
			item.size = info.Size()
		}

		res[shift] = item
		shift++
	}

	for _, f := range files {
		if f.Type().IsDir() {
			continue
		}

		item := &fileinfo{
			parent:   path,
			name:     f.Name(),
			isDir:    false,
			isHidden: strings.HasPrefix(f.Name(), "."),
			isGoCode: strings.HasSuffix(f.Name(), ".go"),
			size:     -1,
		}

		if info, err := f.Info(); err == nil {
			item.size = info.Size()
		}

		res[shift] = item
		shift++

		if item.name == "go.mod" {
			m.gomod = filepath.Join(path, item.name)
		}
	}

	m.list.Title = path
	_ = m.list.SetItems(toBubble(res))
	m.list.ResetSelected()
}
