package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Cache for making decisions as quick as possible
type Cache struct {
	// OpenIssues stores key as user/repo name and value as the
	// open issues count
	OpenIssues map[string]int       `json:"open_issues"`
	Issues     map[string][]GHIssue `json:"issues"`
	XToken     string               `json:"xtoken"`
}

// Config is the config file for focus cli
type Config struct {
	Cache            `json:"cache"`
	Editor           string `json:"editor" toml:"editor"`
	CurrentMilestone string `json:"current_milestone" toml:"current_milestone"`
}

// GetConfig returns a config file which is located inside ~/.focus
// if no focus config is found it creates a default config in this
// path and returns it
func GetConfig() (Config, error) {
	home, _ := os.UserHomeDir()
	configPath := filepath.Join(home, ".focus")
	if f, _ := os.Stat(configPath); f == nil {
		c := Config{
			Editor: "nano",
		}
		b, err := json.Marshal(&c)
		if err != nil {
			return Config{}, err
		}

		err = ioutil.WriteFile(configPath, b, 0755)
		if err != nil {
			return Config{}, nil
		}

		return c, nil
	}

	var c Config
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		return c, err
	}

	err = json.Unmarshal(b, &c)
	if err != nil {
		return c, err
	}

	return c, nil
}
