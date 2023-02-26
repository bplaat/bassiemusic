# Simple build script to build the BassieMusic Windows Client app
# with MinGW (MSYS2 & pacman -S mingw-w64-x86_64-toolchain),
# 7-zip and a minify-xml tool (npm install -g minify-xml)
# to build release use ./build.sh release

if [ "$1" = "clean" ]; then
    rm -rf build
    exit
fi

rm -rf build
mkdir build
minify-xml res/app.manifest > build/app.min.manifest || exit
windres res/resource.rc -o build/resource.o || exit

if [ "$1" = "release" ]; then
    mkdir build/bassiemusic-win64
    gcc -c -Os -IWebView2 bassiemusic.c -o build/bassiemusic.o || exit
    ld -s --subsystem windows build/bassiemusic.o build/resource.o -e _start \
        -L"C:\\Windows\\System32" -lkernel32 -luser32 -lgdi32 -lshell32 -lole32 -lversion -ldwmapi -ld3d11 -o build/bassiemusic-win64/bassiemusic.exe
    cp WebView2/WebView2Loader.dll build/bassiemusic-win64

    cd build
    7z a bassiemusic-win64.zip bassiemusic-win64 > /dev/null
    exit
fi

gcc -c -IWebView2 src/main.c -o build/main.o || exit
gcc -c src/about.c -o build/about.o || exit
gcc -c src/utils.c -o build/utils.o || exit
ld build/main.o build/about.o build/utils.o build/resource.o -e _start \
    -L"C:\\Windows\\System32" -lkernel32 -luser32 -lgdi32 -lshell32 -lole32 -lversion -ldwmapi -ld3d11 -o build/bassiemusic.exe
cp WebView2/WebView2Loader.dll build
./build/bassiemusic
