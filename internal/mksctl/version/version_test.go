package version_test

import (
	"bytes"
	"errors"
	"regexp"
	"testing"
	"time"

	ver "github.com/JulienBreux/mksctl/internal/mksctl/version"

	"github.com/stretchr/testify/assert"
)

func TestVersionDateFailed(t *testing.T) {
	ver.RawDate = "n/a"

	expectedErr := "unable to parse date: n/a"
	_, err := ver.Date()

	assert.Error(t, err, expectedErr)
	assert.Equal(t, expectedErr, err.Error())

	var expectedErrType = err.(*ver.DateParseError)
	assert.True(t, errors.As(err, &expectedErrType))
	assert.Equal(
		t,
		expectedErrType.Unwrap().Error(),
		"parsing time \"n/a\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"n/a\" as \"2006\"",
	)
}

func TestVersionDateSuccess(t *testing.T) {
	ver.RawDate = "1987-01-16T09:00:00Z"

	d, err := ver.Date()

	assert.NoError(t, err)
	assert.Equal(t, d.Year(), 1987)
	assert.Equal(t, d.Month(), time.January)
	assert.Equal(t, d.Day(), 16)
}

func TestPrintVersionJSON(t *testing.T) {
	var r *regexp.Regexp
	w := &bytes.Buffer{}

	ver.Print(w, false, "json")
	r = regexp.MustCompile(`{"version_server":"n/a","version_cli":"dev","commit":"n/a","date":"[0-9T:Z-]+"}`)
	assert.Regexp(t, r, w.String())
}

func TestPrintVersionYAML(t *testing.T) {
	var r *regexp.Regexp
	w := &bytes.Buffer{}

	ver.Print(w, false, "yaml")
	r = regexp.MustCompile(`versionServer: n/a\nversionCLI: dev\ncommit: n/a\ndate: "[0-9T:Z-]+"\n`)
	assert.Regexp(t, r, w.String())
}

func TestPrintVersionText(t *testing.T) {
	var r *regexp.Regexp
	w := &bytes.Buffer{}

	ver.Print(w, false, "")
	r = regexp.MustCompile(`versionServer: n/a\nversionCLI: dev\ncommit: n/a\ndate: "[0-9T:Z-]+"\n`)
	assert.Regexp(t, r, w.String())
}
