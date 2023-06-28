package gomod

import (
	"os"

	"golang.org/x/mod/modfile"
)

type Info struct {
	Path            string
	DirectDepsCnt   int
	IndirectDepsCnt int
}

func Read(path string) (*Info, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	f, err := modfile.Parse(path, content, nil)
	if err != nil {
		return nil, err
	}

	directDeps, indirectDeps := 0, 0
	for _, d := range f.Require {
		if d.Indirect {
			indirectDeps++
		} else {
			directDeps++
		}
	}

	return &Info{
		Path:            f.Module.Mod.Path,
		DirectDepsCnt:   directDeps,
		IndirectDepsCnt: indirectDeps,
	}, nil
}
