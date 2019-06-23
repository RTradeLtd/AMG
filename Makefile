
.PHONY: amg-bindings
amg-bindings:
	abigen -sol AMG.sol -pkg bindings -out ./bindings/bindings.go


.PHONY: tests
tests:
	go test -race -cover ./...