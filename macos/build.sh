find . -name ".DS_Store" -delete
mkdir -p BassieMusic.app/Contents/MacOS BassieMusic.app/Contents/Resources
if [[ $1 = "release" ]]; then
    clang -x objective-c --target=arm64-macos -Os bassiemusic.m -framework Cocoa -framework WebKit -o BassieMusic-arm64 || exit 1
    clang -x objective-c --target=x86_64-macos -Os bassiemusic.m -framework Cocoa -framework WebKit -o BassieMusic-x86_64 || exit 1
    strip BassieMusic-arm64 BassieMusic-x86_64
    lipo BassieMusic-arm64 BassieMusic-x86_64 -create -output BassieMusic.app/Contents/MacOS/BassieMusic
    rm BassieMusic-arm64 BassieMusic-x86_64
else
    clang -x objective-c bassiemusic.m -framework Cocoa -framework WebKit -o BassieMusic.app/Contents/MacOS/BassieMusic || exit 1
fi
cp -r Resources BassieMusic.app/Contents
cp Info.plist BassieMusic.app/Contents
if [[ $1 = "release" ]]; then
    zip -r BassieMusic.app.zip BassieMusic.app
fi
open BassieMusic.app
