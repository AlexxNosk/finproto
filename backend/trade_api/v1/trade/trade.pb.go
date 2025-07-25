// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: trade_api/grpc/tradeapi/v1/trade.proto

package trade

import (
	decimal "google.golang.org/genproto/googleapis/type/decimal"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	side "github.com/alexxnosk/finproto/backend/trade_api/v1/side"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Информация о сделке
type AccountTrade struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Идентификатор сделки
	TradeId string `protobuf:"bytes,1,opt,name=trade_id,json=tradeId,proto3" json:"trade_id,omitempty"`
	// Символ инструмента
	Symbol string `protobuf:"bytes,2,opt,name=symbol,proto3" json:"symbol,omitempty"`
	// Цена исполнения
	Price *decimal.Decimal `protobuf:"bytes,3,opt,name=price,proto3" json:"price,omitempty"`
	// Размер в шт.
	Size *decimal.Decimal `protobuf:"bytes,4,opt,name=size,proto3" json:"size,omitempty"`
	// Сторона сделки (long или short)
	Side side.Side `protobuf:"varint,5,opt,name=side,proto3,enum=grpc.tradeapi.v1.Side" json:"side,omitempty"`
	// Метка времени
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Идентификатор заявки
	OrderId       string `protobuf:"bytes,7,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AccountTrade) Reset() {
	*x = AccountTrade{}
	mi := &file_trade_api_grpc_tradeapi_v1_trade_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AccountTrade) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccountTrade) ProtoMessage() {}

func (x *AccountTrade) ProtoReflect() protoreflect.Message {
	mi := &file_trade_api_grpc_tradeapi_v1_trade_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccountTrade.ProtoReflect.Descriptor instead.
func (*AccountTrade) Descriptor() ([]byte, []int) {
	return file_trade_api_grpc_tradeapi_v1_trade_proto_rawDescGZIP(), []int{0}
}

func (x *AccountTrade) GetTradeId() string {
	if x != nil {
		return x.TradeId
	}
	return ""
}

func (x *AccountTrade) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *AccountTrade) GetPrice() *decimal.Decimal {
	if x != nil {
		return x.Price
	}
	return nil
}

func (x *AccountTrade) GetSize() *decimal.Decimal {
	if x != nil {
		return x.Size
	}
	return nil
}

func (x *AccountTrade) GetSide() side.Side {
	if x != nil {
		return x.Side
	}
	return side.Side(0)
}

func (x *AccountTrade) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *AccountTrade) GetOrderId() string {
	if x != nil {
		return x.OrderId
	}
	return ""
}

var File_trade_api_grpc_tradeapi_v1_trade_proto protoreflect.FileDescriptor

const file_trade_api_grpc_tradeapi_v1_trade_proto_rawDesc = "" +
	"\n" +
	"&trade_api/grpc/tradeapi/v1/trade.proto\x12\x10grpc.tradeapi.v1\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x19google/type/decimal.proto\x1a\x1bgrpc/tradeapi/v1/side.proto\"\x98\x02\n" +
	"\fAccountTrade\x12\x19\n" +
	"\btrade_id\x18\x01 \x01(\tR\atradeId\x12\x16\n" +
	"\x06symbol\x18\x02 \x01(\tR\x06symbol\x12*\n" +
	"\x05price\x18\x03 \x01(\v2\x14.google.type.DecimalR\x05price\x12(\n" +
	"\x04size\x18\x04 \x01(\v2\x14.google.type.DecimalR\x04size\x12*\n" +
	"\x04side\x18\x05 \x01(\x0e2\x16.grpc.tradeapi.v1.SideR\x04side\x128\n" +
	"\ttimestamp\x18\x06 \x01(\v2\x1a.google.protobuf.TimestampR\ttimestamp\x12\x19\n" +
	"\border_id\x18\a \x01(\tR\aorderIdB\x16P\x01Z\x12trade_api/v1/tradeb\x06proto3"

var (
	file_trade_api_grpc_tradeapi_v1_trade_proto_rawDescOnce sync.Once
	file_trade_api_grpc_tradeapi_v1_trade_proto_rawDescData []byte
)

func file_trade_api_grpc_tradeapi_v1_trade_proto_rawDescGZIP() []byte {
	file_trade_api_grpc_tradeapi_v1_trade_proto_rawDescOnce.Do(func() {
		file_trade_api_grpc_tradeapi_v1_trade_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_trade_api_grpc_tradeapi_v1_trade_proto_rawDesc), len(file_trade_api_grpc_tradeapi_v1_trade_proto_rawDesc)))
	})
	return file_trade_api_grpc_tradeapi_v1_trade_proto_rawDescData
}

var file_trade_api_grpc_tradeapi_v1_trade_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_trade_api_grpc_tradeapi_v1_trade_proto_goTypes = []any{
	(*AccountTrade)(nil),          // 0: grpc.tradeapi.v1.AccountTrade
	(*decimal.Decimal)(nil),       // 1: google.type.Decimal
	(side.Side)(0),                // 2: grpc.tradeapi.v1.Side
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_trade_api_grpc_tradeapi_v1_trade_proto_depIdxs = []int32{
	1, // 0: grpc.tradeapi.v1.AccountTrade.price:type_name -> google.type.Decimal
	1, // 1: grpc.tradeapi.v1.AccountTrade.size:type_name -> google.type.Decimal
	2, // 2: grpc.tradeapi.v1.AccountTrade.side:type_name -> grpc.tradeapi.v1.Side
	3, // 3: grpc.tradeapi.v1.AccountTrade.timestamp:type_name -> google.protobuf.Timestamp
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_trade_api_grpc_tradeapi_v1_trade_proto_init() }
func file_trade_api_grpc_tradeapi_v1_trade_proto_init() {
	if File_trade_api_grpc_tradeapi_v1_trade_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_trade_api_grpc_tradeapi_v1_trade_proto_rawDesc), len(file_trade_api_grpc_tradeapi_v1_trade_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_trade_api_grpc_tradeapi_v1_trade_proto_goTypes,
		DependencyIndexes: file_trade_api_grpc_tradeapi_v1_trade_proto_depIdxs,
		MessageInfos:      file_trade_api_grpc_tradeapi_v1_trade_proto_msgTypes,
	}.Build()
	File_trade_api_grpc_tradeapi_v1_trade_proto = out.File
	file_trade_api_grpc_tradeapi_v1_trade_proto_goTypes = nil
	file_trade_api_grpc_tradeapi_v1_trade_proto_depIdxs = nil
}
