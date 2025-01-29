from utils.common import EbroTestCase


class TestInventory(EbroTestCase):

    def test_inventory_is_correct(self):
        commands = ["-inventory", "-i"]
        for command in commands:
            with self.subTest(command):
                exit_code, stdout = self.ebro(command)
                self.assertEqual(exit_code, 0)
                self.assertStdout(
                    stdout,
                    f"""
                    :apt:default:
                      working_directory: {self.workdir}/apt/wd
                      environment:
                        DOCKER_MODULE_LOCATION: docker
                        DOCKER_PLUGINS_MODULE_LOCATION: {self.workdir}/docker_plugins
                        EBRO_ROOT: {self.workdir}
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/apt/wd
                      script: |
                        echo 'Installing apt packages'
                        cat "${{EBRO_ROOT}}/.cache/apt/packages/"*
                      when:
                        output_changes: cat "${{EBRO_ROOT}}/.cache/apt/packages/"*
                    :apt:pre-config:
                      working_directory: {self.workdir}/apt/wd
                      environment:
                        DOCKER_MODULE_LOCATION: docker
                        DOCKER_PLUGINS_MODULE_LOCATION: {self.workdir}/docker_plugins
                        EBRO_ROOT: {self.workdir}
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/apt/wd
                      script: mkdir -p "${{EBRO_ROOT}}/.cache/apt/packages"
                      when:
                        check_fails: test -d "${{EBRO_ROOT}}/.cache/apt/packages"
                    :caddy:default:
                      working_directory: {self.workdir}/caddy
                      environment:
                        DOCKER_MODULE_LOCATION: docker
                        DOCKER_PLUGINS_MODULE_LOCATION: {self.workdir}/docker_plugins
                        EBRO_ROOT: {self.workdir}
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/caddy
                      requires:
                      - :caddy:package
                    :caddy:package:
                      working_directory: {self.workdir}/caddy
                      environment:
                        DOCKER_MODULE_LOCATION: docker
                        DOCKER_PLUGINS_MODULE_LOCATION: {self.workdir}/docker_plugins
                        EBRO_ROOT: {self.workdir}
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/caddy
                      requires:
                      - :caddy:package-apt-config
                    :caddy:package-apt-config:
                      working_directory: {self.workdir}/caddy
                      environment:
                        DOCKER_MODULE_LOCATION: docker
                        DOCKER_PLUGINS_MODULE_LOCATION: {self.workdir}/docker_plugins
                        EBRO_ROOT: {self.workdir}
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/caddy
                      requires:
                      - :apt:pre-config
                      required_by:
                      - :apt:default
                      script: echo caddy > "${{EBRO_ROOT}}/.cache/apt/packages/caddy.txt"
                      when:
                        check_fails: test -f "${{EBRO_ROOT}}/.cache/apt/packages/caddy.txt"
                        output_changes: echo caddy
                    :default:
                      labels:
                        default: "true"
                      working_directory: {self.workdir}
                      environment:
                        DOCKER_MODULE_LOCATION: docker
                        DOCKER_PLUGINS_MODULE_LOCATION: {self.workdir}/docker_plugins
                        EBRO_ROOT: {self.workdir}
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}
                      requires:
                      - :apt:default
                      - :docker:default
                      - :caddy:default
                      script: |
                        echo "Done!"
                    :docker:default:
                      labels:
                        docker.version: 2.0.0
                      working_directory: {self.workdir}/docker
                      environment:
                        DOCKER_APT_VERSION: 2.0.0-1-apt
                        DOCKER_MODULE_LOCATION: docker
                        DOCKER_PLUGINS_MODULE_LOCATION: {self.workdir}/docker_plugins
                        DOCKER_VERSION: 2.0.0
                        EBRO_ROOT: {self.workdir}
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/docker
                      requires:
                      - :docker:package
                    :docker:package:
                      working_directory: {self.workdir}/docker
                      environment:
                        DOCKER_APT_VERSION: 2.0.0-1-apt
                        DOCKER_MODULE_LOCATION: docker
                        DOCKER_PLUGINS_MODULE_LOCATION: {self.workdir}/docker_plugins
                        DOCKER_VERSION: 2.0.0
                        EBRO_ROOT: {self.workdir}
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/docker
                      requires:
                      - :docker:package-apt-config
                    :docker:package-apt-config:
                      working_directory: {self.workdir}/docker
                      environment:
                        DOCKER_APT_VERSION: 2.0.0-1-apt
                        DOCKER_MODULE_LOCATION: docker
                        DOCKER_PLUGINS_MODULE_LOCATION: {self.workdir}/docker_plugins
                        DOCKER_VERSION: 2.0.0
                        EBRO_ROOT: {self.workdir}
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/docker
                      requires:
                      - :apt:pre-config
                      required_by:
                      - :apt:default
                      script: echo "docker==${{DOCKER_APT_VERSION}}" > "${{EBRO_ROOT}}/.cache/apt/packages/docker.txt"
                      when:
                        check_fails: test -f "${{EBRO_ROOT}}/.cache/apt/packages/docker.txt"
                        output_changes: echo "docker==${{DOCKER_APT_VERSION}}"
                    :docker:plugins:default:
                      working_directory: {self.workdir}/docker_plugins
                      environment:
                        DOCKER_APT_VERSION: 2.0.0-1-apt
                        DOCKER_MODULE_LOCATION: docker
                        DOCKER_PLUGINS_MODULE_LOCATION: {self.workdir}/docker_plugins
                        DOCKER_VERSION: 2.0.0
                        EBRO_ROOT: {self.workdir}
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/docker_plugins
                      script: echo Hello
                    :farm:chicken:
                      working_directory: {self.workdir}
                      environment:
                        DOCKER_MODULE_LOCATION: docker
                        DOCKER_PLUGINS_MODULE_LOCATION: {self.workdir}/docker_plugins
                        EBRO_ROOT: {self.workdir}
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}
                      requires:
                      - :farm:egg
                      script: echo Chicken ready
                    :farm:egg:
                      working_directory: {self.workdir}
                      environment:
                        DOCKER_MODULE_LOCATION: docker
                        DOCKER_PLUGINS_MODULE_LOCATION: {self.workdir}/docker_plugins
                        EBRO_ROOT: {self.workdir}
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}
                      script: echo 'Egg ready'
                    :farm:tractor:default:
                      working_directory: {self.workdir}/tractor
                      environment:
                        DOCKER_MODULE_LOCATION: docker
                        DOCKER_PLUGINS_MODULE_LOCATION: {self.workdir}/docker_plugins
                        EBRO_ROOT: {self.workdir}
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/tractor
                      script: echo "Tractor is here"
                    :ignored:
                      working_directory: {self.workdir}
                      environment:
                        DOCKER_MODULE_LOCATION: docker
                        DOCKER_PLUGINS_MODULE_LOCATION: {self.workdir}/docker_plugins
                        EBRO_ROOT: {self.workdir}
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}
                      required_by:
                      - :default
                      script: echo 'I should be ignored'
                    """,
                )

    def test_list_is_correct(self):
        commands = ["-list", "-l"]
        for command in commands:
            with self.subTest(command):
                exit_code, stdout = self.ebro(command)
                self.assertEqual(exit_code, 0)
                self.assertStdout(
                    stdout,
                    f"""
                    :apt:default
                    :apt:pre-config
                    :caddy:default
                    :caddy:package
                    :caddy:package-apt-config
                    :default
                    :docker:default
                    :docker:package
                    :docker:package-apt-config
                    :docker:plugins:default
                    :farm:chicken
                    :farm:egg
                    :farm:tractor:default
                    :ignored
                    """,
                )

    def test_inventory_with_absolute_workdir_is_correct(self):
        exit_code, stdout = self.ebro("-inventory", "--file", "Ebro.workdirs.yaml")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            :child:
              working_directory: /somewhere/absolute/child
              environment:
                ABSTRACT_WORKING_DIRECTORY: /somewhere/absolute/abstract
                EBRO_ROOT: {self.workdir}
                EBRO_TASK_WORKING_DIRECTORY: /somewhere/absolute/child
              script: echo $ABSTRACT_WORKING_DIRECTORY
            :default:
              working_directory: /somewhere/absolute
              environment:
                EBRO_ROOT: {self.workdir}
                EBRO_TASK_WORKING_DIRECTORY: /somewhere/absolute
              script: echo "Hello!"
            :other-absolute:
              working_directory: /other/absolute
              environment:
                EBRO_ROOT: {self.workdir}
                EBRO_TASK_WORKING_DIRECTORY: /other/absolute
              script: echo "Hello from the other absolute side!"
            :other-relative:
              working_directory: /somewhere/absolute/other/relative
              environment:
                EBRO_ROOT: {self.workdir}
                EBRO_TASK_WORKING_DIRECTORY: /somewhere/absolute/other/relative
              script: echo "Hello from the other relative side!"
            :submodule:other:
              working_directory: /somewhere/absolute/submodule
              environment:
                EBRO_ROOT: {self.workdir}
                EBRO_TASK_WORKING_DIRECTORY: /somewhere/absolute/submodule
              script: echo "Hello from the other side!"
            :submodule:other-absolute:
              working_directory: /other/absolute
              environment:
                EBRO_ROOT: {self.workdir}
                EBRO_TASK_WORKING_DIRECTORY: /other/absolute
              script: echo "Hello from the other absolute side!"
            :submodule:other-relative:
              working_directory: /somewhere/absolute/submodule/other/relative
              environment:
                EBRO_ROOT: {self.workdir}
                EBRO_TASK_WORKING_DIRECTORY: /somewhere/absolute/submodule/other/relative
              script: echo "Hello from the other relative side!"
            """,
        )

    def test_inventory_fails_with_task_with_nothing_to_do(self):
        exit_code, stdout = self.ebro(
            "-inventory", "--file", "Ebro.fail_when_nothing_to_do.yaml"
        )
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            f"""
            ███ ERROR: validating root module: validating task default: task has nothing to do (no requires, script, extends nor abstract)
            """,
        )

    def test_unkown_properties_are_not_allowed(self):
        exit_code, stdout = self.ebro(
            "-inventory", "--file", "Ebro.unknown_properties.yaml"
        )
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            f"""
            ███ ERROR: parsing root module: unmarshalling module file: [1:1] unknown field "import"
            >  1 | import:
                   ^
               2 |   apt:
               3 |     from: ./apt
            """,
        )
