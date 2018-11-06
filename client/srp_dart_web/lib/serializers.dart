import 'package:built_value/serializer.dart';
import 'data_model.dart';

part 'serializers.g.dart';

@SerializersFor(const [
  Session,
  ChallengeAnswer
])
final Serializers serializers = _$serializers.toBuilder()
    .build();
