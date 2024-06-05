package niks

import (
	"os"
	"sort"
)

func validBuildDir(buildsDir, id string) bool {
	p := buildsDir + id + "/status.json"
	_, err := os.ReadFile(p)
	return err == nil
}

// Lists build (dirs), ordered latest first. Dirs without status.json are ignored.
func ListBuilds(buildsDir string) ([]Build, error) {
	ls, err := os.ReadDir(buildsDir)
	if err != nil {
		return nil, err
	}

	var bs []Build
	for _, l := range ls {
		if !l.IsDir() {
			continue
		}
		id := l.Name()
		if !validBuildDir(buildsDir, id) {
			continue
		}
		p := buildsDir + id
		bs = append(bs, Build{
			ID:   l.Name(),
			Path: p,
		})
	}
	sort.Slice(bs, func(i, j int) bool { return bs[j].ID < bs[i].ID })
	return bs, nil
}
