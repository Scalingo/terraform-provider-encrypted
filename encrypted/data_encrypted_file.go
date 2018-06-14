package encrypted

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"
	"errors"
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

	ciphertext, err := ioutil.ReadFile(path)
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
		var parsed map[string]string
		err := json.Unmarshal(ciphertext, &parsed)
		if err != nil {
			return err
		}
		d.Set("parsed", parsed)
	}

	d.SetId(path)

	return nil
}
