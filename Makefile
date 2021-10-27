#Defining variables
# It's necessary to set this because some environments don't link sh -> bash.
SHELL := /usr/bin/env bash

OUT_DIR ?= _output
BIN_DIR := $(OUT_DIR)/bin

.PHONY: all
all: sysadm

.PHONY: sysadm
sysadm:
	$(info Now building sysadm package. sysadm binary file will be placed into "$(BIN_DIR)")
	go build -o $(BIN_DIR)/sysadm 

.PHONY: clean 
clean: 
	$(info Cleaning building cached files.....)
	rm $(BIN_DIR)/*
