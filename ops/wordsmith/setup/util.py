import unicodedata

def make_lemma_id(lemma, part_of_speech_id):
    processed_lemma_text = _process_word_text(lemma)
    return "{}-eslm1{}".format(part_of_speech_id, processed_lemma_text)


def make_word_id(word_text, lemma_id):
    processed_word_text = _process_word_text(word_text)
    return "{}-{}".format(lemma_id, processed_word_text)


def make_part_of_speech_id(part_of_speech):
    return "espos1{}".format(part_of_speech.lower())


def make_word_ranking_id(word_ranking_text):
    processed_word_text = _process_word_text(word_ranking_text)
    return "eswr1{}".format(processed_word_text)


_characters_to_accent_numbers = {
    "á": "1",
    "â": "2",
    "ä": "3",
    "å": "4",
    "à": "5",
    "ã": "6",
    "ç": "1",
    "è": "1",
    "é": "2",
    "ê": "3",
    "ë": "4",
    "í": "1",
    "î": "2",
    "ï": "3",
    "ñ": "1",
    "ò": "1",
    "ó": "2",
    "õ": "3",
    "ö": "4",
    "ø": "5",
    "ú": "6",
    "ü": "7",
    "ý": "1",
}

def _process_word_text(word):
    out = []
    for c in word:
        lower_cased = c.lower()
        normalized_character = unicodedata.normalize('NFD', lower_cased) \
            .encode('ascii', 'ignore') \
            .decode('utf-8')
        if not len(normalized_character):
            out.append(lower_cased)
        else:
            out.append(normalized_character)
            if normalized_character != lower_cased:
                out.append(_characters_to_accent_numbers.get(lower_cased, "10"))
    return str("".join(out))
