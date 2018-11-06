import 'package:angular/angular.dart';
import 'package:ng_bootstrap/ng_bootstrap.dart';
import 'package:angular_forms/angular_forms.dart';

import 'package:convert/convert.dart';
import 'package:crypto/crypto.dart';
import 'dart:convert';
import 'dart:html';
import 'dart:io';


import 'login_service.dart';

List<int> computeHKDF(List<int> ikm, List<int> salt) {
  final hmac1 = new Hmac(sha256, salt);
  final prk = hmac1.convert(ikm);

  print(base64.encode(prk.bytes));
  final hmac2 = new Hmac(sha256, prk.bytes);
  final dig = hmac2.convert(utf8.encode(keyinfo) + [1]);
  return dig.bytes.sublist(0, 16);
}


@Component(
  selector: 'main',
  templateUrl: 'app_component.html',
  styleUrls: ['app_component.css'],
  providers: [ClassProvider(LoginService)],
  directives: [formDirectives, bsDirectives],
)
class AppComponent {
  var username;
  var password;

  final LoginService _loginService;

  AppComponent(this._loginService);

  void onRegisterClicked() {
   _loginService.registerUser(username, password);
  }

  void onLoginClicked() async {
    final session = await _loginService.startSession(username);
    print(session);

    final hasher = new Hmac(sha256, encodeBigInt(_loginService.A));

    final B_bytes = base64.decode(session.B);
    final B = decodeBigInt(B_bytes);
    final A = _loginService.A;
    final N = LoginService.N;
    final salt = base64.decode(session.salt);

    final u = decodeBigInt(hasher.convert(B_bytes).bytes);

    final hmacSha256 = new Hmac(sha256, salt);
    final identity = utf8.encode("$username:$password");
    final xBytes = hmacSha256.convert(identity).bytes;
    final x = decodeBigInt(xBytes);


    final S_c = modPow(B - _loginService.k * modPow(BigInt.two, x, N), 
        _loginService.a + u * x, LoginService.N);

    print("S_c: ${S_c.toRadixString(16)}");

    final K_c_bytes = computeHKDF(encodeBigInt(S_c), encodeBigInt(u));
    final K_c = decodeBigInt(K_c_bytes);

    print("K_c ${K_c.toRadixString(16)}");
    final hasher_M = new Hmac(sha256, K_c_bytes);
    final M1_c_bytes = hasher_M.convert(encodeBigInt(modPow(A, B, N))).bytes;
    final M1_c = decodeBigInt(M1_c_bytes);

    print("A: ${A.toRadixString(16)}\n B: ${B.toRadixString(16)}\n");

    try {
      final answer = await _loginService.answerChallenge(session, decodeBigInt(M1_c_bytes));
  
      print (answer);

      final M2_s_bytes = base64.decode(answer.M2);
      final M2_s = decodeBigInt(M2_s_bytes);
    
      final M2_c_bytes = hasher_M.convert(encodeBigInt(modPow(A, M1_c, N))).bytes;
      final M2_c = decodeBigInt(M2_c_bytes);

      if (M2_s == M2_c) {
        window.alert("Authorised!");
      } else {
        window.alert("Unauthorised!");
      }
    } catch (e) {
      if (e is HttpException) {
        window.alert(e.message);
      } else {
        throw e;
      }
    }
  } 
}
