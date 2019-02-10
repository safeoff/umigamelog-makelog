# -*- coding:utf-8 -*-
"""."""
import sqlite3
import re

conn = sqlite3.connect("umigamelog.sqlite")
cur = conn.cursor()

for x in xrange(832, 833):
    hai = u'' + unicode(x) + u''
    sql = '''
        SELECT T.thread
        FROM thread AS T
        WHERE T.thread_id = "%(thr)s"
    ''' % {'thr': hai}
    cur.execute(sql)
    tit = cur.fetchone()[0]
    sql = '''
        SELECT L.handle, L.mail, L.datetime, L.id, L.body, L.log_id
        FROM log AS L
        WHERE L.thread_id = "%(thr)s"
    ''' % {'thr': hai}
    cur.execute(sql)
    cnt = 0
    han = []
    mai = []
    dat = []
    uid = []
    bod = []
    log = []
    for fetch in cur.fetchall():
        cnt += 1
        han.append(unicode(fetch[0]))
        mai.append(unicode(fetch[1]))
        dat.append(unicode(fetch[2]))
        uid.append(unicode(fetch[3]))
        bod.append(unicode(fetch[4]))
        log.append(unicode(fetch[5]))
    que = []
    for i in xrange(0, cnt):
        sql = '''
            SELECT Q.question_id, Q.start_log_ids, Q.end_log_ids
            FROM question AS Q
            WHERE Q.start_log_ids like %(q)s
        ''' % {'q': log[i]}
        cur.execute(sql)
        fe = cur.fetchone()
        if fe is not None:
            que.append(fe)
    que = list(set(que))
    que.sort()
    htm = u'''
        <html>
        <head>
        <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
        <meta http-equiv="Content-Style-Type" content="text/css">
        <meta name="viewport" content="width=device-width">
        <link rel="stylesheet" media="all" type="text/css" href="../css/search.css" />
        <link rel="stylesheet" media="all" type="text/css" href="../css/log.css" />
        <title>%(tit)s - カップスープ</title>
        <script type="text/javascript" src="../js/analyticstracking.js"></script>
        </head>
        <body>
        <div class="wrapper">
            <div class="header">
                <a href="../">過去問情報ナビゲータ（過去ログ保管庫） </a>
                <p class="thr">%(tit)s - カップスープ</p>
            </div>
            <div class="content">
    ''' % {'tit': tit}
    for q in que:
        q_body = u''
        a_body = u''
        handle = u''
        datetime = u''
        s_res = 0
        e_res = 0
        for s in q[1].split(u','):
            sql = '''
                SELECT L.body, L.handle, L.datetime, L.responce_num
                FROM log AS L
                WHERE L.log_id = %(a)s
            ''' % {'a': s}
            cur.execute(sql)
            fet = cur.fetchone()
            q_body += unicode(fet[0])
            if handle == u'':
                handle = fet[1]
            elif fet[1] != u'あなたのうしろに名無しさんが・・・' or fet[1] != u'本当にあった怖い名無し' or fet[1] != u'ウミガメ信者':
                handle = fet[1]
            if datetime == u'':
                datetime = fet[2]
            if s_res == 0:
                s_res = fet[3]
        try:
            for s in q[2].split(u','):
                sql = '''
                    SELECT L.body, L.handle, L.datetime, L.responce_num
                    FROM log AS L
                    WHERE L.log_id = %(a)s
                ''' % {'a': s}
                cur.execute(sql)
            fet = cur.fetchone()
            a_body += unicode(fet[0])
            if handle == u'':
                handle = fet[1]
            elif fet[1] != u'あなたのうしろに名無しさんが・・・' or fet[1] != u'本当にあった怖い名無し' or fet[1] != u'ウミガメ信者':
                handle = fet[1]
            if e_res < fet[3]:
                e_res = fet[3]
        except:
            a_body = u'未解決です'
        htm += u'''
            <div class="box">
            <h1>
            <a href="../log/%(thr)s.html#%(ress)s">%(res)s</a>
            <form name="s%(i)s"  action="../search.cgi" method="post">
            <input type="hidden" name="key" value="%(han)s"></form>
             <a href="../search.cgi" onclick="document.s%(i)s.submit();return false;">%(han)s</a>
             %(dat)s</h1>
        ''' % {'thr': hai,
               'res': unicode(s_res) + u'-' + unicode(e_res),
               'ress': unicode(s_res),
               'han': handle,
               'dat': datetime,
               'i': q[0]}

        q_body = q_body.rstrip(u'　 ')
        q_body = q_body.replace(u'　 \n', u'\n').rstrip(u'　 ')
        if u'　 ' in q_body:
            q_body = u'<p class="aa">' + q_body + u' </p><br>'
        else:
            q_body = u'<p>' + q_body + u' </p><br>'
        a_body = a_body.rstrip(u'　 ')
        a_body = a_body.replace(u'　 \n', u'\n').rstrip(u'　 ')
        if u'　 ' in a_body:
            a_body = u'<p class="aa ans">' + a_body + u' </p>'
        else:
            a_body = u'<p class="ans">' + a_body + u' </p>'
        body = q_body + a_body
        body = body.replace(u'\n', u'<br>')
        body = body.replace(u'\n \n', u'<br><br>').replace(u'\n', u' ')
        body = body.replace(u'""', u'"').replace(u'a href=', u'')
        while u'<br><br> <br><br>' in body:
            body = body.replace(u'<br><br> <br><br>', u'<br><br>')

        tmp = body.split('<br>')
        for j in xrange(len(tmp)):
            mat = re.findall(u'h?ttps?://[\w/:%#\$&\?\(\)~\.=\+\-]+', tmp[j])
            for m in mat:
                tmp[j] = tmp[j].replace(m, u'<a href="h' + m.lstrip(u'h') + u'">h' + m.lstrip(u'h') + u'</a>')
        body = u''
        for j in xrange(len(tmp)):
            body += tmp[j] + u'\n'
        body = body.rstrip(u'\n')
        body = body.replace(u'\n', u'<br>')
        tmp = body.split('<br>')
        body = u''
        for j in xrange(len(tmp)):
            body += tmp[j] + u'\n'
        body = body.rstrip(u'\n')
        body = body.replace(u'\n', u'<br>')
        htm += body + u'</div>'
    htm += u'''
        <a href="../">過去問情報ナビゲータ（過去ログ保管庫） </a>
        </div>
    </div>
    '''
    f = open('cupsoup/' + hai + 'cupsoup.html', 'w')
    f.write(htm.encode('utf-8').replace('\n', '').replace('    ', ''))
    f.close()
    # f.write(htm.encode('utf-8'))
    print hai
cur.close()
conn.close()
