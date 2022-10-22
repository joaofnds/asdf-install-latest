# asdf-install-latest

`asdf-install-latest` installs the latest semver version of all your asdf plugins so you don't have to 😄

https://user-images.githubusercontent.com/9938253/195338904-29dab966-96e7-41e8-b480-305f2d1f8c61.mov

`asdf-install-latest` will:

- ignore desired plugins
- ignore non-semver versions
- set the global version to the latest version after install
- run `~/.config/ail/hooks/{plugin}.sh` after installing a new version of `{plugin}`
- reshim before exiting

# Install

```sh
brew install joaofnds/tap/asdf-install-latest
```

asdf-install-latest is also on aur as [asdf-install-latest-bin](https://aur.archlinux.org/packages/asdf-install-latest-bin)

you can find builds for other systems in the [releases page](https://github.com/joaofnds/asdf-install-latest/releases)

# Default config

```yaml
# ~/.config/ail/config.yml
skip_set_global: false
ignore: []
```
