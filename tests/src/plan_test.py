from utils.common import EbroTestCase


class TestCli(EbroTestCase):

    def test_default_plan_is_correct(self):
        exit_code, stdout = self.ebro("-plan")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            """
            :apt:pre-config
            :docker:package-apt-config
            :apt:default
            :docker:package
            :docker:default
            :default
            """,
        )

    def test_plan_for_different_task_is_correct(self):
        exit_code, stdout = self.ebro("-plan", "chicken")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            """
            :egg
            :chicken
            """,
        )

    def test_plan_for_freestanging_task_is_correct(self):
        exit_code, stdout = self.ebro("-plan", "egg")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            """
            :egg
            """,
        )

    def test_plan_for_cyclic_task_fails_with_correct_help(self):
        exit_code, stdout = self.ebro("-plan", "--file", "Ebro.cyclic.yaml")
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            """
            ERROR: planning could not complete. there could be a cyclic dependency. here is the list of tasks remaining to be planned and their requirements:
            :chicken:
                - :egg
            :default:
                - :chicken
            :egg:
                - :chicken
            """,
        )
