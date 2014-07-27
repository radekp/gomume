package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

type Session struct {
	conn   net.Conn
	kmode  bool
	target string
	door   string
}

func reader(conn net.Conn) {
	var bytes [65535]byte
	buf := bytes[:]
	log, err := os.Create("log.txt")
	if err != nil {
		panic(err)
	}
	defer log.Close()

	for {
		n, _ := conn.Read(buf)
		input := bytes[:n]
		os.Stdout.Write(input)
		go log.Write(input)
	}
}

func handleSimpleCmd(ss *Session, line string) string {
	switch line {
	case "o", "p", "c", "x", "z":
		return line + " " + ss.door
	case ";", "k", "b":
		return line + " " + ss.target
    case "t":
        return "label " + ss.target + " aa";
	case "1":
		ss.target = "*elf*"
	case "2":
		ss.target = "*man*"
	case "3":
		ss.target = "*dwarf*"
	case "4":
		ss.target = "*hobbit*"
	case "5":
		ss.target = "*bear*"
    case "obo", "cbo":
        ss.door = "boulder";
    case "obr", "cbr":
        ss.door = "brush";
    case "oobs", "cobs":
        ss.door = "obsidian";
    case "obu", "cbu":
        ss.door = "bushes";
    case "oca", "cca":
        ss.door = "cask";
    case "oce", "cce":
        ss.door = "ceiling";
    case "oco", "cco":
        ss.door = "corner"
    case "ocr", "ccr":
        ss.door = "crack";
    case "on", "cn":
        ss.door = "exit north";
    case "os", "cs":
        ss.door = "exit south";
    case "oe", "ce":
        ss.door = "exit east";
    case "ow", "cw":
        ss.door = "exit west";
    case "ou", "cu":
        ss.door = "exit up";
    case "od", "cd":
        ss.door = "exit down";
    case "ogr", "cgr":
        ss.door = "grasses";
    case "oha", "cha":
        ss.door = "hatch";
    case "ohe", "che":
        ss.door = "hedge";
    case "oi", "ci":
        ss.door = "icedoor";
    case "ooc", "cooc":
        ss.door = "looserocks";
    case "opa", "cpa":
        ss.door = "panel";
    case "ops", "cps":
        ss.door = "passage";
    case "orf", "crf":
        ss.door = "rockface";
    case "oru", "cru":
        ss.door = "runes";
    case "ose", "cse":
        ss.door = "secret";
    case "ost", "cst":
        ss.door = "statuary";
    case "otd", "ctd":
        ss.door = "trapdoor";
    case "ote", "cte":
        ss.door = "tendrils";
    case "oth", "cth":
        ss.door = "thorns";
    case "oto", "cto":
        ss.door = "thornbushes";
    case "owa", "cwa":
        ss.door = "wall";
	case "n", "s", "e", "w", "u", "d":
		return line
	}
	return line
}

func handleCmdWithArg(ss *Session, line, cmd, arg string) string {
	switch cmd {
	case "t", ";", "k", "b":
		ss.target = arg
		return line
	case "o", "p", "c", "x", "z":
		ss.door = arg
		return line
	}
	return line
}

func main() {

	debug, _ := os.Create("debug.txt")
	defer debug.Close()

	fmt.Printf("Connecting mume.org:23\n")

	conn, err := net.Dial("tcp", "mume.org:23")
	if err != nil {
		panic("connect failed")
	}

	// Session object
	sso := Session{conn, false, "", "exit"}
	ss := &sso

	go reader(conn)

	stdinBuf := make([]byte, 65535)
	newLine := true
	line := ""
	for {
		n, err := os.Stdin.Read(stdinBuf)
		if err != nil {
			panic(err)
		}

		buf := stdinBuf[:n]

		// Enter
		newLine = (buf[0] == 10)
		if newLine {
			ss.kmode = false
			debug.Write([]byte("newline before: " + line + "\r\n"))
			index := strings.Index(line, " ")
			if index > 0 {
				line = handleCmdWithArg(ss, line, line[0:index], line[index+1:])
			} else {
				line = handleSimpleCmd(ss, line)
			}
			debug.Write([]byte("newline after: " + line + "\r\n"))
			conn.Write([]byte(line + "\r\n"))
			line = ""
			continue
		}

		// Backspace
		if buf[0] == 127 {
			newLen := len(line) - 1
			if newLen >= 0 {
				line = line[0:newLen]
				fmt.Printf("\b\b\b   \b\b\b")
			} else {
				fmt.Printf("\b\b  \b\b")
			}
			continue
		}

		// Arrows
		if n >= 3 && buf[0] == 27 && buf[1] == 91 {
			ss.kmode = true
			switch buf[2] {
			case 68:
				conn.Write([]byte("w\n"))
				fmt.Printf("\b\b\b\bw   \n")
				continue
			case 67:
				conn.Write([]byte("e\n"))
				fmt.Printf("\b\b\b\be   \n")
				continue
			case 65:
				conn.Write([]byte("n\n"))
				fmt.Printf("\b\b\b\bn   \n")
				continue
			case 66:
				conn.Write([]byte("s\n"))
				fmt.Printf("\b\b\b\bs   \n")
				continue
			}
		}

		// Binds in key mode
		/*if(ss.kmode) {
		    switch(buf[0]) {
		        case 'n':  conn.Write([]byte("n\n")); fmt.Printf("\bn\n"); continue;
		        case 's':  conn.Write([]byte("s\n")); fmt.Printf("\bs\n"); continue;
		        case 'e':  conn.Write([]byte("e\n")); fmt.Printf("\be\n"); continue;
		        case 'w':  conn.Write([]byte("w\n")); fmt.Printf("\bw\n"); continue;
		        case 'u':  conn.Write([]byte("u\n")); fmt.Printf("\bu\n"); continue;
		        case 'd':  conn.Write([]byte("d\n")); fmt.Printf("\bd\n"); continue;
		        case 'f':  conn.Write([]byte("f\n")); fmt.Printf("\bflee\n"); continue;
		        case '.':  conn.Write([]byte("label aa\n")); fmt.Printf("\blabel aa\n"); continue;
		        case 'v':  conn.Write([]byte("bash aa\n")); fmt.Printf("\bbash aa\n"); continue;            
		        case ';':  conn.Write([]byte("kill " + ss.target + "\n")); fmt.Printf("\bkill " + ss.target + "\n"); continue;            
		        case 'b':  conn.Write([]byte("bash " + ss.target + "\n")); fmt.Printf("\bbash " + ss.target + "\n"); continue;                        
		    }
		}*/

		ss.kmode = false

		line += string(buf[0])
		debug.Write([]byte(line))
		debug.Write([]byte("\r\n"))

		//fmt.Fprintf(debug, "send> %s\n", string(buf));
		fmt.Fprintf(debug, "target> %s door> %s\n", ss.target, ss.door)

		//conn.Write(buf)

		//os.Stdout.Write(buf2)
		for i := 0; i < n; i++ {
			fmt.Fprintf(debug, "key[%d]=%d", i, buf[i])
		}
	}

	//in    := bufio.NewReader(os.Stdin)
	//input := ""

	/*for input != "." {
	    input, err := in.ReadString('\n')
	    if err != nil {
	        panic(err)
	    }
	    buf := []byte(input)
	    conn.Write(buf)
	}*/
}
