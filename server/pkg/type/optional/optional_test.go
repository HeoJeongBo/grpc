package optional

import (
	"testing"
)

func TestSome(t *testing.T) {
	opt := Some(42)

	if !opt.IsSome() {
		t.Error("Some() should return IsSome() = true")
	}

	if opt.IsNone() {
		t.Error("Some() should return IsNone() = false")
	}

	val, ok := opt.Unwrap()

	if !ok {
		t.Error("Some().Unwrap() should return ok = true")
	}
	if val != 42 {
		t.Errorf("Some(42).Unwrap() = %v, want 42", val)
	}
}

func TestNone(t *testing.T) {
	opt := None[int]()

	if opt.IsSome() {
		t.Error("None() should return IsSome() = false")
	}

	if !opt.IsNone() {
		t.Error("None() should return IsNone() = true")
	}

	_, ok := opt.Unwrap()

	if ok {
		t.Error("None().Unwrap() should return ok = false")
	}
}

func TestUnwrapOr(t *testing.T) {
	tests := []struct {
		name         string
		opt          Optional[int]
		defaultValue int
		want         int
	}{
		{
			name:         "Some returns value",
			opt:          Some(42),
			defaultValue: 100,
			want:         42,
		},
		{
			name:         "None returns default",
			opt:          None[int](),
			defaultValue: 100,
			want:         100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.opt.UnwrapOr((tt.defaultValue))
			if v != tt.want {
				t.Errorf("UnwrapOr() = %v, want %v", v, tt.want)
			}
		})
	}
}

func TestUnwrapOrElse(t *testing.T) {
	tests := []struct {
		name string
		opt  Optional[int]
		fn   func() int
		want int
	}{
		{
			name: "Some returns value",
			opt:  Some(42),
			fn:   func() int { return 100 },
			want: 42,
		},
		{
			name: "None calls function",
			opt:  None[int](),
			fn:   func() int { return 100 },
			want: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.opt.UnwrapOrElse(tt.fn)
			if v != tt.want {
				t.Errorf("UnwrapOrElse() = %v, want %v", v, tt.want)
			}
		})
	}
}

func TestMap(t *testing.T) {
	double := func(x int) int { return x * 2 }

	tests := []struct {
		name   string
		opt    Optional[int]
		fn     func(int) int
		wantOk bool
		want   int
	}{
		{
			name:   "Some maps value",
			opt:    Some(21),
			fn:     double,
			wantOk: true,
			want:   42,
		},
		{
			name:   "None returns none",
			opt:    None[int](),
			fn:     double,
			wantOk: false,
			want:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.opt.Map(tt.fn)
			if r.IsSome() != tt.wantOk {
				t.Errorf("Map().IsSome() = %v, want %v", r.IsSome(), tt.wantOk)
			}

			if tt.wantOk {
				v, _ := r.Unwrap()
				if v != tt.want {
					t.Errorf("Map() value = %v, want %v", v, tt.want)
				}
			}
		})
	}
}

func TestFilter(t *testing.T) {
	isEven := func(x int) bool { return x%2 == 0 }

	tests := []struct {
		name      string
		opt       Optional[int]
		predicate func(int) bool
		wantOk    bool
	}{
		{
			name:      "Some with matching predicate",
			opt:       Some(42),
			predicate: isEven,
			wantOk:    true,
		},
		{
			name:      "Some with non-matching predicate",
			opt:       Some(43),
			predicate: isEven,
			wantOk:    false,
		},
		{
			name:      "None returns None",
			opt:       None[int](),
			predicate: isEven,
			wantOk:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.opt.Filter(tt.predicate)
			if v.IsSome() != tt.wantOk {
				t.Errorf("Filter().IsSome() = %v, want %v", v.IsSome(), tt.wantOk)
			}
		})
	}
}

func TestOrElse(t *testing.T) {
	tests := []struct {
		name        string
		opt         Optional[int]
		alternative Optional[int]
		want        int
		wantOk      bool
	}{
		{
			name:        "Some returns original",
			opt:         Some(42),
			alternative: Some(100),
			want:        42,
			wantOk:      true,
		},
		{
			name:        "None returns alternative Some",
			opt:         None[int](),
			alternative: Some(100),
			want:        100,
			wantOk:      true,
		},
		{
			name:        "None with None alternative",
			opt:         None[int](),
			alternative: None[int](),
			want:        0,
			wantOk:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := tt.opt.OrElse(tt.alternative)
			if v.IsSome() != tt.wantOk {
				t.Errorf("OrElse().IsSome() = %v, want %v", v.IsSome(), tt.wantOk)
			}
			if tt.wantOk {
				r, _ := v.Unwrap()
				if r != tt.want {
					t.Errorf("OrElse() value = %v, want %v", r, tt.want)
				}
			}
		})
	}
}
