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
                        EBRO_BIN: {self.bin}
                        EBRO_ROOT: {self.workdir}
                        EBRO_ROOT_FILE: {self.workdir}/Ebro.yaml
                        EBRO_TASK_ID: :apt:default
                        EBRO_TASK_MODULE: :apt
                        EBRO_TASK_NAME: default
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/apt/wd
                      script: |
                        echo 'Installing apt packages'
                        cat "${{EBRO_ROOT}}/.cache/apt/packages/"*
                      when:
                        output_changes: cat "${{EBRO_ROOT}}/.cache/apt/packages/"*
                    :apt:pre-config:
                      working_directory: {self.workdir}/apt/wd
                      environment:
                        EBRO_BIN: {self.bin}
                        EBRO_ROOT: {self.workdir}
                        EBRO_ROOT_FILE: {self.workdir}/Ebro.yaml
                        EBRO_TASK_ID: :apt:pre-config
                        EBRO_TASK_MODULE: :apt
                        EBRO_TASK_NAME: pre-config
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/apt/wd
                      script: mkdir -p "${{EBRO_ROOT}}/.cache/apt/packages"
                      when:
                        check_fails: test -d "${{EBRO_ROOT}}/.cache/apt/packages"
                    :bash:
                      working_directory: {self.workdir}
                      environment:
                        EBRO_BIN: {self.bin}
                        EBRO_ROOT: {self.workdir}
                        EBRO_ROOT_FILE: {self.workdir}/Ebro.yaml
                        EBRO_TASK_ID: :bash
                        EBRO_TASK_MODULE: ":"
                        EBRO_TASK_NAME: bash
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}
                      script: bash
                      interactive: true
                    :caddy:default:
                      working_directory: {self.workdir}/caddy
                      environment:
                        EBRO_BIN: {self.bin}
                        EBRO_ROOT: {self.workdir}
                        EBRO_ROOT_FILE: {self.workdir}/Ebro.yaml
                        EBRO_TASK_ID: :caddy:default
                        EBRO_TASK_MODULE: :caddy
                        EBRO_TASK_NAME: default
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/caddy
                      requires:
                      - :caddy:package
                    :caddy:package:
                      working_directory: {self.workdir}/caddy
                      environment:
                        EBRO_BIN: {self.bin}
                        EBRO_ROOT: {self.workdir}
                        EBRO_ROOT_FILE: {self.workdir}/Ebro.yaml
                        EBRO_TASK_ID: :caddy:package
                        EBRO_TASK_MODULE: :caddy
                        EBRO_TASK_NAME: package
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/caddy
                      requires:
                      - :caddy:package-apt-config
                    :caddy:package-apt-config:
                      working_directory: {self.workdir}/caddy
                      environment:
                        EBRO_BIN: {self.bin}
                        EBRO_ROOT: {self.workdir}
                        EBRO_ROOT_FILE: {self.workdir}/Ebro.yaml
                        EBRO_TASK_ID: :caddy:package-apt-config
                        EBRO_TASK_MODULE: :caddy
                        EBRO_TASK_NAME: package-apt-config
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
                        EBRO_BIN: {self.bin}
                        EBRO_ROOT: {self.workdir}
                        EBRO_ROOT_FILE: {self.workdir}/Ebro.yaml
                        EBRO_TASK_ID: :default
                        EBRO_TASK_MODULE: ":"
                        EBRO_TASK_NAME: default
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
                        EBRO_BIN: {self.bin}
                        EBRO_ROOT: {self.workdir}
                        EBRO_ROOT_FILE: {self.workdir}/Ebro.yaml
                        DOCKER_VERSION: 2.0.0
                        DOCKER_APT_VERSION: 2.0.0-1-apt
                        EBRO_TASK_ID: :docker:default
                        EBRO_TASK_MODULE: :docker
                        EBRO_TASK_NAME: default
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/docker
                      requires:
                      - :docker:package
                    :docker:package:
                      working_directory: {self.workdir}/docker
                      environment:
                        EBRO_BIN: {self.bin}
                        EBRO_ROOT: {self.workdir}
                        EBRO_ROOT_FILE: {self.workdir}/Ebro.yaml
                        DOCKER_VERSION: 2.0.0
                        DOCKER_APT_VERSION: 2.0.0-1-apt
                        EBRO_TASK_ID: :docker:package
                        EBRO_TASK_MODULE: :docker
                        EBRO_TASK_NAME: package
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/docker
                      requires:
                      - :docker:package-apt-config
                    :docker:package-apt-config:
                      working_directory: {self.workdir}/docker
                      environment:
                        EBRO_BIN: {self.bin}
                        EBRO_ROOT: {self.workdir}
                        EBRO_ROOT_FILE: {self.workdir}/Ebro.yaml
                        DOCKER_VERSION: 2.0.0
                        DOCKER_APT_VERSION: 2.0.0-1-apt
                        EBRO_TASK_ID: :docker:package-apt-config
                        EBRO_TASK_MODULE: :docker
                        EBRO_TASK_NAME: package-apt-config
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
                        EBRO_BIN: {self.bin}
                        EBRO_ROOT: {self.workdir}
                        EBRO_ROOT_FILE: {self.workdir}/Ebro.yaml
                        DOCKER_VERSION: 2.0.0
                        DOCKER_APT_VERSION: 2.0.0-1-apt
                        EBRO_TASK_ID: :docker:plugins:default
                        EBRO_TASK_MODULE: :docker:plugins
                        EBRO_TASK_NAME: default
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/docker_plugins
                      script: echo Hello
                    :farm:chicken:
                      working_directory: {self.workdir}
                      environment:
                        EBRO_BIN: {self.bin}
                        EBRO_ROOT: {self.workdir}
                        EBRO_ROOT_FILE: {self.workdir}/Ebro.yaml
                        EBRO_TASK_ID: :farm:chicken
                        EBRO_TASK_MODULE: :farm
                        EBRO_TASK_NAME: chicken
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}
                      requires:
                      - :farm:egg
                      script: echo Chicken ready
                    :farm:egg:
                      working_directory: {self.workdir}
                      environment:
                        EBRO_BIN: {self.bin}
                        EBRO_ROOT: {self.workdir}
                        EBRO_ROOT_FILE: {self.workdir}/Ebro.yaml
                        EBRO_TASK_ID: :farm:egg
                        EBRO_TASK_MODULE: :farm
                        EBRO_TASK_NAME: egg
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}
                      script: echo 'Egg ready'
                    :farm:tractor:default:
                      working_directory: {self.workdir}/tractor
                      environment:
                        EBRO_BIN: {self.bin}
                        EBRO_ROOT: {self.workdir}
                        EBRO_ROOT_FILE: {self.workdir}/Ebro.yaml
                        EBRO_TASK_ID: :farm:tractor:default
                        EBRO_TASK_MODULE: :farm:tractor
                        EBRO_TASK_NAME: default
                        EBRO_TASK_WORKING_DIRECTORY: {self.workdir}/tractor
                      script: echo "Tractor is here"
                    :ignored:
                      working_directory: {self.workdir}
                      environment:
                        EBRO_BIN: {self.bin}
                        EBRO_ROOT: {self.workdir}
                        EBRO_ROOT_FILE: {self.workdir}/Ebro.yaml
                        EBRO_TASK_ID: :ignored
                        EBRO_TASK_MODULE: ":"
                        EBRO_TASK_NAME: ignored
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
                    :bash
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

    def test_inventory_with_malformed_query_fails(self):
        commands = ["-inventory", "-i"]
        for command in commands:
            with self.subTest(command):
                exit_code, stdout = self.ebro(
                    command,
                    "--query",
                    'joins(tasks | filter(.labels.default == "true") | map(.id), "\\n")',
                )
                self.assertStdout(
                    stdout,
                    f"""
███ ERROR: compiling query expression: unknown name joins (1:1)
 | joins(tasks | filter(.labels.default == "true") | map(.id), "\\n")
 | ^
                    """,
                )
                self.assertEqual(exit_code, 1)

    def test_inventory_with_query_works(self):
        commands = ["-inventory", "-i"]
        for command in commands:
            with self.subTest(command):
                exit_code, stdout = self.ebro(
                    command,
                    "--query",
                    'join(tasks | filter(.labels.default == "true") | map(.id), "\\n")',
                )
                self.assertStdout(
                    stdout,
                    f"""
                    :default
                    """,
                )
                self.assertEqual(exit_code, 0)

    def test_inventory_with_query_works_2(self):
        commands = ["-inventory", "-i"]
        for command in commands:
            with self.subTest(command):
                exit_code, stdout = self.ebro(
                    command,
                    "--query",
                    'join(tasks | filter(.module == ":apt") | map(.id), "\\n")',
                )
                self.assertStdout(
                    stdout,
                    f"""
                    :apt:default
                    :apt:pre-config
                    """,
                )
                self.assertEqual(exit_code, 0)

    def test_inventory_with_query_works_3(self):
        commands = ["-inventory", "-i"]
        for command in commands:
            with self.subTest(command):
                exit_code, stdout = self.ebro(
                    command,
                    "--query",
                    'join(tasks | filter(.id == ":apt:default") | map(.id), "\\n")',
                )
                self.assertStdout(
                    stdout,
                    f"""
                    :apt:default
                    """,
                )
                self.assertEqual(exit_code, 0)

    def test_inventory_with_query_works_4(self):
        commands = ["-inventory", "-i"]
        for command in commands:
            with self.subTest(command):
                exit_code, stdout = self.ebro(
                    command,
                    "--query",
                    'join(tasks | filter(.name == "default") | map(.id), "\\n")',
                )
                self.assertStdout(
                    stdout,
                    f"""
                    :apt:default
                    :caddy:default
                    :default
                    :docker:default
                    :docker:plugins:default
                    :farm:tractor:default
                    """,
                )
                self.assertEqual(exit_code, 0)

    def test_inventory_with_query_works_4(self):
        commands = ["-inventory", "-i"]
        for command in commands:
            with self.subTest(command):
                exit_code, stdout = self.ebro(
                    command,
                    "--query",
                    'tasks | filter(.name == "default") | map(.id)',
                )
                self.assertStdout(
                    stdout,
                    f"""
                    - :apt:default
                    - :caddy:default
                    - :default
                    - :docker:default
                    - :docker:plugins:default
                    - :farm:tractor:default
                    """,
                )
                self.assertEqual(exit_code, 0)

    def test_inventory_with_absolute_workdir_is_correct(self):
        exit_code, stdout = self.ebro("-inventory", "--file", "Ebro.workdirs.yaml")
        self.assertEqual(exit_code, 0)
        self.assertStdout(
            stdout,
            f"""
            :child:
              working_directory: /somewhere/absolute/child
              environment:
                EBRO_BIN: {self.bin}
                EBRO_ROOT: {self.workdir}
                EBRO_ROOT_FILE: {self.workdir}/Ebro.workdirs.yaml
                ABSTRACT_WORKING_DIRECTORY: /somewhere/absolute/abstract
                EBRO_TASK_ID: :child
                EBRO_TASK_MODULE: ":"
                EBRO_TASK_NAME: child
                EBRO_TASK_WORKING_DIRECTORY: /somewhere/absolute/child
              script: echo $ABSTRACT_WORKING_DIRECTORY
            :default:
              working_directory: /somewhere/absolute
              environment:
                EBRO_BIN: {self.bin}
                EBRO_ROOT: {self.workdir}
                EBRO_ROOT_FILE: {self.workdir}/Ebro.workdirs.yaml
                EBRO_TASK_ID: :default
                EBRO_TASK_MODULE: ":"
                EBRO_TASK_NAME: default
                EBRO_TASK_WORKING_DIRECTORY: /somewhere/absolute
              script: echo "Hello!"
            :other-absolute:
              working_directory: /other/absolute
              environment:
                EBRO_BIN: {self.bin}
                EBRO_ROOT: {self.workdir}
                EBRO_ROOT_FILE: {self.workdir}/Ebro.workdirs.yaml
                EBRO_TASK_ID: :other-absolute
                EBRO_TASK_MODULE: ":"
                EBRO_TASK_NAME: other-absolute
                EBRO_TASK_WORKING_DIRECTORY: /other/absolute
              script: echo "Hello from the other absolute side!"
            :other-relative:
              working_directory: /somewhere/absolute/other/relative
              environment:
                EBRO_BIN: {self.bin}
                EBRO_ROOT: {self.workdir}
                EBRO_ROOT_FILE: {self.workdir}/Ebro.workdirs.yaml
                EBRO_TASK_ID: :other-relative
                EBRO_TASK_MODULE: ":"
                EBRO_TASK_NAME: other-relative
                EBRO_TASK_WORKING_DIRECTORY: /somewhere/absolute/other/relative
              script: echo "Hello from the other relative side!"
            :submodule:other:
              working_directory: /somewhere/absolute/submodule
              environment:
                EBRO_BIN: {self.bin}
                EBRO_ROOT: {self.workdir}
                EBRO_ROOT_FILE: {self.workdir}/Ebro.workdirs.yaml
                EBRO_TASK_ID: :submodule:other
                EBRO_TASK_MODULE: :submodule
                EBRO_TASK_NAME: other
                EBRO_TASK_WORKING_DIRECTORY: /somewhere/absolute/submodule
              script: echo "Hello from the other side!"
            :submodule:other-absolute:
              working_directory: /other/absolute
              environment:
                EBRO_BIN: {self.bin}
                EBRO_ROOT: {self.workdir}
                EBRO_ROOT_FILE: {self.workdir}/Ebro.workdirs.yaml
                EBRO_TASK_ID: :submodule:other-absolute
                EBRO_TASK_MODULE: :submodule
                EBRO_TASK_NAME: other-absolute
                EBRO_TASK_WORKING_DIRECTORY: /other/absolute
              script: echo "Hello from the other absolute side!"
            :submodule:other-relative:
              working_directory: /somewhere/absolute/submodule/other/relative
              environment:
                EBRO_BIN: {self.bin}
                EBRO_ROOT: {self.workdir}
                EBRO_ROOT_FILE: {self.workdir}/Ebro.workdirs.yaml
                EBRO_TASK_ID: :submodule:other-relative
                EBRO_TASK_MODULE: :submodule
                EBRO_TASK_NAME: other-relative
                EBRO_TASK_WORKING_DIRECTORY: /somewhere/absolute/submodule/other/relative
              script: echo "Hello from the other relative side!"
            """,
        )

    def test_inventory_fails_with_task_with_nothing_to_do(self):
        exit_code, stdout = self.ebro(
            "-inventory", "--file", "Ebro.fail_when_nothing_to_do.yaml"
        )
        self.assertStdout(
            stdout,
            f"""
            ███ ERROR: parsing module: parsing 'tasks': parsing task 'default': task has nothing to do (no requires, script, extends nor abstract)
            """,
        )
        self.assertEqual(exit_code, 1)

    def test_unkown_properties_are_not_allowed(self):
        exit_code, stdout = self.ebro(
            "-inventory", "--file", "Ebro.unknown_properties.yaml"
        )
        self.assertEqual(exit_code, 1)
        self.assertStdout(
            stdout,
            f"""
            ███ ERROR: parsing module: unexpected key 'import'
            """,
        )
