module github.com/nilorg/naas/examples/grpc-casbin-adapter

go 1.14

require (
	github.com/casbin/casbin/v2 v2.17.1
	github.com/nilorg/istio v0.0.0-20200912054551-2792f5b24cc3
	github.com/nilorg/naas v0.0.0
	github.com/nilorg/sdk v0.0.0-20200912025101-a4037e6ee224
)

replace github.com/nilorg/naas => ../../
