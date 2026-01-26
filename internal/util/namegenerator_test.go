package util

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetFirstName(t *testing.T) {
	// Test that we get a valid name
	name := GetFirstName()
	require.NotEmpty(t, name, "Name should not be empty")

	// Name should start with uppercase
	assert.True(t, name[0] >= 'A' && name[0] <= 'Z', "Name should start with uppercase letter")

	// Name should contain only lowercase letters after first character
	matched, err := regexp.MatchString(`^[A-Z][a-z]+$`, name)
	require.NoError(t, err)
	assert.True(t, matched, "Name should match pattern: uppercase letter followed by lowercase letters")
}

func TestGetFirstName_Uniqueness(t *testing.T) {
	// Generate multiple names and check they're not all the same
	names := make(map[string]bool)
	for i := 0; i < 100; i++ {
		name := GetFirstName()
		names[name] = true
	}

	// With 100 random names, we should have some variety
	// (though it's theoretically possible to have duplicates)
	assert.Greater(t, len(names), 1, "Should generate different names")
}

func TestGetFirstName_Length(t *testing.T) {
	// Generate multiple names and check reasonable length
	// Based on genFirstName with syllableCount=5, each syllable is 2 chars
	// So minimum should be around 10 chars (5 syllables * 2)
	for i := 0; i < 50; i++ {
		name := GetFirstName()
		assert.GreaterOrEqual(t, len(name), 4, "Name should be at least 4 characters (2 syllables minimum)")
		assert.LessOrEqual(t, len(name), 20, "Name should be reasonable length (max 5 syllables * 4 chars)")
	}
}
