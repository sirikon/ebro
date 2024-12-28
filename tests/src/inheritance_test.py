from utils.common import EbroTestCase


class TestInheritance(EbroTestCase):

    def test_inventory_is_correct(self):
        exit_code, stdout = self.ebro("-inventory", "--file", "Ebro.inheritance.yaml")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            :b:
                working_directory: {self.workdir}
                environment:
                    BAR: It's BAR
                    EBRO_ROOT: {self.workdir}
                    FOO: It's FOO
                requires:
                    - :d
                    - :c
                required_by:
                    - :default
                script: echo $BAR
                when:
                    check_fails: exit 2
                    output_changes: echo test
            :c:
                working_directory: {self.workdir}
                environment:
                    EBRO_ROOT: {self.workdir}
                script: echo "I'm C"
            :d:
                working_directory: {self.workdir}
                environment:
                    EBRO_ROOT: {self.workdir}
                script: echo "I'm D"
            :default:
                working_directory: {self.workdir}
                environment:
                    EBRO_ROOT: {self.workdir}
                script: echo Hello
            :multi-c:
                working_directory: {self.workdir}
                environment:
                    A: "1"
                    B: "22"
                    C: "3"
                    D: "44"
                    EBRO_ROOT: {self.workdir}
                script: echo multi-b
            :y:
                working_directory: {self.workdir}
                environment:
                    A: "1"
                    B: "2"
                    EBRO_ROOT: {self.workdir}
                script: echo y
                when:
                    check_fails: exit 1
                    output_changes: echo test2
            :z:
                working_directory: {self.workdir}
                environment:
                    A: "1"
                    B: "22"
                    C: "3"
                    EBRO_ROOT: {self.workdir}
                script: echo z
                when:
                    check_fails: exit 1
                    output_changes: echo test2
            """,
        )

    def test_execution_is_correct(self):
        exit_code, stdout = self.ebro("--file", "Ebro.inheritance.yaml")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            ███ [:c] running
            I'm C
            ███ [:d] running
            I'm D
            ███ [:b] running
            It's BAR
            ███ [:default] running
            Hello
            """,
        )
