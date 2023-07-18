package url_test

import (
	"testing"

	"github.com/JulienBreux/mksctl/pkg/validator/url"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestBadURL(t *testing.T) {
	cmd := &cobra.Command{}
	args := []string{"microcks"}

	v := url.ValidArg(0)
	err := v(cmd, args)

	assert.Error(t, err)
}

func TestBadArg(t *testing.T) {
	cmd := &cobra.Command{}
	args := []string{}

	v := url.ValidArg(0)
	err := v(cmd, args)

	assert.Error(t, err)
}

func TestGoodURL(t *testing.T) {
	cmd := &cobra.Command{}
	args := []string{"https://microcks.io/"}

	v := url.ValidArg(0)
	err := v(cmd, args)

	assert.NoError(t, err)
}
