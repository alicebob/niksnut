package niks

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func validBuildDir(p string) bool {
	if !strings.Contains(p, "/runs/") {
		return false
	}
	if strings.Contains(p, "..") {
		return false
	}
	f := fmt.Sprintf("%s/status.json", p)
	_, err := os.ReadFile(f)
	return err == nil
}

// Lists build (dirs), ordered latest first. Dirs without status.json are ignored.
func ListBuilds(root string) ([]Build, error) {
	ls, err := os.ReadDir(root + "/runs/")
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			return nil, nil
		}
		return nil, err
	}

	var bs []Build
	for _, l := range ls {
		if !l.IsDir() {
			continue
		}
		id := l.Name()
		p := buildPath(root, id)
		if !validBuildDir(p) {
			continue
		}
		bs = append(bs, Build{
			ID:   l.Name(),
			Path: p,
		})
	}
	sort.Slice(bs, func(i, j int) bool { return bs[j].ID < bs[i].ID })
	return bs, nil
}
