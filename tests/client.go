package tests

import (
	"database/sql"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type SQLStep struct {
	ClientID int
	SQL      string
}

const (
	SQLBegin    = "begin;"
	SQLCommit   = "commit"
	SQLRollBack = "rollback"

	SQLTypeBegin    = 1
	SQLTypeCommit   = 2
	SQLTypeRollback = 3
	SQLTypeOther    = 4
)

type Client struct {
	txn      *sql.Tx
	db       *sql.DB
	ClientID int
}

func (Client) GetSQLType(s string) int {
	s = strings.ToLower(s)
	switch s {
	case SQLBegin:
		return SQLTypeBegin
	case SQLCommit:
		return SQLTypeCommit
	case SQLRollBack:
		return SQLTypeRollback
	default:
		return SQLTypeOther
	}
}

func (c *Client) Exec(s string) (err error) {
	log.Infof("%d client run sql: %s", c.ClientID, s)

	t := c.GetSQLType(s)
	switch t {
	case SQLTypeBegin:
		if c.txn != nil {
			err := errors.New(fmt.Sprintf("begin when txn not finished."))
			log.Fatal(err)
		} else {
			if c.txn, err = c.db.Begin(); err != nil {
				log.Fatal(err)
			}
		}

		return
	case SQLTypeCommit:
		if c.txn == nil {
			err := errors.New(fmt.Sprintf("commit when txn not begin."))
			log.Fatal(err)
		} else {
			if err := c.txn.Commit(); err != nil {
				log.Fatal("commit transaction failed. ", err)
			}

			c.txn = nil
		}
	case SQLTypeRollback:
		if c.txn == nil {
			err := errors.New(fmt.Sprintf("rollback when txn not begin."))
			log.Fatal(err)
		} else {
			if err := c.txn.Rollback(); err != nil {
				log.Fatal("rollback transaction failed. ", err)
			}
			c.txn = nil
		}
	default:
		if c.txn != nil {
			_, err = c.txn.Exec(s)
		} else {
			_, err = c.db.Exec(s)
		}

		if err != nil {
			log.Fatalf("run sql error. %p, %s, %s", c.txn, s, err)
		}
	}

	return nil
}
