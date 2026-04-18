package packet

import (
        "encoding/json"
        "github.com/miguelpernudo/speakgosy/internal/crypt"
)

// Parse decrypts and deserializes a SPA packet into a Payload.
func Parse(data []byte, key []byte) (Payload, error) {

	plaintext, err := crypt.Decrypt(key, data)
	    if err != nil {
	        return Payload{}, err
	    }

	// []byte → struct
	var p Payload	
	err = json.Unmarshal(plaintext, &p)  
  if err != nil {     
  	return Payload{}, err  
  	}

	return p, nil
}
