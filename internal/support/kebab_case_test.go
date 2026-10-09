package support_test

import (
	"testing"

	"github.com/hettiger/clg/internal/support"
	"github.com/stretchr/testify/require"
)

func TestKebabCase(t *testing.T) {
	tests := []struct {
		str  string
		want string
	}{
		{str: "foo", want: "foo"},
		{str: "foo-bar", want: "foo-bar"},
		{str: "foo_bar", want: "foo-bar"},
		{str: "fooBar", want: "foo-bar"},
		{str: "FooBar", want: "foo-bar"},
		{str: "föö", want: "f"},
		{str: "fööBar", want: "f-bar"},
		{str: "   foo", want: "foo"},
		{str: "   foo   ", want: "foo"},
		{str: "foo   ", want: "foo"},
		{str: "öööö", want: ""},
		{str: "fööööBar", want: "f-bar"},
		{str: "HTTPServer", want: "http-server"},
		{str: "HTTPS€rver", want: "http-s-rver"},
		{str: "HTTP-Server", want: "http-server"},
		{str: "HTTP", want: "http"},
		{str: "HTTP2", want: "http2"},
		{str: "../../escaped", want: "escaped"},
		{str: "", want: ""},
		{str: "---___...", want: ""},
		{str: "FOO ", want: "foo"},
		{str: "HTTP-server", want: "http-server"},
		{str: "fooBAR", want: "foo-bar"},
		{str: "myHTTPServer", want: "my-http-server"},
		{str: "HTTP2Server", want: "http2-server"},
		{str: "HTTPSörver", want: "http-s-rver"},
		{str: "HTTPS🙂rver", want: "http-s-rver"},
		{str: `..\..\escaped`, want: "escaped"},
		{str: `C:\temp\entry`, want: "c-temp-entry"},
		{str: "key:stream", want: "key-stream"},
		{str: "foo\x00bar", want: "foo-bar"},
		{str: "K", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			got := support.KebabCase(tt.str)

			require.Equal(t, tt.want, got)
		})
	}
}
