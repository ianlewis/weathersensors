# Makefile for indoor_mod. Copies the files to the right directory and then
# builds the binary.

# Check that given variables are set and all have non-empty values,
# die with an error otherwise.
#
# Params:
#   1. Variable name(s) to test.
#   2. (optional) Error message to print.
check_defined = \
	$(strip $(foreach 1,$1, \
		$(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
	$(if $(value $1),, \
		$(error Undefined $1$(if $2, ($2))))

export PLATFORM=photon

.PHONY: all firmware program-dfu clean src/config.h

all: firmware

firmware: firmware/build/target/user-part/platform-6-m/user-part.bin

firmware/build/target/user-part/platform-6-m/user-part.bin: src/application.cpp src/config.h src/version.h
	rm -rf firmware/user/src/*
	cp -r src/* firmware/user/src/
	# Delete examples or it will try to build them.
	rm -rf firmware/user/src/third_party/PietteTech_DHT/firmware/examples
	rm -rf firmware/user/src/third_party/MQTT-TLS/examples
	$(MAKE) -C firmware/modules all

src/version.h: ../../.git/HEAD ../../.git/index
	echo "#ifndef __VERSION_H__" > $@
	echo "#define __VERSION_H__" >> $@
	echo "const char *gitversion = \"$(shell git rev-parse --short HEAD)\";" >> $@
	echo "#endif" >> $@

src/config.h:
	@:$(call check_defined, MQTT_HOST)
	@:$(call check_defined, MQTT_PORT)
	@:$(call check_defined, MQTT_USERNAME)
	@:$(call check_defined, MQTT_PASSWORD)
	@:$(call check_defined, LOCATION)
	echo "#ifndef __CONFIG_H__" > $@
	echo "#define __CONFIG_H__" >> $@
	echo "#define CONFIG_MQTT_HOST \"$(MQTT_HOST)\"" >> $@
	echo "#define CONFIG_MQTT_PORT $(MQTT_PORT)" >> $@
	echo "#define CONFIG_MQTT_USERNAME \"$(MQTT_USERNAME)\"" >> $@
	echo "#define CONFIG_MQTT_PASSWORD \"$(MQTT_PASSWORD)\"" >> $@
	echo "#define CONFIG_LOCATION \"$(LOCATION)\"" >> $@
	echo "#endif" >> $@

program-dfu: firmware
	$(MAKE) -C firmware/main program-dfu

flash: firmware
	@:$(call check_defined, DEVICE)
	particle flash $(DEVICE) firmware/build/target/user-part/platform-6-m/user-part.bin

clean:
	rm -f src/version.h
	rm -f src/config.h
	$(MAKE) -C firmware/modules clean
