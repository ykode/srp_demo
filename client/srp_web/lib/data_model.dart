import 'package:built_value/built_value.dart';
import 'package:built_value/serializer.dart';

part 'data_model.g.dart';

abstract class Session implements Built<Session, SessionBuilder> {
  static Serializer<Session> get serializer => _$sessionSerializer;

  @BuiltValueField(wireName: 'B')
  String get B;

  @BuiltValueField(wireName: 'session_id')
  String get sessionId;

  @BuiltValueField(wireName: 'salt')
  String get salt;

  factory Session([updates(SessionBuilder b)]) = _$Session;

  Session._();
}

abstract class ChallengeAnswer implements Built<ChallengeAnswer, ChallengeAnswerBuilder> {
  static Serializer<ChallengeAnswer> get serializer => _$challengeAnswerSerializer;

  @BuiltValueField(wireName: 'sessionId')
  String get sessionId;

  @BuiltValueField(wireName: 'M_s')
  String get M2;

  factory ChallengeAnswer([updates(ChallengeAnswerBuilder b)]) = _$ChallengeAnswer;

  ChallengeAnswer._();
}
