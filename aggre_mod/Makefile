# Makefile to build aggre_mod and related packages

PROJECT=$(shell gcloud config list project | awk 'FNR==2 { print $$3 }')
VERSION=$(shell cat VERSION)

all: aggre_mod

# Build the server for the local architecture
aggre_mod:
	go generate
	CGO_ENABLED=0 GOOS=linux go build -o aggre_mod -a -ldflags '-s' -installsuffix cgo .

# Build a docker image for the local architecture
image: aggre_mod
	docker build -t aggremod .
	$(SHELL) -ec 'docker tag aggremod asia.gcr.io/$(PROJECT)/aggremod:$(VERSION)'

push:
	$(SHELL) -ec 'gcloud docker push asia.gcr.io/$(PROJECT)/aggremod:$(VERSION)'

clean:
	rm -rf aggre_mod
