variable "bootstrap_server_url" {
  type = string
}

resource "aws_iam_instance_profile" "instance_profile" {
  role = "AllowEC2ToUseMSK"
}

##
# Consumers instances
#

resource "aws_instance" "analytics_consumer" {

  ami           = "ami-0a94c8e4ca2674d5a"
  instance_type = "t3.small"

  tags = {
    Name = "simulation-consumer-analytics"
  }

  security_groups = ["default"]

  iam_instance_profile = aws_iam_instance_profile.instance_profile.name
  user_data = <<-EOF
              #!/bin/bash
              sed -i "s/Environment=*/Environment=KAFKA_BOOTSTRAP_SERVER=${var.bootstrap_server_url}/g" /lib/systemd/system/kafka-consumer.service
              systemctl daemon-reload 
              systemctl restart kafka-consumer.service
              EOF
}
