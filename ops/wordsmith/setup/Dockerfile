FROM python:3.8.5

WORKDIR /script

# Get the appropriate files
RUN wget https://www.cs.upc.edu/~nlp/wikicorpus/tagged.es.tgz
RUN mkdir data/
RUN tar -xzvf tagged.es.tgz -C data/

RUN mkdir out/
COPY process.py .
COPY process2.py .
COPY reader.py .
COPY templates.py .
COPY corpus_read_defs.py .
COPY special_tokens.py .
COPY util.py .

CMD [ "python", "-u", "process.py" ]
