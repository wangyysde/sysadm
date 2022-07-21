#!/bin/sh

set -e

/opt/registryctl/bin/registryctl -c /opt/registryctl/conf/registryctl.yaml daemon start
