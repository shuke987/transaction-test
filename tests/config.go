package tests

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
)

const (
	AssertFileName   = "assert.json"
	SetupFileName    = "setup.sql"
	TearDownFileName = "teardown.sql"
)

type Config struct {
	SQLFiles     []string
	AssertFile   string
	SetupFile    string
	TeardownFile string
}

// find all case in dir and sub directories, recursively.
func findAllConfigs(dir string) ([]*Config, error) {
	var cfgs []*Config

	cfg, err := findConfigInCurrentDir(dir)
	if err != nil {
		return nil, err
	}

	if cfg != nil {
		cfgs = append(cfgs, cfg)
	}

	// find in sub directories.
	subDirs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range subDirs {
		if !file.IsDir() {
			continue
		}

		if subConfigFiles, err := findAllConfigs(path.Join(dir, file.Name())); err != nil {
			return nil, err
		} else {
			cfgs = append(cfgs, subConfigFiles...)
		}
	}

	return cfgs, nil
}

func findConfigInCurrentDir(dir string) (*Config, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Errorf("find files in %s failed for %s", dir, err)
		return nil, err
	}

	cfg := &Config{}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		switch file.Name() {
		case TearDownFileName:
			cfg.TeardownFile = path.Join(dir, file.Name())
		case SetupFileName:
			cfg.SetupFile = path.Join(dir, file.Name())
		case AssertFileName:
			cfg.AssertFile = path.Join(dir, file.Name())
		default:
			cfg.SQLFiles = append(cfg.SQLFiles, path.Join(dir, file.Name()))
		}
	}

	if len(cfg.SQLFiles) >= 1 && cfg.AssertFile != "" {
		return cfg, nil
	} else {
		return nil, nil
	}
}
