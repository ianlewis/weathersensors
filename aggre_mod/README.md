# Aggre Mod

This module is a data aggregator. It receives data published from devices and
forwards it on to Fluentd/BigQuery.

# Architecture

![Architecture](https://docs.google.com/drawings/d/1QY_T4k4DTx9b4ChLrcK1LF7cy9I0blKa10raLj2bux0/pub?w=960&amp;h=720)

# Deploying

aggre\_mod is deployed to [Google Container
Engine](https://cloud.google.com/container-engine/) where it runs in a
[pod](http://kubernetes.io/v1.0/docs/user-guide/pods.html) with Fluentd.

The app is deployed via the following steps. You will need to have
make, Docker, and the Google Cloud SDK installed:

1. Create a cluster on Container Engine
1. Create and push the container images:

       $ make push

1. Create the production namespace:

       $ kubectl create -f homesensors-prod-ns.yaml

1. Create the aggre\_mod [Replication Controller](http://kubernetes.io/v1.0/docs/user-guide/replication-controller.html):

       $ kubectl create -f aggremod-rc.yaml --namespace=homesensors-prod
