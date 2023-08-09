FROM golang:1.20 as build-env
WORKDIR /app
COPY ./src .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o binder .

FROM debian:buster-slim
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates
WORKDIR /root/binder-app
COPY --from=build-env /app/binder .
ENTRYPOINT [ "./binder" ]