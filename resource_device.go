package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDevice() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeviceCreate,
		Read:   resourceDeviceRead,
		Update: resourceDeviceCreate,
		Delete: resourceDeviceDelete,

		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"fqdn": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"serial_number": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"system_mac_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"internal_version": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"model_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"wait": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceDeviceCreate(d *schema.ResourceData, meta interface{}) error {
	client := *meta.(*CvpClient).Client
	var containerString string
	var timeout int

	address := d.Get("ip_address").(string)
	container, ok := d.GetOk("container")
	if !ok {
		containerString = meta.(*CvpClient).Container
	} else {
		containerString = container.(string)
	}
	if err := client.AddDevice(address, containerString); err != nil {
		return err
	}

	// Wait X seconds for the device to boot and saves it to inventory
	wait, ok := d.GetOk("wait")
	if !ok {
		timeout = 60
	} else {
		timeout = wait.(int)
	}
	if err := client.SaveCommit(address, timeout); err != nil {
		return err
	}

	d.SetId(address)
	return resourceDeviceRead(d, meta)
}

func resourceDeviceRead(d *schema.ResourceData, meta interface{}) error {
	client := *meta.(*CvpClient).Client

	obj, err := client.GetDevice(d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("fqdn", obj.Fqdn)
	d.Set("key", obj.Key)
	d.Set("version", obj.Version)
	d.Set("serial_number", obj.SerialNumber)
	d.Set("system_mac_address", obj.SystemMacAddress)
	d.Set("internal_version", obj.InternalVersion)
	d.Set("mode_name", obj.ModeName)

	return nil
}

func resourceDeviceDelete(d *schema.ResourceData, meta interface{}) error {
	client := *meta.(*CvpClient).Client

	_, ok := d.GetOk("system_mac_address")
	if !ok {
		resourceDeviceRead(d, meta)
	}
	sysMac := d.Get("system_mac_address").(string)
	if err := client.RemoveDevice(sysMac); err != nil {
		return err
	}
	// d.SetId("") is automatically called assuming delete returns no errors, bu
	// it is added here for explicitness.
	d.SetId("")
	return nil
}
