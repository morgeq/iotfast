package rpc

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io"
	"net/rpc"
	"sync"

	"github.com/xiaodingding/iotfast/plugin/internal/frame"

	"encoding/json"

	"github.com/vmihailenco/msgpack/v5"
	"google.golang.org/protobuf/proto"
)

// Codec represent net/rpc bridge over Goridge socket transfer.
type Codec struct {
	transfer frame.Transfer
	closed   bool
	frame    *frame.Frame
	codec    sync.Map

	bPool sync.Pool
	fPool sync.Pool
}

// NewCodec initiates new server rpc codec over socket connection.
func NewCodec(rwc io.ReadWriteCloser) *Codec {
	return &Codec{
		transfer: frame.NewTransferImp(rwc),
		codec:    sync.Map{},

		bPool: sync.Pool{New: func() interface{} {
			return new(bytes.Buffer)
		}},

		fPool: sync.Pool{New: func() interface{} {
			return frame.NewFrame()
		}},
	}
}

// NewCodecWithRelay initiates new server rpc codec with a transfer of choice.
func NewCodecWithTransfer(transfer frame.Transfer) *Codec {
	return &Codec{transfer: transfer}
}

func (c *Codec) get() *bytes.Buffer {
	return c.bPool.Get().(*bytes.Buffer)
}

func (c *Codec) put(b *bytes.Buffer) {
	b.Reset()
	c.bPool.Put(b)
}

func (c *Codec) getFrame() *frame.Frame {
	return c.fPool.Get().(*frame.Frame)
}

func (c *Codec) putFrame(f *frame.Frame) {
	f.Reset()
	c.fPool.Put(f)
}

// WriteResponse marshals response, byte slice or error to remote party.
func (c *Codec) WriteResponse(r *rpc.Response, body interface{}) error { //nolint:funlen
	const op = ("goridge_write_response")
	fr := c.getFrame()
	defer c.putFrame(fr)

	// SEQ_ID + METHOD_NAME_LEN
	fr.WriteOptions(fr.HeaderPtr(), uint32(r.Seq), uint32(len(r.ServiceMethod)))
	// Write protocol version
	fr.WriteVersion(fr.Header(), frame.Version1)

	// load and delete associated codec to not waste memory
	// because we write it to the fr and don't need more information about it
	codec, ok := c.codec.LoadAndDelete(r.Seq)
	if !ok {
		// fallback codec
		fr.WriteFlags(fr.Header(), frame.CodecGob)
	} else {
		fr.WriteFlags(fr.Header(), codec.(byte))
	}

	// if error returned, we sending it via transfer and return error from WriteResponse
	if r.Error != "" {
		// Append error flag
		return c.handleError(r, fr, r.Error)
	}

	switch {
	case codec.(byte)&frame.CodecProto != 0:
		d, err := proto.Marshal(body.(proto.Message))
		if err != nil {
			return c.handleError(r, fr, err.Error())
		}

		// initialize buffer
		buf := c.get()
		defer c.put(buf)

		buf.Grow(len(d) + len(r.ServiceMethod))
		// writeServiceMethod to the buffer
		buf.WriteString(r.ServiceMethod)
		buf.Write(d)

		fr.WritePayloadLen(fr.Header(), uint32(buf.Len()))
		// copy inside
		fr.WritePayload(buf.Bytes())
		fr.WriteCRC(fr.Header())
		// send buffer
		return c.transfer.Send(fr)
	case codec.(byte)&frame.CodecRaw != 0:
		// initialize buffer
		buf := c.get()
		defer c.put(buf)

		switch data := body.(type) {
		case []byte:
			buf.Grow(len(data) + len(r.ServiceMethod))
			// writeServiceMethod to the buffer
			buf.WriteString(r.ServiceMethod)
			buf.Write(data)

			fr.WritePayloadLen(fr.Header(), uint32(buf.Len()))
			fr.WritePayload(buf.Bytes())
		case *[]byte:
			buf.Grow(len(*data) + len(r.ServiceMethod))
			// writeServiceMethod to the buffer
			buf.WriteString(r.ServiceMethod)
			buf.Write(*data)

			fr.WritePayloadLen(fr.Header(), uint32(buf.Len()))
			fr.WritePayload(buf.Bytes())
		default:
			return c.handleError(r, fr, "unknown Raw payload type")
		}

		// send buffer
		fr.WriteCRC(fr.Header())
		return c.transfer.Send(fr)

	case codec.(byte)&frame.CodecJSON != 0:
		data, err := json.Marshal(body)
		if err != nil {
			return c.handleError(r, fr, err.Error())
		}

		// initialize buffer
		buf := c.get()
		defer c.put(buf)

		buf.Grow(len(data) + len(r.ServiceMethod))
		// writeServiceMethod to the buffer
		buf.WriteString(r.ServiceMethod)
		buf.Write(data)

		fr.WritePayloadLen(fr.Header(), uint32(buf.Len()))
		// copy inside
		fr.WritePayload(buf.Bytes())
		fr.WriteCRC(fr.Header())
		// send buffer
		return c.transfer.Send(fr)

	case codec.(byte)&frame.CodecMsgpack != 0:
		b, err := msgpack.Marshal(body)
		if err != nil {
			return errors.New(op + err.Error())
		}
		// initialize buffer
		buf := c.get()
		defer c.put(buf)

		buf.Grow(len(b) + len(r.ServiceMethod))
		// writeServiceMethod to the buffer
		buf.WriteString(r.ServiceMethod)
		buf.Write(b)

		fr.WritePayloadLen(fr.Header(), uint32(buf.Len()))
		// copy inside
		fr.WritePayload(buf.Bytes())
		fr.WriteCRC(fr.Header())
		// send buffer
		return c.transfer.Send(fr)

	case codec.(byte)&frame.CodecGob != 0:
		// initialize buffer
		buf := c.get()
		defer c.put(buf)

		buf.WriteString(r.ServiceMethod)

		dec := gob.NewEncoder(buf)
		err := dec.Encode(body)
		if err != nil {
			return errors.New(op + err.Error())
		}

		fr.WritePayloadLen(fr.Header(), uint32(buf.Len()))
		// copy inside
		fr.WritePayload(buf.Bytes())
		fr.WriteCRC(fr.Header())
		// send buffer
		return c.transfer.Send(fr)
	default:
		return c.handleError(r, fr, (op + "unknown codec"))
	}
}

func (c *Codec) handleError(r *rpc.Response, fr *frame.Frame, err string) error {
	buf := c.get()
	defer c.put(buf)

	// write all possible errors
	buf.WriteString(r.ServiceMethod)

	const op = ("handle codec error")
	fr.WriteFlags(fr.Header(), frame.ERROR)
	// error should be here
	if err != "" {
		buf.WriteString(err)
	}
	fr.WritePayloadLen(fr.Header(), uint32(buf.Len()))
	fr.WritePayload(buf.Bytes())

	fr.WriteCRC(fr.Header())
	_ = c.transfer.Send(fr)
	return errors.New(op + r.Error)
}

// ReadRequestHeader receives frame with options
// options should have 2 values
// [0] - integer, sequence ID
// [1] - integer, offset for method name
// For example:
// 15Test.Payload
// SEQ_ID: 15
// METHOD_LEN: 12 and we take 12 bytes from the payload as method name
func (c *Codec) ReadRequestHeader(r *rpc.Request) error {
	const op = ("goridge_read_request_header")
	f := c.getFrame()

	err := c.transfer.Receive(f)
	if err != nil {
		if errors.Is(err, io.EOF) {
			c.putFrame(f)
			return err
		}

		c.putFrame(f)
		return err
	}

	// opts[0] sequence ID
	// opts[1] service method name offset from payload in bytes
	opts := f.ReadOptions(f.Header())
	if len(opts) != 2 {
		c.putFrame(f)
		return errors.New(op + ("should be 2 options. SEQ_ID and METHOD_LEN"))
	}

	r.Seq = uint64(opts[0])
	r.ServiceMethod = string(f.Payload()[:opts[1]])
	c.frame = f
	return c.storeCodec(r, f.ReadFlags())
}

func (c *Codec) storeCodec(r *rpc.Request, flag byte) error {
	switch {
	case flag&frame.CodecProto != 0:
		c.codec.Store(r.Seq, frame.CodecProto)
	case flag&frame.CodecJSON != 0:
		c.codec.Store(r.Seq, frame.CodecJSON)
	case flag&frame.CodecRaw != 0:
		c.codec.Store(r.Seq, frame.CodecRaw)
	case flag&frame.CodecMsgpack != 0:
		c.codec.Store(r.Seq, frame.CodecMsgpack)
	case flag&frame.CodecGob != 0:
		c.codec.Store(r.Seq, frame.CodecGob)
	default:
		c.codec.Store(r.Seq, frame.CodecGob)
	}

	return nil
}

// ReadRequestBody fetches prefixed body data and automatically unmarshal it as json. RawBody flag will populate
// []byte lice argument for rpc method.
func (c *Codec) ReadRequestBody(out interface{}) error {
	const op = ("goridge_read_request_body")
	if out == nil {
		return nil
	}

	defer c.putFrame(c.frame)

	flags := c.frame.ReadFlags()

	switch { //nolint:dupl
	case flags&frame.CodecProto != 0:
		opts := c.frame.ReadOptions(c.frame.Header())
		if len(opts) != 2 {
			return errors.New(op + ("should be 2 options. SEQ_ID and METHOD_LEN"))
		}
		payload := c.frame.Payload()[opts[1]:]
		if len(payload) == 0 {
			return nil
		}

		// check if the out message is a correct proto.Message
		// instead send an error
		if pOut, ok := out.(proto.Message); ok {
			err := proto.Unmarshal(payload, pOut)
			if err != nil {
				return errors.New(op + err.Error())
			}
			return nil
		}

		return errors.New(op + ("message type is not a proto"))
	case flags&frame.CodecJSON != 0:
		opts := c.frame.ReadOptions(c.frame.Header())
		if len(opts) != 2 {
			return errors.New(op + ("should be 2 options. SEQ_ID and METHOD_LEN"))
		}
		payload := c.frame.Payload()[opts[1]:]
		if len(payload) == 0 {
			return nil
		}
		return json.Unmarshal(payload, out)
	case flags&frame.CodecGob != 0:
		opts := c.frame.ReadOptions(c.frame.Header())
		if len(opts) != 2 {
			return errors.New(op + ("should be 2 options. SEQ_ID and METHOD_LEN"))
		}
		payload := c.frame.Payload()[opts[1]:]
		if len(payload) == 0 {
			return nil
		}

		buf := c.get()
		defer c.put(buf)

		dec := gob.NewDecoder(buf)
		buf.Write(payload)

		err := dec.Decode(out)
		if err != nil {
			return errors.New(op + err.Error())
		}

		return nil
	case flags&frame.CodecRaw != 0:
		opts := c.frame.ReadOptions(c.frame.Header())
		if len(opts) != 2 {
			return errors.New(op + ("should be 2 options. SEQ_ID and METHOD_LEN"))
		}
		payload := c.frame.Payload()[opts[1]:]
		if len(payload) == 0 {
			return nil
		}

		if raw, ok := out.(*[]byte); ok {
			*raw = append(*raw, payload...)
		}

		return nil
	case flags&frame.CodecMsgpack != 0:
		opts := c.frame.ReadOptions(c.frame.Header())
		if len(opts) != 2 {
			return errors.New(op + ("should be 2 options. SEQ_ID and METHOD_LEN"))
		}
		payload := c.frame.Payload()[opts[1]:]
		if len(payload) == 0 {
			return nil
		}

		return msgpack.Unmarshal(payload, out)
	default:
		return errors.New(op + ("unknown decoder used in frame"))
	}
}

// Close underlying socket.
func (c *Codec) Close() error {
	if c.closed {
		return nil
	}

	c.closed = true
	return c.transfer.Close()
}
