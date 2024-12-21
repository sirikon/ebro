from pathlib import Path

from flask import Flask, render_template, Response
from markdown import Markdown


app = Flask(__name__)


@app.get("/")
def index():
    md = Markdown(
        output_format="html5", extensions=["extra", "meta", "codehilite"], tab_length=2
    )
    html = md.convert(Path("docs/README.md").read_text())
    return render_template("index.html", content=html)


@app.get("/schema.json")
def schema():
    return Response(
        Path("docs/schema.json").read_bytes(),
        200,
        {"Content-Type": "application/json; charset=utf8"},
    )
