build:
	make tailwind-build && make templ-generate && go build -o ./bin/$(APP_NAME) ./cmd/$(APP_NAME)/main.go

run:
	make build && go run ./bin/main

dev:
	air

tailwind-build:
	npx tailwindcss -c ./configs/tailwind.config.js -i ./configs/input.css -o static/css/style.css --minify

tailwind-watch:
	npx tailwindcss -c ./configs/tailwind.config.js -i ./configs/input.css -o static/css/style.css --watch

templ-generate:
	templ generate

templ-watch:
	templ generate --watch

docker-build:
	docker build --no-cache -f Dockerfile -t crowmw/risiti:latest . 

docker-run:
	docker run -d -e SECRET="secretKey" -p 80:80 -v ${HOME}/data:/data crowmw/risiti:latest

docker-push:
	docker push crowmw/risiti