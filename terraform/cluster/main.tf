
locals {
  msk = {
    version                    = "3.5.1"
    instance_size              = "kafka.t3.small"
    ebs_volume_size            = 20
    log_retention_ms           = 604800000 # 7 days
    num_partitions             = 2
    default_replication_factor = 2
  }
}

variable "cluster_name" {
  type = string
}

data "aws_availability_zones" "azs" {
  state = "available"
}

data "aws_vpc" "default_vpc" {
  default = true
}

data "aws_security_group" "sg" {
  id = "sg-39d6e247"
}

data "aws_subnet" "subnet_az1" {
  vpc_id            = data.aws_vpc.default_vpc.id
  id = "subnet-a41657de"
}

data "aws_subnet" "subnet_az2" {
  vpc_id            = data.aws_vpc.default_vpc.id
  id = "subnet-7213161b"
}

resource "aws_s3_bucket" "bucket" {
  bucket = "${var.cluster_name}-bucket"

  force_destroy = true
}

resource "aws_kinesis_firehose_delivery_stream" "test_stream" {
  name        = "${var.cluster_name}-kinesis-firehose-msk-broker-logs-stream"
  destination = "extended_s3"

  extended_s3_configuration {
    role_arn   = aws_iam_role.firehose_role.arn
    bucket_arn = aws_s3_bucket.bucket.arn
  }

  tags = {
    LogDeliveryEnabled = "placeholder"
  }

  lifecycle {
    ignore_changes = [
      tags["LogDeliveryEnabled"],
    ]
  }
}


data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["firehose.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "firehose_role" {
  name               = "events-streaming_firehose_test_role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_msk_cluster" "msk_data_cluster" {
  cluster_name           = "${var.cluster_name}"
  kafka_version          = local.msk.version
  number_of_broker_nodes = 2
  configuration_info {
    arn      = aws_msk_configuration.msk_config.arn
    revision = aws_msk_configuration.msk_config.latest_revision
  }

  broker_node_group_info {
    instance_type   = local.msk.instance_size
    client_subnets  = [data.aws_subnet.subnet_az1.id, data.aws_subnet.subnet_az2.id]
    security_groups = [data.aws_security_group.sg.id]
    storage_info {
      ebs_storage_info {
        volume_size = local.msk.ebs_volume_size
      }
    }
  }

  client_authentication {
    sasl {
      iam = true
    }
  }

  logging_info {
    broker_logs {
      s3 {
        enabled = true
        bucket  = aws_s3_bucket.bucket.id
        prefix  = "logs/msk-"
      }
    }     
  }

  tags = {}

  depends_on = [aws_msk_configuration.msk_config]
}

resource "aws_msk_configuration" "msk_config" {
  name = "${var.cluster_name}-msk-configuration"

  kafka_versions = [local.msk.version]

  server_properties = <<PROPERTIES
    auto.create.topics.enable = true
    delete.topic.enable = true
    log.retention.ms = ${local.msk.log_retention_ms}
    num.partitions = ${local.msk.num_partitions}
    default.replication.factor = ${local.msk.default_replication_factor}
  PROPERTIES
}


output "bootstrap_server_url" {
  value = split(",", aws_msk_cluster.msk_data_cluster.bootstrap_brokers_sasl_iam)[0]
}