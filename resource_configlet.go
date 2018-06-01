package main

import (
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	cvpgo "github.com/networkop/cvpgo/client"
)

func resourceConfiglet() *schema.Resource {
	return &schema.Resource{
		Create: resourceConfigletCreate,
		Read:   resourceConfigletRead,
		Update: resourceConfigletUpdate,
		Delete: resourceConfigletDelete,

		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"config": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceConfigletCreate(d *schema.ResourceData, meta interface{}) error {
	client := *meta.(*CvpClient).Client

	configlet := cvpgo.Configlet{
		Name:   d.Get("name").(string),
		Config: d.Get("config").(string),
	}
	if _, err := client.AddConfiglet(configlet); err != nil {
		return err
	}

	d.SetId(configlet.Name)
	return resourceConfigletRead(d, meta)
}

func resourceConfigletRead(d *schema.ResourceData, meta interface{}) error {
	client := *meta.(*CvpClient).Client

	obj, err := client.GetConfigletByName(d.Id())
	if err != nil {
		d.SetId("")
		return nil
	}

	d.Set("key", obj.Key)

	return nil
}

func resourceConfigletUpdate(d *schema.ResourceData, meta interface{}) error {
	// TODO - Need Updates to Fred's cvpgo
	return nil
}

func resourceConfigletDelete(d *schema.ResourceData, meta interface{}) error {
	client := *meta.(*CvpClient).Client

	// This is to prevent race condition in CVP when device is removed
	// But configlet is still marked as assigned to a device
	time.Sleep(500 * time.Millisecond)

	name := d.Get("name").(string)
	if err := client.DeleteConfiglet(name); err != nil {
		return err
	}
	// d.SetId("") is automatically called assuming delete returns no errors, bu
	// it is added here for explicitness.
	d.SetId("")
	return nil
}
