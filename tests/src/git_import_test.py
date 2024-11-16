from os import environ
from utils.common import EbroTestCase, fake_git_server


class TestGitImport(EbroTestCase):

    @fake_git_server
    def test_git_import_works(self, repository_url):
        exit_code, stdout = self.ebro(
            "--file", "Ebro.fake_git.yaml", env={"TEST_REPOSITORY_URL": repository_url}
        )
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            ███ cloning {repository_url}
            ███ [:fake:default] running
            I am fake
            ███ [:default] running
            Done!
            """,
        )
