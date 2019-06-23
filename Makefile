
.PHONY: amg-bindings
amg-bindings:
	abigen -sol AMG.sol -pkg bindings -out ./bindings/bindings.go


.PHONY: compile
compile:
	solc --bin --abi AMG.sol -o ./bin --overwrite

.PHONY: tests
tests:
	go test -v -race -cover ./...