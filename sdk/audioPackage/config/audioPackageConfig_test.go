/**
 * This file is part of the KeyTone project.
 *
 * Copyright (C) 2024 LuSrackhall
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package audioPackageConfig

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"KeyTone/logger"
)

// ============================================================================
// 测试工具函数
// ============================================================================

// 初始化日志系统（仅在测试中执行一次）
var loggerInitialized = false

func initLogger() {
	if !loggerInitialized {
		tmpLogDir, _ := os.MkdirTemp("", "audioPackageConfig_test_logs_")
		logFilePath := filepath.Join(tmpLogDir, "test.log")
		logger.InitLogger(logFilePath)
		loggerInitialized = true
	}
}

// 创建临时测试目录
func setupTestDir(t *testing.T) string {
	initLogger()
	tmpDir, err := os.MkdirTemp("", "audioPackageConfig_test_")
	if err != nil {
		t.Fatalf("无法创建临时目录: %v", err)
	}
	return tmpDir
}

// 清理测试目录
func teardownTestDir(tmpDir string) {
	os.RemoveAll(tmpDir)
}

// 初始化配置文件用于测试
func initTestConfig(t *testing.T, tmpDir string) {
	LoadConfig(tmpDir, true)
	if Viper == nil {
		t.Fatalf("初始化配置失败: Viper为nil")
	}
	// 等待初始化完成
	time.Sleep(100 * time.Millisecond)
}

// ============================================================================
// 测试序列定义
// ============================================================================

// TestCase 定义了一个通用的测试用例结构
type TestCase struct {
	name        string      // 测试用例名称
	key         string      // 要设置/获取/删除的键
	setValue    interface{} // 要设置的值
	expectedGet interface{} // 预期的获取值
	shouldExist bool        // 操作后该键是否应该存在
	description string      // 测试描述
}

// ============================================================================
// 第一层测试集：简单键值对（0-9, a-f）
// ============================================================================

func TestSimpleKeysSet_0to9(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []TestCase{
		{"set_0", "0", "value_0", "value_0", true, "设置简单键0"},
		{"set_1", "1", "value_1", "value_1", true, "设置简单键1"},
		{"set_2", "2", "value_2", "value_2", true, "设置简单键2"},
		{"set_3", "3", "value_3", "value_3", true, "设置简单键3"},
		{"set_4", "4", "value_4", "value_4", true, "设置简单键4"},
		{"set_5", "5", "value_5", "value_5", true, "设置简单键5"},
		{"set_6", "6", "value_6", "value_6", true, "设置简单键6"},
		{"set_7", "7", "value_7", "value_7", true, "设置简单键7"},
		{"set_8", "8", "value_8", "value_8", true, "设置简单键8"},
		{"set_9", "9", "value_9", "value_9", true, "设置简单键9"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetValue(tc.key, tc.setValue)
			got := GetValue(tc.key)
			if got != tc.expectedGet {
				t.Errorf("SetValue/GetValue 失败\n  键: %s\n  期望: %v\n  实际: %v", tc.key, tc.expectedGet, got)
			}
		})
	}
}

func TestSimpleKeysSet_a_to_f(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []TestCase{
		{"set_a", "a", "value_a", "value_a", true, "设置简单键a"},
		{"set_b", "b", "value_b", "value_b", true, "设置简单键b"},
		{"set_c", "c", "value_c", "value_c", true, "设置简单键c"},
		{"set_d", "d", "value_d", "value_d", true, "设置简单键d"},
		{"set_e", "e", "value_e", "value_e", true, "设置简单键e"},
		{"set_f", "f", "value_f", "value_f", true, "设置简单键f"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetValue(tc.key, tc.setValue)
			got := GetValue(tc.key)
			if got != tc.expectedGet {
				t.Errorf("SetValue/GetValue 失败\n  键: %s\n  期望: %v\n  实际: %v", tc.key, tc.expectedGet, got)
			}
		})
	}
}

// ============================================================================
// 第二层测试集：两层嵌套（使用点号分隔符）
// ============================================================================

func TestNestedKeys_TwoLevels(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []TestCase{
		{"nested_0.0", "0.0", "value_00", "value_00", true, "嵌套键0.0"},
		{"nested_1.1", "1.1", "value_11", "value_11", true, "嵌套键1.1"},
		{"nested_2.2", "2.2", map[string]interface{}{"sub": "value_22"}, nil, true, "嵌套键2.2（对象值）"},
		{"nested_3.3", "3.3", []interface{}{"a", "b", "c"}, nil, true, "嵌套键3.3（数组值）"},
		{"nested_4.a", "4.a", "value_4a", "value_4a", true, "嵌套键4.a"},
		{"nested_5.b", "5.b", 123.0, 123.0, true, "嵌套键5.b（浮点数值）"},
		{"nested_6.c", "6.c", 45.67, 45.67, true, "嵌套键6.c（浮点数值）"},
		{"nested_7.d", "7.d", true, true, true, "嵌套键7.d（布尔值true）"},
		{"nested_8.e", "8.e", false, false, true, "嵌套键8.e（布尔值false）"},
		{"nested_9.f", "9.f", "", "", true, "嵌套键9.f（空字符串）"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetValue(tc.key, tc.setValue)
			got := GetValue(tc.key)

			// 对于对象和数组，只验证是否存在非nil值
			if tc.expectedGet == nil {
				if got == nil {
					t.Errorf("SetValue/GetValue 失败\n  键: %s\n  期望值存在，但实际为nil", tc.key)
				}
			} else if got != tc.expectedGet {
				t.Errorf("SetValue/GetValue 失败\n  键: %s\n  期望: %v\n  实际: %v", tc.key, tc.expectedGet, got)
			}
		})
	}
}

// ============================================================================
// 第三层测试集：三层嵌套
// ============================================================================

func TestNestedKeys_ThreeLevels(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []TestCase{
		{"three_level_0", "0.0.0", "value_000", "value_000", true, "三层嵌套0.0.0"},
		{"three_level_1", "1.1.1", "value_111", "value_111", true, "三层嵌套1.1.1"},
		{"three_level_2", "2.a.2", "value_2a2", "value_2a2", true, "三层嵌套2.a.2"},
		{"three_level_3", "3.b.c", 999.0, 999.0, true, "三层嵌套3.b.c（浮点数值）"},
		{"three_level_4", "level.config.debug", true, true, true, "三层嵌套level.config.debug"},
		{"three_level_5", "database.connection.timeout", 30.0, 30.0, true, "三层嵌套database.connection.timeout"},
		{"three_level_6", "audio.package.volume", 0.85, 0.85, true, "三层嵌套audio.package.volume"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetValue(tc.key, tc.setValue)
			got := GetValue(tc.key)
			if got != tc.expectedGet {
				t.Errorf("SetValue/GetValue 失败\n  键: %s\n  期望: %v\n  实际: %v", tc.key, tc.expectedGet, got)
			}
		})
	}
}

// ============================================================================
// 第四层测试集：四层嵌套
// ============================================================================

func TestNestedKeys_FourLevels(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []TestCase{
		{"four_level_0", "0.0.0.0", "value_0000", "value_0000", true, "四层嵌套0.0.0.0"},
		{"four_level_1", "a.b.c.d", "nested_value", "nested_value", true, "四层嵌套a.b.c.d"},
		{"four_level_2", "settings.audio.playback.quality", "high", "high", true, "四层嵌套settings.audio.playback.quality"},
		{"four_level_3", "system.config.database.host", "localhost", "localhost", true, "四层嵌套system.config.database.host"},
		{"four_level_4", "app.theme.dark.opacity", 0.95, 0.95, true, "四层嵌套app.theme.dark.opacity"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetValue(tc.key, tc.setValue)
			got := GetValue(tc.key)
			if got != tc.expectedGet {
				t.Errorf("SetValue/GetValue 失败\n  键: %s\n  期望: %v\n  实际: %v", tc.key, tc.expectedGet, got)
			}
		})
	}
}

// ============================================================================
// 第五层测试集：五层嵌套（深层结构）
// ============================================================================

func TestNestedKeys_FiveLevels(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []TestCase{
		{"five_level_0", "0.0.0.0.0", "value_00000", "value_00000", true, "五层嵌套0.0.0.0.0"},
		{"five_level_1", "a.b.c.d.e", "deep_nested_value", "deep_nested_value", true, "五层嵌套a.b.c.d.e"},
		{"five_level_2", "database.connection.pool.settings.timeout", 60.0, 60.0, true, "五层嵌套database.connection.pool.settings.timeout"},
		{"five_level_3", "app.modules.audio.effects.reverb.level", 0.5, 0.5, true, "五层嵌套app.modules.audio.effects.reverb.level"},
		{"five_level_4", "config.server.api.auth.oauth.token_expiry", 3600.0, 3600.0, true, "五层嵌套config.server.api.auth.oauth.token_expiry"},
		{"five_level_5", "x.y.z.w.q", "end_of_deep_chain", "end_of_deep_chain", true, "五层嵌套x.y.z.w.q"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetValue(tc.key, tc.setValue)
			got := GetValue(tc.key)
			if got != tc.expectedGet {
				t.Errorf("SetValue/GetValue 失败\n  键: %s\n  期望: %v\n  实际: %v", tc.key, tc.expectedGet, got)
			}
		})
	}
}

// ============================================================================
// 第六层测试集：值的类型覆盖
// ============================================================================

func TestValueTypes(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []TestCase{
		{"string_value", "type.string", "hello world", "hello world", true, "字符串值"},
		{"int_value", "type.int", 42.0, 42.0, true, "整数值（JSON中为浮点数）"},
		{"float_value", "type.float", 3.14159, 3.14159, true, "浮点数值"},
		{"bool_true", "type.bool_true", true, true, true, "布尔值true"},
		{"bool_false", "type.bool_false", false, false, true, "布尔值false"},
		{"empty_string", "type.empty_string", "", "", true, "空字符串"},
		{"zero_int", "type.zero_int", 0.0, 0.0, true, "整数零（JSON中为0.0）"},
		{"negative_int", "type.negative_int", -100.0, -100.0, true, "负整数（JSON中为浮点数）"},
		{"large_number", "type.large_number", 9999999999.0, 9999999999.0, true, "大整数（JSON中为浮点数）"},
		{"small_float", "type.small_float", 0.00001, 0.00001, true, "小浮点数"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetValue(tc.key, tc.setValue)
			time.Sleep(50 * time.Millisecond)
			got := GetValue(tc.key)
			if got != tc.expectedGet {
				t.Errorf("SetValue/GetValue 类型测试失败\n  键: %s\n  期望: %v (类型: %T)\n  实际: %v (类型: %T)",
					tc.key, tc.expectedGet, tc.expectedGet, got, got)
			}
		})
	}
}

// ============================================================================
// 第七层测试集：删除操作测试
// ============================================================================

func TestDeleteValue_SimpleKeys(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []struct {
		name string
		key  string
		desc string
	}{
		{"delete_0", "del.0", "删除简单键del.0"},
		{"delete_1", "del.1", "删除简单键del.1"},
		{"delete_2", "del.2", "删除简单键del.2"},
		{"delete_a", "del.a", "删除简单键del.a"},
		{"delete_b", "del.b", "删除简单键del.b"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 先设置值
			SetValue(tc.key, "test_value")
			time.Sleep(50 * time.Millisecond)

			// 验证值已设置
			if got := GetValue(tc.key); got == nil {
				t.Errorf("预设值失败\n  键: %s", tc.key)
				return
			}

			// 删除值
			DeleteValue(tc.key)
			time.Sleep(100 * time.Millisecond)

			// 验证值已删除（应为nil或不存在）
			if got := GetValue(tc.key); got != nil {
				t.Errorf("DeleteValue 失败\n  键: %s\n  期望: nil\n  实际: %v", tc.key, got)
			}
		})
	}
}

func TestDeleteValue_NestedKeys(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []struct {
		name string
		key  string
		desc string
	}{
		{"delete_nested_0", "nested.del.0", "删除嵌套键nested.del.0"},
		{"delete_nested_1", "nested.del.1.sub", "删除嵌套键nested.del.1.sub"},
		{"delete_nested_2", "a.b.c.d.del", "删除深层嵌套键a.b.c.d.del"},
		{"delete_nested_3", "config.settings.delete.test", "删除嵌套键config.settings.delete.test"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 先设置值
			SetValue(tc.key, "test_nested_value")
			time.Sleep(50 * time.Millisecond)

			// 验证值已设置
			if got := GetValue(tc.key); got == nil {
				t.Errorf("预设值失败\n  键: %s", tc.key)
				return
			}

			// 删除值
			DeleteValue(tc.key)
			time.Sleep(100 * time.Millisecond)

			// 验证值已删除
			if got := GetValue(tc.key); got != nil {
				t.Errorf("DeleteValue 失败\n  键: %s\n  期望: nil\n  实际: %v", tc.key, got)
			}
		})
	}
}

// ============================================================================
// 第八层测试集：修改覆盖测试
// ============================================================================

func TestUpdateValue_Overwrite(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []struct {
		name   string
		key    string
		value1 interface{}
		value2 interface{}
		desc   string
	}{
		{"overwrite_0", "overwrite.0", "old_value", "new_value", "字符串覆盖"},
		{"overwrite_1", "overwrite.1", 10.0, 20.0, "整数覆盖"},
		{"overwrite_2", "overwrite.2", 0.5, 0.8, "浮点数覆盖"},
		{"overwrite_3", "overwrite.3", true, false, "布尔值覆盖"},
		{"overwrite_4", "overwrite.4", "string", 100.0, "类型转换覆盖（字符串→浮点数）"},
		{"overwrite_5", "nested.overwrite.5", "value1", "value2", "嵌套键字符串覆盖"},
		{"overwrite_6", "a.b.c.d.overwrite", 1.0, 2.0, "深层嵌套整数覆盖"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 设置第一个值
			SetValue(tc.key, tc.value1)
			time.Sleep(50 * time.Millisecond)
			got1 := GetValue(tc.key)
			if got1 != tc.value1 {
				t.Errorf("初始值设置失败\n  键: %s\n  期望: %v\n  实际: %v", tc.key, tc.value1, got1)
				return
			}

			// 覆盖为第二个值
			SetValue(tc.key, tc.value2)
			time.Sleep(50 * time.Millisecond)
			got2 := GetValue(tc.key)
			if got2 != tc.value2 {
				t.Errorf("覆盖值失败\n  键: %s\n  期望: %v\n  实际: %v", tc.key, tc.value2, got2)
			}
		})
	}
}

// ============================================================================
// 第九层测试集：边缘情况和特殊场景
// ============================================================================

func TestEdgeCases_SpecialCharacters(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []TestCase{
		{"special_underscore", "edge_case_0", "value_with_underscore", "value_with_underscore", true, "键包含下划线"},
		{"special_hyphen", "edge-case-1", "value-with-hyphen", "value-with-hyphen", true, "键包含连字符"},
		{"special_numbers", "test123case456", "value_123_456", "value_123_456", true, "键包含数字"},
		{"special_mixed", "mix_case-test.value", "mixed_special_value", "mixed_special_value", true, "键混合特殊字符和点"},
		{"special_long_key", "very.long.key.path.with.many.segments.0", "long_key_value", "long_key_value", true, "很长的键路径"},
		{"special_unicode", "键.unicode.值", "unicode_value", "unicode_value", true, "键包含unicode字符"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetValue(tc.key, tc.setValue)
			time.Sleep(50 * time.Millisecond)
			got := GetValue(tc.key)
			if got != tc.expectedGet {
				t.Errorf("特殊字符测试失败\n  键: %s\n  期望: %v\n  实际: %v", tc.key, tc.expectedGet, got)
			}
		})
	}
}

func TestEdgeCases_ValueContent(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []TestCase{
		{"edge_empty_string", "edge.empty", "", "", true, "空字符串值"},
		{"edge_whitespace", "edge.space", "   ", "   ", true, "仅空格值"},
		{"edge_newline", "edge.newline", "line1\nline2\nline3", "line1\nline2\nline3", true, "包含换行符的值"},
		{"edge_tab", "edge.tab", "col1\tcol2\tcol3", "col1\tcol2\tcol3", true, "包含制表符的值"},
		{"edge_quotes", "edge.quotes", "\"quoted\"", "\"quoted\"", true, "包含引号的值"},
		{"edge_special_chars", "edge.special", "!@#$%^&*()", "!@#$%^&*()", true, "特殊符号值"},
		{"edge_json_like", "edge.json", "{\"key\":\"value\"}", "{\"key\":\"value\"}", true, "JSON格式的字符串值"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetValue(tc.key, tc.setValue)
			time.Sleep(50 * time.Millisecond)
			got := GetValue(tc.key)
			if got != tc.expectedGet {
				t.Errorf("值内容测试失败\n  键: %s\n  期望: %v\n  实际: %v", tc.key, tc.expectedGet, got)
			}
		})
	}
}

// ============================================================================
// 第十层测试集：批量操作一致性测试
// ============================================================================

func TestBatchOperations_Consistency(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	// 测试用例序列：按顺序执行多个操作，验证一致性
	sequences := []struct {
		name string
		ops  []struct {
			operation string // "set", "get", "delete"
			key       string
			value     interface{}
			expected  interface{}
		}
	}{
		{
			name: "sequence_0_set_get_delete",
			ops: []struct {
				operation string
				key       string
				value     interface{}
				expected  interface{}
			}{
				{"set", "seq.0.a", "value_a", nil},
				{"set", "seq.0.b", "value_b", nil},
				{"get", "seq.0.a", nil, "value_a"},
				{"get", "seq.0.b", nil, "value_b"},
				{"delete", "seq.0.a", nil, nil},
				{"get", "seq.0.a", nil, nil},
				{"get", "seq.0.b", nil, "value_b"},
			},
		},
		{
			name: "sequence_1_nested_consistency",
			ops: []struct {
				operation string
				key       string
				value     interface{}
				expected  interface{}
			}{
				{"set", "seq.1.level1.level2", "value_l2", nil},
				{"set", "seq.1.level1.level3", "value_l3", nil},
				{"get", "seq.1.level1.level2", nil, "value_l2"},
				{"get", "seq.1.level1.level3", nil, "value_l3"},
				{"set", "seq.1.level1.level2", "updated_l2", nil},
				{"get", "seq.1.level1.level2", nil, "updated_l2"},
			},
		},
	}

	for _, seq := range sequences {
		t.Run(seq.name, func(t *testing.T) {
			for i, op := range seq.ops {
				switch op.operation {
				case "set":
					SetValue(op.key, op.value)
					time.Sleep(50 * time.Millisecond)
				case "get":
					got := GetValue(op.key)
					if got != op.expected {
						t.Errorf("批量操作一致性测试失败 [操作%d]\n  操作: %s\n  键: %s\n  期望: %v\n  实际: %v",
							i, op.operation, op.key, op.expected, got)
					}
				case "delete":
					DeleteValue(op.key)
					time.Sleep(100 * time.Millisecond)
				}
			}
		})
	}
}

// ============================================================================
// 第十一层测试集：数值范围和边界测试
// ============================================================================

func TestBoundaryValues(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []TestCase{
		{"boundary_zero", "boundary.zero", 0.0, 0.0, true, "零值"},
		{"boundary_neg_one", "boundary.neg_one", -1.0, -1.0, true, "负一"},
		{"boundary_pos_one", "boundary.pos_one", 1.0, 1.0, true, "正一"},
		{"boundary_max_int", "boundary.max_int", 9223372036854775807.0, 9223372036854775807.0, true, "最大int64"},
		{"boundary_min_int", "boundary.min_int", -9223372036854775808.0, -9223372036854775808.0, true, "最小int64"},
		{"boundary_max_float", "boundary.max_float", 1.7976931348623157e+308, 1.7976931348623157e+308, true, "最大float64"},
		{"boundary_min_float", "boundary.min_float", 2.2250738585072014e-308, 2.2250738585072014e-308, true, "最小float64"},
		{"boundary_zero_float", "boundary.zero_float", 0.0, 0.0, true, "浮点零"},
		{"boundary_neg_float", "boundary.neg_float", -0.5, -0.5, true, "负浮点数"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			SetValue(tc.key, tc.setValue)
			time.Sleep(50 * time.Millisecond)
			got := GetValue(tc.key)
			if got != tc.expectedGet {
				t.Errorf("边界值测试失败\n  键: %s\n  期望: %v\n  实际: %v", tc.key, tc.expectedGet, got)
			}
		})
	}
}

// ============================================================================
// 第十二层测试集：文件持久化验证
// ============================================================================

func TestPersistence_FileWrite(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	// 设置多个值
	testData := map[string]interface{}{
		"persist.0":                "value_0",
		"persist.1":                100,
		"persist.2.nested":         "nested_value",
		"persist.3.deep.very.deep": "deep_value",
	}

	for key, value := range testData {
		SetValue(key, value)
		time.Sleep(50 * time.Millisecond)
	}

	// 验证文件确实存在并且包含数据
	configFile := filepath.Join(tmpDir, "config.json")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Errorf("配置文件未被创建: %s", configFile)
	}

	// 验证文件内容不为空
	fileInfo, err := os.Stat(configFile)
	if err != nil {
		t.Errorf("无法获取配置文件信息: %v", err)
	}
	if fileInfo.Size() == 0 {
		t.Errorf("配置文件为空，数据未被正确持久化")
	}

	// 重新加载配置
	LoadConfig(tmpDir, false)
	time.Sleep(100 * time.Millisecond)

	// 验证重新加载后数据仍然存在
	for key, expectedValue := range testData {
		got := GetValue(key)
		if got == nil {
			t.Errorf("数据持久化失败 - 重新加载后键不存在\n  键: %s\n  期望: %v", key, expectedValue)
		}
	}
}

// ============================================================================
// 第十五层测试集：JSON文件内容持久化验证 ⭐ 新增
// ============================================================================

// 辅助函数：读取JSON文件并解析为Map
func readConfigFileAsMap(filePath string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// 辅助函数：获取嵌套的JSON值（通过点号路径）
func getNestedValue(data map[string]interface{}, path string) interface{} {
	keys := strings.Split(path, ".")
	var current interface{} = data

	for _, key := range keys {
		switch v := current.(type) {
		case map[string]interface{}:
			current = v[key]
		default:
			return nil
		}
	}

	return current
}

func TestFileContent_SimpleKeysPersistence(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []struct {
		key   string
		value interface{}
		desc  string
	}{
		{"file_0", "file_value_0", "简单键file_0的文件持久化"},
		{"file_1", "file_value_1", "简单键file_1的文件持久化"},
		{"file_a", "file_value_a", "简单键file_a的文件持久化"},
		{"file_b", "file_value_b", "简单键file_b的文件持久化"},
	}

	configFilePath := filepath.Join(tmpDir, "config.json")

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			// 设置值
			SetValue(tc.key, tc.value)
			time.Sleep(150 * time.Millisecond)

			// 读取JSON文件
			fileData, err := readConfigFileAsMap(configFilePath)
			if err != nil {
				t.Fatalf("无法读取配置文件: %v", err)
			}

			// 验证文件中的值
			if got, exists := fileData[tc.key]; !exists {
				t.Errorf("文件中键不存在\n  键: %s\n  描述: %s", tc.key, tc.desc)
			} else if got != tc.value {
				t.Errorf("文件中的值不匹配\n  键: %s\n  期望: %v\n  实际: %v\n  描述: %s",
					tc.key, tc.value, got, tc.desc)
			}
		})
	}
}

func TestFileContent_NestedKeysPersistence_TwoLevels(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []struct {
		key   string
		value interface{}
		desc  string
	}{
		{"nested.file.0", "nested_value_0", "两层嵌套nested.file.0的文件持久化"},
		{"nested.file.1", "nested_value_1", "两层嵌套nested.file.1的文件持久化"},
		{"config.debug", true, "两层嵌套config.debug（布尔值）的文件持久化"},
		{"app.version", "1.0.0", "两层嵌套app.version的文件持久化"},
		{"system.a", 100.0, "两层嵌套system.a（数值）的文件持久化"},
	}

	configFilePath := filepath.Join(tmpDir, "config.json")

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			// 设置值
			SetValue(tc.key, tc.value)
			time.Sleep(150 * time.Millisecond)

			// 读取JSON文件
			fileData, err := readConfigFileAsMap(configFilePath)
			if err != nil {
				t.Fatalf("无法读取配置文件: %v", err)
			}

			// 使用嵌套路径获取值
			got := getNestedValue(fileData, tc.key)
			if got == nil {
				t.Errorf("文件中嵌套键不存在\n  键: %s\n  描述: %s", tc.key, tc.desc)
			} else if got != tc.value {
				t.Errorf("文件中的嵌套值不匹配\n  键: %s\n  期望: %v\n  实际: %v\n  描述: %s",
					tc.key, tc.value, got, tc.desc)
			}
		})
	}
}

func TestFileContent_NestedKeysPersistence_FiveLevels(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []struct {
		key   string
		value interface{}
		desc  string
	}{
		{
			"database.connection.pool.settings.timeout",
			60.0,
			"五层嵌套database.connection.pool.settings.timeout的文件持久化",
		},
		{
			"app.modules.audio.effects.reverb",
			0.5,
			"五层嵌套app.modules.audio.effects.reverb的文件持久化",
		},
		{
			"config.server.api.auth.oauth",
			true,
			"五层嵌套config.server.api.auth.oauth（布尔值）的文件持久化",
		},
		{
			"system.storage.cache.redis.host",
			"127.0.0.1",
			"五层嵌套system.storage.cache.redis.host的文件持久化",
		},
		{
			"component.ui.theme.dark.opacity",
			0.8,
			"五层嵌套component.ui.theme.dark.opacity的文件持久化",
		},
	}

	configFilePath := filepath.Join(tmpDir, "config.json")

	for _, tc := range testCases {
		t.Run(strings.ReplaceAll(tc.key, ".", "_"), func(t *testing.T) {
			// 设置值
			SetValue(tc.key, tc.value)
			time.Sleep(150 * time.Millisecond)

			// 读取JSON文件
			fileData, err := readConfigFileAsMap(configFilePath)
			if err != nil {
				t.Fatalf("无法读取配置文件: %v", err)
			}

			// 使用嵌套路径获取值
			got := getNestedValue(fileData, tc.key)
			if got == nil {
				t.Errorf("文件中五层嵌套键不存在\n  键: %s\n  描述: %s", tc.key, tc.desc)
			} else if got != tc.value {
				t.Errorf("文件中的五层嵌套值不匹配\n  键: %s\n  期望: %v\n  实际: %v\n  描述: %s",
					tc.key, tc.value, got, tc.desc)
			}
		})
	}
}

func TestFileContent_DeletePersistence(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []struct {
		key  string
		desc string
	}{
		{"delete.file.0", "删除文件中的delete.file.0"},
		{"delete.file.1", "删除文件中的delete.file.1"},
		{"system.delete.nested", "删除文件中的system.delete.nested"},
	}

	configFilePath := filepath.Join(tmpDir, "config.json")

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			// 设置值
			SetValue(tc.key, "to_be_deleted")
			time.Sleep(100 * time.Millisecond)

			// 验证文件中值存在
			fileData1, err := readConfigFileAsMap(configFilePath)
			if err != nil {
				t.Fatalf("无法读取配置文件: %v", err)
			}
			got1 := getNestedValue(fileData1, tc.key)
			if got1 == nil {
				t.Errorf("初始值未成功写入文件\n  键: %s", tc.key)
				return
			}

			// 删除值
			DeleteValue(tc.key)
			time.Sleep(150 * time.Millisecond)

			// 验证文件中值已删除
			fileData2, err := readConfigFileAsMap(configFilePath)
			if err != nil {
				t.Fatalf("无法读取配置文件: %v", err)
			}
			got2 := getNestedValue(fileData2, tc.key)
			if got2 != nil {
				t.Errorf("文件中的值未被正确删除\n  键: %s\n  描述: %s\n  残留值: %v",
					tc.key, tc.desc, got2)
			}
		})
	}
}

func TestFileContent_UpdatePersistence(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	testCases := []struct {
		key    string
		value1 interface{}
		value2 interface{}
		desc   string
	}{
		{"update.file.0", "old_value", "new_value", "覆盖文件中的update.file.0（字符串）"},
		{"update.file.1", 10.0, 20.0, "覆盖文件中的update.file.1（数值）"},
		{"config.setting.flag", false, true, "覆盖文件中的config.setting.flag（布尔值）"},
		{"system.database.port", 5432.0, 3306.0, "覆盖文件中的system.database.port"},
	}

	configFilePath := filepath.Join(tmpDir, "config.json")

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			// 设置第一个值
			SetValue(tc.key, tc.value1)
			time.Sleep(100 * time.Millisecond)

			// 验证第一个值在文件中
			fileData1, err := readConfigFileAsMap(configFilePath)
			if err != nil {
				t.Fatalf("无法读取配置文件: %v", err)
			}
			got1 := getNestedValue(fileData1, tc.key)
			if got1 != tc.value1 {
				t.Errorf("第一次写入文件失败\n  键: %s\n  期望: %v\n  实际: %v",
					tc.key, tc.value1, got1)
				return
			}

			// 设置第二个值（覆盖）
			SetValue(tc.key, tc.value2)
			time.Sleep(100 * time.Millisecond)

			// 验证第二个值在文件中
			fileData2, err := readConfigFileAsMap(configFilePath)
			if err != nil {
				t.Fatalf("无法读取配置文件: %v", err)
			}
			got2 := getNestedValue(fileData2, tc.key)
			if got2 != tc.value2 {
				t.Errorf("覆盖写入文件失败\n  键: %s\n  期望: %v\n  实际: %v\n  描述: %s",
					tc.key, tc.value2, got2, tc.desc)
			}
		})
	}
}

func TestFileContent_BatchOperationsPersistence(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	configFilePath := filepath.Join(tmpDir, "config.json")

	// 测试场景：连续的多个操作
	operations := []struct {
		op    string // "set" 或 "delete"
		key   string
		value interface{}
	}{
		{"set", "batch.0", "value_0"},
		{"set", "batch.1", "value_1"},
		{"set", "batch.2.nested", 100.0},
		{"set", "batch.3.deep.nested", true},
		{"set", "config.a.b.c", "deep_value"},
		{"set", "system.x.y.z", 42.0},
	}

	// 执行所有设置操作
	for _, op := range operations {
		SetValue(op.key, op.value)
		time.Sleep(80 * time.Millisecond)
	}

	// 验证所有值都在文件中
	fileData, err := readConfigFileAsMap(configFilePath)
	if err != nil {
		t.Fatalf("无法读取配置文件: %v", err)
	}

	for _, op := range operations {
		got := getNestedValue(fileData, op.key)
		if got != op.value {
			t.Errorf("批量操作文件持久化失败\n  键: %s\n  期望: %v\n  实际: %v",
				op.key, op.value, got)
		}
	}

	// 现在删除部分键
	keysToDelete := []string{"batch.0", "batch.2.nested", "config.a.b.c"}
	for _, key := range keysToDelete {
		DeleteValue(key)
		time.Sleep(100 * time.Millisecond)
	}

	// 验证被删除的键不在文件中，其他键仍存在
	fileDataAfterDelete, err := readConfigFileAsMap(configFilePath)
	if err != nil {
		t.Fatalf("无法读取配置文件: %v", err)
	}

	for _, op := range operations {
		got := getNestedValue(fileDataAfterDelete, op.key)
		isDeleted := false
		for _, delKey := range keysToDelete {
			if op.key == delKey {
				isDeleted = true
				break
			}
		}

		if isDeleted {
			if got != nil {
				t.Errorf("删除操作在文件中未生效\n  键: %s\n  期望: nil\n  实际: %v", op.key, got)
			}
		} else {
			if got != op.value {
				t.Errorf("非删除键在文件中被意外修改\n  键: %s\n  期望: %v\n  实际: %v",
					op.key, op.value, got)
			}
		}
	}
}

func TestFileContent_DataTypeConsistency(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	configFilePath := filepath.Join(tmpDir, "config.json")

	testCases := []struct {
		key         string
		value       interface{}
		description string
	}{
		{"type.string", "hello_world", "字符串类型在文件中的一致性"},
		{"type.float", 3.14159, "浮点数类型在文件中的一致性"},
		{"type.bool_true", true, "布尔真值在文件中的一致性"},
		{"type.bool_false", false, "布尔假值在文件中的一致性"},
		{"type.empty_string", "", "空字符串在文件中的一致性"},
		{"type.zero", 0.0, "零值在文件中的一致性"},
		{"type.negative", -100.0, "负数在文件中的一致性"},
	}

	for _, tc := range testCases {
		t.Run(tc.key, func(t *testing.T) {
			// 设置值
			SetValue(tc.key, tc.value)
			time.Sleep(100 * time.Millisecond)

			// 从文件中读取
			fileData, err := readConfigFileAsMap(configFilePath)
			if err != nil {
				t.Fatalf("无法读取配置文件: %v", err)
			}

			// 从内存中读取（通过GetValue）
			memoryValue := GetValue(tc.key)

			// 从文件中获取
			fileValue := getNestedValue(fileData, tc.key)

			// 验证文件值与内存值一致
			if fileValue != memoryValue {
				t.Errorf("文件值与内存值不一致\n  键: %s\n  文件值: %v (类型: %T)\n  内存值: %v (类型: %T)\n  描述: %s",
					tc.key, fileValue, fileValue, memoryValue, memoryValue, tc.description)
			}

			// 验证文件值与预期值一致
			if fileValue != tc.value {
				t.Errorf("文件值与预期值不一致\n  键: %s\n  期望: %v\n  实际: %v\n  描述: %s",
					tc.key, tc.value, fileValue, tc.description)
			}
		})
	}
}

func TestFileContent_JSONStructureValidity(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	configFilePath := filepath.Join(tmpDir, "config.json")

	// 设置多个不同层级的键
	complexData := []struct {
		key   string
		value interface{}
	}{
		{"root.level1.level2.level3.level4", "value1"},
		{"root.level1.level2.level3.other", "value2"},
		{"root.level1.other", "value3"},
		{"other.branch.deep", "value4"},
		{"simple", "value5"},
	}

	for _, data := range complexData {
		SetValue(data.key, data.value)
		time.Sleep(80 * time.Millisecond)
	}

	// 读取JSON文件
	rawData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		t.Fatalf("无法读取配置文件: %v", err)
	}

	// 验证JSON格式是否有效
	var fileData map[string]interface{}
	err = json.Unmarshal(rawData, &fileData)
	if err != nil {
		t.Errorf("JSON格式无效\n  错误: %v\n  文件内容: %s", err, string(rawData))
		return
	}

	// 验证所有预期的键都存在于正确的嵌套结构中
	for _, data := range complexData {
		got := getNestedValue(fileData, data.key)
		if got != data.value {
			t.Errorf("JSON结构中的值不匹配\n  键: %s\n  期望: %v\n  实际: %v",
				data.key, data.value, got)
		}
	}

	// 打印JSON结构用于调试
	t.Logf("JSON结构内容:\n%s", string(rawData))
}

// ============================================================================
// 综合测试：完整的A-F序列加深层路径
// ============================================================================

func TestComprehensive_FullSequence(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer teardownTestDir(tmpDir)
	initTestConfig(t, tmpDir)

	// 定义完整的测试矩阵
	sequences := []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
		"a", "b", "c", "d", "e", "f",
	}

	// 第一阶段：设置所有简单键
	for _, seq := range sequences {
		SetValue(seq, "simple_"+seq)
		time.Sleep(30 * time.Millisecond)
	}

	// 第二阶段：验证所有简单键
	for _, seq := range sequences {
		got := GetValue(seq)
		expected := "simple_" + seq
		if got != expected {
			t.Errorf("简单键验证失败\n  键: %s\n  期望: %s\n  实际: %v", seq, expected, got)
		}
	}

	// 第三阶段：设置嵌套键（两层）
	for i, seq1 := range sequences {
		for j, seq2 := range sequences {
			if (i+j)%3 == 0 { // 创建部分组合以避免过多文件IO
				key := seq1 + "." + seq2
				SetValue(key, "nested_"+seq1+"_"+seq2)
				time.Sleep(20 * time.Millisecond)
			}
		}
	}

	// 第四阶段：删除部分简单键
	for i, seq := range sequences {
		if i%2 == 0 { // 删除偶数位置的键
			DeleteValue(seq)
			time.Sleep(50 * time.Millisecond)

			// 验证删除
			if got := GetValue(seq); got != nil {
				t.Errorf("删除验证失败\n  键: %s\n  期望: nil\n  实际: %v", seq, got)
			}
		}
	}
}
