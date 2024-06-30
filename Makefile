templ:
	@templ generate

tailwind:
	@npx tailwindcss -i tailwind.css -o ./static/css/styles.css

.PHONY: live/templ
live/templ:
	templ generate --watch --proxy="http://localhost:8080" --open-browser=true

.PHONY: live/server
live/server:
	go run github.com/cosmtrek/air@v1.51.0 \
		--build.cmd "go build -o tmp/bin/yujin ./cmd/yujintest" \
		--build.bin "./tmp/bin/yujin" \
		--build.delay "100" \
		--build.exclude_dir "node_modules" \
		--build.include_ext "go" \
		--build.stop_on_error "false" \
		--misc.clean_on_exit true

.PHONY: live/tailwind
live/tailwind:
	npx tailwindcss -i ./tailwind.css -o ./static/css/styles.css --watch

.PHONY: live/sync_assets
live/sync_assets:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "static" \
	--build.include_ext "js,css"

.PHONY: live
live:
	make -j4 live/templ live/server live/tailwind live/sync_assets
