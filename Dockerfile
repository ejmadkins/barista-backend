FROM golang:1.19 AS build
WORKDIR /go/src/app
COPY *.go go.mod go.sum ./
RUN go build -o app

FROM gcr.io/distroless/base-debian11 AS run
WORKDIR /
COPY --from=build /go/src/app/app /app
ENTRYPOINT ["/app"]