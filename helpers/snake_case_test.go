package helpers

import "testing"

func TestSnakeCase(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"camelCase", args{"camelCase"}, "camel_case"},
		{"PascalCase", args{"PascalCase"}, "pascal_case"},
		{"already_snake_case", args{"already_snake_case"}, "already_snake_case"},
		{"empty string", args{""}, ""},
		{"all lowercase", args{"lowercase"}, "lowercase"},
		{"with number", args{"Field1Value"}, "field1_value"},
		{"multiple uppercase", args{"JSONData"}, "json_data"},
		{"ends with uppercase", args{"fieldA"}, "field_a"},
		{"consecutive capitals", args{"MyHTTPServer"}, "my_http_server"},
		{"single char", args{"A"}, "a"},
		{"no capitals", args{"abc123"}, "abc123"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SnakeCase(tt.args.str); got != tt.want {
				t.Errorf("SnakeCase() = %v, want %v", got, tt.want)
			}
		})
	}
}
