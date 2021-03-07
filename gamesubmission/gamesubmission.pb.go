// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.15.4
// source: gamesubmission/gamesubmission.proto

package gamesubmission

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type SubmitGameRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Version              string       `protobuf:"bytes,1,opt,name=version,proto3" json:"version,omitempty"`
	PlayerID             int32        `protobuf:"varint,2,opt,name=playerID,proto3" json:"playerID,omitempty"`
	PlayerName           string       `protobuf:"bytes,3,opt,name=playerName,proto3" json:"playerName,omitempty"`
	LevelHashMD5         string       `protobuf:"bytes,4,opt,name=levelHashMD5,proto3" json:"levelHashMD5,omitempty"`
	Time                 float32      `protobuf:"fixed32,5,opt,name=time,proto3" json:"time,omitempty"`
	TimeLvll2            float32      `protobuf:"fixed32,6,opt,name=timeLvll2,proto3" json:"timeLvll2,omitempty"`
	TimeLvll3            float32      `protobuf:"fixed32,7,opt,name=timeLvll3,proto3" json:"timeLvll3,omitempty"`
	TimeLvll4            float32      `protobuf:"fixed32,8,opt,name=timeLvll4,proto3" json:"timeLvll4,omitempty"`
	TimeLeviDown         float32      `protobuf:"fixed32,9,opt,name=timeLeviDown,proto3" json:"timeLeviDown,omitempty"`
	TimeOrbDown          float32      `protobuf:"fixed32,10,opt,name=timeOrbDown,proto3" json:"timeOrbDown,omitempty"`
	GemsCollected        int32        `protobuf:"varint,11,opt,name=gemsCollected,proto3" json:"gemsCollected,omitempty"`
	Kills                int32        `protobuf:"varint,12,opt,name=kills,proto3" json:"kills,omitempty"`
	DaggersFired         int32        `protobuf:"varint,13,opt,name=daggersFired,proto3" json:"daggersFired,omitempty"`
	DaggersHit           int32        `protobuf:"varint,14,opt,name=daggersHit,proto3" json:"daggersHit,omitempty"`
	EnemiesAlive         int32        `protobuf:"varint,15,opt,name=enemiesAlive,proto3" json:"enemiesAlive,omitempty"`
	EnemiesAliveMax      int32        `protobuf:"varint,16,opt,name=enemiesAliveMax,proto3" json:"enemiesAliveMax,omitempty"`
	EnemiesAliveMaxTime  float32      `protobuf:"fixed32,17,opt,name=enemiesAliveMaxTime,proto3" json:"enemiesAliveMaxTime,omitempty"`
	LevelGems            int32        `protobuf:"varint,18,opt,name=levelGems,proto3" json:"levelGems,omitempty"`
	HomingDaggers        int32        `protobuf:"varint,19,opt,name=homingDaggers,proto3" json:"homingDaggers,omitempty"`
	HomingDaggersMax     int32        `protobuf:"varint,20,opt,name=homingDaggersMax,proto3" json:"homingDaggersMax,omitempty"`
	HomingDaggersMaxTime float32      `protobuf:"fixed32,21,opt,name=homingDaggersMaxTime,proto3" json:"homingDaggersMaxTime,omitempty"`
	TotalGems            int32        `protobuf:"varint,22,opt,name=totalGems,proto3" json:"totalGems,omitempty"`
	GemsDespawned        int32        `protobuf:"varint,23,opt,name=gemsDespawned,proto3" json:"gemsDespawned,omitempty"`
	GemsEaten            int32        `protobuf:"varint,24,opt,name=gemsEaten,proto3" json:"gemsEaten,omitempty"`
	DaggersEaten         int32        `protobuf:"varint,25,opt,name=daggersEaten,proto3" json:"daggersEaten,omitempty"`
	DeathType            uint32       `protobuf:"varint,26,opt,name=deathType,proto3" json:"deathType,omitempty"`
	IsReplay             bool         `protobuf:"varint,27,opt,name=isReplay,proto3" json:"isReplay,omitempty"`
	ReplayPlayerID       int32        `protobuf:"varint,28,opt,name=replayPlayerID,proto3" json:"replayPlayerID,omitempty"`
	PerEnemyAliveCount   []int32      `protobuf:"varint,29,rep,packed,name=perEnemyAliveCount,proto3" json:"perEnemyAliveCount,omitempty"`
	PerEnemyKillcount    []int32      `protobuf:"varint,30,rep,packed,name=perEnemyKillcount,proto3" json:"perEnemyKillcount,omitempty"`
	Stats                []*StatFrame `protobuf:"bytes,31,rep,name=stats,proto3" json:"stats,omitempty"`
}

func (x *SubmitGameRequest) Reset() {
	*x = SubmitGameRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gamesubmission_gamesubmission_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitGameRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitGameRequest) ProtoMessage() {}

func (x *SubmitGameRequest) ProtoReflect() protoreflect.Message {
	mi := &file_gamesubmission_gamesubmission_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitGameRequest.ProtoReflect.Descriptor instead.
func (*SubmitGameRequest) Descriptor() ([]byte, []int) {
	return file_gamesubmission_gamesubmission_proto_rawDescGZIP(), []int{0}
}

func (x *SubmitGameRequest) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *SubmitGameRequest) GetPlayerID() int32 {
	if x != nil {
		return x.PlayerID
	}
	return 0
}

func (x *SubmitGameRequest) GetPlayerName() string {
	if x != nil {
		return x.PlayerName
	}
	return ""
}

func (x *SubmitGameRequest) GetLevelHashMD5() string {
	if x != nil {
		return x.LevelHashMD5
	}
	return ""
}

func (x *SubmitGameRequest) GetTime() float32 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *SubmitGameRequest) GetTimeLvll2() float32 {
	if x != nil {
		return x.TimeLvll2
	}
	return 0
}

func (x *SubmitGameRequest) GetTimeLvll3() float32 {
	if x != nil {
		return x.TimeLvll3
	}
	return 0
}

func (x *SubmitGameRequest) GetTimeLvll4() float32 {
	if x != nil {
		return x.TimeLvll4
	}
	return 0
}

func (x *SubmitGameRequest) GetTimeLeviDown() float32 {
	if x != nil {
		return x.TimeLeviDown
	}
	return 0
}

func (x *SubmitGameRequest) GetTimeOrbDown() float32 {
	if x != nil {
		return x.TimeOrbDown
	}
	return 0
}

func (x *SubmitGameRequest) GetGemsCollected() int32 {
	if x != nil {
		return x.GemsCollected
	}
	return 0
}

func (x *SubmitGameRequest) GetKills() int32 {
	if x != nil {
		return x.Kills
	}
	return 0
}

func (x *SubmitGameRequest) GetDaggersFired() int32 {
	if x != nil {
		return x.DaggersFired
	}
	return 0
}

func (x *SubmitGameRequest) GetDaggersHit() int32 {
	if x != nil {
		return x.DaggersHit
	}
	return 0
}

func (x *SubmitGameRequest) GetEnemiesAlive() int32 {
	if x != nil {
		return x.EnemiesAlive
	}
	return 0
}

func (x *SubmitGameRequest) GetEnemiesAliveMax() int32 {
	if x != nil {
		return x.EnemiesAliveMax
	}
	return 0
}

func (x *SubmitGameRequest) GetEnemiesAliveMaxTime() float32 {
	if x != nil {
		return x.EnemiesAliveMaxTime
	}
	return 0
}

func (x *SubmitGameRequest) GetLevelGems() int32 {
	if x != nil {
		return x.LevelGems
	}
	return 0
}

func (x *SubmitGameRequest) GetHomingDaggers() int32 {
	if x != nil {
		return x.HomingDaggers
	}
	return 0
}

func (x *SubmitGameRequest) GetHomingDaggersMax() int32 {
	if x != nil {
		return x.HomingDaggersMax
	}
	return 0
}

func (x *SubmitGameRequest) GetHomingDaggersMaxTime() float32 {
	if x != nil {
		return x.HomingDaggersMaxTime
	}
	return 0
}

func (x *SubmitGameRequest) GetTotalGems() int32 {
	if x != nil {
		return x.TotalGems
	}
	return 0
}

func (x *SubmitGameRequest) GetGemsDespawned() int32 {
	if x != nil {
		return x.GemsDespawned
	}
	return 0
}

func (x *SubmitGameRequest) GetGemsEaten() int32 {
	if x != nil {
		return x.GemsEaten
	}
	return 0
}

func (x *SubmitGameRequest) GetDaggersEaten() int32 {
	if x != nil {
		return x.DaggersEaten
	}
	return 0
}

func (x *SubmitGameRequest) GetDeathType() uint32 {
	if x != nil {
		return x.DeathType
	}
	return 0
}

func (x *SubmitGameRequest) GetIsReplay() bool {
	if x != nil {
		return x.IsReplay
	}
	return false
}

func (x *SubmitGameRequest) GetReplayPlayerID() int32 {
	if x != nil {
		return x.ReplayPlayerID
	}
	return 0
}

func (x *SubmitGameRequest) GetPerEnemyAliveCount() []int32 {
	if x != nil {
		return x.PerEnemyAliveCount
	}
	return nil
}

func (x *SubmitGameRequest) GetPerEnemyKillcount() []int32 {
	if x != nil {
		return x.PerEnemyKillcount
	}
	return nil
}

func (x *SubmitGameRequest) GetStats() []*StatFrame {
	if x != nil {
		return x.Stats
	}
	return nil
}

type StatFrame struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GemsCollected      int32   `protobuf:"varint,1,opt,name=gemsCollected,proto3" json:"gemsCollected,omitempty"`
	Kills              int32   `protobuf:"varint,2,opt,name=kills,proto3" json:"kills,omitempty"`
	DaggersFired       int32   `protobuf:"varint,3,opt,name=daggersFired,proto3" json:"daggersFired,omitempty"`
	DaggersHit         int32   `protobuf:"varint,4,opt,name=daggersHit,proto3" json:"daggersHit,omitempty"`
	EnemiesAlive       int32   `protobuf:"varint,5,opt,name=enemiesAlive,proto3" json:"enemiesAlive,omitempty"`
	LevelGems          int32   `protobuf:"varint,6,opt,name=levelGems,proto3" json:"levelGems,omitempty"`
	HomingDaggers      int32   `protobuf:"varint,7,opt,name=homingDaggers,proto3" json:"homingDaggers,omitempty"`
	GemsDespawned      int32   `protobuf:"varint,8,opt,name=gemsDespawned,proto3" json:"gemsDespawned,omitempty"`
	GemsEaten          int32   `protobuf:"varint,9,opt,name=gemsEaten,proto3" json:"gemsEaten,omitempty"`
	TotalGems          int32   `protobuf:"varint,10,opt,name=totalGems,proto3" json:"totalGems,omitempty"`
	DaggersEaten       int32   `protobuf:"varint,11,opt,name=daggersEaten,proto3" json:"daggersEaten,omitempty"`
	PerEnemyAliveCount []int32 `protobuf:"varint,12,rep,packed,name=perEnemyAliveCount,proto3" json:"perEnemyAliveCount,omitempty"`
	PerEnemyKillCount  []int32 `protobuf:"varint,13,rep,packed,name=perEnemyKillCount,proto3" json:"perEnemyKillCount,omitempty"`
}

func (x *StatFrame) Reset() {
	*x = StatFrame{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gamesubmission_gamesubmission_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StatFrame) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StatFrame) ProtoMessage() {}

func (x *StatFrame) ProtoReflect() protoreflect.Message {
	mi := &file_gamesubmission_gamesubmission_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StatFrame.ProtoReflect.Descriptor instead.
func (*StatFrame) Descriptor() ([]byte, []int) {
	return file_gamesubmission_gamesubmission_proto_rawDescGZIP(), []int{1}
}

func (x *StatFrame) GetGemsCollected() int32 {
	if x != nil {
		return x.GemsCollected
	}
	return 0
}

func (x *StatFrame) GetKills() int32 {
	if x != nil {
		return x.Kills
	}
	return 0
}

func (x *StatFrame) GetDaggersFired() int32 {
	if x != nil {
		return x.DaggersFired
	}
	return 0
}

func (x *StatFrame) GetDaggersHit() int32 {
	if x != nil {
		return x.DaggersHit
	}
	return 0
}

func (x *StatFrame) GetEnemiesAlive() int32 {
	if x != nil {
		return x.EnemiesAlive
	}
	return 0
}

func (x *StatFrame) GetLevelGems() int32 {
	if x != nil {
		return x.LevelGems
	}
	return 0
}

func (x *StatFrame) GetHomingDaggers() int32 {
	if x != nil {
		return x.HomingDaggers
	}
	return 0
}

func (x *StatFrame) GetGemsDespawned() int32 {
	if x != nil {
		return x.GemsDespawned
	}
	return 0
}

func (x *StatFrame) GetGemsEaten() int32 {
	if x != nil {
		return x.GemsEaten
	}
	return 0
}

func (x *StatFrame) GetTotalGems() int32 {
	if x != nil {
		return x.TotalGems
	}
	return 0
}

func (x *StatFrame) GetDaggersEaten() int32 {
	if x != nil {
		return x.DaggersEaten
	}
	return 0
}

func (x *StatFrame) GetPerEnemyAliveCount() []int32 {
	if x != nil {
		return x.PerEnemyAliveCount
	}
	return nil
}

func (x *StatFrame) GetPerEnemyKillCount() []int32 {
	if x != nil {
		return x.PerEnemyKillCount
	}
	return nil
}

type SubmitGameReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GameID int32 `protobuf:"varint,1,opt,name=gameID,proto3" json:"gameID,omitempty"`
}

func (x *SubmitGameReply) Reset() {
	*x = SubmitGameReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gamesubmission_gamesubmission_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SubmitGameReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmitGameReply) ProtoMessage() {}

func (x *SubmitGameReply) ProtoReflect() protoreflect.Message {
	mi := &file_gamesubmission_gamesubmission_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmitGameReply.ProtoReflect.Descriptor instead.
func (*SubmitGameReply) Descriptor() ([]byte, []int) {
	return file_gamesubmission_gamesubmission_proto_rawDescGZIP(), []int{2}
}

func (x *SubmitGameReply) GetGameID() int32 {
	if x != nil {
		return x.GameID
	}
	return 0
}

var File_gamesubmission_gamesubmission_proto protoreflect.FileDescriptor

var file_gamesubmission_gamesubmission_proto_rawDesc = []byte{
	0x0a, 0x23, 0x67, 0x61, 0x6d, 0x65, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x2f, 0x67, 0x61, 0x6d, 0x65, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x67, 0x61, 0x6d, 0x65, 0x73, 0x75, 0x62, 0x6d, 0x69,
	0x73, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0xdc, 0x08, 0x0a, 0x11, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74,
	0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49,
	0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49,
	0x44, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x22, 0x0a, 0x0c, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x48, 0x61, 0x73, 0x68, 0x4d, 0x44,
	0x35, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x48, 0x61,
	0x73, 0x68, 0x4d, 0x44, 0x35, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x02, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d,
	0x65, 0x4c, 0x76, 0x6c, 0x6c, 0x32, 0x18, 0x06, 0x20, 0x01, 0x28, 0x02, 0x52, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x4c, 0x76, 0x6c, 0x6c, 0x32, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x4c,
	0x76, 0x6c, 0x6c, 0x33, 0x18, 0x07, 0x20, 0x01, 0x28, 0x02, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x4c, 0x76, 0x6c, 0x6c, 0x33, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x4c, 0x76, 0x6c,
	0x6c, 0x34, 0x18, 0x08, 0x20, 0x01, 0x28, 0x02, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x4c, 0x76,
	0x6c, 0x6c, 0x34, 0x12, 0x22, 0x0a, 0x0c, 0x74, 0x69, 0x6d, 0x65, 0x4c, 0x65, 0x76, 0x69, 0x44,
	0x6f, 0x77, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0c, 0x74, 0x69, 0x6d, 0x65, 0x4c,
	0x65, 0x76, 0x69, 0x44, 0x6f, 0x77, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x74, 0x69, 0x6d, 0x65, 0x4f,
	0x72, 0x62, 0x44, 0x6f, 0x77, 0x6e, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0b, 0x74, 0x69,
	0x6d, 0x65, 0x4f, 0x72, 0x62, 0x44, 0x6f, 0x77, 0x6e, 0x12, 0x24, 0x0a, 0x0d, 0x67, 0x65, 0x6d,
	0x73, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0d, 0x67, 0x65, 0x6d, 0x73, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x6b, 0x69, 0x6c, 0x6c, 0x73, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05,
	0x6b, 0x69, 0x6c, 0x6c, 0x73, 0x12, 0x22, 0x0a, 0x0c, 0x64, 0x61, 0x67, 0x67, 0x65, 0x72, 0x73,
	0x46, 0x69, 0x72, 0x65, 0x64, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x64, 0x61, 0x67,
	0x67, 0x65, 0x72, 0x73, 0x46, 0x69, 0x72, 0x65, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x61, 0x67,
	0x67, 0x65, 0x72, 0x73, 0x48, 0x69, 0x74, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x64,
	0x61, 0x67, 0x67, 0x65, 0x72, 0x73, 0x48, 0x69, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x65, 0x6e, 0x65,
	0x6d, 0x69, 0x65, 0x73, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0c, 0x65, 0x6e, 0x65, 0x6d, 0x69, 0x65, 0x73, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x12, 0x28, 0x0a,
	0x0f, 0x65, 0x6e, 0x65, 0x6d, 0x69, 0x65, 0x73, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x4d, 0x61, 0x78,
	0x18, 0x10, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0f, 0x65, 0x6e, 0x65, 0x6d, 0x69, 0x65, 0x73, 0x41,
	0x6c, 0x69, 0x76, 0x65, 0x4d, 0x61, 0x78, 0x12, 0x30, 0x0a, 0x13, 0x65, 0x6e, 0x65, 0x6d, 0x69,
	0x65, 0x73, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x4d, 0x61, 0x78, 0x54, 0x69, 0x6d, 0x65, 0x18, 0x11,
	0x20, 0x01, 0x28, 0x02, 0x52, 0x13, 0x65, 0x6e, 0x65, 0x6d, 0x69, 0x65, 0x73, 0x41, 0x6c, 0x69,
	0x76, 0x65, 0x4d, 0x61, 0x78, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x65, 0x76,
	0x65, 0x6c, 0x47, 0x65, 0x6d, 0x73, 0x18, 0x12, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x6c, 0x65,
	0x76, 0x65, 0x6c, 0x47, 0x65, 0x6d, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x68, 0x6f, 0x6d, 0x69, 0x6e,
	0x67, 0x44, 0x61, 0x67, 0x67, 0x65, 0x72, 0x73, 0x18, 0x13, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d,
	0x68, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x44, 0x61, 0x67, 0x67, 0x65, 0x72, 0x73, 0x12, 0x2a, 0x0a,
	0x10, 0x68, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x44, 0x61, 0x67, 0x67, 0x65, 0x72, 0x73, 0x4d, 0x61,
	0x78, 0x18, 0x14, 0x20, 0x01, 0x28, 0x05, 0x52, 0x10, 0x68, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x44,
	0x61, 0x67, 0x67, 0x65, 0x72, 0x73, 0x4d, 0x61, 0x78, 0x12, 0x32, 0x0a, 0x14, 0x68, 0x6f, 0x6d,
	0x69, 0x6e, 0x67, 0x44, 0x61, 0x67, 0x67, 0x65, 0x72, 0x73, 0x4d, 0x61, 0x78, 0x54, 0x69, 0x6d,
	0x65, 0x18, 0x15, 0x20, 0x01, 0x28, 0x02, 0x52, 0x14, 0x68, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x44,
	0x61, 0x67, 0x67, 0x65, 0x72, 0x73, 0x4d, 0x61, 0x78, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x47, 0x65, 0x6d, 0x73, 0x18, 0x16, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x09, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x47, 0x65, 0x6d, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x67,
	0x65, 0x6d, 0x73, 0x44, 0x65, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x65, 0x64, 0x18, 0x17, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0d, 0x67, 0x65, 0x6d, 0x73, 0x44, 0x65, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x65,
	0x64, 0x12, 0x1c, 0x0a, 0x09, 0x67, 0x65, 0x6d, 0x73, 0x45, 0x61, 0x74, 0x65, 0x6e, 0x18, 0x18,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x67, 0x65, 0x6d, 0x73, 0x45, 0x61, 0x74, 0x65, 0x6e, 0x12,
	0x22, 0x0a, 0x0c, 0x64, 0x61, 0x67, 0x67, 0x65, 0x72, 0x73, 0x45, 0x61, 0x74, 0x65, 0x6e, 0x18,
	0x19, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x64, 0x61, 0x67, 0x67, 0x65, 0x72, 0x73, 0x45, 0x61,
	0x74, 0x65, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x64, 0x65, 0x61, 0x74, 0x68, 0x54, 0x79, 0x70, 0x65,
	0x18, 0x1a, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x64, 0x65, 0x61, 0x74, 0x68, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x18, 0x1b, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x12, 0x26, 0x0a,
	0x0e, 0x72, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x44, 0x18,
	0x1c, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x72, 0x65, 0x70, 0x6c, 0x61, 0x79, 0x50, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x49, 0x44, 0x12, 0x2e, 0x0a, 0x12, 0x70, 0x65, 0x72, 0x45, 0x6e, 0x65, 0x6d,
	0x79, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x1d, 0x20, 0x03, 0x28,
	0x05, 0x52, 0x12, 0x70, 0x65, 0x72, 0x45, 0x6e, 0x65, 0x6d, 0x79, 0x41, 0x6c, 0x69, 0x76, 0x65,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2c, 0x0a, 0x11, 0x70, 0x65, 0x72, 0x45, 0x6e, 0x65, 0x6d,
	0x79, 0x4b, 0x69, 0x6c, 0x6c, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x1e, 0x20, 0x03, 0x28, 0x05,
	0x52, 0x11, 0x70, 0x65, 0x72, 0x45, 0x6e, 0x65, 0x6d, 0x79, 0x4b, 0x69, 0x6c, 0x6c, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x12, 0x2f, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x73, 0x18, 0x1f, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x19, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x52, 0x05, 0x73,
	0x74, 0x61, 0x74, 0x73, 0x22, 0xd7, 0x03, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x74, 0x46, 0x72, 0x61,
	0x6d, 0x65, 0x12, 0x24, 0x0a, 0x0d, 0x67, 0x65, 0x6d, 0x73, 0x43, 0x6f, 0x6c, 0x6c, 0x65, 0x63,
	0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x67, 0x65, 0x6d, 0x73, 0x43,
	0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6b, 0x69, 0x6c, 0x6c,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x6b, 0x69, 0x6c, 0x6c, 0x73, 0x12, 0x22,
	0x0a, 0x0c, 0x64, 0x61, 0x67, 0x67, 0x65, 0x72, 0x73, 0x46, 0x69, 0x72, 0x65, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x64, 0x61, 0x67, 0x67, 0x65, 0x72, 0x73, 0x46, 0x69, 0x72,
	0x65, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x64, 0x61, 0x67, 0x67, 0x65, 0x72, 0x73, 0x48, 0x69, 0x74,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x64, 0x61, 0x67, 0x67, 0x65, 0x72, 0x73, 0x48,
	0x69, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x65, 0x6e, 0x65, 0x6d, 0x69, 0x65, 0x73, 0x41, 0x6c, 0x69,
	0x76, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0c, 0x65, 0x6e, 0x65, 0x6d, 0x69, 0x65,
	0x73, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x47,
	0x65, 0x6d, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x6c, 0x65, 0x76, 0x65, 0x6c,
	0x47, 0x65, 0x6d, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x68, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x44, 0x61,
	0x67, 0x67, 0x65, 0x72, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0d, 0x68, 0x6f, 0x6d,
	0x69, 0x6e, 0x67, 0x44, 0x61, 0x67, 0x67, 0x65, 0x72, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x67, 0x65,
	0x6d, 0x73, 0x44, 0x65, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x65, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0d, 0x67, 0x65, 0x6d, 0x73, 0x44, 0x65, 0x73, 0x70, 0x61, 0x77, 0x6e, 0x65, 0x64,
	0x12, 0x1c, 0x0a, 0x09, 0x67, 0x65, 0x6d, 0x73, 0x45, 0x61, 0x74, 0x65, 0x6e, 0x18, 0x09, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x09, 0x67, 0x65, 0x6d, 0x73, 0x45, 0x61, 0x74, 0x65, 0x6e, 0x12, 0x1c,
	0x0a, 0x09, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x47, 0x65, 0x6d, 0x73, 0x18, 0x0a, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x09, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x47, 0x65, 0x6d, 0x73, 0x12, 0x22, 0x0a, 0x0c,
	0x64, 0x61, 0x67, 0x67, 0x65, 0x72, 0x73, 0x45, 0x61, 0x74, 0x65, 0x6e, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0c, 0x64, 0x61, 0x67, 0x67, 0x65, 0x72, 0x73, 0x45, 0x61, 0x74, 0x65, 0x6e,
	0x12, 0x2e, 0x0a, 0x12, 0x70, 0x65, 0x72, 0x45, 0x6e, 0x65, 0x6d, 0x79, 0x41, 0x6c, 0x69, 0x76,
	0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x0c, 0x20, 0x03, 0x28, 0x05, 0x52, 0x12, 0x70, 0x65,
	0x72, 0x45, 0x6e, 0x65, 0x6d, 0x79, 0x41, 0x6c, 0x69, 0x76, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x2c, 0x0a, 0x11, 0x70, 0x65, 0x72, 0x45, 0x6e, 0x65, 0x6d, 0x79, 0x4b, 0x69, 0x6c, 0x6c,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x0d, 0x20, 0x03, 0x28, 0x05, 0x52, 0x11, 0x70, 0x65, 0x72,
	0x45, 0x6e, 0x65, 0x6d, 0x79, 0x4b, 0x69, 0x6c, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x29,
	0x0a, 0x0f, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x12, 0x16, 0x0a, 0x06, 0x67, 0x61, 0x6d, 0x65, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x06, 0x67, 0x61, 0x6d, 0x65, 0x49, 0x44, 0x32, 0x62, 0x0a, 0x0c, 0x47, 0x61, 0x6d,
	0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x12, 0x52, 0x0a, 0x0a, 0x53, 0x75, 0x62,
	0x6d, 0x69, 0x74, 0x47, 0x61, 0x6d, 0x65, 0x12, 0x21, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x73, 0x75,
	0x62, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x74, 0x47,
	0x61, 0x6d, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x67, 0x61, 0x6d,
	0x65, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x2e, 0x53, 0x75, 0x62, 0x6d,
	0x69, 0x74, 0x47, 0x61, 0x6d, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x38, 0x5a,
	0x36, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x6c, 0x65, 0x78,
	0x77, 0x69, 0x6c, 0x6b, 0x65, 0x72, 0x73, 0x6f, 0x6e, 0x2f, 0x64, 0x64, 0x73, 0x74, 0x61, 0x74,
	0x73, 0x2d, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2f, 0x67, 0x61, 0x6d, 0x65, 0x73, 0x75, 0x62,
	0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gamesubmission_gamesubmission_proto_rawDescOnce sync.Once
	file_gamesubmission_gamesubmission_proto_rawDescData = file_gamesubmission_gamesubmission_proto_rawDesc
)

func file_gamesubmission_gamesubmission_proto_rawDescGZIP() []byte {
	file_gamesubmission_gamesubmission_proto_rawDescOnce.Do(func() {
		file_gamesubmission_gamesubmission_proto_rawDescData = protoimpl.X.CompressGZIP(file_gamesubmission_gamesubmission_proto_rawDescData)
	})
	return file_gamesubmission_gamesubmission_proto_rawDescData
}

var file_gamesubmission_gamesubmission_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_gamesubmission_gamesubmission_proto_goTypes = []interface{}{
	(*SubmitGameRequest)(nil), // 0: gamesubmission.SubmitGameRequest
	(*StatFrame)(nil),         // 1: gamesubmission.StatFrame
	(*SubmitGameReply)(nil),   // 2: gamesubmission.SubmitGameReply
}
var file_gamesubmission_gamesubmission_proto_depIdxs = []int32{
	1, // 0: gamesubmission.SubmitGameRequest.stats:type_name -> gamesubmission.StatFrame
	0, // 1: gamesubmission.GameRecorder.SubmitGame:input_type -> gamesubmission.SubmitGameRequest
	2, // 2: gamesubmission.GameRecorder.SubmitGame:output_type -> gamesubmission.SubmitGameReply
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_gamesubmission_gamesubmission_proto_init() }
func file_gamesubmission_gamesubmission_proto_init() {
	if File_gamesubmission_gamesubmission_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gamesubmission_gamesubmission_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubmitGameRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gamesubmission_gamesubmission_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StatFrame); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_gamesubmission_gamesubmission_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SubmitGameReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_gamesubmission_gamesubmission_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_gamesubmission_gamesubmission_proto_goTypes,
		DependencyIndexes: file_gamesubmission_gamesubmission_proto_depIdxs,
		MessageInfos:      file_gamesubmission_gamesubmission_proto_msgTypes,
	}.Build()
	File_gamesubmission_gamesubmission_proto = out.File
	file_gamesubmission_gamesubmission_proto_rawDesc = nil
	file_gamesubmission_gamesubmission_proto_goTypes = nil
	file_gamesubmission_gamesubmission_proto_depIdxs = nil
}
