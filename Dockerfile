FROM golang:alpine

ARG APPNAME="kayacredit"
ARG ENV="dev"

ADD . $GOPATH/src/"${APPNAME}"
WORKDIR $GOPATH/src/"${APPNAME}"

RUN apk add --update git;
#  tzdata wget gcc libc-dev make openssl py-pip;

RUN go get -u github.com/golang/dep/cmd/dep

CMD if [ "$ENV" = "dev" ] ; then \
    dep ensure -v \
    && go build -v -o $GOPATH/bin/"${APPNAME}" \
    && "${APPNAME}" run; \
    fi

EXPOSE 8000