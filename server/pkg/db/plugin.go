package db

import (
	"Havoc/pkg/logger"
	"database/sql"
	"errors"
)

type Plugin struct {
	Name    string
	Type    string
	Version string
	Path    string
}

func (db *Database) PluginExists(name string) (bool, error) {
	var (
		stmt *sql.Stmt
		err  error
	)

	// prepare statement for the sql query
	stmt, err = db.sqlite.Prepare("SELECT COUNT(*) FROM Plugins WHERE name = ?")
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

		if NumRows >= 1 {
			return true, err
		} else {
			return false, err
		}
	}

	return false, nil
}

func (db *Database) Plugin(name string) (*Plugin, error) {
	var (
		stmt   *sql.Stmt
		err    error
		exist  bool
		plugin Plugin
	)

	if exist, err = db.PluginExists(name); err != nil {
		return nil, err
	} else if !exist {
		return nil, errors.New("plugin does not exist")
	}

	stmt, err = db.sqlite.Prepare("SELECT name, type, version, path FROM Plugins WHERE name = ?")
	if err != nil {
		logger.DebugError("sqlite.Prepare failed: %v", err)
		return nil, err
	}

	err = stmt.QueryRow(name).Scan(&plugin.Name, &plugin.Type, &plugin.Version, &plugin.Path)
	if err != nil {
		logger.DebugError("stmt.Query failed: %v", err)
		return nil, err
	}

	return &plugin, err
}

func (db *Database) PluginList() ([]Plugin, error) {
	var (
		plugins []Plugin
		query   *sql.Rows
		err     error
	)

	// prepare the query statement
	query, err = db.sqlite.Query("SELECT name, type, version, path FROM Plugins")
	if err != nil {
		logger.DebugError("sqlite.Query failed: %v", err)
		return nil, err
	}

	defer query.Close()

	// iterate over the list of plugins and
	// scan them into the plugin entry structure
	for query.Next() {
		var entry Plugin

		err = query.Scan(&entry.Name, &entry.Type, &entry.Version, &entry.Path)
		if err != nil {
			logger.DebugError("query.Scan failed: %v", err)
			return nil, err
		}

		plugins = append(plugins, entry)
	}

	return plugins, nil
}

func (db *Database) PluginRegister(name, _type, version, path string) error {
	var (
		stmt  *sql.Stmt
		err   error
		exist bool
	)

	// check if the to be inserted listener name already exists
	exist, err = db.PluginExists(name)
	if err != nil {
		return err
	} else if exist {
		return errors.New("plugin already exist")
	}

	// create sql insert statement
	if stmt, err = db.sqlite.Prepare(`
        INSERT INTO Plugins (name, type, version, path) values(?, ?, ?, ?)
    `); err != nil {
		logger.DebugError("sqlite.Prepare failed: %v", err)
		return err
	}

	// insert the data into the created sql statement
	if _, err = stmt.Exec(name, _type, version, path); err != nil {
		logger.DebugError("stmt.Exec failed: %v", err)
		return err
	}

	return nil
}

func (db *Database) PluginUnRegister(name string) error {
	var (
		stmt  *sql.Stmt
		err   error
		exist bool
	)

	// check if the to be inserted listener name already exists
	exist, err = db.PluginExists(name)
	if err != nil {
		return err
	} else if !exist {
		return errors.New("plugin does not exist")
	}

	// create sql insert statement
	if stmt, err = db.sqlite.Prepare(`
        DELETE FROM Plugins WHERE name = ?
    `); err != nil {
		logger.DebugError("sqlite.Prepare failed: %v", err)
		return err
	}

	// insert the data into the created sql statement
	if _, err = stmt.Exec(name); err != nil {
		logger.DebugError("stmt.Exec failed: %v", err)
		return err
	}

	return nil
}
