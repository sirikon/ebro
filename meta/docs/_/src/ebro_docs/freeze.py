from os import getcwd
from os.path import join

from flask_frozen import Freezer
from ebro_docs.app import app

app.config["FREEZER_DESTINATION"] = join(getcwd(), "out", "docs")
freezer = Freezer(app)

if __name__ == "__main__":
    freezer.freeze()
