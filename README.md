# Home Sensors Project

This is my home sensors project.

## Architecture

The system will contain:

1. 2 temperature/humidity sensor modules for inside. One for upstairs and one for downstairs.
1. 1 temperature/humidity/barameter sensor module for outside (maybe also wind & rainfall).

Software will need to be written for each of the types of modules. Each of the
sensor modules will contain a web server which will return the current data when
queried. The Raspberry Pi module will have a client program which polls each
of the sensors and streams that data to BigQuery.

![Architecture](https://docs.google.com/drawings/d/1QY_T4k4DTx9b4ChLrcK1LF7cy9I0blKa10raLj2bux0/pub?w=960&amp;h=720)

## Design

# Indoor sensor modules

The indoor sensor modules will be powered by the [Particle
Photon](https://store.particle.io/?product=particle-photon) module. Photons are
USB powered and have an integrated ARM microcontroller and wifi module. This will
make it easy to connect implement a web server that can be polled by the Raspberry Pi.

The indoor sensor modules will use the [AM2302 temperature & humidity
sensor](https://www.adafruit.com/products/393) and use a [small
enclosure](http://www.alibaba.com/product-detail/Hot-selling-small-electric-plastic-enclosure_60269810523.html?spm=a2700.7724857.35.1.5RfGUV).

![Schematic](https://googledrive.com/host/0ByyQM-rDUkX6VHl1alZNT3I2aUU/indoor_mod_bb.png)

# Outdoor sensor modules

The indoor sensor modules will also be powered by the [Particle
Photon](https://store.particle.io/?product=particle-photon) module. Photons are
USB powered and have an integrated ARM microcontroller and wifi module. This will
make it easy to connect implement a web server that can be polled by the Raspberry Pi.

The indoor sensor modules will use the [AM2315 - Encased I2C
Temperature/Humidity Sensor](https://www.adafruit.com/products/1293). The
outdoor module will also contain the [BMP180 Barometric
Pressure/Temperature/Altitude Sensor](https://www.adafruit.com/products/1603)

## Aggregation App

The [aggregation server](aggre_mod/) will aggregate data and store it in BigQuery.

# Notes

## Local Development

Flashing devices locally requires the
[particle-cli](https://github.com/spark/particle-cli) and
[dfu-util](http://dfu-util.sourceforge.net/).
