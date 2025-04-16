# Text Protocol

## Introduction
`textprotocol` is a project written in Golang. The server allows you to make multiple connection through terminal and execute a commands that execute operations inside the `tcp` server.

---

## Environment Setup
There are two options for your environment setup.
1) Local Setup - Setup your local machine to be able to run the application.
2) Docker Setup - Setup Docker for running the application.

### Configuration
Before running the application make sure the values in the `config.json` are set properly.
```json
{
    "port": "4242",
    "session_active_interval_hours": 2,
    "rate_limit_max_sessions": 5,
    "rate_limit_refill_duration_secods": 15,
    "rate_limit_max_input_per_interval": 5,
    "connection_life_span_minutes": 30
}
```
The code snipped above is an example of what the config.json should look. </br>

| Variable Name | Description | Value Type | Example Value |
|---------------|-------------|------------|---------------|
|    `port`     | Defining the Port which will be allocated from the server | `string` | '4242' |
| `session_active_interval_hours` | This variable defines the live span of the connection if no input is comming from the stream. | `int` | 2 |
| `rate_limit_max_sessions` | Defining the Maximum amount of sessions a user can open from a single IP | `int` | 5 |
| `rate_limit_refill_duration_secods` | After this interval has passed the input per interval variable will be refresh with the `rate_limit_max_input_per_interval` | `int` | 15 |
| `rate_limit_max_input_per_interval` | Defines how much commands can be intered in the specified duration interval from the `rate_limit_refill_duration_secods` variable | `int` | 5 |
| `connection_life_span_minutes` | Defines the connection life span, how much an active connection can last. | 30 | 

There is also an additional `config_test.json` this `json` has the same variables as the `config.json`. However, this file is used for when running the tests.


### Local Setup
In order to run the project you will need to have a [Golang version](https://go.dev/dl/). </br> 
Make sure you follow the process of installing `Golang` on your Local machine. </br>
Choose the downloading Go package based on your OS. </br> 
After the installation process, make sure you clone the project on your machine. </br>
```bash
git clone git@github.com:boPopov/textprotocol.git
```

### Docker Setup
If you make a decision to go with the Docker setup follow the next steps:
#### Windows
Download [Docker Desktop](https://docs.docker.com/get-started/introduction/get-docker-desktop/). Make sure to install the applicat
ion. After the installation is completed open a terminal and run the command below:
```bash
docker --version
```
The result should be:
```
Docker version <Docker_Version>, build <id>
```
#### Linux
Follow the [installation guid](https://docs.docker.com/engine/install/ubuntu/) for installing Docker on Ubuntu or any other Linux OS you have. </br>

After the installation is complted check if docker is successfully installed and setup on your machine:
```bash
docker --version
```

### Install Telnet
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

---

## Running Application
As specified in the `Environment Setup` there are a two options to run the application, one is local and the second is through Docker.

### Local Execution
#### Linux
Open terminal, navigate to the project
```bash
cd PATH/TO/textprotocol
```
Follow with the second command
```bash
go run ./src $(pwd)/config.json
```

#### Windows
Navigate to the folder for `textprotocol`. Open a `PowerShell` or `Cmd` or `Windows Terminal` and execute the command below.
```bash
go run .\\src\\ .\\config.json
```

### Docker Execution
Befor explaining the process for starting the application in Docker make sure the docker engine is running.</br>
For windows start `Docker Desktop`, for Linux try to access the `docker`. </br>
There are no different commands for Windows and Linux:
```bash
# To start the application run
docker-compose up -d
```
The application will be reachable through `Nginx` on `localhost:8080`.

```bash
# To stop the application run
docker-compose down
```

#### Docker Application
There is also a second option that is available. Start the golang application through `Docker`.
```bash
# Build the application
docker build -t textprotocol .
```

```bash
# Run the application
docker run -p 4242:4242 textprotocol
```

```bash
docker ps # Find the textprotocol docker_id
# Stop the application
docker stop <docker_id>
```

The provided commands must be used for Linux, as for Windows you have to use the build command, after that you can do all of the commands, runing the application and stoping the application, inside the `Docker Desktop`.

## Usage Guide
Once you have installed `telnet`, `Golang` and the `TCP` server is up and running. Open a terminal PowerShell, CMD, Windows Terminal or Terminal (Linux) and run:

### If application is running on your Local Machine
```bash
telnet localhost PORT
```

### If application is running on a machine inside your Network
Find your local ip and then execute the command below.
```bash
telnet YOUR.LOCAL.IP.HERE PORT
```

### Port
The `PORT` can have any value of an available port on your machine. Also, if you run the application in docker then it has a different port. </br>

If you run the application without `Docker`, use the value from inside the `config.json` -> `port`. </br>

If you run the application with `docker-compose` use the value `8080`. If you can not use the port `8080`, navigate to `docker-compose.yml` and change the value under `ports` to:
```yml
lbl:
    image: nginx
    container_name: load_balancer
    depends_on:
      - tcp-server
    ports:
      - YOUR_AVAILABLE_PORT:8080
```

And finally, if you use the `Docker Application` process for running the application use the port defined in the `config.json`.

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