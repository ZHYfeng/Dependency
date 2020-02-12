# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: Input.proto

from google.protobuf.internal import enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='Input.proto',
  package='dra',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=b'\n\x0bInput.proto\x12\x03\x64ra\"l\n\x04\x43\x61ll\x12\x0b\n\x03idx\x18\x01 \x01(\r\x12\'\n\x07\x61\x64\x64ress\x18\x02 \x03(\x0b\x32\x16.dra.Call.AddressEntry\x1a.\n\x0c\x41\x64\x64ressEntry\x12\x0b\n\x03key\x18\x01 \x01(\r\x12\r\n\x05value\x18\x02 \x01(\r:\x02\x38\x01\"\xd8\x03\n\x05Input\x12\x0b\n\x03sig\x18\x0b \x01(\t\x12\x0f\n\x07program\x18\x0c \x01(\x0c\x12\"\n\x04\x63\x61ll\x18\r \x03(\x0b\x32\x14.dra.Input.CallEntry\x12\x19\n\x05paths\x18\x01 \x03(\x0b\x32\n.dra.Paths\x12\x0e\n\x06stable\x18\x0e \x01(\r\x12\r\n\x05total\x18\x0f \x01(\r\x12\x1e\n\x04stat\x18\x15 \x01(\x0e\x32\x10.dra.FuzzingStat\x12;\n\x11uncovered_address\x18\x16 \x03(\x0b\x32 .dra.Input.UncoveredAddressEntry\x12\x33\n\rwrite_address\x18\x19 \x03(\x0b\x32\x1c.dra.Input.WriteAddressEntry\x12\x1b\n\x13program_before_mini\x18\x1e \x01(\x0c\x1a\x36\n\tCallEntry\x12\x0b\n\x03key\x18\x01 \x01(\r\x12\x18\n\x05value\x18\x02 \x01(\x0b\x32\t.dra.Call:\x02\x38\x01\x1a\x37\n\x15UncoveredAddressEntry\x12\x0b\n\x03key\x18\x01 \x01(\r\x12\r\n\x05value\x18\x02 \x01(\r:\x02\x38\x01\x1a\x33\n\x11WriteAddressEntry\x12\x0b\n\x03key\x18\x01 \x01(\r\x12\r\n\x05value\x18\x02 \x01(\r:\x02\x38\x01\"#\n\x06Inputs\x12\x19\n\x05input\x18\x01 \x03(\x0b\x32\n.dra.Input\"\x17\n\x04Path\x12\x0f\n\x07\x61\x64\x64ress\x18\x01 \x03(\r\"c\n\x05Paths\x12\"\n\x04path\x18\x01 \x03(\x0b\x32\x14.dra.Paths.PathEntry\x1a\x36\n\tPathEntry\x12\x0b\n\x03key\x18\x01 \x01(\r\x12\x18\n\x05value\x18\x02 \x01(\x0b\x32\t.dra.Path:\x02\x38\x01\"\xb2\x01\n\rUnstableInput\x12\x0b\n\x03sig\x18\x01 \x01(\t\x12\x0f\n\x07program\x18\x02 \x01(\x0c\x12!\n\runstable_path\x18\x0c \x03(\x0b\x32\n.dra.Paths\x12\x30\n\x07\x61\x64\x64ress\x18\r \x03(\x0b\x32\x1f.dra.UnstableInput.AddressEntry\x1a.\n\x0c\x41\x64\x64ressEntry\x12\x0b\n\x03key\x18\x01 \x01(\r\x12\r\n\x05value\x18\x02 \x01(\r:\x02\x38\x01\"\x9a\x01\n\x0eUnstableInputs\x12>\n\x0eunstable_input\x18\x01 \x03(\x0b\x32&.dra.UnstableInputs.UnstableInputEntry\x1aH\n\x12UnstableInputEntry\x12\x0b\n\x03key\x18\x01 \x01(\t\x12!\n\x05value\x18\x02 \x01(\x0b\x32\x12.dra.UnstableInput:\x02\x38\x01*\xd3\x01\n\x0b\x46uzzingStat\x12\x10\n\x0cStatGenerate\x10\x00\x12\x0c\n\x08StatFuzz\x10\x01\x12\x11\n\rStatCandidate\x10\x02\x12\x0e\n\nStatTriage\x10\x03\x12\x10\n\x0cStatMinimize\x10\x04\x12\r\n\tStatSmash\x10\x05\x12\x0c\n\x08StatHint\x10\x06\x12\x0c\n\x08StatSeed\x10\x07\x12\x12\n\x0eStatDependency\x10\x08\x12\x16\n\x12StatDependencyBoot\x10\t\x12\x18\n\x0bStatDefault\x10\xff\xff\xff\xff\xff\xff\xff\xff\xff\x01\x62\x06proto3'
)

_FUZZINGSTAT = _descriptor.EnumDescriptor(
  name='FuzzingStat',
  full_name='dra.FuzzingStat',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='StatGenerate', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='StatFuzz', index=1, number=1,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='StatCandidate', index=2, number=2,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='StatTriage', index=3, number=3,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='StatMinimize', index=4, number=4,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='StatSmash', index=5, number=5,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='StatHint', index=6, number=6,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='StatSeed', index=7, number=7,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='StatDependency', index=8, number=8,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='StatDependencyBoot', index=9, number=9,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='StatDefault', index=10, number=-1,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=1107,
  serialized_end=1318,
)
_sym_db.RegisterEnumDescriptor(_FUZZINGSTAT)

FuzzingStat = enum_type_wrapper.EnumTypeWrapper(_FUZZINGSTAT)
StatGenerate = 0
StatFuzz = 1
StatCandidate = 2
StatTriage = 3
StatMinimize = 4
StatSmash = 5
StatHint = 6
StatSeed = 7
StatDependency = 8
StatDependencyBoot = 9
StatDefault = -1



_CALL_ADDRESSENTRY = _descriptor.Descriptor(
  name='AddressEntry',
  full_name='dra.Call.AddressEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dra.Call.AddressEntry.key', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='dra.Call.AddressEntry.value', index=1,
      number=2, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
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
  serialized_start=82,
  serialized_end=128,
)

_CALL = _descriptor.Descriptor(
  name='Call',
  full_name='dra.Call',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='idx', full_name='dra.Call.idx', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='address', full_name='dra.Call.address', index=1,
      number=2, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_CALL_ADDRESSENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=20,
  serialized_end=128,
)


_INPUT_CALLENTRY = _descriptor.Descriptor(
  name='CallEntry',
  full_name='dra.Input.CallEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dra.Input.CallEntry.key', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='dra.Input.CallEntry.value', index=1,
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
  serialized_start=439,
  serialized_end=493,
)

_INPUT_UNCOVEREDADDRESSENTRY = _descriptor.Descriptor(
  name='UncoveredAddressEntry',
  full_name='dra.Input.UncoveredAddressEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dra.Input.UncoveredAddressEntry.key', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='dra.Input.UncoveredAddressEntry.value', index=1,
      number=2, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
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
  serialized_start=495,
  serialized_end=550,
)

_INPUT_WRITEADDRESSENTRY = _descriptor.Descriptor(
  name='WriteAddressEntry',
  full_name='dra.Input.WriteAddressEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dra.Input.WriteAddressEntry.key', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='dra.Input.WriteAddressEntry.value', index=1,
      number=2, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
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
  serialized_start=552,
  serialized_end=603,
)

_INPUT = _descriptor.Descriptor(
  name='Input',
  full_name='dra.Input',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='sig', full_name='dra.Input.sig', index=0,
      number=11, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='program', full_name='dra.Input.program', index=1,
      number=12, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=b"",
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='call', full_name='dra.Input.call', index=2,
      number=13, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='paths', full_name='dra.Input.paths', index=3,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='stable', full_name='dra.Input.stable', index=4,
      number=14, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='total', full_name='dra.Input.total', index=5,
      number=15, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='stat', full_name='dra.Input.stat', index=6,
      number=21, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='uncovered_address', full_name='dra.Input.uncovered_address', index=7,
      number=22, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='write_address', full_name='dra.Input.write_address', index=8,
      number=25, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='program_before_mini', full_name='dra.Input.program_before_mini', index=9,
      number=30, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=b"",
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_INPUT_CALLENTRY, _INPUT_UNCOVEREDADDRESSENTRY, _INPUT_WRITEADDRESSENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=131,
  serialized_end=603,
)


_INPUTS = _descriptor.Descriptor(
  name='Inputs',
  full_name='dra.Inputs',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='input', full_name='dra.Inputs.input', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
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
  serialized_start=605,
  serialized_end=640,
)


_PATH = _descriptor.Descriptor(
  name='Path',
  full_name='dra.Path',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='address', full_name='dra.Path.address', index=0,
      number=1, type=13, cpp_type=3, label=3,
      has_default_value=False, default_value=[],
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
  serialized_start=642,
  serialized_end=665,
)


_PATHS_PATHENTRY = _descriptor.Descriptor(
  name='PathEntry',
  full_name='dra.Paths.PathEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dra.Paths.PathEntry.key', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='dra.Paths.PathEntry.value', index=1,
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
  serialized_start=712,
  serialized_end=766,
)

_PATHS = _descriptor.Descriptor(
  name='Paths',
  full_name='dra.Paths',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='path', full_name='dra.Paths.path', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_PATHS_PATHENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=667,
  serialized_end=766,
)


_UNSTABLEINPUT_ADDRESSENTRY = _descriptor.Descriptor(
  name='AddressEntry',
  full_name='dra.UnstableInput.AddressEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dra.UnstableInput.AddressEntry.key', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='dra.UnstableInput.AddressEntry.value', index=1,
      number=2, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
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
  serialized_start=82,
  serialized_end=128,
)

_UNSTABLEINPUT = _descriptor.Descriptor(
  name='UnstableInput',
  full_name='dra.UnstableInput',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='sig', full_name='dra.UnstableInput.sig', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='program', full_name='dra.UnstableInput.program', index=1,
      number=2, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=b"",
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='unstable_path', full_name='dra.UnstableInput.unstable_path', index=2,
      number=12, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='address', full_name='dra.UnstableInput.address', index=3,
      number=13, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_UNSTABLEINPUT_ADDRESSENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=769,
  serialized_end=947,
)


_UNSTABLEINPUTS_UNSTABLEINPUTENTRY = _descriptor.Descriptor(
  name='UnstableInputEntry',
  full_name='dra.UnstableInputs.UnstableInputEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dra.UnstableInputs.UnstableInputEntry.key', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='dra.UnstableInputs.UnstableInputEntry.value', index=1,
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
  serialized_start=1032,
  serialized_end=1104,
)

_UNSTABLEINPUTS = _descriptor.Descriptor(
  name='UnstableInputs',
  full_name='dra.UnstableInputs',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='unstable_input', full_name='dra.UnstableInputs.unstable_input', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_UNSTABLEINPUTS_UNSTABLEINPUTENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=950,
  serialized_end=1104,
)

_CALL_ADDRESSENTRY.containing_type = _CALL
_CALL.fields_by_name['address'].message_type = _CALL_ADDRESSENTRY
_INPUT_CALLENTRY.fields_by_name['value'].message_type = _CALL
_INPUT_CALLENTRY.containing_type = _INPUT
_INPUT_UNCOVEREDADDRESSENTRY.containing_type = _INPUT
_INPUT_WRITEADDRESSENTRY.containing_type = _INPUT
_INPUT.fields_by_name['call'].message_type = _INPUT_CALLENTRY
_INPUT.fields_by_name['paths'].message_type = _PATHS
_INPUT.fields_by_name['stat'].enum_type = _FUZZINGSTAT
_INPUT.fields_by_name['uncovered_address'].message_type = _INPUT_UNCOVEREDADDRESSENTRY
_INPUT.fields_by_name['write_address'].message_type = _INPUT_WRITEADDRESSENTRY
_INPUTS.fields_by_name['input'].message_type = _INPUT
_PATHS_PATHENTRY.fields_by_name['value'].message_type = _PATH
_PATHS_PATHENTRY.containing_type = _PATHS
_PATHS.fields_by_name['path'].message_type = _PATHS_PATHENTRY
_UNSTABLEINPUT_ADDRESSENTRY.containing_type = _UNSTABLEINPUT
_UNSTABLEINPUT.fields_by_name['unstable_path'].message_type = _PATHS
_UNSTABLEINPUT.fields_by_name['address'].message_type = _UNSTABLEINPUT_ADDRESSENTRY
_UNSTABLEINPUTS_UNSTABLEINPUTENTRY.fields_by_name['value'].message_type = _UNSTABLEINPUT
_UNSTABLEINPUTS_UNSTABLEINPUTENTRY.containing_type = _UNSTABLEINPUTS
_UNSTABLEINPUTS.fields_by_name['unstable_input'].message_type = _UNSTABLEINPUTS_UNSTABLEINPUTENTRY
DESCRIPTOR.message_types_by_name['Call'] = _CALL
DESCRIPTOR.message_types_by_name['Input'] = _INPUT
DESCRIPTOR.message_types_by_name['Inputs'] = _INPUTS
DESCRIPTOR.message_types_by_name['Path'] = _PATH
DESCRIPTOR.message_types_by_name['Paths'] = _PATHS
DESCRIPTOR.message_types_by_name['UnstableInput'] = _UNSTABLEINPUT
DESCRIPTOR.message_types_by_name['UnstableInputs'] = _UNSTABLEINPUTS
DESCRIPTOR.enum_types_by_name['FuzzingStat'] = _FUZZINGSTAT
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Call = _reflection.GeneratedProtocolMessageType('Call', (_message.Message,), {

  'AddressEntry' : _reflection.GeneratedProtocolMessageType('AddressEntry', (_message.Message,), {
    'DESCRIPTOR' : _CALL_ADDRESSENTRY,
    '__module__' : 'Input_pb2'
    # @@protoc_insertion_point(class_scope:dra.Call.AddressEntry)
    })
  ,
  'DESCRIPTOR' : _CALL,
  '__module__' : 'Input_pb2'
  # @@protoc_insertion_point(class_scope:dra.Call)
  })
_sym_db.RegisterMessage(Call)
_sym_db.RegisterMessage(Call.AddressEntry)

Input = _reflection.GeneratedProtocolMessageType('Input', (_message.Message,), {

  'CallEntry' : _reflection.GeneratedProtocolMessageType('CallEntry', (_message.Message,), {
    'DESCRIPTOR' : _INPUT_CALLENTRY,
    '__module__' : 'Input_pb2'
    # @@protoc_insertion_point(class_scope:dra.Input.CallEntry)
    })
  ,

  'UncoveredAddressEntry' : _reflection.GeneratedProtocolMessageType('UncoveredAddressEntry', (_message.Message,), {
    'DESCRIPTOR' : _INPUT_UNCOVEREDADDRESSENTRY,
    '__module__' : 'Input_pb2'
    # @@protoc_insertion_point(class_scope:dra.Input.UncoveredAddressEntry)
    })
  ,

  'WriteAddressEntry' : _reflection.GeneratedProtocolMessageType('WriteAddressEntry', (_message.Message,), {
    'DESCRIPTOR' : _INPUT_WRITEADDRESSENTRY,
    '__module__' : 'Input_pb2'
    # @@protoc_insertion_point(class_scope:dra.Input.WriteAddressEntry)
    })
  ,
  'DESCRIPTOR' : _INPUT,
  '__module__' : 'Input_pb2'
  # @@protoc_insertion_point(class_scope:dra.Input)
  })
_sym_db.RegisterMessage(Input)
_sym_db.RegisterMessage(Input.CallEntry)
_sym_db.RegisterMessage(Input.UncoveredAddressEntry)
_sym_db.RegisterMessage(Input.WriteAddressEntry)

Inputs = _reflection.GeneratedProtocolMessageType('Inputs', (_message.Message,), {
  'DESCRIPTOR' : _INPUTS,
  '__module__' : 'Input_pb2'
  # @@protoc_insertion_point(class_scope:dra.Inputs)
  })
_sym_db.RegisterMessage(Inputs)

Path = _reflection.GeneratedProtocolMessageType('Path', (_message.Message,), {
  'DESCRIPTOR' : _PATH,
  '__module__' : 'Input_pb2'
  # @@protoc_insertion_point(class_scope:dra.Path)
  })
_sym_db.RegisterMessage(Path)

Paths = _reflection.GeneratedProtocolMessageType('Paths', (_message.Message,), {

  'PathEntry' : _reflection.GeneratedProtocolMessageType('PathEntry', (_message.Message,), {
    'DESCRIPTOR' : _PATHS_PATHENTRY,
    '__module__' : 'Input_pb2'
    # @@protoc_insertion_point(class_scope:dra.Paths.PathEntry)
    })
  ,
  'DESCRIPTOR' : _PATHS,
  '__module__' : 'Input_pb2'
  # @@protoc_insertion_point(class_scope:dra.Paths)
  })
_sym_db.RegisterMessage(Paths)
_sym_db.RegisterMessage(Paths.PathEntry)

UnstableInput = _reflection.GeneratedProtocolMessageType('UnstableInput', (_message.Message,), {

  'AddressEntry' : _reflection.GeneratedProtocolMessageType('AddressEntry', (_message.Message,), {
    'DESCRIPTOR' : _UNSTABLEINPUT_ADDRESSENTRY,
    '__module__' : 'Input_pb2'
    # @@protoc_insertion_point(class_scope:dra.UnstableInput.AddressEntry)
    })
  ,
  'DESCRIPTOR' : _UNSTABLEINPUT,
  '__module__' : 'Input_pb2'
  # @@protoc_insertion_point(class_scope:dra.UnstableInput)
  })
_sym_db.RegisterMessage(UnstableInput)
_sym_db.RegisterMessage(UnstableInput.AddressEntry)

UnstableInputs = _reflection.GeneratedProtocolMessageType('UnstableInputs', (_message.Message,), {

  'UnstableInputEntry' : _reflection.GeneratedProtocolMessageType('UnstableInputEntry', (_message.Message,), {
    'DESCRIPTOR' : _UNSTABLEINPUTS_UNSTABLEINPUTENTRY,
    '__module__' : 'Input_pb2'
    # @@protoc_insertion_point(class_scope:dra.UnstableInputs.UnstableInputEntry)
    })
  ,
  'DESCRIPTOR' : _UNSTABLEINPUTS,
  '__module__' : 'Input_pb2'
  # @@protoc_insertion_point(class_scope:dra.UnstableInputs)
  })
_sym_db.RegisterMessage(UnstableInputs)
_sym_db.RegisterMessage(UnstableInputs.UnstableInputEntry)


_CALL_ADDRESSENTRY._options = None
_INPUT_CALLENTRY._options = None
_INPUT_UNCOVEREDADDRESSENTRY._options = None
_INPUT_WRITEADDRESSENTRY._options = None
_PATHS_PATHENTRY._options = None
_UNSTABLEINPUT_ADDRESSENTRY._options = None
_UNSTABLEINPUTS_UNSTABLEINPUTENTRY._options = None
# @@protoc_insertion_point(module_scope)
