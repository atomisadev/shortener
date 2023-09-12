package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"math/rand"
)

type URLMapping struct {
    db *sql.DB
}

func New() (*URLMapping, error) {
    db, err := sql.Open("sqlite3", "./urlshortener.db")
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS url_mapping (
			url TEXT NOT NULL,
			alias TEXT NOT NULL PRIMARY KEY) WITHOUT ROWID;
		`)

		if err != nil {
			return nil, err
		}

		return &URLMapping{db}, nil
}

func (m *URLMapping) Shorten(url string) (string, error) {
    alias := genAlias(url)

		_, err := m.db.Exec("INSERT INTO url_mapping(url, alias) values(?, ?)", url, alias)

		if err != nil {
			return "", err
		}

		return alias, nil
}

func (m *URLMapping) Delete(alias string) error {
	_, err := m.db.Exec("DELETE FROM url_mapping WHERE alias = ?", alias)

	return err
}

func genAlias(url string) string {
    var sb strings.Builder
    for i := 0; i < 6; i++ {
        sb.WriteByte('a' + byte(rand.Intn(26)))
    }
    return sb.String()
}

func (m *URLMapping) Resolve(alias string) (string, error) {
		var url string
		err := m.db.QueryRow("SELECT url FROM url_mapping WHERE alias = ?", alias).Scan(&url)

		if err != nil {
			if err == sql.ErrNoRows {
				return "", errors.New("alias not found")
			}
			return "", err
		}

		return url, nil
}