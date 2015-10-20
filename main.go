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
func handleDoor(ss *Session, line string, door string) string {
	ss.door = door
	if line[0] == 'o' {
		return "open " + door
	}
	return "close " + door
}

func handleTarget(ss *Session, line string, action string, target string) string {
	ss.target = target
	return action + " " + target
}

func handleSimpleCmd(ss *Session, line string) string {
	switch line {
	case ".", "t":
		ss.target = "aa"
		return "label aa"
	case "o":
		return "open " + ss.door
	case "c":
		return "close " + ss.door
	case "p":
		return "pick " + ss.door
	case "x":
		return "lock " + ss.door
	case "z":
		return "unlock " + ss.door
	case ";", "k", "b":
		return line + " " + ss.target
	case "1":
		return handleTarget(ss, line, "kill", "*elf*");
	case "2":
		return handleTarget(ss, line, "kill", "*man*");
	case "3":
		return handleTarget(ss, line, "kill", "*dwarf*");
	case "4":
		return handleTarget(ss, line, "kill", "*hobbit*");
	case "5":
		return handleTarget(ss, line, "kill", "*bear*");
	case "on", "cn":
		return handleDoor(ss, line, "exit north");
	case "os", "cs":
		return handleDoor(ss, line, "exit south");
	case "oe", "ce":
		return handleDoor(ss, line, "exit east");
	case "ow", "cw":
		return handleDoor(ss, line, "exit west");
	case "ou", "cu":
		return handleDoor(ss, line, "exit up");
	case "od", "cd":
		return handleDoor(ss, line, "exit down");
    case "obo", "cbo":
		return handleDoor(ss, line, "boulder");
    case "obr", "cbr":
		return handleDoor(ss, line, "brush");
    case "oobs", "cobs":
		return handleDoor(ss, line, "obsidian");
    case "obu", "cbu":
		return handleDoor(ss, line, "bushes");
    case "oca", "cca":
		return handleDoor(ss, line, "cask");
    case "oce", "cce":
		return handleDoor(ss, line, "ceiling");
    case "oco", "cco":
		return handleDoor(ss, line, "corner");
    case "ocr", "ccr":
		return handleDoor(ss, line, "crack");
    case "ogr", "cgr":
		return handleDoor(ss, line, "grasses");
    case "oha", "cha":
        return handleDoor(ss, line, "hatch");
    case "ohe", "che":
		return handleDoor(ss, line, "hedge");
    case "oi", "ci":
		return handleDoor(ss, line, "icedoor");
	case "oro", "cro":
		return handleDoor(ss, line, "rock");
    case "ooc", "cooc":
		return handleDoor(ss, line, "looserocks");
    case "opa", "cpa":
		return handleDoor(ss, line, "panel");
    case "ops", "cps":
		return handleDoor(ss, line, "passage");
    case "orf", "crf":
		return handleDoor(ss, line, "rockface");
    case "oru", "cru":
		return handleDoor(ss, line, "runes");
    case "ose", "cse":
		return handleDoor(ss, line, "secret");
    case "ost", "cst":
		return handleDoor(ss, line, "statuary");
    case "otd", "ctd":
		return handleDoor(ss, line, "trapdoor");
    case "ote", "cte":
		return handleDoor(ss, line, "tendrils");
    case "oth", "cth":
		return handleDoor(ss, line, "thorns");
    case "oto", "cto":
		return handleDoor(ss, line, "thornbushes");
    case "owa", "cwa":
		return handleDoor(ss, line, "wall");
	case "n", "s", "e", "w", "u", "d":
		return line
	}
	return line
}

func handleCmdWithArg(ss *Session, line, cmd, arg string) string {
	switch cmd {
	case "t":
		ss.target = arg
		return "label " + arg + " aa"
	case ";", "k", "b":
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
			oldLine := line
			if index > 0 {
				line = handleCmdWithArg(ss, line, line[0:index], line[index+1:])
			} else {
				line = handleSimpleCmd(ss, line)
			}
			if oldLine != line {		// print command result
				fmt.Println(line);
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
