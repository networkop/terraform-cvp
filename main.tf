resource "cvp_device" "Device-A" {
    ip_address = "${var.ceos_1}"
}

resource "cvp_device" "Device-B" {
    ip_address = "${var.ceos_2}"
}

