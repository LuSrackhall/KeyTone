package keySound

import (
	"errors"
	"testing"

	"github.com/gopxl/beep/v2"
)

// fakeStreamSeekCloser 是一个用于测试的假实现，模拟了一个可流式读取、可定位并可关闭的音频流。
// 它以双声道 ([2]float64) 的样本切片作为底层数据，支持读取、查询长度、定位（Seek）、返回错误以及关闭操作。
// 该类型用于在单元测试中替代真实的 beep.StreamSeekCloser，便于精确控制流的行为（如在特定位置注入 seek 错误）。
type fakeStreamSeekCloser struct {
	// samples 保存整个音频流的样本，按帧顺序存储，每帧为左右两个声道的 float64 值。
	samples    [][2]float64
	// position 表示下一个读取样本在 samples 中的索引位置（以帧为单位）。
	position   int
	// closed 表示流是否已被 Close() 调用标记为已关闭。
	closed     bool
	// failOnSeek 可选地在特定的 Seek 位置返回错误。键是目标位置，值是要返回的错误。
	failOnSeek map[int]error
	// err 模拟 Stream 方法可能返回的错误（通过 Err() 暴露）。
	err        error
}

// newFakeStreamSeekCloser 根据单通道数值切片构造一个 fakeStreamSeekCloser。
// 为简化测试，输入的每个 float64 值被复制到左右两个声道中，生成双声道帧。
func newFakeStreamSeekCloser(values []float64) *fakeStreamSeekCloser {
	samples := make([][2]float64, len(values))
	for index, value := range values {
		samples[index] = [2]float64{value, value}
	}
	return &fakeStreamSeekCloser{samples: samples}
}

// Stream 实现了 beep.Streamer 接口的 Stream 方法。
// 它将从当前 position 开始的样本复制到提供的 dst 缓冲区，返回复制的帧数与是否还有剩余样本。
// 如果流已关闭或已到达末尾，则返回 0, false。
func (f *fakeStreamSeekCloser) Stream(dst [][2]float64) (int, bool) {
	if f.closed || f.position >= len(f.samples) {
		return 0, false
	}

	count := copy(dst, f.samples[f.position:])
	f.position += count
	return count, f.position < len(f.samples)
}

// Err 返回 Stream 期间可能发生的错误（如果有）。这是 beep.StreamSeekCloser 的常见约定。
func (f *fakeStreamSeekCloser) Err() error {
	return f.err
}

// Len 返回流中总的帧数（每帧包含两个声道）。
func (f *fakeStreamSeekCloser) Len() int {
	return len(f.samples)
}

// Position 返回当前流的位置（下一次读取将从该索引开始）。
func (f *fakeStreamSeekCloser) Position() int {
	return f.position
}

// Seek 将流的位置移动到指定的帧索引。
// 如果在 failOnSeek 中注册了该位置，会返回对应错误；否则进行边界检查并设置 position。
func (f *fakeStreamSeekCloser) Seek(position int) error {
	if err, ok := f.failOnSeek[position]; ok {
		return err
	}
	if position < 0 || position > len(f.samples) {
		return errors.New("seek out of range")
	}
	f.position = position
	return nil
}

// Close 标记流为已关闭并返回 nil（模拟可关闭的资源）。
func (f *fakeStreamSeekCloser) Close() error {
	f.closed = true
	return nil
}

// collectLeftChannel 从给定的 beep.Streamer 中读取所有样本，并返回左声道（索引 0）的样本序列。
// 该辅助函数用于测试中验证生成片段的样本内容是否符合预期。
func collectLeftChannel(streamer beep.Streamer) []float64 {
	buffer := make([][2]float64, 2)
	values := make([]float64, 0)
	for {
		count, ok := streamer.Stream(buffer)
		for index := 0; index < count; index++ {
			values = append(values, buffer[index][0])
		}
		if !ok {
			break
		}
	}
	return values
}

// TestPreparePlaybackSourceSlicesRequestedRange 验证 preparePlaybackSource 在给定起止毫秒范围时：
// - 正确从底层流中裁切出对应样本（按采样率换算毫秒到帧）
// - 返回的初始音量等于 Cut.Volume
// - 底层流的位置被移动到切片末端
func TestPreparePlaybackSourceSlicesRequestedRange(t *testing.T) {
	stream := newFakeStreamSeekCloser([]float64{0, 1, 2, 3, 4, 5, 6})
	segment, initVolume, err := preparePlaybackSource(stream, beep.SampleRate(1000), &Cut{
		StartMS: 2,
		EndMS:   5,
		Volume:  0.25,
	})
	if err != nil {
		t.Fatalf("preparePlaybackSource returned error: %v", err)
	}
	if initVolume != 0.25 {
		t.Fatalf("unexpected init volume: got %v want %v", initVolume, 0.25)
	}

	values := collectLeftChannel(segment)
	if len(values) != 3 {
		t.Fatalf("unexpected sample count: got %d want %d", len(values), 3)
	}
	for index, expected := range []float64{2, 3, 4} {
		if values[index] != expected {
			t.Fatalf("unexpected sample at %d: got %v want %v", index, values[index], expected)
		}
	}
	if stream.Position() != 5 {
		t.Fatalf("unexpected final position: got %d want %d", stream.Position(), 5)
	}
}

// TestPreparePlaybackSourceClampsEndToLength 验证当 Cut.EndMS 超过流长度时，preparePlaybackSource 会将结束位置限制到流的末尾，避免越界。
func TestPreparePlaybackSourceClampsEndToLength(t *testing.T) {
	stream := newFakeStreamSeekCloser([]float64{0, 1, 2, 3, 4, 5, 6})
	segment, _, err := preparePlaybackSource(stream, beep.SampleRate(1000), &Cut{
		StartMS: 5,
		EndMS:   20,
	})
	if err != nil {
		t.Fatalf("preparePlaybackSource returned error: %v", err)
	}

	values := collectLeftChannel(segment)
	if len(values) != 2 {
		t.Fatalf("unexpected sample count: got %d want %d", len(values), 2)
	}
	for index, expected := range []float64{5, 6} {
		if values[index] != expected {
			t.Fatalf("unexpected sample at %d: got %v want %v", index, values[index], expected)
		}
	}
}

// TestPreparePlaybackSourceRejectsEmptyOrOutOfRangeCut 验证当裁切范围完全在流之外或长度为零时，preparePlaybackSource 返回 errEmptyAudioCut。
func TestPreparePlaybackSourceRejectsEmptyOrOutOfRangeCut(t *testing.T) {
	stream := newFakeStreamSeekCloser([]float64{0, 1, 2, 3, 4})
	_, _, err := preparePlaybackSource(stream, beep.SampleRate(1000), &Cut{
		StartMS: 10,
		EndMS:   20,
	})
	if !errors.Is(err, errEmptyAudioCut) {
		t.Fatalf("unexpected error: got %v want %v", err, errEmptyAudioCut)
	}

	_, _, err = preparePlaybackSource(stream, beep.SampleRate(1000), &Cut{
		StartMS: 3,
		EndMS:   3,
	})
	if !errors.Is(err, errEmptyAudioCut) {
		t.Fatalf("unexpected error for zero-length cut: got %v want %v", err, errEmptyAudioCut)
	}
}

// TestPreparePlaybackSourceReturnsSeekError 验证当底层流的 Seek 操作返回错误时，preparePlaybackSource 会将该错误向上传递而不是将其视为空裁切。
func TestPreparePlaybackSourceReturnsSeekError(t *testing.T) {
	stream := newFakeStreamSeekCloser([]float64{0, 1, 2, 3, 4})
	stream.failOnSeek = map[int]error{2: errors.New("seek failed")}

	_, _, err := preparePlaybackSource(stream, beep.SampleRate(1000), &Cut{
		StartMS: 2,
		EndMS:   4,
	})
	if err == nil {
		t.Fatal("expected seek error, got nil")
	}
	if errors.Is(err, errEmptyAudioCut) {
		t.Fatalf("expected seek error, got empty cut error: %v", err)
	}
}
