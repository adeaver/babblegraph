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

def _process_word_text(word):
    out = []
    for c in word:
        lower_cased = c.lower()
        normalized_character = unicodedata.normalize('NFD', lower_cased) \
            .encode('ascii', 'ignore') \
            .decode('utf-8')
        out.append(normalized_character)
        if normalized_character != lower_cased:
            out.append("1")
    return str("".join(out))
