# Arista CloudVision Terraform Provider

# Caveats

Currently only supports create/read/delete operations on Devices and Configlets

# Using the provider

```
# Configure the CVP provider
provider "cvp" {
  # NOTE: Environment Variables can also be used for authentication

  # cvp_address    = "..."
  # cvp_user       = "..."
  # cvp_pwd        = "..."
  # cvp_container  = "..."
}

# Create a resource group
# Optional wait parameter specifies how long to 
# wait for device's state to become "connected"
# before saving into CVP's inventory
resource "cvp_device" "Device-A" {
    ip_address = "192.168.100.1"
    wait = "60"
}

# Create a Configlet
resource "cvp_configlet" "Test1" {
    name = "Test1"
    config = <<EOF
    \nusername TEST privilege 1 nopassword
    hostname BLA
    EOF
}

# Create a Configlet from template
data "template_file" "init" {
    template = "${file("config.tpl")}"

    vars {
        username = "FOO"
        hostname = "BAR"
    }
} 

resource "cvp_configlet" "Test2" {
    name = "Test2"
    config = "${data.template.test2.rendered"}
}

```