// nolint: dupl
package appcli

import (
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/urfave/cli"
)

func TestNewFlagApp(t *testing.T) {
	cfgText := []byte(`options:
  foo:
    short: f
    default: foovalue

tasks:
  mytask:
    run: echo ${foo}
`)

	flagApp, err := newFlagApp(cfgText)
	if err != nil {
		t.Fatalf(
			"newFlagApp():\nconfig: `%s`\nunexpected err: %s",
			string(cfgText), err,
		)
	}

	args := []string{"tusk", "mytask", "--foo", "other"}
	if err = flagApp.Run(args); err != nil {
		t.Fatalf(
			"flagApp.Run():\nconfig: `%s`\nunexpected err: %s",
			string(cfgText), err,
		)
	}

	command, ok := flagApp.Metadata["command"].(*cli.Command)
	if !ok {
		t.Fatalf(
			"flagApp.Metadata:\nconfig: `%s`\nMetadata command not a *cli.Command: %#v",
			string(cfgText), flagApp.Metadata["command"],
		)
	}

	commandName := command.Name
	commandExpected := "mytask"

	if commandExpected != commandName {
		t.Errorf(
			"flagApp.Metadata[\"command\"] for args(%s):\n expected: %s\nactual: %s",
			args, commandExpected, commandName,
		)
	}

	flagsActual, ok := flagApp.Metadata["flagsPassed"].(map[string]string)
	if !ok {
		t.Fatalf(
			"flagApp.Metadata:\nconfig: `%s`\nMetadata flagsPassed not a map: %#v",
			string(cfgText), flagApp.Metadata["flagsPassed"],
		)
	}

	flagsExpected := map[string]string{
		"foo": "other",
	}

	if !reflect.DeepEqual(flagsExpected, flagsActual) {
		t.Errorf(
			"flagApp.Metadata for args(%s):\n expected: %#v\nactual: %#v",
			args, flagsExpected, flagsActual,
		)
	}
}

func TestNewFlagApp_no_options(t *testing.T) {
	cfgText := []byte(`tasks:
  mytask:
    run: echo foo
`)

	flagApp, err := newFlagApp(cfgText)
	if err != nil {
		t.Fatalf(
			"newFlagApp():\nconfig: `%s`\nunexpected err: %s",
			string(cfgText), err,
		)
	}

	args := []string{"tusk", "mytask"}
	if err = flagApp.Run(args); err != nil {
		t.Fatalf(
			"flagApp.Run():\nconfig: `%s`\nunexpected err: %s",
			string(cfgText), err,
		)
	}

	command, ok := flagApp.Metadata["command"].(*cli.Command)
	if !ok {
		t.Fatalf(
			"flagApp.Metadata:\nconfig: `%s`\nMetadata command not a *cli.Command: %#v",
			string(cfgText), flagApp.Metadata["command"],
		)
	}

	commandName := command.Name
	commandExpected := "mytask"

	if commandExpected != commandName {
		t.Errorf(
			"flagApp.Metadata[\"command\"] for args(%s):\n expected: %s\nactual: %s",
			args, commandExpected, commandName,
		)
	}

	flagsActual, ok := flagApp.Metadata["flagsPassed"].(map[string]string)
	if !ok {
		t.Fatalf(
			"flagApp.Metadata:\nconfig: `%s`\nMetadata flagsPassed not a map: %#v",
			string(cfgText), flagApp.Metadata["flagsPassed"],
		)
	}

	flagsExpected := map[string]string{}

	if !reflect.DeepEqual(flagsExpected, flagsActual) {
		t.Errorf(
			"flagApp.Metadata for args(%s):\n expected: %#v\nactual: %#v",
			args, flagsExpected, flagsActual,
		)
	}
}

func TestGetConfigMetadata_defaults(t *testing.T) {
	args := []string{"tusk"}

	metadata, err := GetConfigMetadata(args)
	if err != nil {
		t.Fatalf(
			"GetConfigMetadata(%s): unexpected err: %s",
			args, err,
		)
	}

	// The project's tuskfile should be found in the project root.
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd(): unexpected err: %s", err)
	}

	directory := path.Dir(wd)
	if metadata.Directory != directory {
		t.Errorf(
			"GetConfigMetadata(%s): expected Directory: %s, actual: %s",
			args, directory, metadata.Directory,
		)
	}

	if metadata.PrintVersion {
		t.Errorf(
			"GetConfigMetadata(%s): expected RunVersion: false, actual: true",
			args,
		)
	}

	if metadata.Quiet {
		t.Errorf(
			"GetConfigMetadata(%s): expected Quiet: false, actual: true",
			args,
		)
	}

	if metadata.Verbose {
		t.Errorf(
			"GetConfigMetadata(%s): expected Verbose: false, actual: true",
			args,
		)
	}
}

func TestGetConfigMetadata_file(t *testing.T) {
	cfgPath := "testdata/example.yml"
	args := []string{"tusk", "--file", cfgPath}

	metadata, err := GetConfigMetadata(args)
	if err != nil {
		t.Fatalf(
			"GetConfigMetadata(%s): unexpected err: %s",
			args, err,
		)
	}

	directory := "testdata"

	if directory != metadata.Directory {
		t.Errorf(
			"GetConfigMetadata(%s): expected Directory: %s, actual: %s",
			args, directory, metadata.Directory,
		)
	}

	cfgText, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		t.Fatalf(
			"ioutil.ReadFile(%s): unexpected err: %s",
			cfgPath, err,
		)
	}

	expected := string(cfgText)
	actual := string(metadata.CfgText)

	if expected != actual {
		t.Errorf(
			"GetConfigMetadata(%s):\nexpected config text: %s\nactual: %s",
			args, expected, actual,
		)
	}
}

func TestGetConfigMetadata_fileNoExist(t *testing.T) {
	args := []string{"tusk", "--file", "fakefile.yml"}

	_, err := GetConfigMetadata(args)
	if !os.IsNotExist(err) {
		t.Errorf(
			"GetConfigMetadata(%s): unexpected err: os.IsNotExist, actual: %s",
			args, err,
		)
	}
}

func TestGetConfigMetadata_version(t *testing.T) {
	args := []string{"tusk", "--version"}

	metadata, err := GetConfigMetadata(args)
	if err != nil {
		t.Fatalf(
			"GetConfigMetadata(%s):\nunexpected err: %s",
			args, err,
		)
	}

	if !metadata.PrintVersion {
		t.Errorf(
			"GetConfigMetadata(%s): expected RunVersion: true, actual: false",
			args,
		)
	}
}

func TestGetConfigMetadata_quiet(t *testing.T) {
	args := []string{"tusk", "--quiet"}

	metadata, err := GetConfigMetadata(args)
	if err != nil {
		t.Fatalf(
			"GetConfigMetadata(%s):\nunexpected err: %s",
			args, err,
		)
	}

	if !metadata.Quiet {
		t.Errorf(
			"GetConfigMetadata(%s): expected Quiet: true, actual: false",
			args,
		)
	}
}

func TestGetConfigMetadata_verbose(t *testing.T) {
	args := []string{"tusk", "--verbose"}

	metadata, err := GetConfigMetadata(args)
	if err != nil {
		t.Fatalf(
			"GetConfigMetadata(%s):\nunexpected err: %s",
			args, err,
		)
	}

	if !metadata.Verbose {
		t.Errorf(
			"GetConfigMetadata(%s): expected Verbose: true, actual: false",
			args,
		)
	}
}
