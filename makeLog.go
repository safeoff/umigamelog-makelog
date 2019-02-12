package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strings"
	"strconv"
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
	LID string
	Handle string
	Mail string
	Date string
	ID string
	Body string
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

// func selectQs(db *sql.DB) []Q {
// 	// 問題のidsを取得
// 	lim := 100000
// 	que := fmt.Sprintf(`
// 	SELECT Q.start_log_ids, Q.end_log_ids, Q.note
// 	FROM question AS Q
// 	WHERE Q.question_id < %d
// 	`, lim)
// 	rows, err := db.Query(que)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer rows.Close()
// 
// 	// 問題をリストにして戻す
// 	qs := []Q{}
// 	for rows.Next() {
// 		tmp := Q{}
// 		rows.Scan(&tmp.Sta, &tmp.End, &tmp.Note)
// 		qs = append(qs, tmp)
// 	}
// 	return qs
// }
// 
// func selectHandle (se SE) string {
// 	ops := se.Stas
// 	if len(se.Ends) != 0 {
// 		ops = append(ops, se.Ends...)
// 	}
// 	for _, op := range ops {
// 		if op.Handle == "あなたのうしろに名無しさんが・・・" {
// 			continue
// 		}
// 		if op.Handle == "本当にあった怖い名無し" {
// 			continue
// 		}
// 		if op.Handle == "ウミガメ信者" {
// 			continue
// 		}
// 		return op.Handle
// 	}
// 	return ops[0].Handle
// }
// 
// func formatData(se SE) Data{
// 	// ハンドル名を選択
// 	handle := selectHandle(se)
// 	
// 	// ログを加工
// 	qBody := ""
// 	for _, sta := range se.Stas {
// 		qBody += sta.Body + "\n"
// 	}
// 	aBody := ""
// 	for _, end := range se.Ends {
// 		aBody += end.Body + "\n"
// 	}
// 
// 	// レスを加工
// 	res := se.Stas[0].Res
// 	if len(se.Ends) != 0 {
// 		res += "-" + se.Ends[len(se.Ends)-1].Res
// 	}
// 	return Data{se.Stas[0].TID, handle, se.Stas[0].Date, res, qBody, aBody, se.Note}
// }
// 
// func insertData(stmt *sql.Stmt, data Data) {
// 	if _, err := stmt.Exec(data.TID, data.Handle, data.Date, data.Res, data.QBody, data.ABody, data.Note); err != nil {
// 		panic(err)
// 	}
// }

func selectThread(db *sql.DB, tID string) string {
	// idでlogを検索する
	que := fmt.Sprintf(`
	SELECT thread
	FROM thread
	WHERE thread_id = %s
	`, tID)
	row := db.QueryRow(que)

	// logをリストにする
	thread := ""
	row.Scan(&thread)
	return thread
}

func selectLog(db *sql.DB, tID string) []Log {
	// idでlogを検索する
	que := fmt.Sprintf(`
    SELECT log_id, handle, mail, datetime, id, body
    FROM log
    WHERE thread_id = "%s"
	`, tID)
	rows, err := db.Query(que)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// logをリストにする
	logs := []Log{}
	for rows.Next() {
		t := Log{}
		rows.Scan(&t.LID, &t.Handle, &t.Mail, &t.Date, &t.ID, &t.Body)
		logs = append(logs, t)
	}

	return logs
}

func searchLIDs(db *sql.DB, logs []Log) []int {
	LIDs := []int{}
	for _, log := range logs {
		// idでlogを検索する
		que := fmt.Sprintf(`
        SELECT start_log_ids, end_log_ids
        FROM question
        WHERE start_log_ids like %s OR end_log_ids like %s
		`, log.LID, log.LID)
		row := db.QueryRow(que)

		// 空ならcontinue
		// 空じゃないならリストに保持
		t1, t2 := "", ""
		err := row.Scan(&t1, &t2)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
		}
		t := strings.Split(t1+","+t2, ",")
		fmt.Print(t)
		//LIDs = append(LIDs, t)
	}

	return LIDs
}

func main() {
	// umigamelogのコネクションを開く
	db, err := sql.Open("sqlite3", "../log.db")
	if err != nil {
		panic(err)
	}

	// スレッド名を取得
	tID := 831
	thread := selectThread(db, strconv.Itoa(tID))
	fmt.Print(thread)
	// 1スレッド分のログを取得
	log := selectLog(db, strconv.Itoa(tID))
	fmt.Print(len(log))
	// 出題のレス番号を取得
	LIDs := searchLIDs(db, log) 
	fmt.Print(len(LIDs))
	// 本文を整形
	// 名前欄を整形

	// ファイルを出力
}
