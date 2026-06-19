templ=$(HOME)/go/bin/templ

main: main.go style/ templates/ public/
	go install
	$(templ) generate
	go run main.go builders.go build_projects.go
	cd build && python -m http.server 8765

build_dev: main.go style/ templates/ public/
	$(templ) generate
	go run main.go builders.go build_projects.go --dev
	
dev: main.go style/ templates/ public/
	go install
	python devserver.py --port 8765 --build-cmd "make build_dev"
