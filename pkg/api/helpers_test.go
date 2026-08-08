package api

import "testing"

func TestValidVersion(t *testing.T) {
	// cases are tied to oldestValidClientVersion in api.go
	tests := []struct {
		version string
		want    bool
	}{
		{"0.5.9", false},
		{"0.6.0", true},
		{"0.6.1", true},
		{"0.7.0", true},
		{"1.0.0", true},
		{"0.0.9", false},
	}
	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			got, err := validVersion(tt.version)
			if err != nil {
				t.Errorf("got error: %v", err.Error())
			}
			if got != tt.want {
				t.Errorf("got %v; want %v", got, tt.want)
			}
		})
	}
}

// func TestUpdateAvailable(t *testing.T) {
// 	tests := []struct {
// 		version string
// 		want    bool
// 	}{
// 		{"0.4.5", false},
// 		{"0.4.4", true},
// 		{"0.4.0", true},
// 		{"0.3.0", true},
// 		{"0.3.1", true},
// 		{"1.3.1", false},
// 		{"0.0.9", true},
// 		{"0.4.6", false},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.version, func(t *testing.T) {
// 			got, err := updateAvailable(tt.version)
// 			if err != nil {
// 				t.Errorf("got error: %v", err.Error())
// 			}
// 			if got != tt.want {
// 				t.Errorf("got %v; want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
