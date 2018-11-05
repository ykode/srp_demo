import 'dart:convert';
import 'dart:typed_data';
import 'dart:math';
import 'dart:async';
import 'dart:io';
import 'package:convert/convert.dart';
import 'package:crypto/crypto.dart';
import 'package:http/http.dart' as http;
import 'package:built_value/standard_json_plugin.dart';
import 'package:built_value/serializer.dart';
import 'package:logging/logging.dart';

import 'data_model.dart';
import 'serializers.dart';

final Uint8List _N_BYTES = Uint8List.fromList([
    0xBD, 0xE5, 0xC8, 0x29, 0xE8, 0xD1, 0xFE, 0x9D, 0xD8, 0x51, 0xB3, 0xE7,
		0xC6, 0x3B, 0xA3, 0x58, 0xDD, 0xDE, 0x32, 0x9B, 0x98, 0x9A, 0x00, 0x49,
		0xAB, 0x00, 0x6A, 0xAD, 0xD8, 0x0A, 0xAC, 0xE8, 0xE3, 0xFF, 0xC2, 0x82,
		0xD8, 0x94, 0xB5, 0x72, 0x5F, 0x2D, 0x72, 0xD5, 0xD9, 0x87, 0x43, 0xFC,
		0xF1, 0xA9, 0xC0, 0x2C, 0x60, 0xB2, 0xED, 0xBD, 0xEA, 0x7B, 0x03, 0x28,
		0xD8, 0xD3, 0x65, 0x5E, 0xD9, 0xB1, 0x82, 0xBE, 0x6C, 0x5B, 0x03, 0xB5,
		0xC8, 0x4B, 0x75, 0x34, 0x40, 0x4D, 0x9A, 0x65, 0xD6, 0xE6, 0x49, 0xDF,
		0x5A, 0x28, 0xF5, 0x2A, 0xEF, 0x35, 0x3C, 0xA5, 0x4A, 0x45, 0x30, 0x14,
		0xFB, 0x37, 0xAE, 0x8F, 0x97, 0xC1, 0x92, 0x9B, 0x01, 0x2B, 0x16, 0xEA,
		0x21, 0xA0, 0x1A, 0xDD, 0xDF, 0xC4, 0xBA, 0x05, 0xBC, 0xC7, 0x4E, 0x8F,
		0x9A, 0x50, 0xE4, 0x22, 0x58, 0x0D, 0xFB, 0xCB]); 
const keyinfo = "SRP Demo Key Information";

final _bigIntFF = new BigInt.from(0xff);

BigInt decodeBigInt(List<int> bytes) {
  var result = BigInt.zero;
  for (int i = 0; i < bytes.length; i++) {
    result += new BigInt.from(bytes[bytes.length - i - 1]) << (8 * i);
  }
  return result;
}

var _byteMask = new BigInt.from(0xff);

/// Encode a BigInt into bytes using big-endian encoding.
Uint8List encodeBigInt(BigInt number) {

  if (number < BigInt.zero) {
    throw new RangeError("Cannot encode negative BigInt number");
  }

  final size = (number.bitLength + 7) >> 3;
  var result = new Uint8List(size);
  for (int i = 0; i < size; i++) {
    result[size - i - 1] = (number & _byteMask).toInt();
    number = number >> 8;
  }
  return result;
}

BigInt modPow(BigInt b, BigInt e, BigInt m) => e < BigInt.one
    ? BigInt.one
    : (b < BigInt.zero || b > m ? (b % m) : b).modPow(e, m);

class LoginService {
  static final N = decodeBigInt(_N_BYTES);
  static final g = BigInt.two;
  static const _MAX_RANDOM_BIT_LENGTH = 1024;
  static final secureRandom = new Random.secure();
  final _serializers =
      (serializers.toBuilder()..addPlugin(new StandardJsonPlugin())).build();

  final BigInt k;
  final BigInt a;
  final BigInt A;
  final Uint8List _infoBits;

  final _client = http.Client();

  BigInt B;

  final Logger _log = new Logger('LoginService');

  LoginService._init(this.k, this.a, this.A, this._infoBits);

  static List<int> randomBytes(int length) =>
      new List.generate(length, (_) => secureRandom.nextInt(254) + 1);

  static BigInt randomBigInt(int bitLength) =>
      decodeBigInt(randomBytes(bitLength >> 3)) % N;

  factory LoginService() {
    final hmacSha256 = new Hmac(sha256, [2]);
    final kDigest = hmacSha256.convert(_N_BYTES);
    final kBytes = kDigest.bytes;
    final a = randomBigInt(_MAX_RANDOM_BIT_LENGTH);
    final A = modPow(g, a, N);
    if (A % N == BigInt.zero) {
      throw new Exception("Illegal parameter, A mod N cannot be 0");
    }
    return new LoginService._init(
        decodeBigInt(kBytes), a, A, utf8.encode(keyinfo));
  }

  Future<Null> registerUser(String userName, String password) async {
    final salt = randomBytes(16);
    final hmacSha256 = new Hmac(sha256, salt);
    final identity = utf8.encode("$userName:$password");
    final xBytes = hmacSha256.convert(identity).bytes;
    final x = decodeBigInt(xBytes);
    final v = modPow(g, x, N);
    final vBase64 = base64.encode(encodeBigInt(v));
    final saltBase64 = base64.encode(salt);

    final params = {
      "user_name": userName,
      "v": vBase64,
      "salt": saltBase64,
    };

    print(params);
    print(base64.encode(xBytes));

    final uri = Uri.parse("http://localhost:4000/identities/");

    final request = new http.Request("POST", uri);

    request.headers['content-type'] = 'application/x-www-form-urlencoded';
    request.bodyFields = params;

    final response = await _client.send(request);
    if (201 != response.statusCode) {
      throw new HttpException("Cannot create user", uri: uri);
    } 
  }

  Future<Session> startSession(String userName) async {
    final A_base64 = base64.encode(encodeBigInt(A));

    final params = {
      "action": "start_session",
      "user_name": userName,
      "A": A_base64,
    };

    print(params);

    final uri = Uri.parse("http://localhost:4000/sessions/");

    final request = new http.Request("POST", uri);

    request.headers['content-type'] = 'application/x-www-form-urlencoded';
    request.bodyFields = params;

    final response = await _client.send(request);
    final decodeMap = await response.stream
        .transform(utf8.decoder)
        .transform(json.decoder)
        .cast<Map<String, dynamic>>()
        .single;
    if (201 != response.statusCode) {
      throw new HttpException("Cannot start session", uri: uri);
    } else {
      return _serializers.deserialize(decodeMap,
          specifiedType: const FullType(Session));
    }
  }

  Future<ChallengeAnswer> answerChallenge(Session s, BigInt mClient) async {
    final String M_c = base64.encode(encodeBigInt(mClient));

    final params = {
      "action":"answer",
      "session_id":s.sessionId,
      "m_client": M_c,
    };

    print(params);

    final uri = Uri.parse("http://localhost:4000/sessions/");

    final request = new http.Request("POST", uri);

    request.headers['content-type'] = 'application/x-www-form-urlencoded';
    request.bodyFields = params;

    final response = await _client.send(request);
    final decodeMap = await response.stream
        .transform(utf8.decoder)
        .transform(json.decoder)
        .cast<Map<String, dynamic>>()
        .single;

    if (200 != response.statusCode) {
      throw new HttpException("Unauthorized", uri:uri);
    } else {
      return _serializers.deserialize(decodeMap,
          specifiedType: const FullType(ChallengeAnswer));
    }
  }
}
