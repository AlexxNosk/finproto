// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v6.31.1
// source: trade_api/grpc/tradeapi/v1/accounts/accounts_service.proto

package accounts_service

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	decimal "google.golang.org/genproto/googleapis/type/decimal"
	interval "google.golang.org/genproto/googleapis/type/interval"
	money "google.golang.org/genproto/googleapis/type/money"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	_ "github.com/alexxnosk/finproto/backend/trade_api/v1/side"
	trade "github.com/alexxnosk/finproto/backend/trade_api/v1/trade"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Запрос получения информации по конкретному аккаунту
type GetAccountRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Идентификатор аккаунта
	AccountId     string `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAccountRequest) Reset() {
	*x = GetAccountRequest{}
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAccountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAccountRequest) ProtoMessage() {}

func (x *GetAccountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAccountRequest.ProtoReflect.Descriptor instead.
func (*GetAccountRequest) Descriptor() ([]byte, []int) {
	return file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDescGZIP(), []int{0}
}

func (x *GetAccountRequest) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

// Информация о конкретном аккаунте
type GetAccountResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Идентификатор аккаунта
	AccountId string `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	// Тип аккаунта
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	// Статус аккаунта
	Status string `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	// Доступные средства плюс стоимость открытых позиций
	Equity *decimal.Decimal `protobuf:"bytes,4,opt,name=equity,proto3" json:"equity,omitempty"`
	// Нереализованная прибыль
	UnrealizedProfit *decimal.Decimal `protobuf:"bytes,5,opt,name=unrealized_profit,json=unrealizedProfit,proto3" json:"unrealized_profit,omitempty"`
	// Позиции. Открытые, плюс теоретические (по неисполненным активным заявкам)
	Positions []*Position `protobuf:"bytes,6,rep,name=positions,proto3" json:"positions,omitempty"`
	// Доступные средства
	Cash          []*money.Money `protobuf:"bytes,7,rep,name=cash,proto3" json:"cash,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetAccountResponse) Reset() {
	*x = GetAccountResponse{}
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetAccountResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetAccountResponse) ProtoMessage() {}

func (x *GetAccountResponse) ProtoReflect() protoreflect.Message {
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetAccountResponse.ProtoReflect.Descriptor instead.
func (*GetAccountResponse) Descriptor() ([]byte, []int) {
	return file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDescGZIP(), []int{1}
}

func (x *GetAccountResponse) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *GetAccountResponse) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *GetAccountResponse) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *GetAccountResponse) GetEquity() *decimal.Decimal {
	if x != nil {
		return x.Equity
	}
	return nil
}

func (x *GetAccountResponse) GetUnrealizedProfit() *decimal.Decimal {
	if x != nil {
		return x.UnrealizedProfit
	}
	return nil
}

func (x *GetAccountResponse) GetPositions() []*Position {
	if x != nil {
		return x.Positions
	}
	return nil
}

func (x *GetAccountResponse) GetCash() []*money.Money {
	if x != nil {
		return x.Cash
	}
	return nil
}

// Запрос получения истории по сделкам
type TradesRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Идентификатор аккаунта
	AccountId string `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	// Лимит количества сделок
	Limit int32 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	// Начало и окончание запрашиваемого периода, Unix epoch time
	Interval      *interval.Interval `protobuf:"bytes,3,opt,name=interval,proto3" json:"interval,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TradesRequest) Reset() {
	*x = TradesRequest{}
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TradesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TradesRequest) ProtoMessage() {}

func (x *TradesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TradesRequest.ProtoReflect.Descriptor instead.
func (*TradesRequest) Descriptor() ([]byte, []int) {
	return file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDescGZIP(), []int{2}
}

func (x *TradesRequest) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *TradesRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *TradesRequest) GetInterval() *interval.Interval {
	if x != nil {
		return x.Interval
	}
	return nil
}

// История по сделкам
type TradesResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Сделки по аккаунту
	Trades        []*trade.AccountTrade `protobuf:"bytes,1,rep,name=trades,proto3" json:"trades,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TradesResponse) Reset() {
	*x = TradesResponse{}
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TradesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TradesResponse) ProtoMessage() {}

func (x *TradesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TradesResponse.ProtoReflect.Descriptor instead.
func (*TradesResponse) Descriptor() ([]byte, []int) {
	return file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDescGZIP(), []int{3}
}

func (x *TradesResponse) GetTrades() []*trade.AccountTrade {
	if x != nil {
		return x.Trades
	}
	return nil
}

// Запрос получения списка транзакций
type TransactionsRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Идентификатор аккаунта
	AccountId string `protobuf:"bytes,1,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	// Лимит количества транзакций
	Limit int32 `protobuf:"varint,2,opt,name=limit,proto3" json:"limit,omitempty"`
	// Начало и окончание запрашиваемого периода, Unix epoch time
	Interval      *interval.Interval `protobuf:"bytes,3,opt,name=interval,proto3" json:"interval,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TransactionsRequest) Reset() {
	*x = TransactionsRequest{}
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransactionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionsRequest) ProtoMessage() {}

func (x *TransactionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionsRequest.ProtoReflect.Descriptor instead.
func (*TransactionsRequest) Descriptor() ([]byte, []int) {
	return file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDescGZIP(), []int{4}
}

func (x *TransactionsRequest) GetAccountId() string {
	if x != nil {
		return x.AccountId
	}
	return ""
}

func (x *TransactionsRequest) GetLimit() int32 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *TransactionsRequest) GetInterval() *interval.Interval {
	if x != nil {
		return x.Interval
	}
	return nil
}

// Список транзакций
type TransactionsResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Транзакции по аккаунту
	Transactions  []*Transaction `protobuf:"bytes,1,rep,name=transactions,proto3" json:"transactions,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TransactionsResponse) Reset() {
	*x = TransactionsResponse{}
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TransactionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TransactionsResponse) ProtoMessage() {}

func (x *TransactionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TransactionsResponse.ProtoReflect.Descriptor instead.
func (*TransactionsResponse) Descriptor() ([]byte, []int) {
	return file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDescGZIP(), []int{5}
}

func (x *TransactionsResponse) GetTransactions() []*Transaction {
	if x != nil {
		return x.Transactions
	}
	return nil
}

// Информация о позиции
type Position struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Символ инструмента
	Symbol string `protobuf:"bytes,1,opt,name=symbol,proto3" json:"symbol,omitempty"`
	// Количество в шт., значение со знаком определяющее (long-short)
	Quantity *decimal.Decimal `protobuf:"bytes,2,opt,name=quantity,proto3" json:"quantity,omitempty"`
	// Средняя цена
	AveragePrice *decimal.Decimal `protobuf:"bytes,3,opt,name=average_price,json=averagePrice,proto3" json:"average_price,omitempty"`
	// Текущая цена
	CurrentPrice  *decimal.Decimal `protobuf:"bytes,4,opt,name=current_price,json=currentPrice,proto3" json:"current_price,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Position) Reset() {
	*x = Position{}
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Position) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Position) ProtoMessage() {}

func (x *Position) ProtoReflect() protoreflect.Message {
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Position.ProtoReflect.Descriptor instead.
func (*Position) Descriptor() ([]byte, []int) {
	return file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDescGZIP(), []int{6}
}

func (x *Position) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *Position) GetQuantity() *decimal.Decimal {
	if x != nil {
		return x.Quantity
	}
	return nil
}

func (x *Position) GetAveragePrice() *decimal.Decimal {
	if x != nil {
		return x.AveragePrice
	}
	return nil
}

func (x *Position) GetCurrentPrice() *decimal.Decimal {
	if x != nil {
		return x.CurrentPrice
	}
	return nil
}

// Информация о транзакции
type Transaction struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Идентификатор транзакции
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Тип транзакции из TransactionCategory
	Category string `protobuf:"bytes,2,opt,name=category,proto3" json:"category,omitempty"`
	// Метка времени
	Timestamp *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// Символ инструмента
	Symbol string `protobuf:"bytes,5,opt,name=symbol,proto3" json:"symbol,omitempty"`
	// Изменение в деньгах
	Change *money.Money `protobuf:"bytes,6,opt,name=change,proto3" json:"change,omitempty"`
	// Информация о сделке
	Trade         *Transaction_Trade `protobuf:"bytes,7,opt,name=trade,proto3" json:"trade,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Transaction) Reset() {
	*x = Transaction{}
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Transaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction) ProtoMessage() {}

func (x *Transaction) ProtoReflect() protoreflect.Message {
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction.ProtoReflect.Descriptor instead.
func (*Transaction) Descriptor() ([]byte, []int) {
	return file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDescGZIP(), []int{7}
}

func (x *Transaction) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Transaction) GetCategory() string {
	if x != nil {
		return x.Category
	}
	return ""
}

func (x *Transaction) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

func (x *Transaction) GetSymbol() string {
	if x != nil {
		return x.Symbol
	}
	return ""
}

func (x *Transaction) GetChange() *money.Money {
	if x != nil {
		return x.Change
	}
	return nil
}

func (x *Transaction) GetTrade() *Transaction_Trade {
	if x != nil {
		return x.Trade
	}
	return nil
}

// Объект заполняется для торговых типов транзакций
type Transaction_Trade struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Количество в шт.
	Size *decimal.Decimal `protobuf:"bytes,1,opt,name=size,proto3" json:"size,omitempty"`
	// Цена сделки за штуку. Цена исполнения/Размер премии по опциону. Это цена заключения, значение берется из сделки.
	Price *decimal.Decimal `protobuf:"bytes,2,opt,name=price,proto3" json:"price,omitempty"`
	// НКД. Заполнено если в сделке есть НКД
	AccruedInterest *decimal.Decimal `protobuf:"bytes,3,opt,name=accrued_interest,json=accruedInterest,proto3" json:"accrued_interest,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *Transaction_Trade) Reset() {
	*x = Transaction_Trade{}
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Transaction_Trade) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transaction_Trade) ProtoMessage() {}

func (x *Transaction_Trade) ProtoReflect() protoreflect.Message {
	mi := &file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transaction_Trade.ProtoReflect.Descriptor instead.
func (*Transaction_Trade) Descriptor() ([]byte, []int) {
	return file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDescGZIP(), []int{7, 0}
}

func (x *Transaction_Trade) GetSize() *decimal.Decimal {
	if x != nil {
		return x.Size
	}
	return nil
}

func (x *Transaction_Trade) GetPrice() *decimal.Decimal {
	if x != nil {
		return x.Price
	}
	return nil
}

func (x *Transaction_Trade) GetAccruedInterest() *decimal.Decimal {
	if x != nil {
		return x.AccruedInterest
	}
	return nil
}

var File_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto protoreflect.FileDescriptor

const file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDesc = "" +
	"\n" +
	":trade_api/grpc/tradeapi/v1/accounts/accounts_service.proto\x12\x19grpc.tradeapi.v1.accounts\x1a\x1cgoogle/api/annotations.proto\x1a\x1fgoogle/protobuf/timestamp.proto\x1a\x19google/type/decimal.proto\x1a\x1agoogle/type/interval.proto\x1a\x17google/type/money.proto\x1a\x1bgrpc/tradeapi/v1/side.proto\x1a\x1cgrpc/tradeapi/v1/trade.proto\"2\n" +
	"\x11GetAccountRequest\x12\x1d\n" +
	"\n" +
	"account_id\x18\x01 \x01(\tR\taccountId\"\xbb\x02\n" +
	"\x12GetAccountResponse\x12\x1d\n" +
	"\n" +
	"account_id\x18\x01 \x01(\tR\taccountId\x12\x12\n" +
	"\x04type\x18\x02 \x01(\tR\x04type\x12\x16\n" +
	"\x06status\x18\x03 \x01(\tR\x06status\x12,\n" +
	"\x06equity\x18\x04 \x01(\v2\x14.google.type.DecimalR\x06equity\x12A\n" +
	"\x11unrealized_profit\x18\x05 \x01(\v2\x14.google.type.DecimalR\x10unrealizedProfit\x12A\n" +
	"\tpositions\x18\x06 \x03(\v2#.grpc.tradeapi.v1.accounts.PositionR\tpositions\x12&\n" +
	"\x04cash\x18\a \x03(\v2\x12.google.type.MoneyR\x04cash\"w\n" +
	"\rTradesRequest\x12\x1d\n" +
	"\n" +
	"account_id\x18\x01 \x01(\tR\taccountId\x12\x14\n" +
	"\x05limit\x18\x02 \x01(\x05R\x05limit\x121\n" +
	"\binterval\x18\x03 \x01(\v2\x15.google.type.IntervalR\binterval\"H\n" +
	"\x0eTradesResponse\x126\n" +
	"\x06trades\x18\x01 \x03(\v2\x1e.grpc.tradeapi.v1.AccountTradeR\x06trades\"}\n" +
	"\x13TransactionsRequest\x12\x1d\n" +
	"\n" +
	"account_id\x18\x01 \x01(\tR\taccountId\x12\x14\n" +
	"\x05limit\x18\x02 \x01(\x05R\x05limit\x121\n" +
	"\binterval\x18\x03 \x01(\v2\x15.google.type.IntervalR\binterval\"b\n" +
	"\x14TransactionsResponse\x12J\n" +
	"\ftransactions\x18\x01 \x03(\v2&.grpc.tradeapi.v1.accounts.TransactionR\ftransactions\"\xca\x01\n" +
	"\bPosition\x12\x16\n" +
	"\x06symbol\x18\x01 \x01(\tR\x06symbol\x120\n" +
	"\bquantity\x18\x02 \x01(\v2\x14.google.type.DecimalR\bquantity\x129\n" +
	"\raverage_price\x18\x03 \x01(\v2\x14.google.type.DecimalR\faveragePrice\x129\n" +
	"\rcurrent_price\x18\x04 \x01(\v2\x14.google.type.DecimalR\fcurrentPrice\"\x9c\x03\n" +
	"\vTransaction\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x1a\n" +
	"\bcategory\x18\x02 \x01(\tR\bcategory\x128\n" +
	"\ttimestamp\x18\x04 \x01(\v2\x1a.google.protobuf.TimestampR\ttimestamp\x12\x16\n" +
	"\x06symbol\x18\x05 \x01(\tR\x06symbol\x12*\n" +
	"\x06change\x18\x06 \x01(\v2\x12.google.type.MoneyR\x06change\x12B\n" +
	"\x05trade\x18\a \x01(\v2,.grpc.tradeapi.v1.accounts.Transaction.TradeR\x05trade\x1a\x9e\x01\n" +
	"\x05Trade\x12(\n" +
	"\x04size\x18\x01 \x01(\v2\x14.google.type.DecimalR\x04size\x12*\n" +
	"\x05price\x18\x02 \x01(\v2\x14.google.type.DecimalR\x05price\x12?\n" +
	"\x10accrued_interest\x18\x03 \x01(\v2\x14.google.type.DecimalR\x0faccruedInterest2\xcc\x03\n" +
	"\x0fAccountsService\x12\x8c\x01\n" +
	"\n" +
	"GetAccount\x12,.grpc.tradeapi.v1.accounts.GetAccountRequest\x1a-.grpc.tradeapi.v1.accounts.GetAccountResponse\"!\x82\xd3\xe4\x93\x02\x1b\x12\x19/v1/accounts/{account_id}\x12\x87\x01\n" +
	"\x06Trades\x12(.grpc.tradeapi.v1.accounts.TradesRequest\x1a).grpc.tradeapi.v1.accounts.TradesResponse\"(\x82\xd3\xe4\x93\x02\"\x12 /v1/accounts/{account_id}/trades\x12\x9f\x01\n" +
	"\fTransactions\x12..grpc.tradeapi.v1.accounts.TransactionsRequest\x1a/.grpc.tradeapi.v1.accounts.TransactionsResponse\".\x82\xd3\xe4\x93\x02(\x12&/v1/accounts/{account_id}/transactionsB*P\x01Z&trade_api/v1/accounts/accounts_serviceb\x06proto3"

var (
	file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDescOnce sync.Once
	file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDescData []byte
)

func file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDescGZIP() []byte {
	file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDescOnce.Do(func() {
		file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDesc), len(file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDesc)))
	})
	return file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDescData
}

var file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_goTypes = []any{
	(*GetAccountRequest)(nil),     // 0: grpc.tradeapi.v1.accounts.GetAccountRequest
	(*GetAccountResponse)(nil),    // 1: grpc.tradeapi.v1.accounts.GetAccountResponse
	(*TradesRequest)(nil),         // 2: grpc.tradeapi.v1.accounts.TradesRequest
	(*TradesResponse)(nil),        // 3: grpc.tradeapi.v1.accounts.TradesResponse
	(*TransactionsRequest)(nil),   // 4: grpc.tradeapi.v1.accounts.TransactionsRequest
	(*TransactionsResponse)(nil),  // 5: grpc.tradeapi.v1.accounts.TransactionsResponse
	(*Position)(nil),              // 6: grpc.tradeapi.v1.accounts.Position
	(*Transaction)(nil),           // 7: grpc.tradeapi.v1.accounts.Transaction
	(*Transaction_Trade)(nil),     // 8: grpc.tradeapi.v1.accounts.Transaction.Trade
	(*decimal.Decimal)(nil),       // 9: google.type.Decimal
	(*money.Money)(nil),           // 10: google.type.Money
	(*interval.Interval)(nil),     // 11: google.type.Interval
	(*trade.AccountTrade)(nil),    // 12: grpc.tradeapi.v1.AccountTrade
	(*timestamppb.Timestamp)(nil), // 13: google.protobuf.Timestamp
}
var file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_depIdxs = []int32{
	9,  // 0: grpc.tradeapi.v1.accounts.GetAccountResponse.equity:type_name -> google.type.Decimal
	9,  // 1: grpc.tradeapi.v1.accounts.GetAccountResponse.unrealized_profit:type_name -> google.type.Decimal
	6,  // 2: grpc.tradeapi.v1.accounts.GetAccountResponse.positions:type_name -> grpc.tradeapi.v1.accounts.Position
	10, // 3: grpc.tradeapi.v1.accounts.GetAccountResponse.cash:type_name -> google.type.Money
	11, // 4: grpc.tradeapi.v1.accounts.TradesRequest.interval:type_name -> google.type.Interval
	12, // 5: grpc.tradeapi.v1.accounts.TradesResponse.trades:type_name -> grpc.tradeapi.v1.AccountTrade
	11, // 6: grpc.tradeapi.v1.accounts.TransactionsRequest.interval:type_name -> google.type.Interval
	7,  // 7: grpc.tradeapi.v1.accounts.TransactionsResponse.transactions:type_name -> grpc.tradeapi.v1.accounts.Transaction
	9,  // 8: grpc.tradeapi.v1.accounts.Position.quantity:type_name -> google.type.Decimal
	9,  // 9: grpc.tradeapi.v1.accounts.Position.average_price:type_name -> google.type.Decimal
	9,  // 10: grpc.tradeapi.v1.accounts.Position.current_price:type_name -> google.type.Decimal
	13, // 11: grpc.tradeapi.v1.accounts.Transaction.timestamp:type_name -> google.protobuf.Timestamp
	10, // 12: grpc.tradeapi.v1.accounts.Transaction.change:type_name -> google.type.Money
	8,  // 13: grpc.tradeapi.v1.accounts.Transaction.trade:type_name -> grpc.tradeapi.v1.accounts.Transaction.Trade
	9,  // 14: grpc.tradeapi.v1.accounts.Transaction.Trade.size:type_name -> google.type.Decimal
	9,  // 15: grpc.tradeapi.v1.accounts.Transaction.Trade.price:type_name -> google.type.Decimal
	9,  // 16: grpc.tradeapi.v1.accounts.Transaction.Trade.accrued_interest:type_name -> google.type.Decimal
	0,  // 17: grpc.tradeapi.v1.accounts.AccountsService.GetAccount:input_type -> grpc.tradeapi.v1.accounts.GetAccountRequest
	2,  // 18: grpc.tradeapi.v1.accounts.AccountsService.Trades:input_type -> grpc.tradeapi.v1.accounts.TradesRequest
	4,  // 19: grpc.tradeapi.v1.accounts.AccountsService.Transactions:input_type -> grpc.tradeapi.v1.accounts.TransactionsRequest
	1,  // 20: grpc.tradeapi.v1.accounts.AccountsService.GetAccount:output_type -> grpc.tradeapi.v1.accounts.GetAccountResponse
	3,  // 21: grpc.tradeapi.v1.accounts.AccountsService.Trades:output_type -> grpc.tradeapi.v1.accounts.TradesResponse
	5,  // 22: grpc.tradeapi.v1.accounts.AccountsService.Transactions:output_type -> grpc.tradeapi.v1.accounts.TransactionsResponse
	20, // [20:23] is the sub-list for method output_type
	17, // [17:20] is the sub-list for method input_type
	17, // [17:17] is the sub-list for extension type_name
	17, // [17:17] is the sub-list for extension extendee
	0,  // [0:17] is the sub-list for field type_name
}

func init() { file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_init() }
func file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_init() {
	if File_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDesc), len(file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_goTypes,
		DependencyIndexes: file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_depIdxs,
		MessageInfos:      file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_msgTypes,
	}.Build()
	File_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto = out.File
	file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_goTypes = nil
	file_trade_api_grpc_tradeapi_v1_accounts_accounts_service_proto_depIdxs = nil
}
