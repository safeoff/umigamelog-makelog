package main

import (
	"database/sql"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

// Req ...
type Req struct {
	Q  string `json:"q"`
	Op string `json:"op"`
}

// Res ...
type Res struct {
	List []Data `json:"list"`
}

// Data ...
type Data struct {
	TID    string `json:"tID"`
	Handle string `json:"handle"`
	Date   string `json:"date"`
	Res    string `json:"res"`
	QBody  string `json:"qBody"`
	ABody  string `json:"aBody"`
	Note   string `json:"note"`
}

func selectSQL(req Req) (Res, error) {
	if req.Q == "" {
		return Res{}, nil
	}
	// データベースのコネクションを開く
	db, err := sql.Open("sqlite3", "q.db")
	if err != nil {
		panic(err)
	}

	// 複数レコード取得
	s := "SELECT tID, handle, date, res, qBody, aBody, note FROM q WHERE "
	qs := strings.Fields(req.Q)
	s += "(tID || handle || qBody || aBody || note) like '%" + qs[0] + "%' "
	for _, q := range qs {
		s += req.Op + "(tID || handle || qBody || aBody || note) like '%" + q + "%' "
	}
	s += " ORDER BY date DESC LIMIT 1000"

	rows, err := db.Query(s)
	if err != nil {
		panic(err)
	}

	// 処理が終わったらカーソルを閉じる
	defer rows.Close()
	if rows == nil {
		return Res{}, nil
	}

	r := Res{}
	for rows.Next() {
		var tID string
		var handle string
		var date string
		var res string
		var qBody string
		var aBody string
		var note string

		// カーソルから値を取得
		if err := rows.Scan(&tID, &handle, &date, &res, &qBody, &aBody, &note); err != nil {
			break
		}

		qBody = strings.Replace(qBody, "\n", "", -1)
		aBody = strings.Replace(aBody, "\n", "", -1)
		r.List = append(r.List, Data{TID: tID, Handle: handle, Date: date, Res: res, QBody: qBody, ABody: aBody, Note: note})
	}
	return r, nil
}

func main() {
	lambda.Start(selectSQL)
}
