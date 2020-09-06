import os, time, io
import psycopg2

lemma_forms = dict()
word_forms = dict()

print("connecting to database...")
pg_user = os.environ["PG_USER"]
pg_host = os.environ["PG_HOST"]
pg_password = os.environ["PG_PASSWORD"]
pg_dbname = os.environ["PG_DB"]

for i in range(10):
    try:
        print("trying to connect to database...")
        conn = psycopg2.connect("dbname='{}' user='{}' host='{}' password='{}'".format(pg_dbname, pg_user, pg_host, pg_password))
    except:
        if i < 9:
            time.sleep(2)
            continue
        raise Exception("Cannot connect to database")

files = os.listdir("./data")
for idx in range(len(files)):
    print("Currently on document {} of {}".format(idx+1, len(files)))
    file_name = files[idx]
    tagged_file = io.open("./data/{}".format(file_name), 'r', encoding='latin1')
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
    return {"lemma": lemma, "part_of_speech": pos}

print("inserting lemmas...")
cur = conn.cursor()
cur.executemany("""INSERT INTO lemmas(lemma, part_of_speech) VALUES (%(lemma)s, %(part_of_speech)s)""", [lemma_to_dict(lemma_key) for lemma_key in lemma_forms.keys()])

print("pulling lemma ids...")
cur = conn.cursor()
cur.execute("""SELECT lemma, pos, _id FROM lemmas""")
lemma_ids = { "{}::{}".format(row[0], row[1]): row[2] for row in cur.fetchall() }

def word_to_dict(word_key, lemma_id):
    word, pos = word_key.split("::")
    return { "word": word, "part_of_speech": pos, "lemma_id": lemma_id }

print("inserting words...")
words_to_insert = []
for word_key, lemma_key in word_forms.items():
    lemma_id = lemma_ids[lemma_key]
    words_to_insert.append(word_to_dict(word_key, lemma_id))

cur = conn.cursor()
cur.executemany("""INSERT INTO words(word, part_of_speech, lemma_id) VALUES (%(word)s, %(part_of_speech)s, %(lemma_id)s)""", words_to_insert)

print("done!")
