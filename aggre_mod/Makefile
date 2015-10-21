# Makefile to build aggre_mod and related packages

PROJECT=$(shell gcloud config list | grep project | awk '{ print $3 }')

help:
	@echo "Targets:"
	@echo ""
	@echo "aggre_mod: 		Build the server for the local architecture"
	@echo "docker: 		Build a docker image for the local architecture"
	@echo "aggre_mod_raspi: 	Build the server for the raspberry pi"
	@echo "docker_raspi: 		Build the docker image for the raspberry pi"

# Build the server for the local architecture
aggre_mod:
	CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo .

# Build a docker image for the local architecture
docker: aggre_mod
	docker build -t aggremod .
	docker tag aggremod asia.gcr.io/:`gcloud config list | grep project | awk '{ print $3 }'`/aggremod:`./aggre_mod -version`

push: docker
	gcloud docker push asia.gcr.io/:`gcloud config list | grep project | awk '{ print $3 }'`/aggremod:`./aggre_mod -version`

clean:
	rm -rf aggre_mod
