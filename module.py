#!/usr/bin/python
# -*- coding:utf-8 -*-
u"""テーブルを作る・見る."""


def select_sql(sql):
    u"""引数はSELECT文。一致したものを返す."""
    import sqlite3
    conn = sqlite3.connect("umigamelog.sqlite")
    cur = conn.cursor()
    cur.execute(sql)
    return cur.fetchone()[0]


def search_between(t, s, e):
    u"""文字列tにおいてsから始まりeで終わる箇所を返す."""
    import re
    match = re.search(s + r'.*?' + e, t)
    if match is None:
        return None
    t = match.group()
    t = t.lstrip(s)
    t = t.rstrip(e)
    return t


def read_file_unicode(t):
    u"""ファイルを読み込んでunicodeに変換して返す."""
    a = open(t).read()
    lookup = (
        'utf-8', 'euc_jp', 'euc_jis_2004', 'euc_jisx0213',
        'shift_jis', 'shift_jis_2004', 'shift_jisx0213',
        'iso2022jp', 'iso2022_jp_1', 'iso2022_jp_2', 'iso2022_jp_3',
        'iso2022_jp_ext', 'cp932'
    )
    for enc in lookup:
        try:
            a = a.decode(enc)
            break
        except:
            pass
    return a


def html2unicode(u):
    u"""渡されたurlからhtmlをunicodeにして返す."""
    import urllib2

    h = urllib2.urlopen(u).read()
    # strをunicodeにする
    return h.decode('utf-8')


def remove_a(t):
    u"""aタグを除去、テキストは残す."""
    import re

    t = re.sub(r'<a .*?>', u'', t)
    t = t.replace(u'</a>', u'')
    return t


def shaping_body(t):
    u"""htmlの本文をsql用に整形."""
    import re

    t = t.replace(u'<br> ', u'\n')
    t = t.replace(u'</dd>', u'')
    t = re.sub(r'=\"mg.*?\">', u'', t)
    t = re.sub(r'<p style=\"margin:.*?>', u'\n\n', t)
    t = t.replace(u'</p> ', u'\n')
    t = re.sub(r'<a .*?>', u'', t)
    t = t.replace(u'</a>', u'')
    t = t.replace(u'"', u'""')
    t = t.strip()
    t = t.rstrip('<br><br>')
    t = t.replace(u'<br>', u'\n')
    return t


def htmlentity2unicode(text):
    u"""正規表現のコンパイル."""
    import re
    import htmlentitydefs

    reference_regex = re.compile(u'&(#x?[0-9a-f]+|[a-z]+);', re.IGNORECASE)
    num16_regex = re.compile(u'#x\d+', re.IGNORECASE)
    num10_regex = re.compile(u'#\d+', re.IGNORECASE)
    result = u''
    i = 0
    while True:
        # 実体参照 or 文字参照を見つける
        match = reference_regex.search(text, i)
        if match is None:
            result += text[i:]
            break
        result += text[i:match.start()]
        i = match.end()
        name = match.group(1)
        # 実体参照
        if name in htmlentitydefs.name2codepoint.keys():
            result += unichr(htmlentitydefs.name2codepoint[name])
        # 文字参照
        elif num16_regex.match(name):
            # 16進数
            result += unichr(int(u'0' + name[1:], 16))
        elif num10_regex.match(name):
            # 10進数
            result += unichr(int(name[1:]))
    return result


if __name__ == '__main__':
    function()
