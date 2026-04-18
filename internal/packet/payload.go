package packet

// Payload holds the data encrypted inside a SPA packet.
type Payload struct {
    ClientIP string
    Nonce    string
    TOTP     string
}
