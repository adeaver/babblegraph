from corpus_read_defs import corpus_reader

import re
from typing import NamedTuple, List, Tuple
import xml.etree.ElementTree as ET

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


def _get_words_data():
    print("reading data")
    observed_lemmas, observed_words, _, _ = corpus_reader.read_data()
    print("post processing")
    lemmas_by_lemma_id = {}
    for lemma_key, lemma_id in observed_lemmas.items():
        lemma, _ = lemma_key.split(",")
        lemmas = lemmas_by_lemma_id.get(lemma_id, [])
        lemmas.append(lemma)
        lemmas_by_lemma_id[lemma_id] = lemmas
    words_to_lemma_id = {}
    for word_key, word_value in observed_words.items():
        word_text, _ = word_key.split(",")
        value_parts = word_value.split(",")
        lemma_key = ",".join(value_parts[1:])
        lemma_ids = words_to_lemma_id.get(word_text, [])
        lemma_ids.append(observed_lemmas[lemma_key])
        words_to_lemma_id[word_text] = lemma_ids
    return words_to_lemma_id, lemmas_by_lemma_id

def _make_lemma_phrases(start, lemmas):
    if len(lemmas) == 0:
        return [start]
    current_lemmas = lemmas[0]
    next_lemmas = lemmas[1:]
    out_lemmas = []
    for lemma in current_lemmas:
        out_lemmas += _make_lemma_phrases(f"{start} {lemma}", next_lemmas)
    return out_lemmas

count = 0
lemma_ids_by_word_text, lemmas_by_lemma_id = _get_words_data()
for w in _get_words_from_xml("./data-defs/es-en.xml"):
    if " " not in w.word_text:
        continue
    words = re.split(r" +", w.word_text)
    lemmas = []
    for word in words:
        lemma_ids = lemma_ids_by_word_text.get(word, None)
        if lemma_ids is None:
            break
        word_lemmas = []
        for lemma_id in lemma_ids:
            word_lemmas += [
                lemma
                for lemma in lemmas_by_lemma_id.get(lemma_id, [])
            ]
        lemmas.append(word_lemmas)
    if len(lemmas) == len(words):
        lemma_phrases = _make_lemma_phrases("", lemmas)
        lemma_phrase_to_phrase = { lemma_phrase.strip(): w for lemma_phrase in lemma_phrases }
        print(lemma_phrase_to_phrase)
    else:
        count += 1

print(f"Missed {count}")

