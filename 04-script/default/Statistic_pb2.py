# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: Statistic.proto

from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


import Input_pb2 as Input__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='Statistic.proto',
  package='dra',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=b'\n\x0fStatistic.proto\x12\x03\x64ra\x1a\x0bInput.proto\"|\n\tStatistic\x12\x1e\n\x04name\x18\x01 \x01(\x0e\x32\x10.dra.FuzzingStat\x12\x12\n\nexecuteNum\x18\x0b \x01(\x04\x12\x0c\n\x04time\x18\x0c \x01(\x01\x12\x16\n\x0enewTestCaseNum\x18\r \x01(\x04\x12\x15\n\rnewAddressNum\x18\x0e \x01(\x04\"5\n\x04Time\x12\x0c\n\x04time\x18\x01 \x01(\x01\x12\x0b\n\x03num\x18\x02 \x01(\x03\x12\x12\n\nexecuteNum\x18\x03 \x01(\x03\"\x83\x01\n\x08\x43overage\x12-\n\x08\x63overage\x18\x01 \x03(\x0b\x32\x1b.dra.Coverage.CoverageEntry\x12\x17\n\x04time\x18\x02 \x03(\x0b\x32\t.dra.Time\x1a/\n\rCoverageEntry\x12\x0b\n\x03key\x18\x01 \x01(\r\x12\r\n\x05value\x18\x02 \x01(\r:\x02\x38\x01\"X\n\x0bUsefulInput\x12\x19\n\x05input\x18\x01 \x01(\x0b\x32\n.dra.Input\x12\x0c\n\x04time\x18\x02 \x01(\x01\x12\x0b\n\x03num\x18\x03 \x01(\x04\x12\x13\n\x0bnew_address\x18\x04 \x03(\r\"\xea\x01\n\nStatistics\x12\x11\n\tsignalNum\x18\x01 \x01(\x04\x12\x1a\n\x12\x62\x61sic_block_number\x18\n \x01(\r\x12\x1f\n\x08\x63overage\x18\x08 \x01(\x0b\x32\r.dra.Coverage\x12\'\n\x04stat\x18\x0b \x03(\x0b\x32\x19.dra.Statistics.StatEntry\x12&\n\x0cuseful_input\x18\x0c \x03(\x0b\x32\x10.dra.UsefulInput\x1a;\n\tStatEntry\x12\x0b\n\x03key\x18\x01 \x01(\x05\x12\x1d\n\x05value\x18\x02 \x01(\x0b\x32\x0e.dra.Statistic:\x02\x38\x01\x62\x06proto3'
  ,
  dependencies=[Input__pb2.DESCRIPTOR,])




_STATISTIC = _descriptor.Descriptor(
  name='Statistic',
  full_name='dra.Statistic',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='dra.Statistic.name', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='executeNum', full_name='dra.Statistic.executeNum', index=1,
      number=11, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='time', full_name='dra.Statistic.time', index=2,
      number=12, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='newTestCaseNum', full_name='dra.Statistic.newTestCaseNum', index=3,
      number=13, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='newAddressNum', full_name='dra.Statistic.newAddressNum', index=4,
      number=14, type=4, cpp_type=4, label=1,
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
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=37,
  serialized_end=161,
)


_TIME = _descriptor.Descriptor(
  name='Time',
  full_name='dra.Time',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='time', full_name='dra.Time.time', index=0,
      number=1, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='num', full_name='dra.Time.num', index=1,
      number=2, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='executeNum', full_name='dra.Time.executeNum', index=2,
      number=3, type=3, cpp_type=2, label=1,
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
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=163,
  serialized_end=216,
)


_COVERAGE_COVERAGEENTRY = _descriptor.Descriptor(
  name='CoverageEntry',
  full_name='dra.Coverage.CoverageEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dra.Coverage.CoverageEntry.key', index=0,
      number=1, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='dra.Coverage.CoverageEntry.value', index=1,
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
  serialized_start=303,
  serialized_end=350,
)

_COVERAGE = _descriptor.Descriptor(
  name='Coverage',
  full_name='dra.Coverage',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='coverage', full_name='dra.Coverage.coverage', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='time', full_name='dra.Coverage.time', index=1,
      number=2, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_COVERAGE_COVERAGEENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=219,
  serialized_end=350,
)


_USEFULINPUT = _descriptor.Descriptor(
  name='UsefulInput',
  full_name='dra.UsefulInput',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='input', full_name='dra.UsefulInput.input', index=0,
      number=1, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='time', full_name='dra.UsefulInput.time', index=1,
      number=2, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='num', full_name='dra.UsefulInput.num', index=2,
      number=3, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='new_address', full_name='dra.UsefulInput.new_address', index=3,
      number=4, type=13, cpp_type=3, label=3,
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
  serialized_start=352,
  serialized_end=440,
)


_STATISTICS_STATENTRY = _descriptor.Descriptor(
  name='StatEntry',
  full_name='dra.Statistics.StatEntry',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='key', full_name='dra.Statistics.StatEntry.key', index=0,
      number=1, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='value', full_name='dra.Statistics.StatEntry.value', index=1,
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
  serialized_start=618,
  serialized_end=677,
)

_STATISTICS = _descriptor.Descriptor(
  name='Statistics',
  full_name='dra.Statistics',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='signalNum', full_name='dra.Statistics.signalNum', index=0,
      number=1, type=4, cpp_type=4, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='basic_block_number', full_name='dra.Statistics.basic_block_number', index=1,
      number=10, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='coverage', full_name='dra.Statistics.coverage', index=2,
      number=8, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='stat', full_name='dra.Statistics.stat', index=3,
      number=11, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='useful_input', full_name='dra.Statistics.useful_input', index=4,
      number=12, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_STATISTICS_STATENTRY, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=443,
  serialized_end=677,
)

_STATISTIC.fields_by_name['name'].enum_type = Input__pb2._FUZZINGSTAT
_COVERAGE_COVERAGEENTRY.containing_type = _COVERAGE
_COVERAGE.fields_by_name['coverage'].message_type = _COVERAGE_COVERAGEENTRY
_COVERAGE.fields_by_name['time'].message_type = _TIME
_USEFULINPUT.fields_by_name['input'].message_type = Input__pb2._INPUT
_STATISTICS_STATENTRY.fields_by_name['value'].message_type = _STATISTIC
_STATISTICS_STATENTRY.containing_type = _STATISTICS
_STATISTICS.fields_by_name['coverage'].message_type = _COVERAGE
_STATISTICS.fields_by_name['stat'].message_type = _STATISTICS_STATENTRY
_STATISTICS.fields_by_name['useful_input'].message_type = _USEFULINPUT
DESCRIPTOR.message_types_by_name['Statistic'] = _STATISTIC
DESCRIPTOR.message_types_by_name['Time'] = _TIME
DESCRIPTOR.message_types_by_name['Coverage'] = _COVERAGE
DESCRIPTOR.message_types_by_name['UsefulInput'] = _USEFULINPUT
DESCRIPTOR.message_types_by_name['Statistics'] = _STATISTICS
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

Statistic = _reflection.GeneratedProtocolMessageType('Statistic', (_message.Message,), {
  'DESCRIPTOR' : _STATISTIC,
  '__module__' : 'Statistic_pb2'
  # @@protoc_insertion_point(class_scope:dra.Statistic)
  })
_sym_db.RegisterMessage(Statistic)

Time = _reflection.GeneratedProtocolMessageType('Time', (_message.Message,), {
  'DESCRIPTOR' : _TIME,
  '__module__' : 'Statistic_pb2'
  # @@protoc_insertion_point(class_scope:dra.Time)
  })
_sym_db.RegisterMessage(Time)

Coverage = _reflection.GeneratedProtocolMessageType('Coverage', (_message.Message,), {

  'CoverageEntry' : _reflection.GeneratedProtocolMessageType('CoverageEntry', (_message.Message,), {
    'DESCRIPTOR' : _COVERAGE_COVERAGEENTRY,
    '__module__' : 'Statistic_pb2'
    # @@protoc_insertion_point(class_scope:dra.Coverage.CoverageEntry)
    })
  ,
  'DESCRIPTOR' : _COVERAGE,
  '__module__' : 'Statistic_pb2'
  # @@protoc_insertion_point(class_scope:dra.Coverage)
  })
_sym_db.RegisterMessage(Coverage)
_sym_db.RegisterMessage(Coverage.CoverageEntry)

UsefulInput = _reflection.GeneratedProtocolMessageType('UsefulInput', (_message.Message,), {
  'DESCRIPTOR' : _USEFULINPUT,
  '__module__' : 'Statistic_pb2'
  # @@protoc_insertion_point(class_scope:dra.UsefulInput)
  })
_sym_db.RegisterMessage(UsefulInput)

Statistics = _reflection.GeneratedProtocolMessageType('Statistics', (_message.Message,), {

  'StatEntry' : _reflection.GeneratedProtocolMessageType('StatEntry', (_message.Message,), {
    'DESCRIPTOR' : _STATISTICS_STATENTRY,
    '__module__' : 'Statistic_pb2'
    # @@protoc_insertion_point(class_scope:dra.Statistics.StatEntry)
    })
  ,
  'DESCRIPTOR' : _STATISTICS,
  '__module__' : 'Statistic_pb2'
  # @@protoc_insertion_point(class_scope:dra.Statistics)
  })
_sym_db.RegisterMessage(Statistics)
_sym_db.RegisterMessage(Statistics.StatEntry)


_COVERAGE_COVERAGEENTRY._options = None
_STATISTICS_STATENTRY._options = None
# @@protoc_insertion_point(module_scope)
