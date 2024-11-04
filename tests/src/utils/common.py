from os import getcwd
from os.path import join, normpath
from subprocess import run, PIPE
from textwrap import dedent
import unittest


class EbroTestCase(unittest.TestCase):
    def __init__(self, *args, **kwargs) -> None:
        self.maxDiff = None
        self.workdir = normpath(join(getcwd(), "..", "playground"))
        super().__init__(*args, **kwargs)

    def setUp(self) -> None:
        run(["rm", "-rf", ".ebro"], cwd=self.workdir, check=True)
        run(["rm", "-rf", ".cache"], cwd=self.workdir, check=True)
        return super().setUp()

    def assertStdout(self, first, second):
        self.assertEqual(first, dedent(second).strip() + "\n")

    def ebro(self, *args):
        result = run(
            [join(getcwd(), "..", "out", "ebro"), *args],
            cwd=self.workdir,
            stdout=PIPE,
            stderr=PIPE,
        )
        return result.returncode, result.stdout.decode("utf-8")
