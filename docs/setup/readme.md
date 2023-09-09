# Setup

Here you can find some docs on how to setup certain things.

## WIFI

To setup wifi on a raspberry, you can use the Raspi Imager to set it all up via UI: <https://www.raspberrypi.org/software/>.

You can also do it manually, see <https://www.raspberrypi.com/documentation/computers/configuration.html#setting-up-a-headless-raspberry-pi>.

TL;DR:

```sh
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
