.DEFAULT_GOAL := roboheart

.PHONY: prepare
prepare:
	mkdir -p output
	go mod tidy
	go mod vendor

.PHONY: roboheart
roboheart: prepare
	go build -o output/roboheart cmd/roboheart/*.go

.PHONY: roboheart-txt
roboheart-txt: prepare
	GOOS=linux GOARCH=arm go build -o output/roboheart-txt cmd/roboheart/*.go

.PHONY: clean
clean:
	rm -rf output

.PHONY: checksvcs
checksvcs:
	go run cmd/checksvcs/*.go

.PHONY: prepare-services
prepare-services:
	go run cmd/setservices/*.go -a
	$(MAKE) checkdeps
