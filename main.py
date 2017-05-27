#!/usr/bin/env python
import webapp2

downloads_html = """<html>
<body>
download the latest build at:
<table>
<tr>
<th>os</th>
<th>download link</th>
</tr>
<tr>
<td>linux amd64</td>
<td><a href='https://go-lanscan.appspot.com/go-lanscan'>https://go-lanscan.appspot.com/go-lanscan</a></td>
</tr>
</table>
Source code at <a href='https://bitbucket.org/johnpfeiffer/go-lanscan'>https://bitbucket.org/johnpfeiffer/go-lanscan</a>
</body>
</html>
"""

class MainHandler(webapp2.RequestHandler):
    def get(self):
        self.response.write(downloads_html)

app = webapp2.WSGIApplication([('/', MainHandler)], debug=True)

