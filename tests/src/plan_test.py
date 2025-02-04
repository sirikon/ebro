from utils.common import EbroTestCase


class TestPlan(EbroTestCase):

    def test_plan_for_cyclic_task_fails_with_correct_help(self):
        exit_code, stdout = self.ebro("-plan", "--file", "Ebro.cyclic.yaml")
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            """
            ███ ERROR: planning could not complete. there could be a cyclic dependency. here is the list of tasks remaining to be planned and their requirements:
            :chicken:
            - :egg
            :default:
            - :chicken
            :egg:
            - :chicken
            """,
        )

    def test_plan_for_wrong_references_1(self):
        exit_code, stdout = self.ebro("-plan", "--file", "Ebro.wrong_references_1.yaml")
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            """
            ███ ERROR: normalizing 'requires' for task ':default': referenced task ':nonexistent' does not exist
            """,
        )
    
    def test_plan_for_wrong_references_2(self):
        exit_code, stdout = self.ebro("-plan", "--file", "Ebro.wrong_references_2.yaml")
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            """
            ███ ERROR: normalizing 'required_by' for task ':default': referenced task ':nonexistent' does not exist
            """,
        )
    
    def test_plan_for_wrong_references_3(self):
        exit_code, stdout = self.ebro("-plan", "--file", "Ebro.wrong_references_3.yaml")
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            """
            ███ ERROR: normalizing 'extends' for task ':default': referenced task ':nonexistent' does not exist
            """,
        )
