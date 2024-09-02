package db

import "database/sql"

func (db *Database) AgentInsert(uuid, _type, parent, status, note string, metadata []byte) error {
	var (
		stmt *sql.Stmt
		err  error
	)

	// create sql insert statement
	if stmt, err = db.sqlite.Prepare(`
        INSERT INTO Agents (uuid, type, parent, status, note, metadata, disabled) values(?, ?, ?, ?, ?, ?, ?)
    `); err != nil {
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
	); err != nil {
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
		return err
	}

	// insert the data into the created sql statement
	if _, err = stmt.Exec(
		parent,
		status,
		note,
		metadata,
		uuid,
	); err != nil {
		return err
	}

	return err
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
		return err
	}

	// insert the data into the created sql statement
	if _, err = stmt.Exec(true, uuid); err != nil {
		return err
	}

	return err
}
