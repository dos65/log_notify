# lognotify [![Build Status](https://travis-ci.org/dos65/lognotify.svg)](https://travis-ci.org/dos65/lognotify)

Lognotify is a poor tail-like tool for log watching on desktop. The main goal is use sound and desktop popups for notification.

#Usage

```
./my_slow_programm | lognotify -e [expression]
```

```
lognotify -f [PATH_TO_LOG_FILE] -e [expression]
```

# Install

* Dependencies

  Lognotify use ogg123 and notify-send.

  For ubuntu:
  ```
  sudo apt-get install vorbis-tools libnotify-bin
  ```
  
* [Install golang](https://golang.org/doc/install)
* Clone project and build

  ```
  make
  make install
  ```
