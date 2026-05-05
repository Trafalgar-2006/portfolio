package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// Project maps to the projects section of content.yaml
type Project struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Tags        []string `yaml:"tags"`
	Status      string   `yaml:"status"`
	GitHubURL   string   `yaml:"github"`
	Highlight   string   `yaml:"highlight"`
}

// Contact maps to the contacts section of content.yaml
type Contact struct {
	Icon  string `yaml:"icon"`
	Label string `yaml:"label"`
	Value string `yaml:"value"`
}

// Content holds the entire loaded content.yaml
type Content struct {
	Projects []Project `yaml:"projects"`
	Contacts []Contact `yaml:"contacts"`
}

// Loaded holds the parsed content — populated by Load()
var Loaded *Content

// Load reads and parses content.yaml from the given path
func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var c Content
	if err := yaml.Unmarshal(data, &c); err != nil {
		return err
	}
	Loaded = &c
	log.Printf("Loaded content.yaml: %d projects, %d contacts", len(c.Projects), len(c.Contacts))
	return nil
}
