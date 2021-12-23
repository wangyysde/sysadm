#Defining variables
# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /usr/bin/env bash
OUT_DIR ?= _output
BIN_DIR := $(OUT_DIR)/bin
PRINT_HELP ?=
PREFIX ?= /usr/local/sysadm

.PHONY: all
ifeq ($(PRINT_HELP),y)
all:
	build/help_info.sh all
else
all: 
	build/build.sh $(WHAT)
endif


.PHONY: sysadm
sysadm:
	$(info Now building sysadm package. sysadm binary file will be placed into "$(BIN_DIR)")
	build/build.sh $(WHAT)
#	go build -o $(BIN_DIR)/sysadm 

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
