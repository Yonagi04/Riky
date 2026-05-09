package logger

import (
	"math"
	"strconv"
	"strings"
	"testing"
	"time"
)

// mockBuffer 用于测试的 Buffer 实现
type mockBuffer struct {
	sb strings.Builder
}

func (b *mockBuffer) Bytes() []byte {
	return []byte(b.sb.String())
}

func (b *mockBuffer) AppendByte(c byte) {
	b.sb.WriteByte(c)
}

func (b *mockBuffer) AppendBytes(bs []byte) {
	b.sb.Write(bs)
}

func (b *mockBuffer) AppendString(s string) {
	b.sb.WriteString(s)
}

func (b *mockBuffer) AppendInt(i int64) {
	b.sb.WriteString(int64ToString(i))
}

func (b *mockBuffer) AppendFloat(f float64, fmt byte, prec, bitSize int) {
	b.sb.WriteString(float64ToString(f, fmt, prec))
}

func (b *mockBuffer) Reset() {
	b.sb.Reset()
}

func int64ToString(i int64) string {
	if i == 0 {
		return "0"
	}
	if i < 0 {
		return "-" + uitoa(uint64(-i))
	}
	return uitoa(uint64(i))
}

func uitoa(n uint64) string {
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

func float64ToString(f float64, fmt byte, prec int) string {
	if prec == -1 {
		prec = -1
	}
	var sb strings.Builder
	if math.IsInf(f, 0) {
		sb.WriteString("Inf")
		return sb.String()
	}
	if math.IsNaN(f) {
		sb.WriteString("NaN")
		return sb.String()
	}
	abs := math.Abs(f)
	if abs == 0 {
		sb.WriteString("0")
		if prec > 0 {
			sb.WriteString(".")
			sb.WriteString(strings.Repeat("0", prec))
		}
		return sb.String()
	}
	return strconv.FormatFloat(f, fmt, prec, 64)
}

// ==================== 正常场景测试 ====================

func TestJSONEncoder_Encode_Basic(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &mockBuffer{}

	enc.Encode(buf, "test message", "info", time.Date(2024, 1, 15, 10, 30, 45, 0, time.UTC), nil)

	result := string(buf.Bytes())
	if !strings.Contains(result, `"level":"info"`) {
		t.Errorf("expected level field, got: %s", result)
	}
	if !strings.Contains(result, `"msg":"test message"`) {
		t.Errorf("expected msg field, got: %s", result)
	}
	if !strings.Contains(result, `"time":"2024-01-15 10:30:45"`) {
		t.Errorf("expected time field, got: %s", result)
	}
	if !strings.HasSuffix(result, "}\n") {
		t.Errorf("expected to end with }\\n, got: %s", result)
	}
}

func TestJSONEncoder_Encode_WithFields(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &mockBuffer{}

	fields := []Field{
		{Key: "user_id", Type: Int64Type, Integer: 12345},
		{Key: "action", Type: StringType, String: "login"},
		{Key: "success", Type: BoolType, Integer: 1},
	}

	enc.Encode(buf, "user logged in", "info", time.Now(), fields)

	result := string(buf.Bytes())
	if !strings.Contains(result, `"user_id":12345`) {
		t.Errorf("expected user_id field, got: %s", result)
	}
	if !strings.Contains(result, `"action":"login"`) {
		t.Errorf("expected action field, got: %s", result)
	}
	if !strings.Contains(result, `"success":true`) {
		t.Errorf("expected success field, got: %s", result)
	}
}

func TestJSONEncoder_AddFields(t *testing.T) {
	enc := NewJSONEncoder()

	fields := []Field{
		{Key: "name", Type: StringType, String: "test"},
		{Key: "count", Type: Int64Type, Integer: 42},
		{Key: "rate", Type: Float64Type, Float: math.Float64bits(3.14)},
		{Key: "enabled", Type: BoolType, Integer: 1},
	}

	enc.AddFields(fields)

	// 验证 AddFields 持久化字段会在 Encode 时出现
	buf := &mockBuffer{}
	enc.Encode(buf, "test", "debug", time.Now(), nil)

	result := string(buf.Bytes())
	if !strings.Contains(result, `"name":"test"`) {
		t.Errorf("expected name field, got: %s", result)
	}
	if !strings.Contains(result, `"count":42`) {
		t.Errorf("expected count field, got: %s", result)
	}
}

func TestJSONEncoder_Clone(t *testing.T) {
	enc := NewJSONEncoder()
	enc.AddFields([]Field{{Key: "original", Type: StringType, String: "value"}})

	cloned := enc.Clone().(*jsonEncoder)

	// 修改克隆的 encoder
	cloned.AddFields([]Field{{Key: "cloned", Type: StringType, String: "new"}})

	// 原始 encoder 不应受影响
	buf1 := &mockBuffer{}
	enc.Encode(buf1, "test", "info", time.Now(), nil)
	result1 := string(buf1.Bytes())

	if strings.Contains(result1, `"cloned"`) {
		t.Errorf("original encoder should not have cloned field, got: %s", result1)
	}
	if !strings.Contains(result1, `"original"`) {
		t.Errorf("original encoder should have original field, got: %s", result1)
	}
}

// ==================== 异常场景测试 ====================

func TestJSONEncoder_Encode_EmptyMessage(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &mockBuffer{}

	enc.Encode(buf, "", "info", time.Now(), nil)

	result := string(buf.Bytes())
	if !strings.Contains(result, `"msg":""`) {
		t.Errorf("expected empty msg field, got: %s", result)
	}
}

func TestJSONEncoder_Encode_EmptyLevel(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &mockBuffer{}

	enc.Encode(buf, "test", "", time.Now(), nil)

	result := string(buf.Bytes())
	if !strings.Contains(result, `"level":""`) {
		t.Errorf("expected empty level field, got: %s", result)
	}
}

func TestJSONEncoder_Encode_ZeroTime(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &mockBuffer{}

	enc.Encode(buf, "test", "info", time.Time{}, nil)

	result := string(buf.Bytes())
	// time.Time{} 是 0001-01-01 00:00:00，但实现没有补零
	if !strings.Contains(result, `"time":"1-01-01`) {
		t.Errorf("expected zero time field, got: %s", result)
	}
}

func TestJSONEncoder_AddFields_Empty(t *testing.T) {
	enc := NewJSONEncoder()
	enc.AddFields([]Field{})

	buf := &mockBuffer{}
	enc.Encode(buf, "test", "info", time.Now(), nil)

	result := string(buf.Bytes())
	// 没有任何字段时，输出应该只有基础字段
	if !strings.Contains(result, `"level"`) {
		t.Errorf("should have basic fields, got: %s", result)
	}
}

func TestJSONEncoder_Encode_SpecialCharactersInMessage(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &mockBuffer{}

	enc.Encode(buf, `test "quoted" \ backslash`, "info", time.Now(), nil)

	result := string(buf.Bytes())
	if !strings.Contains(result, `test \"quoted\" \\ backslash`) {
		t.Errorf("expected escaped special characters, got: %s", result)
	}
}

func TestJSONEncoder_Encode_NewlineInMessage(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &mockBuffer{}

	enc.Encode(buf, "line1\nline2", "info", time.Now(), nil)

	result := string(buf.Bytes())
	if !strings.Contains(result, `line1\nline2`) {
		t.Errorf("expected escaped newline, got: %s", result)
	}
}

func TestJSONEncoder_Encode_TabInMessage(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &mockBuffer{}

	enc.Encode(buf, "col1\tcol2", "info", time.Now(), nil)

	result := string(buf.Bytes())
	if !strings.Contains(result, `col1\tcol2`) {
		t.Errorf("expected escaped tab, got: %s", result)
	}
}

// ==================== 边界场景测试 ====================

func TestJSONEncoder_Encode_MaxInt64(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &mockBuffer{}

	fields := []Field{
		{Key: "max_int", Type: Int64Type, Integer: math.MaxInt64},
	}

	enc.Encode(buf, "test", "info", time.Now(), fields)

	result := string(buf.Bytes())
	if !strings.Contains(result, `"max_int":9223372036854775807`) {
		t.Errorf("expected MaxInt64, got: %s", result)
	}
}

func TestJSONEncoder_Encode_MinInt64(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &mockBuffer{}

	fields := []Field{
		{Key: "min_int", Type: Int64Type, Integer: math.MinInt64},
	}

	enc.Encode(buf, "test", "info", time.Now(), fields)

	result := string(buf.Bytes())
	if !strings.Contains(result, `"min_int":-9223372036854775808`) {
		t.Errorf("expected MinInt64, got: %s", result)
	}
}

func TestJSONEncoder_Encode_Float64_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{"max_float", math.MaxFloat64},
		{"smallest_nonzero", math.SmallestNonzeroFloat64},
		{"positive_zero", 0.0},
		{"negative_zero", math.Copysign(0, -1)},
		{"inf", math.Inf(1)},
		{"neg_inf", math.Inf(-1)},
		{"nan", math.NaN()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := NewJSONEncoder()
			buf := &mockBuffer{}

			fields := []Field{
				{Key: "value", Type: Float64Type, Float: math.Float64bits(tt.value)},
			}

			enc.Encode(buf, "test", "info", time.Now(), fields)

			result := string(buf.Bytes())
			if !strings.Contains(result, `"value":`) {
				t.Errorf("expected value field, got: %s", result)
			}
		})
	}
}

func TestJSONEncoder_Encode_AllFieldTypes(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &mockBuffer{}

	fields := []Field{
		{Key: "str_field", Type: StringType, String: "hello"},
		{Key: "int_field", Type: Int64Type, Integer: 100},
		{Key: "float_field", Type: Float64Type, Float: math.Float64bits(1.5)},
		{Key: "bool_true", Type: BoolType, Integer: 1},
		{Key: "bool_false", Type: BoolType, Integer: 0},
		{Key: "time_field", Type: TimeType, Integer: time.Now().UnixNano()},
		{Key: "duration_field", Type: DurationType, Integer: int64(time.Second)},
		{Key: "error_field", Type: ErrorType, String: "error message"},
	}

	enc.Encode(buf, "test", "info", time.Now(), fields)

	result := string(buf.Bytes())
	if !strings.Contains(result, `"str_field":"hello"`) {
		t.Errorf("missing string field, got: %s", result)
	}
	if !strings.Contains(result, `"int_field":100`) {
		t.Errorf("missing int field, got: %s", result)
	}
	if !strings.Contains(result, `"float_field":1.5`) {
		t.Errorf("missing float field, got: %s", result)
	}
	if !strings.Contains(result, `"bool_true":true`) {
		t.Errorf("missing bool true field, got: %s", result)
	}
	if !strings.Contains(result, `"bool_false":false`) {
		t.Errorf("missing bool false field, got: %s", result)
	}
	if !strings.Contains(result, `"time_field":`) {
		t.Errorf("missing time field, got: %s", result)
	}
	if !strings.Contains(result, `"duration_field":`) {
		t.Errorf("missing duration field, got: %s", result)
	}
	if !strings.Contains(result, `"error_field":"error message"`) {
		t.Errorf("missing error field, got: %s", result)
	}
}

func TestJSONEncoder_Encode_Time_LeadingZeros(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "single digit month",
			time:     time.Date(2024, 1, 5, 3, 4, 5, 0, time.UTC),
			expected: `"time":"2024-01-05 03:04:05"`,
		},
		{
			name:     "single digit day",
			time:     time.Date(2024, 12, 3, 9, 8, 7, 0, time.UTC),
			expected: `"time":"2024-12-03 09:08:07"`,
		},
		{
			name:     "all single digits",
			time:     time.Date(2024, 1, 2, 1, 2, 3, 0, time.UTC),
			expected: `"time":"2024-01-02 01:02:03"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc := NewJSONEncoder()
			buf := &mockBuffer{}

			enc.Encode(buf, "test", "info", tt.time, nil)

			result := string(buf.Bytes())
			if !strings.Contains(result, tt.expected) {
				t.Errorf("expected %s, got: %s", tt.expected, result)
			}
		})
	}
}

func TestJSONEncoder_AddFields_PersistsAcrossEncode(t *testing.T) {
	enc := NewJSONEncoder()
	enc.AddFields([]Field{{Key: "persistent", Type: StringType, String: "value"}})

	// 第一次 Encode
	buf1 := &mockBuffer{}
	enc.Encode(buf1, "first", "info", time.Now(), nil)
	result1 := string(buf1.Bytes())

	if !strings.Contains(result1, `"persistent":"value"`) {
		t.Errorf("first encode missing persistent field, got: %s", result1)
	}

	// 第二次 Encode - 持久化字段应该仍然存在
	buf2 := &mockBuffer{}
	enc.Encode(buf2, "second", "debug", time.Now(), nil)
	result2 := string(buf2.Bytes())

	if !strings.Contains(result2, `"persistent":"value"`) {
		t.Errorf("second encode missing persistent field, got: %s", result2)
	}
}

func TestJSONEncoder_Encode_EmptyFields(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &mockBuffer{}

	enc.Encode(buf, "test", "info", time.Now(), []Field{})

	result := string(buf.Bytes())
	// 基础字段后不应该有多余的逗号
	if strings.Contains(result, `,""`) {
		t.Errorf("unexpected empty field comma, got: %s", result)
	}
}

func TestJSONEncoder_ControlCharacters(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &mockBuffer{}

	// 测试各种控制字符
	testCases := []struct {
		char     byte
		expected string
	}{
		{0x01, `\u0001`},
		{0x1F, `\u001f`},
		{'"', `\"`},
		{'\\', `\\`},
	}

	for _, tc := range testCases {
		s := string(tc.char)
		enc.Encode(buf, s, "info", time.Now(), nil)
		result := string(buf.Bytes())
		if !strings.Contains(result, tc.expected) {
			t.Errorf("char 0x%02x: expected %s, got: %s", tc.char, tc.expected, result)
		}
		buf.sb.Reset()
	}
}

func TestJSONEncoder_NewJSONEncoder_InitialCapacity(t *testing.T) {
	enc := NewJSONEncoder()

	// 验证预分配的 buf 容量
	if cap(enc.buf) != 256 {
		t.Errorf("expected initial capacity 256, got: %d", cap(enc.buf))
	}

	// 验证初始长度为 0
	if len(enc.buf) != 0 {
		t.Errorf("expected initial length 0, got: %d", len(enc.buf))
	}
}

func TestJSONEncoder_LargeFields(t *testing.T) {
	enc := NewJSONEncoder()
	buf := &mockBuffer{}

	// 创建大字符串字段
	largeString := strings.Repeat("a", 10000)
	fields := []Field{
		{Key: "large", Type: StringType, String: largeString},
	}

	enc.Encode(buf, "test", "info", time.Now(), fields)

	result := string(buf.Bytes())
	if len(result) < 10000 {
		t.Errorf("expected large string in output, got length: %d", len(result))
	}
}
