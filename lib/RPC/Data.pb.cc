// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: Data.proto

#include "Data.pb.h"

#include <algorithm>

#include <google/protobuf/stubs/common.h>
#include <google/protobuf/io/coded_stream.h>
#include <google/protobuf/extension_set.h>
#include <google/protobuf/wire_format_lite.h>
#include <google/protobuf/descriptor.h>
#include <google/protobuf/generated_message_reflection.h>
#include <google/protobuf/reflection_ops.h>
#include <google/protobuf/wire_format.h>
// @@protoc_insertion_point(includes)
#include <google/protobuf/port_def.inc>
extern PROTOBUF_INTERNAL_EXPORT_Data_2eproto ::PROTOBUF_NAMESPACE_ID::internal::SCCInfo<0> scc_info_all_data_AddressTimeEntry_DoNotUse_Data_2eproto;
extern PROTOBUF_INTERNAL_EXPORT_Data_2eproto ::PROTOBUF_NAMESPACE_ID::internal::SCCInfo<0> scc_info_data_Data_2eproto;
namespace dra {
class dataDefaultTypeInternal {
 public:
  ::PROTOBUF_NAMESPACE_ID::internal::ExplicitlyConstructed<data> _instance;
} _data_default_instance_;
class all_data_AddressTimeEntry_DoNotUseDefaultTypeInternal {
 public:
  ::PROTOBUF_NAMESPACE_ID::internal::ExplicitlyConstructed<all_data_AddressTimeEntry_DoNotUse> _instance;
} _all_data_AddressTimeEntry_DoNotUse_default_instance_;
class all_dataDefaultTypeInternal {
 public:
  ::PROTOBUF_NAMESPACE_ID::internal::ExplicitlyConstructed<all_data> _instance;
} _all_data_default_instance_;
}  // namespace dra
static void InitDefaultsscc_info_all_data_Data_2eproto() {
  GOOGLE_PROTOBUF_VERIFY_VERSION;

  {
    void* ptr = &::dra::_all_data_default_instance_;
    new (ptr) ::dra::all_data();
    ::PROTOBUF_NAMESPACE_ID::internal::OnShutdownDestroyMessage(ptr);
  }
  ::dra::all_data::InitAsDefaultInstance();
}

::PROTOBUF_NAMESPACE_ID::internal::SCCInfo<2> scc_info_all_data_Data_2eproto =
    {{ATOMIC_VAR_INIT(::PROTOBUF_NAMESPACE_ID::internal::SCCInfoBase::kUninitialized), 2, InitDefaultsscc_info_all_data_Data_2eproto}, {
      &scc_info_all_data_AddressTimeEntry_DoNotUse_Data_2eproto.base,
      &scc_info_data_Data_2eproto.base,}};

static void InitDefaultsscc_info_all_data_AddressTimeEntry_DoNotUse_Data_2eproto() {
  GOOGLE_PROTOBUF_VERIFY_VERSION;

  {
    void* ptr = &::dra::_all_data_AddressTimeEntry_DoNotUse_default_instance_;
    new (ptr) ::dra::all_data_AddressTimeEntry_DoNotUse();
  }
  ::dra::all_data_AddressTimeEntry_DoNotUse::InitAsDefaultInstance();
}

::PROTOBUF_NAMESPACE_ID::internal::SCCInfo<0> scc_info_all_data_AddressTimeEntry_DoNotUse_Data_2eproto =
    {{ATOMIC_VAR_INIT(::PROTOBUF_NAMESPACE_ID::internal::SCCInfoBase::kUninitialized), 0, InitDefaultsscc_info_all_data_AddressTimeEntry_DoNotUse_Data_2eproto}, {}};

static void InitDefaultsscc_info_data_Data_2eproto() {
  GOOGLE_PROTOBUF_VERIFY_VERSION;

  {
    void* ptr = &::dra::_data_default_instance_;
    new (ptr) ::dra::data();
    ::PROTOBUF_NAMESPACE_ID::internal::OnShutdownDestroyMessage(ptr);
  }
  ::dra::data::InitAsDefaultInstance();
}

::PROTOBUF_NAMESPACE_ID::internal::SCCInfo<0> scc_info_data_Data_2eproto =
    {{ATOMIC_VAR_INIT(::PROTOBUF_NAMESPACE_ID::internal::SCCInfoBase::kUninitialized), 0, InitDefaultsscc_info_data_Data_2eproto}, {}};

static ::PROTOBUF_NAMESPACE_ID::Metadata file_level_metadata_Data_2eproto[3];
static constexpr ::PROTOBUF_NAMESPACE_ID::EnumDescriptor const** file_level_enum_descriptors_Data_2eproto = nullptr;
static constexpr ::PROTOBUF_NAMESPACE_ID::ServiceDescriptor const** file_level_service_descriptors_Data_2eproto = nullptr;

const ::PROTOBUF_NAMESPACE_ID::uint32 TableStruct_Data_2eproto::offsets[] PROTOBUF_SECTION_VARIABLE(protodesc_cold) = {
  ~0u,  // no _has_bits_
  PROTOBUF_FIELD_OFFSET(::dra::data, _internal_metadata_),
  ~0u,  // no _extensions_
  ~0u,  // no _oneof_case_
  ~0u,  // no _weak_field_map_
  PROTOBUF_FIELD_OFFSET(::dra::data, time_),
  PROTOBUF_FIELD_OFFSET(::dra::data, date_),
  PROTOBUF_FIELD_OFFSET(::dra::data, uncovered_),
  PROTOBUF_FIELD_OFFSET(::dra::all_data_AddressTimeEntry_DoNotUse, _has_bits_),
  PROTOBUF_FIELD_OFFSET(::dra::all_data_AddressTimeEntry_DoNotUse, _internal_metadata_),
  ~0u,  // no _extensions_
  ~0u,  // no _oneof_case_
  ~0u,  // no _weak_field_map_
  PROTOBUF_FIELD_OFFSET(::dra::all_data_AddressTimeEntry_DoNotUse, key_),
  PROTOBUF_FIELD_OFFSET(::dra::all_data_AddressTimeEntry_DoNotUse, value_),
  0,
  1,
  ~0u,  // no _has_bits_
  PROTOBUF_FIELD_OFFSET(::dra::all_data, _internal_metadata_),
  ~0u,  // no _extensions_
  ~0u,  // no _oneof_case_
  ~0u,  // no _weak_field_map_
  PROTOBUF_FIELD_OFFSET(::dra::all_data, name_),
  PROTOBUF_FIELD_OFFSET(::dra::all_data, address_time_),
  PROTOBUF_FIELD_OFFSET(::dra::all_data, data_),
};
static const ::PROTOBUF_NAMESPACE_ID::internal::MigrationSchema schemas[] PROTOBUF_SECTION_VARIABLE(protodesc_cold) = {
  { 0, -1, sizeof(::dra::data)},
  { 8, 15, sizeof(::dra::all_data_AddressTimeEntry_DoNotUse)},
  { 17, -1, sizeof(::dra::all_data)},
};

static ::PROTOBUF_NAMESPACE_ID::Message const * const file_default_instances[] = {
  reinterpret_cast<const ::PROTOBUF_NAMESPACE_ID::Message*>(&::dra::_data_default_instance_),
  reinterpret_cast<const ::PROTOBUF_NAMESPACE_ID::Message*>(&::dra::_all_data_AddressTimeEntry_DoNotUse_default_instance_),
  reinterpret_cast<const ::PROTOBUF_NAMESPACE_ID::Message*>(&::dra::_all_data_default_instance_),
};

const char descriptor_table_protodef_Data_2eproto[] =
  "\n\nData.proto\022\003dra\"5\n\004data\022\014\n\004time\030\001 \001(\004\022"
  "\014\n\004date\030\002 \001(\t\022\021\n\tuncovered\030\003 \001(\004\"\233\001\n\010all"
  "_data\022\014\n\004name\030\001 \001(\t\0224\n\014address_time\030\002 \003("
  "\0132\036.dra.all_data.AddressTimeEntry\022\027\n\004dat"
  "a\030\003 \003(\0132\t.dra.data\0322\n\020AddressTimeEntry\022\013"
  "\n\003key\030\001 \001(\004\022\r\n\005value\030\002 \001(\004:\0028\001b\006proto3"
  ;
static const ::PROTOBUF_NAMESPACE_ID::internal::DescriptorTable*const descriptor_table_Data_2eproto_deps[1] = {
};
static ::PROTOBUF_NAMESPACE_ID::internal::SCCInfoBase*const descriptor_table_Data_2eproto_sccs[3] = {
  &scc_info_all_data_Data_2eproto.base,
  &scc_info_all_data_AddressTimeEntry_DoNotUse_Data_2eproto.base,
  &scc_info_data_Data_2eproto.base,
};
static ::PROTOBUF_NAMESPACE_ID::internal::once_flag descriptor_table_Data_2eproto_once;
static bool descriptor_table_Data_2eproto_initialized = false;
const ::PROTOBUF_NAMESPACE_ID::internal::DescriptorTable descriptor_table_Data_2eproto = {
  &descriptor_table_Data_2eproto_initialized, descriptor_table_protodef_Data_2eproto, "Data.proto", 238,
  &descriptor_table_Data_2eproto_once, descriptor_table_Data_2eproto_sccs, descriptor_table_Data_2eproto_deps, 3, 0,
  schemas, file_default_instances, TableStruct_Data_2eproto::offsets,
  file_level_metadata_Data_2eproto, 3, file_level_enum_descriptors_Data_2eproto, file_level_service_descriptors_Data_2eproto,
};

// Force running AddDescriptors() at dynamic initialization time.
static bool dynamic_init_dummy_Data_2eproto = (  ::PROTOBUF_NAMESPACE_ID::internal::AddDescriptors(&descriptor_table_Data_2eproto), true);
namespace dra {

// ===================================================================

void data::InitAsDefaultInstance() {
}
class data::HasBitSetters {
 public:
};

#if !defined(_MSC_VER) || _MSC_VER >= 1900
const int data::kTimeFieldNumber;
const int data::kDateFieldNumber;
const int data::kUncoveredFieldNumber;
#endif  // !defined(_MSC_VER) || _MSC_VER >= 1900

data::data()
  : ::PROTOBUF_NAMESPACE_ID::Message(), _internal_metadata_(nullptr) {
  SharedCtor();
  // @@protoc_insertion_point(constructor:dra.data)
}
data::data(const data& from)
  : ::PROTOBUF_NAMESPACE_ID::Message(),
      _internal_metadata_(nullptr) {
  _internal_metadata_.MergeFrom(from._internal_metadata_);
  date_.UnsafeSetDefault(&::PROTOBUF_NAMESPACE_ID::internal::GetEmptyStringAlreadyInited());
  if (from.date().size() > 0) {
    date_.AssignWithDefault(&::PROTOBUF_NAMESPACE_ID::internal::GetEmptyStringAlreadyInited(), from.date_);
  }
  ::memcpy(&time_, &from.time_,
    static_cast<size_t>(reinterpret_cast<char*>(&uncovered_) -
    reinterpret_cast<char*>(&time_)) + sizeof(uncovered_));
  // @@protoc_insertion_point(copy_constructor:dra.data)
}

void data::SharedCtor() {
  ::PROTOBUF_NAMESPACE_ID::internal::InitSCC(&scc_info_data_Data_2eproto.base);
  date_.UnsafeSetDefault(&::PROTOBUF_NAMESPACE_ID::internal::GetEmptyStringAlreadyInited());
  ::memset(&time_, 0, static_cast<size_t>(
      reinterpret_cast<char*>(&uncovered_) -
      reinterpret_cast<char*>(&time_)) + sizeof(uncovered_));
}

data::~data() {
  // @@protoc_insertion_point(destructor:dra.data)
  SharedDtor();
}

void data::SharedDtor() {
  date_.DestroyNoArena(&::PROTOBUF_NAMESPACE_ID::internal::GetEmptyStringAlreadyInited());
}

void data::SetCachedSize(int size) const {
  _cached_size_.Set(size);
}
const data& data::default_instance() {
  ::PROTOBUF_NAMESPACE_ID::internal::InitSCC(&::scc_info_data_Data_2eproto.base);
  return *internal_default_instance();
}


void data::Clear() {
// @@protoc_insertion_point(message_clear_start:dra.data)
  ::PROTOBUF_NAMESPACE_ID::uint32 cached_has_bits = 0;
  // Prevent compiler warnings about cached_has_bits being unused
  (void) cached_has_bits;

  date_.ClearToEmptyNoArena(&::PROTOBUF_NAMESPACE_ID::internal::GetEmptyStringAlreadyInited());
  ::memset(&time_, 0, static_cast<size_t>(
      reinterpret_cast<char*>(&uncovered_) -
      reinterpret_cast<char*>(&time_)) + sizeof(uncovered_));
  _internal_metadata_.Clear();
}

#if GOOGLE_PROTOBUF_ENABLE_EXPERIMENTAL_PARSER
const char* data::_InternalParse(const char* ptr, ::PROTOBUF_NAMESPACE_ID::internal::ParseContext* ctx) {
#define CHK_(x) if (PROTOBUF_PREDICT_FALSE(!(x))) goto failure
  while (!ctx->Done(&ptr)) {
    ::PROTOBUF_NAMESPACE_ID::uint32 tag;
    ptr = ::PROTOBUF_NAMESPACE_ID::internal::ReadTag(ptr, &tag);
    CHK_(ptr);
    switch (tag >> 3) {
      // uint64 time = 1;
      case 1:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::PROTOBUF_NAMESPACE_ID::uint8>(tag) == 8)) {
          time_ = ::PROTOBUF_NAMESPACE_ID::internal::ReadVarint(&ptr);
          CHK_(ptr);
        } else goto handle_unusual;
        continue;
      // string date = 2;
      case 2:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::PROTOBUF_NAMESPACE_ID::uint8>(tag) == 18)) {
          ptr = ::PROTOBUF_NAMESPACE_ID::internal::InlineGreedyStringParserUTF8(mutable_date(), ptr, ctx, "dra.data.date");
          CHK_(ptr);
        } else goto handle_unusual;
        continue;
      // uint64 uncovered = 3;
      case 3:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::PROTOBUF_NAMESPACE_ID::uint8>(tag) == 24)) {
          uncovered_ = ::PROTOBUF_NAMESPACE_ID::internal::ReadVarint(&ptr);
          CHK_(ptr);
        } else goto handle_unusual;
        continue;
      default: {
      handle_unusual:
        if ((tag & 7) == 4 || tag == 0) {
          ctx->SetLastTag(tag);
          goto success;
        }
        ptr = UnknownFieldParse(tag, &_internal_metadata_, ptr, ctx);
        CHK_(ptr != nullptr);
        continue;
      }
    }  // switch
  }  // while
success:
  return ptr;
failure:
  ptr = nullptr;
  goto success;
#undef CHK_
}
#else  // GOOGLE_PROTOBUF_ENABLE_EXPERIMENTAL_PARSER
bool data::MergePartialFromCodedStream(
    ::PROTOBUF_NAMESPACE_ID::io::CodedInputStream* input) {
#define DO_(EXPRESSION) if (!PROTOBUF_PREDICT_TRUE(EXPRESSION)) goto failure
  ::PROTOBUF_NAMESPACE_ID::uint32 tag;
  // @@protoc_insertion_point(parse_start:dra.data)
  for (;;) {
    ::std::pair<::PROTOBUF_NAMESPACE_ID::uint32, bool> p = input->ReadTagWithCutoffNoLastTag(127u);
    tag = p.first;
    if (!p.second) goto handle_unusual;
    switch (::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::GetTagFieldNumber(tag)) {
      // uint64 time = 1;
      case 1: {
        if (static_cast< ::PROTOBUF_NAMESPACE_ID::uint8>(tag) == (8 & 0xFF)) {

          DO_((::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::ReadPrimitive<
                   ::PROTOBUF_NAMESPACE_ID::uint64, ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::TYPE_UINT64>(
                 input, &time_)));
        } else {
          goto handle_unusual;
        }
        break;
      }

      // string date = 2;
      case 2: {
        if (static_cast< ::PROTOBUF_NAMESPACE_ID::uint8>(tag) == (18 & 0xFF)) {
          DO_(::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::ReadString(
                input, this->mutable_date()));
          DO_(::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::VerifyUtf8String(
            this->date().data(), static_cast<int>(this->date().length()),
            ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::PARSE,
            "dra.data.date"));
        } else {
          goto handle_unusual;
        }
        break;
      }

      // uint64 uncovered = 3;
      case 3: {
        if (static_cast< ::PROTOBUF_NAMESPACE_ID::uint8>(tag) == (24 & 0xFF)) {

          DO_((::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::ReadPrimitive<
                   ::PROTOBUF_NAMESPACE_ID::uint64, ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::TYPE_UINT64>(
                 input, &uncovered_)));
        } else {
          goto handle_unusual;
        }
        break;
      }

      default: {
      handle_unusual:
        if (tag == 0) {
          goto success;
        }
        DO_(::PROTOBUF_NAMESPACE_ID::internal::WireFormat::SkipField(
              input, tag, _internal_metadata_.mutable_unknown_fields()));
        break;
      }
    }
  }
success:
  // @@protoc_insertion_point(parse_success:dra.data)
  return true;
failure:
  // @@protoc_insertion_point(parse_failure:dra.data)
  return false;
#undef DO_
}
#endif  // GOOGLE_PROTOBUF_ENABLE_EXPERIMENTAL_PARSER

void data::SerializeWithCachedSizes(
    ::PROTOBUF_NAMESPACE_ID::io::CodedOutputStream* output) const {
  // @@protoc_insertion_point(serialize_start:dra.data)
  ::PROTOBUF_NAMESPACE_ID::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

  // uint64 time = 1;
  if (this->time() != 0) {
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::WriteUInt64(1, this->time(), output);
  }

  // string date = 2;
  if (this->date().size() > 0) {
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::VerifyUtf8String(
      this->date().data(), static_cast<int>(this->date().length()),
      ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::SERIALIZE,
      "dra.data.date");
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::WriteStringMaybeAliased(
      2, this->date(), output);
  }

  // uint64 uncovered = 3;
  if (this->uncovered() != 0) {
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::WriteUInt64(3, this->uncovered(), output);
  }

  if (_internal_metadata_.have_unknown_fields()) {
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormat::SerializeUnknownFields(
        _internal_metadata_.unknown_fields(), output);
  }
  // @@protoc_insertion_point(serialize_end:dra.data)
}

::PROTOBUF_NAMESPACE_ID::uint8* data::InternalSerializeWithCachedSizesToArray(
    ::PROTOBUF_NAMESPACE_ID::uint8* target) const {
  // @@protoc_insertion_point(serialize_to_array_start:dra.data)
  ::PROTOBUF_NAMESPACE_ID::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

  // uint64 time = 1;
  if (this->time() != 0) {
    target = ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::WriteUInt64ToArray(1, this->time(), target);
  }

  // string date = 2;
  if (this->date().size() > 0) {
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::VerifyUtf8String(
      this->date().data(), static_cast<int>(this->date().length()),
      ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::SERIALIZE,
      "dra.data.date");
    target =
      ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::WriteStringToArray(
        2, this->date(), target);
  }

  // uint64 uncovered = 3;
  if (this->uncovered() != 0) {
    target = ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::WriteUInt64ToArray(3, this->uncovered(), target);
  }

  if (_internal_metadata_.have_unknown_fields()) {
    target = ::PROTOBUF_NAMESPACE_ID::internal::WireFormat::SerializeUnknownFieldsToArray(
        _internal_metadata_.unknown_fields(), target);
  }
  // @@protoc_insertion_point(serialize_to_array_end:dra.data)
  return target;
}

size_t data::ByteSizeLong() const {
// @@protoc_insertion_point(message_byte_size_start:dra.data)
  size_t total_size = 0;

  if (_internal_metadata_.have_unknown_fields()) {
    total_size +=
      ::PROTOBUF_NAMESPACE_ID::internal::WireFormat::ComputeUnknownFieldsSize(
        _internal_metadata_.unknown_fields());
  }
  ::PROTOBUF_NAMESPACE_ID::uint32 cached_has_bits = 0;
  // Prevent compiler warnings about cached_has_bits being unused
  (void) cached_has_bits;

  // string date = 2;
  if (this->date().size() > 0) {
    total_size += 1 +
      ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::StringSize(
        this->date());
  }

  // uint64 time = 1;
  if (this->time() != 0) {
    total_size += 1 +
      ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::UInt64Size(
        this->time());
  }

  // uint64 uncovered = 3;
  if (this->uncovered() != 0) {
    total_size += 1 +
      ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::UInt64Size(
        this->uncovered());
  }

  int cached_size = ::PROTOBUF_NAMESPACE_ID::internal::ToCachedSize(total_size);
  SetCachedSize(cached_size);
  return total_size;
}

void data::MergeFrom(const ::PROTOBUF_NAMESPACE_ID::Message& from) {
// @@protoc_insertion_point(generalized_merge_from_start:dra.data)
  GOOGLE_DCHECK_NE(&from, this);
  const data* source =
      ::PROTOBUF_NAMESPACE_ID::DynamicCastToGenerated<data>(
          &from);
  if (source == nullptr) {
  // @@protoc_insertion_point(generalized_merge_from_cast_fail:dra.data)
    ::PROTOBUF_NAMESPACE_ID::internal::ReflectionOps::Merge(from, this);
  } else {
  // @@protoc_insertion_point(generalized_merge_from_cast_success:dra.data)
    MergeFrom(*source);
  }
}

void data::MergeFrom(const data& from) {
// @@protoc_insertion_point(class_specific_merge_from_start:dra.data)
  GOOGLE_DCHECK_NE(&from, this);
  _internal_metadata_.MergeFrom(from._internal_metadata_);
  ::PROTOBUF_NAMESPACE_ID::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

  if (from.date().size() > 0) {

    date_.AssignWithDefault(&::PROTOBUF_NAMESPACE_ID::internal::GetEmptyStringAlreadyInited(), from.date_);
  }
  if (from.time() != 0) {
    set_time(from.time());
  }
  if (from.uncovered() != 0) {
    set_uncovered(from.uncovered());
  }
}

void data::CopyFrom(const ::PROTOBUF_NAMESPACE_ID::Message& from) {
// @@protoc_insertion_point(generalized_copy_from_start:dra.data)
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

void data::CopyFrom(const data& from) {
// @@protoc_insertion_point(class_specific_copy_from_start:dra.data)
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

bool data::IsInitialized() const {
  return true;
}

void data::Swap(data* other) {
  if (other == this) return;
  InternalSwap(other);
}
void data::InternalSwap(data* other) {
  using std::swap;
  _internal_metadata_.Swap(&other->_internal_metadata_);
  date_.Swap(&other->date_, &::PROTOBUF_NAMESPACE_ID::internal::GetEmptyStringAlreadyInited(),
    GetArenaNoVirtual());
  swap(time_, other->time_);
  swap(uncovered_, other->uncovered_);
}

::PROTOBUF_NAMESPACE_ID::Metadata data::GetMetadata() const {
  return GetMetadataStatic();
}


// ===================================================================

all_data_AddressTimeEntry_DoNotUse::all_data_AddressTimeEntry_DoNotUse() {}
all_data_AddressTimeEntry_DoNotUse::all_data_AddressTimeEntry_DoNotUse(::PROTOBUF_NAMESPACE_ID::Arena* arena)
    : SuperType(arena) {}
void all_data_AddressTimeEntry_DoNotUse::MergeFrom(const all_data_AddressTimeEntry_DoNotUse& other) {
  MergeFromInternal(other);
}
::PROTOBUF_NAMESPACE_ID::Metadata all_data_AddressTimeEntry_DoNotUse::GetMetadata() const {
  return GetMetadataStatic();
}
void all_data_AddressTimeEntry_DoNotUse::MergeFrom(
    const ::PROTOBUF_NAMESPACE_ID::Message& other) {
  ::PROTOBUF_NAMESPACE_ID::Message::MergeFrom(other);
}


// ===================================================================

void all_data::InitAsDefaultInstance() {
}
class all_data::HasBitSetters {
 public:
};

#if !defined(_MSC_VER) || _MSC_VER >= 1900
const int all_data::kNameFieldNumber;
const int all_data::kAddressTimeFieldNumber;
const int all_data::kDataFieldNumber;
#endif  // !defined(_MSC_VER) || _MSC_VER >= 1900

all_data::all_data()
  : ::PROTOBUF_NAMESPACE_ID::Message(), _internal_metadata_(nullptr) {
  SharedCtor();
  // @@protoc_insertion_point(constructor:dra.all_data)
}
all_data::all_data(const all_data& from)
  : ::PROTOBUF_NAMESPACE_ID::Message(),
      _internal_metadata_(nullptr),
      data_(from.data_) {
  _internal_metadata_.MergeFrom(from._internal_metadata_);
  address_time_.MergeFrom(from.address_time_);
  name_.UnsafeSetDefault(&::PROTOBUF_NAMESPACE_ID::internal::GetEmptyStringAlreadyInited());
  if (from.name().size() > 0) {
    name_.AssignWithDefault(&::PROTOBUF_NAMESPACE_ID::internal::GetEmptyStringAlreadyInited(), from.name_);
  }
  // @@protoc_insertion_point(copy_constructor:dra.all_data)
}

void all_data::SharedCtor() {
  ::PROTOBUF_NAMESPACE_ID::internal::InitSCC(&scc_info_all_data_Data_2eproto.base);
  name_.UnsafeSetDefault(&::PROTOBUF_NAMESPACE_ID::internal::GetEmptyStringAlreadyInited());
}

all_data::~all_data() {
  // @@protoc_insertion_point(destructor:dra.all_data)
  SharedDtor();
}

void all_data::SharedDtor() {
  name_.DestroyNoArena(&::PROTOBUF_NAMESPACE_ID::internal::GetEmptyStringAlreadyInited());
}

void all_data::SetCachedSize(int size) const {
  _cached_size_.Set(size);
}
const all_data& all_data::default_instance() {
  ::PROTOBUF_NAMESPACE_ID::internal::InitSCC(&::scc_info_all_data_Data_2eproto.base);
  return *internal_default_instance();
}


void all_data::Clear() {
// @@protoc_insertion_point(message_clear_start:dra.all_data)
  ::PROTOBUF_NAMESPACE_ID::uint32 cached_has_bits = 0;
  // Prevent compiler warnings about cached_has_bits being unused
  (void) cached_has_bits;

  address_time_.Clear();
  data_.Clear();
  name_.ClearToEmptyNoArena(&::PROTOBUF_NAMESPACE_ID::internal::GetEmptyStringAlreadyInited());
  _internal_metadata_.Clear();
}

#if GOOGLE_PROTOBUF_ENABLE_EXPERIMENTAL_PARSER
const char* all_data::_InternalParse(const char* ptr, ::PROTOBUF_NAMESPACE_ID::internal::ParseContext* ctx) {
#define CHK_(x) if (PROTOBUF_PREDICT_FALSE(!(x))) goto failure
  while (!ctx->Done(&ptr)) {
    ::PROTOBUF_NAMESPACE_ID::uint32 tag;
    ptr = ::PROTOBUF_NAMESPACE_ID::internal::ReadTag(ptr, &tag);
    CHK_(ptr);
    switch (tag >> 3) {
      // string name = 1;
      case 1:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::PROTOBUF_NAMESPACE_ID::uint8>(tag) == 10)) {
          ptr = ::PROTOBUF_NAMESPACE_ID::internal::InlineGreedyStringParserUTF8(mutable_name(), ptr, ctx, "dra.all_data.name");
          CHK_(ptr);
        } else goto handle_unusual;
        continue;
      // map<uint64, uint64> address_time = 2;
      case 2:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::PROTOBUF_NAMESPACE_ID::uint8>(tag) == 18)) {
          ptr -= 1;
          do {
            ptr += 1;
            ptr = ctx->ParseMessage(&address_time_, ptr);
            CHK_(ptr);
            if (!ctx->DataAvailable(ptr)) break;
          } while (::PROTOBUF_NAMESPACE_ID::internal::UnalignedLoad<::PROTOBUF_NAMESPACE_ID::uint8>(ptr) == 18);
        } else goto handle_unusual;
        continue;
      // repeated .dra.data data = 3;
      case 3:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::PROTOBUF_NAMESPACE_ID::uint8>(tag) == 26)) {
          ptr -= 1;
          do {
            ptr += 1;
            ptr = ctx->ParseMessage(add_data(), ptr);
            CHK_(ptr);
            if (!ctx->DataAvailable(ptr)) break;
          } while (::PROTOBUF_NAMESPACE_ID::internal::UnalignedLoad<::PROTOBUF_NAMESPACE_ID::uint8>(ptr) == 26);
        } else goto handle_unusual;
        continue;
      default: {
      handle_unusual:
        if ((tag & 7) == 4 || tag == 0) {
          ctx->SetLastTag(tag);
          goto success;
        }
        ptr = UnknownFieldParse(tag, &_internal_metadata_, ptr, ctx);
        CHK_(ptr != nullptr);
        continue;
      }
    }  // switch
  }  // while
success:
  return ptr;
failure:
  ptr = nullptr;
  goto success;
#undef CHK_
}
#else  // GOOGLE_PROTOBUF_ENABLE_EXPERIMENTAL_PARSER
bool all_data::MergePartialFromCodedStream(
    ::PROTOBUF_NAMESPACE_ID::io::CodedInputStream* input) {
#define DO_(EXPRESSION) if (!PROTOBUF_PREDICT_TRUE(EXPRESSION)) goto failure
  ::PROTOBUF_NAMESPACE_ID::uint32 tag;
  // @@protoc_insertion_point(parse_start:dra.all_data)
  for (;;) {
    ::std::pair<::PROTOBUF_NAMESPACE_ID::uint32, bool> p = input->ReadTagWithCutoffNoLastTag(127u);
    tag = p.first;
    if (!p.second) goto handle_unusual;
    switch (::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::GetTagFieldNumber(tag)) {
      // string name = 1;
      case 1: {
        if (static_cast< ::PROTOBUF_NAMESPACE_ID::uint8>(tag) == (10 & 0xFF)) {
          DO_(::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::ReadString(
                input, this->mutable_name()));
          DO_(::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::VerifyUtf8String(
            this->name().data(), static_cast<int>(this->name().length()),
            ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::PARSE,
            "dra.all_data.name"));
        } else {
          goto handle_unusual;
        }
        break;
      }

      // map<uint64, uint64> address_time = 2;
      case 2: {
        if (static_cast< ::PROTOBUF_NAMESPACE_ID::uint8>(tag) == (18 & 0xFF)) {
          all_data_AddressTimeEntry_DoNotUse::Parser< ::PROTOBUF_NAMESPACE_ID::internal::MapField<
              all_data_AddressTimeEntry_DoNotUse,
              ::PROTOBUF_NAMESPACE_ID::uint64, ::PROTOBUF_NAMESPACE_ID::uint64,
              ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::TYPE_UINT64,
              ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::TYPE_UINT64,
              0 >,
            ::PROTOBUF_NAMESPACE_ID::Map< ::PROTOBUF_NAMESPACE_ID::uint64, ::PROTOBUF_NAMESPACE_ID::uint64 > > parser(&address_time_);
          DO_(::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::ReadMessageNoVirtual(
              input, &parser));
        } else {
          goto handle_unusual;
        }
        break;
      }

      // repeated .dra.data data = 3;
      case 3: {
        if (static_cast< ::PROTOBUF_NAMESPACE_ID::uint8>(tag) == (26 & 0xFF)) {
          DO_(::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::ReadMessage(
                input, add_data()));
        } else {
          goto handle_unusual;
        }
        break;
      }

      default: {
      handle_unusual:
        if (tag == 0) {
          goto success;
        }
        DO_(::PROTOBUF_NAMESPACE_ID::internal::WireFormat::SkipField(
              input, tag, _internal_metadata_.mutable_unknown_fields()));
        break;
      }
    }
  }
success:
  // @@protoc_insertion_point(parse_success:dra.all_data)
  return true;
failure:
  // @@protoc_insertion_point(parse_failure:dra.all_data)
  return false;
#undef DO_
}
#endif  // GOOGLE_PROTOBUF_ENABLE_EXPERIMENTAL_PARSER

void all_data::SerializeWithCachedSizes(
    ::PROTOBUF_NAMESPACE_ID::io::CodedOutputStream* output) const {
  // @@protoc_insertion_point(serialize_start:dra.all_data)
  ::PROTOBUF_NAMESPACE_ID::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

  // string name = 1;
  if (this->name().size() > 0) {
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::VerifyUtf8String(
      this->name().data(), static_cast<int>(this->name().length()),
      ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::SERIALIZE,
      "dra.all_data.name");
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::WriteStringMaybeAliased(
      1, this->name(), output);
  }

  // map<uint64, uint64> address_time = 2;
  if (!this->address_time().empty()) {
    typedef ::PROTOBUF_NAMESPACE_ID::Map< ::PROTOBUF_NAMESPACE_ID::uint64, ::PROTOBUF_NAMESPACE_ID::uint64 >::const_pointer
        ConstPtr;
    typedef ::PROTOBUF_NAMESPACE_ID::internal::SortItem< ::PROTOBUF_NAMESPACE_ID::uint64, ConstPtr > SortItem;
    typedef ::PROTOBUF_NAMESPACE_ID::internal::CompareByFirstField<SortItem> Less;

    if (output->IsSerializationDeterministic() &&
        this->address_time().size() > 1) {
      ::std::unique_ptr<SortItem[]> items(
          new SortItem[this->address_time().size()]);
      typedef ::PROTOBUF_NAMESPACE_ID::Map< ::PROTOBUF_NAMESPACE_ID::uint64, ::PROTOBUF_NAMESPACE_ID::uint64 >::size_type size_type;
      size_type n = 0;
      for (::PROTOBUF_NAMESPACE_ID::Map< ::PROTOBUF_NAMESPACE_ID::uint64, ::PROTOBUF_NAMESPACE_ID::uint64 >::const_iterator
          it = this->address_time().begin();
          it != this->address_time().end(); ++it, ++n) {
        items[static_cast<ptrdiff_t>(n)] = SortItem(&*it);
      }
      ::std::sort(&items[0], &items[static_cast<ptrdiff_t>(n)], Less());
      for (size_type i = 0; i < n; i++) {
        all_data_AddressTimeEntry_DoNotUse::Funcs::SerializeToCodedStream(2, items[static_cast<ptrdiff_t>(i)].second->first, items[static_cast<ptrdiff_t>(i)].second->second, output);
      }
    } else {
      for (::PROTOBUF_NAMESPACE_ID::Map< ::PROTOBUF_NAMESPACE_ID::uint64, ::PROTOBUF_NAMESPACE_ID::uint64 >::const_iterator
          it = this->address_time().begin();
          it != this->address_time().end(); ++it) {
        all_data_AddressTimeEntry_DoNotUse::Funcs::SerializeToCodedStream(2, it->first, it->second, output);
      }
    }
  }

  // repeated .dra.data data = 3;
  for (unsigned int i = 0,
      n = static_cast<unsigned int>(this->data_size()); i < n; i++) {
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::WriteMessageMaybeToArray(
      3,
      this->data(static_cast<int>(i)),
      output);
  }

  if (_internal_metadata_.have_unknown_fields()) {
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormat::SerializeUnknownFields(
        _internal_metadata_.unknown_fields(), output);
  }
  // @@protoc_insertion_point(serialize_end:dra.all_data)
}

::PROTOBUF_NAMESPACE_ID::uint8* all_data::InternalSerializeWithCachedSizesToArray(
    ::PROTOBUF_NAMESPACE_ID::uint8* target) const {
  // @@protoc_insertion_point(serialize_to_array_start:dra.all_data)
  ::PROTOBUF_NAMESPACE_ID::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

  // string name = 1;
  if (this->name().size() > 0) {
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::VerifyUtf8String(
      this->name().data(), static_cast<int>(this->name().length()),
      ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::SERIALIZE,
      "dra.all_data.name");
    target =
      ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::WriteStringToArray(
        1, this->name(), target);
  }

  // map<uint64, uint64> address_time = 2;
  if (!this->address_time().empty()) {
    typedef ::PROTOBUF_NAMESPACE_ID::Map< ::PROTOBUF_NAMESPACE_ID::uint64, ::PROTOBUF_NAMESPACE_ID::uint64 >::const_pointer
        ConstPtr;
    typedef ::PROTOBUF_NAMESPACE_ID::internal::SortItem< ::PROTOBUF_NAMESPACE_ID::uint64, ConstPtr > SortItem;
    typedef ::PROTOBUF_NAMESPACE_ID::internal::CompareByFirstField<SortItem> Less;

    if (false &&
        this->address_time().size() > 1) {
      ::std::unique_ptr<SortItem[]> items(
          new SortItem[this->address_time().size()]);
      typedef ::PROTOBUF_NAMESPACE_ID::Map< ::PROTOBUF_NAMESPACE_ID::uint64, ::PROTOBUF_NAMESPACE_ID::uint64 >::size_type size_type;
      size_type n = 0;
      for (::PROTOBUF_NAMESPACE_ID::Map< ::PROTOBUF_NAMESPACE_ID::uint64, ::PROTOBUF_NAMESPACE_ID::uint64 >::const_iterator
          it = this->address_time().begin();
          it != this->address_time().end(); ++it, ++n) {
        items[static_cast<ptrdiff_t>(n)] = SortItem(&*it);
      }
      ::std::sort(&items[0], &items[static_cast<ptrdiff_t>(n)], Less());
      for (size_type i = 0; i < n; i++) {
        target = all_data_AddressTimeEntry_DoNotUse::Funcs::SerializeToArray(2, items[static_cast<ptrdiff_t>(i)].second->first, items[static_cast<ptrdiff_t>(i)].second->second, target);
      }
    } else {
      for (::PROTOBUF_NAMESPACE_ID::Map< ::PROTOBUF_NAMESPACE_ID::uint64, ::PROTOBUF_NAMESPACE_ID::uint64 >::const_iterator
          it = this->address_time().begin();
          it != this->address_time().end(); ++it) {
        target = all_data_AddressTimeEntry_DoNotUse::Funcs::SerializeToArray(2, it->first, it->second, target);
      }
    }
  }

  // repeated .dra.data data = 3;
  for (unsigned int i = 0,
      n = static_cast<unsigned int>(this->data_size()); i < n; i++) {
    target = ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::
      InternalWriteMessageToArray(
        3, this->data(static_cast<int>(i)), target);
  }

  if (_internal_metadata_.have_unknown_fields()) {
    target = ::PROTOBUF_NAMESPACE_ID::internal::WireFormat::SerializeUnknownFieldsToArray(
        _internal_metadata_.unknown_fields(), target);
  }
  // @@protoc_insertion_point(serialize_to_array_end:dra.all_data)
  return target;
}

size_t all_data::ByteSizeLong() const {
// @@protoc_insertion_point(message_byte_size_start:dra.all_data)
  size_t total_size = 0;

  if (_internal_metadata_.have_unknown_fields()) {
    total_size +=
      ::PROTOBUF_NAMESPACE_ID::internal::WireFormat::ComputeUnknownFieldsSize(
        _internal_metadata_.unknown_fields());
  }
  ::PROTOBUF_NAMESPACE_ID::uint32 cached_has_bits = 0;
  // Prevent compiler warnings about cached_has_bits being unused
  (void) cached_has_bits;

  // map<uint64, uint64> address_time = 2;
  total_size += 1 *
      ::PROTOBUF_NAMESPACE_ID::internal::FromIntSize(this->address_time_size());
  for (::PROTOBUF_NAMESPACE_ID::Map< ::PROTOBUF_NAMESPACE_ID::uint64, ::PROTOBUF_NAMESPACE_ID::uint64 >::const_iterator
      it = this->address_time().begin();
      it != this->address_time().end(); ++it) {
    total_size += all_data_AddressTimeEntry_DoNotUse::Funcs::ByteSizeLong(it->first, it->second);
  }

  // repeated .dra.data data = 3;
  {
    unsigned int count = static_cast<unsigned int>(this->data_size());
    total_size += 1UL * count;
    for (unsigned int i = 0; i < count; i++) {
      total_size +=
        ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::MessageSize(
          this->data(static_cast<int>(i)));
    }
  }

  // string name = 1;
  if (this->name().size() > 0) {
    total_size += 1 +
      ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::StringSize(
        this->name());
  }

  int cached_size = ::PROTOBUF_NAMESPACE_ID::internal::ToCachedSize(total_size);
  SetCachedSize(cached_size);
  return total_size;
}

void all_data::MergeFrom(const ::PROTOBUF_NAMESPACE_ID::Message& from) {
// @@protoc_insertion_point(generalized_merge_from_start:dra.all_data)
  GOOGLE_DCHECK_NE(&from, this);
  const all_data* source =
      ::PROTOBUF_NAMESPACE_ID::DynamicCastToGenerated<all_data>(
          &from);
  if (source == nullptr) {
  // @@protoc_insertion_point(generalized_merge_from_cast_fail:dra.all_data)
    ::PROTOBUF_NAMESPACE_ID::internal::ReflectionOps::Merge(from, this);
  } else {
  // @@protoc_insertion_point(generalized_merge_from_cast_success:dra.all_data)
    MergeFrom(*source);
  }
}

void all_data::MergeFrom(const all_data& from) {
// @@protoc_insertion_point(class_specific_merge_from_start:dra.all_data)
  GOOGLE_DCHECK_NE(&from, this);
  _internal_metadata_.MergeFrom(from._internal_metadata_);
  ::PROTOBUF_NAMESPACE_ID::uint32 cached_has_bits = 0;
  (void) cached_has_bits;

  address_time_.MergeFrom(from.address_time_);
  data_.MergeFrom(from.data_);
  if (from.name().size() > 0) {

    name_.AssignWithDefault(&::PROTOBUF_NAMESPACE_ID::internal::GetEmptyStringAlreadyInited(), from.name_);
  }
}

void all_data::CopyFrom(const ::PROTOBUF_NAMESPACE_ID::Message& from) {
// @@protoc_insertion_point(generalized_copy_from_start:dra.all_data)
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

void all_data::CopyFrom(const all_data& from) {
// @@protoc_insertion_point(class_specific_copy_from_start:dra.all_data)
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

bool all_data::IsInitialized() const {
  return true;
}

void all_data::Swap(all_data* other) {
  if (other == this) return;
  InternalSwap(other);
}
void all_data::InternalSwap(all_data* other) {
  using std::swap;
  _internal_metadata_.Swap(&other->_internal_metadata_);
  address_time_.Swap(&other->address_time_);
  CastToBase(&data_)->InternalSwap(CastToBase(&other->data_));
  name_.Swap(&other->name_, &::PROTOBUF_NAMESPACE_ID::internal::GetEmptyStringAlreadyInited(),
    GetArenaNoVirtual());
}

::PROTOBUF_NAMESPACE_ID::Metadata all_data::GetMetadata() const {
  return GetMetadataStatic();
}


// @@protoc_insertion_point(namespace_scope)
}  // namespace dra
PROTOBUF_NAMESPACE_OPEN
template<> PROTOBUF_NOINLINE ::dra::data* Arena::CreateMaybeMessage< ::dra::data >(Arena* arena) {
  return Arena::CreateInternal< ::dra::data >(arena);
}
template<> PROTOBUF_NOINLINE ::dra::all_data_AddressTimeEntry_DoNotUse* Arena::CreateMaybeMessage< ::dra::all_data_AddressTimeEntry_DoNotUse >(Arena* arena) {
  return Arena::CreateInternal< ::dra::all_data_AddressTimeEntry_DoNotUse >(arena);
}
template<> PROTOBUF_NOINLINE ::dra::all_data* Arena::CreateMaybeMessage< ::dra::all_data >(Arena* arena) {
  return Arena::CreateInternal< ::dra::all_data >(arena);
}
PROTOBUF_NAMESPACE_CLOSE

// @@protoc_insertion_point(global_scope)
#include <google/protobuf/port_undef.inc>