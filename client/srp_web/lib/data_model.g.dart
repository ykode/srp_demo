// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'data_model.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

// ignore_for_file: always_put_control_body_on_new_line
// ignore_for_file: annotate_overrides
// ignore_for_file: avoid_annotating_with_dynamic
// ignore_for_file: avoid_catches_without_on_clauses
// ignore_for_file: avoid_returning_this
// ignore_for_file: lines_longer_than_80_chars
// ignore_for_file: omit_local_variable_types
// ignore_for_file: prefer_expression_function_bodies
// ignore_for_file: sort_constructors_first
// ignore_for_file: unnecessary_const
// ignore_for_file: unnecessary_new
// ignore_for_file: test_types_in_equals

Serializer<Session> _$sessionSerializer = new _$SessionSerializer();
Serializer<ChallengeAnswer> _$challengeAnswerSerializer =
    new _$ChallengeAnswerSerializer();

class _$SessionSerializer implements StructuredSerializer<Session> {
  @override
  final Iterable<Type> types = const [Session, _$Session];
  @override
  final String wireName = 'Session';

  @override
  Iterable serialize(Serializers serializers, Session object,
      {FullType specifiedType = FullType.unspecified}) {
    final result = <Object>[
      'B',
      serializers.serialize(object.B, specifiedType: const FullType(String)),
      'session_id',
      serializers.serialize(object.sessionId,
          specifiedType: const FullType(String)),
      'salt',
      serializers.serialize(object.salt, specifiedType: const FullType(String)),
    ];

    return result;
  }

  @override
  Session deserialize(Serializers serializers, Iterable serialized,
      {FullType specifiedType = FullType.unspecified}) {
    final result = new SessionBuilder();

    final iterator = serialized.iterator;
    while (iterator.moveNext()) {
      final key = iterator.current as String;
      iterator.moveNext();
      final dynamic value = iterator.current;
      switch (key) {
        case 'B':
          result.B = serializers.deserialize(value,
              specifiedType: const FullType(String)) as String;
          break;
        case 'session_id':
          result.sessionId = serializers.deserialize(value,
              specifiedType: const FullType(String)) as String;
          break;
        case 'salt':
          result.salt = serializers.deserialize(value,
              specifiedType: const FullType(String)) as String;
          break;
      }
    }

    return result.build();
  }
}

class _$ChallengeAnswerSerializer
    implements StructuredSerializer<ChallengeAnswer> {
  @override
  final Iterable<Type> types = const [ChallengeAnswer, _$ChallengeAnswer];
  @override
  final String wireName = 'ChallengeAnswer';

  @override
  Iterable serialize(Serializers serializers, ChallengeAnswer object,
      {FullType specifiedType = FullType.unspecified}) {
    final result = <Object>[
      'sessionId',
      serializers.serialize(object.sessionId,
          specifiedType: const FullType(String)),
      'M_s',
      serializers.serialize(object.M2, specifiedType: const FullType(String)),
    ];

    return result;
  }

  @override
  ChallengeAnswer deserialize(Serializers serializers, Iterable serialized,
      {FullType specifiedType = FullType.unspecified}) {
    final result = new ChallengeAnswerBuilder();

    final iterator = serialized.iterator;
    while (iterator.moveNext()) {
      final key = iterator.current as String;
      iterator.moveNext();
      final dynamic value = iterator.current;
      switch (key) {
        case 'sessionId':
          result.sessionId = serializers.deserialize(value,
              specifiedType: const FullType(String)) as String;
          break;
        case 'M_s':
          result.M2 = serializers.deserialize(value,
              specifiedType: const FullType(String)) as String;
          break;
      }
    }

    return result.build();
  }
}

class _$Session extends Session {
  @override
  final String B;
  @override
  final String sessionId;
  @override
  final String salt;

  factory _$Session([void updates(SessionBuilder b)]) =>
      (new SessionBuilder()..update(updates)).build();

  _$Session._({this.B, this.sessionId, this.salt}) : super._() {
    if (B == null) {
      throw new BuiltValueNullFieldError('Session', 'B');
    }
    if (sessionId == null) {
      throw new BuiltValueNullFieldError('Session', 'sessionId');
    }
    if (salt == null) {
      throw new BuiltValueNullFieldError('Session', 'salt');
    }
  }

  @override
  Session rebuild(void updates(SessionBuilder b)) =>
      (toBuilder()..update(updates)).build();

  @override
  SessionBuilder toBuilder() => new SessionBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is Session &&
        B == other.B &&
        sessionId == other.sessionId &&
        salt == other.salt;
  }

  @override
  int get hashCode {
    return $jf($jc($jc($jc(0, B.hashCode), sessionId.hashCode), salt.hashCode));
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper('Session')
          ..add('B', B)
          ..add('sessionId', sessionId)
          ..add('salt', salt))
        .toString();
  }
}

class SessionBuilder implements Builder<Session, SessionBuilder> {
  _$Session _$v;

  String _B;
  String get B => _$this._B;
  set B(String B) => _$this._B = B;

  String _sessionId;
  String get sessionId => _$this._sessionId;
  set sessionId(String sessionId) => _$this._sessionId = sessionId;

  String _salt;
  String get salt => _$this._salt;
  set salt(String salt) => _$this._salt = salt;

  SessionBuilder();

  SessionBuilder get _$this {
    if (_$v != null) {
      _B = _$v.B;
      _sessionId = _$v.sessionId;
      _salt = _$v.salt;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(Session other) {
    if (other == null) {
      throw new ArgumentError.notNull('other');
    }
    _$v = other as _$Session;
  }

  @override
  void update(void updates(SessionBuilder b)) {
    if (updates != null) updates(this);
  }

  @override
  _$Session build() {
    final _$result =
        _$v ?? new _$Session._(B: B, sessionId: sessionId, salt: salt);
    replace(_$result);
    return _$result;
  }
}

class _$ChallengeAnswer extends ChallengeAnswer {
  @override
  final String sessionId;
  @override
  final String M2;

  factory _$ChallengeAnswer([void updates(ChallengeAnswerBuilder b)]) =>
      (new ChallengeAnswerBuilder()..update(updates)).build();

  _$ChallengeAnswer._({this.sessionId, this.M2}) : super._() {
    if (sessionId == null) {
      throw new BuiltValueNullFieldError('ChallengeAnswer', 'sessionId');
    }
    if (M2 == null) {
      throw new BuiltValueNullFieldError('ChallengeAnswer', 'M2');
    }
  }

  @override
  ChallengeAnswer rebuild(void updates(ChallengeAnswerBuilder b)) =>
      (toBuilder()..update(updates)).build();

  @override
  ChallengeAnswerBuilder toBuilder() =>
      new ChallengeAnswerBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ChallengeAnswer &&
        sessionId == other.sessionId &&
        M2 == other.M2;
  }

  @override
  int get hashCode {
    return $jf($jc($jc(0, sessionId.hashCode), M2.hashCode));
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper('ChallengeAnswer')
          ..add('sessionId', sessionId)
          ..add('M2', M2))
        .toString();
  }
}

class ChallengeAnswerBuilder
    implements Builder<ChallengeAnswer, ChallengeAnswerBuilder> {
  _$ChallengeAnswer _$v;

  String _sessionId;
  String get sessionId => _$this._sessionId;
  set sessionId(String sessionId) => _$this._sessionId = sessionId;

  String _M2;
  String get M2 => _$this._M2;
  set M2(String M2) => _$this._M2 = M2;

  ChallengeAnswerBuilder();

  ChallengeAnswerBuilder get _$this {
    if (_$v != null) {
      _sessionId = _$v.sessionId;
      _M2 = _$v.M2;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ChallengeAnswer other) {
    if (other == null) {
      throw new ArgumentError.notNull('other');
    }
    _$v = other as _$ChallengeAnswer;
  }

  @override
  void update(void updates(ChallengeAnswerBuilder b)) {
    if (updates != null) updates(this);
  }

  @override
  _$ChallengeAnswer build() {
    final _$result =
        _$v ?? new _$ChallengeAnswer._(sessionId: sessionId, M2: M2);
    replace(_$result);
    return _$result;
  }
}
