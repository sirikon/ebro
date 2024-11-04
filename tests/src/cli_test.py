from utils.common import EbroTestCase


class TestCli(EbroTestCase):

    def test_help_is_displayed(self):
        exit_code, stdout = self.ebro("-help")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            """
            Usage: ebro [flag?] [targets?...]
            
            Available flags:
              -config   Display all imported configuration files merged into one
              -catalog  Display complete catalog of tasks with their definitive configuration
              -plan     Display the execution plan
            """,
        )

    def test_unknown_flags_are_reported(self):
        exit_code, stdout = self.ebro("-invent")
        self.assertEqual(exit_code, 1)
        self.assertStdout(stdout, "ERROR: unknown flag: invent")
