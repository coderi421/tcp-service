package packet

import (
	"reflect"
	"testing"
)

func TestConAck_Decode(t *testing.T) {
	type args struct {
		connBody []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ConAckDecodeTest",
			args: args{
				connBody: []byte{'0', '0', '0', '0', '0', '0', '0', '1', 0},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConAck{}
			if err := c.Decode(tt.args.connBody); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConAck_Encode(t *testing.T) {
	type fields struct {
		ID     string
		Result uint8
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "ConAckEncodeTest",
			fields: fields{
				ID:     "00000001",
				Result: 0,
			},
			want:    []byte{'0', '0', '0', '0', '0', '0', '0', '1', 0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &ConAck{
				ID:     tt.fields.ID,
				Result: tt.fields.Result,
			}
			got, err := c.Encode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCon_Decode(t *testing.T) {
	type args struct{ connBody []byte }
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ConDecodeTest",
			args: args{
				connBody: []byte{'0', '0', '0', '0', '0', '0', '0', '1'},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Con{}
			if err := c.Decode(tt.args.connBody); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCon_Encode(t *testing.T) {
	type fields struct {
		ID      string
		Payload []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "ConEncodeTest",
			fields: fields{
				ID:      "00000001",
				Payload: nil,
			},
			want:    []byte{'0', '0', '0', '0', '0', '0', '0', '1'},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Con{
				ID:      tt.fields.ID,
				Payload: tt.fields.Payload,
			}
			got, err := c.Encode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecode(t *testing.T) {
	type args struct{ packet []byte }
	tests := []struct {
		name    string
		args    args
		want    Packet
		wantErr bool
	}{
		{
			name: "ConnDecodeTest",
			args: args{packet: []byte{CommandConn, '0', '0', '0', '0', '0', '0', '0', '1'}},
			want: &Con{
				ID:      "00000001",
				Payload: nil,
			},
			wantErr: false,
		},
		{
			name: "ConnAckDecodeTest",
			args: args{packet: []byte{CommandConnAck, '0', '0', '0', '0', '0', '0', '0', '1', 0}},
			want: &ConAck{
				ID:     "00000001",
				Result: 0,
			},
			wantErr: false,
		},
		{
			name: "SubmitDecodeTest",
			args: args{packet: []byte{CommandSubmit, '0', '0', '0', '0', '0', '0', '0', '1', 'h', 'e', 'l', 'l', 'o'}},
			want: &Submit{
				ID:      "00000001",
				Payload: []byte{'h', 'e', 'l', 'l', 'o'},
			},
			wantErr: false,
		},
		{
			name: "SubmitAckDecodeTest",
			args: args{packet: []byte{CommandSubmitAck, '0', '0', '0', '0', '0', '0', '0', '1', 0}},
			want: &SubmitAck{
				ID:     "00000001",
				Result: 0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.args.packet)
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

func TestEncode(t *testing.T) {
	type args struct{ p Packet }

	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ConnEncodeTest",
			args: args{
				p: &Con{ID: "00000001", Payload: nil},
			},
			want:    []byte{CommandConn, '0', '0', '0', '0', '0', '0', '0', '1'},
			wantErr: false,
		},
		{
			name: "ConnAckEncodeTest",
			args: args{
				p: &ConAck{ID: "00000001", Result: 0},
			},
			want:    []byte{CommandConnAck, '0', '0', '0', '0', '0', '0', '0', '1', 0},
			wantErr: false,
		},
		{
			name: "SubmitEncodeTest",
			args: args{
				p: &Con{ID: "00000001", Payload: []byte{'h', 'e', 'l', 'l', 'o'}},
			},
			want:    []byte{CommandConn, '0', '0', '0', '0', '0', '0', '0', '1', 'h', 'e', 'l', 'l', 'o'},
			wantErr: false,
		},
		{
			name: "SubmitAckEncodeTest",
			args: args{
				p: &ConAck{ID: "00000001", Result: 0},
			},
			want:    []byte{CommandConnAck, '0', '0', '0', '0', '0', '0', '0', '1', 0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encode(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubmitAck_Decode(t *testing.T) {
	type fields struct {
		ID     string
		Result uint8
	}
	type args struct {
		packetBody []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "SubmitAckDecodeTest",
			fields: fields{
				ID:     "00000001",
				Result: 0,
			},
			args: args{
				packetBody: []byte{'0', '0', '0', '0', '0', '0', '0', '1', 0},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SubmitAck{
				ID:     tt.fields.ID,
				Result: tt.fields.Result,
			}
			if err := s.Decode(tt.args.packetBody); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubmitAck_Encode(t *testing.T) {
	type fields struct {
		ID     string
		Result uint8
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "SubmitAckEncodeTest",
			fields: fields{
				ID:     "00000001",
				Result: 0,
			},
			want:    []byte{'0', '0', '0', '0', '0', '0', '0', '1', 0},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SubmitAck{
				ID:     tt.fields.ID,
				Result: tt.fields.Result,
			}
			got, err := s.Encode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubmit_Decode(t *testing.T) {
	type fields struct {
		ID      string
		Payload []byte
	}
	type args struct {
		packetBody []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "SubmitDecodeTest",
			fields: fields{
				ID:      "00000001",
				Payload: []byte{'h', 'e', 'l', 'l', 'o'},
			},
			args: args{
				packetBody: []byte{'0', '0', '0', '0', '0', '0', '0', '1', 'h', 'e', 'l', 'l', 'o'},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Submit{
				ID:      tt.fields.ID,
				Payload: tt.fields.Payload,
			}
			if err := s.Decode(tt.args.packetBody); (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubmit_Encode(t *testing.T) {
	type fields struct {
		ID      string
		Payload []byte
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "SubmitEncodeTest",
			fields: fields{
				ID:      "00000001",
				Payload: []byte{'h', 'e', 'l', 'l', 'o'},
			},
			want:    []byte{'0', '0', '0', '0', '0', '0', '0', '1', 'h', 'e', 'l', 'l', 'o'},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Submit{
				ID:      tt.fields.ID,
				Payload: tt.fields.Payload,
			}
			got, err := s.Encode()
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
