.PHONY: run
default: run ;

run:
	@docker build -t collector . && docker run -e GOOGLE_APPLICATION_CREDENTIALS=/go/src/github.com/DanShu93/vocabulary-collector/google.json collector > output/vocabulary.json
