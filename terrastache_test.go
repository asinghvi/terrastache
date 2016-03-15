package main

import (
	"strings"
	"testing"
)

const template = `
{{^sg_elb}}
resource "aws_security_group" "sg_elb" {
  name = "${var.service_name}-elb-${var.service_env}"
  description = "${var.service_name} service ${var.service_env} ELB"
  vpc_id = "${var.vpc_id}"

  # Allow internal inbound HTTP access over port 80
  ingress {
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = ["10.0.0.0/8", "172.16.0.0/12"]
  }

  # Allow all outbound traffic
  egress {
    from_port = 0
    to_port = 0
    protocol = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }
}
{{/sg_elb}}

resource "aws_elb" "web" {
{{^sg_elb}}
  depends_on = ["aws_security_group.sg_elb"]
{{/sg_elb}}

  name = "${var.service_name}-elb-${var.service_env}"
  subnets  = ["${split(",", var.subnet_ids)}"]
{{#sg_elb}}
  security_groups = ["{{{sg_elb}}}"]
{{/sg_elb}}
{{^sg_elb}}
  security_groups = ["${aws_security_group.sg_elb.id}"]
{{/sg_elb}}

  cross_zone_load_balancing = true
  idle_timeout = 60
  connection_draining = true
  connection_draining_timeout = 15
  listener {
    instance_port = "${var.service_web_port}"
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }
}
`

func TestRenderTemplate(t *testing.T) {
	rendered, err := renderTemplate(template, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}

	expected := `
resource "aws_security_group" "sg_elb" {
  name = "${var.service_name}-elb-${var.service_env}"
  description = "${var.service_name} service ${var.service_env} ELB"
  vpc_id = "${var.vpc_id}"

  # Allow internal inbound HTTP access over port 80
  ingress {
    from_port = 80
    to_port = 80
    protocol = "tcp"
    cidr_blocks = ["10.0.0.0/8", "172.16.0.0/12"]
  }

  # Allow all outbound traffic
  egress {
    from_port = 0
    to_port = 0
    protocol = -1
    cidr_blocks = ["0.0.0.0/0"]
  }

  lifecycle {
    create_before_destroy = true
  }
}


resource "aws_elb" "web" {
  depends_on = ["aws_security_group.sg_elb"]


  name = "${var.service_name}-elb-${var.service_env}"
  subnets  = ["${split(",", var.subnet_ids)}"]

  security_groups = ["${aws_security_group.sg_elb.id}"]


  cross_zone_load_balancing = true
  idle_timeout = 60
  connection_draining = true
  connection_draining_timeout = 15
  listener {
    instance_port = "${var.service_web_port}"
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }
}
`

	expected = strings.TrimSpace(expected)
	rendered = strings.TrimSpace(rendered)
	if expected != rendered {
		t.Errorf("Expected:\n%s\n\nActual:\n%s\n", expected, rendered)
	}
}

func TestRenderTemplateWithSG(t *testing.T) {
	rendered, err := renderTemplate(template, map[string]string{"sg_elb": "sg-blah"})
	if err != nil {
		t.Fatal(err)
	}

	expected := `
resource "aws_elb" "web" {


  name = "${var.service_name}-elb-${var.service_env}"
  subnets  = ["${split(",", var.subnet_ids)}"]
  security_groups = ["sg-blah"]



  cross_zone_load_balancing = true
  idle_timeout = 60
  connection_draining = true
  connection_draining_timeout = 15
  listener {
    instance_port = "${var.service_web_port}"
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }
}
`

	expected = strings.TrimSpace(expected)
	rendered = strings.TrimSpace(rendered)
	if expected != rendered {
		t.Errorf("Expected:\n%s\n\nActual:\n%s\n", expected, rendered)
	}
}
