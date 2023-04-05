package router

import (
	"bytes"
	"encoding/json"
	"github.com/cloakscn/gpt-server/internal/model"
	. "github.com/cloakscn/gpt-server/internal/pkg/error"
	httpx "github.com/cloakscn/gpt-server/internal/pkg/http"
	"io"
	"net/http"
)

const (
	GPT_PREFIX     = "https://api.openai.com/v1/chat"
	OPENAI_API_KEY = "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"
)

func GptRouter() {
	http.HandleFunc("/chat/completions", completionsHandler)
	http.HandleFunc("/chat", chatHandler)
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		if err := httpx.RenderHtml(w, "chat", nil); err != nil {
			Check(err)
			return
		}
		return
	}

	if r.Method == "POST" {
		b := buildMessage(r.FormValue("message"))

		req, err := http.NewRequest("POST", GPT_PREFIX+"/completions", bytes.NewReader(b))
		if err != nil {
			Check(err)
			return
		}
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+OPENAI_API_KEY)

		client := &http.Client{}
		resp, err := client.Do(req)
		defer resp.Body.Close()
		if err != nil {
			Check(err)
			return
		}

		b, err = io.ReadAll(resp.Body)
		if err != nil {
			Check(err)
			return
		}
		body := model.Resp{}
		err = json.Unmarshal(b, &body)
		if err != nil {
			Check(err)
			return
		}

		var messages []string
		for _, choice := range body.Choices {
			messages = append(messages, choice.Message.Content)
		}
		locals := make(map[string]interface{})
		locals["messages"] = messages
		//if err := httpx.RenderHtml(w, "chat", locals); err != nil {
		//	Check(err)
		//	return
		//}
		marshal, _ := json.Marshal(locals)
		io.Copy(w, bytes.NewReader(marshal))
	}
}

func buildMessage(msg string) (body []byte) {
	body, _ = json.Marshal(model.Body{
		Model: "gpt-3.5-turbo",
		Messages: []model.Message{
			{
				Role:    "user",
				Content: msg,
			},
		},
		Temperature: 0.7,
	})
	return
}

func completionsHandler(w http.ResponseWriter, r *http.Request) {
	body := buildMessage(r.FormValue("message"))

	req, err := http.NewRequest("POST", GPT_PREFIX+"/completions", bytes.NewReader(body))
	if err != nil {
		Check(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+OPENAI_API_KEY)
	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		Check(err)
		return
	}
	io.Copy(w, resp.Body)
}
