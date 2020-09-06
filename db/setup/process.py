import os, time, io, uuid

lemma_forms = dict()
word_forms = dict()

files = os.listdir("./data")
for idx in range(len(files)):
    print("Currently on document {} of {}".format(idx+1, len(files)))
    file_name = files[idx]
    with io.open("./data/{}".format(file_name), 'r', encoding='latin1') as tagged_file:
        for line in tagged_file.readlines():
            parts = line.split(" ")
            if len(parts) != 4:
                continue
            word, lemma, pos, _ = parts
            if pos[0] == "F":
                continue
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

def lemma_to_dict(lemma_key):
    lemma, pos = lemma_key.split("::")
    return (lemma, pos, "es")

def word_to_dict(word_key):
    word, pos = word_key.split("::")
    return  (word, pos, "es")

print("writing lemma csv")
lemma_key_to_id = dict()
with io.open("./out/lemmas.csv", "w", encoding="latin1") as f:
    f.write("_id,lemma,part_of_speech,language\n")
    for lemma_key in lemma_forms.keys():
        lemma, part_of_speech, language = lemma_to_dict(lemma_key)
        _id = uuid.uuid4()
        line = "{},{},{},{}\n".format(_id, lemma, part_of_speech, language)
        f.write(line)
        lemma_key_to_id[lemma_key] = _id
    del lemma_forms

print("writing word csv")
with io.open("./out/words.csv", "w", encoding="latin1") as f:
    f.write("lemma_id,word,part_of_speech,language\n")
    for word_key, lemma_key in word_forms.items():
        word, part_of_speech, language = word_to_dict(word_key)
        lemma_id = lemma_key_to_id[lemma_key]
        line = "{},{},{},{}\n".format(lemma_id, word, part_of_speech, language)
        f.write(line)

print("done!")
