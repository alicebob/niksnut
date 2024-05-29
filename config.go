package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type (
	Config struct {
		Projects []Project
	}

	Project struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Git       string `json:"git"`
		Attribute string `json:"attribute"`
		Post      string `json:"post"`
	}
)

func ReadConfig(f string) (*Config, error) {
	nix := fmt.Sprintf("with import %s; builtins.toJSON projects", f)
	cmd := exec.Command("nix-instantiate", "--eval", "-E", nix)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	var pr []Project
	if err := json.Unmarshal(out, &pr); err != nil {
		return nil, err
	}
	return &Config{Projects: pr}, nil
}
