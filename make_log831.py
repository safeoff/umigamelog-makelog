# -*- coding:utf-8 -*-
"""."""
import sqlite3
import re

conn = sqlite3.connect("umigamelog.sqlite")
cur = conn.cursor()

for x in xrange(831, 832):
    sql = '''
    SELECT thread
    FROM thread
    WHERE thread_id = "%(thr)s"
    ''' % {'thr': unicode(x)}
    cur.execute(sql)
    tit = cur.fetchone()[0]
    sql = '''
    SELECT handle, mail, datetime, id, body, log_id
    FROM log
    WHERE thread_id = "%(thr)s"
    ''' % {'thr': unicode(x)}
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
    que = u''
    for i in xrange(0, cnt):
        sql = '''
        SELECT start_log_ids, end_log_ids
        FROM question
        WHERE start_log_ids like %(q)s OR end_log_ids like %(q)s
        ''' % {'q': log[i]}
        cur.execute(sql)
        fe = cur.fetchone()
        if fe is not None:
            for fetch in fe:
                que += fetch + u','
    que = que.split(',')
    que = list(set(que))
    que.remove(u'')
    que = map(int, que)
    que = map(lambda n: n - int(log[0]) + 1, que)
    que.sort()
    while que[0] < 0:
        que.pop(0)
    htm = u''
    for i in xrange(len(bod)):
        bod[i] = bod[i].rstrip(u'　 ')
        bod[i] = bod[i].replace(u'　 \n', u'\n').rstrip(u'　 ')
        if u'　 ' in bod[i]:
            bod[i] = u'<p class="aa">' + bod[i] + u' </p>'
        else:
            bod[i] = u'<p>' + bod[i] + u' </p>'
        bod[i] = bod[i].replace(u'\n', u'<br>')
        bod[i] = bod[i].replace(u'\n \n', u'<br><br>').replace(u'\n', u' ')
        bod[i] = bod[i].replace(u'""', u'"').replace(u'a href=', u'')
        while u'<br><br> <br><br>' in bod[i]:
            bod[i] = bod[i].replace(u'<br><br> <br><br>', u'<br><br>')

        tmp = bod[i].split('<br>')
        for j in xrange(len(tmp)):
            mat = re.findall(u'h?ttps?://[\w/:%#\$&\?\(\)~\.=\+\-]+', tmp[j])
            for m in mat:
                tmp[j] = tmp[j].replace(m, u'<a href="h' + m.lstrip(u'h') + u'">h' + m.lstrip(u'h') + u'</a>')
        bod[i] = u''
        for j in xrange(len(tmp)):
            bod[i] += tmp[j] + u'\n'
        bod[i] = bod[i].rstrip(u'\n')
        bod[i] = bod[i].replace(u'\n', u'<br>')
        tmp = bod[i].split('<br>')
        for j in xrange(len(tmp)):
            mat = re.findall(u'>>[0-9]+[\-[0-9]*]?', tmp[j])
            for m in mat:
                tmp[j] = tmp[j].replace(m, u'<a href="#' + m.lstrip(u'>>') + '">' + m + u'</a>')
        bod[i] = u''
        for j in xrange(len(tmp)):
            bod[i] += tmp[j] + u'\n'
        bod[i] = bod[i].rstrip(u'\n')
        bod[i] = bod[i].replace(u'\n', u'<br>')
        if i + 1 == que[0]:
            htm += u'<div class="box">'
            han[i] = u'''
            <a href="../search/?q=%(han)s&op=and">%(han)s</a>
            ''' % {'han': han[i]}
        htm += u'''
            <h1 id="%(res)s">
                %(res)s %(han)s %(mai)s %(dat)s %(uid)s
            </h1>
            %(bod)s
        ''' % {'res': unicode(i + 1),
               'han': han[i],
               'mai': mai[i],
               'dat': dat[i],
               'uid': uid[i],
               'bod': bod[i]}
        if i + 1 == que[0]:
            htm += u'</div>'
            if len(que) > 1:
                que.pop(0)
    htm += u'''
        <a href="../">過去問情報ナビゲータ（過去ログ保管庫） </a>
        </div>
    </div>
    '''
    f = open('log/' + unicode(x) + '.html', 'w')
    f.write(htm.encode('utf-8').replace('\n', '').replace('    ', ''))
    f.close()
    # f.write(htm.encode('utf-8'))
    print (x)
cur.close()
conn.close()
