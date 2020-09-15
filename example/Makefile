prepare:
	go build -o genembed ../genembed/genembed.go
gen: prepare
	find ./* -name '*_genembed.go' -print0 | xargs -0 rm
	PATH=${PATH}:${PWD} go generate ./...

run:
	go run main.go main_genembed.go
