package result

import (
	"errors"
	"testing"
)

func TestErr(t *testing.T) {
	testErr := errors.New("test error")
	result := Err[int, error](testErr)

	if !result.IsErr() {
		t.Error("Expected IsErr() to return true")
	}

	if result.IsOk() {
		t.Error("Expected IsOk() to return false")
	}

	_, err := result.Unwrap()
	if err != testErr {
		t.Errorf("Expected error %v, got %v", testErr, err)
	}
}

func TestOk(t *testing.T) {
	result := Ok[int, error](42)

	if !result.IsOk() {
		t.Error("Expected IsOk() to return true")
	}

	if result.IsErr() {
		t.Error("Expected IsErr() to return false")
	}

	val, err := result.Unwrap()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if val != 42 {
		t.Errorf("Expected value 42, got %v", val)
	}
}

func TestUnwrap(t *testing.T) {
	testErr := errors.New("test")
	tests := []struct {
		name   string
		result Result[int, error]
		wantV  int
		wantE  error
	}{
		{
			name:   "Result returns value",
			result: Ok[int, error](42),
			wantV:  42,
			wantE:  nil,
		},
		{
			name:   "Result returns err",
			result: Err[int, error](testErr),
			wantV:  0,
			wantE:  testErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotV, gotE := tt.result.Unwrap()
			if gotV != tt.wantV {
				t.Errorf("Unwrap() value = %v, want %v", gotV, tt.wantV)
			}
			if gotE != tt.wantE {
				t.Errorf("Unwrap() error = %v, want %v", gotE, tt.wantE)
			}
		})
	}
}

func TestUnwrapOr(t *testing.T) {
	testErr := errors.New("test")
	tests := []struct {
		name         string
		result       Result[int, error]
		defaultValue int
		want         int
	}{
		{
			name:         "Result returns value",
			result:       Ok[int, error](42),
			defaultValue: 100,
			want:         42,
		},
		{
			name:         "Error returns default",
			result:       Err[int](testErr),
			defaultValue: 100,
			want:         100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotV := tt.result.UnwrapOr(tt.defaultValue)
			if gotV != tt.want {
				t.Errorf("UnwrapOr() = %v, want %v", gotV, tt.want)
			}
		})
	}
}

func TestUnwrapOrElse(t *testing.T) {
	testErr := errors.New("test")
	tests := []struct {
		name   string
		result Result[int, error]
		fn     func() int
		want   int
	}{
		{
			name:   "Result return value",
			result: Ok[int, error](42),
			fn:     func() int { return 100 },
			want:   42,
		},
		{
			name:   "Error return calls function",
			result: Err[int](testErr),
			fn:     func() int { return 100 },
			want:   100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.result.UnwrapOrElse(tt.fn)
			if v != tt.want {
				t.Errorf("UnwrapOrElse = %v, want %v", v, tt.want)
			}
		})
	}
}

func TestExpect(t *testing.T) {
	testErr := errors.New("original error")
	tests := []struct {
		name      string
		result    Result[int, error]
		expectMsg string
		wantValue int
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "Ok result returns value with no error",
			result:    Ok[int, error](42),
			expectMsg: "should not fail",
			wantValue: 42,
			wantErr:   false,
		},
		{
			name:      "Err result returns wrapped error with message",
			result:    Err[int](testErr),
			expectMsg: "operation failed",
			wantValue: 0,
			wantErr:   true,
			errMsg:    "operation failed: original error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := tt.result.Expect(tt.expectMsg)

			if v != tt.wantValue {
				t.Errorf("Expect() value = %v, want %v", v, tt.wantValue)
			}

			if tt.wantErr {
				if err == nil {
					t.Error("Expect() error = nil, want error")
				} else if err.Error() != tt.errMsg {
					t.Errorf("Expect() error = %q, want %q", err.Error(), tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Expect() error = %v, want nil", err)
				}
			}
		})
	}
}
