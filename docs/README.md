<div markdown="1" remove-in-website="1">

# Ebro Documentation

</div>

Ebro is a task runner. Tasks are defined inside YAML files, scripted with Bash, and configured with a name, requirements, when to skip execution, and other details. Tasks can also be [imported from other files](#importing-tasks) or [extend other tasks](#task-inheritance).

Ebro is distributed as a single binary, including the script interpreter ([mvdan/sh](https://github.com/mvdan/sh)).

It's heavily inspired in [go-task/task](https://github.com/go-task/task), but originally built around a personal need for configuring servers, although it's not tied to this use case and remains agnostic.

[TOC]

## Getting started

The format of `Ebro.yaml` files is defined [here](#the-ebroyaml-format). Here is an example:

```yaml
tasks:
  default:
    requires: [echoer, producer]

  echoer:
    script: |
      cat cache/A.txt
      cat cache/B.txt
    when:
      output_changes: |
        cat cache/A.txt
        cat cache/B.txt

  producer:
    requires: [produce_a, produce_b]

  produce_a:
    requires: [cache_dir]
    required_by: [echoer]
    script: echo 'this is A'>cache/A.txt
    when:
      check_fails: test -f cache/A.txt

  produce_b:
    requires: [cache_dir]
    required_by: [echoer]
    script: echo 'this is B'>cache/B.txt
    when:
      check_fails: test -f cache/B.txt

  cache_dir:
    script: mkdir -p cache
    when:
      check_fails: test -d cache
```

To give it a try, create a folder in your system, copy the content above in a file inside it called `Ebro.yaml`, and also download the `ebrow` script alongside it.

```bash
curl --location --output ebrow 'https://github.com/sirikon/ebro/releases/latest/download/ebrow'
```

**Before running it**, read the script and understand what it does (because you shouldn't blindly execute scripts from the internet).

`ebrow` is Ebro's "workspace script". It is a Bash script that contains a reference to an exact Ebro version and is able to download it, verify its integrity, and place it inside the `.ebro` directory created next to itself. The next time you execute it, it will use the already-downloaded Ebro binary. It ensures that the correct binary is present in the workspace.

Now, let's give it execution permissions and execute it:

```bash
chmod +x ebrow
./ebrow
```

Ebro on start will check for a file called `Ebro.yaml` in the working directory and parse it if present, constructing what is called the **inventory**, a collection of every task available with their definitive configuration for running.

You can check the inventory yourself by calling `./ebrow -inventory`, you'll notice extra details like the definitive list of extra environment variables that will be included in each task execution, or the working directory.

Next, it will construct a **plan**, which is an ordered list of all the tasks that will be executed sequentially in order to reach our target task, which by default is `default`, but can be any other by passing it as an argument (`./ebrow echoer`).

Again, check it yourself by running `./ebrow -plan`. This plan is deterministic, which means that given the same configuration, it will always be the same.

Finally, it will execute the plan, running tasks sequentially until the end.

<div class="x-tip" markdown="1">

Before running any Bash script in `script`, `when.output_changes` or `when.check_fails`, Ebro will prepend to the script the lines `set -euo pipefail` to ensure sane defaults:

- `-e`: Exit on error
- `-u`: Usage of unset variables is considered an error
- `-o pipefail`: The pipeline’s return status is the value of the last (rightmost) command to exit with a non-zero status, or zero if all commands exit successfully

More on Bash's documentation: [The Set Builtin](https://www.gnu.org/software/bash/manual/bash.html#The-Set-Builtin).

</div>

During the first execution it will execute everything, with no skips, which should output something like this:

```text
$ ./ebrow
███ [:cache_dir] running
███ [:produce_a] running
███ [:produce_b] running
███ [:echoer] running
this is A
this is B
███ [:producer] satisfied
███ [:default] satisfied
```

But if we execute `./ebrow` again, we'll see this output:

```text
$ ./ebrow
███ [:cache_dir] skipping
███ [:produce_a] skipping
███ [:produce_b] skipping
███ [:echoer] skipping
███ [:producer] satisfied
███ [:default] satisfied
```

Ebro skips tasks whenever possible, and the task definition is what mandates when a task should be skipped. In our example, the task `echoer` is skipped whenever the output of running `cat cache/A.txt` and `cat cache/B.txt` doesn't change. In the case of the task `produce_a`, it skips whenever the command `test -f cache/A.txt` succeeds, because the file `cache/A.txt` already exists.

Now we'll manually edit the file `cache/A.txt`, run `./ebrow` again, and see the result.

```text
$ echo 'hello world!' > cache/A.txt
$ ./ebrow
███ [:cache_dir] skipping
███ [:produce_a] skipping
███ [:produce_b] skipping
███ [:echoer] running
hello world!
this is B
███ [:producer] satisfied
███ [:default] satisfied
```

The `when.output_changes` checker of the `echoer` task detected that running `cat cache/A.txt` and `cat cache/B.txt` produced a different output when compared with the previous execution, hence, the task is executed again.

## Importing tasks

The `Ebro.yaml` file supports importing tasks from other `Ebro.yaml` files by defining the [`imports.*.from` parameter](#the-ebroyaml-format__imports.from). It works like this:

```yaml
# Ebro.yaml
imports:
  something:
    from: ./somewhere # a directory containing an `Ebro.yaml` file.
```

```yaml
# somewhere/Ebro.yaml
tasks:
  default:
    script: echo 'something'
  else:
    script: echo 'something else'
```

```text
$ ebro something
███ [:something:default] running
something
```

```text
$ ebro something:else
███ [:something:else] running
something else
```

It's important to note that the contents of an `Ebro.yaml` file are considered a **module**. When we import another `Ebro.yaml` file, we're creating a new module that hangs from the **root module** and has an explicitly-given name. In this case, `something`.

Targeting a module by its name is equivalent to targeting the module's `default` task. As `something` is a module, it translates to `something:default`.

## Conditional existence

Tasks can be configured to only exist when another task already exists using the `if_tasks_exist` parameter. Additionally, we can `require` tasks only if the referenced task exists by using the `?` suffix, and ignore the requirement otherwise.

With this configuration, as the task `restic` doesn't exist, `configure-backups` will not exist either, but that's okay, because `server`'s reference to it was optional.

```yaml
tasks:
  server:
    requires: [configure-backups?]
    script: |
      echo 'Configuring server'

  configure-backups:
    if_tasks_exist: [restic]
    requires: [restic]
    script: |
      echo 'Configuring backups'
```

```text
$ ebro server
███ [:server] running
Configuring server
```

But as soon as the `restic` task exists, this happens:

```yaml
tasks:
  server:
    requires: [configure-backups?]
    script: |
      echo 'Configuring server'

  configure-backups:
    if_tasks_exist: [restic]
    requires: [restic]
    script: |
      echo 'Configuring backups'

  restic:
    script: |
      echo 'Installing restic'
```

```text
$ ebro server
███ [:restic] running
Installing restic
███ [:configure-backups] running
Configuring backups
███ [:server] running
Configuring server
```

## Task inheritance

Ebro has a system of task inheritance. Tasks can extend other tasks, merging the parent properties with their own. Check the merging strategy in the [schema documentation](#the-ebroyaml-format__tasks.extends).

Here's an example `Ebro.yaml` file and what happens when running `ebro` on it:

```yaml
tasks:
  default:
    script: echo 'Hello World'

  parent:
    abstract: true
    environment:
      FOO: "foo"
    required_by: [default]

  child:
    extends: [parent]
    script: echo $FOO
```

Now, running `ebro default child` (targeting both `default` and `child` tasks) has this output:

```text
$ ebro default child
███ [:child] running
foo
███ [:default] running
Hello World
```

We can see how the `child` task looks after the merging strategy is applied with `ebro -inventory`. It inherited `parent`'s `required_by` and `FOO` environment variable, but kept its own `script`.

```yaml
:child:
  working_directory: /workdir
  environment:
    EBRO_ROOT: /workdir
    FOO: foo
  required_by:
  - :default
  script: echo $FOO
:default:
  working_directory: /workdir
  environment:
    EBRO_ROOT: /workdir
  script: echo 'Hello World'
```

## CLI

Ebro's command line interface is very straightforward, but has a couple of general rules:

- **Commands** define an specific action. At most, there is one command in a call. Absence of a command means the default command of _running_. Commands are prefixed with a single hyphen (`-command`).
- **Flags** depend on the command being executed. Their mere presence can mean a boolean value (`true` or `false`) or be accompanyed with a value. Flags are prefixed with two hyphens (`--flag`).
- **Targets** are the names of the tasks that we want to run. When no target is specified, the task `default` is assumed.

To know Ebro's available commands with their flags and explanations, run `ebro -help` (or `./ebrow -help` if using the workspace script).

```text
ebro [--flags...] [targets...]
  # Run everything
  flags:
    --file value  Specify the file that should be loaded as root module. default: Ebro.yaml
    --force       Ignore when.* conditionals and dont skip any task. default: false
  targets:
    defaults to [default]


ebro -inventory [--flags...]
  or -i
  # Display complete inventory of tasks with their definitive configuration in YAML format
  flags:
    --file value   Specify the file that should be loaded as root module. default: Ebro.yaml
    --query value  Query the inventory using an `expr` expression


ebro -list [--flags...]
  or -l
  # Display only the names of all the tasks in the inventory
  flags:
    --file value  Specify the file that should be loaded as root module. default: Ebro.yaml


ebro -plan [--flags...] [targets...]
  or -p
  # Display the execution plan
  flags:
    --file value  Specify the file that should be loaded as root module. default: Ebro.yaml
  targets:
    defaults to [default]


ebro -version
  or -v
  # Display ebro's version information in YAML format


ebro -help
  or -h
  # Display this help message
```

## Queries

### In `-inventory`

The `-inventory` (or `-i`) command supports a flag called `--query` for using an [Expr expression](https://expr-lang.org/docs/language-definition) that transforms the output at your will:

```text
$ ebro -i --query 'tasks | filter(.name == "default") | map(.id)'
- :apt:default
- :caddy:default
- :default
- :docker:default
- :docker:plugins:default
```

The output is serialized as `YAML` unless the data resulting from the expression is a **string**, in which case the result is printed to `stdout` verbatim.

```text
$ ebro -i --query 'tasks | filter(.name == "default") | map(.id) | join(", ")'
:apt:default, :caddy:default, :default, :docker:default, :docker:plugins:default
```

Here is the environment available during the expression execution apart from [Expr's built-in functions](https://expr-lang.org/docs/language-definition):

- `tasks`: Array of objects with the following properties:
  - `id`: (`string`) Task's ID (ex: `:module:task`)
  - `module`: (`string`) Task's module (ex: `:module`)
  - `name`: (`string`) Task's name (ex: `name`)
  - `working_directory`: (`string`)
  - `environment`: (`string` -> `string` dictionary)
  - `labels`: (`string` -> `string` dictionary)
  - `requires`: (`string` array)
  - `required_by`: (`string` array)
  - `script`: (`string` array)
  - `quiet`: (`bool`)
  - `interactive`: (`bool`)
  - `when`: Object with the following properties:
    - `check_fails`: (`string` array)
    - `output_changes`: (`string` array)
- `modules`: Array of objects with the following properties:
  - `id`: (`string`) Modules's ID (ex: `:module:submodule`)
  - `working_directory`: (`string`)
  - `environment`: (`string` -> `string` dictionary)
  - `labels`: (`string` -> `string` dictionary)

### In `requires` and `required_by`

The task properties [`requires`](#the-ebroyaml-format__tasks.requires) and [`required_by`](#the-ebroyaml-format__tasks.required_by) also support querying by expression, but with a difference: The Task properties `requires` and `required_by` are not available in the expression environment.

It is not possible to define `requires` or `required_by` based on the value of other tasks' `requires` or `required_by`.

## The `Ebro.yaml` format

This is also available as a [JSON Schema](./schema.json).

[Ebro.yaml format explained]

## Versioning

<div class="x-tip is-warning" markdown="1">

⚠️ Ebro is in version **0.x.x**, which means that the API isn't stable and could change at any time. This is covered in [SemVer's 4th spec item](https://semver.org/#spec-item-4):

_Anything MAY change at any time. The public API SHOULD NOT be considered stable._

</div>

Ebro follows [SemVer](https://semver.org/), this means that it's important to clarify what the API is in Ebro, or, "the exposed parts that will only break compatibility on major version releases".

Ebro's API in SemVer terms is composed of, but not limited to:

- The `Ebro.yaml` specification. The main idea is that an `Ebro.yaml` file that worked on Ebro version `1.x.x` **MUST** work exactly the same on any future version under the `1.x.x` major release. This includes, but is not limited to:
  - How `Ebro.yaml` is interpreted.
  - Environment variables added by Ebro during the execution.
    - **New** environment variables prefixed with `EBRO_` **COULD** be added, and this would not be considered a breaking change.
  - How scripts are interpreted and executed.
  - How task names are processed and referred to.
  - In which order are tasks executed.
- The CLI **actions**, meaning how the machine state changes given an `Ebro.yaml` file and a set of commands and flags.
  - How the internals of the `.ebro` directory operate **IS NOT** part of the API, and this could change at any time without being considered a breaking change.
- The CLI **stdout output** (**NOT** stderr) of the `-version`, `-plan`, `-inventory` and `-list` commands. These commands output a structured representation of data. New data could be added in the future, but existing data will remain the same. This doesn't apply to formatting, as it could change as long as it satisfies the output format specification.
  - The CLI stdout and stderr output of any other command apart from the ones explicitly mentioned above are **NOT** considered part of the API and could change at any time without it being considered a breaking change.


Unintended functionality that gets included in Ebro but is considered a bug afterwards will **NOT** be considered part of the API. **Bugs won't be kept for the sake of API stability and will be removed**. If this was the case, a release including a bugfix that is also a breaking change will have proper instructions for updating any affected `Ebro.yaml` file.
