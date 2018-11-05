import 'package:angular/angular.dart';
import 'package:ng_bootstrap/ng_bootstrap.dart';
@Component(
  selector: 'main',
  templateUrl: 'app_component.html',
  styleUrls: ['app_component.css'],
  directives: const [bsDirectives],
)
class AppComponent {
  var name = 'Angular';
}
