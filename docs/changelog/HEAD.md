## Breaking changes

### `required_by` behavior

Referencing a task in `required_by` doesn't add the referenced task to the plan anymore. The referenced task will need to be referenced in a `requires` or directly called as a target. Here's an example:

```yaml
tasks:
  default:
    requires: [b]

  a:
    script: echo A

  b:
    required_by: [a]
    script: echo B
```

Before this release:

```yaml
███ [:b] running
B
███ [:a] running
A
███ [:default] satisfied
```

After this release:

```yaml
███ [:b] running
B
███ [:default] satisfied
```

### Tasks with nothing to do are invalid

From now on, tasks that have nothing to do (no `script` nor `requires`) are considered invalid and Ebro's execution will error during the inventory process.

## Miscellanea

Now, the `.sha256` files produced during the release process are valid files to be checked by `sha256sum` by running `sha256sum --check <binary>.sha256`, easing up the check process. Before this release it just contained the SHA256 checksum of the file.