import unittest
import re


class TestCli(unittest.TestCase):

    # I do not accept any feedback or opinion about this godforsaken
    # test. Whatever you're going to tell me, I already know.
    def test_dockerfile_and_tool_versions_are_in_sync(self):
        expected_go_version = None
        actual_go_version = None

        with open("../.tool-versions", "r") as f:
            content = f.read()
            expected_go_version = match(r"^go ([0-9\.]+)$", content)

        with open("../Dockerfile", "r") as f:
            content = f.read()
            actual_go_version = match(r"^FROM golang:([0-9\.]+)$", content)

        self.assertIsNotNone(expected_go_version)
        self.assertIsNotNone(actual_go_version)
        self.assertEqual(
            expected_go_version,
            actual_go_version,
            "Dockerfile and .tool-versions have different go versions defined",
        )


def match(restr: str, content: str):
    return re.compile(restr, re.MULTILINE).match(content).group(1)
