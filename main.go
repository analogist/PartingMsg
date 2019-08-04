package main
import (
    "fmt"
    "encoding/base64"
    "golang.org/x/crypto/nacl/secretbox"
    sssa "github.com/SSSaaS/sssa-golang"
    "bufio"
    "os"
    "github.com/analogist/partingmsg/partmsg"
)

func main() {
    var decryptnonce [24]byte // nonce for XSalsa, will be read in from partmsg.Encryptedencode
    var reconkey [32]byte // will hold the reconstructed key from SSSS pieces
    var pieces []string // pieces of SSSS to be input from stdin

    reconencrypted, _ := base64.StdEncoding.DecodeString(partmsg.Encryptedencode)
    copy(decryptnonce[:], reconencrypted[:24])

    fmt.Println("Please enter the pieces, one on each line. End")
    fmt.Println("with Ctrl+D, or by typing DONE, on its own line.")
    fmt.Println("==========")
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        if scanner.Text() == "DONE" {
            break
        }
        pieces = append(pieces, string(scanner.Text()))
    }
    if err := scanner.Err(); err != nil {
        panic(err)
    }

    combinekey, err := sssa.Combine(pieces)
    if err != nil {
        panic(err)
    }
    copy(reconkey[:], combinekey)
    decrypted, ok := secretbox.Open(nil, reconencrypted[24:], &decryptnonce, &reconkey)
    if !ok {
        fmt.Println("Not enough or incorrect shares to reconstruct.")
    }
    fmt.Println(string(decrypted))
}
