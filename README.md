# Streaming simulation

This repository contains code and configuration to deploy a lightweight simulation of data streamed into a Kafka cluster on AWS MSK. It uses:
 - Go
 - PostgreSQL
 - Docker (to run locally)
 - Terraform
 - AWS

## Setup

Setup the AWS infrastructure using Terraform:

```
cd terraform/

terraform init
terraform apply
```

## Usage

With Docker:

```
docker-compose up
```

```
go run . generate-user-ids
```

```
go run . start
```