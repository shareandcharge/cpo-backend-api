package configs

var AddressMap = make(map[string]string)


//------- ADD TO THIS FUNCTION TO GET ADDRESS TO NAME  -------------
func Init() {

	AddressMap["0x6b0c56d1ad5144b4d37fa6e27dc9afd5c2435c3b"] = "Faucet"
	AddressMap["0x7b0f2b531c018d4269a95561cfb4e038a7e3c8dc"] = "CPO 1"
	AddressMap["0xf60b71a4d360a42ec9d4e7977d8d9928fd7c8365"] = "MSP"

}

func AddressToName(addr string) string {
	Init()
	if val, ok := AddressMap[addr]; ok {
		return val
	}
	return addr
}
