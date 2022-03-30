import os, io

_START_DOC_LINE = 0
_TEXT_LINE = 1
_EMPTY_LINE = 2
_END_ARTICLE_LINE = 3
_END_DOC_LINE = 4

class Reader:
    def __init__(
        self,
        text_fn,
        get_data_fn,
        start_doc_fn=None,
        empty_line_fn=None,
        end_article_fn=None,
        end_doc_fn=None,
    ):
        self.text_fn = text_fn
        self.start_doc_fn = start_doc_fn
        self.empty_line_fn = empty_line_fn
        self.end_article_fn = end_article_fn
        self.end_doc_fn = end_doc_fn
        self.get_data_fn = get_data_fn

    def read_data(self):
        files = os.listdir("./data")
        for idx in range(len(files)):
            print("Currently on document {} of {}".format(idx+1, len(files)))
            file_name = files[idx]
            self._read_lines_in_file(file_name)
        return self.get_data_fn()

    def _read_lines_in_file(self, file_name):
        with io.open("./data/{}".format(file_name), 'r', encoding='latin1') as tagged_file:
            for line in tagged_file.readlines():
                self._process_line(line)

    def _process_line(self, line):
        line_type = self._classify_line(line)
        if line_type == _END_ARTICLE_LINE and self.end_article_fn is not None:
            self.end_article_fn(line)
        elif line_type == _END_DOC_LINE and self.end_doc_fn is not None:
            self.end_doc_fn(line)
        elif line_type == _START_DOC_LINE and self.start_doc_fn is not None:
            self.start_doc_fn(line)
        elif line_type == _EMPTY_LINE and self.empty_line_fn is not None:
            self.empty_line_fn()
        elif line_type == _TEXT_LINE:
            try:
                word, lemma, pos, _ = line.split(" ")
                self.text_fn(word.lower(), lemma.lower(), pos)
            except:
                print(f"Error on line {line}")

    def _classify_line(self, line):
        if line.startswith("ENDOFARTICLE"):
            return _END_ARTICLE_LINE
        elif line.startswith("</doc>"):
            return _END_DOC_LINE
        elif line.startswith("<doc"):
            return _START_DOC_LINE
        elif len(line.split(" ")) == 1:
            return _EMPTY_LINE
        return _TEXT_LINE
