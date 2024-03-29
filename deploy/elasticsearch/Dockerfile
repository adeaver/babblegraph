FROM ubuntu:18.04

ENV DEBIAN_FRONTEND="noninteractive"
ENV TZ="America/New_York"
ENV ES_VERSION=7.9.3

RUN apt-get update && \
    apt-get -y upgrade && \
    apt-get install -y wget libdigest-sha-perl software-properties-common && \
    apt-get clean

# Download and install elasticsearch
RUN wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-${ES_VERSION}-linux-x86_64.tar.gz
RUN wget https://artifacts.elastic.co/downloads/elasticsearch/elasticsearch-${ES_VERSION}-linux-x86_64.tar.gz.sha512
RUN shasum -a 512 -c elasticsearch-${ES_VERSION}-linux-x86_64.tar.gz.sha512
RUN mkdir elasticsearch && \
    tar -xzf elasticsearch-${ES_VERSION}-linux-x86_64.tar.gz -C elasticsearch --strip-components 1

# Install Java 8
RUN add-apt-repository -y ppa:openjdk-r/ppa
RUN apt-get update && \
    apt-get install -y openjdk-11-jdk && \
    apt-get clean
ENV JAVA_HOME /usr/lib/jvm/java-11-openjdk-amd64/
ENV PATH $PATH:$JAVA_HOME/bin
ENV CLASSPATH $JAVA_HOME/lib/tools.jar
ENV MANPATH $JAVA_HOME/man

COPY limits.conf /etc/security/limits.conf

ENV SERVICE_USER elasticsearch
RUN groupadd $SERVICE_USER && useradd -m -g $SERVICE_USER -l $SERVICE_USER

# Move configuration files and work out permissions
RUN mkdir -p /usr/share/elasticsearch/lib /usr/share/elasticsearch/data/ /usr/share/elasticsearch/log
RUN chmod 777 -R /usr/share/elasticsearch/

COPY --chown=$SERVICE_USER:$SERVICE_USER elasticsearch.yml elasticsearch/config/elasticsearch.yml
COPY --chown=$SERVICE_USER:$SERVICE_USER jvm.options elasticsearch/config/jvm.options
RUN chown -R $SERVICE_USER:$SERVICE_USER /elasticsearch

USER $SERVICE_USER
WORKDIR /elasticsearch

RUN wget https://artifacts.elastic.co/downloads/elasticsearch-plugins/repository-s3/repository-s3-${ES_VERSION}.zip
RUN wget https://artifacts.elastic.co/downloads/elasticsearch-plugins/repository-s3/repository-s3-${ES_VERSION}.zip.sha512
RUN shasum -a 512 -c repository-s3-${ES_VERSION}.zip.sha512

RUN yes | HOST_IP=0.0.0.0 /elasticsearch/bin/elasticsearch-plugin install file:/elasticsearch/repository-s3-${ES_VERSION}.zip

CMD ["/elasticsearch/bin/elasticsearch"]
