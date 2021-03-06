package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
)

// Data ...
type Data struct {
	TID    string
	Handle string
	Date   string
	Res    string
	QBody  string
	ABody  string
	Note   string
}

// Log ...
type Log struct {
	TID    string
	Handle string
	Date   string
	Res    string
	Body   string
}

// SE ...
type SE struct {
	Stas []Log
	Ends []Log
	Note string
}

// Q ...
type Q struct {
	Sta  string
	End  string
	Note string
}

func selectQs(db *sql.DB) []Q {
	// 問題のidsを取得
	lim := 15913
	que := fmt.Sprintf(`
	SELECT Q.start_log_ids, Q.end_log_ids, Q.note
	FROM question AS Q
	WHERE Q.question_id > %d
	`, lim)
	rows, err := db.Query(que)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// 問題をリストにして戻す
	qs := []Q{}
	for rows.Next() {
		tmp := Q{}
		rows.Scan(&tmp.Sta, &tmp.End, &tmp.Note)
		qs = append(qs, tmp)
	}
	return qs
}

func selectLog(db *sql.DB, ids string) []Log {
	logs := []Log{}
	for _, id := range strings.Split(ids, ",") {
		// idでlogを検索する
		que := fmt.Sprintf(`
		SELECT L.thread_id, L.handle, L.datetime, L.responce_num, L.body
		FROM log AS L
		WHERE L.log_id = %s
		`, id)
		row := db.QueryRow(que)

		// logをリストにする
		tmp := Log{}
		row.Scan(&tmp.TID, &tmp.Handle, &tmp.Date, &tmp.Res, &tmp.Body)
		logs = append(logs, tmp)
	}
	return logs
}

func selectHandle(se SE) string {
	ops := se.Stas
	if len(se.Ends) != 0 {
		ops = append(ops, se.Ends...)
	}
	for _, op := range ops {
		if op.Handle == "あなたのうしろに名無しさんが・・・" {
			continue
		}
		if op.Handle == "本当にあった怖い名無し" {
			continue
		}
		if op.Handle == "ウミガメ信者" {
			continue
		}
		return op.Handle
	}
	return ops[0].Handle
}

func formatData(se SE) Data {
	// ハンドル名を選択
	handle := selectHandle(se)

	// ログを加工
	qBody := ""
	for _, sta := range se.Stas {
		qBody += sta.Body + "\n"
	}
	aBody := ""
	for _, end := range se.Ends {
		aBody += end.Body + "\n"
	}

	// レスを加工
	res := se.Stas[0].Res
	if len(se.Ends) != 0 {
		res += "-" + se.Ends[len(se.Ends)-1].Res
	}
	return Data{se.Stas[0].TID, handle, se.Stas[0].Date, res, qBody, aBody, se.Note}
}

func insertData(stmt *sql.Stmt, data Data) {
	if _, err := stmt.Exec(data.TID, data.Handle, data.Date, data.Res, data.QBody, data.ABody, data.Note); err != nil {
		panic(err)
	}
}

func main() {
	// umigamelogのコネクションを開く
	db, err := sql.Open("sqlite3", "../log.db")
	if err != nil {
		panic(err)
	}

	// 問題を取得
	// 解説がないときは空にする
	ses := []SE{}
	qs := selectQs(db)
	for _, q := range qs {
		stas := selectLog(db, q.Sta)
		if q.End != "" {
			ends := selectLog(db, q.End)
			ses = append(ses, SE{stas, ends, q.Note})
		} else {
			ses = append(ses, SE{stas, []Log{}, q.Note})
		}
	}

	// qのデータベースのコネクションを開く
	db2, err := sql.Open("sqlite3", "../q.db")
	if err != nil {
		panic(err)
	}
	tx, err := db2.Begin()
	if err != nil {
		panic(err)
	}
	que := `INSERT INTO q (tID, handle, date, res, qBody, aBody, note) VALUES (?, ?, ?, ?, ?, ?, ?)`
	stmt, err := tx.Prepare(que)
	if err != nil {
		panic(err)
	}
	// 問題を収録
	for _, se := range ses {
		data := formatData(se)
		insertData(stmt, data)
		fmt.Print(data)
	}
	tx.Commit()
}
