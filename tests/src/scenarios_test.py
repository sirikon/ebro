from os import listdir, getcwd
from os.path import join
from utils.common import EbroTestCase
from subprocess import run
import shlex


class TestScenarios(EbroTestCase):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, **kwargs)

    def test_scenarios_work(self):
        scenarios = listdir(join(getcwd(), "scenarios"))
        for scenario in scenarios:
            scenarioPath = join(getcwd(), "scenarios", scenario)

            run(["rm", "-rf", ".ebro"], cwd=scenarioPath, check=True)
            run(["rm", "-rf", ".cache"], cwd=scenarioPath, check=True)

            with open(join(scenarioPath, "TEST.txt"), "r") as f:
                test_data = f.readlines()

            modes = ["meta", "output"]
            mode = 0
            currentTest = {"meta": {}, "output": ""}
            tests = []
            for line in test_data:
                match line:
                    case "---\n":
                        if mode == len(modes) - 1:
                            currentTest = {"meta": {}, "output": ""}
                            mode = 0
                        else:
                            mode = mode + 1
                    case line:
                        if currentTest not in tests:
                            tests.append(currentTest)
                        match modes[mode]:
                            case "meta":
                                if line == "\n":
                                    continue
                                [key, value] = [c.strip() for c in line.split(":", 1)]
                                if key not in currentTest["meta"]:
                                    currentTest["meta"][key] = []
                                currentTest["meta"][key].append(value)
                            case "output":
                                currentTest["output"] += line

            for i, test in enumerate(tests):
                args = [""]
                expected_exit_code = 0
                expected_output = test["output"]

                if "args" in test["meta"]:
                    args = test["meta"]["args"]

                if "exit_code" in test["meta"]:
                    expected_exit_code = int(test["meta"]["exit_code"][0])

                for y, a in enumerate(args):
                    with self.subTest(f"{scenario} {i} {y}"):
                        final_args = []
                        if a != "":
                            final_args = shlex.shlex(a, posix=True) if a != "" else []
                            final_args.whitespace_split = True
                            final_args = list(final_args)
                        exit_code, stdout = self.ebro(*final_args, cwd=scenarioPath)

                        expected_output = expected_output.replace(
                            "{{WORKDIR}}", scenarioPath
                        )
                        expected_output = expected_output.replace("{{BIN}}", self.bin)

                        self.assertEqual(stdout, expected_output)
                        self.assertEqual(exit_code, expected_exit_code)
