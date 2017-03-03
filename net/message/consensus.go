package message

import (
	"GoOnchain/common"
	"GoOnchain/common/serialization"
	"GoOnchain/core/contract/program"
	//"GoOnchain/events"
	"GoOnchain/events"
	. "GoOnchain/net/protocol"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"unsafe"
)

type ConsensusPayload struct {
	Version    uint32
	PrevHash   common.Uint256
	Height     uint32
	MinerIndex uint16
	Timestamp  uint32
	Data       []byte
	Program    *program.Program

	hash common.Uint256
}

type consensus struct {
	msgHdr
	cons ConsensusPayload
	//event *events.Event
	//TBD
}

func (cp *ConsensusPayload) Hash() common.Uint256 {
	return common.Uint256{}
}

func (cp *ConsensusPayload) Verify() error {
	return nil
}

func (cp *ConsensusPayload) InvertoryType() common.InventoryType {
	return common.CONSENSUS
}

func (cp *ConsensusPayload) GetProgramHashes() ([]common.Uint160, error) {
	return nil, nil
}

func (cp *ConsensusPayload) SetPrograms([]*program.Program) {
}

func (cp *ConsensusPayload) GetPrograms() []*program.Program {
	return nil
}

func (cp *ConsensusPayload) GetMessage() []byte {
	//TODO: GetMessage
	return []byte{}
}

func (msg consensus) Handle(node Noder) error {
	common.Trace()
	fmt.Printf("RX Consensus message\n")

	node.LocalNode().GetEvent("consensus").Notify(events.EventNewInventory, &msg.cons)
	return nil
}

func reqConsensusData(node Noder, hash common.Uint256) error {
	var msg dataReq
	msg.dataType = common.CONSENSUS
	// TODO handle the hash array case
	msg.hash = hash

	buf, _ := msg.Serialization()
	go node.Tx(buf)

	return nil
}
func (cp *ConsensusPayload) Type() common.InventoryType {

	//TODO:Temporary add for Interface signature.SignableData use.
	return common.CONSENSUS
}

func (cp *ConsensusPayload) SerializeUnsigned(w io.Writer) error {
	serialization.WriteUint32(w, cp.Version)
	cp.PrevHash.Serialize(w)
	serialization.WriteUint32(w, cp.Height)
	serialization.WriteUint16(w, cp.MinerIndex)
	serialization.WriteUint32(w, cp.Timestamp)
	return nil

}

func (cp *ConsensusPayload) Serialize(w io.Writer) {
	cp.SerializeUnsigned(w)
	serialization.WriteVarBytes(w, cp.Data)
	cp.Program.Serialize(w)
}

func (msg *consensus) Serialization() ([]byte, error) {
	hdrBuf, err := msg.msgHdr.Serialization()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(hdrBuf)
	msg.cons.Serialize(buf)

	return buf.Bytes(), err
}

func (cp *ConsensusPayload) DeserializeUnsigned(r io.Reader) error {
	var err error
	cp.Version, err = serialization.ReadUint32(r)
	if err != nil {
		return errors.New("consensus item Version Deserialize failed.")
	}

	preBlock := new(common.Uint256)
	err = preBlock.Deserialize(r)
	if err != nil {
		return errors.New("consensus item preHash Deserialize failed.")
	}
	cp.PrevHash = *preBlock

	cp.Height, err = serialization.ReadUint32(r)
	if err != nil {
		return errors.New("consensus item Height Deserialize failed.")
	}

	cp.MinerIndex, err = serialization.ReadUint16(r)
	if err != nil {
		return errors.New("consensus item MinerIndex Deserialize failed.")
	}

	cp.Timestamp, err = serialization.ReadUint32(r)
	if err != nil {
		return errors.New("consensus item Timestamp Deserialize failed.")
	}

	cp.Data, err = serialization.ReadVarBytes(r)
	if err != nil {
		return errors.New("consensus item Data Deserialize failed.")
	}
	return nil
}

func (cp *ConsensusPayload) Deserialize(r io.Reader) {
	cp.DeserializeUnsigned(r)

	cp.Program.Deserialize(r)
}

func (msg *consensus) Deserialization(p []byte) error {
	fmt.Printf("The size of messge is %d in deserialization\n",
		uint32(unsafe.Sizeof(*msg)))
	buf := bytes.NewBuffer(p)
	err := binary.Read(buf, binary.LittleEndian, &(msg.msgHdr))
	msg.cons.Deserialize(buf)
	return err
}

func NewConsensus(cp *ConsensusPayload) ([]byte, error) {
	common.Trace()
	var msg consensus
	msg.msgHdr.Magic = NETMAGIC
	cmd := "consensus"
	copy(msg.msgHdr.CMD[0:len(cmd)], cmd)
	tmpBuffer := bytes.NewBuffer([]byte{})
	cp.Serialize(tmpBuffer)
	msg.cons = *cp
	b := new(bytes.Buffer)
	err := binary.Write(b, binary.LittleEndian, tmpBuffer.Bytes())
	if err != nil {
		fmt.Println("Binary Write failed at new Msg")
		return nil, err
	}
	s := sha256.Sum256(b.Bytes())
	s2 := s[:]
	s = sha256.Sum256(s2)
	buf := bytes.NewBuffer(s[:4])
	binary.Read(buf, binary.LittleEndian, &(msg.msgHdr.Checksum))
	msg.msgHdr.Length = uint32(len(b.Bytes()))
	fmt.Printf("The message payload length is %d\n", msg.msgHdr.Length)

	m, err := msg.Serialization()
	if err != nil {
		fmt.Println("Error Convert net message ", err.Error())
		return nil, err
	}

	str := hex.EncodeToString(m)
	fmt.Printf("The message length is %d, %s\n", len(m), str)

	return m, nil
}
