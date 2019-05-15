package main

import (
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
	"transaction-test/tests"
)

func main() {
	log.Println("test start")

	testCases, err := tests.LoadAllTests("./samples")
	if err != nil {
		log.Fatal(err)
	}

	if err := tests.Run(testCases); err != nil {
		log.Fatal(err)
	}

	log.Info("all case passed")
}
