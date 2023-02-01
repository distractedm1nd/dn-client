package client

import (
	"github.com/gogo/protobuf/types"
	"github.com/ipfs/go-cid"
)

type ExtendedHeader struct {
	// This type is missing the rest of the fields, because they are tendermint types which we cannot import here
	DAH *DataAvailabilityHeader `json:"dah"`
}

// NamespacedShares represents all shares with proofs within a specific namespace of an EDS.
type NamespacedShares []NamespacedRow

// NamespacedRow represents all shares with proofs within a specific namespace of a single EDS row.
type NamespacedRow struct {
	Shares [][]byte
	Proof  *Proof
}

type Proof struct {
	Nodes      []cid.Cid
	Start, End int
}

// DataAvailabilityHeader (DAHeader) is copied from celestia-node.
// It contains the row and column roots of the
// erasure coded version of the data in Block.Data. The original Block.Data is
// split into shares and arranged in a square of width squareSize. Then, this
// square is "extended" into an extended data square (EDS) of width 2*squareSize
// by applying Reed-Solomon encoding. For details see Section 5.2 of
// https://arxiv.org/abs/1809.09044 or the Celestia specification:
// https://github.com/celestiaorg/celestia-specs/blob/master/src/specs/data_structures.md#availabledataheader
type DataAvailabilityHeader struct {
	// RowRoot_j = root((M_{j,1} || M_{j,2} || ... || M_{j,2k} ))
	RowsRoots [][]byte `json:"row_roots"`
	// ColumnRoot_j = root((M_{1,j} || M_{2,j} || ... || M_{2k,j} ))
	ColumnRoots [][]byte `json:"column_roots"`
	// hash is the Merkle root of the row and column roots. This field is the
	// memoized result from `Hash()`.
	hash []byte
}

// Types below are copied from celestia-node (or cosmos-sdk dependency of celestia node, to be precise)
// They are needed for proper deserialization.
// It's probably far from the best approach to those types, but it's simple and works.
// Some alternatives:
// 1. Generate types from protobuf definitions (and automate updating of protobuf files)
// 2. Extract common dependency that defines all types used in RPC.

// TxResponse defines a structure containing relevant tx data and metadata. The
// tags are stringified and the log is JSON decoded.
type TxResponse struct {
	// The block height
	Height int64 `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	// The transaction hash.
	TxHash string `protobuf:"bytes,2,opt,name=txhash,proto3" json:"txhash,omitempty"`
	// Namespace for the Code
	Codespace string `protobuf:"bytes,3,opt,name=codespace,proto3" json:"codespace,omitempty"`
	// Response code.
	Code uint32 `protobuf:"varint,4,opt,name=code,proto3" json:"code,omitempty"`
	// Result bytes, if any.
	Data string `protobuf:"bytes,5,opt,name=data,proto3" json:"data,omitempty"`
	// The output of the application's logger (raw string). May be
	// non-deterministic.
	RawLog string `protobuf:"bytes,6,opt,name=raw_log,json=rawLog,proto3" json:"raw_log,omitempty"`
	// The output of the application's logger (typed). May be non-deterministic.
	Logs ABCIMessageLogs `protobuf:"bytes,7,rep,name=logs,proto3,castrepeated=ABCIMessageLogs" json:"logs"`
	// Additional information. May be non-deterministic.
	Info string `protobuf:"bytes,8,opt,name=info,proto3" json:"info,omitempty"`
	// Amount of gas requested for transaction.
	GasWanted int64 `protobuf:"varint,9,opt,name=gas_wanted,json=gasWanted,proto3" json:"gas_wanted,omitempty"`
	// Amount of gas consumed by transaction.
	GasUsed int64 `protobuf:"varint,10,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	// The request transaction bytes.
	Tx *types.Any `protobuf:"bytes,11,opt,name=tx,proto3" json:"tx,omitempty"`
	// Time of the previous block. For heights > 1, it's the weighted median of
	// the timestamps of the valid votes in the block.LastCommit. For height == 1,
	// it's genesis time.
	Timestamp string `protobuf:"bytes,12,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

// ABCIMessageLogs represents a slice of ABCIMessageLog.
type ABCIMessageLogs []ABCIMessageLog

// ABCIMessageLog defines a structure containing an indexed tx ABCI message log.
type ABCIMessageLog struct {
	MsgIndex uint32 `protobuf:"varint,1,opt,name=msg_index,json=msgIndex,proto3" json:"msg_index,omitempty"`
	Log      string `protobuf:"bytes,2,opt,name=log,proto3" json:"log,omitempty"`
	// Events contains a slice of Event objects that were emitted during some
	// execution.
	Events StringEvents `protobuf:"bytes,3,rep,name=events,proto3,castrepeated=StringEvents" json:"events"`
}

// StringAttributes defines a slice of StringEvents objects.
type StringEvents []StringEvent

// StringEvent defines en Event object wrapper where all the attributes
// contain key/value pairs that are strings instead of raw bytes.
type StringEvent struct {
	Type       string      `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Attributes []Attribute `protobuf:"bytes,2,rep,name=attributes,proto3" json:"attributes"`
}

// Attribute defines an attribute wrapper where the key and value are
// strings instead of raw bytes.
type Attribute struct {
	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}
