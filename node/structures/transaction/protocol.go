package transaction

//
type Protocol struct {
	// Doamin is the owner of this protocol. The same protocol can be used by many different domains.
	// For example, a chat app (protocol) can be used by both Facebook and twitter.
	// Domains are stored & synced in a bolt db on all devices. they are initially synced from a Master node.
	// Example domain: com.somedomain.chat
	Domain string
	// Location of origin. Looked up in location database (bolt) and used to prefer closer nodes for better local area perforance
	// Locations are stored & synced in a bolt db on all devices. they are initially synced from a Master node.
	Location uint
	// Type is the Unsigned interger value use to identofy what type of protocol this is for passing to the correct
	// code for processing. Protocol Type are compiled into the binary and are defined in protocoltype.go.
	// ie. Protocoltype.Chat or ProtocolType.Currency
	Type uint
	// Data is a JSON string representing the the data/info needed to for this protocol
	// JSOn is different for each protocol, but have some common properties to all.
	Data string
}

func NewProtocol() Protocol {
	return Protocol{}
}

func (p *Protocol) syncLocations() {}

func (p *Protocol) syncDomains() {}

func (p *Protocol) detectLocation() {}
