build:
	go build -v


release:
	rm -f miti 
	rm -rf dist 
	VERSION=0.3.2 go generate
	git push
	goreleaser release
	