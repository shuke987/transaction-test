package tests

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
)

type GlobalConfig struct {
	DSNSetup string `json:"dsn-setup"`
	Database string `json:"database"`
	DSN      string `json:"dsn"`
}

func (cfg *GlobalConfig) Load(dir string) error {
	filePath := path.Join(dir, "global.json")

	jsonData, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("load file error, ", filePath)
		return err
	}

	err = json.Unmarshal(jsonData, cfg)
	if err != nil {
		log.Error("parse global config error, ", err)
		return err
	}

	if cfg.DSN == "" || cfg.Database == "" || cfg.DSNSetup == "" {
		errStr := fmt.Sprintf("bad global config file")
		log.Error(errStr)
		return errors.New(errStr)
	}

	return nil
}

func (cfg *GlobalConfig) Setup() error {
	db, err := sql.Open("mysql", cfg.DSNSetup)
	if err != nil {
		log.Error("connect db error ", err)
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	sqlCmd := "DROP DATABASE if EXISTS " + cfg.Database
	if _, err := db.Exec(sqlCmd); err != nil {
		log.Error("drop db error, %s, %s", sqlCmd, err)
		return err
	}

	sqlCmd = "CREATE DATABASE " + cfg.Database
	if _, err := db.Exec(sqlCmd); err != nil {
		log.Error("create db error, %s, %s", sqlCmd, err)
		return err
	}

	return nil
}

func (cfg *GlobalConfig) Teardown() error {
	db, err := sql.Open("mysql", cfg.DSNSetup)
	if err != nil {
		log.Error("connect db error ", err)
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	sqlCmd := "DROP DATABASE if EXISTS " + cfg.Database
	if _, err := db.Exec(sqlCmd); err != nil {
		log.Error("drop db error, %s, %s", sqlCmd, err)
		return err
	}

	return nil
}
