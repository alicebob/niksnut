package niks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
)

type (
	Config struct {
		Projects []Project `json:"projects"`
	}

	Project struct {
		ID        string   `json:"id"`
		Name      string   `json:"name"`
		Category  string   `json:"category"`
		Git       string   `json:"git"`
		Attribute string   `json:"attribute"`
		Packages  []string `json:"packages"`
		Post      string   `json:"post"`
	}
)

func ReadConfig(f string) (*Config, error) {
	stderr := &bytes.Buffer{}
	cmd := exec.Command(cmdNixInstantiate, "--strict", "--json", "--read-write-mode", "--eval", f)
	cmd.Stderr = stderr
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("config: %s\n", stderr.String())
		return nil, err
	}

	var c Config
	if err := json.Unmarshal(out, &c); err != nil {
		return nil, err
	}
	return &c, checkConfig(&c)
}

func checkConfig(c *Config) error {
	seen := map[string]bool{}
	for _, p := range c.Projects {
		if !validID(p.ID) {
			return fmt.Errorf("invalid project ID: %q", p.ID)
		}
		if seen[p.ID] {
			return fmt.Errorf("repeated project ID: %s", p.ID)
		}
		seen[p.ID] = true
	}

	return nil
}

func validID(id string) bool {
	match, err := regexp.MatchString(`^[a-zA-Z0-9.-]+$`, id)
	if err != nil {
		panic(err)
	}
	return match
}
