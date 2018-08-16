package tools

import (
	"github.com/rhinoman/couchdb-go"
	"time"
	"encoding/json"
	"github.com/myesui/uuid"
)

type CouchDB struct {
	Conn *couchdb.Connection
	Db   *couchdb.Database
}

//gets a CouchDb connection
func Database(host string, port int) (*CouchDB, error) {

	conn, err := couchdb.NewConnection(host, port, time.Duration(1000*time.Millisecond))
	if err != nil {
		return nil, err
	}

	return &CouchDB{Conn: conn, Db: nil}, nil
}

//selects a database
func (d *CouchDB) SelectDb(database string, username string, password string) error {
	auth := couchdb.BasicAuth{Username: username, Password: password}

	d.Db = d.Conn.SelectDB(database, &auth)
	if err := d.Db.DbExists(); err != nil {
		return err
	}
	return nil
}

// CREATE
func (d *CouchDB) Insert(payload interface{}) (string, error) {
	rev, err := d.Db.Save(payload, GetUuid(), "")
	if err != nil {
		return "", err
	}
	return rev, nil
}

// QUERY reimbursement
// Example:
//type FindResponse struct {
//	Docs []TestDocument `json:"docs"`
//}
//selector := `{"Age": {"$eq":202}}`
//var findResult FindResponse
//err = couchDb.query(selector, &findResult)
//ErrorCheck(err)
//
//for _, doc := range findResult.Docs {
//... process ....
//}

func (d *CouchDB) Query(selector string, where interface{}) error {

	var selectorObj interface{}

	err := json.Unmarshal([]byte(selector), &selectorObj)
	if err != nil {
		return err
	}

	params := couchdb.FindQueryParams{Selector: &selectorObj}

	return d.Db.Find(&where, &params)
}

func GetUuid() string {
	theUuid := uuid.NewV4()
	return uuid.Formatter(theUuid, uuid.FormatHex)
}