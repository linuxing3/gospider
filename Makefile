GOCMD = go
GOBUILD = $(GOCMD) build
GOMOD = $(GOCMD) mod
GOTEST = $(GOCMD) test
BINARY_NAME = goadmin
CLI = adm

all: serve

# adm
adm-install:
	go install github.com/GoAdminGroup/go-admin/adm

adm-generate:
	$(CLI) generate -c adm.ini

# go admin
init-project:
	$(CLI) init web -l cn

init-mod:
	$(GOMOD) init $(module)

install-mod:
	$(GOMOD) tidy

serve:
	$(GOCMD) run .

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ./build/$(BINARY_NAME) -v ./

test: black-box-test user-acceptance-test

black-box-test: ready-for-data
	$(GOTEST) -v -test.run=TestMainBlackBox
	make clean

user-acceptance-test: ready-for-data
	$(GOTEST) -v -test.run=TestMainUserAcceptance
	make clean

ready-for-data:
	cp admin.db admin_test.db

clean:
	rm admin_test.db

# spider
install-spider:
	$(GOCMD) install .

run-spider:
	$(GOCMD) run .

# prisma
prisma-introspect:
	go run github.com/prisma/prisma-client-go introspect

prisma-push:
	go run github.com/prisma/prisma-client-go db push --preview-feature

prisma-generate:
	go run github.com/prisma/prisma-client-go generate

# docker
docker-up:
	# bash ./run.sh
	docker-compose up -d --remove-orphans

db-init:
	deno run -A --unstable ./create_table.ts

spider-douban:
	docker exec -it spider /root/go/bin/gospider spider

.PHONY: all install-cli init-project serve build generate test black-box-test user-acceptance-test ready-for-data clean