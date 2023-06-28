package menu

import (
	"fmt"
	"runtime"
	"runtime/debug"

	tea "github.com/charmbracelet/bubbletea"
)

const info = `GitHub: https://github.com/o7q2ab/copper

runtime.Version():....... %s
runtime.GOOS:............ %s
runtime.GOARCH:.......... %s
runtime.NumCPU():........ %d
runtime.NumGoroutine():.. %d

[]debug.BuildSetting: 
%s
`

func New() *Model {
	return &Model{}
}

type Model struct{}

func (m *Model) Init() tea.Cmd                           { return nil }
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }

func (m *Model) View() string {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return "Call `debug.ReadBuildInfo()` failed."
	}

	buildSettigns := ""
	for _, s := range buildInfo.Settings {
		buildSettigns += fmt.Sprintf("\t%s=%s\n", s.Key, s.Value)
	}

	return fmt.Sprintf(
		info,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		runtime.NumGoroutine(),
		buildSettigns,
	)
}
