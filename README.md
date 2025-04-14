# Text Protocol

## Introduction
`textprotocol` is a project written in Golang. The server allows you to make multiple connection through terminal and execute a commands that have execute operations inside the `tcp` server.

---

## Environment Setup
The current setup of the Application is Local. The plan is to create a Docker version of the application with Docker-compose in order to include nginx.

---

### Configuration
Before running the application make sure the values in the `config.json` are set properly.
```json
{
    "port": "4242",
    "session_active_interval_hours": 2,
    "rate_limit_max_sessions": 5,
    "rate_limit_refill_duration_secods": 15,
    "rate_limit_max_input_per_interval": 5
}
```
The code snipped above is an example of what the config.json should look. </br>

| Variable Name | Description | Value Type | Example Value |
|---------------|-------------|------------|---------------|
|    `port`     | Defining the Port which will be allocated from the server | `string` | '4242' |
| `session_active_interval_hours` | This variable defines the live span of the connection that will be open from the client | `int` | 2 |
| `rate_limit_max_sessions` | Defining the Maximum amount of sessions a user can open from a single IP | `int` | 5 |
| `rate_limit_refill_duration_secods` | After this interval has passed the input per interval variable will be refresh with the `rate_limit_max_input_per_interval` | `int` | 15 |
| `rate_limit_max_input_per_interval` | Defines how much commands can be intered in the specified duration interval from the `rate_limit_refill_duration_secods` variable | `int` | 5 |

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