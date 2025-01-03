import unittest
import glob
import jsonschema
import json
import yaml


class TestSchema(unittest.TestCase):

    def test_schema_validates_all_ebro_yamls_in_playground(self):
        with open("../docs/schema.json") as f:
            schema = json.loads(f.read())

        files = flatten(
            [
                glob.glob(p, recursive=True)
                for p in ["../playground/**/Ebro.*.yaml", "../playground/**/Ebro.yaml"]
            ],
        )

        for file in files:
            if file.endswith("/Ebro.unknown_properties.yaml"):
                continue
            with open(file) as f:
                content = yaml.load(f, yaml.Loader)
            with self.subTest(file):
                jsonschema.validate(content, schema)


def flatten(l):
    r = []
    for i in l:
        r += i
    return r
