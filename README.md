# Home Sensors Project

This is my home sensors project.

## Architecture

The system will contain:

1. 2 temperature/humidity sensor modules for inside. One for upstairs and one for downstairs.
1. 1 temperature/humidity/barameter sensor module for outside (maybe also wind & rainfall).

Software will need to be written for each of the types of modules. Each of the
sensor modules will contain a web server which will return the current data
when queried. The aggregator module will polls data for each of the sensors and
streams that data to BigQuery.

![Architecture](https://docs.google.com/drawings/d/1QY_T4k4DTx9b4ChLrcK1LF7cy9I0blKa10raLj2bux0/pub?w=960&amp;h=720)

## Design

### Indoor sensor modules

The [indoor module](indoor_mod/) will read indoor temperature and humidity.

### Outdoor sensor modules

The [outdoor module](outdoor_mod/) will read outdoor temperature, humidity, atmospheric pressure, wind strength, wind direction, and rainfall.

### Aggregation App

The [aggregation server](aggre_mod/) will aggregate data and store it in BigQuery.

# Notes

## Local Development

Flashing devices locally requires the
[particle-cli](https://github.com/spark/particle-cli) and
[dfu-util](http://dfu-util.sourceforge.net/).
