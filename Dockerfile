# Build tailwindcss
FROM node:21 as tailwind-builder
WORKDIR /build

# Import necessary codebase
COPY /internal/components ./internal/components
COPY /configs ./configs

# Generate tailwindcss output
RUN npx tailwindcss -c configs/tailwind.config.js -i configs/input.css -o style.css --minify

# Build app
FROM golang:1.22.1-bookworm as builder
WORKDIR /build

# Copy code files
COPY . .

# Install templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Generate templates
RUN templ generate

# Build go app
RUN go build -a -ldflags '-w -extldflags "-static"' -o ./bin/risiti ./cmd/main.go & wait

# Run app
FROM scratch AS runner
WORKDIR /bin

# Server binary from builder
COPY --from=builder /build/bin/risiti ./risiti

# Copy public files
COPY /public ./public

# Copy data catalog for files
COPY /data ./data

# Copy built tailwindcss styles
COPY --from=tailwind-builder /build/style.css ./public/css

EXPOSE 2137

# Run the server
ENTRYPOINT ["/bin/risiti"]