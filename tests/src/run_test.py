from utils.common import EbroTestCase


class TestRun(EbroTestCase):

    def test_execution_is_correct(self):
        exit_code, stdout = self.ebro()
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            ███ [:apt:pre-config] running
            + mkdir -p {self.workdir}/.cache/apt/packages
            ███ [:caddy:package-apt-config] running
            + echo caddy
            ███ [:docker:package-apt-config] running
            + echo 'docker==2.0.0-1-apt'
            ███ [:apt:default] running
            + echo 'Installing apt packages'
            Installing apt packages
            + cat {self.workdir}/.cache/apt/packages/caddy.txt {self.workdir}/.cache/apt/packages/docker.txt
            caddy
            docker==2.0.0-1-apt
            ███ [:caddy:package] satisfied
            ███ [:docker:package] satisfied
            ███ [:caddy:default] satisfied
            ███ [:docker:default] satisfied
            ███ [:default] running
            + echo Done!
            Done!
            """,
        )

        # Second execution should cache everything except "Done!"
        exit_code, stdout = self.ebro()
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            ███ [:apt:pre-config] skipping
            ███ [:caddy:package-apt-config] skipping
            ███ [:docker:package-apt-config] skipping
            ███ [:apt:default] skipping
            ███ [:caddy:package] satisfied
            ███ [:docker:package] satisfied
            ███ [:caddy:default] satisfied
            ███ [:docker:default] satisfied
            ███ [:default] running
            + echo Done!
            Done!
            """,
        )

        # Third execution, with force, should look like the first one
        exit_code, stdout = self.ebro("--force")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            ███ [:apt:pre-config] running
            + mkdir -p {self.workdir}/.cache/apt/packages
            ███ [:caddy:package-apt-config] running
            + echo caddy
            ███ [:docker:package-apt-config] running
            + echo 'docker==2.0.0-1-apt'
            ███ [:apt:default] running
            + echo 'Installing apt packages'
            Installing apt packages
            + cat {self.workdir}/.cache/apt/packages/caddy.txt {self.workdir}/.cache/apt/packages/docker.txt
            caddy
            docker==2.0.0-1-apt
            ███ [:caddy:package] satisfied
            ███ [:docker:package] satisfied
            ███ [:caddy:default] satisfied
            ███ [:docker:default] satisfied
            ███ [:default] running
            + echo Done!
            Done!
            """,
        )

    def test_scripts_fail_asap(self):
        exit_code, stdout = self.ebro("--file", "Ebro.fail_asap.yaml")
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            f"""
            ███ [:default] running
            + echo 'This should print'
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
            + echo 'This should print all the time'
            This should print all the time
            + exit 1
            ███ ERROR: task :default returned status code 1
            """,
        )

        exit_code, stdout = self.ebro("--file", "Ebro.fail_with_output_change.yaml")
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            f"""
            ███ [:default] running
            + echo 'This should print all the time'
            This should print all the time
            + exit 1
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
            + echo B
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
            + echo Running
            Running
            ███ [:never_fails] running
            + echo Running
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
            + echo Running
            Running
            ███ [:never_fails] skipping
            ███ [:default] satisfied
            """,
        )
