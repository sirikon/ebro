from utils.common import EbroTestCase


class TestCli(EbroTestCase):

    def test_help_is_displayed(self):
        exit_code, stdout = self.ebro("-help")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            """
            Usage: ebro [-command?] [--flags?...] [targets?...]

            Available commands:
              -help
              -version
              -config
              -catalog
              -plan

            Available flags:
              --file
            """,
        )

    def test_unknown_commands_are_reported(self):
        exit_code, stdout = self.ebro("-invent")
        self.assertEqual(exit_code, 1)
        self.assertStdout(stdout, "ERROR: unknown command: invent")

    def test_unknown_flags_are_reported(self):
        exit_code, stdout = self.ebro("-plan", "--invent")
        self.assertEqual(exit_code, 1)
        self.assertStdout(stdout, "ERROR: unknown flag: invent")

    def test_version_is_displayed(self):
        exit_code, stdout = self.ebro("-version")
        self.assertEqual(exit_code, 0)
        self.assertStdout(stdout, "dev")
