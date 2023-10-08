FROM golang:latest

RUN mkdir komp-registry
COPY . komp-registry/

RUN cd komp-registry && go build  -o komp-registry-backend *.go

ENTRYPOINT ["/go/komp-registry/komp-registry-backend"]