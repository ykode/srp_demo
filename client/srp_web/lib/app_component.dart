import 'package:angular/angular.dart';
import 'package:ng_bootstrap/ng_bootstrap.dart';
import 'package:angular_forms/angular_forms.dart';


import 'login_service.dart';


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
 } 
}
