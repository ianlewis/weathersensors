# Makefile to build aggre_mod and related packages

# Build the server for the local architecture
aggre_mod:
	CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo .

# Build a docker image for the local architecture
docker: aggre_mod
	docker build -t aggre_mod .

# Build the docker image for the raspberry pi
docker_raspi: aggre_mod_arm
	docker build -t aggre_mod_raspi .

# Build the server for the raspberry pi
aggre_mod_arm:
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build -a -ldflags '-s' -installsuffix cgo .

clean:
	rm -rf aggre_mod
