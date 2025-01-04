## New features

- New command `-list` for listing just the names of the tasks in the inventory.
- Short version for commands have been added: `-i`, `-l`, `-p`, `-v` and `-h`.

## Breaking changes

- Unknown properties in `Ebro.yaml` files are not allowed anymore.
- When the target in an Ebro execution is prefixed with `:`, another `:` will not be added.
  - Before, calling `ebro :foo` would try to run the task `::foo`
  - Now, calling `ebro :foo` is equivalent to calling `ebro foo`.
