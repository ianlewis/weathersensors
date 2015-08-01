# Makefile for indoor_mod. Copies the files to the right directory and then
# builds the binary.

export PLATFORM=photon

all:
	rm -f firmware/user/src/*
	cp -r src/* firmware/user/src/
	$(MAKE) -C firmware/main

program-dfu: all
	$(MAKE) -C firmware/main program-dfu

clean:
	$(MAKE) -C firmware/main clean
