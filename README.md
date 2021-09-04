# Radiocorder

Simply record Internet radio or any other streamed audio by url

## Usage

Available Commands:

- `broadcast` Record broadcast
- `help` Help about any command
- `now` Record broadcast from now

You also can use command to get information about all commands

```sh
./radiocorder help [optionally command name]
```

Example:

```sh
./radiocorder broadcast <stream url> "Fri, 23:59" "Sat, 6:00"
```

Example 2:

```sh
./radiocorder now <stream url> "2h13m7s"
```

file with recorded audio has saved to directory where program has launched

## Pull requests

Pull requests are welcome.
If you got any questions, there is [discussions page](https://github.com/fhodun/radiocorder/discussions) enabled

## License

See [LICENSE](LICENSE) file
