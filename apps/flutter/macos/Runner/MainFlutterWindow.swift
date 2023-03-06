import Cocoa
import AVFoundation
import FlutterMacOS

@objc class MainFlutterWindow: NSWindow {
  override func awakeFromNib() {
    let controller = FlutterViewController.init()
    let windowFrame = self.frame
    self.contentViewController = controller
    self.setFrame(windowFrame, display: true)

    var player: AVPlayer?

    let playerChannel = FlutterMethodChannel(name: "bassiemusic.plaatsoft.nl/player",
                                             binaryMessenger: controller.engine.binaryMessenger)
    playerChannel.setMethodCallHandler({
      (call: FlutterMethodCall, result: @escaping FlutterResult) -> Void in

      if call.method == "start" {
        if let musicUrl = call.arguments as? String {
          print(musicUrl)
          player = AVPlayer(url: URL(string: musicUrl)!)
          player!.play()
          result(true);
        }
        result(false);
      }

      if call.method == "play" {
        player?.play()
        result(true);
      }

      if call.method == "pause" {
        player?.pause()
        result(true);
      }
    })

    RegisterGeneratedPlugins(registry: controller)

    super.awakeFromNib()
  }
}
