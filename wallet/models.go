package wallet

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"time"

	"github.com/tonkeeper/tongo"
	"github.com/tonkeeper/tongo/boc"
	"github.com/tonkeeper/tongo/tlb"
)

type Version int

const (
	V1R1 Version = iota
	V1R2
	V1R3
	V2R1
	V2R2
	V3R1
	V3R2
	V4R1
	V4R2
	HighLoadV1R1
	HighLoadV1R2
	HighLoadV2
	HighLoadV2R1
	HighLoadV2R2
	// TODO: maybe add lockup wallet
)

const (
	DefaultSubWalletIdV3V4 = 698983191
	DefaultMessageLifetime = time.Minute * 3
	DefaultMessageMode     = 3
)

var codes = map[Version]string{
	V1R1:         "te6cckEBAQEARAAAhP8AIN2k8mCBAgDXGCDXCx/tRNDTH9P/0VESuvKhIvkBVBBE+RDyovgAAdMfMSDXSpbTB9QC+wDe0aTIyx/L/8ntVEH98Ik=",
	V1R2:         "te6cckEBAQEAUwAAov8AIN0gggFMl7qXMO1E0NcLH+Ck8mCBAgDXGCDXCx/tRNDTH9P/0VESuvKhIvkBVBBE+RDyovgAAdMfMSDXSpbTB9QC+wDe0aTIyx/L/8ntVNDieG8=",
	V1R3:         "te6cckEBAQEAXwAAuv8AIN0gggFMl7ohggEznLqxnHGw7UTQ0x/XC//jBOCk8mCBAgDXGCDXCx/tRNDTH9P/0VESuvKhIvkBVBBE+RDyovgAAdMfMSDXSpbTB9QC+wDe0aTIyx/L/8ntVLW4bkI=",
	V2R1:         "te6cckEBAQEAVwAAqv8AIN0gggFMl7qXMO1E0NcLH+Ck8mCDCNcYINMf0x8B+CO78mPtRNDTH9P/0VExuvKhA/kBVBBC+RDyovgAApMg10qW0wfUAvsA6NGkyMsfy//J7VShNwu2",
	V2R2:         "te6cckEBAQEAYwAAwv8AIN0gggFMl7ohggEznLqxnHGw7UTQ0x/XC//jBOCk8mCDCNcYINMf0x8B+CO78mPtRNDTH9P/0VExuvKhA/kBVBBC+RDyovgAApMg10qW0wfUAvsA6NGkyMsfy//J7VQETNeh",
	V3R1:         "te6cckEBAQEAYgAAwP8AIN0gggFMl7qXMO1E0NcLH+Ck8mCDCNcYINMf0x/TH/gjE7vyY+1E0NMf0x/T/9FRMrryoVFEuvKiBPkBVBBV+RDyo/gAkyDXSpbTB9QC+wDo0QGkyMsfyx/L/8ntVD++buA=",
	V3R2:         "te6cckEBAQEAcQAA3v8AIN0gggFMl7ohggEznLqxn3Gw7UTQ0x/THzHXC//jBOCk8mCDCNcYINMf0x/TH/gjE7vyY+1E0NMf0x/T/9FRMrryoVFEuvKiBPkBVBBV+RDyo/gAkyDXSpbTB9QC+wDo0QGkyMsfyx/L/8ntVBC9ba0=",
	V4R1:         "te6cckECFQEAAvUAART/APSkE/S88sgLAQIBIAIDAgFIBAUE+PKDCNcYINMf0x/THwL4I7vyY+1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8REhMUA+7QAdDTAwFxsJFb4CHXScEgkVvgAdMfIYIQcGx1Z70ighBibG5jvbAighBkc3RyvbCSXwPgAvpAMCD6RAHIygfL/8nQ7UTQgQFA1yH0BDBcgQEI9ApvoTGzkl8F4ATTP8glghBwbHVnupEx4w0kghBibG5juuMABAYHCAIBIAkKAFAB+gD0BDCCEHBsdWeDHrFwgBhQBcsFJ88WUAP6AvQAEstpyx9SEMs/AFL4J28ighBibG5jgx6xcIAYUAXLBSfPFiT6AhTLahPLH1Iwyz8B+gL0AACSghBkc3Ryuo41BIEBCPRZMO1E0IEBQNcgyAHPFvQAye1UghBkc3Rygx6xcIAYUATLBVjPFiL6AhLLassfyz+UEDRfBOLJgED7AAIBIAsMAFm9JCtvaiaECAoGuQ+gIYRw1AgIR6STfSmRDOaQPp/5g3gSgBt4EBSJhxWfMYQCAVgNDgARuMl+1E0NcLH4AD2ynftRNCBAUDXIfQEMALIygfL/8nQAYEBCPQKb6ExgAgEgDxAAGa3OdqJoQCBrkOuF/8AAGa8d9qJoQBBrkOuFj8AAbtIH+gDU1CL5AAXIygcVy//J0Hd0gBjIywXLAiLPFlAF+gIUy2sSzMzJcfsAyEAUgQEI9FHypwIAbIEBCNcYyFQgJYEBCPRR8qeCEG5vdGVwdIAYyMsFywJQBM8WghAF9eEA+gITy2oSyx/JcfsAAgBygQEI1xgwUgKBAQj0WfKn+CWCEGRzdHJwdIAYyMsFywJQBc8WghAF9eEA+gIUy2oTyx8Syz/Jc/sAAAr0AMntVEap808=",
	V4R2:         "te6cckECFAEAAtQAART/APSkE/S88sgLAQIBIAIDAgFIBAUE+PKDCNcYINMf0x/THwL4I7vyZO1E0NMf0x/T//QE0VFDuvKhUVG68qIF+QFUEGT5EPKj+AAkpMjLH1JAyx9SMMv/UhD0AMntVPgPAdMHIcAAn2xRkyDXSpbTB9QC+wDoMOAhwAHjACHAAuMAAcADkTDjDQOkyMsfEssfy/8QERITAubQAdDTAyFxsJJfBOAi10nBIJJfBOAC0x8hghBwbHVnvSKCEGRzdHK9sJJfBeAD+kAwIPpEAcjKB8v/ydDtRNCBAUDXIfQEMFyBAQj0Cm+hMbOSXwfgBdM/yCWCEHBsdWe6kjgw4w0DghBkc3RyupJfBuMNBgcCASAICQB4AfoA9AQw+CdvIjBQCqEhvvLgUIIQcGx1Z4MesXCAGFAEywUmzxZY+gIZ9ADLaRfLH1Jgyz8gyYBA+wAGAIpQBIEBCPRZMO1E0IEBQNcgyAHPFvQAye1UAXKwjiOCEGRzdHKDHrFwgBhQBcsFUAPPFiP6AhPLassfyz/JgED7AJJfA+ICASAKCwBZvSQrb2omhAgKBrkPoCGEcNQICEekk30pkQzmkD6f+YN4EoAbeBAUiYcVnzGEAgFYDA0AEbjJftRNDXCx+AA9sp37UTQgQFA1yH0BDACyMoHy//J0AGBAQj0Cm+hMYAIBIA4PABmtznaiaEAga5Drhf/AABmvHfaiaEAQa5DrhY/AAG7SB/oA1NQi+QAFyMoHFcv/ydB3dIAYyMsFywIizxZQBfoCFMtrEszMyXP7AMhAFIEBCPRR8qcCAHCBAQjXGPoA0z/IVCBHgQEI9FHyp4IQbm90ZXB0gBjIywXLAlAGzxZQBPoCFMtqEssfyz/Jc/sAAgBsgQEI1xj6ANM/MFIkgQEI9Fnyp4IQZHN0cnB0gBjIywXLAlAFzxZQA/oCE8tqyx8Syz/Jc/sAAAr0AMntVGliJeU=",
	HighLoadV1R1: "te6ccgEBBgEAhgABFP8A9KQT9KDyyAsBAgEgAgMCAUgEBQC88oMI1xgg0x/TH9Mf+CMTu/Jj7UTQ0x/TH9P/0VEyuvKhUUS68qIE+QFUEFX5EPKj9ATR+AB/jhghgBD0eG+hb6EgmALTB9QwAfsAkTLiAbPmWwGkyMsfyx/L/8ntVAAE0DAAEaCZL9qJoa4WPw==",
	HighLoadV1R2: "te6ccgEBCAEAmQABFP8A9KQT9LzyyAsBAgEgAgMCAUgEBQC88oMI1xgg0x/TH9Mf+CMTu/Jj7UTQ0x/TH9P/0VEyuvKhUUS68qIE+QFUEFX5EPKj9ATR+AB/jhghgBD0eG+hb6EgmALTB9QwAfsAkTLiAbPmWwGkyMsfyx/L/8ntVAAE0DACAUgGBwAXuznO1E0NM/MdcL/4ABG4yX7UTQ1wsfg=",
	HighLoadV2:   "te6ccgEBCQEA5QABFP8A9KQT9LzyyAsBAgEgAgcCAUgDBAAE0DACASAFBgAXvZznaiaGmvmOuF/8AEG+X5dqJoaY+Y6Z/p/5j6AmipEEAgegc30JjJLb/JXdHxQB6vKDCNcYINMf0z/4I6ofUyC58mPtRNDTH9M/0//0BNFTYIBA9A5voTHyYFFzuvKiB/kBVBCH+RDyowL0BNH4AH+OFiGAEPR4b6UgmALTB9QwAfsAkTLiAbPmW4MlochANIBA9EOK5jEByMsfE8s/y//0AMntVAgANCCAQPSWb6VsEiCUMFMDud4gkzM2AZJsIeKz",
	HighLoadV2R1: "te6ccgEBBwEA1gABFP8A9KQT9KDyyAsBAgEgAgMCAUgEBQHu8oMI1xgg0x/TP/gjqh9TILnyY+1E0NMf0z/T//QE0VNggED0Dm+hMfJgUXO68qIH+QFUEIf5EPKjAvQE0fgAf44YIYAQ9HhvoW+hIJgC0wfUMAH7AJEy4gGz5luDJaHIQDSAQPRDiuYxyBLLHxPLP8v/9ADJ7VQGAATQMABBoZfl2omhpj5jpn+n/mPoCaKkQQCB6BzfQmMktv8ld0fFADgggED0lm+hb6EyURCUMFMDud4gkzM2AZIyMOKz",
	HighLoadV2R2: "te6ccgEBCQEA6QABFP8A9KQT9LzyyAsBAgEgAgMCAUgEBQHu8oMI1xgg0x/TP/gjqh9TILnyY+1E0NMf0z/T//QE0VNggED0Dm+hMfJgUXO68qIH+QFUEIf5EPKjAvQE0fgAf44YIYAQ9HhvoW+hIJgC0wfUMAH7AJEy4gGz5luDJaHIQDSAQPRDiuYxyBLLHxPLP8v/9ADJ7VQIAATQMAIBIAYHABe9nOdqJoaa+Y64X/wAQb5fl2omhpj5jpn+n/mPoCaKkQQCB6BzfQmMktv8ld0fFAA4IIBA9JZvoW+hMlEQlDBTA7neIJMzNgGSMjDisw==",
}

type blockchain interface {
	GetSeqno(ctx context.Context, account tongo.AccountID) (uint32, error)
	SendMessage(ctx context.Context, payload []byte) (uint32, error)
	GetAccountState(ctx context.Context, accountID tongo.AccountID) (tlb.ShardAccount, error)
}

func GetCodeByVer(ver Version) *boc.Cell {
	c, err := boc.DeserializeBocBase64(codes[ver])
	if err != nil {
		panic("invalid wallet hardcoded code")
	}
	if len(c) != 1 {
		panic("code must have one root cell")
	}
	return c[0]
}

func GetCodeHashByVer(ver Version) tlb.Bits256 {
	code := GetCodeByVer(ver)
	h, err := code.Hash()
	if err != nil {
		panic("can not calc hash for hardcoded code")
	}
	var hash tlb.Bits256
	copy(hash[:], h[:])
	return hash
}

func GetVerByCodeHash(hash tlb.Bits256) Version {
	// TODO: implement
	return 0
}

func (v Version) ToString() string {
	names := []string{"v1R1", "v1R2", "v1R3", "v2R1", "v2R2", "v3R1", "v3R2", "v4R1", "v4R2", "highload_v1R1", "highload_v1R2", "highload_v2", "highload_v2R1", "highload_v2R2"}
	if int(v) > len(names) {
		panic("to string conversion for this ver not supported")
	}
	return names[v]
}

type Wallet struct {
	key         ed25519.PrivateKey
	address     tongo.AccountID
	ver         Version
	subWalletId uint32
	blockchain  blockchain
}

// GetAddress returns current wallet address but you can also call function GenerateWalletAddress
// which returns same address but doesn't require blockchain connection for calling
func (w *Wallet) GetAddress() tongo.AccountID {
	return w.address
}

type DataV1V2 struct {
	Seqno     uint32
	PublicKey tlb.Bits256
}

type DataV3 struct {
	Seqno       uint32
	SubWalletId uint32
	PublicKey   tlb.Bits256
}

type DataV4 struct {
	Seqno       uint32
	SubWalletId uint32
	PublicKey   tlb.Bits256
	PluginDict  tlb.HashmapE[tlb.Bits264, tlb.Any] // TODO: find type and check size
}

type MessageV3 struct {
	SubWalletId uint32
	ValidUntil  uint32
	Seqno       uint32
	Payload     PayloadV1toV4
}

type MessageV4 struct {
	// Op: 0 - simple send, 1 - deploy and install plugin, 2 - install plugin, 3 - remove plugin
	SubWalletId uint32
	ValidUntil  uint32
	Seqno       uint32
	Op          int8
	Payload     PayloadV1toV4
}

type PayloadV1toV4 []RawMessage

type Sendable interface {
	ToInternal() (tlb.Message, uint8, error)
}

type SimpleTransfer struct {
	Amount  tlb.Grams
	Address tongo.AccountID
	Comment string
}

func (m SimpleTransfer) ToInternal() (message tlb.Message, mode uint8, err error) {
	info := tlb.CommonMsgInfo{
		SumType: "IntMsgInfo",
	}

	info.IntMsgInfo.IhrDisabled = true
	info.IntMsgInfo.Src = (*tongo.AccountID)(nil).ToMsgAddress()
	info.IntMsgInfo.Dest = m.Address.ToMsgAddress()
	info.IntMsgInfo.Value.Grams = m.Amount
	info.IntMsgInfo.Bounce = false

	intMsg := tlb.Message{
		Info: info,
	}

	if m.Comment != "" {
		body := boc.NewCell()
		err := tlb.Marshal(body, TextComment(m.Comment))
		if err != nil {
			return tlb.Message{}, 0, err
		}
		intMsg.Body.IsRight = true //todo: check length and
		intMsg.Body.Value = tlb.Any(*body)
	}
	return intMsg, DefaultMessageMode, nil
}

type Message struct {
	Amount  tlb.Grams
	Address tongo.AccountID
	Body    *boc.Cell
	Code    *boc.Cell
	Data    *boc.Cell
	Bounce  bool
	Mode    uint8
}

func (m Message) ToInternal() (message tlb.Message, mode uint8, err error) {
	info := tlb.CommonMsgInfo{
		SumType: "IntMsgInfo",
	}

	info.IntMsgInfo.IhrDisabled = true
	info.IntMsgInfo.Src = (*tongo.AccountID)(nil).ToMsgAddress()
	info.IntMsgInfo.Dest = m.Address.ToMsgAddress()
	info.IntMsgInfo.Value.Grams = m.Amount
	info.IntMsgInfo.Bounce = m.Bounce

	intMsg := tlb.Message{
		Info: info,
	}

	if m.Body != nil {
		intMsg.Body.IsRight = true //todo: check length and
		intMsg.Body.Value = tlb.Any(*m.Body)
	}
	if m.Code != nil && m.Data != nil {
		intMsg.Init.Exists = true
		intMsg.Init.Value.IsRight = true
		intMsg.Init.Value.Value.Code.Exists = true
		intMsg.Init.Value.Value.Data.Exists = true
		intMsg.Init.Value.Value.Code.Value.Value = *m.Code
		intMsg.Init.Value.Value.Data.Value.Value = *m.Data
	}

	return intMsg, m.Mode, nil
}

type RawMessage struct {
	Message *boc.Cell
	Mode    byte
}

type TextComment string

func (t TextComment) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error { // TODO: implement for binary comment
	err := c.WriteUint(0, 32) // text comment tag
	if err != nil {
		return err
	}
	return tlb.Marshal(c, tlb.Text(t))
}

func (t *TextComment) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error { // TODO: implement for binary comment
	val, err := c.ReadUint(32) // text comment tag
	if err != nil {
		return err
	}
	if val != 0 {
		return fmt.Errorf("not a text comment")
	}
	var text tlb.Text
	err = tlb.Unmarshal(c, &text)
	if err != nil {
		return err
	}
	*t = TextComment(text)
	return nil
}

func (p PayloadV1toV4) MarshalTLB(c *boc.Cell, encoder *tlb.Encoder) error {
	if len(p) > 4 {
		return fmt.Errorf("PayloadV1toV4 supports only up to 4 messages")
	}
	for _, msg := range p {
		err := c.WriteUint(uint64(msg.Mode), 8)
		if err != nil {
			return err
		}
		err = c.AddRef(msg.Message)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *PayloadV1toV4) UnmarshalTLB(c *boc.Cell, decoder *tlb.Decoder) error {
	for {
		ref, err := c.NextRef()
		if err != nil {
			break
		}
		mode, err := c.ReadUint(8)
		if err != nil {
			return err
		}
		msg := RawMessage{
			Message: ref,
			Mode:    byte(mode),
		}
		*p = append(*p, msg)
	}
	return nil
}
