package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

// Data ...
type Data struct {
	TID    string `json:"tID"`
	Handle string `json:"handle"`
	Date   string `json:"date"`
	Res    string `json:"res"`
	QBody  string `json:"qBody"`
	ABody  string `json:"aBody"`
}

func sqdb() {
	// データベースのコネクションを開く
	db, err := sql.Open("sqlite3", os.Args[1])
	if err != nil {
		panic(err)
	}

	// 複数レコード取得
	s := "SELECT tID, handle, date, res, qBody, aBody FROM q"

	rows, err := db.Query(s)
	if err != nil {
		panic(err)
	}

	// 処理が終わったらカーソルを閉じる
	defer rows.Close()

	ds := []Data{}
	for rows.Next() {
		d := Data{}

		// カーソルから値を取得
		if err := rows.Scan(&d.TID, &d.Handle, &d.Date, &d.Res, &d.QBody, &d.ABody); err != nil {
			break
		}
		ds = append(ds, d)
	}
	for _, d := range ds {
		fmt.Println(d)
	}
}

func main() {
	sqdb()
}
