from utils.common import EbroTestCase


class TestPlan(EbroTestCase):

    def test_default_plan_is_correct(self):
        commands = ["-plan", "-p"]
        for command in commands:
            with self.subTest(command):
                exit_code, stdout = self.ebro(command)
                self.assertEqual(exit_code, 0)
                self.assertStdout(
                    stdout,
                    """
                    :apt:pre-config
                    :caddy:package-apt-config
                    :docker:package-apt-config
                    :apt:default
                    :caddy:package
                    :docker:package
                    :caddy:default
                    :docker:default
                    :default
                    """,
                )

    def test_plan_for_different_task_is_correct(self):
        exit_code, stdout = self.ebro("-plan", "farm:chicken")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            """
            :farm:egg
            :farm:chicken
            """,
        )

    def test_plan_for_freestanging_task_is_correct(self):
        exit_code, stdout = self.ebro("-plan", "farm:egg")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            """
            :farm:egg
            """,
        )

    def test_plan_for_cyclic_task_fails_with_correct_help(self):
        exit_code, stdout = self.ebro("-plan", "--file", "Ebro.cyclic.yaml")
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            """
            ███ ERROR: validating task chicken: checking reference chain: cyclic reference detected:
            :egg -> :chicken -> :egg
            """,
        )

    def test_plan_for_wrong_references(self):
        exit_code, stdout = self.ebro("-plan", "--file", "Ebro.wrong_references.yaml")
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            """
            ███ ERROR: validating task default: checking reference chain: task :nonexistent does not exist
            """,
        )
