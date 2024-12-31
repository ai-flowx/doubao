package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Provider struct {
	Name  string
	URL   string
	Model string
	Key   string
}

type Prompt struct {
	Input []string `json:"input"`
}

type Response struct {
	Id      string `json:"id"`
	Model   string `json:"model"`
	Created int    `json:"created"`
	Object  string `json:"object"`
	Data    []Data `json:"data"`
	Usage   Usage  `json:"usage"`
}

type Data struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
	Object    string    `json:"object"`
}

type Usage struct {
	PromptTokens int `json:"prompt_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

var (
	providerFile string
	promptFile   string
)

var rootCmd = &cobra.Command{
	Use:   "embedding",
	Short: "doubao embedding",
	Long:  "doubao embedding",
	Run: func(cmd *cobra.Command, args []string) {
		var provider Provider
		var prompt Prompt
		var err error
		ctx := context.Background()
		if provider, err = initProvider(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		if prompt, err = initPrompt(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
		if err := runModel(ctx, provider, prompt); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(1)
		}
	},
}

func main() {
	Execute()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// nolint:gochecknoinits
func init() {
	rootCmd.PersistentFlags().StringVarP(&providerFile, "provider-file", "p", "", "provider file")
	rootCmd.PersistentFlags().StringVarP(&promptFile, "prompt-file", "m", "", "prompt file")

	rootCmd.Root().CompletionOptions.DisableDefaultCmd = true
}

func initProvider() (Provider, error) {
	var provider Provider

	if providerFile == "" {
		return Provider{}, errors.New("invalid provider file\n")
	}

	content, err := os.ReadFile(providerFile)
	if err != nil {
		return Provider{}, err
	}

	if err = yaml.Unmarshal(content, &provider); err != nil {
		return Provider{}, err
	}

	return provider, nil
}

func initPrompt() (Prompt, error) {
	var prompt Prompt

	if promptFile == "" {
		return Prompt{}, errors.New("invalid prompt file\n")
	}

	content, err := os.ReadFile(promptFile)
	if err != nil {
		return Prompt{}, err
	}

	if err = json.Unmarshal(content, &prompt); err != nil {
		return Prompt{}, err
	}

	return prompt, nil
}

func runModel(_ context.Context, provider Provider, prompt Prompt) error {
	var res Response

	p, err := json.Marshal(prompt.Input)
	if err != nil {
		return err
	}

	buf := fmt.Sprintf(`{
		"model": "%s",
		"input": %s
	}`, provider.Model, p)

	req, err := http.NewRequest(http.MethodPost, provider.URL, bytes.NewBuffer([]byte(buf)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+provider.Key)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	fmt.Printf("Response status: %s\n", resp.Status)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal((body), &res); err != nil {
		return err
	}

	fmt.Printf("Response body: %v\n", res)

	return nil
}
