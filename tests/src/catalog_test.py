from utils.common import EbroTestCase


class TestInventory(EbroTestCase):

    def test_inventory_is_correct(self):
        exit_code, stdout = self.ebro("-inventory")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            :apt:default:
                working_directory: {self.workdir}/apt/wd
                environment:
                    EBRO_ROOT: {self.workdir}
                script: |
                    echo 'Installing apt packages'
                    cat "${{EBRO_ROOT}}/.cache/apt/packages/"*
                when:
                    output_changes: cat "${{EBRO_ROOT}}/.cache/apt/packages/"*
            :apt:pre-config:
                working_directory: {self.workdir}/apt/wd
                environment:
                    EBRO_ROOT: {self.workdir}
                script: mkdir -p "${{EBRO_ROOT}}/.cache/apt/packages"
                when:
                    check_fails: test -d "${{EBRO_ROOT}}/.cache/apt/packages"
            :caddy:default:
                working_directory: {self.workdir}/caddy
                environment:
                    EBRO_ROOT: {self.workdir}
                requires:
                    - :caddy:package
            :caddy:package:
                working_directory: {self.workdir}/caddy
                environment:
                    EBRO_ROOT: {self.workdir}
                requires:
                    - :apt:default
                    - :caddy:package-apt-config
            :caddy:package-apt-config:
                working_directory: {self.workdir}/caddy
                environment:
                    EBRO_ROOT: {self.workdir}
                requires:
                    - :apt:pre-config
                required_by:
                    - :apt:default
                script: echo caddy > "${{EBRO_ROOT}}/.cache/apt/packages/caddy.txt"
                when:
                    check_fails: test -f "${{EBRO_ROOT}}/.cache/apt/packages/caddy.txt"
                    output_changes: echo caddy
            :chicken:
                working_directory: {self.workdir}
                environment:
                    EBRO_ROOT: {self.workdir}
                requires:
                    - :egg
                script: echo Chicken ready
            :default:
                working_directory: {self.workdir}
                environment:
                    EBRO_ROOT: {self.workdir}
                requires:
                    - :docker:default
                    - :caddy:default
                script: |
                    echo "Done!"
            :docker:default:
                working_directory: {self.workdir}/docker
                environment:
                    DOCKER_APT_VERSION: 2.0.0-1-apt
                    DOCKER_VERSION: 2.0.0
                    EBRO_ROOT: {self.workdir}
                requires:
                    - :docker:package
            :docker:package:
                working_directory: {self.workdir}/docker
                environment:
                    DOCKER_APT_VERSION: 2.0.0-1-apt
                    DOCKER_VERSION: 2.0.0
                    EBRO_ROOT: {self.workdir}
                requires:
                    - :apt:default
                    - :docker:package-apt-config
            :docker:package-apt-config:
                working_directory: {self.workdir}/docker
                environment:
                    DOCKER_APT_VERSION: 2.0.0-1-apt
                    DOCKER_VERSION: 2.0.0
                    EBRO_ROOT: {self.workdir}
                requires:
                    - :apt:pre-config
                required_by:
                    - :apt:default
                script: echo "docker==${{DOCKER_APT_VERSION}}" > "${{EBRO_ROOT}}/.cache/apt/packages/docker.txt"
                when:
                    check_fails: test -f "${{EBRO_ROOT}}/.cache/apt/packages/docker.txt"
                    output_changes: echo "docker==${{DOCKER_APT_VERSION}}"
            :egg:
                working_directory: {self.workdir}
                environment:
                    EBRO_ROOT: {self.workdir}
                script: echo 'Egg ready'
            """,
        )

    def test_inventory_with_absolute_workdir_is_correct(self):
        exit_code, stdout = self.ebro(
            "-inventory", "--file", "Ebro.workdirs.yaml"
        )
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            :default:
                working_directory: /somewhere/absolute
                environment:
                    EBRO_ROOT: {self.workdir}
                script: echo "Hello!"
            :other-absolute:
                working_directory: /other/absolute
                environment:
                    EBRO_ROOT: {self.workdir}
                script: echo "Hello from the other absolute side!"
            :other-relative:
                working_directory: /somewhere/absolute/other/relative
                environment:
                    EBRO_ROOT: {self.workdir}
                script: echo "Hello from the other relative side!"
            """,
        )
