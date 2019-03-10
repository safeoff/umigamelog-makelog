package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"regexp"
	"strconv"
	"strings"
	"html"
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

func selectThread(db *sql.DB, tID string) string {
	// idでlogを検索する
	que := fmt.Sprintf(`
	SELECT thread
	FROM thread
	WHERE thread_id = "%s"
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
		p = strings.TrimRight(p, "　")

		p = html.EscapeString(p)
		p = strings.Replace(p, "{", "&#123;", -1)
		p = strings.Replace(p, "}", "&#125;", -1)

		if strings.Index(p, "　 ") != -1 {
			p = "<div class=\"aa\">" + p + "</div>"
		} else {
			p = "<div>" + p + "</div>"
		}

		p = strings.Replace(p, "\n", "<br>", -1)

		p = strings.Replace(p, "\"\"", "\"", -1)

		p = strings.Replace(p, "a href=", "", -1)

		p = strings.Replace(p, "<br><br> <br><br>", "<br><br>", -1)

		rows := strings.Split(p, "<br>")
		p = ""
		for rowi, row := range rows {
			if rowi != 0 {
				p += "<br>"
			}
			re, _ := regexp.Compile(`h?ttps?://[\w/:%#\$&\?\(\)~\.=\+\-]+`)
			match := re.FindAllStringSubmatch(row, -1)
			for _, m := range match {
				row = strings.Replace(row, m[0], "<a href=\"h"+strings.TrimLeft(m[0], "h")+"\">h"+strings.TrimLeft(m[0], "h")+"</a>", -1)
			}
			re, _ = regexp.Compile(`&gt;&gt;[0-9]+[\-[0-9]*]?`)
			match = re.FindAllStringSubmatch(row, -1)
			for _, m := range match {
				row = strings.Replace(row, m[0], "<a href=\"#"+strings.TrimLeft(m[0], "&gt;&gt;")+"\">"+m[0]+"</a> ", -1)
			}
			p += row
		}
		p = strings.Replace(p, "&gt;", ">", -1)

		lid, _ := strconv.Atoi(res.LID)
		n := res.Handle
		h += "<div class=\"box"
		if contains(QIDs, lid) {
			h += " qa"
			n = "<a href=\"../../search/?q=" + n + "&op=and\">" + n + "</a>"
		}
		h += "\"><div class=\"footnote\" id=\"" + strconv.Itoa(i+1) + "\">" + strconv.Itoa(i+1) + " " + n + " " + res.Mail + " " + res.Date + " " + res.ID + "</div>" + p + "</div>\n"
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

func getDate(date1 string) string {
	s := strings.Split(date1, " ")
	s = strings.Split(s[0], "/")
	year := s[0]
	// アンゴルモア歴とかを置換
	rep := regexp.MustCompile(`[^0-9]`)
	year = rep.ReplaceAllString(year, "")
	// 2665を2005に置換 エイプリルフールめんどくさ
	if year == "2665" {
		year = "2005"
	}
	// 20**ではない文字列を置換
	if year[:1] != "2" {
		year = "20" + year
	}
	month := s[1]
	day := strings.Split(s[2], "(")[0]
	// 時刻はきめうち
	date2 := year + "-" + month + "-" + day + "T12:00:00+09:00"
	return date2
}

func writeMD(tID string, thread string, date string) {
	t, _ := strconv.Atoi(tID)
	be := strconv.Itoa(t-1)
	af := strconv.Itoa(t+1)

	s := "---\ntitle: " + thread +
		"\ndate: " + date +
		"\ntags: [" + strings.Split(date, "-")[0] +
		",オカルト板]" +
		"\n---" +
		"\n<div class=\"th_left\"><a href=\"../" + be + "\"><< " + be + "</a></div>" +
		"\n<div class=\"th_right\"><a href=\"../" + af + "\">" + af + " >></a></div>" +
		"\n<br><br>" +
		"\n<script src=\"../../js/cupsoup.js\"></script>" +
		"\n<form>" +
		"\n<input type=\"button\" value=\"問題と解説のみにする\" onClick=\"toggleCupsoup()\">" +
		"\n</form>" +
		"\n{{< " + tID + " >}}" +
		"\n<div class=\"th_left\"><a href=\"../" + be + "\"><< " + be + "</a></div>" +
		"\n<div class=\"th_right\"><a href=\"../" + af + "\">" + af + " >></a></div>"
	f, _ := os.Create("../umigamelog-hugo/content/posts/" + tID + ".md")
	defer f.Close()
	f.Write(([]byte)(s))
}

func writeHTML(tID string, body string) {
	f, _ := os.Create("../umigamelog-hugo/layouts/shortcodes/" + tID + ".html")
	//f, _ := os.Create(tID + ".html")
	defer f.Close()
	f.Write(([]byte)(body))
}

func main() {
	// umigamelogのコネクションを開く
	db, err := sql.Open("sqlite3", "../log.db")
	if err != nil {
		panic(err)
	}

	for tID := 829; tID < 833; tID++ {
		s_tID := strconv.Itoa(tID)
		// 全出題のレス番号を配列を取得
		QIDs := selectQIDs(db)
		// スレッド名を取得
		thread := selectThread(db, s_tID)
		// 1スレッド分のログを取得
		log := selectLog(db, s_tID)
		// 本文・名前欄を整形
		body := makeBody(db, log, QIDs)
		// 更新日を取得
		date := getDate(log[len(log)-1].Date)
		// ファイルを出力
		writeMD(s_tID, thread, date)
		writeHTML(s_tID, body)
		fmt.Print(s_tID + ", ")
	}
}
