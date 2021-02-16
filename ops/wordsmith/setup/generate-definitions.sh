#!/bin/bash
set -euo pipefail

if [ -d "./out/definitions*" ]; then
    rm ./out/definitions*
fi

docker build -f Dockerfile.definitions -t definitions-data .
docker run -v $(pwd)/out/:/out definitions-data
