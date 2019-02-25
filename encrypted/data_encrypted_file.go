package encrypted

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceEncryptedFile() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_path": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parsed": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"array": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
		Read: dataSourceEncryptedFileRead,
	}
}

func dataSourceEncryptedFileRead(d *schema.ResourceData, meta interface{}) error {
	keyS := meta.(string)
	path := d.Get("path").(string)

	key, err := hex.DecodeString(keyS)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(string(raw))
	if err != nil {
		return err
	}

	if len(ciphertext) < aes.BlockSize {
		return errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	d.Set("value", ciphertext)
	if d.Get("content_type").(string) == "json" {
		var parsed map[string]interface{}
		var parsedArray []interface{}
		err := json.Unmarshal(ciphertext, &parsed)
		if err != nil {
			return err
		}
		dataPath := d.Get("data_path").([]interface{})
		if dataPath != nil {
			for i, segment := range dataPath {
				if v, ok := parsed[segment.(string)].([]interface{}); i == len(dataPath)-1 && ok {
					parsedArray = v
				} else {
					v, ok := parsed[segment.(string)].(map[string]interface{})
					if ok {
						parsed = v
					} else {
						return fmt.Errorf("invalid data_path %v", dataPath)
					}
				}
			}
		}
		if parsedArray != nil {
			d.Set("array", parsedArray)
		} else {
			parsed = flatten(parsed)
			d.Set("parsed", parsed)
		}
	}

	d.SetId(path)

	return nil
}

func flatten(m map[string]interface{}) map[string]interface{} {
	dest := map[string]interface{}{}
	flattenToDest(dest, m, "")
	return dest
}

func flattenToDest(dest map[string]interface{}, m map[string]interface{}, prefix string) {
	for key, value := range m {
		if submap, ok := value.(map[string]interface{}); ok {
			flattenToDest(dest, submap, prefix+key+"_")
		} else {
			dest[prefix+key] = fmt.Sprint(value)
		}
	}
}
