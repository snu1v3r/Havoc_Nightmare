package db

import (
	"Havoc/pkg/logger"
	"database/sql"
	"errors"
)

type Listener struct {
	Name     string
	Protocol string
	Status   string
	Config   []byte
}

func (db *Database) ListenerInsert(name, protocol, status string, config []byte) error {
	var (
		stmt  *sql.Stmt
		err   error
		exist bool
	)

	// check if the to be inserted listener name already exists
	exist, err = db.ListenerExists(name)
	if err != nil {
		return err
	} else if exist {
		return errors.New("listener already exist")
	}

	// create sql insert statement
	if stmt, err = db.sqlite.Prepare(`
        INSERT INTO Listeners (name, protocol, status, config) values(?, ?, ?, ?)
    `); err != nil {
		logger.DebugError("sqlite.Prepare failed: %v", err)
		return err
	}

	// insert the data into the created sql statement
	if _, err = stmt.Exec(name, protocol, status, config); err != nil {
		logger.DebugError("stmt.Exec failed: %v", err)
		return err
	}

	return err
}

func (db *Database) ListenerUpdate(name, status string, config []byte) error {
	var (
		stmt *sql.Stmt
		err  error
	)

	// create sql insert statement
	if stmt, err = db.sqlite.Prepare(`
        UPDATE Listeners SET status = ?, config = ? WHERE name = ?
    `); err != nil {
		logger.DebugError("sqlite.Prepare failed: %v", err)
		return err
	}

	// insert the data into the created sql statement
	if _, err = stmt.Exec(status, config, name); err != nil {
		logger.DebugError("stmt.Exec failed: %v", err)
		return err
	}

	return err
}

func (db *Database) ListenerSetConfig(name string, config []byte) error {
	var (
		stmt *sql.Stmt
		err  error
	)

	// create sql insert statement
	if stmt, err = db.sqlite.Prepare(`
        UPDATE Listeners SET config = ? WHERE name = ?
    `); err != nil {
		logger.DebugError("sqlite.Prepare failed: %v", err)
		return err
	}

	// insert the data into the created sql statement
	if _, err = stmt.Exec(config, name); err != nil {
		logger.DebugError("stmt.Exec failed: %v", err)
		return err
	}

	return err
}

func (db *Database) ListenerSetStatus(name string, status string) error {
	var (
		stmt *sql.Stmt
		err  error
	)

	// create sql insert statement
	if stmt, err = db.sqlite.Prepare(`
        UPDATE Listeners SET status = ? WHERE name = ?
    `); err != nil {
		logger.DebugError("sqlite.Prepare failed: %v", err)
		return err
	}

	// insert the data into the created sql statement
	if _, err = stmt.Exec(status, name); err != nil {
		logger.DebugError("stmt.Exec failed: %v", err)
		return err
	}

	return err
}

func (db *Database) ListenerExists(name string) (bool, error) {
	var (
		stmt *sql.Stmt
		err  error
	)

	// prepare statement for the sql query
	stmt, err = db.sqlite.Prepare("SELECT COUNT(*) FROM Listeners WHERE name = ?")
	if err != nil {
		logger.DebugError("sqlite.Prepare failed: %v", err)
		return false, err
	}

	// execute statement
	query, err := stmt.Query(name)
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

func (db *Database) ListenerRemove(name string) error {
	var (
		stmt *sql.Stmt
		err  error
	)

	// prepare some statement to delete
	// the listener entry from the database
	stmt, err = db.sqlite.Prepare("DELETE FROM Listeners WHERE name = ?")
	defer stmt.Close()
	if err != nil {
		logger.DebugError("sqlite.Prepare failed: %v", err)
		return err
	}

	// execute statement with the listener name
	_, err = stmt.Exec(name)

	if err != nil {
		logger.DebugError("stmt.Exec failed: %v", err)
		return err
	}

	return nil
}

func (db *Database) ListenerList() ([]Listener, error) {
	var (
		listeners []Listener
		query     *sql.Rows
		err       error
	)

	// prepare the query statement
	query, err = db.sqlite.Query("SELECT name, protocol, status, config FROM Listeners")
	if err != nil {
		logger.DebugError("sqlite.Query failed: %v", err)
		return nil, err
	}

	defer query.Close()

	// iterate over the list of agents and
	// scan them into the agent entry structure
	for query.Next() {
		var entry Listener

		err = query.Scan(&entry.Name, &entry.Protocol, &entry.Status, &entry.Config)
		if err != nil {
			logger.DebugError("query.Scan failed: %v", err)
			return nil, err
		}

		listeners = append(listeners, entry)
	}

	return listeners, nil
}
