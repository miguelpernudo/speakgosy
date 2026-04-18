package packet

import (
	"encoding/json"
	"github.com/miguelpernudo/speakgosy/internal/crypt"
)

// Build serializes and encrypts a Payload into a SPA packet.
func Build(p Payload, key []byte) ([]byte, error) {

	// Struct to []byte (JSON)
	data, err := json.Marshal(p)
	if err != nil {	
		return nil, err
	}

	ciphertext, err := crypt.Encrypt(key, data) 
	    if err != nil {
	        return nil, err
	    }

	return ciphertext, nil
}
