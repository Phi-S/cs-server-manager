npm install --prefix frontend/
npm run build --prefix frontend/

cp -r frontend/dist/* backend/web

go mod tidy -C backend/
go mod verify -C backend/
go mod download -C backend/
go build -C backend/ -v -o ../cs-server-manager