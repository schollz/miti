build:
	go build -v


release:
	rm -f miti 
	rm -rf dist 
	go generate
	git push
	goreleaser release
