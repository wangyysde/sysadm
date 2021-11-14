#Defining variables
# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /usr/bin/env bash
OUT_DIR ?= _output/bin
BIN_DIR := $(OUT_DIR)/bin
PRINT_HELP ?=

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

.PHONY: clean 
clean: 
	$(info Cleaning building cached files.....)
	rm $(BIN_DIR)/*
	