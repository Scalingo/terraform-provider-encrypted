package encrypted

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:      schema.TypeString,
				Sensitive: true,
				Required:  true,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"encrypted_file": dataSourceEncryptedFile(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(data *schema.ResourceData) (interface{}, error) {
	key := data.Get("key").(string)

	return key, nil
}
