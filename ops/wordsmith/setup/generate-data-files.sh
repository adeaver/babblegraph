#!/bin/bash
docker build -f Dockerfile -t initial-data .
docker run -v $(pwd)/out/:/out initial-data /bin/bash -c "python -u process2.py"
