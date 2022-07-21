FROM centos:centos7.9.2009 

RUN mkdir -p /opt/infrastructure/{bin,conf,logs}
COPY ./install_cert.sh /home/sysadm
COPY ./entrypoint.sh /home/sysadm
COPY ./registry /usr/bin/registry_DO_NOT_USE_GC
COPY ./config.yml /etc/registry/

RUN chown -R sysadm:sysadm /etc/pki/tls/certs \
    && chown sysadm:sysadm /home/sysadm/entrypoint.sh && chmod u+x /home/sysadm/entrypoint.sh \
    && chown sysadm:sysadm /home/sysadm/install_cert.sh && chmod u+x /home/sysadm/install_cert.sh \
    && chown sysadm:sysadm /usr/bin/registry_DO_NOT_USE_GC && chmod u+x /usr/bin/registry_DO_NOT_USE_GC

HEALTHCHECK CMD curl --fail -s http://localhost:5000 || curl -k --fail -s https://localhost:5443 || exit 1

USER sysadm

ENTRYPOINT ["/home/sysadm/entrypoint.sh"]

VOLUME ["/storage"]
