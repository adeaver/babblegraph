#!/bin/bash
set -euo pipefail

if [ -d "./out/phrase-definitions*" ]; then
    rm ./out/phrase-definitions*
fi

docker build -f Dockerfile.definitions2 -t definitions-data .
docker run -v $(pwd)/out/:/out definitions-data
