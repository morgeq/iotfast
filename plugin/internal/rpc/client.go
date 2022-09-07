package rpc

import (
	"bytes"
	"encoding/gob"
	"io"
	"net/rpc"
	"sync"

	"github.com/morgeq/iotfast/plugin/internal/frame"

	"encoding/json"

	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
	"google.golang.org/protobuf/proto"
)

// ClientCodec is codec for goridge connection.
type ClientCodec struct {
	// bytes sync.Pool
	bPool sync.Pool
	fPool sync.Pool

	transfer frame.Transfer
	closed   bool
	frame    *frame.Frame
}

// NewClientCodec initiates new server rpc codec over socket connection.
func NewClientCodec(rwc io.ReadWriteCloser) *ClientCodec {
	return &ClientCodec{
		bPool: sync.Pool{New: func() interface{} {
			return new(bytes.Buffer)
		}},

		fPool: sync.Pool{New: func() interface{} {
			return frame.NewFrame()
		}},

		transfer: frame.NewTransferImp(rwc),
	}
}

func (c *ClientCodec) get() *bytes.Buffer {
	return c.bPool.Get().(*bytes.Buffer)
}

func (c *ClientCodec) put(b *bytes.Buffer) {
	b.Reset()
	c.bPool.Put(b)
}

func (c *ClientCodec) getFrame() *frame.Frame {
	return c.fPool.Get().(*frame.Frame)
}

func (c *ClientCodec) putFrame(f *frame.Frame) {
	f.Reset()
	c.fPool.Put(f)
}

// WriteRequest writes request to the connection. Sequential.
func (c *ClientCodec) WriteRequest(r *rpc.Request, body interface{}) error {
	const op = "goridge_write_request"

	// get a frame from the pool
	fr := c.getFrame()
	defer c.putFrame(fr)

	// get a buffer from the pool
	buf := c.get()
	defer c.put(buf)

	// writeServiceMethod to the buffer
	buf.WriteString(r.ServiceMethod)
	// use fallback as gob
	fr.WriteFlags(fr.Header(), frame.CodecGob)

	if body != nil {
		// if body is proto message, use proto codec
		switch m := body.(type) {
		// check if message is PROTO
		case proto.Message:
			fr.WriteFlags(fr.Header(), frame.CodecProto)
			b, err := proto.Marshal(m)
			if err != nil {
				return errors.New(op + err.Error())
			}
			buf.Write(b)
		default:
			enc := gob.NewEncoder(buf)
			// write data to the gob
			err := enc.Encode(body)
			if err != nil {
				return errors.New(op + err.Error())
			}
		}
	}

	// SEQ_ID + METHOD_NAME_LEN
	fr.WriteOptions(fr.HeaderPtr(), uint32(r.Seq), uint32(len(r.ServiceMethod)))
	fr.WriteVersion(fr.Header(), frame.Version1)

	fr.WritePayloadLen(fr.Header(), uint32(buf.Len()))
	fr.WritePayload(buf.Bytes())
	fr.WriteCRC(fr.Header())

	err := c.transfer.Send(fr)
	if err != nil {
		return errors.New(op + err.Error())
	}
	return nil
}

// ReadResponseHeader reads response from the connection.
func (c *ClientCodec) ReadResponseHeader(r *rpc.Response) error {
	const op = "client_read_response_header"

	// get a frame from sync.Pool
	fr := c.getFrame()

	err := c.transfer.Receive(fr)
	if err != nil {
		return errors.New(op + err.Error())
	}
	if !fr.VerifyCRC(fr.Header()) {
		return errors.New(op + ("CRC verification failed"))
	}

	// save the frame after CRC verification
	c.frame = fr

	opts := fr.ReadOptions(fr.Header())
	if len(opts) != 2 {
		return errors.New(op + ("should be 2 options. SEQ_ID and METHOD_LEN"))
	}

	// check for error
	if fr.ReadFlags()&frame.ERROR != 0 {
		r.Error = string(fr.Payload()[opts[1]:])
	}

	r.Seq = uint64(opts[0])
	r.ServiceMethod = string(fr.Payload()[:opts[1]])

	return nil
}

// ReadResponseBody response from the connection.
func (c *ClientCodec) ReadResponseBody(out interface{}) error {
	const op = "client_read_response_body"

	// put frame after response was sent
	defer c.putFrame(c.frame)
	// if there is no out interface to unmarshall the body, skip
	if out == nil {
		return nil
	}

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
		return errors.New(op + "unknown decoder used in frame")
	}
}

// Close closes the client connection.
func (c *ClientCodec) Close() error {
	if c.closed {
		return nil
	}

	c.closed = true
	return c.transfer.Close()
}
