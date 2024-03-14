build:
	make tailwind-build && make templ-generate && go build -o ./bin/$(APP_NAME) ./cmd/$(APP_NAME)/main.go

dev:
	make tailwind-build && make templ-generate && go run cmd/main.go

tailwind-build:
	npx tailwindcss --config ./configs/tailwind.config.js -i configs/input.css -o static/css/style.css --minify

tailwind-watch:
	npx tailwindcss --config ./configs/tailwind.config.js -i configs/input.css -o static/css/style.css --watch

templ-generate:
	templ generate

templ-watch:
	templ generate --watch

docker-build:
	docker build --no-cache --progress=plain -f Dockerfile -t risiti:latest . 

docker-run:
	docker run -d -p 2137:2137 risiti:latest