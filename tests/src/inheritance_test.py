# from utils.common import EbroTestCase


# class TestInheritance(EbroTestCase):

#     def test_inventory_is_correct(self):
#         exit_code, stdout = self.ebro("-inventory", "--file", "Ebro.inheritance.yaml")
#         self.assertStdout(
#             stdout,
#             f"""
#             :b:
#               working_directory: {self.workdir}
#               environment:
#                 EBRO_BIN: {self.bin}
#                 EBRO_ROOT: {self.workdir}
#                 EBRO_ROOT_FILE: {self.workdir}/Ebro.inheritance.yaml
#                 FOO: It's FOO
#                 BAR: It's BAR
#                 EBRO_TASK_ID: :b
#                 EBRO_TASK_MODULE: ":"
#                 EBRO_TASK_NAME: b
#                 EBRO_TASK_WORKING_DIRECTORY: {self.workdir}
#               requires:
#               - :d
#               - :c
#               required_by:
#               - :default
#               script: echo $BAR
#               when:
#                 check_fails: exit 2
#                 output_changes: echo test
#             :c:
#               working_directory: {self.workdir}
#               environment:
#                 EBRO_BIN: {self.bin}
#                 EBRO_ROOT: {self.workdir}
#                 EBRO_ROOT_FILE: {self.workdir}/Ebro.inheritance.yaml
#                 EBRO_TASK_ID: :c
#                 EBRO_TASK_MODULE: ":"
#                 EBRO_TASK_NAME: c
#                 EBRO_TASK_WORKING_DIRECTORY: {self.workdir}
#               script: echo "I'm C"
#             :d:
#               working_directory: {self.workdir}
#               environment:
#                 EBRO_BIN: {self.bin}
#                 EBRO_ROOT: {self.workdir}
#                 EBRO_ROOT_FILE: {self.workdir}/Ebro.inheritance.yaml
#                 EBRO_TASK_ID: :d
#                 EBRO_TASK_MODULE: ":"
#                 EBRO_TASK_NAME: d
#                 EBRO_TASK_WORKING_DIRECTORY: {self.workdir}
#               script: echo "I'm D"
#             :default:
#               working_directory: {self.workdir}
#               environment:
#                 EBRO_BIN: {self.bin}
#                 EBRO_ROOT: {self.workdir}
#                 EBRO_ROOT_FILE: {self.workdir}/Ebro.inheritance.yaml
#                 EBRO_TASK_ID: :default
#                 EBRO_TASK_MODULE: ":"
#                 EBRO_TASK_NAME: default
#                 EBRO_TASK_WORKING_DIRECTORY: {self.workdir}
#               script: echo Hello
#             :multi-c:
#               working_directory: {self.workdir}
#               environment:
#                 A: "1"
#                 EBRO_BIN: {self.bin}
#                 EBRO_ROOT: {self.workdir}
#                 EBRO_ROOT_FILE: {self.workdir}/Ebro.inheritance.yaml
#                 B: "22"
#                 C: "3"
#                 EBRO_TASK_ID: :multi-c
#                 EBRO_TASK_MODULE: ":"
#                 EBRO_TASK_NAME: multi-c
#                 EBRO_TASK_WORKING_DIRECTORY: {self.workdir}
#                 D: "44"
#               script: echo multi-b
#               quiet: true
#             :y:
#               labels:
#                 label.A: "1"
#                 label.B: "2"
#               working_directory: {self.workdir}
#               environment:
#                 EBRO_BIN: {self.bin}
#                 EBRO_ROOT: {self.workdir}
#                 EBRO_ROOT_FILE: {self.workdir}/Ebro.inheritance.yaml
#                 EBRO_TASK_ID: :y
#                 EBRO_TASK_MODULE: ":"
#                 EBRO_TASK_NAME: "y"
#                 EBRO_TASK_WORKING_DIRECTORY: {self.workdir}
#                 A: "1"
#                 B: "2"
#               script: echo y
#               quiet: true
#               when:
#                 check_fails: exit 1
#                 output_changes: echo test2
#             :z:
#               labels:
#                 label.A: "1"
#                 label.B: "22"
#                 label.C: "33"
#               working_directory: {self.workdir}
#               environment:
#                 EBRO_BIN: {self.bin}
#                 EBRO_ROOT: {self.workdir}
#                 EBRO_ROOT_FILE: {self.workdir}/Ebro.inheritance.yaml
#                 A: "1"
#                 EBRO_TASK_ID: :z
#                 EBRO_TASK_MODULE: ":"
#                 EBRO_TASK_NAME: z
#                 EBRO_TASK_WORKING_DIRECTORY: {self.workdir}
#                 B: "22"
#                 C: "3"
#               script: echo z
#               quiet: false
#               when:
#                 check_fails: exit 1
#                 output_changes: echo test2
#             """,
#         )
#         self.assertEqual(exit_code, 0)

#     def test_execution_is_correct(self):
#         exit_code, stdout = self.ebro("--file", "Ebro.inheritance.yaml", "default", "b")
#         self.assertEqual(exit_code, 0)
#         self.assertStdout(
#             stdout,
#             f"""
#             ███ [:c] running
#             I'm C
#             ███ [:d] running
#             I'm D
#             ███ [:b] running
#             It's BAR
#             ███ [:default] running
#             Hello
#             """,
#         )
