main: main.go style/ templates/
	templ generate
	go run main.go
	cd public && python -m http.server