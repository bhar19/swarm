FROM ubuntu:15.10
COPY ./swarm_bhar /usr/bin
#RUN apt-get update && \
#    apt-get install -y golang && \
#    apt-get clean

RUN apt-get install -y gccgo-5 && \
    update-alternatives --set go /usr/bin/go-5

#COPY /usr/lib/x86_64-linux-gnu/libgo.so.7 /usr/lib/x86_64-linux-gnu/libgo.so.7

#CMD ["ls /usr/bin"]
#CMD ls /usr/bin/go
#CMD which ldd
#CMD ldd /usr/bin/swarm_bhar && \
#    ls -ld /usr/bin/go && \
#    ls -ld /usr/lib/x86_64-linux-gnu/libgo.so.7.0.0
#CMD ls -ld /usr/bin/go
#CMD /usr/bin/swarm_bhar $@
ENTRYPOINT ["swarm_bhar"]
