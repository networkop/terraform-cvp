resource "cvp_device" "Device-A" {
    ip_address = "192.168.100.1"
    wait = "5"
    container = "NEW_CONTAINER"
    reconcile = true
    configlets = [{
        name = "${cvp_configlet.test1.name}"
        push = true
    },{
        name = "${cvp_configlet.test2.name}"
        push = true
    }
    ]
    depends_on = ["cvp_configlet.test1"]
}

resource "cvp_device" "Device-B" {
    ip_address = "192.168.100.2"
    wait = "20"
    reconcile = true
    configlets = [{
        name = "${cvp_configlet.test2.name}"
        push = true
    }]
    depends_on = ["cvp_configlet.test2"]
}

resource "cvp_configlet" "test1" {
    name = "Test1"
    config = <<EOF
    username TEST1 privilege 1 nopassword
    EOF
}

data "template_file" "test2" {
    template = "${file("config.tpl")}"

    vars {
        username = "FOO"
    }
} 

resource "cvp_configlet" "test2" {
    name = "Test2"
    config = "${data.template_file.test2.rendered}"
}