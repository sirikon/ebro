from datetime import datetime, timezone
import json
import re
from os import getcwd, listdir
from os.path import join
from subprocess import run, PIPE

from flask import Flask, render_template, send_file
from markdown import Markdown
from markdown.extensions import Extension
from markdown.treeprocessors import Treeprocessor
from markdown.postprocessors import Postprocessor

NOW = int(datetime.now(timezone.utc).timestamp())


class RemoveElementsProcessor(Treeprocessor):
    def run(self, root):
        parent_map = dict((c, p) for p in root.iter() for c in p)
        # Remove elements that should be explicitly removed from the website
        els = root.findall(".//div[@remove-in-website]")
        for el in els:
            parent_map[el].remove(el)


class FixLinksProcessor(Treeprocessor):
    def run(self, root):
        # Fix links pointing to .md files. Point to .html instead.
        links = root.findall(".//a")
        for link in links:
            if link.attrib["href"].endswith(".md"):
                link.attrib["href"] = link.attrib["href"].removesuffix(".md") + ".html"


class IncludeEbroYamlFormatProcessor(Postprocessor):
    REPATTERN = r"<p>\[Ebro\.yaml format explained\]</p>"

    def run(self, text):
        if re.findall(self.REPATTERN, text):
            with open("docs/schema.json", "r") as f:
                schema = json.loads(f.read())
            return re.sub(
                self.REPATTERN,
                render_template("_ebro-format.html", schema=schema),
                text,
            )
        return text


class MyExtension(Extension):
    def extendMarkdown(self, md):
        md.treeprocessors.register(
            RemoveElementsProcessor(), "remove-elements", 9999999
        )
        md.treeprocessors.register(FixLinksProcessor(), "fix-links", -1)
        md.postprocessors.register(
            IncludeEbroYamlFormatProcessor(), "include-ebroyaml", 9999999
        )


app = Flask(__name__)
md = Markdown(
    output_format="html5",
    extensions=["extra", "meta", "codehilite", "md_in_html", "toc", MyExtension()],
    extension_configs={
        "toc": {
            "anchorlink": True,
            "anchorlink_class": "x-anchorlink",
            "title": "Table of Contents",
            "toc_class": "x-toc",
            "title_class": "x-toc-title",
        }
    },
    tab_length=2,
)
md_changelog = Markdown(
    output_format="html5",
    extensions=["extra", "meta", "codehilite", "md_in_html", MyExtension()],
    tab_length=2,
)


MENU = [
    {"id": "home", "name": "Home", "url": "/"},
    {"id": "install", "name": "Install", "url": "/install.html"},
    {"id": "changelog", "name": "Changelog", "url": "/changelog.html"},
    {"name": "Source Code", "url": "https://github.com/sirikon/ebro"},
]


@app.context_processor
def global_variables():
    return {
        "NOW": NOW,
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
            html = md_changelog.convert(f.read())
        versions.append({"tag": tag, "notes": html})
    return render_template("changelog.html", versions=versions, active_menu="changelog")


@app.get("/schema.json")
def schema():
    return send_file(join(getcwd(), "docs", "schema.json"))


def get_tagged_versions():
    result = []
    for file in listdir(join("docs", "changelog")):
        if file.endswith(".md") and file != "HEAD.md":
            result.append(file.removesuffix(".md"))
    result.sort(reverse=True, key=lambda x: tuple([int(n) for n in x.split(".")]))
    return result
