package tools

import "time"

// type XLocation struct {
// 	ScID     Location   `json:",string"`
// }

type XLocation map[string]Location

type Location struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	City        string `json:"city"`
	PostalCode  string `json:"postal_code"`
	Country     string `json:"country"`
	Coordinates struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	} `json:"coordinates"`
	Evses      []Evse `json:"evses"`
	Directions []struct {
		Language string `json:"language"`
		Text     string `json:"text"`
	} `json:"directions"`
	Operator struct {
		Name string `json:"name"`
	} `json:"operator"`
	LastUpdated time.Time `json:"last_updated"`
}

type Evse struct {
	UID        string `json:"uid"`
	EvseID     string `json:"evse_id"`
	Status     string `json:"status"`
	Connectors []struct {
		ID          string    `json:"id"`
		Standard    string    `json:"standard"`
		PowerType   string    `json:"power_type"`
		Voltage     int       `json:"voltage"`
		Amperage    int       `json:"amperage"`
		TariffID    string    `json:"tariff_id"`
		LastUpdated time.Time `json:"last_updated"`
	} `json:"connectors"`
	FloorLevel        string    `json:"floor_level"`
	PhysicalReference string    `json:"physical_reference"`
	LastUpdated       time.Time `json:"last_updated"`
}

// type Tariff struct {
// 	ID       string `json:"id"`
// 	Currency string `json:"currency"`
// 	Elements []struct {
// 		PriceComponents []struct {
// 			Type     string  `json:"type"`
// 			Price    float64 `json:"price"`
// 			StepSize int     `json:"step_size"`
// 		} `json:"price_components"`
// 	} `json:"elements"`
// 	LastUpdated time.Time `json:"last_updated"`
// }

type XTariff map[string]Tariff

type Tariff struct {
	ID            string `json:"id"`
	Currency      string `json:"currency"`
	TariffAltText []struct {
		Language string `json:"language"`
		Text     string `json:"text"`
	} `json:"tariff_alt_text,omitempty"`
	TariffAltUrl string `json:"tariff_alt_url,omitempty"`
	Elements     []struct {
		PriceComponents []struct {
			Type     string  `json:"type"`
			Price    float64 `json:"price"`
			StepSize float64     `json:"step_size,omitempty"`
		} `json:"price_components"`
		Restrictions struct {
			StartTime 		string 		`json:"start_time,omitempty"`
			EndTime 		string 		`json:"end_time,omitempty"`
			StartDate 		string 		`json:"start_date,omitempty"`
			EndDate 		string 		`json:"end_date,omitempty"`
			MinKwh			string 		`json:"min_kwh,omitempty"`
			MaxKwh 			string 		`json:"max_kwh,omitempty"`
			MinPower 		string 		`json:"min_power,omitempty"`
			MaxPower 		string 		`json:"max_power,omitempty"`
			MinDuration 	string 		`json:"min_duration,omitempty"`
			MaxDuration 	string 		`json:"max_duration,omitempty"`
			DayOfWeek 		[]string 	`json:"day_of_week,omitempty"`
		} `json:"restrictions,omitempty"`
	} `json:"elements"`
	LastUpdated string `json:"last_updated"`
}

type TxReceiptResponse struct {
	BlockHash         string      `json:"blockHash" db:"blockHash"`
	BlockNumber       int         `json:"blockNumber" db:"blockNumber"`
	ContractAddress   interface{} `json:"contractAddress" db:"contractAddress"`
	CumulativeGasUsed string      `json:"cumulativeGasUsed" db:"cumulativeGasUsed"`
	GasUsed           string      `json:"gasUsed" db:"gasUsed"`
	LogsNumber        string      `json:"logs" db:"logs_number"`
	LogsBloom         string      `json:"logsBloom" db:"logsBloom"`
	Root              interface{} `json:"root" db:"root"`
	Status            string      `json:"status" db:"status"`
	TransactionHash   string      `json:"transactionHash" db:"transactionHash"`
	TransactionIndex  string      `json:"transactionIndex" db:"transactionIndex"`
	Timestamp         uint64      `json:"timestamp" db:"timestamp"`
}

type TxLog struct {
	Address             string   `json:"address"`
	BlockHash           string   `json:"blockHash"`
	BlockNumber         string   `json:"blockNumber"`
	Data                string   `json:"data"`
	LogIndex            string   `json:"logIndex"`
	Removed             bool     `json:"removed"`
	Topics              []string `json:"topics"`
	TransactionHash     string   `json:"transactionHash"`
	TransactionIndex    string   `json:"transactionIndex"`
	TransactionLogIndex string   `json:"transactionLogIndex"`
	Type                string   `json:"type"`
}

//when query the blockchain, the response
type BlockResponse struct {
	Author           string          `json:"author"`
	Difficulty       string          `json:"difficulty"`
	ExtraData        string          `json:"extraData"`
	GasLimit         string          `json:"gasLimit"`
	GasUsed          string          `json:"gasUsed"`
	Hash             string          `json:"hash"`
	LogsBloom        string          `json:"logsBloom"`
	Miner            string          `json:"miner"`
	Number           string          `json:"number"`
	ParentHash       string          `json:"parentHash"`
	ReceiptsRoot     string          `json:"receiptsRoot"`
	SealFields       []string        `json:"sealFields"`
	Sha3Uncles       string          `json:"sha3Uncles"`
	Signature        string          `json:"signature"`
	Size             string          `json:"size"`
	StateRoot        string          `json:"stateRoot"`
	Step             string          `json:"step"`
	Timestamp        string          `json:"timestamp"`
	TotalDifficulty  string          `json:"totalDifficulty"`
	Transactions     []TxTransaction `json:"transactions"`
	TransactionsRoot string          `json:"transactionsRoot"`
	Uncles           []interface{}   `json:"uncles"`
}

type TxTransaction struct {
	BlockHash        string      `json:"blockHash" db:"blockHash"`
	BlockNumber      int         `json:"blockNumber" db:"blockNumber"`
	ChainID          string      `json:"chainId" db:"chainId"`
	Condition        interface{} `json:"condition" db:"x_condition"`
	Creates          interface{} `json:"creates" db:"creates"`
	From             string      `json:"from" db:"from_addr"`
	Gas              string      `json:"gas" db:"gas"`
	GasPrice         string      `json:"gasPrice" db:"gasPrice"`
	Hash             string      `json:"hash" db:"hash"`
	Input            string      `json:"input" db:"x_input"`
	Nonce            string      `json:"nonce" db:"nonce"`
	PublicKey        string      `json:"publicKey" db:"publicKey"`
	R                string      `json:"r" db:"r"`
	Raw              string      `json:"raw" db:"raw"`
	S                string      `json:"s" db:"s"`
	StandardV        string      `json:"standardV" db:"standardV"`
	To               string      `json:"to" db:"to_addr"`
	TransactionIndex string      `json:"transactionIndex" db:"transactionIndex"`
	V                string      `json:"v" db:"v"`
	Value            string      `json:"value" db:"x_value"`
	Timestamp        uint64      `json:"timestamp" db:"timestamp"`
}

type Reimbursement struct {
	Index           int    `json:"index"`
	Id              int    `json:"id" db:"id"`
	MspName         string `json:"msp_name" db:"msp_name"`
	CpoName         string `json:"cpo_name" db:"cpo_name"`
	Amount          int    `json:"amount" db:"amount"`
	Currency        string `json:"currency" db:"currency"`
	Timestamp       int    `json:"timestamp" db:"timestamp"`
	Status          string `json:"status" db:"status"`
	ReimbursementId string `json:"reimbursement_id" db:"reimbursement_id"`
	CdrRecords      string `json:"cdr_records" db:"cdr_records"`
	ServerAddr      string `json:"server_addr" db:"server_addr"`
	TxNumber        string `json:"txs_number" db:"txs_number"`
	TokenAddress    string `json:"token_address" db:"token_address"`
}

type CDR struct {
	EvseID           string `json:"evseId"`
	ScID             string `json:"scId"`
	LocationName     string `json:"location_name"`
	LocationAddress  string `json:"location_address"`
	Controller       string `json:"controller"`
	Start            string `json:"start"`
	End              string `json:"end"`
	FinalPrice       string `json:"finalPrice"`
	TokenContract    string `json:"tokenContract"`
	Tariff           string `json:"tariff"`
	ChargedUnits     string `json:"chargedUnits"`
	ChargingContract string `json:"chargingContract"`
	TransactionHash  string `json:"transactionHash"`
	Currency         string `json:"currency"`
}
