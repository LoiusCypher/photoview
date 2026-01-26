docker compose -f dev-compose.yaml run --remove-orphans ui bash -c 'npm install --global npm@10.8.2 --no-audit --no-fund ; npm install --no-audit --no-fund ; npm run genSchemaTypes ; npm ci --no-audit --no-fund ; npm run build:dev --no-audit --no-fund -- --base="/" '
docker compose -f dev-compose.yaml run --remove-orphans api bash -c 'go get github.com/loiuscypher/photoview/api/graphql ; go get github.com/loiuscypher/photoview/api/scanner ; go generate'
# sudo docker compose -f dev-compose.yaml --profile mysql up --remove-orphans
