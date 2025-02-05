## New features

- Added the environment variable `EBRO_ROOT_FILE` for all tasks. This points to the root `Ebro.yaml` file.
- Added querying capabilities to `requires` and `required_by` fields.

## Bug fixes

- Fixed an issue in which missing the value of a flag would always report `expected value after --file flag` regardless of the actual flag name.

## Miscellanea

- Updated Go dependencies.
