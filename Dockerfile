FROM golang:alpine

ARG APPNAME="kayacredit"
ARG ENV="dev"

ADD . $GOPATH/src/"${APPNAME}"
WORKDIR $GOPATH/src/"${APPNAME}"

RUN apk add --update git gcc libc-dev;
#  tzdata wget gcc libc-dev make openssl py-pip;

RUN go get -d -v ./...
RUN go install -v ./...

CMD if [ "${ENV}" = "dev" ] ; then \
    go build -v -o $GOPATH/bin/"${APPNAME}" \
    && "${APPNAME}" run "borrower"; \
    fi

EXPOSE 8000