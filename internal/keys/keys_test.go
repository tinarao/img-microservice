package keys

import (
	"fmt"
	"testing"
)

func TestKeysWorkflow(t *testing.T) {
	email := "tinarao228@gmail.com"

	kp, err := NewKeysPair(email)
	if err != nil {
		t.Fatalf("failed to generate keys: %v", err)
	}

	fmt.Printf("public key: %s\n\n", kp.PublicKey)
	fmt.Printf("private key: %s\n\n", kp.PrivateKey)

	content, err := kp.Unmarshal()
	fmt.Printf("content: %s\n\n", *content)
	if *content != email {
		t.Fatal("email does not match")
	}

}
