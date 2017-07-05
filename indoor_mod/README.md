# Indoor sensor modules

The indoor sensor modules will be powered by the [Particle
Photon](https://particle.io/) module. Photons are USB powered and have an
integrated ARM microcontroller and wifi module. The device will use the
Particle Pub/Sub API to send new data to the aggregator module.

The indoor sensor modules will use the [AM2302 temperature & humidity
sensor](https://www.adafruit.com/products/393) and use a enclosure.

![Schematic](https://raw.githubusercontent.com/IanLewis/weathersensors/master/indoor_mod/schematic/indoor_mod_bb.png)

## Build the Firmware

Building the firmware by running make in the firmware directory. You can learn more about building the photon firmware by reading the firmware [Getting Started](https://github.com/spark/firmware/blob/develop/docs/gettingstarted.md) page.

```
$ cd firmware
$ make
```

If you run into errors building the firmware sometimes they will resolve themselves by building from a clean version.

```
$ make clean firmware
```

## Flash the Firmware Locally

Flashing requires the `dfu-util` tool.

```
$ make program-dfu
```

## Flash the Firmware Remotely

Flashing the firmware remotely requires the [particle-cli](https://docs.particle.io/guide/tools-and-features/cli/photon/) tool.

TODO: Add instructions.

## Datasheets

Find info about the components used in the [datasheets](datasheets/) directory.
