FROM python:3.8.5

# Get the appropriate files
RUN wget https://www.cs.upc.edu/~nlp/wikicorpus/tagged.es.tgz
RUN mkdir data/
RUN tar -xzvf tagged.es.tgz -C data/

RUN pip install psycopg2

COPY process.py .

CMD [ "python", "process.py" ]
