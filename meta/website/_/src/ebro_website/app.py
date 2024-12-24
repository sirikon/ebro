import json
from os import getcwd, listdir
from os.path import join
from subprocess import run, PIPE

from flask import Flask, render_template, send_file, request
from markdown import Markdown


app = Flask(__name__)
md = Markdown(
    output_format="html5",
    extensions=["extra", "meta", "codehilite", "md_in_html"],
    tab_length=2,
)


@app.template_filter("markdown")
def markdown(text: str) -> str:
    return md.convert(text)


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
    for tag in get_tagged_versions():
        with open(join("docs", "changelog", tag + ".md"), "r") as f:
            html = md.convert(f.read())
        versions.append({"tag": tag, "notes": html})
    return render_template("changelog.html", versions=versions, active_menu="changelog")


@app.get("/schema.json")
def schema():
    return send_file(join(getcwd(), "docs", "schema.json"))


def get_tagged_versions():
    command_result = run(
        ["git", "tag", "--list"], check=True, stdin=None, stdout=PIPE, stderr=PIPE
    )
    result = command_result.stdout.decode().splitlines()
    result.sort(reverse=True)
    return result
