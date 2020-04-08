FROM golang:1.14-alpine as build
RUN apk add --update upx build-base git
RUN go get -u github.com/gobuffalo/packr/packr
WORKDIR /go/src/github.com/botto/stc
COPY . ./
RUN go get .
RUN packr
RUN GGO_ENABLED=0 GOOS=linux go build -a -ldflags="-s -w" -o dist/stc
RUN upx dist/stc

FROM alpine:3.11
WORKDIR /root/
COPY --from=build /go/src/github.com/botto/stc/dist/stc .
CMD ["./stc"]