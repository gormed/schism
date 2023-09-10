# Device Setup

## Setup SD Card

Install Raspberry PI Imager: <https://www.raspberrypi.org/software/>

With the imager, create a SD card with the following options:

- Choose OS: Raspberry Pi OS (32-bit)
- Choose SD card: Choose the SD card you want to use
- Configure WIFI (see [wifi config](../readme.md#wifi))
- Configure SSH (see [ssh config](../readme.md#ssh))
- Always set username and password to `pi` and `raspberry` (we can change that later), not setting it will cause problems with the setup scripts (see <https://github.com/raspberrypi/rpi-imager/issues/540>)
- Configure your SSH key

## Terminal and co

See [terminal setup](../readme.md#terminal-and-co).

## Install golang

```sh
# See https://github.com/canha/golang-tools-install-script#fast_forward-install
# Add go to path
export GOPATH=$HOME/go
export GOROOT=/usr/local/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
# Install go
wget -q -O - https://git.io/vQhTU | bash
```

## Generate SSH Key

See [ssh key docs](../readme.md#generate-ssh-key).
