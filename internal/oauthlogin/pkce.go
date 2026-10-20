// Package oauthlogin implements the client side of rolewatch's MCP OAuth 2.1
// sign-in (rolewatch#3114): PKCE generation, dynamic client registration, the
// local loopback callback listener, and the authorization-code token exchange.
// It lets `rolewatch auth login` obtain a bearer token by opening a browser
// instead of requiring a pasted API key — the resulting token is stored the
// same way (config.Creds.Token) and sent the same way (Authorization: Bearer),
// so every other package in this CLI, and rolewatch-mcp reading the same
// ~/.rolewatch/creds.json, needs no changes at all.
package oauthlogin

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// generateVerifier mints a high-entropy PKCE code_verifier: 32 random bytes,
// unpadded base64url — 43 characters, the minimum RFC 7636 asks for and the
// size every major OAuth client uses in practice.
func generateVerifier() (string, error) {
	var b [32]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b[:]), nil
}

// challengeFor computes the S256 code_challenge for a verifier — the only
// method rolewatch's authorization server accepts.
func challengeFor(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

// generateState mints the opaque CSRF-style state value round-tripped through
// the authorize redirect, so the local callback can refuse a response that
// does not belong to the login attempt that is waiting for it.
func generateState() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b[:]), nil
}
