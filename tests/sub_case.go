package tests

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type SubCase struct {
	DSN           string
	ClientMap     map[int]*Client
	SetupSQLs     []string
	TeardownSQLs  []string
	SQLs          []SQLStep
	Verifications []Verify
}

func (subCase *SubCase) Run(nClients int) (err error) {
	if err := subCase.Setup(); err != nil {
		return err
	}
	defer func() {
		if err2 := subCase.TearDown(err != nil); err2 != nil {
			if err != nil {
				err = err2
			}
		}
	}()

	// init multiple client.
	for i := 0; i < nClients; i++ {
		c := &Client{
			ClientID: i,
		}

		db, err := sql.Open("mysql", subCase.DSN)
		if err != nil {
			log.Error("connect db error ", err)
			return err
		} else {
			c.db = db
		}

		subCase.ClientMap[i] = c
	}

	// run sqls
	for _, o := range subCase.SQLs {
		c := subCase.ClientMap[o.ClientID]
		if err := c.Exec(o.SQL); err != nil {
			return err
		}
	}

	return nil
}

func (subCase *SubCase) Setup() error {
	db, err := sql.Open("mysql", subCase.DSN)
	if err != nil {
		log.Error("connect db error ", err)
		return err
	}
	defer db.Close()

	if len(subCase.SetupSQLs) > 0 {
		for _, o := range subCase.SetupSQLs {
			if _, err := db.Exec(o); err != nil {
				log.Errorf("setup error, %s, %s", o, err)
				return err
			}
		}
	}

	return nil
}

func (subCase *SubCase) TearDown(isError bool) error {
	// clean up clients
	for k, v := range subCase.ClientMap {
		if err := v.db.Close(); err != nil {
			log.Errorf("clean up db connection error, %s", err)
			return err
		}

		delete(subCase.ClientMap, k)
	}

	// clean up when no error occurs
	if !isError && len(subCase.TeardownSQLs) > 0 {
		db, err := sql.Open("mysql", subCase.DSN)
		if err != nil {
			log.Error("connect db error ", err)
			return err
		}
		defer db.Close()

		if len(subCase.SetupSQLs) > 0 {
			for _, o := range subCase.TeardownSQLs {
				if _, err := db.Exec(o); err != nil {
					log.Errorf("tear down error, %s, %s", o, err)
					return err
				}
			}
		}
	}

	return nil
}
