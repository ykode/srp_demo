import 'package:angular/angular.dart';
import 'package:ng_bootstrap/ng_bootstrap.dart';
import 'package:angular_forms/angular_forms.dart';

import 'package:convert/convert.dart';
import 'package:crypto/crypto.dart';
import 'dart:convert';


import 'login_service.dart';

List<int> computeHKDF(List<int> ikm, List<int> salt) {
  final hmac1 = new Hmac(sha256, salt);
  final prk = hmac1.convert(ikm);
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
    final salt = base64.decode(session.salt);

    final u = decodeBigInt(hasher.convert(B_bytes).bytes);

    final hmacSha256 = new Hmac(sha256, salt);
    final identity = utf8.encode("$username:$password");
    final xBytes = hmacSha256.convert(identity).bytes;
    final x = decodeBigInt(xBytes);


    final S_c = modPow(B - _loginService.k * modPow(BigInt.two, x, LoginService.N), 
        _loginService.a + u * x, LoginService.N);

    final K_c_bytes = computeHKDF(encodeBigInt(S_c), encodeBigInt(u));
    final K_c = decodeBigInt(K_c_bytes);
    final hasher_M = new Hmac(sha256, K_c_bytes);
    final M1_c_bytes = hasher_M.convert(encodeBigInt(modPow(_loginService.A, B, LoginService.N))).bytes;

    final answer = await _loginService.answerChallenge(session, decodeBigInt(M1_c_bytes));

    print (answer);
  } 
}
