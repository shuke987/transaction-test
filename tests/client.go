package tests

import "database/sql"

type SQLInfo struct {
	ClientID int
	SQL      string
}

type Client struct {
	db       *sql.DB
	ClientID int
}
