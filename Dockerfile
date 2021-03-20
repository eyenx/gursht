FROM golang:alpine as go
FROM go as build

WORKDIR /

COPY . /

RUN apk add --no-cache \
    build-base \
    git \
    && go build 

FROM alpine

ENTRYPOINT ["/gursht"]

COPY --from=build /gursht /gursht
