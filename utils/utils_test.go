package utils

import "testing"

func TestQueryProcessName(t *testing.T) {
	type test struct {
		name string
		pid uint32
		want string
	}
	tests := []test {
		{"pid 0", 0, "System"},
		{"pid 1", 1, "System"},
		{"pid 4", 4, "System"},
		{"pid 124", 124, "unknown"},
		{"pid 19024", 19024, "msedge.exe"},
	}
	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			result := QueryProcessName(c.pid)
			if result != c.want {
				t.Errorf("excepted: %v result:%v", c.want, result)
			}
		})
	}
}

func BenchmarkQueryProcessName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		QueryProcessName(19024)
	}
}
