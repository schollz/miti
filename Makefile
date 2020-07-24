build:
	go build -v


release:
	rm -f miti 
	rm -rf dist 
	go generate
	git push
	goreleaser release

linux:
	go build -v -o miti
	zip miti_linux_amd64.zip miti README.md LICENSE 
	.github/uploadrelease.sh miti_linux_amd64.zip
