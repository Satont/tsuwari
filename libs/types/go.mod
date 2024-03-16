module github.com/satont/twir/libs/types

go 1.21.5

require (
	github.com/tkrajina/typescriptify-golang-structs v0.1.11
	github.com/twirapp/twir/libs/bus-core v0.0.0-20240225024146-742838c78cea
)

require (
	github.com/stretchr/testify v1.9.0 // indirect
	github.com/tkrajina/go-reflector v0.5.6 // indirect
)

replace (
	github.com/satont/twir/libs/gomodels => ../gomodels
	github.com/twirapp/twir/libs/bus-core => ../bus-core
)
