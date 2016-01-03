# Local Fluentd

This directory contains config for running a local fluentd that outputs to
stdout. It is inteded to be used for local debugging.

# Running Fluentd

It's easiest to run fluentd in a container. Run it in the current directory
with the following command using Docker.

    $ docker run -p 24224:24224 -v `pwd`:/fluentd/etc fluent/fluentd
