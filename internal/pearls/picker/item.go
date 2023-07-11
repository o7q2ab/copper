package picker

import (
	"github.com/charmbracelet/bubbles/list"
)

type fileinfo struct {
	parent   string
	name     string
	size     int64
	isDir    bool
	isHidden bool
	isGoCode bool
}

func (i *fileinfo) Title() string       { return i.name }
func (i *fileinfo) Description() string { return "" }
func (i *fileinfo) FilterValue() string { return i.name }

func toBubble(in []*fileinfo) []list.Item {
	res := make([]list.Item, len(in))
	for i, f := range in {
		res[i] = f
	}
	return res
}
