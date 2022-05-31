package encrypted

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"os"

	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	ctyyaml "github.com/zclconf/go-cty-yaml"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

var _ tfsdk.DataSourceType = encryptedFileDataSourceType{}
var _ tfsdk.DataSource = encryptedFileDataSource{}

type encryptedFileDataSourceType struct{}

func (t encryptedFileDataSourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Encrypted File data source",

		Attributes: map[string]tfsdk.Attribute{
			"path": {
				Type:     types.StringType,
				Required: true,
			},
			"data_path": {
				Type:     types.ListType{ElemType: types.StringType},
				Optional: true,
			},
			"content_type": {
				Type:     types.StringType,
				Optional: true,
			},
			"value": {
				Type:     types.StringType,
				Computed: true,
			},
			"parsed": {
				Type:     types.MapType{ElemType: types.StringType},
				Computed: true,
			},
			"array": {
				Type:     types.ListType{ElemType: types.StringType},
				Computed: true,
			},
		},
	}, nil
}

func (t encryptedFileDataSourceType) NewDataSource(ctx context.Context, in tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	provider, diags := convertProviderType(in)

	return encryptedFileDataSource{
		provider: provider,
	}, diags
}

type encryptedFileDataSourceData struct {
	Path        types.String      `tfsdk:"path"`
	DataPath    []string          `tfsdk:"data_path"`
	ContentType types.String      `tfsdk:"content_type"`
	Value       types.String      `tfsdk:"value"`
	Parsed      map[string]string `tfsdk:"parsed"`
	Array       []interface{}     `tfsdk:"array"`
}

type encryptedFileDataSource struct {
	provider encrypted
}

func (d encryptedFileDataSource) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	data := encryptedFileDataSourceData{}
	config := encryptedData{}

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	key := stringFromConfigOrEnv(config.Key, "ENCRYPTION_KEY", "")
	dataSourceEncryptedFileRead(&data, key)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

func dataSourceEncryptedFileRead(d *encryptedFileDataSourceData, keyS string) error {
	path := d.Path.Value

	key, err := hex.DecodeString(keyS)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	raw, err := os.ReadFile(path)
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

	d.Value = types.String{Value: string(ciphertext)}

	parse(d, ciphertext)

	return nil
}

func parse(d *encryptedFileDataSourceData, data []byte) error {
	funcs := map[string]function.Function{
		"json": stdlib.JSONDecodeFunc,
		"yaml": ctyyaml.YAMLDecodeFunc,
	}

	ctyValues := []cty.Value{
		cty.StringVal(string(data)),
	}

	value, err := funcs[d.ContentType.Value].Call(ctyValues)
	if err != nil {
		return err
	}

	valueType := value.Type()
	if valueType.IsTupleType() || valueType.IsListType() || valueType.IsSetType() {
		d.Array = ConfigValueFromHCL2(value).([]interface{})
	} else {
		d.Parsed = FlatmapValueFromHCL2(value)
	}

	return nil
}

func stringFromConfigOrEnv(value types.String, env string, def string) string {
	if value.Unknown || value.Null || value.Value == "" {
		value := os.Getenv(env)

		if value != "" {
			return value
		}
	}

	if value.Value == "" {
		return def
	}

	return value.Value
}
