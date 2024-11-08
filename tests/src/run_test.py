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
            docker==1.0.0-1-apt
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