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
	"fmt"
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
	ms := []*sproto.Payload_Metric{}

	for i, m := range p.Metrics {
		sm := sproto.Payload_Metric{}

		sm.Name = &p.Metrics[i].Name

		dt := m.DataType.toUint32()
		sm.Datatype = &dt
		sm.IsHistorical = proto.Bool(m.IsHistorical)
		sm.Metadata = encodeMetadata(m.Metadata)

		switch m.DataType {
		case TypeInt:
			iv, err := strconv.ParseUint(m.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			sm.Value = &sproto.Payload_Metric_IntValue{
				IntValue: uint32(iv),
			}

		case TypeFloat:
			fv, err := strconv.ParseFloat(m.Value, 32)
			if err != nil {
				return nil, err
			}
			sm.Value = &sproto.Payload_Metric_FloatValue{
				FloatValue: float32(fv),
			}

		case TypeBool:
			bv, err := strconv.ParseBool(m.Value)
			if err != nil {
				return nil, err
			}
			sm.Value = &sproto.Payload_Metric_BooleanValue{
				BooleanValue: bv,
			}

		case TypeString:
			sm.Value = &sproto.Payload_Metric_StringValue{
				StringValue: m.Value,
			}
		}

		ms = append(ms, &sm)
	}
	//fmt.Println("---------")

	sp := sproto.Payload{}

	if !isDeathPayload {
		// Set Payload timestamp
		tn := uint64(now)
		sp.Timestamp = &tn
		// Set Payload sequence
		sp.Seq = &p.Seq
	}
	sp.Metrics = ms
	return proto.Marshal(&sp)
}

func (p *Payload) DecodePayload(bytes []byte) error {
	pl := sproto.Payload{}
	proto.Unmarshal(bytes, &pl)
	// fmt.Println("Payload is ", pl.String())
	if pl.Timestamp != nil {
		p.Timestamp = time.UnixMilli(int64(*pl.Timestamp))
	}
	p.Metrics = make([]Metric, len(pl.Metrics))
	for i := range pl.Metrics {
		p.Metrics[i].Name = *pl.Metrics[i].Name
		p.Metrics[i].DataType = DataType(*pl.Metrics[i].Datatype)
		// Set the Value according to DataType
		switch p.Metrics[i].DataType {
		case TypeInt:
			p.Metrics[i].Value = strconv.FormatUint(uint64(pl.Metrics[i].GetIntValue()), 10)
		case TypeFloat:
			p.Metrics[i].Value = fmt.Sprintf("%f", pl.Metrics[i].GetFloatValue())
		case TypeBool:
			p.Metrics[i].Value = strconv.FormatBool(pl.Metrics[i].GetBooleanValue())
		case TypeString:
			p.Metrics[i].Value = pl.Metrics[i].GetStringValue()
		}
	}
	return nil
}
