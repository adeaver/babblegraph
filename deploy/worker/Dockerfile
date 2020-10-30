FROM golang:1.14.9 AS builder

ENV SERVICE_USER somebody
ENV SERVICE_HOME /usr/local/go/src/babblegraph

WORKDIR $SERVICE_HOME

COPY . .

RUN go mod download
RUN go mod verify

RUN go install -v babblegraph/services/worker

FROM debian:buster

ENV SERVICE_USER somebody
ENV SERVICE_HOME /usr/local/go/src/babblegraph

RUN mkdir -p /service 

WORKDIR /service
COPY --from=builder /usr/local/go/bin/worker .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

CMD ["./worker"]
