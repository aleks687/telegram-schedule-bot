package database

import (
    "database/sql"
    _ "modernc.org/sqlite"
)

type Database struct {
    Conn *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
    conn, err := sql.Open("sqlite", dbPath)
    if err != nil {
        return nil, err
    }
    
    // Создаем таблицы
    createTables(conn)
    
    return &Database{Conn: conn}, nil
}

func createTables(db *sql.DB) {
    sql := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        telegram_id INTEGER UNIQUE,
        group_name TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
    
    db.Exec(sql)
}

func (db *Database) Close() error {
    return db.Conn.Close()
}