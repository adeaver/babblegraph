import math, json, os
from corpus_read_defs import corpus_reader, handle_end_of_doc

corpus_reader.end_corpus_doc_fn = handle_end_of_doc
corpus_reader.read_data()

file_prefixes = ["observed_lemmas_", "lemma_counts_", "observed_words_", "observed_parts_of_speech_", "bigram_counts_"]

files_by_prefix = {}

for file in os.listdir("/out"):
    for prefix in file_prefixes:
        if file.startswith(prefix):
            files = files_by_prefix.get(prefix, [])
            files.append(file)
            files_by_prefix[prefix] = files
            break

CHUNK_SIZE = 2.0

number_of_files = sum([len(files) for files in files_by_prefix.values()])
while number_of_files > len(file_prefixes):
    print(files_by_prefix)
    for prefix, files in files_by_prefix.items():
        files_chunked = [ files[int(i*CHUNK_SIZE):min(len(files), int(i*CHUNK_SIZE + CHUNK_SIZE))] for i in range(int(math.ceil(len(files) / CHUNK_SIZE))) ]
        chunks = []
        for chunk in files_chunked:
            data = None
            for file in chunk:
                with open(f"/out/{file}", "r") as f:
                    file_data = json.loads(f.read())
                    if data is None:
                        data = file_data
                    elif isinstance(file_data, list):
                        data += file_data
                    elif isinstance(file_data, dict):
                        if "counts" in prefix:
                            for key, value in file_data.items():
                                data[key] = data.get(key, 0) + int(value)
                        else:
                            data.update(file_data)
            file_name = f"{prefix}_chunked_{len(chunks)}.json"
            print(f"Writing {file_name}")
            with open(f"/out/{file_name}", "w") as f:
                json.dump(data, f)
            chunks.append(file_name)
        files_by_prefix[prefix] = chunks
    number_of_files = sum([len(files) for files in files_by_prefix.values()])


