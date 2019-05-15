package tests

import (
	log "github.com/sirupsen/logrus"
)

var (
	GlobalCfg GlobalConfig
)

type TestCase struct {
	SQLs [][]string
	DSN  string
}

func (testCase *TestCase) Load(cfg Config) error {
	return nil
}

func (testCase *TestCase) Run() error {

	return nil
}

func Run(testCase []*TestCase) error {
	if err := GlobalCfg.Setup(); err != nil {
		return err
	}

	for _, t := range testCase {
		if err := t.Run(); err != nil {
			return err
		}
	}

	if err := GlobalCfg.Teardown(); err != nil {
		return err
	}

	return nil
}

func LoadAllTests(dir string) ([]*TestCase, error) {
	if err := GlobalCfg.Load(dir); err != nil {
		return nil, err
	}

	cfgs, err := findAllConfigs(dir)
	if err != nil {
		return nil, err
	}

	log.Infof("%d configs found", len(cfgs))

	cases := make([]*TestCase, 0, len(cfgs))

	for _, o := range cfgs {
		testCase := &TestCase{}
		if err := testCase.Load(*o); err != nil {
			return nil, err
		}

		testCase.DSN = GlobalCfg.DSN
		cases = append(cases, testCase)
	}

	return cases, nil
}
