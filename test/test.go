package test

import (
	strongpassword "strong_password/app/usecase/strong_password"
	"strong_password/tu"
	"testing"
)

func TestCalculateSteps(t *testing.T) {
	tests := []struct {
		password  string
		wantSteps int
	}{
		{"aA1", 3},
		{"1445D1cd", 0},
		{"aaa", 4},
		{"", 6},
		{"A1b", 3},
		{"aA1aa", 1},
		{"aA1aaA1aaA1aaA1aaA1", 0},
		{"aA1aaA1aaA1aaA1aaA1a1", 1},
		{"111111111", 6},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.password, func(t *testing.T) {
			t.Parallel()
			tc := tu.Setup()
			defer tc.Teardown()

			req := strongpassword.StrongPasswordReq{
				InitPassword: tt.password,
			}

			got, err := strongpassword.StrongPasswordSteps(tc.Ctx(), req, tc.DB)
			if err != nil {
				t.Fatalf("StrongPasswordSteps(%v) returned error: %v", tt.password, err)
			}
			if got.NumOfSteps != tt.wantSteps {
				t.Errorf("CalculateSteps(%v) = %v, want %v", tt.password, got.NumOfSteps, tt.wantSteps)
			}
		})
	}
}
