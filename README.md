# The graphical elevator simulator

This is a simple elevator simulator written for the Sanntidslab at NTNU.
The simulator is compatible with the drivers from the TTK4145 course.

Feedback is appreciated and errors might occur, after all it's software...
Leave an issue if you are having trouble or if something should be changed!

## Installation

## Simple

The simplest way to get going is to download one of the precompiled binaries from the releases tab.
Run the binary in your console and you should be good to go.

> Not sure which binary to download? Grab the latest version and take the one having amd64 in it's name. (Darwin = MAC)

## Build from source

Have a look at the Makefile it should be pretty self-explanatory.

## Usage

The simulator acts as a web server, and the user interface consits of a web page.
When started the server will print how to open the page in the console window.

### Start the server

To start the service, for example:

Linux:

```bash
./elevator-simulator-linux-amd64
```

Mac:

```bash
./elevator-simulator-darwin-amd64
```

### HELP

It won't run?

Check if you have downloaded and started the right version. Both the operating system and architecture must match.
Most modern systems are amd64 architectures (yes, also when it is a indel processor).

The permissions are probably missing. The following will resolve the issue:

```bash
sudo chmod +x ./write_the_name_of_the_executable_here
```

### Flags, configuration and more

The simulator is configurable through flags. Run the program with "--help" to see the available options.

### How to run multiple instances

When trying to run several elevators you may have noticed that the simulator isn't able to bind the ports.
This is because only one application can listen to a port at a time. To solve this run the simulators on different ports (use --help to see how to configure the simulator to do this).
