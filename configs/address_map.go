package configs

var AddressMap = make(map[string]string)
var TopicsMap = make(map[string]string)

//------- ADD TO THIS FUNCTION TO GET ADDRESS TO NAME  -------------
func Init() {

	AddressMap["0x6b0c56d1ad5144b4d37fa6e27dc9afd5c2435c3b"] = "Faucet"
	AddressMap["0x7b0f2b531c018d4269a95561cfb4e038a7e3c8dc"] = "CPO 1"
	AddressMap["0xf60b71a4d360a42ec9d4e7977d8d9928fd7c8365"] = "MSP"

	TopicsMap["msptokenmint"] = "0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885"
	TopicsMap["cdr"] = "0xaee5526e82ba0d7be2d0181caa5c4ac1ddbcef1917331b78a4a6cfdd93f3126a"
	TopicsMap["error"] = "0x57cf7a55e859b30b6bfeb9a7dd14411606106cb3e082f2cda387ec3b4b90be1c"
	TopicsMap["transferToken"] = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
}

func AddressToName(addr string) string {
	Init()
	if val, ok := AddressMap[addr]; ok {
		return val
	}
	return addr
}
