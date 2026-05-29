package credentials

import "testing"

func TestEncryptDecrypt(t *testing.T) {
	key := []byte("test-key-32-bytes-1234567890abc")
	original := []byte(`{"phone":"+79991234567","password":"secret"}`)

	enc, err := Encrypt(original, key)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}

	dec, err := Decrypt(enc, key)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}

	if string(dec) != string(original) {
		t.Errorf("decrypted mismatch: got %s, want %s", dec, original)
	}
}

func TestEncrypt_DifferentCiphertext(t *testing.T) {
	key := []byte("test-key-32-bytes-1234567890abc")
	data := []byte("hello")

	enc1, _ := Encrypt(data, key)
	enc2, _ := Encrypt(data, key)

	if enc1 == enc2 {
		t.Error("expected different ciphertexts (random nonce)")
	}
}

func TestDecrypt_InvalidHex(t *testing.T) {
	_, err := Decrypt("invalid!", []byte("key-32-bytes-1234567890abcdef"))
	if err == nil {
		t.Error("expected error for invalid hex")
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	orig := []byte("secret data")
	key1 := []byte("key-32-bytes-aaaaaaaaaaaaaaaaaaaaaa")
	key2 := []byte("key-32-bytes-bbbbbbbbbbbbbbbbbbbbbb")

	enc, _ := Encrypt(orig, key1)
	_, err := Decrypt(enc, key2)
	if err == nil {
		t.Error("expected error for wrong key")
	}
}

func TestDeriveKey_Padding(t *testing.T) {
	short := []byte("short")
	key, err := deriveKey(short)
	if err != nil {
		t.Fatalf("deriveKey: %v", err)
	}
	if len(key) != 32 {
		t.Errorf("expected 32 bytes, got %d", len(key))
	}
}

func TestDeriveKey_Empty(t *testing.T) {
	_, err := deriveKey([]byte{})
	if err != ErrInvalidKey {
		t.Errorf("expected ErrInvalidKey, got %v", err)
	}
}

func TestDeriveKey_Long(t *testing.T) {
	long := make([]byte, 64)
	for i := range long {
		long[i] = byte(i)
	}
	key, err := deriveKey(long)
	if err != nil {
		t.Fatalf("deriveKey: %v", err)
	}
	if len(key) != 32 {
		t.Errorf("expected 32 bytes, got %d", len(key))
	}
}
