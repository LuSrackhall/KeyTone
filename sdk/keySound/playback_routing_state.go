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

// =============================
// Playback Routing（播放路由）说明
// =============================
//
// 目标：在“播放热路径”（KeyEvent -> KeySoundHandler）中做到零磁盘 IO。
// 实现方式：
//   1) 前端在关键时机调用 /keytone_pkg/apply_playback_routing
//   2) 后端在该 API 中一次性读取专辑配置（支持 plain / legacy-hex / core 三种形态）
//   3) 将解析后的只读 Viper 实例（Snapshot）写入内存态 playbackState
//   4) 播放热路径仅通过 GetPlaybackState() 获取快照指针并读取键值
//
// 关键约束：
//   - playbackState 允许并发读取（RWMutex RLock），播放热路径只读。
//   - apply / mode 切换会整体替换内存态，避免读写同一结构导致竞态。
//   - editor 模式不使用只读快照，直接读取 audioPackageConfig.Viper（编辑器可写配置）。

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	audioPackageConfig "KeyTone/audioPackage/config"
	"KeyTone/audioPackage/enc"

	"github.com/spf13/viper"
)

const (
	// SourceModeRouteUnified 表示使用统一路由快照（键盘/鼠标共用一个专辑）。
	SourceModeRouteUnified = "route-unified"
	// SourceModeRouteSplit 表示使用分离路由快照（键盘/鼠标各自专辑）。
	SourceModeRouteSplit = "route-split"
	// SourceModeEditor 表示进入编辑试听模式，直接读取 audioPackageConfig.Viper。
	SourceModeEditor = "editor"
)

type AlbumSnapshot struct {
	AlbumPath string
	AlbumUUID string
	Viper     *viper.Viper
}

func (s *AlbumSnapshot) GetValue(key string) any {
	if s == nil || s.Viper == nil {
		return nil
	}
	return s.Viper.Get(key)
}

func (s *AlbumSnapshot) AudioPkgUUID() string {
	if s == nil {
		return ""
	}
	// 兼容：若配置内已有 audio_pkg_uuid，优先使用。
	if value, ok := s.GetValue("audio_pkg_uuid").(string); ok && strings.TrimSpace(value) != "" {
		return value
	}
	// 回退：使用专辑目录名（albumUUID）。
	return s.AlbumUUID
}

type RoutingSourceStatus struct {
	// RequestedPath 为前端提交的路径/UUID。
	RequestedPath string `json:"requestedPath"`
	// ResolvedPath 为 SDK 解析后的真实路径。
	ResolvedPath string `json:"resolvedPath"`
	// Loaded 表示该来源是否成功加载为快照。
	Loaded bool `json:"loaded"`
	// Error 为加载失败时的错误信息。
	Error string `json:"error"`
}

type ApplyRoutingResult struct {
	// Mode 为请求的路由模式（unified/split）。
	Mode     string              `json:"mode"`
	Unified  RoutingSourceStatus `json:"unified"`
	Keyboard RoutingSourceStatus `json:"keyboard"`
	Mouse    RoutingSourceStatus `json:"mouse"`
}

type PlaybackRoutingState struct {
	// Mode 为 routing 子状态（unified/split）。
	Mode              string
	UnifiedAlbumPath  string
	KeyboardAlbumPath string
	MouseAlbumPath    string
	UnifiedSnapshot   *AlbumSnapshot
	KeyboardSnapshot  *AlbumSnapshot
	MouseSnapshot     *AlbumSnapshot
}

type PlaybackState struct {
	// SourceMode 为当前播放来源模式（route-unified/route-split/editor）。
	SourceMode string
	Routing    PlaybackRoutingState
	// EditorAlbumPath 仅用于记录编辑模式的专辑路径，便于诊断与状态展示。
	EditorAlbumPath string
}

var playbackStateLock sync.RWMutex
var playbackState = PlaybackState{
	SourceMode: SourceModeRouteUnified,
	Routing: PlaybackRoutingState{
		Mode: "unified",
	},
}

func GetPlaybackState() PlaybackState {
	// 只读快照：拷贝当前播放状态，避免热路径竞争。
	playbackStateLock.RLock()
	defer playbackStateLock.RUnlock()
	return playbackState
}

func SetPlaybackSourceMode(mode string, editorAlbumPath string) (PlaybackState, error) {
	mode = strings.TrimSpace(mode)
	if mode == "" {
		return PlaybackState{}, errors.New("mode is required")
	}
	playbackStateLock.Lock()
	defer playbackStateLock.Unlock()

	switch mode {
	case "editor":
		// editor 模式：播放使用可写配置（audioPackageConfig.Viper），用于编辑页“试听/预览”。
		// 注意：这里只记录 editorAlbumPath 便于诊断；真正的播放读取逻辑在 KeySoundHandler。
		playbackState.SourceMode = SourceModeEditor
		if strings.TrimSpace(editorAlbumPath) != "" {
			playbackState.EditorAlbumPath = editorAlbumPath
		}
	case "route":
		// route 模式：播放使用只读快照（Snapshot）。
		// route 的具体形态由 playbackState.Routing.Mode 决定：unified => route-unified；split => route-split。
		if strings.EqualFold(playbackState.Routing.Mode, "split") {
			playbackState.SourceMode = SourceModeRouteSplit
		} else {
			playbackState.SourceMode = SourceModeRouteUnified
		}
	default:
		return PlaybackState{}, fmt.Errorf("unsupported mode: %s", mode)
	}
	return playbackState, nil
}

func ApplyPlaybackRouting(mode, unifiedPath, keyboardPath, mousePath string) (ApplyRoutingResult, error) {
	mode = strings.TrimSpace(mode)
	if mode != "unified" && mode != "split" {
		return ApplyRoutingResult{}, fmt.Errorf("unsupported routing mode: %s", mode)
	}

	result := ApplyRoutingResult{Mode: mode}

	if mode == "unified" {
		// 统一模式：单快照复用至键盘/鼠标。
		// 兼容性目标：老用户升级后默认行为与旧版本一致（一个专辑同时用于键盘+鼠标）。
		snapshot, resolved, err := loadSnapshotWithResolve(unifiedPath)
		result.Unified = buildStatus(unifiedPath, resolved, snapshot, err)
		result.Keyboard = result.Unified
		result.Mouse = result.Unified

		// 以“整块写入”的方式更新内存态，保证读侧（热路径）始终看到一致结构。
		playbackStateLock.Lock()
		playbackState.Routing.Mode = "unified"
		playbackState.Routing.UnifiedAlbumPath = resolved
		playbackState.Routing.KeyboardAlbumPath = resolved
		playbackState.Routing.MouseAlbumPath = resolved
		playbackState.Routing.UnifiedSnapshot = snapshot
		playbackState.Routing.KeyboardSnapshot = snapshot
		playbackState.Routing.MouseSnapshot = snapshot
		if playbackState.SourceMode == SourceModeRouteSplit || playbackState.SourceMode == SourceModeRouteUnified {
			playbackState.SourceMode = SourceModeRouteUnified
		}
		playbackStateLock.Unlock()

		return result, nil
	}

	// 分离模式：分别加载键盘/鼠标快照。
	// 说明：键盘与鼠标可以来自不同专辑；若其中一侧失败，另一侧仍可继续工作（返回 partial）。
	keyboardSnapshot, keyboardResolved, keyboardErr := loadSnapshotWithResolve(keyboardPath)
	mouseSnapshot, mouseResolved, mouseErr := loadSnapshotWithResolve(mousePath)

	result.Keyboard = buildStatus(keyboardPath, keyboardResolved, keyboardSnapshot, keyboardErr)
	result.Mouse = buildStatus(mousePath, mouseResolved, mouseSnapshot, mouseErr)
	result.Unified = result.Keyboard

	// 以“整块写入”的方式更新内存态（注意 unified 字段会被清空）。
	playbackStateLock.Lock()
	playbackState.Routing.Mode = "split"
	playbackState.Routing.KeyboardAlbumPath = keyboardResolved
	playbackState.Routing.MouseAlbumPath = mouseResolved
	playbackState.Routing.UnifiedAlbumPath = ""
	playbackState.Routing.KeyboardSnapshot = keyboardSnapshot
	playbackState.Routing.MouseSnapshot = mouseSnapshot
	playbackState.Routing.UnifiedSnapshot = nil
	if playbackState.SourceMode == SourceModeRouteSplit || playbackState.SourceMode == SourceModeRouteUnified {
		playbackState.SourceMode = SourceModeRouteSplit
	}
	playbackStateLock.Unlock()

	// 两侧都失败才认为整体不可用；否则允许 partial（前端可根据 result 诊断）。
	if keyboardErr != nil && mouseErr != nil {
		return result, fmt.Errorf("keyboard and mouse snapshots failed")
	}

	return result, nil
}

func buildStatus(requestedPath, resolvedPath string, snapshot *AlbumSnapshot, err error) RoutingSourceStatus {
	status := RoutingSourceStatus{
		RequestedPath: strings.TrimSpace(requestedPath),
		ResolvedPath:  strings.TrimSpace(resolvedPath),
		Loaded:        snapshot != nil && snapshot.Viper != nil && err == nil,
	}
	if err != nil {
		status.Error = err.Error()
	}
	return status
}

func loadSnapshotWithResolve(input string) (*AlbumSnapshot, string, error) {
	resolved := resolveAlbumPath(input)
	if strings.TrimSpace(resolved) == "" {
		return nil, resolved, errors.New("album path is required")
	}
	// resolved 后的路径是实际磁盘目录（AudioPackagePath/UUID 或用户传入的绝对路径）。
	snapshot, err := loadAlbumSnapshot(resolved)
	return snapshot, resolved, err
}

func resolveAlbumPath(input string) string {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return ""
	}
	// 若包含路径分隔符，则视为完整路径，否则视为 UUID。
	if strings.Contains(trimmed, "/") || strings.Contains(trimmed, "\\") {
		return trimmed
	}
	return filepath.Join(audioPackageConfig.AudioPackagePath, trimmed)
}

func loadAlbumSnapshot(albumPath string) (*AlbumSnapshot, error) {
	// loadAlbumSnapshot 会将专辑配置加载为“只读快照”：
	//   - 使用 viper.New() 创建独立实例（不共享全局 audioPackageConfig.Viper）
	//   - 不启用 WatchConfig / SSE
	//   - 支持三种形态：
	//       1) plain JSON（直接 package.json 内容）
	//       2) legacy-hex（历史 hex 密文，需解密为 JSON）
	//       3) core stub（stub 指向 core 密文文件，需解密为 JSON）
	//   - 返回的 Snapshot 仅用于播放读取（热路径）
	//
	// 尝试解析 core stub（加密专辑）或原始 package.json。
	stubInfo, raw, err := audioPackageConfig.ReadCoreStubInfo(albumPath)
	if err != nil {
		return nil, err
	}
	albumUUID := filepath.Base(albumPath)
	var plainJSON string

	if stubInfo != nil {
		// core stub：读取 core 密文并解密。
		corePath := filepath.Join(albumPath, stubInfo.Core)
		cipherBytes, readErr := os.ReadFile(corePath)
		if readErr != nil {
			return nil, readErr
		}
		plainJSON, err = enc.DecryptConfigBytes(cipherBytes, albumUUID)
		if err != nil {
			return nil, err
		}
	} else {
		if raw == nil {
			return nil, errors.New("package.json not found")
		}
		trimmed := strings.TrimSpace(string(raw))
		if enc.IsLikelyHexCipher(raw) {
			// legacy-hex：历史版本将密文直接写在 package.json 内。
			plainJSON, err = enc.DecryptConfigHex(trimmed, albumUUID)
			if err != nil {
				return nil, err
			}
		} else {
			// plain JSON：无需解密。
			plainJSON = trimmed
		}
	}

	// 快速 JSON 校验：尽早给出可诊断错误。
	if err := enc.ValidateJSONFast(plainJSON); err != nil {
		return nil, err
	}

	// 使用独立 viper 实例读取 JSON，保证播放侧快照不与编辑器（可写配置）互相影响。
	v := viper.New()
	v.SetConfigType("json")
	if err := v.ReadConfig(bytes.NewBufferString(plainJSON)); err != nil {
		return nil, err
	}

	return &AlbumSnapshot{
		AlbumPath: albumPath,
		AlbumUUID: albumUUID,
		Viper:     v,
	}, nil
}
