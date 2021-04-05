FROM alpine:3.11.5 AS frontend-builder

# npm
RUN apk update && apk add npm

WORKDIR /usr/app

# Install dependencies
COPY website/package.json .
COPY website/package-lock.json .
RUN npm install

RUN mkdir ./dist

COPY website/dist/index.html ./dist
COPY website/dist/logo.png ./dist
COPY website/src ./src
COPY website/webpack.common.js .
COPY website/webpack.prod.js .
COPY website/babel.config.js .
COPY website/tsconfig.json .

RUN npm run build

FROM golang:1.14.9 AS backend-builder

ENV SERVICE_USER somebody
ENV SERVICE_HOME /usr/local/go/src/babblegraph

WORKDIR $SERVICE_HOME

COPY backend/src/babblegraph .

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
COPY --chown=$SERVICE_USER:$SERVICE_USER --from=backend-builder /usr/local/go/src/babblegraph/actions/email/templates ./templates

RUN mkdir -p ./static && chown -R $SERVICE_USER:$SERVICE_USER ./static
COPY --chown=$SERVICE_USER:$SERVICE_USER --from=frontend-builder /usr/app/dist/bundle.js ./static
COPY --chown=$SERVICE_USER:$SERVICE_USER --from=frontend-builder /usr/app/dist/index.html ./static
COPY --chown=$SERVICE_USER:$SERVICE_USER --from=frontend-builder /usr/app/dist/logo.png ./static

COPY --chown=$SERVICE_USER:$SERVICE_USER --from=backend-builder /usr/local/go/bin/web .
COPY --from=backend-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

USER $SERVICE_USER

EXPOSE 8080

CMD ["./web"]
