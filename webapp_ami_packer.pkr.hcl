variable "aws_region" {
  type = string
}

variable "source_ami" {
  type = string
}

variable "instance_type" {
  type = string
}

variable "ami_user" {
  type = string
}

variable "aws_profile" {
  type = string
}

variable "ssh_username" {
  type = string
}

variable "ami_name" {
  type = string
}

packer {
  required_plugins {
    amazon = {
      version = ">= 1"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

source "amazon-ebs" "webapp" {
  region                      = var.aws_region
  source_ami                  = var.source_ami
  instance_type               = var.instance_type
  ssh_username                = var.ssh_username
  ami_name                    = var.ami_name
  associate_public_ip_address = true
  profile                     = var.aws_profile
  ami_users                   = [var.ami_user]
  ssh_timeout                 = "30m"
}

build {
  sources = ["source.amazon-ebs.webapp"]

  provisioner "shell" {
    inline = [
      "sudo apt-get update",
      "sudo apt-get install -y curl",
      "curl https://amazoncloudwatch-agent.s3.amazonaws.com/ubuntu/amd64/latest/amazon-cloudwatch-agent.deb -o amazon-cloudwatch-agent.deb",
      "sudo dpkg -i amazon-cloudwatch-agent.deb",
      "sudo systemctl enable amazon-cloudwatch-agent",
    ]
  }

  provisioner "file" {
    source      = "./webapp"
    destination = "/tmp/webapp"
  }

  provisioner "file" {
    source      = "./amazon-cloudwatch-agent.json"
    destination = "/tmp/amazon-cloudwatch-agent.json"
  }

  provisioner "shell" {
    inline = [
      "sudo mv /tmp/webapp /opt/webapp",
      "sudo groupadd csye6225",
      "sudo useradd -g csye6225 -s /usr/sbin/nologin csye6225",
      "sudo chown csye6225:csye6225 /opt/webapp",
      "sudo chmod 750 /opt/webapp",
    ]
  }

  provisioner "shell" {
    inline = [
      "sudo mv /tmp/amazon-cloudwatch-agent.json /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json",
      "sudo chown csye6225:csye6225 /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json",
      "sudo chmod 750 /opt/aws/amazon-cloudwatch-agent/etc/amazon-cloudwatch-agent.json",
    ]
  }

  provisioner "shell" {
    inline = [
      "sudo touch /var/log/webapp.log",
      "sudo chown csye6225:csye6225 /var/log/webapp.log",
    ]
  }

  provisioner "file" {
    source      = "./webapp.service"
    destination = "/tmp/webapp.service"
  }

  provisioner "shell" {
    inline = [
      "sudo mv /tmp/webapp.service /etc/systemd/system/webapp.service",
      "sudo systemctl daemon-reload",
      "sudo systemctl enable webapp.service"
    ]
  }

  provisioner "shell" {
    inline = [
      "sudo systemctl daemon-reload",
      "sudo systemctl enable amazon-cloudwatch-agent.service"
    ]
  }
}
  