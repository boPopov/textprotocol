# Text Protocol

## Introduction
`textprotocol` is a project written in Golang. The server allows you to make multiple connection through terminal and execute a commands that have execute operations inside the `tcp` server.

---

## Environment Setup
The current setup of the Application is Local. The plan is to create a Docker version of the application with Docker-compose in order to include nginx.
### Local Setup
In order to run the project you will need to have a [Golang version](https://go.dev/dl/). </br> 
Make sure you follow the process of installing `Golang` on your Local machine. </br>
After the installation process, make sure you clone the project on your machine. </br>

Telnet commands:
#### Windows
In order to install telnet on Windows, follow the instructions below:
1. Navigate to control panel
2. Navigate to Programs
3. Navigate to Enable/Disabled Feature.
4. Enable the telnet below the list.

#### Linux
In order to install telnet on Ubuntu/Debian run the following commands.
```bash
sudo apt update
sudo apt install telnet
```

## Running Application
In order to run the application. 
### Local Run
```bash
go run ./src $(pwd)/config.json
```

## Usage Guide
Once you have installed `telnet`, `Golang` and the `TCP` server is up and running. Open a terminal PowerShell, CMD or Terminal (Linux) and run:

### If application is running on your Local Machine
```bash
telnet localhost 4242
```

### If application is running on a machine inside your Network
Find your local ip.
```bash
telnet YOUR.LOCAL.IP.HERE 4242
```

---
There are a couple of commands avaible for usage. Please look at the table bellow.

| Command Name | Description | Responses |
|--------------|-------------|-----------|
| `EHLO` | The command EHLO, expects a string after the command. After the EHLO name command is entered you can execute the DATE command | 250 Pleased to meet you NAME |
| `DATE` | This command returns the current date and time if the command EHLO was successfully executed before. | Positive response - 250 21/10/2016T14:13:08; Negative Response - 550 Bad state |
| `QUIT` | This command stops the connection between your Machine and the Server | 221 Bye | 

## Testing
There are multiple unit tests that check the command behavior as well as the rate limit for connections. </br>

You can execute the test with the following commands. </br>

```bash
cd tests
CONFIG_PATH=/PATH/TO/PROJECT/textprotocol/config.json go test .
```