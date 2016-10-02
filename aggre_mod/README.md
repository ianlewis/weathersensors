# Aggregator Module

This module is a data aggregator. It receives data published from devices and
forwards it on to Fluentd/BigQuery.

# Architecture

![Architecture](https://docs.google.com/drawings/d/1QY_T4k4DTx9b4ChLrcK1LF7cy9I0blKa10raLj2bux0/pub?w=960&amp;h=720)

# Setup

aggre\_mod is deployed to [Google Container
Engine](https://cloud.google.com/container-engine/) where it runs in a
[pod](http://kubernetes.io/v1.0/docs/user-guide/pods.html) with Fluentd.

The app is deployed via the following steps. You will need to have
make, Docker, and the Google Cloud SDK installed:

1. Create a BigQuery dataset.

        bq mk --description "Sensors Dataset" weathersensors

1. Create a BigQuery table.

        bq mk --description "Sensor data" weathersensors.sensordata schema.json

1. Create and push the container images (make sure you have the Google Cloud SDK installed and configured for your project):

        make clean image push

1. [Create a service account](https://cloud.google.com/iam/docs/creating-managing-service-accounts) and get the JSON key. Save it to the file 'service-account.json'.
1. Create the configmap:

        kubectl create configmap aggremod-conf \
            --from-literal=project-id=$(gcloud config list project | awk 'FNR==2 { print $3 }') \
            --from-literal=bigquery-dataset=weathersensors \
            --from-literal=bigquery-table=sensordata

1. Create the secret:

        kubectl create secret generic aggremod-secret \
            --from-literal=token=[particle.io token] \
            --from-file=service-account.json

1. Create the aggre\_mod Deployment:

        kubectl create -f deploy.yaml
