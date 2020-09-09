.PHONY: debug
debug:
	go build -o cli
	./cli --content=/Users/dave/Development/go/ons/dp-zebedee-content/testing --zebedee=/Users/dave/Development/go/ons/dp-zebedee-content/testing