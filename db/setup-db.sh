if [ ! -f "./setup/out/lemmas.csv" ]; then
    echo "Lemmas does not exist, run `generate-csvs.sh` from setup dir"
    exit 1
fi

if [ ! -f "./setup/out/words.csv" ]; then
    echo "Words does not exist, run `generate-csvs.sh` from setup dir"
    exit 1
fi

echo "Removing"
(docker stop setup-dev-db || true) && (docker rm setup-dev-db || true)

echo "Building"
docker build -f Dockerfile.dev -t dev-db .
docker run -v $(pwd)/data:/var/lib/postgresql/data -v $(pwd)/setup:/setup --name setup-dev-db -d dev-db
until (docker exec setup-dev-db psql -U dev -d babblegraph -a -f /setup/populate_word_data.sql); do
    echo "Waiting for container to be up"
    sleep 0.1;
done;
echo "Stopping"
docker stop setup-dev-db
