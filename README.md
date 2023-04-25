# GoinGo - A Strategic Board Game

GoinGo is an implementation of the Go board game developed in the Go programming language using [ebiten](https://github.com/hajimehoshi/ebiten), [donburi](https://github.com/yohamta/donburi) and [furex](https://github.com/yohamta/furex). The game allows you to play in local multiplayer ot play in single-player using [GNUGo](https://www.gnu.org/software/gnugo/).

## Getting Started

### Prerequisites

If you wanna play single player, make sure that GNUGo installed on your computer.

### Installation

You can download the latest release of GoinGo from the [Releases page](https://github.com/joelschutz/GoinGo/releases). There are pre-built binaries available for Windows and Linux. Alternatively, you can clone the repository and build the game yourself:

```
git clone https://github.com/joelschutz/GoinGo.git
cd GoinGo
go build
```

### Running the game

To start a game, run the following command:

```
// Linux
./goingo-amd64-linux
// Windows
./goingo-amd64.exe
```

## Roadmap
 - Finish main Menu
 - Improve board ui
   - Show Last move
   - Show points better
   - Add helpers for liberties, hints and Atari
 - Add endgame behaviour
 - Allow human player to start as white
 - Add Handcap
 - Clean the code

## Contributing

GoinGo is open to contributions from the community. Feel free to submit a pull request with your changes or improvements.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.

## Acknowledgments

 - Sounds by [MPOOMAN](https://www.youtube.com/channel/UCNEIrnARXCFNpVWZlwg2pzg) 
 - Learn how to play with [In Sente](https://www.youtube.com/watch?v=NJ9QIiWgLWw)
 - This README.md was created using ChatGPT.
