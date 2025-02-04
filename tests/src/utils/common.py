import json
from time import sleep
import unittest
from os import getcwd, environ
from os.path import join, normpath
from subprocess import run, PIPE, STDOUT
from textwrap import dedent
import http.client


def fake_git_server(func):
    def wrapper(*args, **kwargs):
        result = _docker(
            "run",
            "-d",
            "-p",
            "80",
            "-v",
            join(getcwd(), "..", "playground") + ":/content:ro",
            "-v",
            join(getcwd(), "./fake_git_server/entrypoint.sh") + ":/entrypoint.sh:ro",
            "ynohat/git-http-backend",
            "/entrypoint.sh",
        )
        container_id = result.stdout.decode("utf-8").strip()
        container_info = _docker_inspect(container_id)
        local_port = container_info[0]["NetworkSettings"]["Ports"]["80/tcp"][0][
            "HostPort"
        ]
        repository_url = f"http://127.0.0.1:{local_port}/git/fake.git"
        retries = 300
        while retries > 0:
            retries = retries - 1
            try:
                sleep(0.01)
                conn = http.client.HTTPConnection(f"127.0.0.1:{local_port}")
                conn.request("GET", "/")
                conn.getresponse()
                conn.close()
                break
            except Exception:
                pass
        try:
            func(*args, repository_url, **kwargs)
        finally:
            _docker("rm", "-f", container_id)

    return wrapper


def _docker_inspect(id):
    result = _docker("inspect", id)
    return json.loads(result.stdout)


def _docker(*args):
    return run(
        ["docker", *args],
        check=True,
        stdin=None,
        stdout=PIPE,
        stderr=PIPE,
    )


class EbroTestCase(unittest.TestCase):
    def __init__(self, *args, **kwargs) -> None:
        self.maxDiff = None
        self.workdir = normpath(join(getcwd(), "..", "playground"))
        self.bin = normpath(join(getcwd(), "..", "out", "ebro-e2e"))
        super().__init__(*args, **kwargs)

    def setUp(self) -> None:
        run(["rm", "-rf", ".ebro"], cwd=self.workdir, check=True)
        run(["rm", "-rf", ".cache"], cwd=self.workdir, check=True)
        return super().setUp()

    def assertStdout(self, first: str, second: str):
        self.assertEqual(dedent(second).strip() + "\n", first)

    def assertStdoutStrict(self, first: str, second: str):
        self.assertEqual(second, first)

    def ebro(self, *args, env=dict()):
        result = run(
            [self.bin, *args],
            cwd=self.workdir,
            stdout=PIPE,
            stderr=STDOUT,
            env=environ | dict(GOCOVERDIR=join(getcwd(), ".coverage")) | env,
        )
        return result.returncode, result.stdout.decode("utf-8")
