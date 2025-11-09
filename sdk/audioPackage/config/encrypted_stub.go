package audioPackageConfig

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Core 文件与指示 JSON 的常量定义
const (
	CoreFileName      = "core"
	coreSchemaVersion = "kcfg-core-v1"
)

// coreStubMetadata 描述 `package.json` 在加密模式下的最小指示信息
type coreStubMetadata struct {
	Encrypted bool   `json:"_keytone_encrypted"`
	Core      string `json:"_keytone_core"`
	Schema    string `json:"_keytone_schema,omitempty"`
	UpdatedAt string `json:"_keytone_updated_at,omitempty"`
}

// newCoreStub 生成默认的指示 JSON 元数据
func newCoreStub() coreStubMetadata {
	return coreStubMetadata{
		Encrypted: true,
		Core:      CoreFileName,
		Schema:    coreSchemaVersion,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// parseCoreStub 解析字节内容，若为指示 JSON 则返回元数据
func parseCoreStub(data []byte) (*coreStubMetadata, bool, error) {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" || !strings.HasPrefix(trimmed, "{") {
		return nil, false, nil
	}

	var stub coreStubMetadata
	if err := json.Unmarshal(data, &stub); err != nil {
		return nil, false, nil
	}

	if !stub.Encrypted || strings.TrimSpace(stub.Core) == "" {
		return nil, false, nil
	}

	if stub.Schema != "" && stub.Schema != coreSchemaVersion {
		return &stub, true, errors.New("unsupported core schema version")
	}

	if filepath.IsAbs(stub.Core) || strings.Contains(stub.Core, "..") {
		return &stub, true, errors.New("invalid core filename")
	}

	if stub.UpdatedAt == "" {
		stub.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	return &stub, true, nil
}

// readCoreStub 从指定 albumPath 下读取 package.json，返回元数据（若存在）及原始字节
func readCoreStub(albumPath string) (*coreStubMetadata, []byte, error) {
	pkgPath := filepath.Join(albumPath, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	stub, ok, parseErr := parseCoreStub(data)
	if !ok {
		return nil, data, parseErr
	}
	return stub, data, parseErr
}

// writeCoreStub 将指示 JSON 写回 albumPath/package.json，使用原子写入
func writeCoreStub(albumPath string, stub *coreStubMetadata) error {
	pkgPath := filepath.Join(albumPath, "package.json")
	meta := stub
	if meta == nil {
		defaultStub := newCoreStub()
		meta = &defaultStub
	}
	meta.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	payload, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}

	tmp := pkgPath + ".tmp"
	if err := os.WriteFile(tmp, payload, 0644); err != nil {
		return err
	}
	if err := os.Rename(tmp, pkgPath); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

// CoreStubInfo 暴露给其他包的指示 JSON 元数据
type CoreStubInfo struct {
	Core      string
	Schema    string
	UpdatedAt string
}

// ReadCoreStubInfo 读取并返回指示 JSON 元数据
func ReadCoreStubInfo(albumPath string) (*CoreStubInfo, []byte, error) {
	stub, raw, err := readCoreStub(albumPath)
	if err != nil {
		return nil, raw, err
	}
	if stub == nil {
		return nil, raw, nil
	}
	info := &CoreStubInfo{
		Core:      stub.Core,
		Schema:    stub.Schema,
		UpdatedAt: stub.UpdatedAt,
	}
	return info, raw, nil
}

// WriteCoreStubFile 写入默认的指示 JSON
func WriteCoreStubFile(albumPath string) error {
	return writeCoreStub(albumPath, nil)
}
