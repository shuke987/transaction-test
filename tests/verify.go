package tests

import (
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type AssertInterface interface {
	Assert() error
}

type Verify struct {
	Type        string                 `json:"type"`
	SQL         string                 `json:"sql"`
	ExpectedMap map[string]interface{} `json:"expect,omitempty"`
}

func (v *Verify) Run() {

}

func (v *Verify) AssertAdminCheck(db *sql.DB) error {
	log.Println("verify: %s", v.SQL)

	if _, err := db.Exec(v.SQL); err != nil {
		log.Fatal("verify failed. %s, %s", v.SQL, err)
	}

	return nil
}

func (v *Verify) AssertSQLResult(db *sql.DB) error {
	return nil

}
