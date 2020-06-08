// Code generated by protoc-gen-go. DO NOT EDIT.
// source: types.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type ImageFilter int32

const (
	// DEFAULT represents the default filter.
	ImageFilter_DEFAULT ImageFilter = 0
	// NEAREST represents nearest (crisp-edged) filter
	ImageFilter_NEAREST ImageFilter = 1
	// LINEAR represents linear filter
	ImageFilter_LINEAR ImageFilter = 2
)

var ImageFilter_name = map[int32]string{
	0: "DEFAULT",
	1: "NEAREST",
	2: "LINEAR",
}

var ImageFilter_value = map[string]int32{
	"DEFAULT": 0,
	"NEAREST": 1,
	"LINEAR":  2,
}

func (x ImageFilter) String() string {
	return proto.EnumName(ImageFilter_name, int32(x))
}

func (ImageFilter) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{0}
}

type AnimationClipMode int32

const (
	AnimationClipMode_ONCE          AnimationClipMode = 0
	AnimationClipMode_LOOP          AnimationClipMode = 1
	AnimationClipMode_PING_PONG     AnimationClipMode = 2
	AnimationClipMode_CLAMP_FOREVER AnimationClipMode = 4
)

var AnimationClipMode_name = map[int32]string{
	0: "ONCE",
	1: "LOOP",
	2: "PING_PONG",
	4: "CLAMP_FOREVER",
}

var AnimationClipMode_value = map[string]int32{
	"ONCE":          0,
	"LOOP":          1,
	"PING_PONG":     2,
	"CLAMP_FOREVER": 4,
}

func (x AnimationClipMode) String() string {
	return proto.EnumName(AnimationClipMode_name, int32(x))
}

func (AnimationClipMode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{1}
}

type AtlasFile struct {
	Images               [][]byte                  `protobuf:"bytes,1,rep,name=images,proto3" json:"images,omitempty"`
	Filters              []ImageFilter             `protobuf:"varint,2,rep,packed,name=filters,proto3,enum=pb.ImageFilter" json:"filters,omitempty"`
	Frames               map[string]*Frame         `protobuf:"bytes,3,rep,name=frames,proto3" json:"frames,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Clips                map[string]*AnimationClip `protobuf:"bytes,4,rep,name=clips,proto3" json:"clips,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Animations           map[string]*Animation     `protobuf:"bytes,5,rep,name=animations,proto3" json:"animations,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                  `json:"-"`
	XXX_unrecognized     []byte                    `json:"-"`
	XXX_sizecache        int32                     `json:"-"`
}

func (m *AtlasFile) Reset()         { *m = AtlasFile{} }
func (m *AtlasFile) String() string { return proto.CompactTextString(m) }
func (*AtlasFile) ProtoMessage()    {}
func (*AtlasFile) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{0}
}

func (m *AtlasFile) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AtlasFile.Unmarshal(m, b)
}
func (m *AtlasFile) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AtlasFile.Marshal(b, m, deterministic)
}
func (m *AtlasFile) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AtlasFile.Merge(m, src)
}
func (m *AtlasFile) XXX_Size() int {
	return xxx_messageInfo_AtlasFile.Size(m)
}
func (m *AtlasFile) XXX_DiscardUnknown() {
	xxx_messageInfo_AtlasFile.DiscardUnknown(m)
}

var xxx_messageInfo_AtlasFile proto.InternalMessageInfo

func (m *AtlasFile) GetImages() [][]byte {
	if m != nil {
		return m.Images
	}
	return nil
}

func (m *AtlasFile) GetFilters() []ImageFilter {
	if m != nil {
		return m.Filters
	}
	return nil
}

func (m *AtlasFile) GetFrames() map[string]*Frame {
	if m != nil {
		return m.Frames
	}
	return nil
}

func (m *AtlasFile) GetClips() map[string]*AnimationClip {
	if m != nil {
		return m.Clips
	}
	return nil
}

func (m *AtlasFile) GetAnimations() map[string]*Animation {
	if m != nil {
		return m.Animations
	}
	return nil
}

type Frame struct {
	Image                uint32   `protobuf:"varint,1,opt,name=image,proto3" json:"image,omitempty"`
	X                    uint32   `protobuf:"varint,2,opt,name=x,proto3" json:"x,omitempty"`
	Y                    uint32   `protobuf:"varint,3,opt,name=y,proto3" json:"y,omitempty"`
	W                    uint32   `protobuf:"varint,4,opt,name=w,proto3" json:"w,omitempty"`
	H                    uint32   `protobuf:"varint,5,opt,name=h,proto3" json:"h,omitempty"`
	Ox                   int32    `protobuf:"varint,6,opt,name=ox,proto3" json:"ox,omitempty"`
	Oy                   int32    `protobuf:"varint,7,opt,name=oy,proto3" json:"oy,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Frame) Reset()         { *m = Frame{} }
func (m *Frame) String() string { return proto.CompactTextString(m) }
func (*Frame) ProtoMessage()    {}
func (*Frame) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{1}
}

func (m *Frame) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Frame.Unmarshal(m, b)
}
func (m *Frame) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Frame.Marshal(b, m, deterministic)
}
func (m *Frame) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Frame.Merge(m, src)
}
func (m *Frame) XXX_Size() int {
	return xxx_messageInfo_Frame.Size(m)
}
func (m *Frame) XXX_DiscardUnknown() {
	xxx_messageInfo_Frame.DiscardUnknown(m)
}

var xxx_messageInfo_Frame proto.InternalMessageInfo

func (m *Frame) GetImage() uint32 {
	if m != nil {
		return m.Image
	}
	return 0
}

func (m *Frame) GetX() uint32 {
	if m != nil {
		return m.X
	}
	return 0
}

func (m *Frame) GetY() uint32 {
	if m != nil {
		return m.Y
	}
	return 0
}

func (m *Frame) GetW() uint32 {
	if m != nil {
		return m.W
	}
	return 0
}

func (m *Frame) GetH() uint32 {
	if m != nil {
		return m.H
	}
	return 0
}

func (m *Frame) GetOx() int32 {
	if m != nil {
		return m.Ox
	}
	return 0
}

func (m *Frame) GetOy() int32 {
	if m != nil {
		return m.Oy
	}
	return 0
}

type Animation struct {
	Name                 string           `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Clips                []*AnimationClip `protobuf:"bytes,2,rep,name=clips,proto3" json:"clips,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Animation) Reset()         { *m = Animation{} }
func (m *Animation) String() string { return proto.CompactTextString(m) }
func (*Animation) ProtoMessage()    {}
func (*Animation) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{2}
}

func (m *Animation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Animation.Unmarshal(m, b)
}
func (m *Animation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Animation.Marshal(b, m, deterministic)
}
func (m *Animation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Animation.Merge(m, src)
}
func (m *Animation) XXX_Size() int {
	return xxx_messageInfo_Animation.Size(m)
}
func (m *Animation) XXX_DiscardUnknown() {
	xxx_messageInfo_Animation.DiscardUnknown(m)
}

var xxx_messageInfo_Animation proto.InternalMessageInfo

func (m *Animation) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Animation) GetClips() []*AnimationClip {
	if m != nil {
		return m.Clips
	}
	return nil
}

type AnimationClip struct {
	Name                 string            `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Fps                  float32           `protobuf:"fixed32,2,opt,name=fps,proto3" json:"fps,omitempty"`
	ClipMode             AnimationClipMode `protobuf:"varint,3,opt,name=clip_mode,json=clipMode,proto3,enum=pb.AnimationClipMode" json:"clip_mode,omitempty"`
	Frames               []*AnimFrame      `protobuf:"bytes,4,rep,name=frames,proto3" json:"frames,omitempty"`
	EndedEvent           *AnimationEvent   `protobuf:"bytes,5,opt,name=ended_event,json=endedEvent,proto3" json:"ended_event,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *AnimationClip) Reset()         { *m = AnimationClip{} }
func (m *AnimationClip) String() string { return proto.CompactTextString(m) }
func (*AnimationClip) ProtoMessage()    {}
func (*AnimationClip) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{3}
}

func (m *AnimationClip) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AnimationClip.Unmarshal(m, b)
}
func (m *AnimationClip) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AnimationClip.Marshal(b, m, deterministic)
}
func (m *AnimationClip) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AnimationClip.Merge(m, src)
}
func (m *AnimationClip) XXX_Size() int {
	return xxx_messageInfo_AnimationClip.Size(m)
}
func (m *AnimationClip) XXX_DiscardUnknown() {
	xxx_messageInfo_AnimationClip.DiscardUnknown(m)
}

var xxx_messageInfo_AnimationClip proto.InternalMessageInfo

func (m *AnimationClip) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AnimationClip) GetFps() float32 {
	if m != nil {
		return m.Fps
	}
	return 0
}

func (m *AnimationClip) GetClipMode() AnimationClipMode {
	if m != nil {
		return m.ClipMode
	}
	return AnimationClipMode_ONCE
}

func (m *AnimationClip) GetFrames() []*AnimFrame {
	if m != nil {
		return m.Frames
	}
	return nil
}

func (m *AnimationClip) GetEndedEvent() *AnimationEvent {
	if m != nil {
		return m.EndedEvent
	}
	return nil
}

type AnimFrame struct {
	FrameName            string          `protobuf:"bytes,1,opt,name=frame_name,json=frameName,proto3" json:"frame_name,omitempty"`
	Event                *AnimationEvent `protobuf:"bytes,2,opt,name=event,proto3" json:"event,omitempty"`
	XXX_NoUnkeyedLiteral struct{}        `json:"-"`
	XXX_unrecognized     []byte          `json:"-"`
	XXX_sizecache        int32           `json:"-"`
}

func (m *AnimFrame) Reset()         { *m = AnimFrame{} }
func (m *AnimFrame) String() string { return proto.CompactTextString(m) }
func (*AnimFrame) ProtoMessage()    {}
func (*AnimFrame) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{4}
}

func (m *AnimFrame) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AnimFrame.Unmarshal(m, b)
}
func (m *AnimFrame) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AnimFrame.Marshal(b, m, deterministic)
}
func (m *AnimFrame) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AnimFrame.Merge(m, src)
}
func (m *AnimFrame) XXX_Size() int {
	return xxx_messageInfo_AnimFrame.Size(m)
}
func (m *AnimFrame) XXX_DiscardUnknown() {
	xxx_messageInfo_AnimFrame.DiscardUnknown(m)
}

var xxx_messageInfo_AnimFrame proto.InternalMessageInfo

func (m *AnimFrame) GetFrameName() string {
	if m != nil {
		return m.FrameName
	}
	return ""
}

func (m *AnimFrame) GetEvent() *AnimationEvent {
	if m != nil {
		return m.Event
	}
	return nil
}

type AnimationEvent struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AnimationEvent) Reset()         { *m = AnimationEvent{} }
func (m *AnimationEvent) String() string { return proto.CompactTextString(m) }
func (*AnimationEvent) ProtoMessage()    {}
func (*AnimationEvent) Descriptor() ([]byte, []int) {
	return fileDescriptor_d938547f84707355, []int{5}
}

func (m *AnimationEvent) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AnimationEvent.Unmarshal(m, b)
}
func (m *AnimationEvent) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AnimationEvent.Marshal(b, m, deterministic)
}
func (m *AnimationEvent) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AnimationEvent.Merge(m, src)
}
func (m *AnimationEvent) XXX_Size() int {
	return xxx_messageInfo_AnimationEvent.Size(m)
}
func (m *AnimationEvent) XXX_DiscardUnknown() {
	xxx_messageInfo_AnimationEvent.DiscardUnknown(m)
}

var xxx_messageInfo_AnimationEvent proto.InternalMessageInfo

func (m *AnimationEvent) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AnimationEvent) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterEnum("pb.ImageFilter", ImageFilter_name, ImageFilter_value)
	proto.RegisterEnum("pb.AnimationClipMode", AnimationClipMode_name, AnimationClipMode_value)
	proto.RegisterType((*AtlasFile)(nil), "pb.AtlasFile")
	proto.RegisterMapType((map[string]*Animation)(nil), "pb.AtlasFile.AnimationsEntry")
	proto.RegisterMapType((map[string]*AnimationClip)(nil), "pb.AtlasFile.ClipsEntry")
	proto.RegisterMapType((map[string]*Frame)(nil), "pb.AtlasFile.FramesEntry")
	proto.RegisterType((*Frame)(nil), "pb.Frame")
	proto.RegisterType((*Animation)(nil), "pb.Animation")
	proto.RegisterType((*AnimationClip)(nil), "pb.AnimationClip")
	proto.RegisterType((*AnimFrame)(nil), "pb.AnimFrame")
	proto.RegisterType((*AnimationEvent)(nil), "pb.AnimationEvent")
}

func init() { proto.RegisterFile("types.proto", fileDescriptor_d938547f84707355) }

var fileDescriptor_d938547f84707355 = []byte{
	// 573 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x54, 0x51, 0x6f, 0xda, 0x30,
	0x10, 0xae, 0x1d, 0x42, 0x9b, 0x4b, 0x43, 0x53, 0xab, 0x9b, 0xbc, 0x4a, 0xd5, 0xa2, 0x4c, 0xd3,
	0xb2, 0x3e, 0x20, 0x2d, 0xbc, 0x4c, 0x95, 0xf6, 0x80, 0x68, 0xe8, 0xd0, 0x28, 0x20, 0x8f, 0xed,
	0x15, 0x85, 0x62, 0xd6, 0x68, 0x90, 0x20, 0x92, 0xb6, 0xe4, 0x07, 0xee, 0x07, 0xec, 0x1f, 0x4d,
	0xb6, 0x03, 0x4b, 0x56, 0xf6, 0x76, 0x9f, 0xbf, 0xfb, 0x3e, 0x9f, 0xef, 0x2e, 0x01, 0x33, 0xcb,
	0x57, 0x3c, 0x6d, 0xae, 0xd6, 0x49, 0x96, 0x10, 0xbc, 0x9a, 0xba, 0xbf, 0x35, 0x30, 0xda, 0xd9,
	0x22, 0x4c, 0xbb, 0xd1, 0x82, 0x93, 0x97, 0x50, 0x8f, 0x96, 0xe1, 0x0f, 0x9e, 0x52, 0xe4, 0x68,
	0xde, 0x31, 0x2b, 0x10, 0x79, 0x0f, 0x87, 0xf3, 0x68, 0x91, 0xf1, 0x75, 0x4a, 0xb1, 0xa3, 0x79,
	0x0d, 0xff, 0xa4, 0xb9, 0x9a, 0x36, 0x7b, 0x82, 0xec, 0xca, 0x73, 0xb6, 0xe5, 0xc9, 0x07, 0xa8,
	0xcf, 0xd7, 0xe1, 0x92, 0xa7, 0x54, 0x73, 0x34, 0xcf, 0xf4, 0x5f, 0x89, 0xcc, 0xdd, 0x0d, 0xcd,
	0xae, 0xe4, 0x82, 0x38, 0x5b, 0xe7, 0xac, 0x48, 0x24, 0x4d, 0xd0, 0xef, 0x16, 0xd1, 0x2a, 0xa5,
	0x35, 0xa9, 0xa0, 0x55, 0x45, 0x47, 0x50, 0x4a, 0xa0, 0xd2, 0xc8, 0x27, 0x80, 0x30, 0x8e, 0x96,
	0x61, 0x16, 0x25, 0x71, 0x4a, 0x75, 0x29, 0xba, 0xa8, 0x8a, 0xda, 0x3b, 0x5e, 0x29, 0x4b, 0x82,
	0xf3, 0x6b, 0x30, 0x4b, 0x55, 0x10, 0x1b, 0xb4, 0x9f, 0x3c, 0xa7, 0xc8, 0x41, 0x9e, 0xc1, 0x44,
	0x48, 0x5e, 0x83, 0xfe, 0x18, 0x2e, 0x1e, 0x38, 0xc5, 0x0e, 0xf2, 0x4c, 0xdf, 0x10, 0xd6, 0x52,
	0xc1, 0xd4, 0xf9, 0x15, 0xfe, 0x88, 0xce, 0xbf, 0x00, 0xfc, 0xad, 0x6c, 0x8f, 0xc9, 0xbb, 0xaa,
	0xc9, 0xa9, 0xac, 0x6f, 0x5b, 0x84, 0x50, 0x96, 0xcd, 0xfa, 0x70, 0xf2, 0x4f, 0xc5, 0x7b, 0x1c,
	0xdf, 0x54, 0x1d, 0xad, 0x8a, 0x63, 0xc9, 0xcd, 0x7d, 0x00, 0x5d, 0x96, 0x4b, 0xce, 0x40, 0x97,
	0x03, 0x94, 0x2e, 0x16, 0x53, 0x80, 0x1c, 0x03, 0xda, 0x48, 0x0f, 0x8b, 0xa1, 0x8d, 0x40, 0x39,
	0xd5, 0x14, 0xca, 0x05, 0x7a, 0xa2, 0x35, 0x85, 0x9e, 0x04, 0xba, 0xa7, 0xba, 0x42, 0xf7, 0xa4,
	0x01, 0x38, 0xd9, 0xd0, 0xba, 0x83, 0x3c, 0x9d, 0xe1, 0x64, 0x23, 0x71, 0x4e, 0x0f, 0x0b, 0x9c,
	0xbb, 0x9f, 0xc1, 0xd8, 0x95, 0x43, 0x08, 0xd4, 0xe2, 0x70, 0xc9, 0x8b, 0xfa, 0x65, 0x2c, 0x5a,
	0xa2, 0xe6, 0x8c, 0xe5, 0xc8, 0xf6, 0xb5, 0x44, 0xf2, 0xee, 0x2f, 0x04, 0x56, 0x85, 0xd8, 0x6b,
	0x67, 0x83, 0x36, 0x97, 0x66, 0xc8, 0xc3, 0x4c, 0x84, 0xc4, 0x07, 0x43, 0x18, 0x4c, 0x96, 0xc9,
	0x8c, 0xcb, 0x37, 0x35, 0xfc, 0x17, 0xcf, 0x2e, 0xb9, 0x4d, 0x66, 0x9c, 0x1d, 0xdd, 0x15, 0x11,
	0x79, 0xbb, 0xdb, 0x57, 0xb5, 0x7d, 0xbb, 0xb6, 0xaa, 0x89, 0x6f, 0x77, 0xb4, 0x05, 0x26, 0x8f,
	0x67, 0x7c, 0x36, 0xe1, 0x8f, 0x3c, 0xce, 0x64, 0x53, 0x4c, 0x9f, 0x54, 0xcc, 0x03, 0xc1, 0x30,
	0x90, 0x69, 0x32, 0x76, 0xc7, 0xaa, 0x23, 0x6a, 0x18, 0x17, 0x00, 0xd2, 0x6b, 0x52, 0x7a, 0x88,
	0x21, 0x4f, 0x06, 0x82, 0xf6, 0x40, 0x57, 0xd6, 0xf8, 0xbf, 0xd6, 0x2a, 0xc1, 0xbd, 0x82, 0x46,
	0x95, 0xd8, 0xdb, 0x9d, 0xb3, 0xf2, 0xb6, 0x18, 0xc5, 0x7a, 0x5c, 0xb6, 0xc0, 0x2c, 0x7d, 0xb5,
	0xc4, 0x84, 0xc3, 0xeb, 0xa0, 0xdb, 0xfe, 0xd6, 0x1f, 0xdb, 0x07, 0x02, 0x0c, 0x82, 0x36, 0x0b,
	0xbe, 0x8e, 0x6d, 0x44, 0x00, 0xea, 0xfd, 0x9e, 0x80, 0x36, 0xbe, 0xec, 0xc1, 0xe9, 0xb3, 0x0e,
	0x92, 0x23, 0xa8, 0x0d, 0x07, 0x9d, 0xc0, 0x3e, 0x10, 0x51, 0x7f, 0x38, 0x1c, 0xd9, 0x88, 0x58,
	0x60, 0x8c, 0x7a, 0x83, 0x9b, 0xc9, 0x68, 0x38, 0xb8, 0xb1, 0x31, 0x39, 0x05, 0xab, 0xd3, 0x6f,
	0xdf, 0x8e, 0x26, 0xdd, 0x21, 0x0b, 0xbe, 0x07, 0xcc, 0xae, 0x4d, 0xeb, 0xf2, 0xcf, 0xd3, 0xfa,
	0x13, 0x00, 0x00, 0xff, 0xff, 0xbb, 0x47, 0xf5, 0x1b, 0x88, 0x04, 0x00, 0x00,
}
