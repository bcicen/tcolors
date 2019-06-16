NAME=tcolors
VERSION=$(shell cat VERSION)
LDFLAGS="-w -X main.version=$(VERSION) -X main.build=$(shell git log -1 --format=%cd.%h --date=short|tr -d -)"

clean:
	rm -rvf build/ release/ arch-release/

deps:
	go mod download

build: deps
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -o $(NAME)

build-all: deps
	mkdir -p build
	GOOS=darwin GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/$(NAME)-$(VERSION)-darwin-amd64
	GOOS=linux GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/$(NAME)-$(VERSION)-linux-amd64
	GOOS=linux GOARCH=arm go build -ldflags $(LDFLAGS) -o build/$(NAME)-$(VERSION)-linux-arm
	GOOS=freebsd GOARCH=amd64 go build -ldflags $(LDFLAGS) -o build/$(NAME)-$(VERSION)-freebsd-amd64

release:
	mkdir release
	go get github.com/progrium/gh-release/...
	cp build/* release
	gh-release create bcicen/$(NAME) $(VERSION) \
		$(shell git rev-parse --abbrev-ref HEAD) $(VERSION)

arch-release:
	mkdir -p arch-release
	go get github.com/seletskiy/go-makepkg/...
	cd arch-release && \
		go-makepkg -p version "Commandline color picker and palette builder" git://github.com/bcicen/tcolors.git; \
		git clone ssh://aur@aur.archlinux.org/tcolors.git; \
		cp build/* tcolors/
	cd arch-release/tcolors/ && \
		mksrcinfo
