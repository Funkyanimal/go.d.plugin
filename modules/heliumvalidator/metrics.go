package heliumvalidator

// API docs: https://github.com/energicryptocurrency/core-api-documentation

type heliumvalidatorInfo struct {
	BlockHeight *blockheightInfo `stm:"block_height"`
	BlockAge    *blockageInfo    `stm:"block_age"`
	In_Consensus    *inconsensusInfo    `stm:"in_consensus"`
}

type blockheightInfo struct {
	Height     float64 `stm:"block_height" json:"height"`
}

type blockageInfo struct {
	Age      float64 `stm:"block_age" json:"block_age"`
}

