FROM scratch

COPY ./swarm_io /swarm_io
COPY ./certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY ./tmp /tmp

ENV SWARM_HOST :2375
EXPOSE 2375

VOLUME /.swarm_io

ENTRYPOINT ["/swarm_io"]
CMD ["--help"]
