if [ -d "./out" ]; then
    rm -rf ./out
fi
mkdir ./out
docker build -f Dockerfile -t initial-data-csvs .
docker run -v $(pwd)/out/:/out initial-data-csvs
