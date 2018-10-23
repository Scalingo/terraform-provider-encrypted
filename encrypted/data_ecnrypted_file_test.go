package encrypted

import "reflect"
import "testing"

func TestFlatten(t *testing.T) {
	cases := []struct {
		Name string
		Map  map[string]interface{}
		Res  map[string]interface{}
	}{
		{
			Name: "no change to flat map",
			Map:  map[string]interface{}{"pipo": "molo", "hello": "world"},
			Res:  map[string]interface{}{"pipo": "molo", "hello": "world"},
		}, {
			Name: "it should flatten sub map with '_'",
			Map:  map[string]interface{}{"pipo": "molo", "sub": map[string]interface{}{"hello": "world"}},
			Res:  map[string]interface{}{"pipo": "molo", "sub_hello": "world"},
		}, {
			Name: "it should flatten sub sub map also",
			Map: map[string]interface{}{
				"pipo": "molo",
				"sub": map[string]interface{}{
					"hello": "world",
					"foo": map[string]interface{}{
						"bar": "value",
					},
				},
			},
			Res: map[string]interface{}{"pipo": "molo", "sub_hello": "world", "sub_foo_bar": "value"},
		},
	}

	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			if !reflect.DeepEqual(c.Res, flatten(c.Map)) {
				t.Errorf("%v should be equal to %v", flatten(c.Map), c.Res)
			}
		})
	}
}
