package main

import (
	"bufio"
	"bytes"
	"flag"
	"github.com/gorcon/rcon"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

func main() {
	pub := flag.String("p", "localhost:7000", "Public addrss")
	secret := flag.String("s", "", "Secret password")
	ms := flag.String("ms", "", "Minecraft server")
	mp := flag.String("mp", "", "Minecraft password")
	flag.Parse()

	serv := NewServ(*ms, *mp, rcon.SetDeadline(time.Hour*24))
	serv.secret = *secret

	log.Println("[LISTEN]", *pub)
	l, err := net.Listen("tcp", *pub)
	if err != nil {
		log.Fatal(err)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("[ERROR]", err)
			continue
		}
		log.Println("[NEW ADDR]", c.RemoteAddr())
		go serv.client(c)
	}
}

type Serv struct {
	c      *rcon.Conn
	secret string
	sync.Mutex
}

func NewServ(address, password string, options ...rcon.Option) Serv {
	c, err := rcon.Dial(address, password, options...)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for range time.Tick(time.Minute) {
			c.Execute("help")
		}
	}()
	return Serv{c: c}
}

func (s *Serv) client(c net.Conn) {
	defer c.Close()
	defer c.Write([]byte("> Good bye\r\n"))
	r := bufio.NewReader(c)

	c.Write([]byte("password: "))
	l, _, err := r.ReadLine()
	if err == io.EOF {
		log.Println("[CLOSE]", c.RemoteAddr())
		return
	} else if err != nil {
		log.Println("[CONNEXION ERROR]", c.RemoteAddr(), err)
		return
	}
	if string(l) != s.secret {
		log.Println("[WRONG AUTH]", c.RemoteAddr())
		c.Write([]byte("  bad password"))
		return
	}

	for {
		c.Write([]byte("% "))
		l, _, err = r.ReadLine()
		if err == io.EOF {
			log.Println("[CLOSE]", c.RemoteAddr())
			return
		} else if err != nil {
			log.Println("[CONNEXION ERROR]", c.RemoteAddr(), err)
			return
		}

		switch string(bytes.TrimSpace(l)) {
		case "":
			continue
		case "exit":
			return
		default:
			c.Write([]byte(
				"> " + s.Execute(string(l)) + "\r\n",
			))
		}
	}
}

// Execute a command, return erorr code
func (s *Serv) Execute(cmd string) string {
	s.Lock()
	defer s.Unlock()

	rep, err := s.c.Execute(cmd)
	if err != nil {
		return "Error: " + err.Error()
	}
	return rep
}
