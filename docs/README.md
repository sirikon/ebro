<div markdown="1" remove-in-website="1">

# Ebro Documentation

</div>

Ebro is a task runner. Tasks are defined inside YAML files, scripted with Bash, and configured with a name, requirements, when to skip execution, and other details. Tasks can also be [imported from other files](#importing-tasks) or [extend other tasks](#task-inheritance).

Ebro is distributed as a single binary, including the script interpreter ([mvdan/sh](https://github.com/mvdan/sh)).

It's heavily inspired in [go-task/task](https://github.com/go-task/task), but originally built around a personal need for configuring servers, although it's not tied to this use case and remains agnostic.

[TOC]

## Getting started

The format of Ebro.yaml files is defined [here](#the-ebroyaml-format). Here is an example:

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

<details markdown="1">
<summary>
All scripts include <code>set -euo pipefail</code>.
</summary>

Before running any Bash script in `script`, `when.output_changes` or `when.check_fails`, Ebro will prepend to the script the lines `set -euo pipefail` to ensure sane defaults:

- `-e`: Exit on error
- `-u`: Usage of unset variables is considered an error
- `-o pipefail`: The pipeline’s return status is the value of the last (rightmost) command to exit with a non-zero status, or zero if all commands exit successfully

More on Bash's documentation: [The Set Builtin](https://www.gnu.org/software/bash/manual/bash.html#The-Set-Builtin).

</details>

During the first execution it will execute everything, with no skips, which sould output something like this:

```
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

```
███ [:cache_dir] skipping
███ [:produce_a] skipping
███ [:produce_b] skipping
███ [:echoer] skipping
███ [:producer] satisfied
███ [:default] satisfied
```

Ebro skips tasks whenever possible, and the task definition is what mandates when a task should be skipped. In our example, the task `echoer` is skipped whenever the output of running `cat cache/A.txt` and `cat cache/B.txt` doesn't change. In the case of the task `produce_a`, it skips whenever the command `test -f cache/A.txt` succeeds, because the file `cache/A.txt` already exists.

Now we'll manually edit the file `cache/A.txt`, run `./ebrow` again, and see the result.

```bash
echo 'hello world!' > cache/A.txt
./ebrow
```

```
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

The `Ebro.yaml` format supports importing tasks from other `Ebro.yaml` files by defining the [`imports.*.from` parameter](#the-ebroyaml-format__imports.from). It works like this:

```yaml
# Ebro.yaml
import:
  something:
    from: ./something
```

```yaml
# something/Ebro.yaml
tasks:
  default:
    script: echo 'something'
  else:
    script: echo 'something else'
```

Now, running `ebro something` has this output:

```text
███ [:something:default] running
something
```

And running `ebro something:else` has this output:

```text
███ [:something:else] running
something else
```

Targeting a module by its name is equivalent to calling the module's `default` task. `something` translates to `something:default`.

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

```
███ [:child] running
foo
███ [:default] running
Hello World
```

## CLI

Ebro's command line interface is very straightforward, but has a couple of general rules:

- **Commands** define an specific action. At most, there is one command in a call at most. Absence of a command means the default command of "running". Commands are prefixed with a single hyphen (`-command`).
- **Flags** depend on the command being executed. Their mere presence can mean a boolean value (`true` or `false`) or be accompanyed with a value. Flags are prefixed with two hyphens (`--flag`).
- **Targets** are the names of the tasks that we want to run. When no target is specified, the task `default` is assumed.

To know Ebro's available commands with their flags and explanations, run `ebro -help` (or `ebrow -help` if using the workspace script).

```text
ebro [--flags...] [targets...]
  # Run everything
  flags:
    --file value  Specify the file that should be loaded as root module. default: Ebro.yaml
    --force       Ignore when.* conditionals and dont skip any task. default: false
  targets:
    defaults to [default]


ebro -inventory [--flags...]
  # Display complete inventory of tasks with their definitive configuration
  flags:
    --file value  Specify the file that should be loaded as root module. default: Ebro.yaml


ebro -plan [--flags...] [targets...]
  # Display the execution plan
  flags:
    --file value  Specify the file that should be loaded as root module. default: Ebro.yaml
  targets:
    defaults to [default]


ebro -version
  # Display ebro's version


ebro -help
  # Display this help message
```

## The `Ebro.yaml` format

This is also available as a [JSON Schema](./schema.json).

[Ebro.yaml format explained]

## Versioning

**NOTE**: Ebro is in version 0.x.x, which means that the API isn't stable and could change at any time.

Ebro follows [SemVer](https://semver.org/), this means that it's important to clarify what the API is in Ebro, or, "the exposed parts that will only break compatibility on major version releases".

Ebro's "API" in SemVer terms is composed of, but not limited to:

- The `Ebro.yaml` specification. The main idea is that an `Ebro.yaml` file that worked on Ebro version `1.X.X` **MUST** work exactly the same on any future version under the `1.X.X` major release. This includes, but is not limited to:
  - How `Ebro.yaml` is interpreted.
  - Environment variables added by Ebro during the execution.
  - How scripts are interpreted and executed.
  - How task names are processed and referred to.
  - In which order are tasks executed.
- The CLI **actions**, meaning how the machine state changes given a set of commands and flags.
- The CLI **output** (stdout and stderr) of the `-version`, `-plan` and `-inventory` commands. These commands output an structured representation of data. New data could be added in the future, but existing data will remain the same. This doesn't apply to formatting, as it could change as long as it satisfies the output format specification.

Now, there are some things that are **NOT** considered part of the API and can change at any time, including but not limited to:

- The CLI **output** of any other command apart from the ones explicitly mentioned above.
- Unintended functionality that gets included in Ebro but is considered a bug afterwards. Bugs won't be kept for the sake of API stability and will be removed. If this was the case, a release including a combo of bugfix/breaking change will have proper instructions for updating affected `Ebro.yaml` files.
