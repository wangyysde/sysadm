#!/bin/bash

SYSADM_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)"
SYSADM_BIN=${SYSADM_ROOT}/_output/bin/sysadm 
REGISTRYCTL_BIN=${SYSADM_ROOT}/_output/bin/registryctl
REGISTRY_BIN=${SYSADM_ROOT}/_output/bin/registry
INFRASTRUCTURE_BIN=${SYSADM_ROOT}/_output/bin/infrastructure
AGENT_BIN=${SYSADM_ROOT}/_output/bin/agent
SYSADM_CONF=${SYSADM_ROOT}/_output/conf/config.yaml
REGISTRYCTL_CONF=${SYSADM_ROOT}/_output/conf/registryctl.yaml
REGISTRY_CONF=${SYSADM_ROOT}/_output/conf/registry.yml
INFRASTRUCTURE_CONF=${SYSADM_ROOT}/_output/conf/infrastructure.yaml
AGENT_CONF=${SYSADM_ROOT}/_output/conf/agent.yaml

LOG_DIR="/var/log/sysadm/"
nohup ${REGISTRY_BIN} serve ${REGISTRY_CONF} 2>&1 >>"${LOG_DIR}registry.log" &
nohup  ${REGISTRYCTL_BIN} daemon start -c ${REGISTRYCTL_CONF} 2>&1 >>"${LOG_DIR}registryctl.log" &
nohup ${INFRASTRUCTURE_BIN} start -c ${INFRASTRUCTURE_CONF} 2>&1 >>"${LOG_DIR}infrastructure.log" &
nohup ${SYSADM_BIN}  daemon start -c ${SYSADM_CONF} 2>&1 >>"${LOG_DIR}sysadm.log" &


