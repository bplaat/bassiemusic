# BassieMusic Flutter App

## Platform channels
This Flutter app uses some custom native platform channels:

### Music Player
Name: `bassiemusic.plaatsoft.nl/player`

#### Flutter -> Native:
```dart
bool start(String url);
bool stop();
bool play();
bool pause();
bool seek(double position);
```

#### Native -> Flutter:
```dart
void loaded(double durtion);
void positionChange(double position);
void ended();
```
