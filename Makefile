.PHONY: install
install:
	go install github.com/ONSdigital/dp-zebedee-content

.PHONY: debug
debug:
	go build -o cli
	./cli generate -c=<WHERE_YOU_WANT_THE_CONTENT_TO_BE_CREATED>
