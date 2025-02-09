## New Features

- Now the task properties `script`, `when.output_changes` and `when.check_fails` support a sequence of strings instead of a single string.

## Breaking Changes

- Now scripts (`script`, `when.output_changes` and `when.check_fails`) are concatenated instead of replaced during task extension.
- Now unset environment variables are not allowed during environment resolution.

## Miscellanea

- Updated Go `1.23.4` -> `1.23.6`.
- Updated Go dependencies.
