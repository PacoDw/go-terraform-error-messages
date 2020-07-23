package errtm

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func TestBasicMessageError(t *testing.T) {
	var (
		ms *Config = SetProviderName("MyProvider")

		got      = ms.ToError()
		expected = errors.New("error in MyProvider")
	)

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}

	got = ms.SetType(Setting).ToError()
	expected = errors.New("error setting an attribute in MyProvider")

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}

	got = SetError("this is my error").ToError()
	expected = errors.New("error: this is my error")

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}
}

func TestUsingNewErrorMessage_DifferentCRUDTypes(t *testing.T) {
	errTypes := []errType{Creating, Reading, Updating, Deleting}

	for _, errType := range errTypes {
		errType := errType
		t.Run(string(errType), func(t *testing.T) {
			t.Parallel()

			expected := fmt.Errorf("error %s TerraformProvider PeeringConnection (5456543433545656): Error processing your request", strings.ToLower(string(errType)))

			got := NewErrorMessage(&Config{
				ID:           "5456543433545656",
				ProviderName: "TerraformProvider",
				ResourceName: "PeeringConnection",
				Error:        "Error processing your request",
				Type:         errType,
			})

			if diff := deep.Equal(got, expected); diff != nil {
				t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
			}
		})
	}
}

func TestUsingSetProviderName_DifferentCRUDTypes(t *testing.T) {
	var (
		errTypes = []errType{Creating, Reading, Updating, Deleting}

		// You can create a template setting some attribute to be used to create other errors
		err = SetProviderName("MyProvider").SetResourceName("Network Peering Connection")
	)

	for _, errType := range errTypes {
		errType := errType
		t.Run(string(errType), func(t *testing.T) {
			t.Parallel()

			expected := fmt.Errorf("error %s MyProvider Network Peering Connection: error 503 server", strings.ToLower(string(errType)))

			got := err.FillMessage(&Config{
				Type:  errType,
				Error: "error 503 server",
			})

			if diff := deep.Equal(got, expected); diff != nil {
				t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
			}
		})
	}

	for _, errType := range errTypes {
		errType := errType
		t.Run(string(errType), func(t *testing.T) {
			t.Parallel()

			var (
				expected = fmt.Errorf("error %s MyProvider Network Peering Connection: error 503 server", strings.ToLower(string(errType)))
				got      = err.SetType(errType).SetError("error 503 server").ToError()
			)

			if diff := deep.Equal(got, expected); diff != nil {
				t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
			}
		})
	}
}

func TestUsingSettingType(t *testing.T) {
	// With Full configuration error
	expected := errors.New("error setting attribute `vm_id` in TFProvider VM (5456543433545656): nil pointer")

	got := NewErrorMessage(&Config{
		ID:           "5456543433545656",
		ProviderName: "TFProvider",
		ResourceName: "VM",
		Error:        "nil pointer",
		Type:         Setting,
		Attribute:    "vm_id",
	})

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}

	// You can use a predeterminate global variable to set a default attributes
	globarVar := SetProviderName("TFProvider").SetResourceName("VM")

	got = globarVar.FillMessage(&Config{
		ID:        "5456543433545656",
		Error:     "nil pointer",
		Type:      Setting,
		Attribute: "vm_id",
	})

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}

	// Also you can use ToError function to retrieve the message error
	got = globarVar.SetID("5456543433545656").SetError("nil pointer").SetType(Setting).SetAttribute("vm_id").ToError()

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}
}

func TestUsingSomeAttributes_DifferentCRUDTypes(t *testing.T) {
	var (
		errTypes = []errType{Creating, Reading, Updating, Deleting}

		// You can create a template setting some attribute to be used to create other errors
		err = SetProviderName("MyProvider").SetResourceName("Network Peering Connection")
	)

	for _, errType := range errTypes {
		errType := errType
		t.Run(string(errType), func(t *testing.T) {
			t.Parallel()

			expected := fmt.Errorf("error %s MyProvider Network Peering Connection: error", strings.ToLower(string(errType)))

			got := err.FillMessage(&Config{
				Error: "error",
				Type:  errType,
			})

			if diff := deep.Equal(got, expected); diff != nil {
				t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
			}
		})
	}

	for _, errType := range errTypes {
		errType := errType
		t.Run(string(errType), func(t *testing.T) {
			t.Parallel()

			var (
				expected = fmt.Errorf("error %s MyProvider Network Peering Connection: error", strings.ToLower(string(errType)))
				got      = err.SetType(errType).SetError("error").ToError()
			)

			if diff := deep.Equal(got, expected); diff != nil {
				t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
			}
		})
	}
}

func TestUsingSomeAttributes_SettingType(t *testing.T) {
	// With Full configuration error
	expected := errors.New("error setting an attribute in TFProvider VM: nil pointer")

	got := NewErrorMessage(&Config{
		ProviderName: "TFProvider",
		ResourceName: "VM",
		Error:        "nil pointer",
		Type:         Setting,
	})

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}

	// You can use a predeterminate global variable to set a default attributes
	globarVar := SetProviderName("TFProvider").SetResourceName("VM")

	got = globarVar.FillMessage(&Config{
		Error: "nil pointer",
		Type:  Setting,
	})

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}

	// Also you can use ToError function to retrieve the message error
	got = globarVar.SetError("nil pointer").SetType(Setting).ToError()

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}
}