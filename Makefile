.PHONY: install
install:
	go install github.com/ONSdigital/dp-zebedee-content

.PHONY: debug
debug:
	go build -o cli
	./cli -c=/Users/dave/Development/go/ons/dp-zebedee-content/testing -z=/Users/dave/Development/go/ons/dp-zebedee-content/testing