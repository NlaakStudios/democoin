package structures

//
type Protocol struct {
	// Doamin is the owner of this protocol. The same protocol can be used by many different domains.
	// For example, a chat app (protocol) can be used by both Facebook and twitter.
	Domain string
	// Type is the Unsigned interger value use to identofy what type of protocol this is for passing to the correct
	// code for processing.
	Type uint
	// Data is a JSON string representing the the data/info needed to for this protocol
	Data string
}
