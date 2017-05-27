#!/usr/bin/env python
import webapp2

class MainHandler(webapp2.RequestHandler):
    def get(self):
        self.response.write('download the latest build at https://go-lanscan.appspot.com/go-lanscan')

app = webapp2.WSGIApplication([('/', MainHandler)], debug=True)

