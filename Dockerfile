FROM golang:1.9.3-alpine3.7 AS build
COPY . /go/src/github.com/docker/perfkit
RUN go install github.com/docker/perfkit

FROM alpine:3.7
COPY --from=build /go/bin/perfkit bin/
ENTRYPOINT ["perfkit"]
CMD ["--help"]