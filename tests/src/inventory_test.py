from utils.common import EbroTestCase


class TestInventory(EbroTestCase):

    def test_inventory_fails_with_task_with_nothing_to_do(self):
        exit_code, stdout = self.ebro(
            "-inventory", "--file", "Ebro.fail_when_nothing_to_do.yaml"
        )
        self.assertStdout(
            stdout,
            f"""
            ███ ERROR: parsing module: parsing 'tasks': parsing task 'default': task has nothing to do (no requires, script, extends nor abstract)
            """,
        )
        self.assertEqual(exit_code, 1)

    def test_unkown_properties_are_not_allowed(self):
        exit_code, stdout = self.ebro(
            "-inventory", "--file", "Ebro.unknown_properties.yaml"
        )
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            f"""
            ███ ERROR: parsing module: unexpected key 'import'
            """,
        )
