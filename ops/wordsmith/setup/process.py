import io, uuid
from corpus_read_defs import corpus_reader
from templates import SQLTemplate
from heapq import heapify, heappop
from util import make_word_ranking_id

LANGUAGE = "es"
CORPUS = ("upc-wiki-corpus", "escrp1upc-wiki-corpus")

print("reading data")
observed_lemmas, observed_words, observed_parts_of_speech, observed_word_bigram_counts = corpus_reader.read_data()

# HELPER FUNCTIONS
def _word_data_for_word_key(word_key):
    external_key = observed_words[word_key]
    word_text, part_of_speech = word_key.split(",")
    external_key_parts = external_key.split(",")
    lemma_key = ",".join(external_key_parts[1:])
    word_id = external_key_parts[0]
    return word_text, part_of_speech, word_id, lemma_key

# GENERATORS
def _make_part_of_speech_line_generator():
    for part_of_speech in observed_parts_of_speech:
        _id = observed_parts_of_speech[part_of_speech]
        yield "{},{},{},{}".format(_id, LANGUAGE, CORPUS[1], part_of_speech)

def _make_lemma_generator():
    for lemma_key in observed_lemmas:
        _id = observed_lemmas[lemma_key]
        lemma_text, part_of_speech = lemma_key.split(",")
        part_of_speech_id = observed_parts_of_speech[part_of_speech]
        yield "{},{},{},{},{}".format(_id, CORPUS[1], LANGUAGE, lemma_text, part_of_speech_id)

def _make_word_generator():
    for word_key in observed_words:
        word_text, part_of_speech, _id, lemma_key = _word_data_for_word_key(word_key)
        lemma_id = observed_lemmas[lemma_key]
        part_of_speech_id = observed_parts_of_speech[part_of_speech]
        yield "{},{},{},{},{},{}".format(_id, LANGUAGE, CORPUS[1], part_of_speech_id, lemma_id, word_text)

def _make_bigram_generator():
    for bigram_key in observed_word_bigram_counts:
        bigram_key_parts = bigram_key.split(",")
        word_key_partition_idx = len(bigram_key_parts) // 2
        first_word_key = ",".join(bigram_key_parts[:word_key_partition_idx])
        second_word_key = ",".join(bigram_key_parts[word_key_partition_idx:])
        first_word_text, _, first_word_id, first_word_lemma_key = _word_data_for_word_key(first_word_key)
        first_word_lemma_id =  = observed_lemmas.get(first_word_lemma_key, None)
        second_word_text, _, second_word_id, second_word_lemma_key = _word_data_for_word_key(second_word_key)
        second_word_lemma_id =  = observed_lemmas.get(second_word_lemma_key, None)
        if first_word_lemma_id is None or second_word_lemma_id is None:
            continue
        bigram_id = "{}-{}".format(first_word_id, second_word_id)
        count = observed_word_bigram_counts[bigram_key]
        yield "{},{},{},{},{},{},{},{}".format(bigram_id, LANGUAGE, CORPUS[1], first_word_text, first_word_lemma_id, second_word_text, second_word_lemma_id, count)

part_of_speech_template = SQLTemplate(
    "parts_of_speech",
    "_id,language,corpus_id,code"
)
print("writing part of speech files")
part_of_speech_template.write_files_for_template(_make_part_of_speech_line_generator, 500)

lemma_template = SQLTemplate(
    "lemmas",
    "_id,corpus_id,language,lemma_text,part_of_speech_id"
)
print("writing lemma files")
lemma_template.write_files_for_template(
    _make_lemma_generator,
    300000
)

word_template = SQLTemplate(
    "words",
    "_id,language,corpus_id,part_of_speech_id,lemma_id,word_text"
)
print("writing word files")
word_template.write_files_for_template(
    _make_word_generator,
    300000
)

word_bigram_counts_template = SQLTemplate(
    "word_bigram_counts",
    "_id,language,corpus_id,first_word_text,first_word_lemma_id,second_word_text,second_word_lemma_id,count"
)
word_bigram_counts_template.write_files_for_template(
    _make_bigram_generator,
    300000
)

print("writing sql file")
with io.open("/out/populate_db.sql", "w", encoding="latin1") as f:
    f.write("INSERT INTO public.\"languages\" (code) VALUES ('es');\n\n")
    f.write("INSERT INTO public.\"corpora\" (_id, language, name) VALUES ('{}','es','{}');\n\n".format(CORPUS[1], CORPUS[0]))
    for template in part_of_speech_template.yield_sql_template():
        f.write(template)
    for template in lemma_template.yield_sql_template():
        f.write(template)
    for template in word_template.yield_sql_template():
        f.write(template)
    for template in word_bigram_counts_template.yield_sql_template():
        f.write(template)

print("done!")
