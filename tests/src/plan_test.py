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
