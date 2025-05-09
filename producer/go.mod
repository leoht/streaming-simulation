module leohetsch.com/simulation/producer

go 1.24.2

replace leohetsch.com/simulation/simulation => ../simulation

require (
	github.com/aws/aws-msk-iam-sasl-signer-go v1.0.1
	github.com/magiconair/properties v1.8.6
	leohetsch.com/simulation/simulation v0.0.0-00010101000000-000000000000
)

require (
	github.com/google/uuid v1.6.0 // indirect
	github.com/looplab/fsm v1.0.2 // indirect
)

require (
	github.com/aws/aws-sdk-go-v2 v1.32.4 // indirect
	github.com/aws/aws-sdk-go-v2/config v1.28.2 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.17.43 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.4 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.32.4 // indirect
	github.com/aws/smithy-go v1.22.0 // indirect
	github.com/confluentinc/confluent-kafka-go/v2 v2.3.0
)
