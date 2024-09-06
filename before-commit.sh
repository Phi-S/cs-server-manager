set -e
swag init --dir backend -o . -ot json
npx widdershins -v --code --summary --expandBody --omitHeader -o api-documentation.md swagger.json