FROM golang:latest AS build
ADD . /go/src/github.com/frozzare/statscoll
RUN cd /go/src/github.com/frozzare/statscoll && CGO_ENABLED=0 GOOS=linux go build -o statscoll

FROM scratch
WORKDIR /app
COPY --from=build /go/src/github.com/frozzare/statscoll/statscoll /app/
ENTRYPOINT ["/app/statscoll"]