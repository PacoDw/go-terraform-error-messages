package errtm

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/spf13/cast"
)

type errState string

const (
	// Creating represents the creating state error
	Creating = errState("CREATING")
	// Reading represents the read state error
	Reading = errState("READING")
	// Updating represents the updating state error
	Updating = errState("UPDATING")
	// Deleting represents the deleting state error
	Deleting = errState("DELETING")
	// Setting represents the setting state error
	Setting = errState("SETTING")
)

// Config represents the struct to create the Terraform error message
type Config struct {
	ID           string   // Describes the resource id
	ProviderName string   // Name of the provider
	ResourceName string   // Describes the resource name
	ErrorMessage string   // This is the error gave by the server
	Attribute    string   // This is the attribute that doesn't set correctly
	State        errState // It could be one of the following: CREATING, SETTING, DELETING, UPDATING or SETTING

	// this is a reusable config, is must used with the Set Methods, this allows to create a temporaly values to
	// create a message error, so this will avoid modify the values of the configuration error that were saved/set
	// with the Save methods
	// if the user want o use just one method set to create an error avoiding extra info, this avoid it.
	storageConfig *Config

	// storageConfig all the Save's methods storage the data in this attribute, so the first level of methods will remplace to
	// volatil storage I mean will be storageConfig

	partialConfig *Config

	copyConfig *Config

	PartialMode bool
}

// EnablePartialMode enables a partial time where is possible to use the Set Methods to create a temporaly error Template
// so when EnablePartialMode disables all the configuration set after the partial will remove
func (c *Config) EnablePartialMode(pm bool) *Config {
	c.PartialMode = pm

	if !pm {
		c.partialConfig = &Config{}
		return c
	}

	c.copyConfig = c

	return c.partialConfig
}

// SaveID saves the ID in the template configuration to be used in all error messages
func (c *Config) SaveID(id string) *Config {
	if c.storageConfig == nil {
		c.storageConfig = &Config{}
	}
	c.storageConfig.ID = id
	return c
}

// SaveProviderName saves a provider name in the template configuration to be used in all error messages
func (c *Config) SaveProviderName(pn string) *Config {
	if c.storageConfig == nil {
		c.storageConfig = &Config{}
	}
	c.storageConfig.ProviderName = pn
	return c
}

// SaveResourceName saves a resource name in the template configuration to be used in all error messages
func (c *Config) SaveResourceName(rn string) *Config {
	if c.storageConfig == nil {
		c.storageConfig = &Config{}
	}
	c.storageConfig.ResourceName = rn
	return c
}

// SaveError saves a error message gave by the API into the template configuration to be used in all error messages
func (c *Config) SaveError(e interface{}) *Config {
	if c.storageConfig == nil {
		c.storageConfig = &Config{}
	}
	c.storageConfig.ErrorMessage = cast.ToString(e)
	return c
}

// SaveAttribute saves the attribute that does't set correctly in terraform into the template configuration to be used in all error messages
func (c *Config) SaveAttribute(a string) *Config {
	if c.storageConfig == nil {
		c.storageConfig = &Config{}
	}
	c.storageConfig.Attribute = a
	return c
}

// SaveState saves an error type in the template configuration to be used in all error messages,
// the value depending on which method/circumstance
// it occurs, we recommend use one of the follows const:
// Creating, Reading, Updating, Deleting or Setting
func (c *Config) SaveState(a errState) *Config {
	if c.storageConfig == nil {
		c.storageConfig = &Config{}
	}
	c.storageConfig.State = a
	return c
}

// SetID sets the ID as a temporaly value to create one message error
func (c *Config) SetID(id string) *Config {
	c.ID = id
	return c
}

// SetProviderName sets a provider name of the current provider
func (c *Config) SetProviderName(pn string) *Config {
	c.ProviderName = pn
	return c
}

// SetResourceName sets a resource name
func (c *Config) SetResourceName(rn string) *Config {
	c.ResourceName = rn
	return c
}

// SetError sets a error message gave by the API
func (c *Config) SetError(e interface{}) *Config {
	c.ErrorMessage = cast.ToString(e)
	return c
}

// SetAttribute a attribute that does't set correctly in terraform
func (c *Config) SetAttribute(a string) *Config {
	c.Attribute = a
	return c
}

// SetState sets an error type, this depending on which method/circumstance
// it occurs, we recommend use one of the follows const:
// Creating, Reading, Updating, Deleting or Setting
func (c *Config) SetState(s errState) *Config {
	c.State = s
	return c
}

// ToError builds and returns the error message
// Note: if you have created a previous Conf Error to be used in all the message
// and later is used the Set's Method, it will create a new message error
// copy first the Configuration Error and ovewrite it with the Temporaly configuration Error
func (c *Config) Error() string {
	var temp Config = *c

	if c.PartialMode {
		c = c.copyConfig
		temp = *c.partialConfig
	}

	if !reflect.DeepEqual(c.storageConfig, &Config{}) {
		temp = *fillMessage(c.storageConfig, &temp)
	}

	c.cleanConfig()

	return NewErrorMessage(&temp).Error()
}

func (c *Config) cleanConfig() {
	c.ID = ""
	c.ProviderName = ""
	c.ResourceName = ""
	c.ErrorMessage = ""
	c.Attribute = ""
	c.State = ""
}

func fillMessage(c, newConfig *Config) *Config {
	var temp Config = *c

	if newConfig.ID != "" {
		temp.ID = newConfig.ID
	}

	if newConfig.ProviderName != "" {
		temp.ProviderName = newConfig.ProviderName
	}

	if newConfig.ResourceName != "" {
		temp.ResourceName = newConfig.ResourceName
	}

	if newConfig.ErrorMessage != "" {
		temp.ErrorMessage = newConfig.ErrorMessage
	}

	if newConfig.Attribute != "" {
		temp.Attribute = newConfig.Attribute
	}

	if newConfig.State != "" {
		temp.State = newConfig.State
	}

	return &temp
}

// FillMessage sets the missing config attributes and returns the error message
func (c *Config) FillMessage(config *Config) error {
	var temp Config = *c

	if c.PartialMode {
		c = c.copyConfig
		temp = *c.partialConfig
	}

	t := fillMessage(fillMessage(c.storageConfig, &temp), config)

	c.cleanConfig()

	return NewErrorMessage(t)
}

// SaveProviderName retunrs a Config struct setting the provider name
func SaveProviderName(pn string) *Config {
	return &Config{storageConfig: &Config{ProviderName: pn}}
}

// SetError retunrs a Config struct setting the provider name
func SetError(err string) *Config {
	return &Config{ErrorMessage: err, storageConfig: &Config{}}
}

// NewConfigurationError creates a new Configuration Error
func NewConfigurationError() *Config {
	return &Config{storageConfig: &Config{}}
}

// NewErrorMessage builds and creates the error message returning it as error type
func NewErrorMessage(c *Config) error {
	var (
		err   string
		words = []string{"error", strings.ToLower(string(c.State))}
	)

	if c.State == "" {
		words = words[:1]
		if c.ProviderName != "" || c.ResourceName != "" {
			words = append(words, "in")
		}
	}

	if c.ProviderName != "" {
		words = append(words, c.ProviderName)
	}

	if c.ResourceName != "" {
		words = append(words, c.ResourceName)
	}

	if c.ID != "" {
		words = append(words, fmt.Sprintf("(%s)", c.ID))
	}

	if c.Attribute != "" || c.State == Setting {
		if c.State == "" {
			c.State = Setting
		}

		attribute := fmt.Sprintf("attribute `%s`", c.Attribute)
		if c.Attribute == "" {
			attribute = "an attribute"
		}

		words = append([]string{words[0], strings.ToLower(string(c.State)), attribute, "in"}, words[2:]...)
	}

	err = strings.Join(words, " ")

	if c.ErrorMessage != "" {
		err = fmt.Sprintf("%s: %s", err, c.ErrorMessage)
	}

	return errors.New(err)
}
