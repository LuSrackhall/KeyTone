package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

// Hardcoded mask for XOR operation - must match the one in sdk/signature/authorization.go
// This is a simple obfuscation, not strong encryption.
var realMask = []byte{0x55, 0xAA, 0x33, 0xCC, 0x99, 0x66, 0x11, 0xEE, 0x77, 0xBB, 0x22, 0xDD, 0x88, 0x44, 0xFF, 0x00}

func main() {
	keyPtr := flag.String("key", "", "The plaintext key to obfuscate (32 bytes)")
	flag.Parse()

	if *keyPtr == "" {
		// NOTE: stdout is reserved for the machine-consumable hex output.
		// Any human-facing messages MUST go to stderr to avoid polluting build flags.
		fmt.Fprintln(os.Stderr, "Please provide a key using -key flag")
		os.Exit(1)
	}

	key := []byte(*keyPtr)
	if len(key) != 32 {
		// NOTE: keep warning on stderr so stdout stays pure hex.
		fmt.Fprintf(os.Stderr, "Warning: Key length is %d, expected 32 bytes.\n", len(key))
	}

	obfuscated := make([]byte, len(key))
	for i, b := range key {
		obfuscated[i] = b ^ realMask[i%len(realMask)]
	}

	// Output as hex string for easy injection via ldflags or copy-paste
	fmt.Printf("%s", hex.EncodeToString(obfuscated))
}
