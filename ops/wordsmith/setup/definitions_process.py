import re
from typing import NamedTuple, List, Tuple
import xml.etree.ElementTree as ET

from definitions_part_of_speech import get_wordsmith_part_of_speech_id
from util import make_lemma_id


corpus = ("escrp1mananoreboton-definitions", "mananoreboton-definitions", "es")

class InsertableDefinition(NamedTuple):
    part_of_speech_id: str
    lemma_id: str
    english_definition: str
    extra_part_of_speech_info: str
    language: str
    corpus: str

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


def _make_insertable_wordsmith_definition(w: WordDefinition) -> List[InsertableDefinition]:
    wordsmith_part_of_speech_ids = _process_part_of_speech(w.part_of_speech)
    out = []
    for wordsmith_part_of_speech_id, extra_part_of_speech_info in wordsmith_part_of_speech_ids:
        lemma_id = make_lemma_id(w.word_text, wordsmith_part_of_speech_id)
        out.append(InsertableDefinition(
            part_of_speech_id=wordsmith_part_of_speech_id,
            lemma_id=lemma_id,
            english_definition=w.definition,
            extra_part_of_speech_info=extra_part_of_speech_info,
            language="es",
            corpus=corpus[0]
        ))
    return out

def _process_part_of_speech(pos: str) -> List[Tuple[str, str]]:
    """Gets wordsmith part of speech ids for raw definitions part of speech"""
    part_of_speech_tags, extra_info = _clean_part_of_speech(pos)
    wordsmith_ids = {}
    for tag in part_of_speech_tags:
        wordsmith_id = get_wordsmith_part_of_speech_id(tag)
        if wordsmith_id is None:
            continue
        wordsmith_ids[wordsmith_id] = True
    return [[_id, extra_info] for _id in wordsmith_ids.keys()]


def _clean_part_of_speech(raw_part_of_speech: str) -> Tuple[List[str], str]:
    """Separates out the definitions part of speech text and the extra information on it"""
    raw_part_of_speech_matches = re.findall("\{.*\}", raw_part_of_speech)
    if len(raw_part_of_speech_matches) == 0:
        print("no part of speech")
        return [], raw_part_of_speech
    elif len(raw_part_of_speech_matches) == 1:
        extras = raw_part_of_speech.replace(raw_part_of_speech_matches[0], "").strip()
        return [raw_part_of_speech_matches[0]], extras
    print("multiple parts of speech")
    extras = raw_part_of_speech
    for match in raw_part_of_speech_matches:
        extras = extras.replace(match, "").strip()
    return raw_part_of_speech_matches, extras


def _make_definitions_postgres_file(idx: int, rows: List[InsertableDefinition]):
    psql = []
    for r in rows:
        psql.append(f"(\'{r.language}\', \'{r.corpus}\', $${r.lemma_id}$$, $${r.english_definition}$$, \'{r.part_of_speech_id}\', $${r.extra_part_of_speech_info}$$)")
    with open(f"/out/definitions-{idx}.sql", "w") as f:
        if idx == 1:
            corpus_string = f"INSERT INTO \"public\".corpora (_id, name, language) VALUES ('{corpus[0]}', '{corpus[1]}', '{corpus[2]}') ON CONFLICT DO NOTHING;\n\n"
            f.write(corpus_string)
        value_string = ", ".join(psql)
        insert_sql = f"""INSERT INTO definition_mappings (language, corpus_id, lemma_id, english_definition, part_of_speech_id, extra_info)
        SELECT def_rows.*
        FROM (VALUES {value_string}) AS def_rows(language, corpus_id, lemma_id, english_definition, part_of_speech_id, extra_info)
        JOIN lemmas ON def_rows.lemma_id = lemmas._id"""
        f.write(insert_sql)

CHUNK_SIZE = 250000
CURRENT_CHUNK = 1
insertable = []
for w in _get_words_from_xml("./data/es-en.xml"):
    insertable_definitions = _make_insertable_wordsmith_definition(w)
    for d in insertable_definitions:
        insertable.append(d)
        if len(insertable) >= CHUNK_SIZE:
            _make_definitions_postgres_file(CURRENT_CHUNK, insertable)
            insertable = []
            CURRENT_CHUNK += 1

if len(insertable) > 0:
    _make_definitions_postgres_file(CURRENT_CHUNK, insertable)
