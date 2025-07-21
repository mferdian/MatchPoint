package helpers

import "testing"

func TestIsValidEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Valid email with domain",
			args: args{email: "user@example.com"},
			want: true,
		},
		{
			name: "Valid email with subdomain",
			args: args{email: "user@mail.example.co.id"},
			want: true,
		},
		{
			name: "Invalid email missing @",
			args: args{email: "userexample.com"},
			want: false,
		},
		{
			name: "Invalid email missing domain",
			args: args{email: "user@"},
			want: false,
		},
		{
			name: "Invalid email with space",
			args: args{email: "user @example.com"},
			want: false,
		},
		{
			name: "Invalid email with special character",
			args: args{email: "user!@example.com"},
			want: false,
		},
		{
			name: "Empty email",
			args: args{email: ""},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidEmail(tt.args.email); got != tt.want {
				t.Errorf("IsValidEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}
