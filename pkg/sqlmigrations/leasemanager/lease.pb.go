// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: sqlmigrations/leasemanager/lease.proto

package leasemanager

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import hlc "github.com/cockroachdb/cockroach/pkg/util/hlc"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type LeaseVal struct {
	// An opaque string that should be unique per client to identify which client
	// owns the lease.
	Owner string `protobuf:"bytes,1,opt,name=owner" json:"owner"`
	// The expiration time of the lease.
	Expiration hlc.Timestamp `protobuf:"bytes,2,opt,name=expiration" json:"expiration"`
}

func (m *LeaseVal) Reset()         { *m = LeaseVal{} }
func (m *LeaseVal) String() string { return proto.CompactTextString(m) }
func (*LeaseVal) ProtoMessage()    {}
func (*LeaseVal) Descriptor() ([]byte, []int) {
	return fileDescriptor_lease_fab267a03cfd3cca, []int{0}
}
func (m *LeaseVal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *LeaseVal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalTo(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (dst *LeaseVal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LeaseVal.Merge(dst, src)
}
func (m *LeaseVal) XXX_Size() int {
	return m.Size()
}
func (m *LeaseVal) XXX_DiscardUnknown() {
	xxx_messageInfo_LeaseVal.DiscardUnknown(m)
}

var xxx_messageInfo_LeaseVal proto.InternalMessageInfo

func init() {
	proto.RegisterType((*LeaseVal)(nil), "cockroach.client.LeaseVal")
}
func (m *LeaseVal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *LeaseVal) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintLease(dAtA, i, uint64(len(m.Owner)))
	i += copy(dAtA[i:], m.Owner)
	dAtA[i] = 0x12
	i++
	i = encodeVarintLease(dAtA, i, uint64(m.Expiration.Size()))
	n1, err := m.Expiration.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	return i, nil
}

func encodeVarintLease(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *LeaseVal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Owner)
	n += 1 + l + sovLease(uint64(l))
	l = m.Expiration.Size()
	n += 1 + l + sovLease(uint64(l))
	return n
}

func sovLease(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozLease(x uint64) (n int) {
	return sovLease(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *LeaseVal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowLease
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: LeaseVal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: LeaseVal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLease
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthLease
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Expiration", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowLease
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthLease
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Expiration.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipLease(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthLease
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
func skipLease(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowLease
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
					return 0, ErrIntOverflowLease
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowLease
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthLease
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowLease
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipLease(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthLease = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowLease   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("sqlmigrations/leasemanager/lease.proto", fileDescriptor_lease_fab267a03cfd3cca)
}

var fileDescriptor_lease_fab267a03cfd3cca = []byte{
	// 223 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x52, 0x2b, 0x2e, 0xcc, 0xc9,
	0xcd, 0x4c, 0x2f, 0x4a, 0x2c, 0xc9, 0xcc, 0xcf, 0x2b, 0xd6, 0xcf, 0x49, 0x4d, 0x2c, 0x4e, 0xcd,
	0x4d, 0xcc, 0x4b, 0x4c, 0x4f, 0x2d, 0x82, 0x70, 0xf4, 0x0a, 0x8a, 0xf2, 0x4b, 0xf2, 0x85, 0x04,
	0x92, 0xf3, 0x93, 0xb3, 0x8b, 0xf2, 0x13, 0x93, 0x33, 0xf4, 0x92, 0x73, 0x32, 0x53, 0xf3, 0x4a,
	0xa4, 0x24, 0x4a, 0x4b, 0x32, 0x73, 0xf4, 0x33, 0x72, 0x92, 0xf5, 0x4b, 0x32, 0x73, 0x53, 0x8b,
	0x4b, 0x12, 0x73, 0x0b, 0x20, 0x6a, 0xa5, 0x44, 0xd2, 0xf3, 0xd3, 0xf3, 0xc1, 0x4c, 0x7d, 0x10,
	0x0b, 0x22, 0xaa, 0x94, 0xcd, 0xc5, 0xe1, 0x03, 0x32, 0x30, 0x2c, 0x31, 0x47, 0x48, 0x8a, 0x8b,
	0x35, 0xbf, 0x3c, 0x2f, 0xb5, 0x48, 0x82, 0x51, 0x81, 0x51, 0x83, 0xd3, 0x89, 0xe5, 0xc4, 0x3d,
	0x79, 0x86, 0x20, 0x88, 0x90, 0x90, 0x33, 0x17, 0x57, 0x6a, 0x45, 0x41, 0x26, 0xc4, 0x45, 0x12,
	0x4c, 0x0a, 0x8c, 0x1a, 0xdc, 0x46, 0xb2, 0x7a, 0x08, 0xeb, 0x41, 0xd6, 0xea, 0x65, 0xe4, 0x24,
	0xeb, 0x85, 0xc0, 0xac, 0x85, 0xea, 0x47, 0xd2, 0xe6, 0xa4, 0x77, 0xe2, 0xa1, 0x1c, 0xc3, 0x89,
	0x47, 0x72, 0x8c, 0x17, 0x1e, 0xc9, 0x31, 0xde, 0x78, 0x24, 0xc7, 0xf8, 0xe0, 0x91, 0x1c, 0xe3,
	0x84, 0xc7, 0x72, 0x0c, 0x17, 0x1e, 0xcb, 0x31, 0xdc, 0x78, 0x2c, 0xc7, 0x10, 0xc5, 0x83, 0xec,
	0x55, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0x31, 0xb7, 0xb2, 0xb1, 0x07, 0x01, 0x00, 0x00,
}
