from utils.common import EbroTestCase


class TestRun(EbroTestCase):

    def test_scripts_fail_asap(self):
        exit_code, stdout = self.ebro("--file", "Ebro.fail_asap.yaml")
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            f"""
            ███ [:default] running
            This should print
            UNBOUND_VARIABLE: unbound variable
            ███ ERROR: task :default returned status code 1
            """,
        )

    def test_failing_scripts_are_not_cached_by_output_changes(self):
        exit_code, stdout = self.ebro("--file", "Ebro.fail_with_output_change.yaml")
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            f"""
            ███ [:default] running
            This should print all the time
            ███ ERROR: task :default returned status code 1
            """,
        )

        exit_code, stdout = self.ebro("--file", "Ebro.fail_with_output_change.yaml")
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            f"""
            ███ [:default] running
            This should print all the time
            ███ ERROR: task :default returned status code 1
            """,
        )

    def test_required_by_does_not_include_referenced_task_in_plan(self):
        exit_code, stdout = self.ebro("--file", "Ebro.required_by_not_includes.yaml")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            ███ [:b] running
            B
            ███ [:default] satisfied
            """,
        )

    def test_when_checkers_behave_correctly(self):
        exit_code, stdout = self.ebro("--file", "Ebro.when_checkers_are_OR.yaml")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            ███ [:always_fails] running
            Running
            ███ [:never_fails] running
            Running
            ███ [:default] satisfied
            """,
        )

        exit_code, stdout = self.ebro("--file", "Ebro.when_checkers_are_OR.yaml")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            ███ [:always_fails] running
            Running
            ███ [:never_fails] skipping
            ███ [:default] satisfied
            """,
        )

    def test_env_interpolation_works_with_external_env_vars(self):
        exit_code, stdout = self.ebro(
            "--file",
            "Ebro.env.yaml",
            env=dict(EXTERNAL_MESSAGE="This is the external message"),
        )
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            ███ [:default] running
            This is the external message
            """,
        )

    def test_cmd_exec_in_env_is_not_allowed(self):
        exit_code, stdout = self.ebro("--file", "Ebro.fail_on_env_cmd.yaml")
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            f"""
            ███ ERROR: resolving task ':default' environment: expanding $(pwd): unexpected command substitution at 1:1
            """,
        )

    def test_bad_names_are_handled(self):
        exit_code, stdout = self.ebro("--file", "Ebro.bad_names.yaml")
        self.assertStdout(
            stdout,
            f"""
            ███ ERROR: parsing module: parsing 'tasks': validating task name 'dëfault': name does not match the following regex: ^[a-zA-Z0-9-_\\.]+$
            """,
        )
        self.assertEqual(exit_code, 1)

    def test_bad_names_are_handled_2(self):
        exit_code, stdout = self.ebro("--file", "Ebro.bad_names_2.yaml")
        self.assertStdout(
            stdout,
            f"""
            ███ ERROR: parsing module: parsing 'modules': validating module name 'ñodule': name does not match the following regex: ^[a-zA-Z0-9-_\\.]+$
            """,
        )
        self.assertEqual(exit_code, 1)

    def test_conditional_existence_works_1(self):
        exit_code, stdout = self.ebro(
            "--file", "Ebro.conditional_existence_1.yaml", "server"
        )
        self.assertStdout(
            stdout,
            f"""
            ███ [:server] running
            Configuring server
            """,
        )
        self.assertEqual(exit_code, 0)

    def test_conditional_existence_works_2(self):
        exit_code, stdout = self.ebro(
            "--file", "Ebro.conditional_existence_2.yaml", "server"
        )
        self.assertStdout(
            stdout,
            f"""
            ███ [:restic] running
            Installing restic
            ███ [:configure-backups] running
            Configuring backups
            ███ [:server] running
            Configuring server
            """,
        )
        self.assertEqual(exit_code, 0)
