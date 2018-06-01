# Arista CloudVision Terraform Provider

# Caveats

* Currently only supports create/read/delete operations on Devices and Configlets


# Configure the CVP provider

```
provider "cvp" {
  # NOTE: Environment Variables can also be used for authentication

  cvp_address    = "..."
  cvp_user       = "..."
  cvp_pwd        = "..."
  cvp_container  = "..."
}
```

# Device resource
Device resource creates a device inside CVP. The following options are available:

* **ip_address** (Required) - defines the IP address of EOS device for CVP to connect to.
* **wait** (Optional, Default is 60) - defines how long to wait for device to change state to "Connected". Reconcile and configlets defined below assume that device is "Connected".
* **container** (Optional, Default is 'Tenant') - CVP container to put the device into.
* **reconcile** (Optional, Default is false) - if set to true will attempt to reconcile the existing device configuration.
* **configlets** (Optional) - a list of configlets to assign and optionally push to a device. If **push** is ommitted, this simply creates a pending task.


```
resource "cvp_device" "Device-A" {
    ip_address = "192.168.100.1"
    wait = "60"
    container = "Tenant"
    reconcile = true
    configlets = [{
        name = "${cvp_configlet.test1.name}"
        push = true
    }]
}
```

# Configlet resource
Creates a configlet inside CVP, accepts the following parameters:

* **name** (Required) - configlet name
* **config** (Required) - configlet configuration


## Create a configlet with inline config
```
resource "cvp_configlet" "Test1" {
    name = "Test1"
    config = <<EOF
    \nusername TEST privilege 1 nopassword
    hostname BLA
    EOF
}
```

## Create a Configlet from template

```
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