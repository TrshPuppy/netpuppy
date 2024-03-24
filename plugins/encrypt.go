package plugins

import (
    "bufio"
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "os"
    "strings"
	"netpuppy/utils"
)

func init() {
    Register("encrypt", &Encrypt{})
}

type Encrypt struct{}

func (ep *Encrypt) Execute(pluginDataChan chan<- string) {
    fmt.Println("[Encrypt] type a message followed by a flag: -md5, -sha1, -sha256, or -binary. Type 'exit' to quit.")
    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("[Encrypt] ")
        scanner.Scan()
        input := scanner.Text()

        if strings.TrimSpace(input) == "exit" {
            fmt.Println("[Encrypt] Goodbye!")
            break
        }

        parts := strings.Fields(input)
        if len(parts) < 2 {
            fmt.Println(utils.Color("[Encrypt Error] Unsupported encryption/encoding flag. Use '-md5', '-sha1', '-sha256', or '-binary'.", utils.Red))
            continue
        }

        text := strings.Join(parts[:len(parts)-1], " ")
        method := parts[len(parts)-1]

        var result string

        switch method {
        case "-md5":
            result = ep.MD5(text)
        case "-sha1":
            result = ep.SHA1(text)
        case "-sha256":
            result = ep.SHA256(text)
        case "-binary":
            result = ep.Binary(text)
        default:
            fmt.Println(utils.Color("[Encrypt Error] Unsupported encryption/encoding flag. Use '-md5', '-sha1', '-sha256', or '-binary'.", utils.Red))
            continue
        }

		encodedLabel := utils.Color(fmt.Sprintf("[Encoded %s]", method), utils.Green)
		encodedMessage := fmt.Sprintf("%s %s", encodedLabel, result)
		pluginDataChan <- encodedMessage
		fmt.Println(encodedMessage)
    }
}

func (ep *Encrypt) MD5(text string) string {
    hasher := md5.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

func (ep *Encrypt) SHA1(text string) string {
    hasher := sha1.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}

func (ep *Encrypt) SHA256(text string) string {
    hasher := sha256.New()
    hasher.Write([]byte(text))
    return hex.EncodeToString(hasher.Sum(nil))
}


func (ep *Encrypt) Binary(text string) string {
    var binaryString strings.Builder
    for _, char := range text {
        binaryString.WriteString(fmt.Sprintf("%08b ", char))
    }
    return binaryString.String()
}
