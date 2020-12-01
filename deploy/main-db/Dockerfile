FROM ubuntu:18.04

ENV DEBIAN_FRONTEND="noninteractive"
ENV TZ="America/New_York"

RUN apt-get update && \
    apt-get -y upgrade && \
    apt-get install -y curl ca-certificates gnupg && \
    apt-get clean

RUN apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys B97B0AFCAA1A47F044F244A07FCC7D46ACCC4CF8
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ bionic-pgdg main" > /etc/apt/sources.list.d/pgdg.list

RUN apt-get update && \
    apt-get install -y postgresql-12 && \
    apt-get clean

RUN mkdir -p /home/postgres/backups && chown -R postgres:postgres /home/postgres/backups

USER postgres

COPY conf/hba.conf /etc/postgresql/12/main/pg_hba.conf
COPY conf/postgresql.conf /etc/postgresql/12/main/postgresql.conf

CMD ["/home/postgres/scripts/run"]
