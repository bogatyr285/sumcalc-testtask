package sha256

import (
	"fmt"
	"testing"
)

func TestNewSHA256(t *testing.T) {
	for i, tt := range []struct {
		in  []byte
		out string
	}{
		{[]byte(""), "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{[]byte("abc"), "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"},
		{[]byte("hello"), "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{[]byte{10}, "01ba4719c80b6fe911b091a7c05124b64eeece964e09c058ef8f9805daca546b"},
		{[]byte("10"), "4a44dc15364204a80fe80e9039455cc1608281820fe2b24f1e5233ade6af1dd5"},
	} {
		t.Run(fmt.Sprintf("%v", i), func(t *testing.T) {
			result := SHA256Byte(tt.in)
			if result != tt.out {
				t.Errorf("want %v; got %v", tt.out, result)
			}
		})
	}
}
