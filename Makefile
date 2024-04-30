all:
	go build

run:
	go run .

test:
	go test .

upload:
	rsync --archive markdown-gopher sibirocobombus:/home/alex/bin
