## New Features

- Now the task property `script` supports a sequence of strings instead of a single string.

## Breaking Changes

- Now scripts are concatenated instead of replaced during task extension.
- Now unset environment variables are not allowed during environment resolution.

## Miscellanea

- Updated Go `1.23.4` -> `1.23.6`.
