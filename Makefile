PKG := "github.com/GeekQk/devcloud-mini"

gen: ## Gen code
	@protoc -I=. --go_out=. --go-grpc_out=. --go-grpc_opt=module="github.com/GeekQk/devcloud-mini" --go_opt=module="github.com/GeekQk/devcloud-mini" cmdb/apps/*/pb/*.proto
	@protoc-go-inject-tag -input="cmdb/apps/*/*.pb.go"

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'