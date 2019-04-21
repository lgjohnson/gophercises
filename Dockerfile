FROM golang:1.12.4

ARG EXERCISE

WORKDIR $GOPATH/src/app

COPY ./exercise-$EXERCISE .

RUN go get -d -v ./...

RUN go install -v ./...

ENTRYPOINT $GOPATH/bin/app -random -csv=some.csv
