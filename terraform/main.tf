provider "aws" {
  region = "ap-south-1"
}


// new aws instance
resource "aws_instance" "bot_instance" {
  ami           = "ami-02b8269d5e85954ef" // amazon ami image id (ubuntu 64 bit architecture)
  instance_type = "t2.micro"
  # user_data = <- EOF

  // dashboard instance name
  tags = {
    Name = "tg-bot"
  }
}

// defining our security group here for allowed traffic on our server
resource "aws_security_group" "bot_firewall" {
  name        = "bot-firewall"
  description = "Allow SSH and Web traffic"
  vpc_id      = aws_vpc.main.id
}

// open ssh port
resource "aws_vpc_security_group_ingress_rule" "allow_ssh" {
  security_group_id = aws_security_group.bot_firewall.id

  # ssh allowed only from my ip
  cidr_ipv4   = "${var.my_ip}/32"
  from_port   = 22
  ip_protocol = "tcp"
  to_port     = 22
}

// inbound rule for all ipv4
resource "aws_vpc_security_group_ingress_rule" "allow_http" {
  security_group_id = aws_security_group.bot_firewall.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 80
  ip_protocol       = "tcp"
  to_port           = 80
}

# outbound Rule allow all outbound traffic
resource "aws_vpc_security_group_egress_rule" "allow_all_outbound" {
  security_group_id = aws_security_group.bot_firewall.id
  cidr_ipv4         = "0.0.0.0/0"
  ip_protocol       = "-1" # -1 means "all protocols"
}
