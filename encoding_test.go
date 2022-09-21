//go:build !integration
// +build !integration

package valuefirst

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
	}{
		{
			name: "encode special characters, *, #, %, <, >, +",
			args: args{
				"* # % < > +",
			},
			wantRes: "%2A %23 %25 %3C %3E %2B",
		},
		{
			name: "encode below the 32",
			args: args{
				msg: "hello: " + string(rune(17)) + " world",
			},
			wantRes: "hello: %11 world",
		},
		{
			name: "encode after the 128",
			args: args{
				msg: "hello: " + string(rune(137)) + " world",
			},
			wantRes: "hello: %89 world",
		},
		{
			name: "encode below the 16 - 1",
			args: args{
				msg: "hello: " + string(rune(10)) + "Test",
			},
			wantRes: "hello: %0ATest",
		},
		{
			name: "encode below the 16 - 2",
			args: args{
				msg: "hello: " + string(rune(15)) + "Test",
			},
			wantRes: "hello: %0FTest",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantRes, Encode(tt.args.msg), "Encode(%v)", tt.args.msg)
		})
	}
}

func TestEncodeHTML(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
	}{
		{
			name: "encode HTML characters, &#39;, &#32;, &#34;, &#62;, &#60;, &#13;, &#10;, &#9;",
			args: args{
				msg: "&#39; &#32; &#34; &#62; &#60; &#13; &#10; &#9;",
			},
			wantRes: "&apos&#032&#032&#032&quot&#032&gt&#032&lt&#032&#013&#032&#010&#032&#009",
		},
		{
			name: "encode HTML characters, &#39;, &#32;, &#34;, &#62;, &#60;, &#13;, &#10;, &#9;",
			args: args{
				msg: "The value of x can be represented as, x > 128 and x < 32.",
			},
			wantRes: "The&#032value&#032of&#032x&#032can&#032be&#032represented&#032as,&#032x&#032&gt&#032128&#032and&#032x&#032&lt&#03232.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.wantRes, EncodeHTML(tt.args.msg), "EncodeHTML(%v)", tt.args.msg)
		})
	}
}
