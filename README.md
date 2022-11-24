# Maze Game
Simple Maze game coded in Go with genetic algorithms option to search for the best solutions.

## Objective:
Cross the screen and reach the empty space on last column at the right.

**Human** | **Genetic Algorithms**
:-------------------------:|:-------------------------:
<img width="430" alt="horizontal" src="https://github.com/cassianoperin/Maze_Game/blob/main/Images/maze-human.png">  |  <img width="430" alt="vertical" src="https://github.com/cassianoperin/Maze_Game/blob/main/Images/maze-automations.gif">

## Usage
1)  After the first execution, the program will create an ini file named '.maze.ini' into user home folder
  - To execute the game, set the value 'Automation' to false, otherwise, it will start in simulation mode
  - Select the map from 0 to 3
2) Define the genetic altorithm configuration:
  - Number of generations (Generations)
  - Population size (Population_size)
  - Number of genes (Gene_number)
  - Number of participants of tournament for parents selection (K)
  - Crossover rate (Crossover_rate)
  - Mutation rate (Mutation_rate)
  - Elitism percentual (Elitism_percentual)
3) Run the program

## Next steps:
- Improve score considering the individual that got the best result in less movements.
- After finish, show the path of winner
- Key to reset
- Binary for Windows
- Clean code
- Show the time spent after the execution
- Translate the individual into arrows
- Reactivate the background drawing (game.go) it efficiently
- Put the debug into the down side of screen
- Add the spritemap inside the binary

## Compile

### MAC

`go build -ldflags="-s -w"`

#### Compress binaries
`brew install upx`
`upx <binary_file>`

#### Instructions to pack into Mac executable:
1) Baixar o binário
https://github.com/kindlychung/genicon

2) Install dependency:
brew install imagemagick

3) Create the icon based on a PNG image
`./genicon icon.png tmp_folder`

4) Rename the folder mv AppIcon.appiconset AppIcon.iconset

5) Create the icon in the format needed by Mac executable
iconutil -c icns -o icon.icns AppIcon.iconset


### Windows

GO allows to create a Windows executable file using a MacOS:

- Install mingw-w64 (support the GCC compiler on Windows systems):
`brew install mingw-w64`

- Prepare the icon for the binary:
`go install github.com/tc-hib/go-winres@latest`
`sudo go-winres init`
`sudo chown -R $(id -un) winres && chmod 755 winres`

Edit and replace the APP section of the file : `winres/winres.json`

`"APP": {
    "0000": [
      "../Images/AppIcon.iconset/icon_128x128.png",
      "../Images/AppIcon.iconset/icon_64x64.png",
      "../Images/AppIcon.iconset/icon_32x32.png",
      "../Images/AppIcon.iconset/icon_16x16.png"
    ]
  }`
`go-winres make`

- 32 bits:
`env GOOS="windows" GOARCH="386"   CGO_ENABLED="1" CC="i686-w64-mingw32-gcc"   go build -ldflags="-s -w"`

- 64 bits:
`env GOOS="windows" GOARCH="amd64" CGO_ENABLED="1" CC="x86_64-w64-mingw32-gcc" go build -ldflags="-s -w"`

* If you receive the message when running the executable, you need to ensure that the video drivers supports OpenGL (or the virtual driver in the case of virtualization).

* If you receive this message : "APIUnavailable: WGL: The driver does not appear to support OpenGL", please update your graphics driver os just copy the Mesa3D library from https://fdossena.com/?p=mesa/index.frag  (opengl32.dll) to the executable folder.

#### Compress binaries
`brew install upx`
`upx <binary_file>`


### Linux

PENDING


## Documentation:

- Pixel:

https://github.com/faiface/pixel/wiki

- Game programing:

https://www.codingdream.com/index.php/simple-pacman-in-using-go-and-pixelgl-part-1

- Mac APPs

https://medium.com/@mattholt/packaging-a-go-application-for-macos-f7084b00f6b5

- Go Windows Binary Icon:
https://stackoverflow.com/questions/25602600/how-do-you-set-the-application-icon-in-golang
https://github.com/mxre/winres
