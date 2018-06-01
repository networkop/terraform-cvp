package main

import (
	"log"

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
				Default:  60,
			},
			"reconcile": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"configlets": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"push": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}

func resourceDeviceCreate(d *schema.ResourceData, meta interface{}) error {
	client := *meta.(*CvpClient).Client
	var containerString string

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

	// Wait X seconds for the device to boot and save it to inventory
	timeout := d.Get("wait").(int)
	if err := client.SaveCommit(address, timeout); err != nil {
		log.Printf("[INFO] Could not add/save the device into inventory: %+v", err)
	}

	d.SetId(address)
	// From here on, assuming that the device has been created

	if err := resourceDeviceRead(d, meta); err == nil {
		// Reconcile existing configuration
		if reconcile := d.Get("reconcile").(bool); reconcile {
			log.Printf("[INFO] Trying to reconcile existing configuration")

			if err := reconcileDeviceConfiglet(d, meta); err != nil {
				return err
			}
		}
		// Assign configlets
		if _, ok := d.GetOk("configlets"); ok {
			if err := assignDeviceConfiglets(d, meta); err != nil {
				return err
			}
		}
	}

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

func reconcileDeviceConfiglet(d *schema.ResourceData, meta interface{}) error {
	client := *meta.(*CvpClient).Client

	mac := d.Get("system_mac_address").(string)

	log.Printf("[INFO] Trying to generate reconcile configlet for %s", mac)
	r, err := client.ValidateCompareCfglt(mac, []string{})
	if err != nil {
		return err
	}

	log.Printf("[INFO] Trying to create reconcile configlet with name %s", r.ReconciledConfig.Name)
	err = client.UpdateReconcile(mac, r.ReconciledConfig.Name, r.ReconciledConfig.Config)
	if err != nil {
		return err
	}

	if err := assignDeviceConfiglet(d, meta, r.ReconciledConfig.Name, true); err != nil {
		return err
	}

	return nil
}

func assignDeviceConfiglets(d *schema.ResourceData, meta interface{}) error {
	confs := d.Get("configlets").(*schema.Set).List()

	for _, conf := range confs {
		cName := conf.(map[string]interface{})["name"].(string)
		push := conf.(map[string]interface{})["push"].(bool)

		if err := assignDeviceConfiglet(d, meta, cName, push); err != nil {
			return err
		}
	}

	return nil
}

func assignDeviceConfiglet(d *schema.ResourceData, meta interface{}, cName string, push bool) error {
	client := *meta.(*CvpClient).Client

	ip := d.Get("ip_address").(string)
	fqdn := d.Get("fqdn").(string)
	mac := d.Get("system_mac_address").(string)

	log.Printf("[INFO] Trying to add %s configlet to the device %s", cName, fqdn)
	sdata, err := client.ApplyConfigletToDevice(ip, fqdn, mac, []string{cName}, true)
	if err != nil {
		return err
	}

	if push {
		log.Printf("[INFO] Trying to execute pending tasks")
		taskIds := sdata.Data.TaskIds
		if err = client.ExecuteTasks(taskIds); err != nil {
			return err
		}

		log.Printf("[INFO] Trying to check if the tasks have been completed")
		if err = client.CheckTasks(taskIds, 10); err != nil {
			log.Printf("[INFO] Could not verify that the task has been pushed, check CVP tasks")
		}
	}
	return nil
}
