package picker

import (
	"os"
	"strings"

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

func readDir(path string) ([]*fileinfo, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	res := make([]*fileinfo, len(files))

	for i, f := range files {
		item := &fileinfo{
			parent:   path,
			name:     f.Name(),
			isDir:    f.Type().IsDir(),
			isHidden: strings.HasPrefix(f.Name(), "."),
			isGoCode: strings.HasSuffix(f.Name(), ".go"),
			size:     -1,
		}

		if info, err := f.Info(); err == nil {
			item.size = info.Size()
		}

		res[i] = item
	}
	return res, nil
}

func toBubble(in []*fileinfo) []list.Item {
	res := make([]list.Item, len(in))
	for i, f := range in {
		res[i] = f
	}
	return res
}
