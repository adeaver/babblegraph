FROM golang:1.14.9 AS builder

ENV SERVICE_USER somebody
ENV SERVICE_HOME /usr/local/go/src/babblegraph

WORKDIR $SERVICE_HOME

COPY . .

RUN go mod download
RUN go mod verify

RUN go install -v babblegraph/services/web

FROM debian:buster

ENV SERVICE_USER somebody
ENV SERVICE_HOME /usr/local/go/src/babblegraph

RUN groupadd $SERVICE_USER && useradd -m -g $SERVICE_USER -l $SERVICE_USER
RUN mkdir -p /service && chown -R $SERVICE_USER:$SERVICE_USER /service

WORKDIR /service
RUN mkdir -p ./templates && chown -R $SERVICE_USER:$SERVICE_USER ./templates
COPY --chown=$SERVICE_USER:$SERVICE_USER --from=builder /usr/local/go/src/babblegraph/util/email/templates ./templates

COPY --chown=$SERVICE_USER:$SERVICE_USER --from=builder /usr/local/go/bin/web .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER $SERVICE_USER

EXPOSE 8080

CMD ["./web"]
