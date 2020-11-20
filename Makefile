build-mac:
	rm -f bin/*
	cd app && env GOOS=darwin GOARCH=amd64 go build -o ../bin/main main.go

# NOT TESTED
build-windows:
	rm -f bin/*
	cd app && env GOOS=windows GOARCH=amd64 go build -o bin/main main.go

run:
	bin/main