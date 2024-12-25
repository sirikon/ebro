<div markdown="1" remove-in-website="1">

# Ebro Documentation

</div>

Ebro is a task runner. Tasks are defined inside YAML files, scripted with Bash, and configured with a name, requirements, when to skip execution, and other details.

Ebro is distributed as a single binary, including the script interpreter ([mvdan/sh](https://github.com/mvdan/sh)).

It's heavily inspired in [go-task/task](https://github.com/go-task/task), but originally built around a personal need for configuring servers, although it's not tied to this use case and remains agnostic.

The format of Ebro.yaml files is defined [here](./ebro-format.md). Here is an example:

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
