import io

POPULATE_TEMPLATE = """
COPY public.\"{}\"({})
FROM '{}'
DELIMITER ','
CSV HEADER
ENCODING 'latin1';\n\n
"""

def _write_lines_to_file(file_name_prefix, lines_per_file, line_generator, title_line=None, file_type="csv"):
    file_count = 1
    make_file_name = lambda: "/out/{}-{}.{}".format(file_name_prefix, file_count, file_type)
    file_names_written = []
    lines_to_write = []
    for line in line_generator():
        lines_to_write.append(line)
        if len(lines_to_write) == lines_per_file:
            file_name = make_file_name()
            _write_file(file_name, lines_to_write, title_line)
            file_names_written.append(file_name)
            file_count += 1
            lines_to_write = []
    if len(lines_to_write) > 0:
        file_name = make_file_name()
        _write_file(file_name, lines_to_write, title_line)
        file_names_written.append(file_name)
    return file_names_written

def _write_file(file_name, lines, title_line):
    with io.open(file_name, "w", encoding="latin1") as f:
        if title_line is not None:
            f.write("{}\n".format(title_line))
        for line in lines:
            f.write("{}\n".format(line))
    print("wrote file {}".format(file_name))

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
            yield POPULATE_TEMPLATE.format(self.table_name, self.column_headers, f)
