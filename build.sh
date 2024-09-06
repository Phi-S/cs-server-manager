set -e
npm install --prefix frontend/
npm run build --prefix frontend/

cp -R frontend/dist/* backend/web

swag init --dir backend -o . -ot json
cp swagger-ui/* backend/swagger-ui
cp swagger.json backend/swagger-ui/swagger.json

go mod download -C backend/
go build -C backend/ -v -o ../cs-server-manager