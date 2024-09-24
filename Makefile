.SHELLFLAGS := -ec

test:
	go test -C backend/ -v ./...

doc:
	swag init --dir backend -o . -ot json
	npx widdershins -v --code --summary --expandBody --omitHeader -o api-documentation.md swagger.json

ready: tidy test doc clear

tidy:
	go mod tidy -C backend/

backend:
	go run -C backend/ .

frontend:
	npm run dev --prefix frontend/

build:
	npm install --prefix frontend/
	npm run build --prefix frontend/

	cp -R frontend/dist/* backend/web

	swag init --dir backend -o . -ot json
	cp swagger-ui/* backend/swagger-ui
	cp swagger.json backend/swagger-ui/swagger.json

	go mod download -C backend/
	go build -C backend/ -v -o ../cs-server-manager

clear:
	rm -f cs-server-manager
	rm -rf frontend/dist/*
	rm -rf backend/web/*
	rm -rf backend/swagger-ui/*
	git restore --staged --worktree backend/swagger-ui/index.html
	git restore --staged --worktree backend/web/index.html

.PHONY: test doc ready tidy backend frontend build clear