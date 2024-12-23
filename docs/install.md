## Using the `ebrow` workspace script

```bash
curl --locaton --output ebrow 'https://github.com/sirikon/ebro/releases/latest/download/ebrow'
chmod +x ebrow
```

`ebrow` is a Bash script that contains a reference to an exact Ebro version and is able to download it, verify its integrity, and place it inside the `.ebro` directory, created next to itself. 

This is the recommended way of installing Ebro, as the `ebrow` file can be committed to a code repository, helping to maintain a consistent environment for all the collaborators in a project.

## Manually

- Go to the [latest release](https://github.com/sirikon/ebro/releases/latest), or [any release](https://github.com/sirikon/ebro/releases).
- Download the binary appropriate for your operating system and processor architecture (check both using `uname -s` and `uname -m`) and its accompanying `.sha256` file.
- Verify the binary integrity by running `sha256sum -c <binary>.sha256`.
- Place the verified binary wherever you want in your system.
