package hashgenerator

import (
	"fmt"
)

func ExampleGenerate() {
	// Example data
	body := []byte("example body content")
	key := "secret_key"

	// Generate hash
	hash := Generate(body, key)

	// Output result
	fmt.Println(hash)

	// Output:
	// ce13299ad108f3b0538a68704ab755865975b6598b7de6467b32bd0d7de05fab
}
