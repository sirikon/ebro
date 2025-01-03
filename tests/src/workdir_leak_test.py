import unittest

from os import getcwd
from os.path import join, normpath
from subprocess import run, PIPE


class TestWorkdirLeak(unittest.TestCase):

    def test_workdir_is_not_present_anywhere_in_the_project(self):
        cmd = run(
            ["git", "grep", normpath(join(getcwd(), ".."))],
            stdin=None,
            stdout=PIPE,
        )
        self.assertEqual(cmd.stdout.decode("utf-8"), "")
        self.assertEqual(cmd.returncode, 1)
