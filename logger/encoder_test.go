package logger

import (
	"math"
	"strconv"
	"testing"
	"time"
)

// ==================== Buffer 接口测试 ====================

// testBuffer 用于测试的 Buffer 实现
type testBuffer struct {
	data []byte
}

func (b *testBuffer) Bytes() []byte {
	return b.data
}

func (b *testBuffer) AppendByte(c byte) {
	b.data = append(b.data, c)
}

func (b *testBuffer) AppendBytes(bs []byte) {
	b.data = append(b.data, bs...)
}

func (b *testBuffer) AppendString(s string) {
	b.data = append(b.data, s...)
}

func (b *testBuffer) AppendInt(i int64) {
	b.data = append(b.data, []byte(intToString(i))...)
}

func (b *testBuffer) AppendFloat(f float64, fmt byte, prec, bitSize int) {
	b.data = append(b.data, []byte(floatToString(f, prec))...)
}

func (b *testBuffer) Reset() {
	b.data = b.data[:0]
}

func intToString(i int64) string {
	if i == 0 {
		return "0"
	}
	if i < 0 {
		return "-" + uintToString(uint64(-i))
	}
	return uintToString(uint64(i))
}

func uintToString(n uint64) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte(n%10) + '0'
		n /= 10
	}
	return string(buf[i:])
}

func floatToString(f float64, prec int) string {
	if prec == -1 {
		prec = 10
	}
	return strconv.FormatFloat(f, 'f', prec, 64)
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && index(s, substr) >= 0
}

func index(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// ==================== FieldType 测试 ====================

func TestFieldType_Values(t *testing.T) {
	tests := []struct {
		ft     FieldType
		expect int
		name   string
	}{
		{Int64Type, 0, "Int64Type"},
		{StringType, 1, "StringType"},
		{BoolType, 2, "BoolType"},
		{Float64Type, 3, "Float64Type"},
		{TimeType, 4, "TimeType"},
		{DurationType, 5, "DurationType"},
		{ErrorType, 6, "ErrorType"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.ft) != tt.expect {
				t.Errorf("expected %s = %d, got %d", tt.name, tt.expect, tt.ft)
			}
		})
	}
}

func TestFieldType_Count(t *testing.T) {
	// 验证 FieldType 数量
	expected := 7
	if ErrorType+1 != FieldType(expected) {
		t.Errorf("expected %d FieldTypes, got %d", expected, ErrorType+1)
	}
}

// ==================== Field 结构体测试 ====================

func TestField_String(t *testing.T) {
	f := Field{
		Key:    "name",
		Type:   StringType,
		String: "test value",
	}

	if f.Key != "name" {
		t.Errorf("expected key 'name', got '%s'", f.Key)
	}
	if f.Type != StringType {
		t.Errorf("expected type StringType, got %v", f.Type)
	}
	if f.String != "test value" {
		t.Errorf("expected string 'test value', got '%s'", f.String)
	}
}

func TestField_Int64(t *testing.T) {
	f := Field{
		Key:     "count",
		Type:    Int64Type,
		Integer: 12345,
	}

	if f.Key != "count" {
		t.Errorf("expected key 'count', got '%s'", f.Key)
	}
	if f.Type != Int64Type {
		t.Errorf("expected type Int64Type, got %v", f.Type)
	}
	if f.Integer != 12345 {
		t.Errorf("expected integer 12345, got %d", f.Integer)
	}
}

func TestField_Bool(t *testing.T) {
	f := Field{
		Key:     "enabled",
		Type:    BoolType,
		Integer: 1,
	}

	if f.Type != BoolType {
		t.Errorf("expected type BoolType, got %v", f.Type)
	}
	if f.Integer != 1 {
		t.Errorf("expected integer 1, got %d", f.Integer)
	}
}

func TestField_Float64(t *testing.T) {
	f := Field{
		Key:   "rate",
		Type:  Float64Type,
		Float: math.Float64bits(3.14),
	}

	if f.Type != Float64Type {
		t.Errorf("expected type Float64Type, got %v", f.Type)
	}
	if math.Float64frombits(f.Float) != 3.14 {
		t.Errorf("expected float 3.14, got %f", math.Float64frombits(f.Float))
	}
}

func TestField_Time(t *testing.T) {
	ts := time.Now().UnixNano()
	f := Field{
		Key:     "timestamp",
		Type:    TimeType,
		Integer: ts,
	}

	if f.Type != TimeType {
		t.Errorf("expected type TimeType, got %v", f.Type)
	}
	if f.Integer != ts {
		t.Errorf("expected timestamp %d, got %d", ts, f.Integer)
	}
}

func TestField_Duration(t *testing.T) {
	f := Field{
		Key:     "latency",
		Type:    DurationType,
		Integer: int64(time.Second),
	}

	if f.Type != DurationType {
		t.Errorf("expected type DurationType, got %v", f.Type)
	}
	if f.Integer != int64(time.Second) {
		t.Errorf("expected duration %d, got %d", int64(time.Second), f.Integer)
	}
}

func TestField_Error(t *testing.T) {
	f := Field{
		Type:   ErrorType,
		Key:    "error",
		String: "something went wrong",
	}

	if f.Type != ErrorType {
		t.Errorf("expected type ErrorType, got %v", f.Type)
	}
	if f.String != "something went wrong" {
		t.Errorf("expected error 'something went wrong', got '%s'", f.String)
	}
}

// ==================== Buffer 接口测试 ====================

func TestBuffer_AppendByte(t *testing.T) {
	buf := &testBuffer{}

	buf.AppendByte('a')
	buf.AppendByte('b')
	buf.AppendByte('c')

	result := string(buf.Bytes())
	if result != "abc" {
		t.Errorf("expected 'abc', got '%s'", result)
	}
}

func TestBuffer_AppendBytes(t *testing.T) {
	buf := &testBuffer{}

	buf.AppendBytes([]byte("hello"))
	buf.AppendBytes([]byte(" world"))

	result := string(buf.Bytes())
	if result != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", result)
	}
}

func TestBuffer_AppendString(t *testing.T) {
	buf := &testBuffer{}

	buf.AppendString("test")
	buf.AppendString(" string")

	result := string(buf.Bytes())
	if result != "test string" {
		t.Errorf("expected 'test string', got '%s'", result)
	}
}

func TestBuffer_AppendInt(t *testing.T) {
	buf := &testBuffer{}

	buf.AppendInt(0)
	buf.AppendInt(123)
	buf.AppendInt(-456)

	result := string(buf.Bytes())
	if result != "0123-456" {
		t.Errorf("expected '0123-456', got '%s'", result)
	}
}

func TestBuffer_AppendFloat(t *testing.T) {
	buf := &testBuffer{}

	buf.AppendFloat(3.14, 'f', -1, 64)

	result := string(buf.Bytes())
	if result == "" {
		t.Errorf("expected float output, got empty")
	}
}

// ==================== Encoder 接口测试 ====================

func TestEncoder_Interface(t *testing.T) {
	// 验证 jsonEncoder 实现了 Encoder 接口
	var _ Encoder = (*jsonEncoder)(nil)
}

func TestEncoder_Clone(t *testing.T) {
	enc := NewJSONEncoder()
	enc.AddFields([]Field{
		{Key: "field1", Type: StringType, String: "value1"},
	})

	cloned := enc.Clone()

	// 克隆应该返回 Encoder 接口类型
	if cloned == nil {
		t.Error("clone returned nil")
	}

	// 克隆的 encoder 应该有相同的字段
	buf1 := &testBuffer{}
	now := time.Now()

	cloned.Encode(buf1, "test", "info", now, nil)

	if !contains(string(buf1.Bytes()), "field1") {
		t.Errorf("clone should have field1")
	}
}

func TestEncoder_AddFields(t *testing.T) {
	enc := NewJSONEncoder()

	// 添加多个字段
	fields := []Field{
		{Key: "field1", Type: StringType, String: "value1"},
		{Key: "field2", Type: Int64Type, Integer: 42},
		{Key: "field3", Type: BoolType, Integer: 1},
	}

	enc.AddFields(fields)

	// 验证字段已添加
	buf := &testBuffer{}
	enc.Encode(buf, "test", "info", time.Now(), nil)

	result := string(buf.Bytes())
	if !contains(result, "field1") || !contains(result, "field2") || !contains(result, "field3") {
		t.Errorf("expected all fields, got: %s", result)
	}
}

func TestEncoder_Encode_BasicOutput(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &testBuffer{}

	enc.Encode(buf, "message", "info", time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC), nil)

	result := string(buf.Bytes())

	// 验证 JSON 结构
	if !contains(result, "{") || !contains(result, "}") {
		t.Errorf("expected JSON object, got: %s", result)
	}
	if !contains(result, "level") {
		t.Errorf("expected level field, got: %s", result)
	}
	if !contains(result, "msg") {
		t.Errorf("expected msg field, got: %s", result)
	}
	if !contains(result, "time") {
		t.Errorf("expected time field, got: %s", result)
	}
	if !contains(result, "\n") {
		t.Errorf("expected newline ending, got: %s", result)
	}
}

// ==================== 边界情况测试 ====================

func TestField_ZeroValue(t *testing.T) {
	// 测试 Field 零值
	f := Field{}

	if f.Key != "" {
		t.Errorf("expected empty key, got '%s'", f.Key)
	}
	if f.Type != 0 {
		t.Errorf("expected zero type, got %v", f.Type)
	}
	if f.Integer != 0 {
		t.Errorf("expected zero integer, got %d", f.Integer)
	}
	if f.String != "" {
		t.Errorf("expected empty string, got '%s'", f.String)
	}
}

func TestEncoder_Encode_NilFields(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &testBuffer{}

	// 传入 nil fields 应该不崩溃
	enc.Encode(buf, "test", "info", time.Now(), nil)

	result := string(buf.Bytes())
	if result == "" {
		t.Error("expected output, got empty")
	}
}

func TestEncoder_Encode_EmptyKey(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &testBuffer{}

	fields := []Field{
		{Key: "", Type: StringType, String: "value"},
	}

	enc.Encode(buf, "test", "info", time.Now(), fields)

	result := string(buf.Bytes())
	// 空 key 应该仍然被处理
	if !contains(result, `"":`) {
		t.Errorf("expected empty key field, got: %s", result)
	}
}

func TestEncoder_AddFields_Empty(t *testing.T) {
	enc := NewJSONEncoder()

	// 添加空切片
	enc.AddFields([]Field{})

	buf := &testBuffer{}
	enc.Encode(buf, "test", "info", time.Now(), nil)

	result := string(buf.Bytes())
	if result == "" {
		t.Error("expected output, got empty")
	}
}
