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

import 'data_model.dart';
import 'serializers.dart';

const _N_HEX =
    'bde5c829e8d1fe9dd851b3e7c63ba358ddde329b989a0049ab006aadd80aace8e3ffc282d894b5725f2d72d5d98743fcf1a9c02c60b2edbdea7b0328d8d3655ed9b182be6c5b03b5c84b7534404d9a65d6e649df5a28f52aef353ca54a453014fb37ae8f97c1929b012b16ea21a01adddfc4ba05bcc74e8f9a50e422580dfbcb';
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
  static final _N = BigInt.parse(_N_HEX, radix: 16);
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

  LoginService._init(this.k, this.a, this.A, this._infoBits);

  static List<int> randomBytes(int length) =>
      new List.generate(length, (_) => secureRandom.nextInt(254) + 1);

  static BigInt randomBigInt(int bitLength) =>
      decodeBigInt(randomBytes(bitLength >> 3)) % _N;

  factory LoginService() {
    final hmacSha256 = new Hmac(sha256, [2]);
    final kDigest = hmacSha256.convert(hex.decode(_N_HEX));
    final kBytes = kDigest.bytes;
    final a = randomBigInt(_MAX_RANDOM_BIT_LENGTH);
    final A = modPow(g, a, _N);
    if (A % _N == BigInt.zero) {
      throw new Exception("Illegal parameter, A mod N cannot be 0");
    }
    return new LoginService._init(
        decodeBigInt(kBytes), a, A, utf8.encode(keyinfo));
  }

  Future<Null> registerUser(String userName, String password) async {
    final salt = randomBytes(16);
    final hmacSha256 = new Hmac(sha256, salt);
    final xBytes = hmacSha256.convert(salt).bytes;
    final x = decodeBigInt(xBytes);
    final v = modPow(g, x, _N);
    final vBase64 = base64.encode(encodeBigInt(v));
    final saltBase64 = base64.encode(salt);

    final params = {
      "user_name": userName,
      "v": vBase64,
      "salt": saltBase64,
    };

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
}
