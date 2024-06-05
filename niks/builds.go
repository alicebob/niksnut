package niks

import (
	"os"
	"sort"
)

func validBuildDir(id string) bool {
	root := buildsDir // FIXME
	p := root + id + "/status.json"
	_, err := os.ReadFile(p)
	return err == nil
}

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
		id := l.Name()
		if !validBuildDir(id) {
			continue
		}
		p := root + id
		bs = append(bs, Build{
			ID:   l.Name(),
			Path: p,
		})
	}
	sort.Slice(bs, func(i, j int) bool { return bs[j].ID < bs[i].ID })
	return bs, nil
}
