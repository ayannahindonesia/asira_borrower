FROM golang:alpine

ARG APPNAME="asira"
ARG APPMODE=""
ARG ENV="dev"

ADD . $GOPATH/src/"${APPNAME}"
WORKDIR $GOPATH/src/"${APPNAME}"

RUN apk add --update git gcc libc-dev;
#  tzdata wget gcc libc-dev make openssl py-pip;

RUN go get -u github.com/golang/dep/cmd/dep

CMD if [ "${ENV}" = "dev" ] ; then \
    dep ensure -v \
    && go build -v -o $GOPATH/bin/"${APPNAME}" \
    # run app mode
    && "${APPNAME}" run "${APPMODE}" \
    # update db structure
    && "${APPNAME}" migrate up; \
    fi

EXPOSE 8000