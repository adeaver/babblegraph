import os, time, io, uuid
from reader import Reader


lemma_forms = dict()
word_forms = dict()

def process_line(word, lemma, pos):
    if pos[0] == "F":
        return
    lemma_key = "{}::{}".format(lemma, pos[0])
    word_key = "{}::{}".format(word, pos)
    current_word_entry = word_forms.get(word_key, None)
    lemma_forms[lemma_key] = True
    if current_word_entry != None:
        if current_word_entry != lemma_key:
            print("Error on: {}".format(word_key))
            print("Got non-duplicate: {} and {}".format(current_word_entry, lemma_key))
        continue
    word_forms[word_key] = lemma_key

r = Reader(process_line)
r.read_data()

def lemma_to_dict(lemma_key):
    lemma, pos = lemma_key.split("::")
    return (lemma, pos, "es")

def word_to_dict(word_key):
    word, pos = word_key.split("::")
    return  (word, pos, "es")

bad_characters = ["_", ".", ",", "\"", ":", "$", "\'"]
def has_bad_character(word):
    for char in word:
        if char in bad_characters:
            return True
    return False

print("extract lines for lemma csvs")
lemma_key_to_id = dict()
file_lines = []
current_chunk = []
chunk_size = 250000
for lemma_key in lemma_forms.keys():
    lemma, part_of_speech, language = lemma_to_dict(lemma_key)
    _id = uuid.uuid4()
    if has_bad_character(lemma) or has_bad_character(part_of_speech):
        continue
    line = "{},{},{},{}".format(_id, lemma, part_of_speech, language)
    lemma_key_to_id[lemma_key] = _id
    current_chunk.append(line)
    if len(current_chunk) == chunk_size:
        file_lines.append(current_chunk)
        current_chunk = []
if len(current_chunk) > 0:
    file_lines.append(current_chunk)
del lemma_forms
del current_chunk

print("writing lines to csvs")
num_lemma_files = len(file_lines)
for idx in range(num_lemma_files):
    file_name = "./out/lemmas-{}.csv".format(idx+1)
    with io.open(file_name, "w", encoding="latin1") as f:
        f.write("_id,lemma,part_of_speech,language\n")
        f.write("\n".join(file_lines[idx]))
del file_lines

print("getting word lines")
word_lines = []
word_chunk = []
for word_key, lemma_key in word_forms.items():
    word, part_of_speech, language = word_to_dict(word_key)
    lemma_id = lemma_key_to_id.get(lemma_key, None)
    if has_bad_character(word) or has_bad_character(part_of_speech) or lemma_id is None:
        continue
    line = "{},{},{},{}".format(lemma_id, word, part_of_speech, language)
    word_chunk.append(line)
    if len(word_chunk) == chunk_size:
        word_lines.append(word_chunk)
        word_chunk = []
if len(word_chunk) > 0:
    word_lines.append(word_chunk)

print("writing word csv")
num_word_files = len(word_lines)
for idx in range(num_word_files):
    file_name = "./out/words-{}.csv".format(idx+1)
    with io.open(file_name, "w", encoding="latin1") as f:
        f.write("lemma_id,word,part_of_speech,language\n")
        f.write("\n".join(word_lines[idx]))

print("writing sql file")
with io.open("./out/populate_db.sql", "w", encoding="latin1") as f:
    f.write("INSERT INTO languages (code) VALUES ('es');\n\n")
    for idx in range(num_lemma_files):
        f.write(populate_lemmas_template.format(idx+1))
    for idx in range(num_word_files):
        f.write(populate_words_template.format(idx+1))

print("done!")
