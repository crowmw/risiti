build:
	make tailwind-build && make templ-generate && go build -o ./bin/$(APP_NAME) ./cmd/$(APP_NAME)/main.go

dev:
	make tailwind-build && make templ-generate && go run cmd/main.go

tailwind-build:
	npx tailwindcss -c ./configs/tailwind.config.js -i ./configs/input.css -o public/css/style.css --minify

tailwind-watch:
	npx tailwindcss -c ./configs/tailwind.config.js -i ./configs/input.css -o public/css/style.css --watch

templ-generate:
	templ generate

templ-watch:
	templ generate --watch

docker-build:
	docker build --no-cache -f Dockerfile -t crowmw/risiti:latest . 

docker-run:
	docker run -d -p 2137:2137 -v ${HOME}/data:/data crowmw/risiti:latest