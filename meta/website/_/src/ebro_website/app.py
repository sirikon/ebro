import json
from os import getcwd
from os.path import join
from subprocess import run, PIPE

from flask import Flask, render_template, send_file, request
from markdown import Markdown
from markdown.extensions import Extension
from markdown.treeprocessors import Treeprocessor


class MyTreeprocessor(Treeprocessor):
    def run(self, root):
        parent_map = dict((c, p) for p in root.iter() for c in p)

        # Fix links pointing to .md files. Point to .html instead.
        links = root.findall(".//a")
        for link in links:
            if link.attrib["href"].endswith(".md"):
                link.attrib["href"] = link.attrib["href"].removesuffix(".md") + ".html"

        # Remove h1 titles, so we can have those in the documents themselves as
        # titles when reading the source Markdown content, but are removed in
        # the website version.
        els = root.findall(".//h1")
        for el in els:
            parent_map[el].remove(el)


class MyExtension(Extension):
    def extendMarkdown(self, md):
        md.treeprocessors.register(MyTreeprocessor(), "mytreeprocessor", 1)


app = Flask(__name__)
md = Markdown(
    output_format="html5",
    extensions=["extra", "meta", "codehilite", "md_in_html", MyExtension()],
    tab_length=2,
)


MENU = [
    {"id": "home", "name": "Home", "url": "/"},
    {"id": "ebro-format", "name": "<code>Ebro.yaml</code>", "url": "/ebro-format.html"},
    {"id": "install", "name": "Install", "url": "/install.html"},
    {"id": "changelog", "name": "Changelog", "url": "/changelog.html"},
    {"name": "Source Code", "url": "https://github.com/sirikon/ebro"},
]


@app.context_processor
def global_variables():
    return {
        "MENU": MENU,
    }


@app.template_filter("markdown")
def markdown(text: str) -> str:
    return md.convert(text)


@app.get("/")
def index():
    with open("docs/README.md", "r") as f:
        html = md.convert(f.read())
    return render_template("index.html", content=html, active_menu="home")


@app.get("/ebro-format.html")
def ebro_format():
    with open("docs/ebro-format.md", "r") as f:
        html = md.convert(f.read())
    with open("docs/schema.json", "r") as f:
        schema = json.loads(f.read())
    return render_template(
        "ebro-format.html", content=html, schema=schema, active_menu="ebro-format"
    )


@app.get("/install.html")
def install():
    with open("docs/install.md", "r") as f:
        html = md.convert(f.read())
    return render_template("install.html", content=html, active_menu="install")


@app.get("/changelog.html")
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
