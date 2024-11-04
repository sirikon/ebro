from utils.common import EbroTestCase


class TestConfig(EbroTestCase):

    def test_config_is_correct(self):
        exit_code, stdout = self.ebro("-config")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            working_directory: {self.workdir}
            tasks:
                chicken:
                    requires:
                        - egg
                    script: echo Chicken ready
                default:
                    requires:
                        - docker
                    script: |
                        echo "Done!"
                egg:
                    script: echo 'Egg ready'
            modules:
                apt:
                    working_directory: {self.workdir}/apt
                    tasks:
                        default:
                            script: |
                                echo 'Installing apt packages'
                                cat "${{EBRO_ROOT}}/.cache/apt/packages.txt"
                            sources:
                                - ${{EBRO_ROOT}}/.cache/apt/packages.txt
                        pre-config:
                            script: mkdir -p "${{EBRO_ROOT}}/.cache/apt"
                            skip_if: test -d "${{EBRO_ROOT}}/.cache/apt"
                docker:
                    working_directory: {self.workdir}/docker
                    environment:
                        DOCKER_APT_VERSION: ${{DOCKER_VERSION:-1.0.0}}-1-apt
                    tasks:
                        default:
                            requires:
                                - package
                        package:
                            requires:
                                - :apt
                                - package-apt-config
                        package-apt-config:
                            requires:
                                - :apt:pre-config
                            required_by:
                                - :apt
                            script: echo "docker==${{DOCKER_APT_VERSION}}" > "${{EBRO_ROOT}}/.cache/apt/packages.txt"
            """,
        )
