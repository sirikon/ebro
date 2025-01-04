## Using the `ebrow` workspace script

```bash
curl --location --output ebrow 'https://github.com/sirikon/ebro/releases/latest/download/ebrow'
chmod +x ebrow
```

`ebrow` is a Bash script that contains a reference to an exact Ebro version and is able to download it, verify its integrity, and place it inside the `.ebro` directory, created next to itself. [Read the base script](https://github.com/sirikon/ebro/blob/master/scripts/ebrow).

This is the recommended way of installing and using Ebro, as the `ebrow` file can be committed to a code repository, helping to maintain a consistent environment for all the collaborators in a project.

Each release includes it's own `ebrow` file with the correct version and checksums included. This way, for updating Ebro, all you need to do is replace your current `ebrow` with the one from a new release.

The script depends on:

- Bash >= 4
- `sha256sum`
- `curl`

Any regular Linux distro should satisfy these dependencies already, but **Mac** users will need to install `sha256sum` (with `brew install coreutils`) and a newer Bash (with `brew install bash`) because Mac's stock Bash is `3.2.x`.

## Manually

- Go to the [latest release](https://github.com/sirikon/ebro/releases/latest), or [any release](https://github.com/sirikon/ebro/releases).
- Download the appropriate binary for your operating system and processor architecture (check both using `uname -s` and `uname -m`) and its accompanying `.sha256` file.
- Verify the binary integrity by running `sha256sum -c <binary>.sha256`.
  - Mac users will need to install `sha256sum` with `brew install coreutils`.
- Place the verified binary wherever you want in your system.
