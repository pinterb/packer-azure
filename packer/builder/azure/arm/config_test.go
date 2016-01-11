// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See the LICENSE file in the project root for license information.

package arm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

// List of configuration parameters that are required by the ARM builder.
var requiredConfigValues = []string{
	"capture_name_prefix",
	"capture_container_name",
	"client_id",
	"client_secret",
	"image_offer",
	"image_publisher",
	"image_sku",
	"location",
	"storage_account",
	"subscription_id",
	"tenant_id",
}

func TestConfigShouldProvideReasonableDefaultValues(t *testing.T) {
	c, _, err := newConfig(getArmBuilderConfiguration(), getPackerConfiguration())

	if err != nil {
		t.Errorf("Expected configuration creation to succeed, but it failed!\n")
		t.Fatalf(" errors: %s\n", err)
	}

	if c.UserName == "" {
		t.Errorf("Expected 'UserName' to be populated, but it was empty!")
	}

	if c.VMSize == "" {
		t.Errorf("Expected 'VMSize' to be populated, but it was empty!")
	}
}

func TestConfigShouldDefaultVMSizeToStandardA1(t *testing.T) {
	c, _, _ := newConfig(getArmBuilderConfiguration(), getPackerConfiguration())

	if c.VMSize != "Standard_A1" {
		t.Errorf("Expected 'VMSize' to default to 'Standard_A1', but got '%s'.", c.VMSize)
	}
}

func TestUserShouldProvideRequiredValues(t *testing.T) {
	builderValues := make(map[string]string)
	var lastElement = requiredConfigValues[len(requiredConfigValues)-1]

	for _, v := range requiredConfigValues {
		builderValues[v] = "--ignored--"

		_, _, err := newConfig(getArmBuilderConfigurationFromMap(builderValues), getPackerConfiguration())
		if lastElement != v && err == nil {
			t.Errorf("Expected configuration creation to fail, but it succeeded!\n")
			t.Fatalf(" -> %+v\n", builderValues)
		}
	}

	_, _, err := newConfig(getArmBuilderConfigurationFromMap(builderValues), getPackerConfiguration())
	if err != nil {
		t.Errorf("Expected configuration creation to succeed, but it failed!")
		t.Errorf(" errors: %s", err)
		t.Fatalf(" -> %+v", builderValues)
	}
}

func TestSystemShouldDefineRuntimeValues(t *testing.T) {
	c, _, _ := newConfig(getArmBuilderConfiguration(), getPackerConfiguration())
	if c.tmpComputeName == "" {
		t.Errorf("Expected tmpComputeName to not be empty, but it was '%s'!", c.tmpComputeName)
	}

	if c.tmpDeploymentName == "" {
		t.Errorf("Expected tmpDeploymentName to not be empty, but it was '%s'!", c.tmpDeploymentName)
	}

	if c.tmpResourceGroupName == "" {
		t.Errorf("Expected tmpResourceGroupName to not be empty, but it was '%s'!", c.tmpResourceGroupName)
	}

	if c.tmpOSDiskName == "" {
		t.Errorf("Expected tmpOSDiskName to not be empty, but it was '%s'!", c.tmpOSDiskName)
	}
}

func TestConfigShouldTransformToTemplateParameters(t *testing.T) {
	c, _, _ := newConfig(getArmBuilderConfiguration(), getPackerConfiguration())
	templateParameters := c.toTemplateParameters()

	if templateParameters.AdminUsername.Value != c.UserName {
		t.Errorf("Expected AdminUsername to be equal to config's AdminUsername, but they were '%s' and '%s' respectively.", templateParameters.AdminUsername.Value, c.UserName)
	}

	if templateParameters.DnsNameForPublicIP.Value != c.tmpComputeName {
		t.Errorf("Expected DnsNameForPublicIP to be equal to config's DnsNameForPublicIP, but they were '%s' and '%s' respectively.", templateParameters.DnsNameForPublicIP.Value, c.tmpComputeName)
	}

	if templateParameters.ImageOffer.Value != c.ImageOffer {
		t.Errorf("Expected ImageOffer to be equal to config's ImageOffer, but they were '%s' and '%s' respectively.", templateParameters.ImageOffer.Value, c.ImageOffer)
	}

	if templateParameters.ImagePublisher.Value != c.ImagePublisher {
		t.Errorf("Expected ImagePublisher to be equal to config's ImagePublisher, but they were '%s' and '%s' respectively.", templateParameters.ImagePublisher.Value, c.ImagePublisher)
	}

	if templateParameters.ImageSku.Value != c.ImageSku {
		t.Errorf("Expected ImageSku to be equal to config's ImageSku, but they were '%s' and '%s' respectively.", templateParameters.ImageSku.Value, c.ImageSku)
	}

	if templateParameters.Location.Value != c.Location {
		t.Errorf("Expected Location to be equal to config's Location, but they were '%s' and '%s' respectively.", templateParameters.Location.Value, c.Location)
	}

	if templateParameters.OSDiskName.Value != c.tmpOSDiskName {
		t.Errorf("Expected OSDiskName to be equal to config's OSDiskName, but they were '%s' and '%s' respectively.", templateParameters.OSDiskName.Value, c.tmpOSDiskName)
	}

	if templateParameters.StorageAccountName.Value != c.StorageAccount {
		t.Errorf("Expected StorageAccountName to be equal to config's StorageAccountName, but they were '%s' and '%s' respectively.", templateParameters.StorageAccountName.Value, c.StorageAccount)
	}

	if templateParameters.VMName.Value != c.tmpComputeName {
		t.Errorf("Expected VMName to be equal to config's VMName, but they were '%s' and '%s' respectively.", templateParameters.VMName.Value, c.tmpComputeName)
	}

	if templateParameters.VMSize.Value != c.VMSize {
		t.Errorf("Expected VMSize to be equal to config's VMSize, but they were '%s' and '%s' respectively.", templateParameters.VMSize.Value, c.VMSize)
	}
}

func TestConfigShouldTransformToVirtualMachineCaptureParameters(t *testing.T) {
	c, _, _ := newConfig(getArmBuilderConfiguration(), getPackerConfiguration())
	parameters := c.toVirtualMachineCaptureParameters()

	if *parameters.DestinationContainerName != c.CaptureContainerName {
		t.Errorf("Expected DestinationContainerName to be equal to config's CaptureContainerName, but they were '%s' and '%s' respectively.", *parameters.DestinationContainerName, c.CaptureContainerName)
	}

	if *parameters.VhdPrefix != c.CaptureNamePrefix {
		t.Errorf("Expected DestinationContainerName to be equal to config's CaptureContainerName, but they were '%s' and '%s' respectively.", *parameters.VhdPrefix, c.CaptureNamePrefix)
	}

	if *parameters.OverwriteVhds != false {
		t.Error("Expected OverwriteVhds to be false, but it was not.")
	}
}

func getArmBuilderConfiguration() interface{} {
	m := make(map[string]string)
	for _, v := range requiredConfigValues {
		m[v] = fmt.Sprintf("%s00", v)
	}

	return getArmBuilderConfigurationFromMap(m)
}

func getArmBuilderConfigurationFromMap(kvp map[string]string) interface{} {
	bs := bytes.NewBufferString("{")

	for k, v := range kvp {
		bs.WriteString(fmt.Sprintf("\"%s\": \"%s\",\n", k, v))
	}

	// remove the trailing ",\n" because JSON
	bs.Truncate(bs.Len() - 2)
	bs.WriteString("}")

	var config interface{}
	json.Unmarshal([]byte(bs.String()), &config)

	return config
}

func getPackerConfiguration() interface{} {
	var doc = `{
		"packer_user_variables": {
			"sa": "my_storage_account"
		},
		"packer_build_name": "azure-arm-vm",
		"packer_builder_type": "azure-arm-vm",
		"packer_debug": "false",
		"packer_force": "false",
		"packer_template_path": "/home/jenkins/azure-arm-vm/template.json"
	}`

	var config interface{}
	json.Unmarshal([]byte(doc), &config)

	return config
}
