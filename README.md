# Home Sensors Project

This is my home sensors project.

## Architecture

The system will contain:

1. 2 temperature/humidity sensor modules for inside. One for upstairs and one for downstairs.
1. 1 temperature/humidity/barameter sensor module for outside.
1. 1 Raspberry Pi for aggregating data.

Software will need to be written for each of the types of modules. Each of the
sensor modules will contain a web server which will return the current data when
queried. The Raspberry Pi module will have a client program which polls each
of the sensors and streams that data to BigQuery.

## Design

# Indoor sensor modules

The indoor sensor modules will be powered by the [Particle
Photon](https://store.particle.io/?product=particle-photon) module. Photons are
USB powered and have an integrated ARM microcontroller and wifi module. This will
make it easy to connect implement a web server that can be polled by the Raspberry Pi.

The indoor sensor modules will use the [AM2302 temperature & humidity
sensor](https://www.adafruit.com/products/393) and use a [small
enclosure](http://www.alibaba.com/product-detail/Hot-selling-small-electric-plastic-enclosure_60269810523.html?spm=a2700.7724857.35.1.5RfGUV).

## Web Server API

TODO

# Outdoor sensor modules

The indoor sensor modules will also be powered by the [Particle
Photon](https://store.particle.io/?product=particle-photon) module. Photons are
USB powered and have an integrated ARM microcontroller and wifi module. This will
make it easy to connect implement a web server that can be polled by the Raspberry Pi.

The indoor sensor modules will use the [AM2315 - Encased I2C
Temperature/Humidity Sensor](https://www.adafruit.com/products/1293) and use
the same [small
enclosure](http://www.alibaba.com/product-detail/Hot-selling-small-electric-plastic-enclosure_60269810523.html?spm=a2700.7724857.35.1.5RfGUV)
as the indoor modules.

The outdoor module will also contain a barameter [BMP180 Barometric Pressure/Temperature/Altitude Sensor]()

## Web Server API

TODO

# Raspberry Pi

The Raspberry Pi will run an aggregation program written in Go. The app will
poll each of the sensors and forward the data to BigQuery. One way to forward
the data would be to use fluentd and the BigQuery plugin to forward the data
but it may just be easier to have the Go program do that.
