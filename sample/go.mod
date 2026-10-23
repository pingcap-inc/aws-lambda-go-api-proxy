module github.com/awslabs/aws-lambda-go-api-proxy-sample

go 1.12

require (
	github.com/aws/aws-lambda-go v1.10.0
	github.com/awslabs/aws-lambda-go-api-proxy v0.3.0
	github.com/gin-gonic/gin v1.9.1
	github.com/google/uuid v1.1.1
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.19.0 // indirect
)

replace (
	gopkg.in/yaml.v2 v2.2.2 => gopkg.in/yaml.v2 v2.2.8
	gopkg.in/yaml.v2 v2.2.4 => gopkg.in/yaml.v2 v2.2.8
)
