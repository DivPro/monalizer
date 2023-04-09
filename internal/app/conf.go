package app

import (
	"fmt"
	"github.com/divpro/monalizer/internal/render"
	"github.com/divpro/monalizer/internal/source"
	"gopkg.in/yaml.v3"
	"os"
)

type Conf struct {
	Source    *source.Conf `yaml:"source"`
	Render    *render.Conf `yaml:"render"`
	Whitelist []string     `yaml:"whitelist"`
}

func Load(filePath string) (*Conf, error) {
	c := new(Conf)
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("read app file [ %s ]: %w", filePath, err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return nil, fmt.Errorf("parse app file [ %s ]: %w", filePath, err)
	}

	return c, nil
}
