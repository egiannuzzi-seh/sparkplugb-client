/*
Sparkplug 3.0.0
Note: Complies to v3.0.0 of the Sparkplug specification

	to the extent needed for Winsonic DataIO and other industrial 4.0 products.

Copyright (c) 2023 Winsonic Electronics, Taiwan
@author David Lee

* This program and the accompanying materials are made available under the
* terms of the Eclipse Public License 2.0 which is available at
* http://www.eclipse.org/legal/epl-2.0.
*/
package sparkplug

import (
	"strconv"
	"time"

	"github.com/egiannuzzi-seh/sparkplugb-client/sproto"
	"google.golang.org/protobuf/proto"
)

const namespace = "spBv1.0"
const state = "STATE"

type MessageType string

const (
	// Node message types
	MESSAGETYPE_NBIRTH = "NBIRTH"
	MESSAGETYPE_NDEATH = "NDEATH"
	MESSAGETYPE_NDATA  = "NDATA"
	MESSAGETYPE_NCMD   = "NCMD"
	// Device message types
	MESSAGETYPE_DBIRTH = "DBIRTH"
	MESSAGETYPE_DDEATH = "DDEATH"
	MESSAGETYPE_DDATA  = "DDATA"
	MESSAGETYPE_DCMD   = "DCMD"
)

type Payload struct {
	Timestamp time.Time
	Seq       uint64
	Metrics   []Metric
}

func encodeMetadata(m Metadata) *sproto.Payload_MetaData {
	return &sproto.Payload_MetaData{
		IsMultiPart: proto.Bool(m.IsMultiPart),
		ContentType: proto.String(m.ContentType),
		Size:        proto.Uint64(m.Size),
		Seq:         proto.Uint64(m.Seq),
		FileName:    proto.String(m.FileName),
		FileType:    proto.String(m.FileType),
		Md5:         proto.String(m.Md5),
		Description: proto.String(m.Description),
	}
}

func decodeMetadata(md *sproto.Payload_MetaData) Metadata {
	if md == nil {
		return Metadata{}
	}

	return Metadata{
		IsMultiPart: md.GetIsMultiPart(),
		ContentType: md.GetContentType(),
		Size:        md.GetSize(),
		Seq:         md.GetSeq(),
		FileName:    md.GetFileName(),
		FileType:    md.GetFileType(),
		Md5:         md.GetMd5(),
		Description: md.GetDescription(),
	}
}

func (p *Payload) EncodePayload(isDeathPayload bool) ([]byte, error) {
	now := time.Now().UnixMilli()
	ms := make([]*sproto.Payload_Metric, 0, len(p.Metrics))

	for i := range p.Metrics {
		// Use reference to avoid copying the struct
		m := &p.Metrics[i]
		sm := &sproto.Payload_Metric{}

		sm.Name = &m.Name

		dt := m.DataType.toUint32()
		sm.Datatype = &dt
		sm.IsHistorical = proto.Bool(m.IsHistorical)
		sm.Metadata = encodeMetadata(m.Metadata)

		switch m.DataType {
		case TypeInt8:
			v, err := strconv.ParseInt(m.Value, 10, 8)
			if err != nil {
				return nil, err
			}
			sm.Value = &sproto.Payload_Metric_IntValue{IntValue: uint32(v)}

		case TypeInt16:
			v, err := strconv.ParseInt(m.Value, 10, 16)
			if err != nil {
				return nil, err
			}
			sm.Value = &sproto.Payload_Metric_IntValue{IntValue: uint32(v)}

		case TypeInt32:
			v, err := strconv.ParseInt(m.Value, 10, 32)
			if err != nil {
				return nil, err
			}
			sm.Value = &sproto.Payload_Metric_IntValue{IntValue: uint32(v)}

		case TypeUInt8:
			v, err := strconv.ParseUint(m.Value, 10, 8)
			if err != nil {
				return nil, err
			}
			sm.Value = &sproto.Payload_Metric_IntValue{IntValue: uint32(v)}

		case TypeUInt16:
			v, err := strconv.ParseUint(m.Value, 10, 16)
			if err != nil {
				return nil, err
			}
			sm.Value = &sproto.Payload_Metric_IntValue{IntValue: uint32(v)}

		case TypeUInt32:
			v, err := strconv.ParseUint(m.Value, 10, 32)
			if err != nil {
				return nil, err
			}
			sm.Value = &sproto.Payload_Metric_IntValue{IntValue: uint32(v)}

		case TypeInt64:
			v, err := strconv.ParseInt(m.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			sm.Value = &sproto.Payload_Metric_LongValue{LongValue: uint64(v)}

		case TypeUInt64:
			v, err := strconv.ParseUint(m.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			sm.Value = &sproto.Payload_Metric_LongValue{LongValue: v}

		case TypeFloat:
			v, err := strconv.ParseFloat(m.Value, 32)
			if err != nil {
				return nil, err
			}
			sm.Value = &sproto.Payload_Metric_FloatValue{FloatValue: float32(v)}

		case TypeBool:
			v, err := strconv.ParseBool(m.Value)
			if err != nil {
				return nil, err
			}
			sm.Value = &sproto.Payload_Metric_BooleanValue{BooleanValue: v}

		case TypeString:
			sm.Value = &sproto.Payload_Metric_StringValue{StringValue: m.Value}

		case TypeBytes:
			sm.Value = &sproto.Payload_Metric_BytesValue{BytesValue: []byte(m.Value)}
		}

		ms = append(ms, sm)
	}

	sp := sproto.Payload{}

	if !isDeathPayload {
		tn := uint64(now)
		sp.Timestamp = &tn
		sp.Seq = &p.Seq
	}

	sp.Metrics = ms
	return proto.Marshal(&sp)
}

func (p *Payload) DecodePayload(bytes []byte) error {
	pl := sproto.Payload{}
	if err := proto.Unmarshal(bytes, &pl); err != nil {
		return err
	}

	if pl.Timestamp != nil {
		p.Timestamp = time.UnixMilli(int64(*pl.Timestamp))
	}

	p.Metrics = make([]Metric, len(pl.Metrics))

	for i, pm := range pl.Metrics {
		p.Metrics[i].Name = pm.GetName()
		p.Metrics[i].DataType = DataType(pm.GetDatatype())
		p.Metrics[i].IsHistorical = pm.GetIsHistorical()
		p.Metrics[i].Metadata = decodeMetadata(pm.GetMetadata())

		switch p.Metrics[i].DataType {
		case TypeInt8:
			p.Metrics[i].Value = strconv.FormatInt(int64(int8(pm.GetIntValue())), 10)
		case TypeInt16:
			p.Metrics[i].Value = strconv.FormatInt(int64(int16(pm.GetIntValue())), 10)
		case TypeInt32:
			p.Metrics[i].Value = strconv.FormatInt(int64(int32(pm.GetIntValue())), 10)
		case TypeUInt8:
			p.Metrics[i].Value = strconv.FormatUint(uint64(uint8(pm.GetIntValue())), 10)
		case TypeUInt16:
			p.Metrics[i].Value = strconv.FormatUint(uint64(uint16(pm.GetIntValue())), 10)
		case TypeUInt32:
			p.Metrics[i].Value = strconv.FormatUint(uint64(pm.GetIntValue()), 10)
		case TypeInt64:
			p.Metrics[i].Value = strconv.FormatInt(int64(pm.GetLongValue()), 10)
		case TypeUInt64:
			p.Metrics[i].Value = strconv.FormatUint(pm.GetLongValue(), 10)
		case TypeFloat:
			p.Metrics[i].Value = strconv.FormatFloat(float64(pm.GetFloatValue()), 'f', -1, 32)
		case TypeBool:
			p.Metrics[i].Value = strconv.FormatBool(pm.GetBooleanValue())
		case TypeString:
			p.Metrics[i].Value = pm.GetStringValue()
		case TypeBytes:
			p.Metrics[i].Value = string(pm.GetBytesValue())
		}
	}
	return nil
}
