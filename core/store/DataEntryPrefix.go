package store

// DataEntryPrefix
type DataEntryPrefix byte

const (
	// DATA
	DATA_BlockHash		DataEntryPrefix = 0x00
	DATA_Block 		DataEntryPrefix = 0x01
	DATA_Transaction 	DataEntryPrefix = 0x02
	DATA_Contract 		DataEntryPrefix = 0x03

	// INDEX
	IX_HeaderHashList	DataEntryPrefix = 0x80
	IX_Enrollment		DataEntryPrefix = 0x84
	IX_Unspent		DataEntryPrefix = 0x90
	IX_Unclaimed		DataEntryPrefix = 0x91
	IX_Vote			DataEntryPrefix = 0x94

	// ASSET
	ST_Info			DataEntryPrefix = 0x40
	ST_QuantityIssued	DataEntryPrefix = 0x41
	ST_Contract		DataEntryPrefix = 0x50
	ST_Storage		DataEntryPrefix = 0x70

	//SYSTEM
	SYS_CurrentBlock	DataEntryPrefix = 0xc0
	SYS_CurrentHeader	DataEntryPrefix = 0xc1

	//CONFIG
	CFG_Version		DataEntryPrefix = 0xf0
)
