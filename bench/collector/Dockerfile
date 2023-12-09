
FROM golang:1.20

COPY plantdcollector .
COPY conf.yaml .

EXPOSE 4317 8888 8889

CMD ["/bin/sh", "-c", "./plantdcollector --config=conf.yaml"]

