// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: side/btclightclient/params.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Params defines the parameters for the module.
type Params struct {
	// Only accept blocks sending from these addresses
	Senders []string `protobuf:"bytes,1,rep,name=senders,proto3" json:"senders,omitempty"`
	// The minimum number of confirmations required for a block to be accepted
	Confirmations int32 `protobuf:"varint,2,opt,name=confirmations,proto3" json:"confirmations,omitempty"`
	// Indicates the maximum depth or distance from the latest block up to which transactions are considered for acceptance.
	MaxAcceptableBlockDepth uint64 `protobuf:"varint,3,opt,name=max_acceptable_block_depth,json=maxAcceptableBlockDepth,proto3" json:"max_acceptable_block_depth,omitempty"`
	// the denomanation of the voucher
	BtcVoucherDenom string `protobuf:"bytes,4,opt,name=btc_voucher_denom,json=btcVoucherDenom,proto3" json:"btc_voucher_denom,omitempty"`
	// the address to which the voucher is sent
	BtcVoucherAddress []string `protobuf:"bytes,5,rep,name=btc_voucher_address,json=btcVoucherAddress,proto3" json:"btc_voucher_address,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b47f8b78acf6f6e, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetSenders() []string {
	if m != nil {
		return m.Senders
	}
	return nil
}

func (m *Params) GetConfirmations() int32 {
	if m != nil {
		return m.Confirmations
	}
	return 0
}

func (m *Params) GetMaxAcceptableBlockDepth() uint64 {
	if m != nil {
		return m.MaxAcceptableBlockDepth
	}
	return 0
}

func (m *Params) GetBtcVoucherDenom() string {
	if m != nil {
		return m.BtcVoucherDenom
	}
	return ""
}

func (m *Params) GetBtcVoucherAddress() []string {
	if m != nil {
		return m.BtcVoucherAddress
	}
	return nil
}

// Bitcoin Block Header
type BlockHeader struct {
	Version           uint64 `protobuf:"varint,1,opt,name=version,proto3" json:"version,omitempty"`
	Hash              string `protobuf:"bytes,2,opt,name=hash,proto3" json:"hash,omitempty"`
	Height            uint64 `protobuf:"varint,3,opt,name=height,proto3" json:"height,omitempty"`
	PreviousBlockHash string `protobuf:"bytes,4,opt,name=previous_block_hash,json=previousBlockHash,proto3" json:"previous_block_hash,omitempty"`
	MerkleRoot        string `protobuf:"bytes,5,opt,name=merkle_root,json=merkleRoot,proto3" json:"merkle_root,omitempty"`
	Nonce             uint64 `protobuf:"varint,6,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Bits              string `protobuf:"bytes,7,opt,name=bits,proto3" json:"bits,omitempty"`
	Time              uint64 `protobuf:"varint,8,opt,name=time,proto3" json:"time,omitempty"`
	Ntx               uint64 `protobuf:"varint,9,opt,name=ntx,proto3" json:"ntx,omitempty"`
}

func (m *BlockHeader) Reset()         { *m = BlockHeader{} }
func (m *BlockHeader) String() string { return proto.CompactTextString(m) }
func (*BlockHeader) ProtoMessage()    {}
func (*BlockHeader) Descriptor() ([]byte, []int) {
	return fileDescriptor_3b47f8b78acf6f6e, []int{1}
}
func (m *BlockHeader) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *BlockHeader) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_BlockHeader.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *BlockHeader) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BlockHeader.Merge(m, src)
}
func (m *BlockHeader) XXX_Size() int {
	return m.Size()
}
func (m *BlockHeader) XXX_DiscardUnknown() {
	xxx_messageInfo_BlockHeader.DiscardUnknown(m)
}

var xxx_messageInfo_BlockHeader proto.InternalMessageInfo

func (m *BlockHeader) GetVersion() uint64 {
	if m != nil {
		return m.Version
	}
	return 0
}

func (m *BlockHeader) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *BlockHeader) GetHeight() uint64 {
	if m != nil {
		return m.Height
	}
	return 0
}

func (m *BlockHeader) GetPreviousBlockHash() string {
	if m != nil {
		return m.PreviousBlockHash
	}
	return ""
}

func (m *BlockHeader) GetMerkleRoot() string {
	if m != nil {
		return m.MerkleRoot
	}
	return ""
}

func (m *BlockHeader) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *BlockHeader) GetBits() string {
	if m != nil {
		return m.Bits
	}
	return ""
}

func (m *BlockHeader) GetTime() uint64 {
	if m != nil {
		return m.Time
	}
	return 0
}

func (m *BlockHeader) GetNtx() uint64 {
	if m != nil {
		return m.Ntx
	}
	return 0
}

func init() {
	proto.RegisterType((*Params)(nil), "side.btclightclient.Params")
	proto.RegisterType((*BlockHeader)(nil), "side.btclightclient.BlockHeader")
}

func init() { proto.RegisterFile("side/btclightclient/params.proto", fileDescriptor_3b47f8b78acf6f6e) }

var fileDescriptor_3b47f8b78acf6f6e = []byte{
	// 421 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x52, 0x3f, 0x6f, 0xd3, 0x40,
	0x14, 0xcf, 0x91, 0x3f, 0x25, 0x57, 0x21, 0xc8, 0xb5, 0x82, 0x53, 0x07, 0x63, 0x55, 0x0c, 0x11,
	0x43, 0x3c, 0x74, 0x64, 0x6a, 0xd5, 0x81, 0x05, 0x09, 0x79, 0x60, 0x60, 0xb1, 0xce, 0xe7, 0x47,
	0x7c, 0xaa, 0x7d, 0xcf, 0xba, 0xbb, 0x44, 0xe6, 0x5b, 0xf0, 0xb1, 0x18, 0x3b, 0x32, 0x42, 0xf2,
	0x15, 0xf8, 0x00, 0xe8, 0x9e, 0x13, 0xa0, 0x5d, 0xa2, 0xdf, 0xbf, 0xcb, 0xfb, 0x3d, 0xeb, 0xf1,
	0xd4, 0x9b, 0x0a, 0xb2, 0x32, 0xe8, 0xc6, 0xac, 0xeb, 0xf8, 0x0b, 0x36, 0x64, 0x9d, 0x72, 0xaa,
	0xf5, 0xab, 0xce, 0x61, 0x40, 0x71, 0x16, 0x13, 0xab, 0x87, 0x89, 0x8b, 0xf3, 0x35, 0xae, 0x91,
	0xfc, 0x2c, 0xa2, 0x21, 0x7a, 0xf9, 0x8b, 0xf1, 0xd9, 0x47, 0x7a, 0x2b, 0x24, 0x3f, 0xf1, 0x60,
	0x2b, 0x70, 0x5e, 0xb2, 0x74, 0xbc, 0x9c, 0xe7, 0x47, 0x2a, 0xde, 0xf0, 0x67, 0x1a, 0xed, 0x17,
	0xe3, 0x5a, 0x15, 0x0c, 0x5a, 0x2f, 0x9f, 0xa4, 0x6c, 0x39, 0xcd, 0x1f, 0x8a, 0xe2, 0x1d, 0xbf,
	0x68, 0x55, 0x5f, 0x28, 0xad, 0xa1, 0x0b, 0xaa, 0x6c, 0xa0, 0x28, 0x1b, 0xd4, 0x77, 0x45, 0x05,
	0x5d, 0xa8, 0xe5, 0x38, 0x65, 0xcb, 0x49, 0xfe, 0xaa, 0x55, 0xfd, 0xf5, 0xdf, 0xc0, 0x4d, 0xf4,
	0x6f, 0xa3, 0x2d, 0xde, 0xf2, 0x45, 0x19, 0x74, 0xb1, 0xc5, 0x8d, 0xae, 0xc1, 0x15, 0x15, 0x58,
	0x6c, 0xe5, 0x24, 0x65, 0xcb, 0x79, 0xfe, 0xbc, 0x0c, 0xfa, 0xd3, 0xa0, 0xdf, 0x46, 0x59, 0xac,
	0xf8, 0xd9, 0xff, 0x59, 0x55, 0x55, 0x0e, 0xbc, 0x97, 0x53, 0x2a, 0xbd, 0xf8, 0x97, 0xbe, 0x1e,
	0x8c, 0xcb, 0xdf, 0x8c, 0x9f, 0xd2, 0xa8, 0xf7, 0xa0, 0x2a, 0x70, 0x71, 0xd1, 0x2d, 0x38, 0x6f,
	0xd0, 0x4a, 0x46, 0xad, 0x8e, 0x54, 0x08, 0x3e, 0xa9, 0x95, 0xaf, 0x69, 0xbf, 0x79, 0x4e, 0x58,
	0xbc, 0xe4, 0xb3, 0x1a, 0xe2, 0x77, 0x3c, 0xac, 0x70, 0x60, 0xb1, 0x45, 0xe7, 0x60, 0x6b, 0x70,
	0xe3, 0x0f, 0x8b, 0xd2, 0xd3, 0xa1, 0xf3, 0xe2, 0x68, 0x0d, 0x73, 0xe3, 0xff, 0xbc, 0xe6, 0xa7,
	0x2d, 0xb8, 0xbb, 0x06, 0x0a, 0x87, 0x18, 0xe4, 0x94, 0x72, 0x7c, 0x90, 0x72, 0xc4, 0x20, 0xce,
	0xf9, 0xd4, 0xa2, 0xd5, 0x20, 0x67, 0x34, 0x67, 0x20, 0xb1, 0x52, 0x69, 0x82, 0x97, 0x27, 0x43,
	0xa5, 0x88, 0xa3, 0x16, 0x4c, 0x0b, 0xf2, 0x29, 0x05, 0x09, 0x8b, 0x17, 0x7c, 0x6c, 0x43, 0x2f,
	0xe7, 0x24, 0x45, 0x78, 0xf3, 0xe1, 0xfb, 0x2e, 0x61, 0xf7, 0xbb, 0x84, 0xfd, 0xdc, 0x25, 0xec,
	0xdb, 0x3e, 0x19, 0xdd, 0xef, 0x93, 0xd1, 0x8f, 0x7d, 0x32, 0xfa, 0x7c, 0xb5, 0x36, 0xa1, 0xde,
	0x94, 0x2b, 0x8d, 0x6d, 0x16, 0x4f, 0x85, 0x4e, 0x41, 0x63, 0x43, 0x24, 0xeb, 0x1f, 0xdf, 0x56,
	0xf8, 0xda, 0x81, 0x2f, 0x67, 0x94, 0xba, 0xfa, 0x13, 0x00, 0x00, 0xff, 0xff, 0x68, 0x1a, 0xad,
	0x49, 0x7f, 0x02, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.BtcVoucherAddress) > 0 {
		for iNdEx := len(m.BtcVoucherAddress) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.BtcVoucherAddress[iNdEx])
			copy(dAtA[i:], m.BtcVoucherAddress[iNdEx])
			i = encodeVarintParams(dAtA, i, uint64(len(m.BtcVoucherAddress[iNdEx])))
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.BtcVoucherDenom) > 0 {
		i -= len(m.BtcVoucherDenom)
		copy(dAtA[i:], m.BtcVoucherDenom)
		i = encodeVarintParams(dAtA, i, uint64(len(m.BtcVoucherDenom)))
		i--
		dAtA[i] = 0x22
	}
	if m.MaxAcceptableBlockDepth != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.MaxAcceptableBlockDepth))
		i--
		dAtA[i] = 0x18
	}
	if m.Confirmations != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.Confirmations))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Senders) > 0 {
		for iNdEx := len(m.Senders) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Senders[iNdEx])
			copy(dAtA[i:], m.Senders[iNdEx])
			i = encodeVarintParams(dAtA, i, uint64(len(m.Senders[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *BlockHeader) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *BlockHeader) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *BlockHeader) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Ntx != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.Ntx))
		i--
		dAtA[i] = 0x48
	}
	if m.Time != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.Time))
		i--
		dAtA[i] = 0x40
	}
	if len(m.Bits) > 0 {
		i -= len(m.Bits)
		copy(dAtA[i:], m.Bits)
		i = encodeVarintParams(dAtA, i, uint64(len(m.Bits)))
		i--
		dAtA[i] = 0x3a
	}
	if m.Nonce != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.Nonce))
		i--
		dAtA[i] = 0x30
	}
	if len(m.MerkleRoot) > 0 {
		i -= len(m.MerkleRoot)
		copy(dAtA[i:], m.MerkleRoot)
		i = encodeVarintParams(dAtA, i, uint64(len(m.MerkleRoot)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.PreviousBlockHash) > 0 {
		i -= len(m.PreviousBlockHash)
		copy(dAtA[i:], m.PreviousBlockHash)
		i = encodeVarintParams(dAtA, i, uint64(len(m.PreviousBlockHash)))
		i--
		dAtA[i] = 0x22
	}
	if m.Height != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Hash) > 0 {
		i -= len(m.Hash)
		copy(dAtA[i:], m.Hash)
		i = encodeVarintParams(dAtA, i, uint64(len(m.Hash)))
		i--
		dAtA[i] = 0x12
	}
	if m.Version != 0 {
		i = encodeVarintParams(dAtA, i, uint64(m.Version))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Senders) > 0 {
		for _, s := range m.Senders {
			l = len(s)
			n += 1 + l + sovParams(uint64(l))
		}
	}
	if m.Confirmations != 0 {
		n += 1 + sovParams(uint64(m.Confirmations))
	}
	if m.MaxAcceptableBlockDepth != 0 {
		n += 1 + sovParams(uint64(m.MaxAcceptableBlockDepth))
	}
	l = len(m.BtcVoucherDenom)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if len(m.BtcVoucherAddress) > 0 {
		for _, s := range m.BtcVoucherAddress {
			l = len(s)
			n += 1 + l + sovParams(uint64(l))
		}
	}
	return n
}

func (m *BlockHeader) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Version != 0 {
		n += 1 + sovParams(uint64(m.Version))
	}
	l = len(m.Hash)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.Height != 0 {
		n += 1 + sovParams(uint64(m.Height))
	}
	l = len(m.PreviousBlockHash)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	l = len(m.MerkleRoot)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.Nonce != 0 {
		n += 1 + sovParams(uint64(m.Nonce))
	}
	l = len(m.Bits)
	if l > 0 {
		n += 1 + l + sovParams(uint64(l))
	}
	if m.Time != 0 {
		n += 1 + sovParams(uint64(m.Time))
	}
	if m.Ntx != 0 {
		n += 1 + sovParams(uint64(m.Ntx))
	}
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Senders", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Senders = append(m.Senders, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Confirmations", wireType)
			}
			m.Confirmations = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Confirmations |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxAcceptableBlockDepth", wireType)
			}
			m.MaxAcceptableBlockDepth = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaxAcceptableBlockDepth |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BtcVoucherDenom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BtcVoucherDenom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BtcVoucherAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BtcVoucherAddress = append(m.BtcVoucherAddress, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *BlockHeader) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: BlockHeader: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BlockHeader: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Version", wireType)
			}
			m.Version = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Version |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Hash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Hash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PreviousBlockHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PreviousBlockHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MerkleRoot", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MerkleRoot = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			m.Nonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Nonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Bits", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Bits = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 8:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Time", wireType)
			}
			m.Time = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Time |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 9:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Ntx", wireType)
			}
			m.Ntx = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Ntx |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowParams
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowParams
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)