
# Prerequisites:
# npm
# go-bindata
# go-bindata-assetfs
# go

systems := darwin windows linux
architectures := 386 amd64

local_goos = $(shell go env GOOS)
local_goarch := $(shell go env GOARCH) 
.PHONY: install-frontend-dependencies install-backend-dependencies build-frontend encode-frontend build-backend build-total
install-frontend-dependencies:
	cd frontend; npm install

install-backend-dependencies:
	cd backend; dep ensure -vendor-only

build-frontend:
	rm -rf frontend/out/
	mkdir -p frontend/out/
	cd frontend; npm run build
	cp frontend/public/* frontend/out/

encode-frontend:
	rm -f backend/src/handlers/http/bindata_assetfs.go
	go-bindata-assetfs -pkg http -prefix ./frontend/ ./frontend/out/
	mv bindata_assetfs.go backend/src/handlers/http

build-backend:
	rm -rf backend/out/
	mkdir -p backend/out/
	CGO_ENABLED=0 GOOS=$(local_goos) GOARCH=$(local_goarch) go build -o ./backend/out/elevator-simulator-$(local_goos)-$(local_goarch) backend/src/main.go; \
	
build-backend-all:
	rm -rf backend/out/
	mkdir -p backend/out/
	for GOOS in $(systems); do \
		for GOARCH in $(architectures); do \
			CGO_ENABLED=0 GOOS=$$GOOS GOARCH=$$GOARCH go build -o ./backend/out/elevator-simulator-$$GOOS-$$GOARCH backend/src/main.go; \
		done \
	done

build-total: build-frontend encode-frontend build-backend

build-total-all: build-frontend encode-frontend build-backend-all