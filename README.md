# ebro

> [!WARNING]
> Work in progress. Undocumented.

ebro is a tool for defining tasks with their dependencies and executing them in the correct order.

It's configured using Yaml files (sorry) and the tasks are shell scripts interpreted with [mvdan/sh](https://github.com/mvdan/sh).

```yaml
tasks:
  default:
    requires: [echoer, producer]

  echoer:
    script: |
      cat cache/A.txt
      cat cache/B.txt
    when:
      output_changes: cat cache/A.txt && cache/B.txt

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

It's heavily inspired in [go-task/task](https://github.com/go-task/task), but built around a personal need for configuring servers, although it's not tied to this use case.
