// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v3.21.12
// source: api/proto/music.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type StreamRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MusicId       string                 `protobuf:"bytes,1,opt,name=music_id,json=musicId,proto3" json:"music_id,omitempty"`
	StartPosition int32                  `protobuf:"varint,2,opt,name=start_position,json=startPosition,proto3" json:"start_position,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamRequest) Reset() {
	*x = StreamRequest{}
	mi := &file_api_proto_music_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamRequest) ProtoMessage() {}

func (x *StreamRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_music_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamRequest.ProtoReflect.Descriptor instead.
func (*StreamRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_music_proto_rawDescGZIP(), []int{0}
}

func (x *StreamRequest) GetMusicId() string {
	if x != nil {
		return x.MusicId
	}
	return ""
}

func (x *StreamRequest) GetStartPosition() int32 {
	if x != nil {
		return x.StartPosition
	}
	return 0
}

type AudioChunk struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	Data           []byte                 `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	SequenceNumber int32                  `protobuf:"varint,2,opt,name=sequence_number,json=sequenceNumber,proto3" json:"sequence_number,omitempty"`
	Type           string                 `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *AudioChunk) Reset() {
	*x = AudioChunk{}
	mi := &file_api_proto_music_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AudioChunk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AudioChunk) ProtoMessage() {}

func (x *AudioChunk) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_music_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AudioChunk.ProtoReflect.Descriptor instead.
func (*AudioChunk) Descriptor() ([]byte, []int) {
	return file_api_proto_music_proto_rawDescGZIP(), []int{1}
}

func (x *AudioChunk) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *AudioChunk) GetSequenceNumber() int32 {
	if x != nil {
		return x.SequenceNumber
	}
	return 0
}

func (x *AudioChunk) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type SearchRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Query         string                 `protobuf:"bytes,1,opt,name=query,proto3" json:"query,omitempty"`
	Page          int32                  `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	PageSize      int32                  `protobuf:"varint,3,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchRequest) Reset() {
	*x = SearchRequest{}
	mi := &file_api_proto_music_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchRequest) ProtoMessage() {}

func (x *SearchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_music_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchRequest.ProtoReflect.Descriptor instead.
func (*SearchRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_music_proto_rawDescGZIP(), []int{2}
}

func (x *SearchRequest) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

func (x *SearchRequest) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *SearchRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type SearchResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MusicList     []*Music               `protobuf:"bytes,1,rep,name=music_list,json=musicList,proto3" json:"music_list,omitempty"`
	Total         int32                  `protobuf:"varint,2,opt,name=total,proto3" json:"total,omitempty"`
	MetadataList  []*MusicMetadata       `protobuf:"bytes,3,rep,name=metadata_list,json=metadataList,proto3" json:"metadata_list,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SearchResponse) Reset() {
	*x = SearchResponse{}
	mi := &file_api_proto_music_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SearchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchResponse) ProtoMessage() {}

func (x *SearchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_music_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchResponse.ProtoReflect.Descriptor instead.
func (*SearchResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_music_proto_rawDescGZIP(), []int{3}
}

func (x *SearchResponse) GetMusicList() []*Music {
	if x != nil {
		return x.MusicList
	}
	return nil
}

func (x *SearchResponse) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

func (x *SearchResponse) GetMetadataList() []*MusicMetadata {
	if x != nil {
		return x.MetadataList
	}
	return nil
}

type Music struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Title         string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Artist        string                 `protobuf:"bytes,3,opt,name=artist,proto3" json:"artist,omitempty"`
	Album         string                 `protobuf:"bytes,4,opt,name=album,proto3" json:"album,omitempty"`
	Duration      int32                  `protobuf:"varint,5,opt,name=duration,proto3" json:"duration,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Music) Reset() {
	*x = Music{}
	mi := &file_api_proto_music_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Music) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Music) ProtoMessage() {}

func (x *Music) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_music_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Music.ProtoReflect.Descriptor instead.
func (*Music) Descriptor() ([]byte, []int) {
	return file_api_proto_music_proto_rawDescGZIP(), []int{4}
}

func (x *Music) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Music) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Music) GetArtist() string {
	if x != nil {
		return x.Artist
	}
	return ""
}

func (x *Music) GetAlbum() string {
	if x != nil {
		return x.Album
	}
	return ""
}

func (x *Music) GetDuration() int32 {
	if x != nil {
		return x.Duration
	}
	return 0
}

type UploadRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Types that are valid to be assigned to Data:
	//
	//	*UploadRequest_Metadata
	//	*UploadRequest_ChunkData
	Data          isUploadRequest_Data `protobuf_oneof:"data"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadRequest) Reset() {
	*x = UploadRequest{}
	mi := &file_api_proto_music_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadRequest) ProtoMessage() {}

func (x *UploadRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_music_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadRequest.ProtoReflect.Descriptor instead.
func (*UploadRequest) Descriptor() ([]byte, []int) {
	return file_api_proto_music_proto_rawDescGZIP(), []int{5}
}

func (x *UploadRequest) GetData() isUploadRequest_Data {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *UploadRequest) GetMetadata() *MusicMetadata {
	if x != nil {
		if x, ok := x.Data.(*UploadRequest_Metadata); ok {
			return x.Metadata
		}
	}
	return nil
}

func (x *UploadRequest) GetChunkData() []byte {
	if x != nil {
		if x, ok := x.Data.(*UploadRequest_ChunkData); ok {
			return x.ChunkData
		}
	}
	return nil
}

type isUploadRequest_Data interface {
	isUploadRequest_Data()
}

type UploadRequest_Metadata struct {
	Metadata *MusicMetadata `protobuf:"bytes,1,opt,name=metadata,proto3,oneof"`
}

type UploadRequest_ChunkData struct {
	ChunkData []byte `protobuf:"bytes,2,opt,name=chunk_data,json=chunkData,proto3,oneof"`
}

func (*UploadRequest_Metadata) isUploadRequest_Data() {}

func (*UploadRequest_ChunkData) isUploadRequest_Data() {}

type MusicMetadata struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Title         string                 `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	Artist        string                 `protobuf:"bytes,2,opt,name=artist,proto3" json:"artist,omitempty"`
	Album         string                 `protobuf:"bytes,3,opt,name=album,proto3" json:"album,omitempty"`
	Type          string                 `protobuf:"bytes,4,opt,name=type,proto3" json:"type,omitempty"`
	Year          int32                  `protobuf:"varint,5,opt,name=year,proto3" json:"year,omitempty"`                                       // Ano de lançamento
	Genre         string                 `protobuf:"bytes,6,opt,name=genre,proto3" json:"genre,omitempty"`                                      // Gênero musical
	Composer      string                 `protobuf:"bytes,7,opt,name=composer,proto3" json:"composer,omitempty"`                                // Compositor
	Label         string                 `protobuf:"bytes,8,opt,name=label,proto3" json:"label,omitempty"`                                      // Gravadora
	AlbumArt      []byte                 `protobuf:"bytes,9,opt,name=album_art,json=albumArt,proto3" json:"album_art,omitempty"`                // Capa do álbum
	AlbumArtType  string                 `protobuf:"bytes,10,opt,name=album_art_type,json=albumArtType,proto3" json:"album_art_type,omitempty"` //tipo do arquivo da capa
	Comments      string                 `protobuf:"bytes,11,opt,name=comments,proto3" json:"comments,omitempty"`                               // Comentários
	Isrc          string                 `protobuf:"bytes,12,opt,name=isrc,proto3" json:"isrc,omitempty"`                                       // Código ISRC
	Url           string                 `protobuf:"bytes,13,opt,name=url,proto3" json:"url,omitempty"`                                         // URL para informações adicionais
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MusicMetadata) Reset() {
	*x = MusicMetadata{}
	mi := &file_api_proto_music_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MusicMetadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MusicMetadata) ProtoMessage() {}

func (x *MusicMetadata) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_music_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MusicMetadata.ProtoReflect.Descriptor instead.
func (*MusicMetadata) Descriptor() ([]byte, []int) {
	return file_api_proto_music_proto_rawDescGZIP(), []int{6}
}

func (x *MusicMetadata) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *MusicMetadata) GetArtist() string {
	if x != nil {
		return x.Artist
	}
	return ""
}

func (x *MusicMetadata) GetAlbum() string {
	if x != nil {
		return x.Album
	}
	return ""
}

func (x *MusicMetadata) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *MusicMetadata) GetYear() int32 {
	if x != nil {
		return x.Year
	}
	return 0
}

func (x *MusicMetadata) GetGenre() string {
	if x != nil {
		return x.Genre
	}
	return ""
}

func (x *MusicMetadata) GetComposer() string {
	if x != nil {
		return x.Composer
	}
	return ""
}

func (x *MusicMetadata) GetLabel() string {
	if x != nil {
		return x.Label
	}
	return ""
}

func (x *MusicMetadata) GetAlbumArt() []byte {
	if x != nil {
		return x.AlbumArt
	}
	return nil
}

func (x *MusicMetadata) GetAlbumArtType() string {
	if x != nil {
		return x.AlbumArtType
	}
	return ""
}

func (x *MusicMetadata) GetComments() string {
	if x != nil {
		return x.Comments
	}
	return ""
}

func (x *MusicMetadata) GetIsrc() string {
	if x != nil {
		return x.Isrc
	}
	return ""
}

func (x *MusicMetadata) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type UploadResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	MusicId       string                 `protobuf:"bytes,1,opt,name=music_id,json=musicId,proto3" json:"music_id,omitempty"`
	Success       bool                   `protobuf:"varint,2,opt,name=success,proto3" json:"success,omitempty"`
	Message       string                 `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UploadResponse) Reset() {
	*x = UploadResponse{}
	mi := &file_api_proto_music_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UploadResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UploadResponse) ProtoMessage() {}

func (x *UploadResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_proto_music_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UploadResponse.ProtoReflect.Descriptor instead.
func (*UploadResponse) Descriptor() ([]byte, []int) {
	return file_api_proto_music_proto_rawDescGZIP(), []int{7}
}

func (x *UploadResponse) GetMusicId() string {
	if x != nil {
		return x.MusicId
	}
	return ""
}

func (x *UploadResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *UploadResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_api_proto_music_proto protoreflect.FileDescriptor

var file_api_proto_music_proto_rawDesc = string([]byte{
	0x0a, 0x15, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x75, 0x73, 0x69,
	0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x22, 0x51,
	0x0a, 0x0d, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x19, 0x0a, 0x08, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x49, 0x64, 0x12, 0x25, 0x0a, 0x0e, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x5f, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0d, 0x73, 0x74, 0x61, 0x72, 0x74, 0x50, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x22, 0x5d, 0x0a, 0x0a, 0x41, 0x75, 0x64, 0x69, 0x6f, 0x43, 0x68, 0x75, 0x6e, 0x6b, 0x12,
	0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x12, 0x27, 0x0a, 0x0f, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x5f,
	0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x73, 0x65,
	0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x22, 0x56, 0x0a, 0x0d, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x70,
	0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x8e, 0x01, 0x0a, 0x0e, 0x53, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x0a, 0x6d,
	0x75, 0x73, 0x69, 0x63, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x0c, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x2e, 0x4d, 0x75, 0x73, 0x69, 0x63, 0x52, 0x09, 0x6d,
	0x75, 0x73, 0x69, 0x63, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x12, 0x39,
	0x0a, 0x0d, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x2e, 0x4d, 0x75,
	0x73, 0x69, 0x63, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x52, 0x0c, 0x6d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x4c, 0x69, 0x73, 0x74, 0x22, 0x77, 0x0a, 0x05, 0x4d, 0x75, 0x73,
	0x69, 0x63, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x72, 0x74, 0x69,
	0x73, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x72, 0x74, 0x69, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x22, 0x6c, 0x0a, 0x0d, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x12, 0x32, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x2e, 0x4d, 0x75,
	0x73, 0x69, 0x63, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52, 0x08, 0x6d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1f, 0x0a, 0x0a, 0x63, 0x68, 0x75, 0x6e, 0x6b,
	0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x09, 0x63,
	0x68, 0x75, 0x6e, 0x6b, 0x44, 0x61, 0x74, 0x61, 0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x22, 0xc8, 0x02, 0x0a, 0x0d, 0x4d, 0x75, 0x73, 0x69, 0x63, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x72, 0x74, 0x69,
	0x73, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x72, 0x74, 0x69, 0x73, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x79, 0x65,
	0x61, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x79, 0x65, 0x61, 0x72, 0x12, 0x14,
	0x0a, 0x05, 0x67, 0x65, 0x6e, 0x72, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x67,
	0x65, 0x6e, 0x72, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x73, 0x65, 0x72,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x73, 0x65, 0x72,
	0x12, 0x14, 0x0a, 0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x12, 0x1b, 0x0a, 0x09, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x5f,
	0x61, 0x72, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x61, 0x6c, 0x62, 0x75, 0x6d,
	0x41, 0x72, 0x74, 0x12, 0x24, 0x0a, 0x0e, 0x61, 0x6c, 0x62, 0x75, 0x6d, 0x5f, 0x61, 0x72, 0x74,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x61, 0x6c, 0x62,
	0x75, 0x6d, 0x41, 0x72, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x63, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x63, 0x6f, 0x6d,
	0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x73, 0x72, 0x63, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x69, 0x73, 0x72, 0x63, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c,
	0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x5f, 0x0a, 0x0e, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x19, 0x0a,
	0x08, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x07, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0xc8, 0x01, 0x0a,
	0x0c, 0x4d, 0x75, 0x73, 0x69, 0x63, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3a, 0x0a,
	0x0b, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x4d, 0x75, 0x73, 0x69, 0x63, 0x12, 0x14, 0x2e, 0x6d,
	0x75, 0x73, 0x69, 0x63, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x11, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x2e, 0x41, 0x75, 0x64, 0x69, 0x6f,
	0x43, 0x68, 0x75, 0x6e, 0x6b, 0x22, 0x00, 0x30, 0x01, 0x12, 0x3c, 0x0a, 0x0b, 0x53, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x4d, 0x75, 0x73, 0x69, 0x63, 0x12, 0x14, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63,
	0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15,
	0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x3e, 0x0a, 0x0b, 0x55, 0x70, 0x6c, 0x6f, 0x61,
	0x64, 0x4d, 0x75, 0x73, 0x69, 0x63, 0x12, 0x14, 0x2e, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x2e, 0x55,
	0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x6d,
	0x75, 0x73, 0x69, 0x63, 0x2e, 0x55, 0x70, 0x6c, 0x6f, 0x61, 0x64, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x42, 0x31, 0x5a, 0x2f, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x65, 0x72, 0x69, 0x6c, 0x6f, 0x71, 0x75, 0x65, 0x69,
	0x72, 0x6f, 0x7a, 0x2f, 0x6d, 0x75, 0x73, 0x69, 0x63, 0x2d, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x2f, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
})

var (
	file_api_proto_music_proto_rawDescOnce sync.Once
	file_api_proto_music_proto_rawDescData []byte
)

func file_api_proto_music_proto_rawDescGZIP() []byte {
	file_api_proto_music_proto_rawDescOnce.Do(func() {
		file_api_proto_music_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_api_proto_music_proto_rawDesc), len(file_api_proto_music_proto_rawDesc)))
	})
	return file_api_proto_music_proto_rawDescData
}

var file_api_proto_music_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_api_proto_music_proto_goTypes = []any{
	(*StreamRequest)(nil),  // 0: music.StreamRequest
	(*AudioChunk)(nil),     // 1: music.AudioChunk
	(*SearchRequest)(nil),  // 2: music.SearchRequest
	(*SearchResponse)(nil), // 3: music.SearchResponse
	(*Music)(nil),          // 4: music.Music
	(*UploadRequest)(nil),  // 5: music.UploadRequest
	(*MusicMetadata)(nil),  // 6: music.MusicMetadata
	(*UploadResponse)(nil), // 7: music.UploadResponse
}
var file_api_proto_music_proto_depIdxs = []int32{
	4, // 0: music.SearchResponse.music_list:type_name -> music.Music
	6, // 1: music.SearchResponse.metadata_list:type_name -> music.MusicMetadata
	6, // 2: music.UploadRequest.metadata:type_name -> music.MusicMetadata
	0, // 3: music.MusicService.StreamMusic:input_type -> music.StreamRequest
	2, // 4: music.MusicService.SearchMusic:input_type -> music.SearchRequest
	5, // 5: music.MusicService.UploadMusic:input_type -> music.UploadRequest
	1, // 6: music.MusicService.StreamMusic:output_type -> music.AudioChunk
	3, // 7: music.MusicService.SearchMusic:output_type -> music.SearchResponse
	7, // 8: music.MusicService.UploadMusic:output_type -> music.UploadResponse
	6, // [6:9] is the sub-list for method output_type
	3, // [3:6] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_api_proto_music_proto_init() }
func file_api_proto_music_proto_init() {
	if File_api_proto_music_proto != nil {
		return
	}
	file_api_proto_music_proto_msgTypes[5].OneofWrappers = []any{
		(*UploadRequest_Metadata)(nil),
		(*UploadRequest_ChunkData)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_api_proto_music_proto_rawDesc), len(file_api_proto_music_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_proto_music_proto_goTypes,
		DependencyIndexes: file_api_proto_music_proto_depIdxs,
		MessageInfos:      file_api_proto_music_proto_msgTypes,
	}.Build()
	File_api_proto_music_proto = out.File
	file_api_proto_music_proto_goTypes = nil
	file_api_proto_music_proto_depIdxs = nil
}
