.PHONY: dep test
NAME = $(notdir $(shell pwd))


dep:
	go get -v && go mod tidy
test:
	go test -v $(NAME) -count=1 -run=$(run)

