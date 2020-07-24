package errtm

import (
	"errors"
	"fmt"
	"strings"
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
	Error        string   // This is the error gave by the server
	Attribute    string   // This is the attribute that doesn't set correctly
	State        errState // It could be one of the following: CREATING, SETTING, DELETING, UPDATING or SETTING
}

// SetID sets the ID of the resource
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
func (c *Config) SetError(r string) *Config {
	c.Error = r
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
func (c *Config) SetState(a errState) *Config {
	c.State = a
	return c
}

// ToError builds and returns the error message
func (c *Config) ToError() error {
	return NewErrorMessage(c)
}

// FillMessage sets the missing config attributes and returns the error message
func (c *Config) FillMessage(config *Config) error {
	if config.ID != "" {
		c.ID = config.ID
	}

	if config.ProviderName != "" {
		c.ProviderName = config.ProviderName
	}

	if config.ResourceName != "" {
		c.ResourceName = config.ResourceName
	}

	if config.Error != "" {
		c.Error = config.Error
	}

	if config.Attribute != "" {
		c.Attribute = config.Attribute
	}

	if config.State != "" {
		c.State = config.State
	}
	return NewErrorMessage(c)
}

// SetProviderName retunrs a Config struct setting the provider name
func SetProviderName(pn string) *Config {
	return &Config{ProviderName: pn}
}

// SetError retunrs a Config struct setting the message error
func SetError(err string) *Config {
	return &Config{Error: err}
}

// NewErrorMessage builds and creates the error message returning it as error type
func NewErrorMessage(c *Config) error {
	var err string
	words := []string{"error", strings.ToLower(string(c.State))}

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
			c.SetType(Setting)
		}

		attribute := fmt.Sprintf("attribute `%s`", c.Attribute)
		if c.Attribute == "" {
			attribute = "an attribute"
		}

		words = append([]string{words[0], strings.ToLower(string(c.State)), attribute, "in"}, words[2:]...)
	}

	err = strings.Join(words, " ")

	if c.Error != "" {
		err = fmt.Sprintf("%s: %s", err, c.Error)
	}

	return errors.New(err)
}
