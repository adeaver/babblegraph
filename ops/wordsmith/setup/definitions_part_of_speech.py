from typing import Optional
from util import make_part_of_speech_id

PART_OF_SPEECH_NOUN = "noun"
PART_OF_SPEECH_PREPOSITION = "preposition"
PART_OF_SPEECH_VERB = "verb"
PART_OF_SPEECH_ADJECTIVE = "adjective"
PART_OF_SPEECH_ADVERB = "adverb"
PART_OF_SPEECH_INTERJECTION = "interjection"
PART_OF_SPEECH_CONJUNCTION = "conjunction"
PART_OF_SPEECH_PRONOUN = "pronoun"
PART_OF_SPEECH_DETERMINER = "determiner"

PARTS_OF_SPEECH = [
    PART_OF_SPEECH_NOUN,
    PART_OF_SPEECH_PREPOSITION,
    PART_OF_SPEECH_VERB,
    PART_OF_SPEECH_ADJECTIVE,
    PART_OF_SPEECH_ADVERB,
    PART_OF_SPEECH_INTERJECTION,
    PART_OF_SPEECH_CONJUNCTION,
    PART_OF_SPEECH_PRONOUN,
    PART_OF_SPEECH_DETERMINER,
]

_definitions_part_of_speech_mapping = {
    PART_OF_SPEECH_NOUN: ["{m}", "{f}", "{mp}", "{mf}", "{mfp}", "{prop}", "{n}", "{fp}", "{propm}", "{propf}", "{cardinal num}", "{num}", "{numm}", "{vm}"],
    PART_OF_SPEECH_PREPOSITION: ["{prep}"],
    PART_OF_SPEECH_VERB: ["{v}", "{vr}", "{vt}", "{vi}", "{vp}", "{vtr}", "{vir}", "{vitr}", "{vrr}", "{vit}"],
    PART_OF_SPEECH_ADJECTIVE: ["{adj}", "{adjmf}", "{adjf}", "{adjm}"],
    PART_OF_SPEECH_ADVERB: ["{adv}", "{advm}"],
    PART_OF_SPEECH_INTERJECTION: ["{interj}"],
    PART_OF_SPEECH_CONJUNCTION: ["{conj}"],
    PART_OF_SPEECH_PRONOUN: ["{pron}"],
    PART_OF_SPEECH_DETERMINER: ["{determiner}", "{art}"],
}

_wordsmith_part_of_speech_category_mapping = {
    PART_OF_SPEECH_NOUN: "N",
    PART_OF_SPEECH_PREPOSITION: "S",
    PART_OF_SPEECH_VERB: "V",
    PART_OF_SPEECH_ADJECTIVE: "A",
    PART_OF_SPEECH_ADVERB: "R",
    PART_OF_SPEECH_INTERJECTION: "I",
    PART_OF_SPEECH_CONJUNCTION: "C",
    PART_OF_SPEECH_PRONOUN: "P",
    PART_OF_SPEECH_DETERMINER: "D",
}

def get_wordsmith_part_of_speech_id(definitions_part_of_speech) -> Optional[str]:
    category = _find_wordsmith_part_of_speech_category(definitions_part_of_speech)
    if category is None:
        return None
    return make_part_of_speech_id(category)

def _find_wordsmith_part_of_speech_category(definitions_part_of_speech) -> Optional[str]:
    for part_of_speech in PARTS_OF_SPEECH:
        if definitions_part_of_speech in _definitions_part_of_speech_mapping[part_of_speech]:
            return _wordsmith_part_of_speech_category_mapping[part_of_speech]
    return None
