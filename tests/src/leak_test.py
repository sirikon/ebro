import unittest

from os import getcwd
from os.path import join, normpath
from subprocess import run, PIPE

CHECK_HISTORY = False


class TestLeak(unittest.TestCase):

    def test_no_leaks_are_present_anywhere_in_the_project(self):
        sensitive_grep_patterns = [
            normpath(join(getcwd(), "..")),
            "/" + "Users",
            "/" + "home",
        ]
        revisions = [None] + (
            (
                run(
                    ["git", "rev-list", "--all"],
                    cwd=join(getcwd(), ".."),
                    check=True,
                    stdin=None,
                    stdout=PIPE,
                )
                .stdout.decode("utf-8")
                .splitlines()
            )
            if CHECK_HISTORY
            else []
        )

        for pattern in sensitive_grep_patterns:
            with self.subTest(pattern):
                for revision in revisions:
                    cmd = run(
                        [
                            "git",
                            "grep",
                            pattern,
                            *([revision] if revision is not None else []),
                        ],
                        cwd=join(getcwd(), ".."),
                        stdin=None,
                        stdout=PIPE,
                    )
                    self.assertEqual(cmd.stdout.decode("utf-8"), "")
                    self.assertEqual(cmd.returncode, 1)
