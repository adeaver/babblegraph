FROM python:3.8.5

WORKDIR /script

# Get the appropriate files
RUN wget https://www.cs.upc.edu/~nlp/wikicorpus/tagged.es.tgz
RUN mkdir data/
RUN tar -xzvf tagged.es.tgz -C data/

RUN mkdir out/
RUN mkdir data-dump/
COPY process.py .
COPY analyze.py .
COPY reader.py .

CMD [ "python", "-u", "process.py" ]
