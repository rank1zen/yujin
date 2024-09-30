.PHONY: dev

dev:
	templ generate --watch --cmd="go run ."
