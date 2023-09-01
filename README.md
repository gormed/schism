# Schism

Collects data from distributed sensors pushed via the [pushit](https://gitlab.void-ptr.org/go/pushit) service.

Checkout the repository via `git clone && cd schism`.

Config is done via `.envrc`, you need to install [direnv](https://direnv.net) on your dev machine and on the raspberry.

To install dependencies run `go mod download`.

## Setup

### Raspberry

Quick setup on a raspberry pi 4 running raspbian OS lite 64bit:

```sh
# Setup zsh and stuff we want for ease of use
sudo apt update
sudo apt install zsh git ca-certificates curl gnupg figlet direnv
# Install antigen, fzf and a custom .zshrc
curl -L git.io/antigen > ~/antigen.zsh
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf
~/.fzf/install
# Install .zshrc
curl -L "https://gist.githubusercontent.com/gormed/c339a3448c5530471586bc238d44b106/raw/52a11e3309a4c0d86c1589bbe2de35b9a6513d27/.zshrc" > ~/.zshrc
# Change shell to zsh
chsh -s $(which zsh)
# Now we can use zsh
# by just `zsh` or logging out and in again
```

#### Docker

```sh
# Add Dockers GPG key
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
```

```sh
# Use the following command to set up the repository:
echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
```

```sh
# Install Docker Engine
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin
```

See <https://docs.docker.com/engine/install/debian/> for more.

See <https://docs.docker.com/engine/install/raspberry-pi-os/> for the 32bit installation.

Promote the current user to the docker group:

```sh
sudo groupadd docker
sudo usermod -aG docker $USER
# Log out and log back in so that your group membership is re-evaluated
# Verify that you can run docker commands without sudo
docker run hello-world
```

---

See <https://docs.docker.com/engine/install/linux-postinstall/> for more.

## Production

Build and run via:

Images are build via [github](https://github.com/gormed/schism) and [docker.com](https://hub.docker.com/r/gormed/schism) - pull them via `docker pull gormed/schism:latest` to get the latest master build.

```sh
deployer pull production
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
