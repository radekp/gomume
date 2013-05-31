package main

import (
    "os"
    "fmt"
    "net"
)

type Session struct {
    target string
    door string
}

func reader(conn net.Conn) {
    var bytes [65535]byte
    buf := bytes[:]    
    log, err := os.Create("log.txt")
    if err != nil { panic(err) }
    defer log.Close()
    
    for {
        n, _ := conn.Read(buf)
        input := bytes[:n]
        os.Stdout.Write(input)
        go log.Write(input)
    }    
}

func handleSimpleCmd(cmd string) bool {
    fmt.Printf("handleSimpleCmd %s", cmd);
    return false
}

func handleCmdWithArg(cmd, arg string) bool {

    fmt.Printf("handleCmdWithArg %s %s", cmd, arg);
    
    switch cmd {
        case "t": return true
    }
    return false
}

func main() {

    debug, _ := os.Create("debug.txt")
    defer debug.Close()
    
    fmt.Printf("Connecting mume.org:23\n")
    
    conn, err := net.Dial("tcp", "mume.org:23")
    if err != nil {
        panic("connect failed")
    }
    go reader(conn)
    
    stdinBuf := make([] byte, 65535)
    newLine := true
    line := ""
    for {
        n, err := os.Stdin.Read(stdinBuf)
        if err != nil { panic(err) }

        buf := stdinBuf[:n]
        
        if newLine {
            if buf[0] == 49 {
                //conn.Write([]byte("kill *elf*\r\n"))
                cmd := make([] byte, 2)
                cmd[0] = 108
                cmd[1] = 10
                conn.Write(cmd)
                continue
            }
        }
        
        newLine = (buf[0] == 10)
        if newLine {
            if(len(line) > 1) {
                if(line[1] == ' ') {
                    if(handleCmdWithArg(line[0:1], line[2:])) {
                        line = ""
                        continue
                    }
                }
            } else if(len(line) == 1) {
                if(handleSimpleCmd(line)) {
                    line = ""
                    continue
                }
            }
            debug.Write([]byte("newline \r\n"))
            line = ""
        } else {
            line += string(buf[0])
        }
        debug.Write([]byte(line))
        debug.Write([]byte("\r\n"))

        fmt.Fprintf(debug, "send> %s\n", string(buf));
        
        conn.Write(buf)
        
        
        
        //os.Stdout.Write(buf2)
        //for i := 0 ; i < n; i++ {
            //fmt.Printf("key[%d]=%d", i, buf[i])
        //}            
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