# Schism

Collects data from distributed sensors pushed via the [pushit](https://gitlab.void-ptr.org/go/pushit) service.

Chechout the repository via `git clone && cd schism`.

Config is done via `.envrc`, you need to install [direnv](https://direnv.net) on your dev machine and on the raspberry.

To install dependencies run `go mod download`.

## Production

Build and run via:

```sh
deployer build production
deployer push production
deplyoer up production

# View logs
deployer logs production

# Stop production
deployer down production
```

## Development

Build and run with air (see `./build/air.conf`) for hot-reloading code.

```sh
deployer build
deployer push
deplyoer up

# View logs
deployer logs

# Stop development
deployer down

# Rebuid and update is done via air
```

### Debug

Build and run with delve on port 2345

```sh
deployer build debug
deployer push debug
deplyoer up debug

# View logs
deployer logs debug

# Stop debugging
deployer down

# Rebuid and update
deployer build debug && deployer push debug && deployer up debug
```

In vscode run the debug config `Delve into Docker` (do not forget to set a breakpoint).
