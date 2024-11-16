from os import environ
from utils.common import EbroTestCase, fake_git_server


class TestGitImport(EbroTestCase):

    @fake_git_server
    def test_git_import_works(self, repository_url):
        exit_code, stdout = self.ebro(
            "--file",
            "Ebro.git_import_test.yaml",
            env={"TEST_REPOSITORY_URL": repository_url},
        )
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            ███ cloning {repository_url}
            ███ [:apt:pre-config] running
            ███ [:caddy:package-apt-config] running
            ███ [:apt:default] running
            Installing apt packages
            caddy
            ███ [:caddy:package] satisfied
            ███ [:caddy:default] satisfied
            ███ [:default] running
            Done!
            """,
        )

        # second call should not clone again
        exit_code, stdout = self.ebro(
            "--file",
            "Ebro.git_import_test.yaml",
            env={"TEST_REPOSITORY_URL": repository_url},
        )
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            ███ [:apt:pre-config] skipping
            ███ [:caddy:package-apt-config] skipping
            ███ [:apt:default] skipping
            ███ [:caddy:package] satisfied
            ███ [:caddy:default] satisfied
            ███ [:default] running
            Done!
            """,
        )
