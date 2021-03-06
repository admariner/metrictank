package importer

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *ChunkWriteRequest) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "ChunkWriteRequestPayload":
			err = z.ChunkWriteRequestPayload.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "ChunkWriteRequestPayload")
				return
			}
		case "Archive":
			err = z.Archive.DecodeMsg(dc)
			if err != nil {
				err = msgp.WrapError(err, "Archive")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *ChunkWriteRequest) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "ChunkWriteRequestPayload"
	err = en.Append(0x82, 0xb8, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64)
	if err != nil {
		return
	}
	err = z.ChunkWriteRequestPayload.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "ChunkWriteRequestPayload")
		return
	}
	// write "Archive"
	err = en.Append(0xa7, 0x41, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65)
	if err != nil {
		return
	}
	err = z.Archive.EncodeMsg(en)
	if err != nil {
		err = msgp.WrapError(err, "Archive")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ChunkWriteRequest) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "ChunkWriteRequestPayload"
	o = append(o, 0x82, 0xb8, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64)
	o, err = z.ChunkWriteRequestPayload.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "ChunkWriteRequestPayload")
		return
	}
	// string "Archive"
	o = append(o, 0xa7, 0x41, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65)
	o, err = z.Archive.MarshalMsg(o)
	if err != nil {
		err = msgp.WrapError(err, "Archive")
		return
	}
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ChunkWriteRequest) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "ChunkWriteRequestPayload":
			bts, err = z.ChunkWriteRequestPayload.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "ChunkWriteRequestPayload")
				return
			}
		case "Archive":
			bts, err = z.Archive.UnmarshalMsg(bts)
			if err != nil {
				err = msgp.WrapError(err, "Archive")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *ChunkWriteRequest) Msgsize() (s int) {
	s = 1 + 25 + z.ChunkWriteRequestPayload.Msgsize() + 8 + z.Archive.Msgsize()
	return
}
