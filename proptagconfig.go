package proptagconfig

import (
	"reflect"

	"errors"

	"strconv"

	"github.com/magiconair/properties"
	"github.com/Experticity/tagconfig"
)

// PropTagConfig implements the TagValueGetter to allow for populating a
// struct via "prop" struct tagconfig. It also implements TagValueSetter which
// allows for the struct to populate properties, which is basically the inverse
// of the TagValueGetter interface.
type PropTagConfig struct {
	*properties.Properties
}

var ErrUnsupportedSetType = errors.New("Received an unsupported type to set")

func (pv *PropTagConfig) TagName() string {
	return "props"
}

func (pv *PropTagConfig) Get(key string, _ reflect.StructField) string {
	return pv.Properties.GetString(key, "")
}

// Set is used to set properties based on values provided by the struct. For
// this iteration, it only supports a couple of major types, but if there's a
// need for this to grow, we can make this more sophisticated.
func (pv *PropTagConfig) Set(key string, value interface{}, _ reflect.StructField) error {
	var s string
	switch v := value.(type) {
	case string:
		s = v
	case bool:
		s = strconv.FormatBool(v)
	case int:
		s = strconv.Itoa(v)
	default:
		return ErrUnsupportedSetType
	}

	_, _, err := pv.Properties.Set(key, s)
	return err
}

// PopulatePropertiesFromStruct is used to create *properties.Properties based
// off of fields/data within a struct.
func PopulatePropertiesFromStruct(value interface{}) (*properties.Properties, error) {
	p := properties.NewProperties()
	pt := &PropTagConfig{p}
	err := tagconfig.PopulateExternalSource(pt, value)
	return p, err
}
