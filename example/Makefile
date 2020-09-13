gen:
	find ./* -name '*_embeded.go' -print0 | xargs -0 rm
	go generate ./...
	find ./* -name '*_embeded.go.tmp' -print0 | xargs -0 rm
	find ./* -name '*_embeded.go' -exec gofmt -w {} +

run:
	go run main.go main_embeded.go
