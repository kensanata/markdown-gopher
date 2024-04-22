all:
	go build

run:
	go run .

upload:
	rsync --archive markdown-gopher sibirocobombus:/home/alex/bin
