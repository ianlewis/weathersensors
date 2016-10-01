# Makefile to build aggre_mod and related packages

PROJECT=$(shell gcloud config list project | awk 'FNR==2 { print $$3 }')
VERSION=$(shell cat VERSION)

all: devicemonitor

# Build the server for the local architecture
devicemonitor:
	go generate
	CGO_ENABLED=0 GOOS=linux go build -o devicemonitor -a -ldflags '-s' -installsuffix cgo .

# Build a docker image for the local architecture
image: devicemonitor
	docker build -t asia.gcr.io/$(PROJECT)/devicemonitor:$(VERSION) .

push:
	$(SHELL) -ec 'gcloud docker -- push asia.gcr.io/$(PROJECT)/devicemonitor:$(VERSION)'

clean:
	rm -rf devicemonitor
