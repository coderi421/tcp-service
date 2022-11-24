package frame

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"
)

func TestNewCodec(t *testing.T) {
	tests := []struct {
		name string
		want StreamFrameCodec
	}{
		{name: "new", want: NewCodec()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCodec(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCodec() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestNewMyFrameCodec(t *testing.T)
// 	codec := NewCodec()
// 	if codec == nil {
// 		t.Errorf("want non-nil, actual nil")
// 	}
// }

func TestCodec_Encode(t *testing.T) {
	wantRes := []byte{0x0, 0x0, 0x0, 0x9, 'h', 'e', 'l', 'l', 'o'}
	data := []byte("hello")
	tests := []struct {
		name         string
		framePayload Payload
		wantW        string
		wantErr      bool
	}{
		{name: "EncodeHasErr", framePayload: data, wantW: string(wantRes), wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Codec{}
			// var w io.Writer
			w := &bytes.Buffer{}
			err := c.Encode(w, tt.framePayload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("Encode() gotW = %v, want %v", gotW, tt.wantW)
			}
			// var totalLen int32
			// _ = binary.Read(w, binary.BigEndian, &totalLen)
			// fmt.Println(string(w.Bytes()))
		})
	}
}

// func TestEncode(t *testing.T) {
// 	codec := NewCodec()
// 	buf := make([]byte, 0, 128)
// 	rw := bytes.NewBuffer(buf)
//
// 	err := codec.Encode(rw, []byte("hello"))
// 	if err != nil {
// 		t.Errorf("want nil, actual %s", err.Error())
// 	}
// 	// 验证Encode的正确性
// 	var totalLen int32
// 	err = binary.Read(rw, binary.BigEndian, &totalLen)
// 	if err != nil {
// 		t.Errorf("want nil, actual %s", err.Error())
// 	}
//
// 	if totalLen != 9 {
// 		t.Errorf("want 9, actual %d", totalLen)
// 	}
//
// 	left := rw.Bytes()
// 	if string(left) != "hello" {
// 		t.Errorf("want hello, actual %s", string(left))
// 	}
// }

func TestCodec_Decode(t *testing.T) {
	// 模拟一个 io.Reader
	data := []byte{0x0, 0x0, 0x0, 0x9, 'h', 'e', 'l', 'l', 'o'}
	wantRes := []byte{'h', 'e', 'l', 'l', 'o'}

	tests := []struct {
		name    string
		r       io.Reader
		want    Payload
		wantErr bool
	}{
		{name: "Decode", r: bytes.NewReader(data), want: wantRes, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Codec{}
			got, err := c.Decode(tt.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestDecode(t *testing.T) {
// 	codec := NewMyFrameCodec()
// 	data := []byte{0x0, 0x0, 0x0, 0x9, 'h', 'e', 'l', 'l', 'o'}
//
// 	payload, err := codec.Decode(bytes.NewReader(data))
// 	if err != nil {
// 		t.Errorf("want nil, actual %s", err.Error())
// 	}
//
// 	if string(payload) != "hello" {
// 		t.Errorf("want hello, actual %s", string(payload))
// 	}
// }

type ReturnErrorWriter struct {
	Writer           io.Writer
	NumberWriteError int // 第几次调用Write返回错误 Wn  NumberWriteError
	writeCount       int // 写操作次数计数 wc
}

func (w *ReturnErrorWriter) Write(p []byte) (n int, err error) {
	// 当前是第几次
	w.writeCount++
	if w.writeCount >= w.NumberWriteError {
		return 0, errors.New("write error")
	}
	return w.Writer.Write(p)
}

type ReturnErrorReader struct {
	Reader          io.Reader
	NumberReadError int // 第几次调用Read返回错误 Rn NumberReadError
	readCount       int // 读操作次数计数 rc
}

func (r *ReturnErrorReader) Read(p []byte) (n int, err error) {
	// 当前是第几次
	r.readCount++
	if r.readCount >= r.NumberReadError {
		return 0, errors.New("read error")
	}
	return r.Reader.Read(p)
}

func TestCodec_EncodeWithWriteFail(t *testing.T) {
	// wantRes := []byte{0x0, 0x0, 0x0, 0x9, 'h', 'e', 'l', 'l', 'o'}
	data := []byte("hello")
	tests := []struct {
		name         string
		framePayload Payload
		wantErr      bool
		writer       io.Writer
	}{
		// 测试第一次验证 binary.Write 写入错误
		{
			name:         "BinaryWriteErr",
			framePayload: data,
			wantErr:      true,
			writer: &ReturnErrorWriter{
				Writer:           &bytes.Buffer{},
				NumberWriteError: 1,
			},
		},
		// 测试第二次模拟w.Write返回错误
		{
			name:         "WriterWriteErr",
			framePayload: data,
			wantErr:      true,
			writer: &ReturnErrorWriter{
				Writer:           &bytes.Buffer{},
				NumberWriteError: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Codec{}
			// var w io.Writer
			err := c.Encode(tt.writer, tt.framePayload)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestCodec_DecodeWithReadFail(t *testing.T) {
	// 模拟一个 io.Reader
	data := []byte{0x0, 0x0, 0x0, 0x9, 'h', 'e', 'l', 'l', 'o'}
	// wantRes := []byte{'h', 'e', 'l', 'l', 'o'}

	tests := []struct {
		name    string
		r       io.Reader
		wantErr bool
	}{
		// 测试 binary.Read 返回错误
		{
			name: "BinaryReadErr",
			r: &ReturnErrorReader{
				Reader:          bytes.NewReader(data),
				NumberReadError: 1,
			},
			wantErr: true,
		},
		{
			name: "ReaderReadErr",
			r: &ReturnErrorReader{
				Reader:          bytes.NewReader(data),
				NumberReadError: 2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Codec{}
			_, err := c.Decode(tt.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
