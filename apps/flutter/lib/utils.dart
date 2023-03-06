import 'dart:io' show Platform;

String userAgent() {
  if (Platform.isAndroid) {
    return 'BassieMusic Android App v0.1.0';
  }
  if (Platform.isIOS) {
    return 'BassieMusic iOS App v0.1.0';
  }
  if (Platform.isMacOS) {
    return 'BassieMusic macOS App v0.1.0';
  }
  if (Platform.isWindows) {
    return 'BassieMusic Windows App v0.1.0';
  }
  if (Platform.isLinux) {
    return 'BassieMusic Linux App v0.1.0';
  }
  return 'BassieMusic Flutter App v0.1.0';
}
