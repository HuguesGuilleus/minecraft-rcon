# minecraft-rcon
Make a server between telnet and a minecraft server.

```txt
+--------+     +----------------+     +------------------+
| telnet | <-> | minecraft-rcon | <-> | minecraft server |
+--------+     +----------------+     +------------------+
```

## Build

Use relase binary or use go build yourself:

```txt
git clone https://github.com/HuguesGuilleus/minecraft-ron.git
cd minecraft-ron
go get -v -d .
go build
```

## Run
```txt
  -mp string
      Minecraft password
  -ms string
      Minecraft server
  -p string
      Public addrss (default "localhost:7000")
  -s string
      Secret password
```

Then use telnet to connect to these server and write standard minecraft command (without '/' prefix). Type `exit` to exit and close connexion.
