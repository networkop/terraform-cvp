# Arista CloudVision Terraform Provider

# Caveats

Currently only supports CVP devices.

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
resource "cvp_device" "Device-A" {
    ip_address = "192.168.100.1"
}

```