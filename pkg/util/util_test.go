package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenRandomString(t *testing.T) {
	assert.Equal(t, len(GenRandomString(16)), 16, "GenRandomString(16)")

	//lint:ignore SA4000 test
	assert.Equal(t, GenRandomString(16) != GenRandomString(16), true, "GenRandomString not equal")
}

func TestGenID(t *testing.T) {
	assert.Equal(t, len(GenID()), 16, "GenID")
}

func TestPostgresArray(t *testing.T) {

	arr, err := DecodePostgresArray("{a,b,c}")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, arr, []string{"a", "b", "c"}, "DecodePostgresArray")
	assert.Equal(t, EncodePostgresArray(arr), "{a,b,c}", "DecodePostgresArray")
}

func TestIntersectionStrings(t *testing.T) {

	assert.Equal(t, []string{"a", "b"}, IntersectionStrings([]string{"a", "b", "c"}, []string{"a", "b"}))
	assert.Equal(t, []string{"b"}, IntersectionStrings([]string{"a", "b", "c"}, []string{"a", "b"}, []string{"b", "c"}))
	assert.Equal(t, []string{}, IntersectionStrings([]string{"a", "c"}, []string{"b"}, []string{"a"}))
}

func TestGetFirstStringInMap(t *testing.T) {
	assert.Nil(t, GetFirstStringInMap(map[string]string{}))
	assert.Nil(t, GetFirstStringInMap(nil))
	assert.Equal(t, "b", GetFirstStringInMap(map[string]string{
		"a": "b",
	}))
	assert.Equal(t, "en", GetFirstStringInMap(map[string]string{
		"en": "en",
		"vi": "vi",
	}))
}

func TestIsWebBrowserAgent(t *testing.T) {
	assert.False(t, IsWebBrowserAgent(""))
	assert.True(t, IsWebBrowserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/77.0.3865.90 Safari/537.36"))
	assert.True(t, IsWebBrowserAgent("Mozilla/5.0 (Android; Mobile; rv:13.0) Gecko/13.0 Firefox/13.0"))
	assert.True(t, IsWebBrowserAgent("Mozilla/5.0 (Windows Phone 10.0; Android 6.0.1; Xbox; Xbox One) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Mobile Safari/537.36 Edge/16.16299"))

}
