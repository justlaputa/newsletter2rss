BINARY = news2rss
IMAGENAME = laputa/news2rss
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: clean docker-image build deploy

$(BINARY): $(SRC)
	env GOOS=linux GOARCH=amd64 go build -o $(BINARY)

clean:
	$(RM) $(BINARY)

docker-image: $(BINARY)
	docker build -t $(IMAGENAME) .

build: docker-image

deploy: build
	docker push $(IMAGENAME)
