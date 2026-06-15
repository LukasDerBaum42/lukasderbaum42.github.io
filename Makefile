main: main.go style/ templates/ public/
	templ generate
	go run main.go build_projects.go
	cd build && python -m http.server