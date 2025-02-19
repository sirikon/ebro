from utils.common import EbroTestCase


class TestInventory(EbroTestCase):

    def test_inventory_works_with_task_with_nothing_to_do(self):
        exit_code, stdout = self.ebro(
            "-inventory", "--file", "Ebro.fail_when_nothing_to_do.yaml"
        )
        self.assertStdout(
            stdout,
            f"""
            :default:
              working_directory: {self.workdir}
              environment:
                EBRO_BIN: {self.bin}
                EBRO_ROOT: {self.workdir}
                EBRO_ROOT_FILE: {self.workdir}/Ebro.fail_when_nothing_to_do.yaml
                EBRO_TASK_ID: :default
                EBRO_TASK_MODULE: ":"
                EBRO_TASK_NAME: default
                EBRO_TASK_WORKING_DIRECTORY: {self.workdir}
            """,
        )
        self.assertEqual(exit_code, 0)

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
