mkdir -p ~/.local/bin
gcc -s -Os bassiemusic.c $(pkg-config --cflags --libs gtk+-3.0 webkit2gtk-4.0) \
    -o ~/.local/bin/nl.plaatsoft.BassieMusic || exit

mkdir -p ~/.local/share/applications
cp nl.plaatsoft.BassieMusic.desktop ~/.local/share/applications
mkdir -p ~/.local/share/icons/hicolor/scalable/apps
cp nl.plaatsoft.BassieMusic.svg ~/.local/share/icons/hicolor/scalable/apps

~/.local/bin/nl.plaatsoft.BassieMusic

# sudo gtk-update-icon-cache -f /usr/share/icons/hicolor/
