from util import (
    make_lemma_id,
    make_word_id,
    make_part_of_speech_id
)

class SpecialToken:
    def __init__(self, token_text, part_of_speech_id, lemma_id):
        self._token_text = token_text
        self._part_of_speech_id = part_of_speech_id
        self._lemma_id = lemma_id
        self._word_id = make_word_id(self._token_text, self._lemma_id, self._part_of_speech_id)

    def get_token(self):
        return self._token_text

    def get_word_key(self):
        """Word key is (word),(part of speech)

        Here the token text is both"""
        return "{},{}".format(self._token_text, self._token_text)

    def get_word_key_value(self):
        """Word key value is (word_id,lemma_key)"""
        return "{},{}".format(self._word_id, self.get_lemma_key())

    def get_lemma_key(self):
        """Lemma key is (lemma),(part of speech)

        Here the token text is both"""
        return "{},{}".format(self._token_text, self._token_text)

    def get_lemma_id(self):
        return self._lemma_id

    def get_part_of_speech_id(self):
        return self._part_of_speech_id

class StartToken(SpecialToken):
    def __init__(self):
        part_of_speech_id = make_part_of_speech_id("<<START>>")
        super().__init__(
            "<<START>>",
            part_of_speech_id,
            make_lemma_id("<<START>>", part_of_speech_id)
        )
