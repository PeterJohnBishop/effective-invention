# effective-invention

# notes

deploying postgres on Docker:

docker run --name go-postgres \
    -e POSTGRES_PASSWORD=postgres \
    -p 5432:5432 \
    -d postgres
