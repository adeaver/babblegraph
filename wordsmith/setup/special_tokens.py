import uuid

class StartToken:
    def __init__(self):
        self._token_text = "<<START>>"
        self._lemma_id = uuid.uuid4()
        self._word_id = uuid.uuid4()
        self._part_of_speech_id = uuid.uuid4()

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
