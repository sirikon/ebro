import json
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


@app.get("/ebro-format")
def ebro_format():
    with open("docs/ebro-format.md", "r") as f:
        html = md.convert(f.read())
    with open("docs/schema.json", "r") as f:
        schema = json.loads(f.read())
    return render_template(
        "ebro-format.html", content=html, schema=schema, active_menu="ebro-format"
    )


@app.get("/install")
def install():
    with open("docs/install.md", "r") as f:
        html = md.convert(f.read())
    return render_template("install.html", content=html, active_menu="install")


@app.get("/changelog")
def version():
    versions = []
    items = listdir("docs/changelog")
    items.sort(reverse=True)
    for item in items:
        tag = item.removesuffix(".md")
        with open(join("docs", "changelog", item), "r") as f:
            html = md.convert(f.read())
        versions.append({"tag": tag, "notes": html})
    return render_template("changelog.html", versions=versions, active_menu="changelog")


@app.get("/schema.json")
def schema():
    return send_file(join(getcwd(), "docs", "schema.json"))
