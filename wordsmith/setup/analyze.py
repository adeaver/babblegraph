import io
from reader import Reader

word_counts = dict()
document_counts = dict()
current_document_hash_set = dict()

def process_text_line(word, _, pos):
    if pos[0] == "F":
        return
    word_counts[word] = word_counts.get(word, 0) + 1
    if word in current_document_hash_set:
        return
    current_document_hash_set[word] = True

def process_end_of_document(_):
    for key in current_document_hash_set:
        document_counts[key] = document_counts.get(key, 0) + 1
    current_document_hash_set.clear()

r = Reader(process_text_line, end_doc_fn=process_end_of_document)
r.read_data()

bad_characters = ["_", ".", ",", "\"", ":", "$", "\'"]
def has_bad_character(word):
    for char in word:
        if char in bad_characters:
            return True
    return False

with io.open("./data-dump/word_counts.csv", "w", encoding="latin1") as f:
    f.write("word,count\n")
    for word, count in word_counts.items():
        if has_bad_character(word):
            continue
        f.write("{},{}\n".format(word, str(count)))

with io.open("./data-dump/document_counts.csv", "w", encoding="latin1") as f:
    f.write("word,count\n")
    for word, count in document_counts.items():
        if has_bad_character(word):
            continue
        f.write("{},{}\n".format(word, str(count)))
