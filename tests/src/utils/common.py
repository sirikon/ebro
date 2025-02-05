import unittest
from os import getcwd, environ
from os.path import join, normpath
from subprocess import run, PIPE, STDOUT
from textwrap import dedent


class EbroTestCase(unittest.TestCase):
    def __init__(self, *args, **kwargs) -> None:
        self.maxDiff = None
        self.workdir = normpath(join(getcwd(), "..", "playground"))
        self.bin = normpath(join(getcwd(), "..", "out", "ebro-e2e"))
        super().__init__(*args, **kwargs)

    def setUp(self) -> None:
        run(["rm", "-rf", ".ebro"], cwd=self.workdir, check=True)
        run(["rm", "-rf", ".cache"], cwd=self.workdir, check=True)
        return super().setUp()

    def assertStdout(self, first: str, second: str):
        self.assertEqual(dedent(second).strip() + "\n", first)

    def assertStdoutStrict(self, first: str, second: str):
        self.assertEqual(second, first)

    def ebro(self, *args, cwd=None, env=dict()):
        result = run(
            [self.bin, *args],
            cwd=cwd or self.workdir,
            stdout=PIPE,
            stderr=STDOUT,
            env=environ | dict(GOCOVERDIR=join(getcwd(), ".coverage")) | env,
        )
        return result.returncode, result.stdout.decode("utf-8")
