#!/usr/bin/env python3

from flask import Flask, request

app = Flask(__name__)

@app.route('/')
def index():
    """
    Render the index page.
    """
    return dict(request.headers)

if __name__ == "__main__":
    app.run(debug=True)
