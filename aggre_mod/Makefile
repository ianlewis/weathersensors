# Makefile to build aggre_mod and related packages

PROJECT=$(shell gcloud config list project | awk 'FNR==2 { print $$3 }')

help:
	@echo "Targets:"
	@echo ""
	@echo "aggre_mod: 		Build the server for the local architecture"
	@echo "docker: 		Build a docker image for the local architecture"
	@echo "push: 			Push the docker image to the repository"

# Build the server for the local architecture
aggre_mod:
	CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo .

# Build a docker image for the local architecture
docker: aggre_mod
	docker build -t aggremod .
	$(SHELL) -ec 'docker tag aggremod asia.gcr.io/$(PROJECT)/aggremod:`./aggre_mod -version`'

push:
	$(SHELL) -ec 'gcloud docker push asia.gcr.io/$(PROJECT)/aggremod:`./aggre_mod -version`'

clean:
	rm -rf aggre_mod
