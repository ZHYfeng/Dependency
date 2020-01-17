# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: Task.proto

from google.protobuf.internal import enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='Task.proto',
  package='dra',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=b'\n\nTask.proto\x12\x03\x64ra\"\x94\x02\n\x0brunTimeData\x12\x10\n\x08priority\x18\x01 \x01(\r\x12\x19\n\x11\x63ondition_address\x18\x02 \x01(\r\x12\x0f\n\x07\x61\x64\x64ress\x18\x04 \x01(\r\x12\x1c\n\x14right_branch_address\x18\x06 \x03(\r\x12$\n\x0btask_status\x18\n \x01(\x0e\x32\x0f.dra.taskStatus\x12\x0f\n\x07program\x18\x0c \x01(\x0c\x12\x0b\n\x03idx\x18\r \x01(\r\x12\x16\n\x0ercursive_count\x18\x0e \x01(\r\x12\x16\n\x0e\x63heckCondition\x18\x15 \x01(\x08\x12\x14\n\x0c\x63heckAddress\x18\x16 \x01(\x08\x12\x1f\n\x17\x63heckRightBranchAddress\x18\x17 \x01(\x08\"\x95\x03\n\x0fTaskRunTimeData\x12\x0c\n\x04hash\x18\x01 \x01(\t\x12\x0f\n\x07program\x18\x02 \x01(\x0c\x12\x11\n\twrite_idx\x18\x05 \x01(\r\x12\x15\n\rcondition_idx\x18\x06 \x01(\r\x12\x1b\n\x13\x63heck_write_address\x18\n \x01(\x08\x12\x45\n\x11uncovered_address\x18\x15 \x03(\x0b\x32*.dra.TaskRunTimeData.UncoveredAddressEntry\x12\x41\n\x0f\x63overed_address\x18\x17 \x03(\x0b\x32(.dra.TaskRunTimeData.CoveredAddressEntry\x1aI\n\x15UncoveredAddressEntry\x12\x0b\n\x03key\x18\x01 \x01(\r\x12\x1f\n\x05value\x18\x02 \x01(\x0b\x32\x10.dra.runTimeData:\x02\x38\x01\x1aG\n\x13\x43overedAddressEntry\x12\x0b\n\x03key\x18\x01 \x01(\r\x12\x1f\n\x05value\x18\x02 \x01(\x0b\x32\x10.dra.runTimeData:\x02\x38\x01\"\xc3\x04\n\x04Task\x12\x0b\n\x03sig\x18\x01 \x01(\t\x12\r\n\x05index\x18\x02 \x01(\r\x12\x0f\n\x07program\x18\x03 \x01(\x0c\x12\x0c\n\x04kind\x18\x04 \x01(\r\x12\x10\n\x08priority\x18\x05 \x01(\x05\x12\x0c\n\x04hash\x18\x06 \x01(\t\x12\r\n\x05\x63ount\x18\x07 \x01(\r\x12\x11\n\twrite_sig\x18\x0b \x01(\t\x12\x13\n\x0bwrite_index\x18\x0c \x01(\r\x12\x15\n\rwrite_program\x18\r \x01(\x0c\x12\x15\n\rwrite_address\x18\x0e \x01(\r\x12$\n\x0btask_status\x18\x18 \x01(\x0e\x32\x0f.dra.taskStatus\x12\x1b\n\x13\x63heck_write_address\x18\x19 \x01(\x08\x12:\n\x11uncovered_address\x18\x15 \x03(\x0b\x32\x1f.dra.Task.UncoveredAddressEntry\x12\x36\n\x0f\x63overed_address\x18\x17 \x03(\x0b\x32\x1d.dra.Task.CoveredAddressEntry\x12\x30\n\x12task_run_time_data\x18\x1f \x03(\x0b\x32\x14.dra.TaskRunTimeData\x1aI\n\x15UncoveredAddressEntry\x12\x0b\n\x03key\x18\x01 \x01(\r\x12\x1f\n\x05value\x18\x02 \x01(\x0b\x32\x10.dra.runTimeData:\x02\x38\x01\x1aG\n\x13\x43overedAddressEntry\x12\x0b\n\x03key\x18\x01 \x01(\r\x12\x1f\n\x05value\x18\x02 \x01(\x0b\x32\x10.dra.runTimeData:\x02\x38\x01\"\xb7\x01\n\x05Tasks\x12\x0c\n\x04name\x18) \x01(\t\x12\x1b\n\x04kind\x18* \x01(\x0e\x32\r.dra.TaskKind\x12)\n\x08task_map\x18\x01 \x03(\x0b\x32\x17.dra.Tasks.TaskMapEntry\x12\x1d\n\ntask_array\x18\x02 \x03(\x0b\x32\t.dra.Task\x1a\x39\n\x0cTaskMapEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12\x18\n\x05value\x18\x02 \x01(\x0b\x32\t.dra.Task:\x02\x38\x01*\x9a\x01\n\ntaskStatus\x12\x0c\n\x08untested\x10\x00\x12\x12\n\x0eunstable_write\x10\x0b\x12\x16\n\x12unstable_condition\x10\x0c\x12\x13\n\x0funstable_insert\x10\r\x12\x0c\n\x08unstable\x10\x0e\x12\n\n\x06tested\x10\x15\x12\x0b\n\x07\x63overed\x10\x16\x12\r\n\trecursive\x10\x1f\x12\x07\n\x03out\x10 **\n\x08TaskKind\x12\n\n\x06Normal\x10\x00\x12\x08\n\x04High\x10\x01\x12\x08\n\x04\x42oot\x10\x02\x62\x06proto3'
)

_TASKSTATUS = _descriptor.EnumDescriptor(
  name='taskStatus',
  full_name='dra.taskStatus',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='untested', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='unstable_write', index=1, number=11,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='unstable_condition', index=2, number=12,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='unstable_insert', index=3, number=13,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='unstable', index=4, number=14,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='tested', index=5, number=21,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='covered', index=6, number=22,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='recursive', index=7, number=31,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='out', index=8, number=32,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=1475,
  serialized_end=1629,
)
_sym_db.RegisterEnumDescriptor(_TASKSTATUS)

taskStatus = enum_type_wrapper.EnumTypeWrapper(_TASKSTATUS)
_TASKKIND = _descriptor.EnumDescriptor(
  name='TaskKind',
  full_name='dra.TaskKind',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='Normal', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='High', index=1, number=1,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='Boot', index=2, number=2,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=1631,
  serialized_end=1673,
)
_sym_db.RegisterEnumDescriptor(_TASKKIND)

TaskKind = enum_type_wrapper.EnumTypeWrapper(_TASKKIND)
untested = 0
unstable_write = 11
unstable_condition = 12
unstable_insert = 13
unstable = 14
tested = 21
covered = 22
recursive = 31
out = 32
Normal = 0
High = 1
Boot = 2



_RUNTIMEDATA = _descriptor.Descriptor(
  name='runTimeData',
  full_name='dra.runTimeData',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='priority', full_name='dra.runTimeData.priority', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='condition_address', full_name='dra.runTimeData.condition_address', index=1,
      number=2, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='address', full_name='dra.runTimeData.address', index=2,
      number=4, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='right_branch_address', full_name='dra.runTimeData.right_branch_address', index=3,
      number=6, type=13, cpp_type=3, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='task_status', full_name='dra.runTimeData.task_status', index=4,
      number=10, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='program', full_name='dra.runTimeData.program', index=5,
      number=12, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=b"",
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='idx', full_name='dra.runTimeData.idx', index=6,
      number=13, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='rcursive_count', full_name='dra.runTimeData.rcursive_count', index=7,
      number=14, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='checkCondition', full_name='dra.runTimeData.checkCondition', index=8,
      number=21, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='checkAddress', full_name='dra.runTimeData.checkAddress', index=9,
      number=22, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='checkRightBranchAddress', full_name='dra.runTimeData.checkRightBranchAddress', index=10,
      number=23, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=20,
  serialized_end=296,
)


_TASKRUNTIMEDATA_UNCOVEREDADDRESSENTRY = _descriptor.Descriptor(
  name='UncoveredAddressEntry',
  full_name='dra.TaskRunTimeData.UncoveredAddressEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dra.TaskRunTimeData.UncoveredAddressEntry.key', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='dra.TaskRunTimeData.UncoveredAddressEntry.value', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=558,
  serialized_end=631,
)

_TASKRUNTIMEDATA_COVEREDADDRESSENTRY = _descriptor.Descriptor(
  name='CoveredAddressEntry',
  full_name='dra.TaskRunTimeData.CoveredAddressEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dra.TaskRunTimeData.CoveredAddressEntry.key', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='dra.TaskRunTimeData.CoveredAddressEntry.value', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=633,
  serialized_end=704,
)

_TASKRUNTIMEDATA = _descriptor.Descriptor(
  name='TaskRunTimeData',
  full_name='dra.TaskRunTimeData',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='hash', full_name='dra.TaskRunTimeData.hash', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='program', full_name='dra.TaskRunTimeData.program', index=1,
      number=2, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=b"",
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='write_idx', full_name='dra.TaskRunTimeData.write_idx', index=2,
      number=5, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='condition_idx', full_name='dra.TaskRunTimeData.condition_idx', index=3,
      number=6, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='check_write_address', full_name='dra.TaskRunTimeData.check_write_address', index=4,
      number=10, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='uncovered_address', full_name='dra.TaskRunTimeData.uncovered_address', index=5,
      number=21, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='covered_address', full_name='dra.TaskRunTimeData.covered_address', index=6,
      number=23, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_TASKRUNTIMEDATA_UNCOVEREDADDRESSENTRY, _TASKRUNTIMEDATA_COVEREDADDRESSENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=299,
  serialized_end=704,
)


_TASK_UNCOVEREDADDRESSENTRY = _descriptor.Descriptor(
  name='UncoveredAddressEntry',
  full_name='dra.Task.UncoveredAddressEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dra.Task.UncoveredAddressEntry.key', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='dra.Task.UncoveredAddressEntry.value', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=558,
  serialized_end=631,
)

_TASK_COVEREDADDRESSENTRY = _descriptor.Descriptor(
  name='CoveredAddressEntry',
  full_name='dra.Task.CoveredAddressEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dra.Task.CoveredAddressEntry.key', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='dra.Task.CoveredAddressEntry.value', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=633,
  serialized_end=704,
)

_TASK = _descriptor.Descriptor(
  name='Task',
  full_name='dra.Task',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='sig', full_name='dra.Task.sig', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='index', full_name='dra.Task.index', index=1,
      number=2, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='program', full_name='dra.Task.program', index=2,
      number=3, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=b"",
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='kind', full_name='dra.Task.kind', index=3,
      number=4, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='priority', full_name='dra.Task.priority', index=4,
      number=5, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='hash', full_name='dra.Task.hash', index=5,
      number=6, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='count', full_name='dra.Task.count', index=6,
      number=7, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='write_sig', full_name='dra.Task.write_sig', index=7,
      number=11, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='write_index', full_name='dra.Task.write_index', index=8,
      number=12, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='write_program', full_name='dra.Task.write_program', index=9,
      number=13, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=b"",
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='write_address', full_name='dra.Task.write_address', index=10,
      number=14, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='task_status', full_name='dra.Task.task_status', index=11,
      number=24, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='check_write_address', full_name='dra.Task.check_write_address', index=12,
      number=25, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='uncovered_address', full_name='dra.Task.uncovered_address', index=13,
      number=21, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='covered_address', full_name='dra.Task.covered_address', index=14,
      number=23, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='task_run_time_data', full_name='dra.Task.task_run_time_data', index=15,
      number=31, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_TASK_UNCOVEREDADDRESSENTRY, _TASK_COVEREDADDRESSENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=707,
  serialized_end=1286,
)


_TASKS_TASKMAPENTRY = _descriptor.Descriptor(
  name='TaskMapEntry',
  full_name='dra.Tasks.TaskMapEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dra.Tasks.TaskMapEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='dra.Tasks.TaskMapEntry.value', index=1,
      number=2, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=b'8\001',
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1415,
  serialized_end=1472,
)

_TASKS = _descriptor.Descriptor(
  name='Tasks',
  full_name='dra.Tasks',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='dra.Tasks.name', index=0,
      number=41, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='kind', full_name='dra.Tasks.kind', index=1,
      number=42, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='task_map', full_name='dra.Tasks.task_map', index=2,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='task_array', full_name='dra.Tasks.task_array', index=3,
      number=2, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_TASKS_TASKMAPENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1289,
  serialized_end=1472,
)

_RUNTIMEDATA.fields_by_name['task_status'].enum_type = _TASKSTATUS
_TASKRUNTIMEDATA_UNCOVEREDADDRESSENTRY.fields_by_name['value'].message_type = _RUNTIMEDATA
_TASKRUNTIMEDATA_UNCOVEREDADDRESSENTRY.containing_type = _TASKRUNTIMEDATA
_TASKRUNTIMEDATA_COVEREDADDRESSENTRY.fields_by_name['value'].message_type = _RUNTIMEDATA
_TASKRUNTIMEDATA_COVEREDADDRESSENTRY.containing_type = _TASKRUNTIMEDATA
_TASKRUNTIMEDATA.fields_by_name['uncovered_address'].message_type = _TASKRUNTIMEDATA_UNCOVEREDADDRESSENTRY
_TASKRUNTIMEDATA.fields_by_name['covered_address'].message_type = _TASKRUNTIMEDATA_COVEREDADDRESSENTRY
_TASK_UNCOVEREDADDRESSENTRY.fields_by_name['value'].message_type = _RUNTIMEDATA
_TASK_UNCOVEREDADDRESSENTRY.containing_type = _TASK
_TASK_COVEREDADDRESSENTRY.fields_by_name['value'].message_type = _RUNTIMEDATA
_TASK_COVEREDADDRESSENTRY.containing_type = _TASK
_TASK.fields_by_name['task_status'].enum_type = _TASKSTATUS
_TASK.fields_by_name['uncovered_address'].message_type = _TASK_UNCOVEREDADDRESSENTRY
_TASK.fields_by_name['covered_address'].message_type = _TASK_COVEREDADDRESSENTRY
_TASK.fields_by_name['task_run_time_data'].message_type = _TASKRUNTIMEDATA
_TASKS_TASKMAPENTRY.fields_by_name['value'].message_type = _TASK
_TASKS_TASKMAPENTRY.containing_type = _TASKS
_TASKS.fields_by_name['kind'].enum_type = _TASKKIND
_TASKS.fields_by_name['task_map'].message_type = _TASKS_TASKMAPENTRY
_TASKS.fields_by_name['task_array'].message_type = _TASK
DESCRIPTOR.message_types_by_name['runTimeData'] = _RUNTIMEDATA
DESCRIPTOR.message_types_by_name['TaskRunTimeData'] = _TASKRUNTIMEDATA
DESCRIPTOR.message_types_by_name['Task'] = _TASK
DESCRIPTOR.message_types_by_name['Tasks'] = _TASKS
DESCRIPTOR.enum_types_by_name['taskStatus'] = _TASKSTATUS
DESCRIPTOR.enum_types_by_name['TaskKind'] = _TASKKIND
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

runTimeData = _reflection.GeneratedProtocolMessageType('runTimeData', (_message.Message,), {
  'DESCRIPTOR' : _RUNTIMEDATA,
  '__module__' : 'Task_pb2'
  # @@protoc_insertion_point(class_scope:dra.runTimeData)
  })
_sym_db.RegisterMessage(runTimeData)

TaskRunTimeData = _reflection.GeneratedProtocolMessageType('TaskRunTimeData', (_message.Message,), {

  'UncoveredAddressEntry' : _reflection.GeneratedProtocolMessageType('UncoveredAddressEntry', (_message.Message,), {
    'DESCRIPTOR' : _TASKRUNTIMEDATA_UNCOVEREDADDRESSENTRY,
    '__module__' : 'Task_pb2'
    # @@protoc_insertion_point(class_scope:dra.TaskRunTimeData.UncoveredAddressEntry)
    })
  ,

  'CoveredAddressEntry' : _reflection.GeneratedProtocolMessageType('CoveredAddressEntry', (_message.Message,), {
    'DESCRIPTOR' : _TASKRUNTIMEDATA_COVEREDADDRESSENTRY,
    '__module__' : 'Task_pb2'
    # @@protoc_insertion_point(class_scope:dra.TaskRunTimeData.CoveredAddressEntry)
    })
  ,
  'DESCRIPTOR' : _TASKRUNTIMEDATA,
  '__module__' : 'Task_pb2'
  # @@protoc_insertion_point(class_scope:dra.TaskRunTimeData)
  })
_sym_db.RegisterMessage(TaskRunTimeData)
_sym_db.RegisterMessage(TaskRunTimeData.UncoveredAddressEntry)
_sym_db.RegisterMessage(TaskRunTimeData.CoveredAddressEntry)

Task = _reflection.GeneratedProtocolMessageType('Task', (_message.Message,), {

  'UncoveredAddressEntry' : _reflection.GeneratedProtocolMessageType('UncoveredAddressEntry', (_message.Message,), {
    'DESCRIPTOR' : _TASK_UNCOVEREDADDRESSENTRY,
    '__module__' : 'Task_pb2'
    # @@protoc_insertion_point(class_scope:dra.Task.UncoveredAddressEntry)
    })
  ,

  'CoveredAddressEntry' : _reflection.GeneratedProtocolMessageType('CoveredAddressEntry', (_message.Message,), {
    'DESCRIPTOR' : _TASK_COVEREDADDRESSENTRY,
    '__module__' : 'Task_pb2'
    # @@protoc_insertion_point(class_scope:dra.Task.CoveredAddressEntry)
    })
  ,
  'DESCRIPTOR' : _TASK,
  '__module__' : 'Task_pb2'
  # @@protoc_insertion_point(class_scope:dra.Task)
  })
_sym_db.RegisterMessage(Task)
_sym_db.RegisterMessage(Task.UncoveredAddressEntry)
_sym_db.RegisterMessage(Task.CoveredAddressEntry)

Tasks = _reflection.GeneratedProtocolMessageType('Tasks', (_message.Message,), {

  'TaskMapEntry' : _reflection.GeneratedProtocolMessageType('TaskMapEntry', (_message.Message,), {
    'DESCRIPTOR' : _TASKS_TASKMAPENTRY,
    '__module__' : 'Task_pb2'
    # @@protoc_insertion_point(class_scope:dra.Tasks.TaskMapEntry)
    })
  ,
  'DESCRIPTOR' : _TASKS,
  '__module__' : 'Task_pb2'
  # @@protoc_insertion_point(class_scope:dra.Tasks)
  })
_sym_db.RegisterMessage(Tasks)
_sym_db.RegisterMessage(Tasks.TaskMapEntry)


_TASKRUNTIMEDATA_UNCOVEREDADDRESSENTRY._options = None
_TASKRUNTIMEDATA_COVEREDADDRESSENTRY._options = None
_TASK_UNCOVEREDADDRESSENTRY._options = None
_TASK_COVEREDADDRESSENTRY._options = None
_TASKS_TASKMAPENTRY._options = None
# @@protoc_insertion_point(module_scope)
