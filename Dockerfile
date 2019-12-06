 # === Lintas Arta's Dockerfile ===
FROM golang:alpine  AS build-env

ARG APPNAME="asira_borrower"
ARG CONFIGPATH="/go/src"

#RUN adduser -D -g '' golang
#USER root

ADD . $GOPATH/src/"${APPNAME}"
WORKDIR $GOPATH/src/"${APPNAME}"

RUN apk add --update git gcc libc-dev;
RUN apk --no-cache add curl
#  tzdata wget gcc libc-dev make openssl py-pip;
RUN go get -u github.com/golang/dep/cmd/dep

RUN cd $GOPATH/src/"${APPNAME}"
RUN cp deploy/conf.yaml config.yaml
RUN dep ensure -v
RUN go build -v -o "${APPNAME}-res"

RUN ls -alh $GOPATH/src/
RUN ls -alh $GOPATH/src/"${APPNAME}"
RUN ls -alh $GOPATH/src/"${APPNAME}"/vendor
RUN pwd

FROM alpine

WORKDIR /go/src/
COPY --from=build-env /go/src/asira_borrower/asira_borrower-res /go/src/asira_borrower
COPY --from=build-env /go/src/asira_borrower/deploy/conf.yaml /go/src/config.yaml
COPY --from=build-env /go/src/asira_borrower/permissions.yaml /go/src/permissions.yaml
COPY --from=build-env /go/src/asira_borrower/migration/ /go/src/migration/
RUN chmod -R 775 migration

RUN pwd
#ENTRYPOINT /app/asira_borrower-res
CMD ["/go/src/asira_borrower","run"]

EXPOSE 8000

# === DEFAULT ===
# FROM golang:alpine

# ARG APPNAME="asira_borrower"
# ARG CONFIGPATH="$$GOPATH/src/asira_borrower"

# ADD . $GOPATH/src/"${APPNAME}"
# WORKDIR $GOPATH/src/"${APPNAME}"

# RUN apk add --update git gcc libc-dev tzdata;
# #  tzdata wget gcc libc-dev make openssl py-pip;

# ENV TZ=Asia/Jakarta

# RUN go get -u github.com/golang/dep/cmd/dep

# CMD if [ "${APPENV}" = "staging" ] || [ "${APPENV}" = "production" ] ; then \
#         cp deploy/conf.yaml config.yaml ; \
#     elif [ "${APPENV}" = "dev" ] ; then \
#         cp deploy/dev-config.yaml config.yaml ; \
#     fi \
#     && dep ensure -v \
#     && go build -v -o $GOPATH/bin/"${APPNAME}" \
#     && "${APPNAME}" run \
#     && "${APPNAME}" migrate up \
# EXPOSE 8000