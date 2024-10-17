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

variable "mysql_pwd" {
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
      "echo 'mysql-server mysql-server/root_password password ${var.mysql_pwd}' | sudo debconf-set-selections",
      "echo 'mysql-server mysql-server/root_password_again password ${var.mysql_pwd}' | sudo debconf-set-selections",
      "sudo apt-get install -y mysql-server",
      "sudo systemctl enable mysql",
      "sudo systemctl start mysql"
    ]
  }

  provisioner "shell" {
    inline = [
      "sudo systemctl start mysql",
      "mysql -u root -p${var.mysql_pwd} -e 'CREATE DATABASE test;'"
    ]
  }

  provisioner "shell" {
    inline = [
      "echo 'Restarting MySQL service...'",
      "sudo systemctl restart mysql"
    ]
  }

  provisioner "shell" {
    inline = [
      "echo 'mysql_pwd=${var.mysql_pwd}' | sudo tee -a /etc/environment"
    ]
  }


  provisioner "file" {
    source      = "./webapp"
    destination = "/tmp/webapp"
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
}
  