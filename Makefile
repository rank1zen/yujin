APP := yujin
BUILD_DATE := `date +%FT%T%z`

STATIC_DIR := internal/ui/static
CSS_SRC := internal/ui/static/css/tailwind.css
CSS_OUT := internal/ui/static/css/styles.css

# .PHONY: \
# 	yujin \
# 	live \
# 	live/templ \
# 	live/server \
# 	live/tailwind \
# 	live/sync_assets

yujin: build-templ build-tailwind
	@ go build -o $(APP) main.go

build-templ:
	@ templ generate

build-tailwind:
	@ npx tailwindcss -i $(STATIC_DIR)/css/tailwind.css -o $(STATIC_DIR)/css/styles.css --minify

dev:
	wgo -file=.go -file=.templ -xfile=_templ.go \
		templ generate :: npx tailwindcss -i $(STATIC_DIR)/css/tailwind.css -o $(STATIC_DIR)/css/styles.css --minify :: go run .
