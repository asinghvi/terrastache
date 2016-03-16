# Terrastache (Terraform + mustache)

Terrastache allows you to generate a Terraform config using a mustache template and existing an Terraform vars file.

The use case that inspired this project was to use a common config across environments, but gracefully handle the fact that some environments had existing resources which should be used and not managed by Terraform. We had a production environment with an ELB and a number of standalone instances as a web tier. And a staging environment with a separate ELB and web instance, but which shared the production security groups. We were moving to Terraform to automate the creation of these environments, since doing so manually was error prone and tedious, leading to such bad pratices as sharing security groups.

The new configuration would use ASGs in place of standalone instances. The staging environment could be replaced completely, and replacing stateless web nodes would be fine for production. But moving to a new ELB would require warming up, and adjusting existing security groups seemed needlessly risky. What we wanted was to take a common Terraform template, and with environment-specific var files control whether existing resources we reused, or if Terraform would create them.

Consider this example for defining the ASG - depending on whether or not a value is provided for `elb`, it will either reference that ELB, or have Terraform create a new one:
```hcl
resource "aws_autoscaling_group" "asg" {
{{#elb}}
  depends_on = ["aws_launch_configuration.asg_conf"]
{{/elb}}
{{^elb}}
  depends_on = ["aws_elb.web", "aws_launch_configuration.asg_conf"]
{{/elb}}

  name = "${var.service_name}-${var.service_env}-web"
  availability_zones = ["${split(",", var.availability_zones)}"]
  vpc_zone_identifier  = ["${split(",", var.subnet_ids)}"]
  min_size = 0
  max_size = 5
  desired_capacity = 0
  launch_configuration = "${aws_launch_configuration.asg_conf.name}"
{{#elb}}
  load_balancers = ["{{{elb}}}"]
{{/elb}}
{{^elb}}
  load_balancers = ["${aws_elb.web.name}"]
{{/elb}}
}
```

Assuming the environment-specific variables were stored in files named `$(env).tfvars`, you would generate and run Terraform using the following commands:
```
$ terrastache -var-file $(env).tfvars -template service.mustache > service.tf
$ terraform plan -var-file $(env).tfvars -state=$(env).tfstate
```
Integrating mustache with the Terraform config parser allows a common vars file to be used to control Terraform and Terrastache.
