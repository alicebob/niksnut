package niks

import (
	"encoding/json"
	"os/exec"
)

type (
	Config struct {
		Projects []Project `json:"projects"`
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
	// nix := fmt.Sprintf("with import %s; builtins.toJSON projects", f)
	// cmd := exec.Command("nix-instantiate", "--eval", "-E", nix)
	cmd := exec.Command(cmdNixInstantiate, "--eval", "--strict", "--json", f)
	// fmt.Printf("running: %s\n", cmd.String())
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	// fmt.Printf("config json: %s\n", out)

	var c Config
	return &c, json.Unmarshal(out, &c)
}
