package niks

import (
	"fmt"
	"os"
	"sort"
)

// Lists build (dirs), ordered latest first. Dirs without status.json are ignored.
func ListBuilds() ([]Build, error) {
	root := buildsDir // FIXME
	ls, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}

	var bs []Build
	for _, l := range ls {
		if !l.IsDir() {
			continue
		}
		p := root + l.Name()
		fmt.Printf("P: %s\n", p)
		st := p + "/status.json"
		if _, err := os.ReadFile(st); err != nil {
			continue
		}
		bs = append(bs, Build{
			ID:         l.Name(),
			Path:       p,
			StdoutFile: p + "/stdout.txt",
			StatusFile: p + "/status.json",
		})
	}
	sort.Slice(bs, func(i, j int) bool { return bs[j].ID < bs[i].ID })
	return bs, nil
}
