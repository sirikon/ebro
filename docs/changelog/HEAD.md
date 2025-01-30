## Breaking Changes

- Now the `interactive` and `labels` properties introduced in version `0.10.0` are included in `extends` operations. 

### Variable expansion updates
Expanding environment variables now supports referencing variables that were defined **before** the current variable in the same map.

```yaml
environment:
  VERSION: "1.0.0"

tasks:
  default:
    environment:
      VERSION: "2.0.0"
      VERSION_RC: "${VERSION}-rc"
    script: echo "$VERSION_RC"
```

Running the `Ebro.yaml` file above would result in:

- before: `1.0.0-rc`
- now: `2.0.0-rc`

## Miscellanea

- Change on the `ebrow` workspace script: Now using `shasum` instead of `sha256sum` so the built-in command behaves the same in both Linux and MacOs.
