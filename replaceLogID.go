package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"regexp"
	"strconv"
	//"strings"
)

// Data ...
type Question struct {
	TID  string
	Res  string
	Note string
}

// Q ...
type STAEND struct {
	Sta string
	End string
}

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

// log.dbから、start_log_idsとend_log_idsを取得する
func getOldSTAENDs(db *sql.DB) []STAEND {
	que := fmt.Sprintf(`SELECT start_log_ids, end_log_ids FROM question`)
	rows, err := db.Query(que)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	oldSTAENDs := []STAEND{}
	for rows.Next() {
		t := STAEND{}
		rows.Scan(&t.Sta, &t.End)
		oldSTAENDs = append(oldSTAENDs, t)
	}

	return oldSTAENDs
}

// レスのlog_id取得
func getLID(db *sql.DB, tid string, res string, column string) string {
	// logからlog_idを取得
	que := fmt.Sprintf(`
	SELECT log_id FROM log
	WHERE thread_id="%s" AND responce_num=%s
	`, tid, res)
	row := db.QueryRow(que)
	LID := ""
	row.Scan(&LID)

	return LID
}

func calcDiff(idstring string) []int {
	// log_idsの差分数値の配列を作成
	ns := []int{}
	rep := regexp.MustCompile(`\s*,\s*`)
	ids := rep.Split(idstring, -1)
	for i, _ := range ids {
		if i == 0 {
			ns = append(ns, 0)
			continue
		}
		origin, _ := strconv.Atoi(ids[0])
		n, _ := strconv.Atoi(ids[i])
		ns = append(ns, n-origin)
	}

	return ns
}

// 問題と解説のlog_idを取得する
func getSTAENDs(db *sql.DB, q Question, o STAEND) STAEND {
	// 問題レス番と解説レス番を取得
	rep := regexp.MustCompile(`\s*-\s*`)
	res := rep.Split(q.Res, -1)

	// 問題レスのlog_id取得
	sta := getLID(db, q.TID, res[0], "start_log_ids")

	// log_idsが複数ある場合は、差分でlog_idsを生成
	diffsta := calcDiff(o.Sta)
	stas := ""
	originsta, _ := strconv.Atoi(sta)
	for i, diff := range diffsta {
		t := strconv.Itoa(originsta + diff)
		stas += t
		if i != len(diffsta)-1 {
			stas += ","
		}
	}
	fmt.Print(diffsta)

	// mikaiketsuじゃなければ、解説レスのlog_id取得
	ends := ""
	if q.Note != "mikaiketsu" && q.Note != "mikaiketsu " {
		end := getLID(db, q.TID, res[1], "end_log_ids")

		// log_idsが複数ある場合は、差分でlog_idsを生成
		diffend := calcDiff(o.End)
		fmt.Print(diffend)
		originend, _ := strconv.Atoi(end)
		for i := len(diffend) - 1; i >= 0; i-- {
			t := strconv.Itoa(originend - diffend[i])
			ends += t
			if i != 0 {
				ends += ","
			}
		}
	}

	staend := STAEND{stas, ends}
	return staend
}

// questionにstart_log_idsとend_log_idsとnoteを入れる
func updateLIDs(db *sql.DB, old STAEND, lids STAEND, note string) {

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

	// log.dbから、start_log_idsとend_log_idsを取得する
	oldSTAENDs := getOldSTAENDs(logdb)

	// 問題の配列でループ
	for i, _ := range questions {
		// 問題と解説のlog_idを取得する
		LIDs := getSTAENDs(logdb, questions[i], oldSTAENDs[i])
		// questionにstart_log_idsとend_log_idsとnoteを入れる
		updateLIDs(logdb, oldSTAENDs[i], LIDs, questions[i].Note)
		fmt.Print(questions[i])
		fmt.Println(LIDs)
	}
}
