import re
from utils.common import EbroTestCase


class TestCli(EbroTestCase):

    def test_version_is_displayed(self):
        commands = ["-version", "-v"]
        for command in commands:
            with self.subTest(command):
                exit_code, stdout = self.ebro(command)
                self.assertEqual(exit_code, 0)
                version = searchline("^version: (.+)$", stdout).group(1)
                commit = searchline("^commit: (.+)$", stdout).group(1)
                date = searchline("^date: (.+)$", stdout).group(1)
                self.assertStdout(
                    stdout,
                    f"""
                    version: {version}
                    commit: {commit}
                    date: {date}
                    """,
                )


def searchline(pattern, text):
    return re.search(pattern, text, flags=re.MULTILINE)
