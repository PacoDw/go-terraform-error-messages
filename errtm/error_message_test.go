package errtm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-test/deep"
)

func TestBasicMessageError(t *testing.T) {
	var (
		ms *Config = SaveProviderName("MyProvider")

		got      = ms.Error()
		expected = "error in MyProvider"
	)

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}

	got = ms.SetState(Setting).Error()
	expected = "error setting an attribute in MyProvider"

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}

	got = ms.SetError("this is my error").Error()
	expected = "error in MyProvider: this is my error"

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}
}

func TestUsingNewErrorMessage_DifferentCRUDTypes(t *testing.T) {
	errStates := []errState{Creating, Reading, Updating, Deleting}

	for _, errState := range errStates {
		errState := errState
		t.Run(string(errState), func(t *testing.T) {
			t.Parallel()

			expected := fmt.Sprintf("error %s TerraformProvider PeeringConnection (5456543433545656): Error processing your request", strings.ToLower(string(errState)))

			got := NewErrorMessage(&Config{
				ID:           "5456543433545656",
				ProviderName: "TerraformProvider",
				ResourceName: "PeeringConnection",
				ErrorMessage: "Error processing your request",
				State:        errState,
			}).Error()

			if diff := deep.Equal(got, expected); diff != nil {
				t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
			}
		})
	}
}

func TestUsingSetProviderName_DifferentCRUDTypes(t *testing.T) {
	var (
		errTypes = []errState{Creating, Reading, Updating, Deleting}

		// You can create a template setting some attribute to be used to create other errors
		// Noted the SaveProviderName func will do but SetResourceName will save the value
		// temporaly, so this means if it uses:  err.SetState(errState).SetError("error 503 server")
		// the error will be: error `creating` MyProvider: error 503 server
		// check the above second test
		err = SaveProviderName("MyProvider")
	)

	for _, errState := range errTypes {
		errState := errState
		t.Run(string(errState), func(t *testing.T) {
			expected := fmt.Sprintf("error %s MyProvider Network Peering Connection: error 503 server", strings.ToLower(string(errState)))

			got := err.FillMessage(&Config{
				ResourceName: "Network Peering Connection",
				State:        errState,
				ErrorMessage: "error 503 server",
			}).Error()

			if diff := deep.Equal(got, expected); diff != nil {
				t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
			}
		})
	}

	for _, errState := range errTypes {
		errState := errState
		t.Run(string(errState), func(t *testing.T) {
			var (
				expected = fmt.Sprintf("error %s MyProvider: error 503 server", strings.ToLower(string(errState)))
				got      = err.SetState(errState).SetError("error 503 server").Error()
			)

			if diff := deep.Equal(got, expected); diff != nil {
				t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
			}
		})
	}

	for _, errState := range errTypes {
		errState := errState
		t.Run(string(errState), func(t *testing.T) {
			var (
				expected = fmt.Sprintf("error %s MyProvider Virtual Machine: error 503 server", strings.ToLower(string(errState)))
				got      = err.SetResourceName("Virtual Machine").SetState(errState).SetError("error 503 server").Error()
			)

			if diff := deep.Equal(got, expected); diff != nil {
				t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
			}
		})
	}
}

func TestUsingSettingType(t *testing.T) {
	// With Full configuration error
	expected := "error setting attribute `vm_id` in TFProvider VM (5456543433545656): nil pointer"

	got := NewErrorMessage(&Config{
		ID:           "5456543433545656",
		ProviderName: "TFProvider",
		ResourceName: "VM",
		ErrorMessage: "nil pointer",
		State:        Setting,
		Attribute:    "vm_id",
	}).Error()

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}

	// You can use a predeterminate global variable to set a default attributes
	globarVar := SaveProviderName("TFProvider").SaveResourceName("VM")

	got = globarVar.FillMessage(&Config{
		ID:           "5456543433545656",
		ErrorMessage: "nil pointer",
		State:        Setting,
		Attribute:    "vm_id",
	}).Error()

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}

	// Also you can use ToError function to retrieve the message error
	got = globarVar.SetID("5456543433545656").SetError("nil pointer").SetState(Setting).SetAttribute("vm_id").Error()

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}
}

func TestUsingSomeAttributes_DifferentCRUDTypes(t *testing.T) {
	var (
		errTypes = []errState{Creating, Reading, Updating, Deleting}

		// You can create a template setting some attribute to be used to create other errors
		err = SaveProviderName("MyProvider").SaveResourceName("Network Peering Connection")
	)

	for _, errState := range errTypes {
		errState := errState
		t.Run(string(errState), func(t *testing.T) {
			expected := fmt.Sprintf("error %s MyProvider Network Peering Connection: error", strings.ToLower(string(errState)))

			got := err.FillMessage(&Config{
				ErrorMessage: "error",
				State:        errState,
			}).Error()

			if diff := deep.Equal(got, expected); diff != nil {
				t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
			}
		})
	}

	for _, errState := range errTypes {
		errState := errState
		t.Run(string(errState), func(t *testing.T) {
			var (
				expected = fmt.Sprintf("error %s MyProvider VM: error", strings.ToLower(string(errState)))
				got      = err.SetState(errState).SetResourceName("VM").SetError("error").Error()
			)

			if diff := deep.Equal(got, expected); diff != nil {
				t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
			}
		})
	}
}

func TestUsingSomeAttributes_SettingType(t *testing.T) {
	// With Full configuration error
	expected := "error setting an attribute in TFProvider VM: nil pointer"

	got := NewErrorMessage(&Config{
		ProviderName: "TFProvider",
		ResourceName: "VM",
		ErrorMessage: "nil pointer",
		State:        Setting,
	}).Error()

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}

	// You can use a predeterminate global variable to set a default attributes
	globarVar := SaveProviderName("TFProvider").SaveResourceName("VM")

	got = globarVar.FillMessage(&Config{
		ErrorMessage: "nil pointer",
		State:        Setting,
	}).Error()

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}

	// Also you can use ToError function to retrieve the message error
	got = globarVar.SetError("nil pointer").SetState(Setting).Error()

	if diff := deep.Equal(got, expected); diff != nil {
		t.Errorf("Diff:\n got=%#v\nwant=%#v \n\ndiff=%#v", got, expected, diff)
	}
}
