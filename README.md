# Radiocorder

Repeatedly record Internet radio on specified time and date (audition time)

## Usage

For now broadcast date and time settings are constantly assigned in source code

```sh
./radiocorder <stream url>
```

Example:

```sh
$ ./radiocorder <stream url>
INFO[0000] Found next broadcast date                     broadcastDuration=0s date="2021-08-04 00:52:59.19727516 +0200 CEST"
INFO[0000] Recording started                             currentTime="2021-08-04 00:52:59.670034944 +0200 CEST"
INFO[0003] Recorded audio saved to file                  fileName=stream.ogg
```

file with recorded audio has saved to directory where program has launched

## License

See [LICENSE](LICENSE) file
