import io

POPULATE_TEMPLATE = """
COPY public.\"{}\"({})
FROM '{}'
DELIMITER ','
CSV HEADER
ENCODING 'latin1';\n\n
"""

def _write_lines_to_file(file_name_prefix, lines_per_file, line_generator, title_line=None, file_type="csv"):
    out = []
    file_count = 1
    make_file_name = lambda: "./out/{}-{}.{}".format(file_name_prefix, file_count, file_type)
    counter = _write_file(make_file_name(), line_generator, lines_per_file, title_line)
    out.append(make_file_name())
    file_count += 1
    while counter == lines_per_file:
        counter = _write_file(make_file_name(), line_generator, lines_per_file, title_line)
        out.append(make_file_name())
        file_count += 1
    return out

def _write_file(file_name, line_generator, max_lines, title_line):
    counter = 0
    with io.open(file_name, "w", encoding="latin1") as f:
        if title_line is not None:
            f.write("{}\n".format(title_line))
        for line in line_generator():
            f.write("{}\n".format(line))
            counter += 1
            if counter == max_lines:
                break
    return counter

class SQLTemplate:
    def __init__(self, table_name, column_headers):
        self.table_name = table_name
        self.column_headers = column_headers
        self._files = []

    def write_files_for_template(self, line_generator, max_lines):
        self._files = _write_lines_to_file(self.table_name, max_lines, line_generator, self.column_headers)

    def yield_sql_template(self):
        if len(self._files) == 0:
            raise Exception("Must have written files")
        for f in self._files:
            yield POPULATE_TEMPLATE.format(self.table_name, self.column_names, f)
