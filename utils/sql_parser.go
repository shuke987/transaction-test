package utils

import (
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

func LoadFile(filePath string) ([]string, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error("open sql file failed, %s", err)
		return nil, err
	}

	var sqls []string
	tmp := strings.Split(string(file), ";")
	for _, o := range tmp {
		sql := strings.Trim(o, " \r\n\t")
		if sql != "" {
			sqls = append(sqls, sql+";")
		}
	}

	return sqls, nil
}
