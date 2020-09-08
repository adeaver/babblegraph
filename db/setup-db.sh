if [ ! -f "./setup/out/lemmas-1.csv" ]; then
    echo "Lemmas does not exist, run generate-csvs.sh from setup dir"
    exit 1
fi

if [ ! -f "./setup/out/words-1.csv" ]; then
    echo "Words does not exist, run generate-csvs.sh from setup dir"
    exit 1
fi

echo "Removing running container"
(docker stop setup-dev-db || true) && (docker rm setup-dev-db || true)

if [ -d "./data/" ]; then
    echo "Removing data from postgres container"
    rm -rf ./data
fi

echo "Building"
docker build -f Dockerfile.dev -t dev-db .
docker run -v $(pwd)/data:/var/lib/postgresql/data -v $(pwd)/setup:/setup --name setup-dev-db -d dev-db
until (docker exec setup-dev-db psql -U dev -d babblegraph -c 'SELECT * FROM lemmas'); do
    echo "Waiting for container to be up"
    sleep 1;
done;

echo "Applying file"
docker exec setup-dev-db psql -U dev -d babblegraph -a -f /setup/out/populate_db.sql

echo "Stopping"
docker stop setup-dev-db
