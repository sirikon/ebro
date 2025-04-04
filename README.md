# Ebro

[📚 Documentation](https://ebro.sirikon.me)

Ebro is a task runner. Tasks are defined inside YAML files, scripted with Bash, and configured with a name, requirements, when to skip execution, and other details.

Ebro is distributed as a single binary, including the script interpreter ([mvdan/sh](https://github.com/mvdan/sh)).

It's heavily inspired in [go-task/task](https://github.com/go-task/task), but originally built around a personal need for configuring servers, although it's not tied to this use case and remains agnostic.

## Install the `ebrow` workspace script

```bash
curl --location --output ebrow 'https://github.com/sirikon/ebro/releases/latest/download/ebrow'
chmod +x ebrow
```

[More information](http://ebro.sirikon.me/install.html).
