build:
	make tailwind-build && make templ-generate && go build -o ./bin/$(APP_NAME) ./cmd/$(APP_NAME)/main.go

dev:
	make tailwind-build && make templ-generate && go run cmd/main.go

tailwind-build:
	npx tailwindcss --config ./tailwind.config.js -i static/css/input.css -o static/css/style.css --minify

tailwind-watch:
	npx tailwindcss --config ./tailwind.config.js -i static/css/input.css -o static/css/style.css --watch

templ-generate:
	templ generate

templ-watch:
	templ generate --watch

