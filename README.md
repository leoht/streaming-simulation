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

## Usage (backend server)

With Docker (runs the webserver and a local DB):

```
docker-compose up
```

Without Docker:

```
go run cmd/web/main.go
```

Alternatively, run commands (without the web server):

```
go run . generate-user-ids
```

```
go run . start
```

## Start the frontend web UI

```
cd web/
npm start
```