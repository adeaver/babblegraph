import uuid, json
from phrase_definitions_filters import DEFINITIONS_TO_FILTER
from corpus_read_defs import corpus_reader, TOKEN_AL, TOKEN_DEL

import re
from typing import NamedTuple, List, Tuple
import xml.etree.ElementTree as ET

LANGUAGE = "es"
CORPUS = ("upc-wiki-corpus", "escrp1upc-wiki-corpus")
DEFINITION_CORPUS = ("escrp1mananoreboton-definitions", "mananoreboton-definitions", "es")

class WordDefinition(NamedTuple):
    word_text: str
    part_of_speech: str
    definition: str

def _get_words_from_xml(file_name):
    tree = ET.parse(file_name)
    root = tree.getroot()
    for word in root.findall('./l/w'):
        text, pos, defn = None, None, None
        for child in word:
            if child.tag == "c":
                text = child.text
            elif child.tag == "d":
                defn = child.text
            elif child.tag == "t":
                pos = child.text
            else:
                print(f"Ignoring tag {child.tag}")
        yield WordDefinition(
            word_text=text,
            part_of_speech=pos,
            definition=defn,
        )

MINIMUM_COUNT = 7
FORMER_MINIMUM_COUNT = 10

def _get_lemmas_by_lemma_id():
    with open("/out/observed_lemmas__chunked_0.json", "r") as f:
        observed_lemmas = json.loads(f.read())
    with open("/out/lemma_counts__chunked_0.json", "r") as f:
        lemma_counts = json.loads(f.read())
    lemmas_by_lemma_id = {}
    new_lemmas, filtered_lemmas = {}, {}
    for lemma_key, lemma_id in observed_lemmas.items():
        count_for_lemma = lemma_counts.get(lemma_key, 0)
        if count_for_lemma < MINIMUM_COUNT:
            continue
        elif count_for_lemma >= MINIMUM_COUNT or count_for_lemma < FORMER_MINIMUM_COUNT:
            new_lemmas[lemma_key] = lemma_id
        filtered_lemmas[lemma_key] = lemma_id
        lemma, _ = lemma_key.split(",")
        if lemma == "<<START>>":
            continue
        lemmas = lemmas_by_lemma_id.get(lemma_id, [])
        lemmas.append(lemma)
        lemmas_by_lemma_id[lemma_id] = lemmas
    return lemmas_by_lemma_id, filtered_lemmas, new_lemmas

def _get_words_to_lemma_id(observed_parts_of_speech, filtered_lemmas, new_lemmas):
    with open("/out/observed_words__chunked_0.json", "r") as f:
        observed_words = json.loads(f.read())
    words_to_lemma_id = {}
    new_words = {}
    for word_key, word_value in observed_words.items():
        word_text, _ = word_key.split(",")
        value_parts = word_value.split(",")
        lemma_key = ",".join(value_parts[1:])
        if filtered_lemmas.get(lemma_key, None) is None:
            continue
        if new_lemmas.get(lemma_key, None) is not None:
            new_words[word_key] = word_value
        lemma_ids = words_to_lemma_id.get(word_text, [])
        lemma_ids.append(filtered_lemmas[lemma_key])
        words_to_lemma_id[word_text] = lemma_ids
    with open("/out/phrase-definitions-0.sql", "a") as f:
        for word_key, word_value in new_words.items():
            external_key = observed_words.get(word_key, None)
            if external_key is None:
                continue
            word_text, part_of_speech = word_key.split(",")
            external_key_parts = external_key.split(",")
            lemma_key = ",".join(external_key_parts[1:])
            word_id = external_key_parts[0]
            lemma_id = filtered_lemmas[lemma_key]
            if part_of_speech not in observed_parts_of_speech:
                continue
            part_of_speech_id = observed_parts_of_speech[part_of_speech]
            f.write(f"""INSERT INTO \"public\".words (
                _id, language, corpus_id, part_of_speech_id, lemma_id, word_text
            ) VALUES (
                {word_id}, {LANGUAGE}, {CORPUS[1]}, {part_of_speech_id}, {lemma_id}, {word_text}
            )\n\n""")
    return words_to_lemma_id

def _get_words_data(observed_parts_of_speech):
    lemmas_by_lemma_id, filtered_lemmas, new_lemmas = _get_lemmas_by_lemma_id()
    words_to_lemma_id = _get_words_to_lemma_id(observed_parts_of_speech, filtered_lemmas, new_lemmas)
    return lemmas_by_lemma_id, words_to_lemma_id, new_lemmas

def _make_lemma_phrases(start, lemmas):
    if len(lemmas) == 0:
        return [start]
    current_lemmas = lemmas[0]
    next_lemmas = lemmas[1:]
    out_lemmas = []
    for lemma in current_lemmas:
        out_lemmas += _make_lemma_phrases(f"{start} {lemma}", next_lemmas)
    return out_lemmas


with open("/out/phrase-definitions-0.sql", "w") as f:
    for special_token in [ TOKEN_AL, TOKEN_DEL ]:
        f.write(f"""INSERT INTO
            \"public\".words (
                _id, language, corpus_id, part_of_speech_id, lemma_id, word_text
            ) VALUES (
                '{special_token.get_word_id()}', '{LANGUAGE}', '{CORPUS[1]}',
                '{special_token.get_part_of_speech_id()}', '{special_token.get_lemma_id()}',
                '{special_token.get_token()}'
            )\n""")

definition_ids = {}

def _process_phrase(file_number, words_to_lemma_id, lemmas_by_lemma_id, phrase):
    words = re.split(r" +", phrase.word_text)
    lemmas = []
    for word in words:
        lemma_ids = words_to_lemma_id.get(word, None)
        if lemma_ids is None:
            return phrase, 0
        word_lemmas = []
        for lemma_id in lemma_ids:
            word_lemmas += [
                lemma
                for lemma in lemmas_by_lemma_id.get(lemma_id, [])
            ]
        lemmas.append(word_lemmas)
    if len(lemmas) != len(words):
        return phrase, 0
    lemma_phrases = set(_make_lemma_phrases("", lemmas))
    with open(f"/out/phrase-definitions-{file_number}.sql", "a") as f:
        phrase_definition = phrase.definition
        try:
            escaped_definition = re.sub(r"\'", "\\'",  phrase_definition)
            escaped_phrase = re.sub(r"\'", "\\'", phrase.word_text)
        except:
            print(f"Error on {phrase_definition} or {phrase.word_text}")
            return phrase, 0
        definition_id = str(uuid.uuid4())
        while definition_id in definition_ids:
            definition_id = str(uuid.uuid4())
        definition_ids[definition_id] = True
        f.write(f"""INSERT INTO
            \"public\".phrase_definitions (
                _id, language, corpus_id, phrase, definition
            ) VALUES (
                '{definition_id}', '{LANGUAGE}', '{DEFINITION_CORPUS[0]}', '{escaped_phrase}', '{escaped_definition}'
            )\n\n""")
        inserted_lines = 0
        for lemma_phrase in lemma_phrases:
            escaped_lemma_phrase = re.sub(r"\'", "\\'", lemma_phrase.strip())
            inserted_lines += 1
            f.write(f"""INSERT INTO
                \"public\".lemma_phrase_definition_mappings (
                    language, corpus_id, lemma_phrase, phrase_definition_id
                ) VALUES (
                    '{LANGUAGE}', '{CORPUS[1]}', '{escaped_lemma_phrase}', '{definition_id}'
                )\n\n""")
        return None, inserted_lines

inserted_lines = 0
MAX_CHUNK = 20000

count = 0
phrases = []
for w in _get_words_from_xml("./data-defs/es-en.xml"):
    if " " not in w.word_text:
        continue
    text, definition = DEFINITIONS_TO_FILTER.get(w.word_text, (w.word_text, w.definition))
    if text is None:
        continue
    phrases.append(WordDefinition(
        word_text=text,
        definition=definition,
        part_of_speech=w.part_of_speech,
    ))

with open("/out/observed_parts_of_speech__chunked_0.json", "r") as f:
    observed_parts_of_speech = json.loads(f.read())
lemmas_by_lemma_id, words_to_lemma_id, new_lemmas = _get_words_data(observed_parts_of_speech)
with open("/out/phrase-definitions-0.sql", "a") as f:
    for lemma_key, _id in new_lemmas.items():
        lemma_text, part_of_speech = lemma_key.split(",")
        if observed_parts_of_speech.get(part_of_speech, None) is None:
            continue
        part_of_speech_id = observed_parts_of_speech[part_of_speech]
        f.write(f"""INSERT INTO
            \"public\".lemmas (
                _id, corpus_id, language, lemma_text, part_of_speech_id
            ) VALUES (
                {_id}, {CORPUS[1]}, {LANGUAGE}, {lemma_text}, {part_of_speech_id})
            )\n\n""")

new_phrases = []
for phrase in phrases:
    file_number = inserted_lines // MAX_CHUNK
    p, added_lines = _process_phrase(file_number, words_to_lemma_id, lemmas_by_lemma_id, phrase)
    if p is not None:
        new_phrases.append(p)
    inserted_lines += added_lines
