package plugins

import (
    "bufio"
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"
	"regexp"
    "netpuppy/utils"
)

func init() {
    Register("chat", &ChatGPT{})
}

type ChatGPT struct {
    Messages []map[string]string // Store the conversation history
}

func (cg *ChatGPT) Execute(pluginDataChan chan<- string) {
    fmt.Println(utils.Color("Ask me anything! Type 'exit' to quit.", utils.Yellow))
    scanner := bufio.NewScanner(os.Stdin)

    cg.Messages = append(cg.Messages, map[string]string{"role": "system", "content": "You are a helpful assistant."})

    for {
        fmt.Print(utils.Color("> ", utils.Blue))
        scanner.Scan()
        userPrompt := scanner.Text()

        if strings.TrimSpace(userPrompt) == "exit" {
            fmt.Println(utils.Color("[Chat] Exiting...", utils.Red))
            break
        }

        // Add user's message to the conversation history
        cg.Messages = append(cg.Messages, map[string]string{"role": "user", "content": userPrompt})

        response, err := cg.sendPromptToOpenAI()
        if err != nil {
            fmt.Printf(utils.Color("[Chat] Error: %v\n", utils.Red), err)
            continue
        }

        highlightedResponse := highlightKeywords(response)

        cg.Messages = append(cg.Messages, map[string]string{"role": "assistant", "content": response})

		pluginDataChan <- highlightedResponse
		// fmt.Println(utils.Color(highlightedResponse)
		
    }
}

func (cg *ChatGPT) sendPromptToOpenAI() (string, error) {
    requestBody, err := json.Marshal(map[string]interface{}{
        "model":    "gpt-3.5-turbo",
        "messages": cg.Messages,
    })
    if err != nil {
        return "", err
    }

    req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(requestBody))
    if err != nil {
        return "", err
    }

    // Replace "YOUR_API_KEY_HERE" with your actual OpenAI API key
    req.Header.Set("Authorization", "Bearer YOUR_API_KEY_HERE")
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    responseBody, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var responseObj map[string]interface{}
    if err := json.Unmarshal(responseBody, &responseObj); err != nil {
        return "", err
    }

    if choices, ok := responseObj["choices"].([]interface{}); ok && len(choices) > 0 {
        if firstChoice, ok := choices[0].(map[string]interface{}); ok {
            if message, ok := firstChoice["message"].(map[string]interface{}); ok {
                if content, ok := message["content"].(string); ok {
                    return content, nil
                }
            }
        }
    }

    return "No response or unable to parse response.", nil
}

// This is a mess right now, if anyone can fix it that would be great x
func highlightKeywords(text string) string {
    codeBlockRegex := regexp.MustCompile("`{3}([^`]*)`{3}")

    patterns := map[*regexp.Regexp]string{
        regexp.MustCompile(`\b(if|else|switch|for|while|do|case|default|try|catch|finally|break|continue|return)\b`): utils.Yellow,
        regexp.MustCompile(`\b(func|function|def|class|interface|struct|enum)\b`): utils.Cyan,
        regexp.MustCompile(`\b(int|float|bool|string|char|void|var|let|const|static|public|private|protected)\b`): utils.Magenta,
        regexp.MustCompile(`\b(new|delete|throw|extends|implements|instanceof|typeof|sizeof|this|super)\b`): utils.Red,
        regexp.MustCompile(`\b\$\w+\b`): utils.Green,
        regexp.MustCompile(`\b\d+(\.\d+)?\b`): utils.Blue,
        regexp.MustCompile(`(['"]).*?\\1`): utils.Magenta,
        regexp.MustCompile(`\b([a-zA-Z_]\w*)\s*\(`): utils.Cyan,
        regexp.MustCompile(`<(?P<tag>\w+)(?P<attributes>(?:\s+\w+(?:\s*=\s*(?:"[^"]*"|'[^']*'|[^'">\s]+))?)+\s*|\s*)\/?>`): utils.Green,
        regexp.MustCompile(`<\/(?P<tag>\w+)>`): utils.Purple,
        regexp.MustCompile(`\b(std::cout|std::cin|cout|cin)\b|\b<<\b|>>`): utils.Blue,
        regexp.MustCompile(`(<\?php|\?>)`): utils.Red,
        regexp.MustCompile(`[=+\-*/]`): utils.Magenta,
        regexp.MustCompile(`"([^"]*)"|'([^']*)'`): utils.Magenta,
    }

    colorizeCodeBlock := func(blockMatch string) string {
        codeContent := blockMatch
        for pattern, color := range patterns {
            codeContent = pattern.ReplaceAllStringFunc(codeContent, func(match string) string {
                return utils.Color(match, color)
            })
        }
        return codeContent
    }

    modifiedText := codeBlockRegex.ReplaceAllStringFunc(text, colorizeCodeBlock)
    return modifiedText
}