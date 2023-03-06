FROM golang:alpine as build
WORKDIR /opt/app
COPY go.mod ./
COPY go.sum ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server

FROM scratch
WORKDIR /opt/app
COPY --from=build /opt/app/server server
CMD ["./server"]