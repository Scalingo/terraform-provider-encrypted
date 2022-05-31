package encrypted

import (
	"reflect"
	"testing"

	"github.com/zclconf/go-cty/cty"
)

func TestFlatten(t *testing.T) {
	cases := []struct {
		Name string
		Map  cty.Value
		Res  map[string]string
	}{
		{
			Name: "no change to flat map",
			Map: cty.ObjectVal(map[string]cty.Value{
				"pipo":  cty.StringVal("molo"),
				"hello": cty.StringVal("world"),
			}),
			Res: map[string]string{"pipo": "molo", "hello": "world"},
		}, {
			Name: "it should flatten sub map with '.'",
			Map: cty.ObjectVal(map[string]cty.Value{
				"pipo": cty.StringVal("molo"),
				"sub": cty.MapVal(map[string]cty.Value{
					"hello": cty.StringVal("world"),
				}),
			}),
			Res: map[string]string{"pipo": "molo", "sub.hello": "world"},
		}, {
			Name: "it should flatten sub sub map also",
			Map: cty.ObjectVal(map[string]cty.Value{
				"pipo": cty.StringVal("molo"),
				"sub": cty.ObjectVal(map[string]cty.Value{
					"hello": cty.StringVal("world"),
					"foo": cty.MapVal(map[string]cty.Value{
						"bar": cty.StringVal("value"),
					}),
				}),
			}),
			Res: map[string]string{"pipo": "molo", "sub.hello": "world", "sub.foo.bar": "value"},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if !reflect.DeepEqual(c.Res, FlatmapValueFromHCL2(c.Map)) {
				t.Errorf("%v should be equal to %v", FlatmapValueFromHCL2(c.Map), c.Res)
			}
		})
	}
}
