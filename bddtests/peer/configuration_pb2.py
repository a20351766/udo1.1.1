# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: peer/configuration.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf import descriptor_pb2
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='peer/configuration.proto',
  package='protos',
  syntax='proto3',
  serialized_pb=_b('\n\x18peer/configuration.proto\x12\x06protos\"7\n\x0b\x41nchorPeers\x12(\n\x0c\x61nchor_peers\x18\x01 \x03(\x0b\x32\x12.protos.AnchorPeer\"(\n\nAnchorPeer\x12\x0c\n\x04host\x18\x01 \x01(\t\x12\x0c\n\x04port\x18\x02 \x01(\x05\x42O\n\"org.hyperledger.udo.protos.peerZ)github.com/hyperledger/udo/protos/peerb\x06proto3')
)
_sym_db.RegisterFileDescriptor(DESCRIPTOR)




_ANCHORPEERS = _descriptor.Descriptor(
  name='AnchorPeers',
  full_name='protos.AnchorPeers',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='anchor_peers', full_name='protos.AnchorPeers.anchor_peers', index=0,
      number=1, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=36,
  serialized_end=91,
)


_ANCHORPEER = _descriptor.Descriptor(
  name='AnchorPeer',
  full_name='protos.AnchorPeer',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='host', full_name='protos.AnchorPeer.host', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
    _descriptor.FieldDescriptor(
      name='port', full_name='protos.AnchorPeer.port', index=1,
      number=2, type=5, cpp_type=1, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      options=None),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=93,
  serialized_end=133,
)

_ANCHORPEERS.fields_by_name['anchor_peers'].message_type = _ANCHORPEER
DESCRIPTOR.message_types_by_name['AnchorPeers'] = _ANCHORPEERS
DESCRIPTOR.message_types_by_name['AnchorPeer'] = _ANCHORPEER

AnchorPeers = _reflection.GeneratedProtocolMessageType('AnchorPeers', (_message.Message,), dict(
  DESCRIPTOR = _ANCHORPEERS,
  __module__ = 'peer.configuration_pb2'
  # @@protoc_insertion_point(class_scope:protos.AnchorPeers)
  ))
_sym_db.RegisterMessage(AnchorPeers)

AnchorPeer = _reflection.GeneratedProtocolMessageType('AnchorPeer', (_message.Message,), dict(
  DESCRIPTOR = _ANCHORPEER,
  __module__ = 'peer.configuration_pb2'
  # @@protoc_insertion_point(class_scope:protos.AnchorPeer)
  ))
_sym_db.RegisterMessage(AnchorPeer)


DESCRIPTOR.has_options = True
DESCRIPTOR._options = _descriptor._ParseOptions(descriptor_pb2.FileOptions(), _b('\n\"org.hyperledger.udo.protos.peerZ)github.com/hyperledger/udo/protos/peer'))
try:
  # THESE ELEMENTS WILL BE DEPRECATED.
  # Please use the generated *_pb2_grpc.py files instead.
  import grpc
  from grpc.framework.common import cardinality
  from grpc.framework.interfaces.face import utilities as face_utilities
  from grpc.beta import implementations as beta_implementations
  from grpc.beta import interfaces as beta_interfaces
except ImportError:
  pass
# @@protoc_insertion_point(module_scope)
