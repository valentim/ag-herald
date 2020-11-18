package slack

import (
	"io/ioutil"
	"testing"
)

func TestProjectConfigurationContent(t *testing.T) {

	json, err := ioutil.ReadFile("project_configuration_fixture.json")

	if err != nil {
		t.Error("Reading json file error")
	}

	projectConfigurationContent := ProjectConfigurationContent(json)

	t.Log(projectConfigurationContent)
}
