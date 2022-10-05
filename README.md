# asdf-install-latest
`asdf-install-latest` installs the latest semver version of all your asdf plugins so you don't have to 😄

https://user-images.githubusercontent.com/9938253/194134283-09df3e14-4994-424a-9aef-3d831edc44fe.mov

`asdf-install-latest` will:
- ignore all plugins found in `~/.config/ail/ignore`
- ignore non-semver versions
- set the global version to the latest version after install
- run (if present) `~/.config/ail/hooks/{plugin}.sh` after installing a new version of `{plugin}`
- reshim before exiting

# Install
```sh
brew install joaofnds/tap/asdf-install-latest
```
[we also provide binaries for other systems](https://github.com/joaofnds/asdf-install-latest/releases)
