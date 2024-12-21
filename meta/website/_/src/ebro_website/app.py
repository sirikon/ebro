from pathlib import Path
from os import getcwd
from os.path import join

from flask import Flask, render_template, Response, send_file
from markdown import Markdown


app = Flask(__name__)


@app.get("/")
def index():
    md = Markdown(
        output_format="html5", extensions=["extra", "meta", "codehilite"], tab_length=2
    )
    with open("docs/README.md", "r") as f:
        html = md.convert(f.read())
    return render_template("_base.html", content=html)


@app.get("/versions/")
def version():
    return render_template("_base.html", content="")


@app.get("/schema.json")
def schema():
    return send_file(join(getcwd(), "docs", "schema.json"))
