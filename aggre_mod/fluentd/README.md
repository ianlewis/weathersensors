# aggre\_mod Fluentd Docker image

Build the docker image:

    $ docker build -t aggre_mod-fluentd .

Put the service account key in the conf directory:

    $ cp <keyname>.p12 conf/private_key.p12

Run the docker image like so:

    $ docker run -d --name=aggre_mod-fluentd
      -v `pwd`/conf:/fluentd/etc/ \
      -e GCP_PROJECT=<project> \
      -e GCP_SERVICE_ACCOUNT_EMAIL=<email> \
      -e GCP_BIGQUERY_DATASET=<dataset> \
      -e GCP_BIGQUERY_TABLE=<table> \
      -p 24224:24224 \
      aggre_mod-fluentd

Push the image:

    $ docker tag asia.gcr.io/<project>/aggre_mod-fluentd:v1
    $ gcloud docker push asia.gcr.io/<project>/aggre_mod-fluentd:v1
