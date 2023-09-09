# Setup

Here you can find some docs on how to setup certain basic things for Raspberry Pi's.

## WIFI

To setup wifi on a raspberry, you can use the Raspi Imager to set it all up via UI: <https://www.raspberrypi.org/software/>.

You can also do it manually, see <https://www.raspberrypi.com/documentation/computers/configuration.html#setting-up-a-headless-raspberry-pi>.

TL;DR:

```conf
# Create a file called `wpa_supplicant.conf` in the boot partition of the SD card with the following content:
ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1
country=DE

network={
  ssid="SSID"
  psk="PASSWORD"
  priority=1
  id_str="MY_ID"
}
```

## SSH

To setup ssh on a raspberry, you can use the Raspi Imager to set it all up via UI: <https://www.raspberrypi.org/software/>.

To do it manually, create a file called `ssh` in the boot partition of the SD card.

## Terminal and co

This setup should work for raspian lite 32bit and 64bit.

```sh
# Setup zsh and stuff we want for ease of use
sudo apt update
sudo apt install zsh git ca-certificates curl gnupg figlet direnv
# Install antigen, fzf and a custom .zshrc
curl -L git.io/antigen > ~/antigen.zsh
git clone --depth 1 https://github.com/junegunn/fzf.git ~/.fzf
~/.fzf/install
# Install .zshrc for raspberry 4 and with docker
curl -L "https://gist.githubusercontent.com/gormed/c339a3448c5530471586bc238d44b106/raw/52a11e3309a4c0d86c1589bbe2de35b9a6513d27/.zshrc" > ~/.zshrc
# Install .zshrc for raspberry zero
curl -L "https://gist.githubusercontent.com/gormed/323eb5bd288b4c6129c881f02aa9b85d/raw/0481d37b1518fb837dc252667998c827f3dcea20/.zshrc" > ~/.zshrc
# Change shell to zsh
chsh -s $(which zsh)
# Now we can use zsh
# by just `zsh` or logging out and in again
```
