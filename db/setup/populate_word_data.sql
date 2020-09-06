COPY lemmas(_id, lemma, part_of_speech, language)
FROM '/setup/out/lemmas.csv'
DELIMITER ','
CSV HEADER
ENCODING 'latin1';

COPY words(lemma_id, word, part_of_speech, language)
FROM '/setup/out/words.csv'
DELIMITER ','
CSV HEADER
ENCODING 'latin1';
