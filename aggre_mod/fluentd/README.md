# aggre\_mod Fluentd Docker image

# Build the Docker Image.

Build the docker image:

    $ docker build -t aggremod-fluentd .

Tag it and push it to Google Cloud Registry:

    $ docker tag aggremod-fluentd gcr.io/<project-id>/aggremod-fluentd:<version>

# Run the container

The container is run as a sidecar with the aggre\_mod application. It requires
that there be a [Kubernetes
secrets](http://kubernetes.io/v1.0/docs/user-guide/secrets.html) mounted in the
"/secrets" folder with the private key (p12) for the Google service account
stored in a file called "private-key".

See [aggremod-rc.yaml](../aggremod-rc.yaml) for details.
