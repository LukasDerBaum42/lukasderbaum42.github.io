templ=$(HOME)/go/bin/templ

main: main.go style/ templates/ public/
	# go install
	$(templ) generate
	go run .
	cd build && python -m http.server 8765

build_dev: main.go style/ templates/ public/
	# go install
	$(templ) generate
	go run . --dev
	
dev: main.go style/ templates/ public/
	# go install
	python devserver.py --port 8765 --build-cmd "make build_dev"
