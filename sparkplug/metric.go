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

import "fmt"

type Metadata struct {
	IsMultiPart bool
	ContentType string
	Size        uint64
	Seq         uint64
	FileName    string
	FileType    string
	Md5         string
	Description string
}

type Metric struct {
	Name     string
	DataType DataType
	// IntValue    int
	// FloatValue  float32
	// BoolValue   bool
	// StringValue string
	Value        string
	IsHistorical bool
	Metadata     Metadata
}

type DataType uint32

const (
	TypeInt8   DataType = 1
	TypeInt16  DataType = 2
	TypeInt32  DataType = 3
	TypeInt64  DataType = 4
	TypeUInt8  DataType = 5
	TypeUInt16 DataType = 6
	TypeUInt32 DataType = 7
	TypeUInt64 DataType = 8
	TypeFloat  DataType = 9
	TypeBool   DataType = 11
	TypeString DataType = 12
	TypeBytes  DataType = 17
)

func (d *DataType) String() string {
	switch *d {
	case TypeInt8:
		return "TypeInt8"
	case TypeInt16:
		return "TypeInt"
	case TypeInt32:
		return "TypeInt32"
	case TypeInt64:
		return "TypeInt64"
	case TypeUInt8:
		return "TypeUInt8"
	case TypeUInt16:
		return "TypeUInt"
	case TypeUInt32:
		return "TypeUInt32"
	case TypeUInt64:
		return "TypeUInt64"
	case TypeFloat:
		return "TypeFloat"
	case TypeBool:
		return "TypeBool"
	case TypeString:
		return "TypeString"
	case TypeBytes:
		return "TypeBytes"
	}

	fmt.Println(int(d.toUint32()))
	return "error"
}

func (d DataType) toUint32() uint32 {
	return uint32(d)
}
