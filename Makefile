.PHONY: install
install:
	go install github.com/ONSdigital/dp-zebedee-content

.PHONY: debug
debug:
	go build -o cli
	./cli generate -c=<REPLACE_WITH_CONTENT_PATH> -z=<REPLACE_WITH_ZEBEDEE_PROJECT_PATH>
