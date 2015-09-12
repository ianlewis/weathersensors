# Makefile for indoor_mod. Copies the files to the right directory and then
# builds the binary.

export PLATFORM=photon

all: firmware

firmware: modules
	$(MAKE) -C firmware/main all

modules:
	rm -f firmware/user/src/*
	cp -r src/* firmware/user/src/
	$(MAKE) -C firmware/modules all

program-dfu: firmware
	$(MAKE) -C firmware/main program-dfu

clean:
	$(MAKE) -C firmware/modules clean
	$(MAKE) -C firmware/main clean
