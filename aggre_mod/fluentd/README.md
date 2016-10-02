# aggre\_mod Fluentd Docker image

This directory holds the build artifacts for the Fluentd that runs with aggre\_mod.

## Build the Docker Image

Build the docker image:

```shell
$ docker build -t gcr.io/<project-id>/aggremod-fluentd:<version> .
```

or alternatively run make:

```shell
$ make
```

## Push the Docker Image

Push the Docker image to GCR:

```shell
$ gcloud docker -- push gcr.io/<project-id>/aggremod-fluentd:<version>
```

or alternatively with make:

```shell
$ make push
```

## Run the container

The container is meant to run as a sidecar with the aggre\_mod application. It requires that there be a [Kubernetes secrets](http://kubernetes.io/docs/user-guide/secrets/) with the path specified in the config.

The following environment variables are used in fluentd's config.

**GCP\_PROJECT**: The Google Cloud project id
**GCP\_SERVICE\_ACCOUNT\_KEY\_PATH**: The path to the service account key JSON file.
**GCP\_BIGQUERY\_DATASET**: The BigQuery dataset to write to.
**GCP\_BIGQUERY\_TABLE**: The BigQuery table to write to.

See [deploy.yaml](../deploy.yaml) for details about how it's used.
