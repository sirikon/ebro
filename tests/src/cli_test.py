from utils.common import EbroTestCase


class TestCli(EbroTestCase):

    def test_help_is_displayed(self):
        exit_code, stdout = self.ebro("-help")
        self.assertEqual(exit_code, 0)
        self.assertStdoutStrict(
            stdout,
            """
  ebro [--flags...] [targets...]
    # Run everything
    flags:
      --file value  Specify the file that should be loaded as root module. default: Ebro.yaml
      --force       Ignore when.* conditionals and dont skip any task. default: false
    targets:
      defaults to [default]


  ebro -config [--flags...]
    # Display all imported configuration files merged into one
    flags:
      --file value  Specify the file that should be loaded as root module. default: Ebro.yaml


  ebro -catalog [--flags...]
    # Display complete catalog of tasks with their definitive configuration
    flags:
      --file value  Specify the file that should be loaded as root module. default: Ebro.yaml


  ebro -plan [--flags...] [targets...]
    # Display the execution plan
    flags:
      --file value  Specify the file that should be loaded as root module. default: Ebro.yaml
    targets:
      defaults to [default]


  ebro -version
    # Display ebro's version


  ebro -help
    # Display this help message

""",
        )

    def test_unknown_commands_are_reported(self):
        exit_code, stdout = self.ebro("-invent")
        self.assertEqual(exit_code, 1)
        self.assertStdout(stdout, "███ ERROR: unknown command: invent")

    def test_unknown_flags_are_reported(self):
        exit_code, stdout = self.ebro("-plan", "--invent")
        self.assertEqual(exit_code, 1)
        self.assertStdout(stdout, "███ ERROR: unknown flag: invent")

    def test_version_is_displayed(self):
        exit_code, stdout = self.ebro("-version")
        self.assertEqual(exit_code, 0)
        self.assertStdout(stdout, "dev")

    def test_file_flag_handles_missing_arg_correctly(self):
        exit_code, stdout = self.ebro("--file")
        self.assertEqual(exit_code, 1)
        self.assertStdout(stdout, "███ ERROR: expected value after --file flag")
