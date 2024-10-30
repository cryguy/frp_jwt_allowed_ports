# frp_jwt_allowed_ports


frp server plugin to support multiple users for [frp](https://github.com/fatedier/frp).


frp_jwt_allowed_ports will run as one single process and accept HTTP requests from frps.

### Features

* Support multiple user authentication by jwt secret key - it allows a server or application to generate a jwt key with grants of ports to the user.

### JWT Format

the jwt payload should contain the ports allowed, this is still to be revised:
    
    {
      "sub": "user2",
      "ports": {
         "tcp": [6001]
      },
      "iat": 1516239022
    }

### Download

Download frp_jwt_allowed_ports binary file from [Release](https://github.com/cryguy/frp_jwt_allowed_ports/releases).

### Requirements

frp version >= v0.31.0

### Example Usage

1. Create file `secret` which contains the jwt secret.

    ``` EXAMPLE ONLY! PLEASE CHANGE!
    61d371c34edebe1b1f8003cd95129415c46d2cae729bb2a455f237dfb264fb42
    ```

2. Run fp-multiuser:

    `./frp_jwt_allowed_ports -l 127.0.0.1:7200 -k ./secret`

3. Register plugin in frps.

   INI:

    ```ini
    # frps.ini
    [common]
    bind_port = 7000

    [plugin.frp_jwt_allowed_ports]
    addr = 127.0.0.1:7200
    path = /handler
    ops = Jwt
    ```

    TOML:

    ```toml
    # frps.toml
    bindPort = 7000

    [[frp_jwt_allowed_ports]]
    addr = "127.0.0.1:7200"
    path = "/handler"
    ops = ["Jwt"]
    ```

4. Specify username and meta_token in frpc configure file.

    For user1 INI:

    ```ini
    # frpc.ini
    [common]
    server_addr = x.x.x.x
    server_port = 7000
    jwt = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMSIsInBvcnRzIjp7InRjcCI6WzYwMDBdfSwiaWF0IjoxNTE2MjM5MDIyfQ.WVRo6Upcw71pQZHGHnAPRVVz5BXZk3l2kWy252Q5YJ8

    [ssh]
    type = tcp
    local_port = 22
    remote_port = 6000
    ```

    For user1 TOML:

    ```toml
    serverAddr = "x.x.x.x"
    serverPort = 7000
    jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMSIsInBvcnRzIjp7InRjcCI6WzYwMDBdfSwiaWF0IjoxNTE2MjM5MDIyfQ.WVRo6Upcw71pQZHGHnAPRVVz5BXZk3l2kWy252Q5YJ8"

    [[proxies]]
    type = "tcp"
    localPort = 22
    remotePort = 6000
    ```

    For user2 INI:

    ```ini
    # frpc.ini
    [common]
    server_addr = x.x.x.x
    server_port = 7000
    jwt = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMiIsInBvcnRzIjp7InRjcCI6WzYwMDFdfSwiaWF0IjoxNTE2MjM5MDIyfQ.jp2jka_m7MMtfhKJDbUtKRJ8lCe01S2seHSOBu08s5o

    [ssh]
    type = tcp
    local_port = 22
    remote_port = 6001
    ```

    For user2 TOML:

    ```toml
    serverAddr = "x.x.x.x"
    serverPort = 7000
    jwt = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ1c2VyMiIsInBvcnRzIjp7InRjcCI6WzYwMDFdfSwiaWF0IjoxNTE2MjM5MDIyfQ.jp2jka_m7MMtfhKJDbUtKRJ8lCe01S2seHSOBu08s5o"

    [[proxies]]
    type = "tcp"
    localPort = 22
    remotePort = 6001
    ```


