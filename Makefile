ifndef PACKAGECLOUD_TOKEN
  $(error PACKAGECLOUD_TOKEN is not set)
endif

test:
	go test -v ./...

generate: assets/distributions.json
	go generate -x ./...

assets/distributions.json:
	curl -L https://$(PACKAGECLOUD_TOKEN):@packagecloud.io/api/v1/distributions.json | jq . >$@

.PHONY: assets/distributions.json
