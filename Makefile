.PHONY: run
default: run ;

run:
	@docker build -t collector . && docker run collector
