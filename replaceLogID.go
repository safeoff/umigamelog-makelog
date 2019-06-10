package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"regexp"
	//"strings"
)

// Data ...
type Question struct {
	TID    string
	Res    string
	Note   string
}

//// Q ...
//type Q struct {
//	Sta  string
//	End  string
//	Note string
//}

// q.dbから、thread_idとresとnoteを取得する
func getQuestions(db *sql.DB) []Question {
	que := fmt.Sprintf(`SELECT tID, res, note FROM q`)
	rows, err := db.Query(que)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	questions := []Question{}
	for rows.Next() {
		t := Question{}
		rows.Scan(&t.TID, &t.Res, &t.Note)
		questions = append(questions, t)
	}

	return questions
}

// レスのlog_ids取得
func getLIDs(db *sql.DB, tid string, res string, column string) string {
	// logからlog_idを取得
	que := fmt.Sprintf(`
	SELECT log_id FROM log
	WHERE thread_id="%s" AND responce_num=%s
	`, tid, res)
	row := db.QueryRow(que)
	LID := ""
	row.Scan(&LID)

	// log_idsが複数ある場合は、差分でlog_idsを生成


	return LID
}

// 問題と解説のlog_idを取得する
func getSTAENDs(db *sql.DB, q Question) string {
	// 問題レス番と解説レス番を取得
	rep := regexp.MustCompile(`\s*-\s*`)
	res := rep.Split(q.Res, -1)

	// 問題レスのlog_id取得
	sta := getLID(db, q.TID, res[0], "start_log_ids")

	// mikaiketsuじゃなければ、解説レスのlog_id取得
	end := ""
	if q.Note != "mikaiketsu" && q.Note != "mikaiketsu " {
		end = getLID(db, q.TID, res[1], "end_log_ids")
		fmt.Print(end)
	}

	// start_log_idsが複数ある場合は、差分でstart_log_idsを生成
	// end_log_idsが複数ある場合は、差分でend_log_idsを生成

	return sta
}

// log.dbのquestionのidを振り直す
func main() {
	// コネクションを開く
	qdb, err := sql.Open("sqlite3", "../q.db")
	if err != nil {
		panic(err)
	}
	logdb, err := sql.Open("sqlite3", "../log.db")
	if err != nil {
		panic(err)
	}

	// q.dbから、thread_idとresとnoteを取得する
	questions := getQuestions(qdb)

	// 問題の配列でループ
	for _, question := range questions {
		// 問題と解説のlog_idを取得する
		LIDs := getSTAENDs(logdb, question)
		// questionにstart_log_idsとend_log_idsとnoteを入れる
		fmt.Print(question)
		fmt.Print(LIDs)
		fmt.Print("\n")
	}
}
