from utils.common import EbroTestCase


class TestInventory(EbroTestCase):

    def test_inventory_with_absolute_workdir_is_correct(self):
        exit_code, stdout = self.ebro("-inventory", "--file", "Ebro.workdirs.yaml")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            :child:
              working_directory: /somewhere/absolute/child
              environment:
                EBRO_BIN: {self.bin}
                EBRO_ROOT: {self.workdir}
                EBRO_ROOT_FILE: {self.workdir}/Ebro.workdirs.yaml
                ABSTRACT_WORKING_DIRECTORY: /somewhere/absolute/abstract
                EBRO_TASK_ID: :child
                EBRO_TASK_MODULE: ":"
                EBRO_TASK_NAME: child
                EBRO_TASK_WORKING_DIRECTORY: /somewhere/absolute/child
              script: echo $ABSTRACT_WORKING_DIRECTORY
            :default:
              working_directory: /somewhere/absolute
              environment:
                EBRO_BIN: {self.bin}
                EBRO_ROOT: {self.workdir}
                EBRO_ROOT_FILE: {self.workdir}/Ebro.workdirs.yaml
                EBRO_TASK_ID: :default
                EBRO_TASK_MODULE: ":"
                EBRO_TASK_NAME: default
                EBRO_TASK_WORKING_DIRECTORY: /somewhere/absolute
              script: echo "Hello!"
            :other-absolute:
              working_directory: /other/absolute
              environment:
                EBRO_BIN: {self.bin}
                EBRO_ROOT: {self.workdir}
                EBRO_ROOT_FILE: {self.workdir}/Ebro.workdirs.yaml
                EBRO_TASK_ID: :other-absolute
                EBRO_TASK_MODULE: ":"
                EBRO_TASK_NAME: other-absolute
                EBRO_TASK_WORKING_DIRECTORY: /other/absolute
              script: echo "Hello from the other absolute side!"
            :other-relative:
              working_directory: /somewhere/absolute/other/relative
              environment:
                EBRO_BIN: {self.bin}
                EBRO_ROOT: {self.workdir}
                EBRO_ROOT_FILE: {self.workdir}/Ebro.workdirs.yaml
                EBRO_TASK_ID: :other-relative
                EBRO_TASK_MODULE: ":"
                EBRO_TASK_NAME: other-relative
                EBRO_TASK_WORKING_DIRECTORY: /somewhere/absolute/other/relative
              script: echo "Hello from the other relative side!"
            :submodule:other:
              working_directory: /somewhere/absolute/submodule
              environment:
                EBRO_BIN: {self.bin}
                EBRO_ROOT: {self.workdir}
                EBRO_ROOT_FILE: {self.workdir}/Ebro.workdirs.yaml
                EBRO_TASK_ID: :submodule:other
                EBRO_TASK_MODULE: :submodule
                EBRO_TASK_NAME: other
                EBRO_TASK_WORKING_DIRECTORY: /somewhere/absolute/submodule
              script: echo "Hello from the other side!"
            :submodule:other-absolute:
              working_directory: /other/absolute
              environment:
                EBRO_BIN: {self.bin}
                EBRO_ROOT: {self.workdir}
                EBRO_ROOT_FILE: {self.workdir}/Ebro.workdirs.yaml
                EBRO_TASK_ID: :submodule:other-absolute
                EBRO_TASK_MODULE: :submodule
                EBRO_TASK_NAME: other-absolute
                EBRO_TASK_WORKING_DIRECTORY: /other/absolute
              script: echo "Hello from the other absolute side!"
            :submodule:other-relative:
              working_directory: /somewhere/absolute/submodule/other/relative
              environment:
                EBRO_BIN: {self.bin}
                EBRO_ROOT: {self.workdir}
                EBRO_ROOT_FILE: {self.workdir}/Ebro.workdirs.yaml
                EBRO_TASK_ID: :submodule:other-relative
                EBRO_TASK_MODULE: :submodule
                EBRO_TASK_NAME: other-relative
                EBRO_TASK_WORKING_DIRECTORY: /somewhere/absolute/submodule/other/relative
              script: echo "Hello from the other relative side!"
            """,
        )

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
