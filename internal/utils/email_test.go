package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateEmail(t *testing.T) {
	t.Run("Email validates correctly", func(t *testing.T) {
		require.Nil(t, ValidateEmail("myemail@example.com"), "should not have an error")
	})

	t.Run("Email does not validate", func(t *testing.T) {
		require.Errorf(t, ValidateEmail("myemail"), "invalid email format", "should not validate email without @ and .")
		require.Errorf(t, ValidateEmail("myemail@example"), "invalid email format", "should not validate email without .")
		require.Errorf(t, ValidateEmail("myemail.com"), "invalid email format", "should not validate email without @")
	})
}
