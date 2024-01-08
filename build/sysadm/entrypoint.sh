#!/bin/sh

set -e
if [ "X${BINGETURL}" != "X" ]; then
    /usr/bin/curl -Rfs "${BINGETURL}/sysadm" -o /tmp/sysadm
    if [ $? -eq 0 ]; then
      if [ "/tmp/sysadm" -nt "/opt/sysadm/bin/sysadm" ]; then
          /usr/bin/cp -Rpf /tmp/sysadm /opt/sysadm/bin/sysadm
          /usr/bin/chmod +x /opt/sysadm/bin/sysadm
      fi
    fi

fi

/opt/sysadm/bin/sysadm -c /opt/sysadm/conf/config.yaml daemon start
