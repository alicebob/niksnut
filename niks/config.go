package niks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
)

type (
	Config struct {
		Projects []Project `json:"projects"`
	}

	Project struct {
		ID        string   `json:"id"`
		Name      string   `json:"name"`
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
	// fmt.Printf("running: %s\n", cmd.String())
	out, err := cmd.Output()
	if err != nil {
		fmt.Printf("config: %s\n", stderr.String())
		return nil, err
	}
	// fmt.Printf("config json: %s\n", out)

	var c Config
	return &c, json.Unmarshal(out, &c)
}
