# Device Monitor

This directory contains a device monitor that monitors [Particle](https://particle.io/) devices and creates error reports in [Stackdriver Error Reporting](https://cloud.google.com/error-reporting/) if they are offline for too long. The device monitor monitors the online status of each device or a subset of devices by polling the [Particle Devices API](https://docs.particle.io/reference/api/#devices).

## Usage

### Command line arguments

```shell
$ ./devicemonitor -help
Usage of ./devicemonitor:
  -access-token string
    	The path to a file containing the Particle API access token.
  -device-list string
    	The path to a text file of device IDs (one per line) to monitor. If not specified, all devices are monitored.
  -device-timeout int
    	The time that a device can be offline before an error is produced. (default 300)
  -host string
    	The web server address for health checks. (default ":8080")
  -poll-interval int
    	API polling interval in seconds. (default 30)
  -project string
    	The Google Cloud Platform project ID for the Error Reporting API.
  -version
    	Print the version and exit.
```

### Environment Variables

Some environment variables can be set to configure the monitor.

- **ADDRESS**: The address of the health check web server to bind to. This is overridden by the `-host` command line argument.
- **DEVICE_LIST_PATH**: The path to a text file of device IDs (one per line) to monitor. If not specified, all devices are monitored. This is overridden by the `-device-list` command line argument.
- **DEVICE_TIMEOUT**: The time in seconds that a device can be offline before an error is produced. This is overridden by the `-device-timeout` command line argument.
- **ACCESS_TOKEN_PATH**: The path to a file containing the Particle API access token. This is overridden by the `-access-token` command line argument.
- **POLL_INTERVAL**: API polling interval in seconds. This is overridden by the `-poll-interval` command line argument.
- **GCP_PROJECT**: The Google Cloud Platform project ID for the Error Reporting API. This is overridden by the `-project` command line argument.
- **GOOGLE_APPLICATION_CREDENTIALS**: The path to the service account JSON file. **Required.**

## Setup

This section will take you through how to set up the required files for the device monitor.

### Create the Service Account

Make sure you have the [Google Cloud Platform SDK](https://cloud.google.com/sdk/) installed and authenticated.

[Create a service account](https://cloud.google.com/iam/docs/creating-managing-service-accounts) for the device monitor.

```shell
$ gcloud iam service-accounts create devicemonitor --display-name "devicemonitor"
```

Create a service account key. This will output a `service-account.json` file.

```shell
$ gcloud iam service-accounts keys create service-account.json --iam-account devicemonitor@my-project.iam.gserviceaccount.com
```

### Create the Token file

Copy your [Particle access token](https://docs.particle.io/reference/api/#authentication) to the `token` file.

```shell
$ echo 9876987698769876987698769876987698769876 > token
```

### Create a Device List (Optional)

Create a list of device you want to monitor. This is optional. If you don't add a device list then all devices returned by the API will be monitored. The path to this file can be passed to the `DEVICE_LIST_PATH` environment variable or to the `-device-list` command line argument.

```shell
$ echo 0123456789abc > devices
$ echo 53ff6f0650723 >> devices
$ echo 53ff291839887 >> devices
```

### Build and Run the App

Building the app will require a relatively recent Go toolchain.

```shell
$ go generate
$ go build .
```

Alternatively, you can use make.

```shell
$ make clean devicemonitor
```

You can run the app from anywhere given yout pass in the right command line arguments.

```shell
$ GOOGLE_APPLICATION_CREDENTIALS=service-account.json ./devicemonitor -project-id=my-project -access-token=token
```

## Deploy on Kubernetes

There are Kubernetes deployment files included in the repository.

### Create the Secrets

Create the necessary secret.

```shell
$ kubectl create secret generic devicemonitor-secret --from-file=service-account.json --from-file=token
secret "devicemonitor-secret" created
```

### Create the ConfigMap

Create the ConfigMap by giving it a path to the `devices` file and passing it your Google Cloud Platform project ID.

```shell
$ kubectl create configmap devicemonitor-conf --from-file=devices --from-literal=project-id=my-project
configmap "devicemonitor-conf" created
```

### Deploy the Device Monitor

Create the deployment in Kubernetes.

```shell
$ kubectl create -f deploy.yaml
deployment "devicemonitor" created.
```
