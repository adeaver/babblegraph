import io, uuid
from corpus_read_defs import corpus_reader
from templates import SQLTemplate

LANGUAGE = "es"
CORPUS = ("upc-wiki-corpus", uuid.uuid4())

print("reading data")
observed_lemmas, observed_words, observed_parts_of_speech, part_of_speech_trigram_counts, word_part_of_speech_counts = corpus_reader.read_data()

def _make_part_of_speech_line_generator():
    for part_of_speech, _id in observed_parts_of_speech:
        yield "{},{},{},{}".format(_id, LANGUAGE, CORPUS[1], part_of_speech)

def _make_lemma_generator():
    for lemma_key, _id in observed_lemmas:
        lemma_text, part_of_speech = lemma_key.split(",")
        part_of_speech_id = observed_parts_of_speech[part_of_speech]
        yield "{},{},{},{},{}".format(_id, CORPUS[1], LANGUAGE, lemma_text, part_of_speech_id)

def _make_word_generator():
    for word_key, external_key in observed_words:
        word_text, part_of_speech = word_key.split(",")
        _id, lemma_key = external_key.split(",")
        lemma_id = obsered_lemmas[lemma_key]
        part_of_speech_id = observed_parts_of_speech[part_of_speech]
        yield "{},{},{},{},{},{}".format(_id, LANGUAGE, CORPUS[1], part_of_speech_id, lemma_id, word_text)

def _make_trigram_count_generator():
    for trigram_key, count in part_of_speech_trigram_counts:
        parts_of_speech = trigram_key.split(",")
        part_of_speech_ids = [observed_parts_of_speech[part_of_speech] for part_of_speech in parts_of_speech]
        yield "{},{},{},{},{},{},{}".format(uuid.uuid4(), LANGUAGE, CORPUS[1], part_of_speech_ids[0], part_of_speech_ids[1], part_of_speech_ids[2], count)

def _make_word_part_of_speech_count_generator():
    for word_key, count in word_part_of_speech_counts:
        external_key = observed_words[word_key]
        word_id, _ = external_key.split(",")
        yield "{},{},{},{},{}".format(uuid.uuid4(), LANGUAGE, CORPUS[1], word_id, count)

part_of_speech_template = SQLTemplate(
    "part_of_speech",
    "_id,language,corpus_id,code"
)
print("writing part of speech files")
part_of_speech_template.write_files_for_template(_make_part_of_speech_line_generator, 500)

lemma_template = SQLTemplate(
    "lemma",
    "_id,corpus_id,language,lemma_text,part_of_speech_id"
)
print("writing lemma files")
lemma_template.write_files_for_template(
    _make_lemma_generator,
    30000
)

word_template = SQLTemplate(
    "words",
    "_id,language,corpus_id,part_of_speech_id,lemma_id,word_text"
)
print("writing word files")
word_template.write_files_for_template(
    _make_word_generator,
    30000
)

part_of_speech_trigram_template =  SQLTemplate(
    "part_of_speech_trigram_counts",
    "_id,language,corpus_id,first_token_id,second_token_id,third_token_id,ocurrences"
)
print("writing part of speech trigram files")
part_of_speech_trigram_template.write_files_for_template(
    _make_trigram_count_generator,
    30000
)

word_part_of_speech_template = SQLTemplate(
    "word_part_of_speech_counts",
    "_id,language,corpus_id,word_id,occurrences"
)
print("writing word counts files")
word_part_of_speech_template.write_files_for_template(
    _make_word_part_of_speech_count_generator,
    30000
)

print("writing sql file")
with io.open("./out/populate_db.sql", "w", encoding="latin1") as f:
    f.write("INSERT INTO languages (code) VALUES ('es');\n\n")
    f.write("INSERT INTO corpora (_id, language, name) VALUES ({},'es',{});\n\n".format(CORPUS[1], CORPUS[0]))
    for template in part_of_speech_template.yield_sql_template():
        f.write(template)
    for template in lemma_template.yield_sql_template():
        f.write(template)
    for template in word_template.yield_sql_template():
        f.write(template)
    for template in part_of_speech_trigram_template.yield_sql_template():
        f.write(template)
    for template in word_part_of_speech_template.yield_sql_template():
        f.write(template)

print("done!")
