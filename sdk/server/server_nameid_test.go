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

package server

import (
	"regexp"
	"testing"
)

// TestGenerateAudioSourceNameID_FormatAndUniqueness
//
// 验证目标：
// 1) 生成值符合 UUID v4 字符串格式；
// 2) 批量生成时不出现重复（用于保障 name_id 不会复用）。
func TestGenerateAudioSourceNameID_FormatAndUniqueness(t *testing.T) {
	uuidV4Pattern := regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)
	seen := make(map[string]struct{})

	for i := 0; i < 500; i++ {
		id, err := generateAudioSourceNameID()
		if err != nil {
			t.Fatalf("generateAudioSourceNameID returned error: %v", err)
		}
		if !uuidV4Pattern.MatchString(id) {
			t.Fatalf("generated id is not UUID v4 format: %s", id)
		}
		if _, exists := seen[id]; exists {
			t.Fatalf("generated duplicate id: %s", id)
		}
		seen[id] = struct{}{}
	}
}
