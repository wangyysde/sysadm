#Defining variables
# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /usr/bin/env bash
OUT_DIR ?= _output
BIN_DIR := $(OUT_DIR)/bin
PRINT_HELP ?=
PREFIX ?= /usr/local/sysadm
REGISTRYvER ?= v2.7.0
BUILD_IMAGE ?= 
IMAGEVER ?= v1.4

.PHONY: all
ifeq ($(PRINT_HELP),y)
all:
	build/help_info.sh all
else
all: 
	build/build.sh "$(WHAT)" "$(BUILD_IMAGE)" "$(IMAGEVER)"
endif


.PHONY: sysadm
sysadm:
	$(info Now building sysadm package. sysadm binary file will be placed into "$(BIN_DIR)")
	build/build.sh "sysadm" "$(BUILD_IMAGE)" "$(IMAGEVER)"
#	go build -o $(BIN_DIR)/sysadm 
	
.PHONY: registryctl
registryctl:
	$(info Now building registryctl package. registryctl binary file will be placed into "$(BIN_DIR)")
	build/build.sh "registryctl" "$(BUILD_IMAGE)" "$(IMAGEVER)"

.PHONY: agent
agent:
	$(info Now building agent package. agent binary file will be placed into "$(BIN_DIR)")
	build/build.sh "agent" "$(BUILD_IMAGE)" "$(IMAGEVER)"

.PHONY: registry
registry: 
	$(info Now building registry binary package. registry binary package file will be placed into "$(BIN_DIR)")
	build/build_registry_binary.sh $(REGISTRYvER)
ifeq ($(BUILD_IMAGE),y)
	build/build_registry_image.sh $(IMAGEVER)
endif

.PHONY: infrastructure
infrastructure:
	$(info Now building infrastructure package. infrastructure binary file will be placed into "$(BIN_DIR)")
	build/build_infrastructure_image.sh "infrastructure" "$(BUILD_IMAGE)" "$(IMAGEVER)"

.PHONY: install 
install: 
	test -d '$(PREFIX)/bin' || mkdir -p '$(PREFIX)/bin'
	install '$(BIN_DIR)/sysadm' '$(PREFIX)/bin/'
	test -d '$(PREFIX)/conf' || mkdir -p '$(PREFIX)/conf'
	install '$(OUT_DIR)/conf/config.yaml' '$(PREFIX)/conf/'
	test -d '$(PREFIX)/logs' || mkdir -p '$(PREFIX)/logs'
	test -d '$(PREFIX)/formstmpl' || mkdir -p '$(PREFIX)/formstmpl'
	install '$(OUT_DIR)/formstmpl/*' '$(PREFIX)/formstmpl'
	test -d '$(PREFIX)/html' || mkdir -p '$(PREFIX)/html'
	install '$(OUT_DIR)/html/*' '$(PREFIX)/html'

.PHONY: clean 
clean: 
	$(info Cleaning building cached files.....)
	rm $(BIN_DIR)/*
