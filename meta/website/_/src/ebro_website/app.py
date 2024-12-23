from pathlib import Path
from os import getcwd, listdir
from os.path import join

from flask import Flask, render_template, send_file, request
from markdown import Markdown


app = Flask(__name__)
md = Markdown(
    output_format="html5", extensions=["extra", "meta", "codehilite"], tab_length=2
)


@app.get("/")
def index():
    with open("docs/README.md", "r") as f:
        html = md.convert(f.read())
    return render_template("index.html", content=html, active_menu="home")


@app.get("/install/")
def install():
    with open("docs/install.md", "r") as f:
        html = md.convert(f.read())
    return render_template("install.html", content=html, active_menu="install")


@app.get("/changelog/")
def version():
    versions = []
    for item in listdir("docs/changelog"):
        tag = item.removesuffix(".md")
        with open(join("docs", "changelog", item), "r") as f:
            html = md.convert(f.read())
        versions.append({"tag": tag, "notes": html})
    return render_template("changelog.html", versions=versions, active_menu="changelog")


@app.get("/schema.json")
def schema():
    return send_file(join(getcwd(), "docs", "schema.json"))
