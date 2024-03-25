package plugins

import (
    "bufio"
    "crypto/md5"
    "crypto/sha1"
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "os"
    "strconv"
    "strings"
)

func init() {
    Register("encrypt", &Encrypt{})
}

type Encrypt struct{}

// ANSI color codes
const (
    Red   = "\033[31m"
    Green = "\033[32m"
    Reset = "\033[0m"
)

func (ep *Encrypt) Execute(pluginDataChan chan<- string) {
    fmt.Println("[Encrypt] Type a message followed by a flag: -md5, -sha1, -sha256, -binary, or -d (for decrypt). Type 'exit' to quit.")
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
            fmt.Println(Red + "[Encrypt Error] Unsupported encryption/decoding flag. Use '-md5', '-sha1', '-sha256', '-binary', or '-d'." + Reset)
            continue
        }

        text := strings.Join(parts[:len(parts)-1], " ")
        method := parts[len(parts)-1]

        var result string
        var err error

        switch method {
        case "-md5":
            result = ep.MD5(text)
        case "-sha1":
            result = ep.SHA1(text)
        case "-sha256":
            result = ep.SHA256(text)
        case "-binary":
            result = ep.Binary(text)
        case "-d":
            result, err = ep.BinaryToString(text)
            if err != nil {
                fmt.Println(Red + "[Decrypt Error] Invalid binary string." + Reset)
                continue
            }
        default:
            fmt.Println(Red + "[Encrypt Error] Unsupported encryption/decoding flag. Use '-md5', '-sha1', '-sha256', '-binary', or '-d'." + Reset)
            continue
        }

        var encodedLabel string
        if method == "-d" {
            encodedLabel = Green + "[Decrypted Binary]" + Reset
        } else {
            encodedLabel = Green + fmt.Sprintf("[Encoded %s]", method) + Reset
        }
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

func (ep *Encrypt) BinaryToString(binaryStr string) (string, error) {
    var text strings.Builder
    binaryParts := strings.Fields(binaryStr)

    for _, part := range binaryParts {
        num, err := strconv.ParseInt(part, 2, 64)
        if err != nil {
            return "", err
        }
        text.WriteByte(byte(num))
    }

    return text.String(), nil
}
