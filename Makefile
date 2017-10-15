BINARY = news2rss
IMAGENAME = laputa/news2rss
NAME = news2rss
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
	-ssh vultr docker stop $(NAME)
	-ssh vultr docker rm $(NAME)
	ssh vultr docker pull $(IMAGENAME)
	ssh vultr docker run -d --restart=always \
	  --name $(NAME) \
		-v '$HOME/feedit/config.json:/config.json' \
		-v '$HOME/feedit/data':/data \
		-v letsencrypt:/etc/letsencrypt \
		-p 25:2525 \
		-p 8100:3000 \
		$(IMAGENAME)

test-deploy:
	ssh vultr docker run -d --restart=always \
		--name $(NAME) \
		-v '$$HOME/feedit/config.json':/config.json \
		-v '$$HOME/feedit/data':/data \
		-v letsencrypt:/etc/letsencrypt \
		-p 25:2525 \
		-p 8100:3000 \
		$(IMAGENAME)
