package main

import (
	"database/sql"
	"time"
)

func CreateTable(db *sql.DB) error {
	SQLCreateThought :=
		`CREATE TABLE thought (
		"tid" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,		
		"title" TEXT NOT NULL,
		"date" DATETIME NOT NULL,
		"html" TEXT	NOT NULL,
		"markdown" TEXT	NOT NULL
	  );`

	SQLCreateConfig :=
		`CREATE TABLE config (
		"key" TEXT NOT NULL PRIMARY KEY,		
		"value" TEXT NOT NULL
	  );`
	_, err := db.Exec(SQLCreateThought)
	if err != nil {
		return err
	}
	_, err = db.Exec(SQLCreateConfig)
	return err
}

func UpdateConfig(db *sql.DB, conf *Config) error {
	m := conf.Map()
	for k, v := range m {
		_, err := db.Exec(`REPLACE INTO config(key, value) VALUES (?, ?);`, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetThoughtsByPage(db *sql.DB, size, page int64) ([]*Thought, error) {
	rows, err := db.Query(`SELECT tid, title, date, html, markdown FROM thought ORDER BY tid DESC LIMIT ? OFFSET ?;`, size, size*(page-1))
	if err != nil {
		if err == sql.ErrNoRows {
			if page == 1 {
				return []*Thought{{
					TID:      0,
					Title:    "Hello World!",
					Date:     time.Now(),
					HTML:     "Create your thoughts here!",
					Markdown: "Create your thoughts here!",
				}}, nil
			}
			return nil, nil
		}
		return nil, err
	}
	var thts []*Thought
	for rows.Next() {
		tht := &Thought{}
		err := rows.Scan(&tht.TID, &tht.Title, &tht.Date, &tht.HTML, &tht.Markdown)
		if err != nil {
			return nil, err
		}
		thts = append(thts, tht)
	}
	return thts, nil

}

func GetThoughtsNumbers(db *sql.DB) (int64, error) {
	var num int64
	err := db.QueryRow(`SELECT count (tid) FROM thought;`).Scan(&num)
	if err != nil {
		return 0, err
	}
	return num, nil

}

func InsertThought(db *sql.DB, tht *Thought) (int64, error) {
	SQLInsertPost := `INSERT INTO thought(title, date, html, markdown) VALUES (?, ?, ?, ?);`
	result, err := db.Exec(SQLInsertPost, tht.Title, tht.Date, tht.HTML, tht.Markdown)
	if err != nil {
		return 0, err
	}
	tid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return tid, nil
}

func GetConfigMap(db *sql.DB) (map[string]string, error) {
	m := make(map[string]string)
	rows, err := db.Query(`SELECT key, value FROM config;`)
	if err != nil {
		if err == sql.ErrNoRows {
			return m, nil
		}
		return m, err
	}
	for rows.Next() {
		var k, v string
		if err := rows.Scan(&k, &v); err != nil {
			return m, err
		}
		m[k] = v
	}
	return m, nil
}
