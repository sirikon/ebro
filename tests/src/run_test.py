from utils.common import EbroTestCase


class TestRun(EbroTestCase):

    def test_execution_is_correct(self):
        exit_code, stdout = self.ebro()
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            ███ [:apt:pre-config] running
            ███ [:caddy:package-apt-config] running
            ███ [:docker:package-apt-config] running
            ███ [:apt:default] running
            Installing apt packages
            caddy
            docker==2.0.0-1-apt
            ███ [:caddy:package] satisfied
            ███ [:docker:package] satisfied
            ███ [:caddy:default] satisfied
            ███ [:docker:default] satisfied
            ███ [:default] running
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
            ███ [:caddy:package-apt-config] running
            ███ [:docker:package-apt-config] running
            ███ [:apt:default] running
            Installing apt packages
            caddy
            docker==2.0.0-1-apt
            ███ [:caddy:package] satisfied
            ███ [:docker:package] satisfied
            ███ [:caddy:default] satisfied
            ███ [:docker:default] satisfied
            ███ [:default] running
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
