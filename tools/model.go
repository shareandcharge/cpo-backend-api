package tools

import "time"


type XLocation struct {
	ScID     string   `json:"scId"`
	Location Location `json:"data"`
}

type Connector struct {
	ID          string    `json:"id"`
	Standard    string    `json:"standard"`
	Format      string    `json:"format"`
	PowerType   string    `json:"power_type"`
	Voltage     int       `json:"voltage"`
	Amperage    int       `json:"amperage"`
	TariffID    string    `json:"tariff_id"`
	LastUpdated time.Time `json:"last_updated"`
}

type Evse struct {
	UID               string        `json:"uid"`
	EvseID            string        `json:"evse_id"`
	Status            string        `json:"status"`
	StatusSchedule    []interface{} `json:"status_schedule"`
	Capabilities      []interface{} `json:"capabilities"`
	Connectors        []Connector   `json:"connectors"`
	PhysicalReference string        `json:"physical_reference"`
	FloorLevel        string        `json:"floor_level"`
	LastUpdated       time.Time     `json:"last_updated"`
}

type Location struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Name       string `json:"name"`
	Address    string `json:"address"`
	City       string `json:"city"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
	Coordinates struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	} `json:"coordinates"`
	Evses []Evse `json:"evses"`
}



type Tariff struct {
	ID       string `json:"id"`
	Currency string `json:"currency"`
	Elements []struct {
		PriceComponents []struct {
			Type     string  `json:"type"`
			Price    float64 `json:"price"`
			StepSize int     `json:"step_size"`
		} `json:"price_components"`
	} `json:"elements"`
	LastUpdated time.Time `json:"last_updated"`
}

type TxReceiptResponse struct {
	BlockHash         string `json:"blockHash" db:"blockHash"`
	BlockNumber       int    `json:"blockNumber" db:"blockNumber"`
	ContractAddress   interface{} `json:"contractAddress" db:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed" db:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed" db:"gasUsed"`
	LogsNumber        string `json:"logs" db:"logs_number"`
	LogsBloom         string `json:"logsBloom" db:"logsBloom"`
	Root              interface{} `json:"root" db:"root"`
	Status            string `json:"status" db:"status"`
	TransactionHash   string `json:"transactionHash" db:"transactionHash"`
	TransactionIndex  string `json:"transactionIndex" db:"transactionIndex"`
	Timestamp         uint64 `json:"timestamp" db:"timestamp"`
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
