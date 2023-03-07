import Cocoa
import AVFoundation
import FlutterMacOS

class MainFlutterWindow: NSWindow {
    var player: AVPlayer?

    func playerMethodCallHandler(call: FlutterMethodCall, result: @escaping FlutterResult) -> Void {
        if call.method == "start" {
            if let musicUrl = call.arguments as? String {
                player = AVPlayer(url: URL(string: musicUrl)!)
                // player.item.duration
                // https://stackoverflow.com/questions/29386531/how-to-detect-when-avplayer-video-ends-playing
                player!.play()
                result(true)
            }
            result(false)
        }

        if call.method == "stop" {
            if player != nil {
                player = nil
                result(true)
            } else {
                result(false)
            }
        }

        if call.method == "play" {
            if player != nil {
                player!.play()
                result(true)
            } else {
                result(false)
            }
        }

        if call.method == "pause" {
            if player != nil {
                player!.pause()
                result(true)
            } else {
                result(false)
            }
        }

        if call.method == "seek" {
            if let position = call.arguments as? Double {
                if player != nil {
                    player!.seek(to: CMTimeMake(value: Int64(position * 1000), timescale: 1000))
                    result(true)
                } else {
                    result(false)
                }
            } else {
                result(false)
            }
        }
    }

    override func awakeFromNib() {
        let controller = FlutterViewController.init()
        let windowFrame = self.frame
        self.contentViewController = controller
        self.setFrame(windowFrame, display: true)

        let playerChannel = FlutterMethodChannel(name: "bassiemusic.plaatsoft.nl/player",
                                                binaryMessenger: controller.engine.binaryMessenger)
        playerChannel.setMethodCallHandler(self.playerMethodCallHandler)

        RegisterGeneratedPlugins(registry: controller)

        super.awakeFromNib()
    }
}
