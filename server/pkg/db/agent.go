package db

import (
	"Havoc/pkg/logger"
	"database/sql"
	"errors"
)

type Agent struct {
	Uuid     string
	Type     string
	Parent   string
	Status   string
	Note     string
	Metadata []byte
}

func (db *Database) AgentInsert(uuid, _type, parent, status, note string, metadata []byte) error {
	var (
		stmt  *sql.Stmt
		err   error
		exist bool
	)

	// check if the to be inserted agent uuid already exists
	exist, err = db.AgentExist(uuid)
	if err != nil {
		return err
	} else if exist {
		return errors.New("agent already exist")
	}

	// create sql insert statement
	if stmt, err = db.sqlite.Prepare(`
        INSERT INTO Agents (uuid, type, parent, status, note, metadata, disabled, hide) values(?, ?, ?, ?, ?, ?, ?, ?)
    `); err != nil {
		logger.DebugError("sqlite.Prepare failed: %v", err)
		return err
	}

	// insert the data into the created sql statement
	if _, err = stmt.Exec(
		uuid,
		_type,
		parent,
		status,
		note,
		metadata,
		false,
		false,
	); err != nil {
		logger.DebugError("stmt.Exec failed: %v", err)
		return err
	}

	return err
}

func (db *Database) AgentUpdate(uuid, parent, status, note string, metadata []byte) error {
	var (
		stmt *sql.Stmt
		err  error
	)

	// create sql insert statement
	if stmt, err = db.sqlite.Prepare(`
        UPDATE Agents SET parent = ?, status = ?, note = ?, metadata = ? WHERE uuid = ?
    `); err != nil {
		logger.DebugError("sqlite.Prepare failed: %v", err)
		return err
	}

	// insert the data into the created sql statement
	if _, err = stmt.Exec(parent, status, note, metadata, uuid); err != nil {
		logger.DebugError("stmt.Exec failed: %v", err)
		return err
	}

	return err
}

func (db *Database) AgentExist(uuid string) (bool, error) {
	var (
		stmt *sql.Stmt
		err  error
	)

	// prepare statement for the sql query
	stmt, err = db.sqlite.Prepare("SELECT COUNT(*) FROM Agents WHERE uuid = ?")
	if err != nil {
		logger.DebugError("sqlite.Prepare failed: %v", err)
		return false, err
	}

	// execute statement
	query, err := stmt.Query(uuid)
	defer query.Close()
	if err != nil {
		logger.DebugError("stmt.Query failed: %v", err)
		return false, err
	}

	for query.Next() {
		var NumRows int

		err = query.Scan(&NumRows)
		if err != nil {
			logger.DebugError("query.Scan failed: %v", err)
			return false, err
		}

		if NumRows > 1 {
			return true, err
		} else {
			return false, err
		}
	}

	return false, nil
}

func (db *Database) AgentDisable(uuid string) error {
	var (
		stmt *sql.Stmt
		err  error
	)

	// create sql insert statement
	if stmt, err = db.sqlite.Prepare(`
        UPDATE Agents SET disabled = ? WHERE uuid = ?
    `); err != nil {
		logger.DebugError("sqlite.Prepare failed: %v", err)
		return err
	}

	// insert the data into the created sql statement
	if _, err = stmt.Exec(true, uuid); err != nil {
		logger.DebugError("stmt.Exec failed: %v", err)
		return err
	}

	return err
}

func (db *Database) AgentSetHide(uuid string, hide bool) error {
	var (
		stmt *sql.Stmt
		err  error
	)

	// create sql insert statement
	if stmt, err = db.sqlite.Prepare(`
        UPDATE Agents SET hide = ? WHERE uuid = ?
    `); err != nil {
		logger.DebugError("sqlite.Prepare failed: %v", err)
		return err
	}

	// insert the data into the created sql statement
	if _, err = stmt.Exec(hide, uuid); err != nil {
		logger.DebugError("stmt.Exec failed: %v", err)
		return err
	}

	return err
}

func (db *Database) AgentLists() ([]Agent, error) {
	var (
		agents []Agent
		query  *sql.Rows
		err    error
	)

	// prepare the query statement
	query, err = db.sqlite.Query("SELECT uuid, type, parent, status, note, metadata FROM Agents WHERE disabled = 0")
	if err != nil {
		logger.DebugError("sqlite.Query failed: %v", err)
		return nil, err
	}

	defer query.Close()

	// iterate over the list of agents and
	// scan them into the agent entry structure
	for query.Next() {
		var entry Agent

		err = query.Scan(&entry.Uuid, &entry.Type, &entry.Parent, &entry.Status, &entry.Note, &entry.Metadata)
		if err != nil {
			logger.DebugError("query.Scan failed: %v", err)
			return nil, err
		}

		agents = append(agents, entry)
	}

	return agents, nil
}

func (db *Database) AgentConsole(uuid string) ([][]byte, error) {
	var (
		entry []byte
		list  [][]byte
		query *sql.Rows
		err   error
	)

	query, err = db.sqlite.Query("SELECT data FROM AgentConsole WHERE uuid = ?", uuid)
	if err != nil {
		logger.DebugError("sqlite.Query failed: %v", err)
		return nil, err
	}

	defer query.Close()

	for query.Next() {
		if err = query.Scan(&entry); err != nil {
			logger.DebugError("query.Scan failed: %v", err)
			return nil, err
		}

		list = append(list, entry)
	}

	return list, nil
}

func (db *Database) AgentConsoleInsert(uuid string, data []byte) error {
	var (
		stmt *sql.Stmt
		err  error
	)

	// create sql insert statement
	if stmt, err = db.sqlite.Prepare(`
        INSERT INTO AgentConsole (uuid, data) values(?, ?)
    `); err != nil {
		logger.DebugError("sqlite.Prepare failed: %v", err)
		return err
	}

	// insert the data into the created sql statement
	if _, err = stmt.Exec(uuid, data); err != nil {
		logger.DebugError("stmt.Exec failed: %v", err)
		return err
	}

	return err
}
