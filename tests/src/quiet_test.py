from utils.common import EbroTestCase


class TestQuiet(EbroTestCase):

    def test_inventory_is_correct(self):
        exit_code, stdout = self.ebro("-inventory", "--file", "Ebro.quiet.yaml")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            :fails:
              working_directory: {self.workdir}
              environment:
                EBRO_ROOT: {self.workdir}
              script: echo 'oh no' && exit 1
              quiet: true
            :works:
              working_directory: {self.workdir}
              environment:
                EBRO_ROOT: {self.workdir}
              script: echo hello
              quiet: true
            """,
        )

    def test_execution_is_correct_with_working_task(self):
        exit_code, stdout = self.ebro("--file", "Ebro.quiet.yaml", "works")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            ███ [:works] running
            """,
        )

    def test_execution_is_correct_with_failing_task(self):
        exit_code, stdout = self.ebro("--file", "Ebro.quiet.yaml", "fails")
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            f"""
            ███ [:fails] running
            oh no
            ███ ERROR: task :fails returned status code 1
            """,
        )
