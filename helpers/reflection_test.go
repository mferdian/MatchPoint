package helpers

import (
	"reflect"
	"testing"
)

type TestStruct struct {
	Name    string
	Age     int
	Email   string
	Active  bool
	Score   float64
	private string
}

type NestedStruct struct {
	User   TestStruct
	Config map[string]interface{}
}

func TestGetFieldValue(t *testing.T) {
	testUser := TestStruct{
		Name:    "John Doe",
		Age:     30,
		Email:   "john@example.com",
		Active:  true,
		Score:   95.5,
		private: "hidden",
	}

	nestedData := NestedStruct{
		User: testUser,
		Config: map[string]interface{}{
			"theme": "dark",
			"lang":  "en",
		},
	}

	type args struct {
		data  any
		field string
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		// Test cases untuk struct dengan berbagai tipe data
		{
			name: "get string field from struct",
			args: args{
				data:  testUser,
				field: "Name",
			},
			want:    "John Doe",
			wantErr: false,
		},
		{
			name: "get int field from struct",
			args: args{
				data:  testUser,
				field: "Age",
			},
			want:    30,
			wantErr: false,
		},
		{
			name: "get bool field from struct",
			args: args{
				data:  testUser,
				field: "Active",
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "get float field from struct",
			args: args{
				data:  testUser,
				field: "Score",
			},
			want:    95.5,
			wantErr: false,
		},
		
		// Test cases untuk pointer ke struct
		{
			name: "get field from struct pointer",
			args: args{
				data:  &testUser,
				field: "Email",
			},
			want:    "john@example.com",
			wantErr: false,
		},
		
		// Test cases untuk nested struct
		{
			name: "get nested struct field",
			args: args{
				data:  nestedData,
				field: "User",
			},
			want:    testUser,
			wantErr: false,
		},
		
		// Test cases untuk map - tidak didukung, akan error
		{
			name: "get value from map should fail",
			args: args{
				data: map[string]interface{}{
					"name": "Jane",
					"age":  25,
				},
				field: "name",
			},
			want:    nil,
			wantErr: true,
		},
		
		// Error test cases
		{
			name: "field not found in struct",
			args: args{
				data:  testUser,
				field: "NonExistentField",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unsupported data type (map)",
			args: args{
				data: map[string]interface{}{
					"name": "Jane",
				},
				field: "name",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "nil data",
			args: args{
				data:  nil,
				field: "Name",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty field name",
			args: args{
				data:  testUser,
				field: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unsupported data type (slice)",
			args: args{
				data:  []string{"a", "b", "c"},
				field: "0",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "unsupported data type (string)",
			args: args{
				data:  "just a string",
				field: "length",
			},
			want:    nil,
			wantErr: true,
		},
		
		// Test cases untuk private field (jika fungsi mendukung)
		{
			name: "access private field should fail",
			args: args{
				data:  testUser,
				field: "private",
			},
			want:    nil,
			wantErr: true,
		},
		
		// Test cases untuk nested struct field yang berisi map
		{
			name: "get nested Config map field",
			args: args{
				data:  nestedData,
				field: "Config",
			},
			want: map[string]interface{}{
				"theme": "dark",
				"lang":  "en",
			},
			wantErr: false,
		},
		
		// Test cases untuk interface{}
		{
			name: "get field from interface{} containing struct",
			args: args{
				data:  interface{}(testUser),
				field: "Name",
			},
			want:    "John Doe",
			wantErr: false,
		},
		
		// Test cases untuk zero values
		{
			name: "get zero value field",
			args: args{
				data: TestStruct{
					Name: "",
					Age:  0,
				},
				field: "Age",
			},
			want:    0,
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFieldValue(tt.args.data, tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFieldValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFieldValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Test helper untuk memverifikasi tipe data yang dikembalikan
func TestGetFieldValueTypes(t *testing.T) {
	testUser := TestStruct{
		Name:   "John Doe",
		Age:    30,
		Active: true,
		Score:  95.5,
	}
	
	tests := []struct {
		field        string
		expectedType string
	}{
		{"Name", "string"},
		{"Age", "int"},
		{"Active", "bool"},
		{"Score", "float64"},
	}
	
	for _, tt := range tests {
		t.Run("type_"+tt.field, func(t *testing.T) {
			got, err := GetFieldValue(testUser, tt.field)
			if err != nil {
				t.Fatalf("GetFieldValue() error = %v", err)
			}
			
			actualType := reflect.TypeOf(got).String()
			if actualType != tt.expectedType {
				t.Errorf("GetFieldValue() returned type %v, want %v", actualType, tt.expectedType)
			}
		})
	}
}