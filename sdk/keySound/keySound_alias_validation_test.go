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

package keySound

import "testing"

// TestAudioFileAliasExists
//
// 该测试直接覆盖“严格三元引用校验”的核心判定：
// - sha256 + name_id + type 全部匹配时返回 true；
// - 任一维度缺失或不匹配时返回 false。
func TestAudioFileAliasExists(t *testing.T) {
	get := func(key string) any {
		switch key {
		case "audio_files.sha_ok.type":
			return ".wav"
		case "audio_files.sha_ok.name.alias_ok":
			return "display_name"
		default:
			return nil
		}
	}

	if !audioFileAliasExists(get, "sha_ok", "alias_ok", ".wav") {
		t.Fatalf("expected alias exists when sha256 + name_id + type are all valid")
	}
	if audioFileAliasExists(get, "sha_ok", "alias_missing", ".wav") {
		t.Fatalf("expected alias missing when name_id does not exist")
	}
	if audioFileAliasExists(get, "sha_ok", "alias_ok", ".mp3") {
		t.Fatalf("expected alias missing when type does not match")
	}
	if audioFileAliasExists(get, "", "alias_ok", ".wav") {
		t.Fatalf("expected alias missing when sha256 is empty")
	}
	if audioFileAliasExists(get, "sha_ok", "", ".wav") {
		t.Fatalf("expected alias missing when name_id is empty")
	}
}
