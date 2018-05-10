resource "cvp_device" "Device-A" {
    ip_address = "172.19.0.2"
    wait = "5"
}

resource "cvp_device" "Device-B" {
    ip_address = "172.19.0.3"
    wait = "20"
}

resource "cvp_configlet" "Test1" {
    name = "Test1"
    config = <<EOF
    \nusername TEST privilege 1 nopassword
    hostname BLA
    EOF
}

data "template_file" "test2" {
    template = "${file("config.tpl")}"

    vars {
        username = "FOO"
        hostname = "BAR"
    }
} 

resource "cvp_configlet" "Test2" {
    name = "Test2"
    config = "${data.template_file.test2.rendered}"
}