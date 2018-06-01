package main

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	var p *schema.Provider
	p = &schema.Provider{
		Schema: map[string]*schema.Schema{
			"cvp_address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CVP_ADDRESS", ""),
			},

			"cvp_user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CVP_USER", ""),
			},

			"cvp_pwd": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CVP_PWD", ""),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"cvp_device":    resourceDevice(),
			"cvp_configlet": resourceConfiglet(),
		},
	}

	p.ConfigureFunc = providerConfigure(p)

	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		config := &CvpInfo{
			CvpAddress: d.Get("cvp_address").(string),
			CvpUser:    d.Get("cvp_user").(string),
			CvpPwd:     d.Get("cvp_pwd").(string),
		}

		if err := validateConfig(config); err != nil {
			return nil, err
		}

		client, err := getCvpClient(config)
		if err != nil {
			return nil, err
		}

		return client, nil
	}
}

func validateConfig(c *CvpInfo) error {
	if c.CvpAddress == "" {
		return fmt.Errorf("Please provider CVP address")
	}
	if c.CvpUser == "" {
		return fmt.Errorf("Please provider CVP username")
	}
	if c.CvpPwd == "" {
		return fmt.Errorf("Please provider CVP password")
	}
	return nil
}
