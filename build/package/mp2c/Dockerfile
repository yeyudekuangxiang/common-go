FROM yeyudekuangxiang/alpine:3.17-tzsh

WORKDIR  /app
COPY mp2c .
COPY ../../static ./static

RUN chmod +x mp2c

EXPOSE 80

ENTRYPOINT ["./mp2c"]
