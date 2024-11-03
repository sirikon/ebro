from subprocess import run, PIPE
from os import getcwd, listdir
from os.path import join, isdir
import unittest


class TestEbro(unittest.TestCase):

    def __init__(self, methodName: str = "runTest") -> None:
        self.maxDiff = None
        super().__init__(methodName)

    def test_cases(self):
        cases_dir = join(getcwd(), "cases")
        for case in listdir(cases_dir):
            case_path = join(cases_dir, case)
            workdir_path = join(case_path, "workdir")
            with self.subTest(case + " (config)"):
                result = ebro(["-config"], workdir_path)
                actual_stdout = result.stdout.decode("utf-8")
                with open(join(case_path, "expected_config.txt"), "r") as f:
                    expected_stdout = f.read().replace("__WORKDIR__", workdir_path)
                self.assertEqual(actual_stdout, expected_stdout)

            with self.subTest(case + " (catalog)"):
                result = ebro(["-catalog"], workdir_path)
                actual_stdout = result.stdout.decode("utf-8")
                with open(join(case_path, "expected_catalog.txt"), "r") as f:
                    expected_stdout = f.read().replace("__WORKDIR__", workdir_path)
                self.assertEqual(actual_stdout, expected_stdout)

            plans_dir = join(case_path, "plans")
            for plan in listdir(plans_dir) if isdir(plans_dir) else []:
                target = plan.removesuffix(".txt")
                with self.subTest(case + " (plan) (" + target + ")"):
                    plan_path = join(plans_dir, plan)
                    if target == "default":
                        result = ebro(["-plan"], workdir_path)
                    else:
                        result = ebro(["-plan", target], workdir_path)
                    actual_stdout = result.stdout.decode("utf-8")
                    with open(plan_path, "r") as f:
                        expected_stdout = f.read()
                    self.assertEqual(actual_stdout, expected_stdout)


def ebro(args, cwd):
    return run(
        [join(getcwd(), "..", "out", "ebro"), *args],
        cwd=cwd,
        check=True,
        stdout=PIPE,
        stderr=PIPE,
    )
