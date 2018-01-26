package main

// Hello (EHLO message)
func Hello(iIn Identity, sIn Signature, messageIn string) (iOut Identity, sOut Signature, messageOut string) {
	// ? Is there a need to establish some form of a basic session here, after verification of the signature?
	return Identity{}, Signature{}, "Hello back how can I help you?"

}
