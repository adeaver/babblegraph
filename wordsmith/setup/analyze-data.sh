if [ -d "./data-dump" ]; then
    rm -rf ./data-dump
fi
mkdir ./data-dump
docker build -f Dockerfile -t analyze-words .
docker run -v $(pwd)/data-dump/:/data-dump analyze-words /bin/bash -c "python -u analyze.py"
