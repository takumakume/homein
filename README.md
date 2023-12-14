# homein

Detecting homograph domains using Levenshtein and Hamming distance.

## How to use

### CLI

```sh
$ homein "M1CR0Z0FT.COM" "MICROSOFT.COM"
levenshtein distance: 4, levenshtein percent: 69.23%
image hash distance: 4
```

```sh
$ homein "M1CR0Z0FT.COM" "MICROSOFT.COM" --enable-output-images
levenshtein distance: 4, levenshtein percent: 69.23%
save image: M1CR0Z0FT.COM.png
save image: MICROSOFT.COM.png
image hash distance: 4

$ ls *.png
M1CR0Z0FT.COM.png
MICROSOFT.COM.png
```

### Docker

```sh
$ docker run -it takumakume/homein M1CR0Z0FT.COM MICROSOFT.COM
levenshtein distance: 4, levenshtein percent: 69.23%
image hash distance: 4
```
