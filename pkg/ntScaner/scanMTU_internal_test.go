package ntScaner

import "testing"

func TestFindLargestPayload(t *testing.T) {
	const threshold = 1472
	minimumSeen := threshold
	largest, err := findLargestPayload(8972, func(size int) (bool, error) {
		if size < minimumSeen {
			minimumSeen = size
		}
		return size <= threshold, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if largest != threshold {
		t.Fatalf("largest payload=%d, want %d", largest, threshold)
	}
	if minimumSeen < 0 {
		t.Fatalf("probe received negative payload %d", minimumSeen)
	}
}

func TestFindLargestPayloadUnreachable(t *testing.T) {
	probes := 0
	_, err := findLargestPayload(1472, func(size int) (bool, error) {
		probes++
		return false, nil
	})
	if err == nil {
		t.Fatal("expected unreachable destination error")
	}
	if probes != 1 {
		t.Fatalf("unreachable destination used %d probes, want 1", probes)
	}
}
