import uuid
from reader import Reader

START_TOKEN = "<<START>>"

part_of_speech_trigrams = [START_TOKEN, START_TOKEN, START_TOKEN]
part_of_speech_trigram_counts = dict()
observed_parts_of_speech = dict()
observed_lemmas = dict()
observed_words = dict()
word_part_of_speech_counts = dict()

def get_data_from_read():
    return observed_lemmas, observed_words, observed_parts_of_speech, part_of_speech_trigram_counts, word_part_of_speech_counts

def process_text_line(word, lemma, pos):
    if not _is_text_line_valid(word, lemma, pos):
        # print("Filtering out line ({}, {}, {})".format(word, lemma, pos))
        return
    processed_part_of_speech = _process_part_of_speech(pos)
    _handle_part_of_speech(processed_part_of_speech)
    lemma_key = _handle_lemma(lemma, pos)


def process_empty_line():
    _reset_part_of_speech_trigrams()

def process_start_doc(_):
    _reset_part_of_speech_trigrams()

# Dates, numbers, and punctuation should be processed out
invalid_part_of_speech_categories = ["F", "Z", "W"]

def _is_text_line_valid(word, lemma, pos):
    part_of_speech_category = pos[0]
    if part_of_speech_category in invalid_part_of_speech_categories: # Punctuation
        return False
    elif _contains_invalid_character(lemma) or _contains_invalid_character(word):
        return False
    return True

invalid_characters = ["_", ".", ",", "\"", ":", "$", "\'"]
def _contains_invalid_character(word):
    for char in word:
        if char in invalid_characters:
            return True
    return False

def _handle_lemma(lemma, part_of_speech):
    part_of_speech_category = part_of_speech[0]
    lemma_key = "{},{}".format(lemma, part_of_speech_category)
    if lemma_key not in observed_lemmas:
        observed_lemmas[lemma_key] = uuid.uuid4()
    return lemma_key

def _handle_word(word, part_of_speech, lemma_key):
    word_key = "{},{}".format(word, part_of_speech)
    if word_key not in observed_words:
        observed_words[word_key] = "{},{}".format(uuid.uuid4(), lemma_key)
    word_part_of_speech_counts[word_key] = word_part_of_speech_counts.get(word_key, 0) + 1

def _handle_part_of_speech(processed_part_of_speech):
    if processed_part_of_speech not in observed_parts_of_speech:
        observed_parts_of_speech[processed_part_of_speech] = uuid.uuid4()
    if processed_part_of_speech[0] not in observed_parts_of_speech:
        observed_parts_of_speech[processed_part_of_speech[0]] = uuid.uuid4()
    part_of_speech_trigrams.pop(0)
    part_of_speech_trigrams.append(processed_part_of_speech)
    trigram_key = "{},{},{}".format(part_of_speech_trigrams[0], part_of_speech_trigrams[1], part_of_speech_trigrams[2])
    part_of_speech_trigram_counts[trigram_key] = part_of_speech_trigram_counts.get(trigram_key, 0) + 1

def _reset_part_of_speech_trigrams():
    part_of_speech_trigrams = [START_TOKEN, START_TOKEN, START_TOKEN]

def _process_part_of_speech(pos):
    """Processes parts of speech to get a smaller list of tags
    The goal of this function is to create a list of tags that helps filter out
    documents and create accurate part of speech tags without overtraining it"""
    # https://freeling-user-manual.readthedocs.io/en/latest/tagsets/tagset-es/
    category = pos[0]
    if category == "A": # Adjective
        """Adjectives include the category, gender, and number. Possibilities should be:
        AMS -> singular male, AFS -> singular female, ACS -> common singular, AxP -> gender plural, AxN -> gender invariable"""
        return pos[0] + pos[3] + pos[4]
    elif category == "C": # Conjunction
        """Conjunctions include category and type (coordinating or subordinating).
        Possibilities are CC and CS"""
        return pos[0:2]
    elif category == "D": # Determiner
        """Determiners include category, type, and number. Possibilities are:
        Type Possibilities are A (article), D (demonstrative), I (indefinite), P (possessive),
        T (interrogative), E (exclamative).
        Number Possibilities are: S (singular), P (plural), N (invariable)
        """
        return pos[0] + pos[1] + pos[4]
    elif category == "N": # Noun
        """Nouns include category, gender, number. Possibilities are:
        NFS (female singular), NMS (male singular), NCS (common singular), NxP (gender plural), NxC (gender invariant)
        """
        return pos[0] + pos[2] + pos[3]
    elif category == "P": # Pronoun
        """Pronouns include category, type, gender, and number"""
        return pos[0] + pos[1] + pos[3] + pos[4]
    elif category == "R": # Adverb
        """Adverbs only contain category"""
        return pos[0]
    elif category == "S": # Adposition (general name for preposition and postposition)
        """Spanish only has prepositions, so include the entire tag"""
        return pos
    elif category == "V": # Verb
        """Verbs include category, mood, tense, person, number"""
        return pos[0] + pos[2] + pos[3] + pos[4] + pos[5]
    elif category == "I": # Interjection
        """Returns only the category"""
        return pos[0]
    else:
        print("Unknown category: {}".format(category))
        return ""
        """Return empty , but log"""

corpus_reader = Reader(
    text_fn=process_text_line,
    get_data_fn=get_data_from_read,
    start_doc_fn=process_start_doc,
    empty_line_fn=process_empty_line,
)
