package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"strconv"
	"strings"
	"regexp"
	"os"
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
	LID    string
	Handle string
	Mail   string
	Date   string
	ID     string
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

func selectQIDs(db *sql.DB) []int {
	// すべてのlog_idsをquestionテーブルから取得して配列に入れる
	que := fmt.Sprintf(` SELECT start_log_ids, end_log_ids FROM question `)
	rows, err := db.Query(que)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	QIDs := []int{}
	for rows.Next() {
		t1, t2 := "", ""
		err := rows.Scan(&t1, &t2)
		if err != nil {
			if err == sql.ErrNoRows {
				continue
			}
		}
		t := strings.Split(t1+","+t2, ",")
		a := make([]int, len(t))
		for i := range a {
			a[i], _ = strconv.Atoi(t[i])
		}
		QIDs = append(QIDs, a...)
	}

	return QIDs
}

func makeBody(db *sql.DB, log []Log, QIDs []int) string {
	// 全角スペースを除去
	// AAなら本文をaaに
	// 改行を<br>に
	// ""を"に
	// aタグを除去
	// 改行を抑制
	// httpをリンクに
	// レスアンカーをリンクに
	// log_idがQIDsに存在するならdivをbox＆searchのリンク
	// 名前欄を作成
	h := ""
	for i, res := range log {
		p := res.Body
		p = strings.Replace(p, "　", "", -1)

		if (strings.Index(p, "　 ") != -1) {
			p = "<p class=\"aa\"" + p + "</p>"
		} else {
			p = "<p>" + p + "</p>"
		}

		p = strings.Replace(p, "\n", "<br>", -1)

		p = strings.Replace(p, "\"\"", "\"", -1)

		p = strings.Replace(p, "a href=", "", -1)

		p = strings.Replace(p, "<br><br> <br><br>", "<br><br>", -1)

		rows := strings.Split(p, "<br>")
		p = ""
		for _, row := range rows {
			re, _ := regexp.Compile(`h?ttps?://[\w/:%#\$&\?\(\)~\.=\+\-]+`)
			match := re.FindAllStringSubmatch(row, -1)
			for _, m := range match {
				row = strings.Replace(row, m[0], "<a href=\"h" + strings.TrimLeft(m[0], "h") + "\">h" + strings.TrimLeft(m[0], "h") + "</a>", -1)
			}
			re, _ = regexp.Compile(`>>[0-9]+[\-[0-9]*]?`)
			match = re.FindAllStringSubmatch(row, -1)
			for _, m := range match {
				row = strings.Replace(row, m[0], "<a href=\"#" + strings.TrimLeft(m[0], ">>") + "\">" + m[0] + "</a>", -1)
			}
			p += row
		}

		lid, _ := strconv.Atoi(res.LID)
		n := res.Handle
		if contains(QIDs, lid) {
			p = "<div class=\"box\">" + p
			n = "<a href=\"../../search/?=" + n + "&op=and\">" + n + "</a>"
		}
		h += "<h1 id=\""+ strconv.Itoa(i+1) + "\">" + strconv.Itoa(i+1) + " " + n + " " + res.Mail + " " + res.Date + " " + res.ID + "</h1>" + p + "</div>"
	}
	return h
}

func contains(s []int, e int) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}
	// LID    string
	// Handle string
	// Mail   string
	// Date   string
	// ID     string
	// Body   string

func writeHTML(tID string, thread string, body string) {
	f, _ := os.Create(tID + ".html")
	defer f.Close()
	f.Write(([]byte)(body))
	fmt.Print(thread)
}

func main() {
	// umigamelogのコネクションを開く
	db, err := sql.Open("sqlite3", "../log.db")
	if err != nil {
		panic(err)
	}

	// 全出題のレス番号を配列を取得
	QIDs := selectQIDs(db)
	// スレッド名を取得
	tID := 831
	thread := selectThread(db, strconv.Itoa(tID))
	// 1スレッド分のログを取得
	log := selectLog(db, strconv.Itoa(tID))
	// 本文・名前欄を整形
	body := makeBody(db, log, QIDs)
	// ファイルを出力
	writeHTML(strconv.Itoa(tID), thread, body)
}
