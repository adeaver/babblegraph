START_TOKEN = "<<START>>"

part_of_speech_trigrams = [START_TOKEN, START_TOKEN, START_TOKEN]
part_of_speech_trigram_counts = dict()

def process_text_line(word, lemma, pos):
    # skip punctuation
    if pos[0] == "F":
        return
    processed_part_of_speech = _process_part_of_speech(pos)
    _handle_part_of_speech(processed_part_of_speech)

def _handle_part_of_speech(processed_part_of_speech):
    part_of_speech_trigrams.pop(0)
    part_of_speech_trigrams.append(processed_part_of_speech)
    trigram_key = "{}, {}, {}".format(part_of_speech_trigrams[0], part_of_speech_trigrams[1], part_of_speech_trigrams[2])
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
    elif category == "Z": # Number
        """Returns only the category"""
        return pos[0]
    elif category == "W": # Date
        """Returns only the category"""
        return pos[0]
    elif category == "I": # Interjection
        """Returns only the category"""
        return pos[0]
    else:
        print("Unknown category: {}".format(category))
        return ""
        """Return empty , but log"""


