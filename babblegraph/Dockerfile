FROM golang:1.14.9 AS builder

ARG SERVICE_NAME

ENV SERVICE_USER somebody
ENV SERVICE_HOME /usr/local/go/src/babblegraph

WORKDIR $SERVICE_HOME

COPY . .

RUN go mod download
RUN go mod verify

RUN go install -v babblegraph/services/$SERVICE_NAME

FROM debian:buster

ENV SERVICE_USER somebody
ENV SERVICE_HOME /usr/local/go/src/babblegraph

RUN groupadd $SERVICE_USER && useradd -m -g $SERVICE_USER -l $SERVICE_USER
RUN mkdir -p /service && chown -R $SERVICE_USER:$SERVICE_USER /service

WORKDIR /service
COPY --chown=0.0 --from=builder /usr/local/go/bin/$SERVICE_NAME ./run-service

USER $SERVICE_USER

CMD ["./run-service"]
