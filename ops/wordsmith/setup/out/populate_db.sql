INSERT INTO public."languages" (code) VALUES ('es');

INSERT INTO public."corpora" (_id, language, name) VALUES ('e293ef3c-db3f-447f-b59d-e0975c327f95','es','upc-wiki-corpus');


COPY public."parts_of_speech"(_id,language,corpus_id,code)
FROM '/home/postgres/wordsmith-data/parts_of_speech-1.csv'
DELIMITER ','
CSV HEADER
ENCODING 'latin1';



COPY public."lemmas"(_id,corpus_id,language,lemma_text,part_of_speech_id)
FROM '/home/postgres/wordsmith-data/lemmas-1.csv'
DELIMITER ','
CSV HEADER
ENCODING 'latin1';



COPY public."words"(_id,language,corpus_id,part_of_speech_id,lemma_id,word_text)
FROM '/home/postgres/wordsmith-data/words-1.csv'
DELIMITER ','
CSV HEADER
ENCODING 'latin1';



COPY public."words"(_id,language,corpus_id,part_of_speech_id,lemma_id,word_text)
FROM '/home/postgres/wordsmith-data/words-2.csv'
DELIMITER ','
CSV HEADER
ENCODING 'latin1';



COPY public."part_of_speech_trigram_counts"(_id,language,corpus_id,first_token_id,second_token_id,third_token_id,occurrences)
FROM '/home/postgres/wordsmith-data/part_of_speech_trigram_counts-1.csv'
DELIMITER ','
CSV HEADER
ENCODING 'latin1';



COPY public."word_part_of_speech_counts"(_id,language,corpus_id,word_id,occurrences)
FROM '/home/postgres/wordsmith-data/word_part_of_speech_counts-1.csv'
DELIMITER ','
CSV HEADER
ENCODING 'latin1';



COPY public."word_part_of_speech_counts"(_id,language,corpus_id,word_id,occurrences)
FROM '/home/postgres/wordsmith-data/word_part_of_speech_counts-2.csv'
DELIMITER ','
CSV HEADER
ENCODING 'latin1';


