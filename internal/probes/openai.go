package probes

import (
	"context"
	"log/slog"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func OpenAI(prompt, endpoint, key, deployment string) (string, error) {
	creds := azcore.NewKeyCredential(key)
	client, err := azopenai.NewClientWithKeyCredential(endpoint, creds, nil)

	if err != nil {
		slog.Error("creating OpenAI client", "error", err)
		return "", err
	}

	messages := []azopenai.ChatRequestMessageClassification{
		&azopenai.ChatRequestSystemMessage{Content: to.Ptr("You are a helpful assistant for Microsoft Azure and Terraform.")},
		&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent("Does Terraform support Azure?")},
		&azopenai.ChatRequestAssistantMessage{Content: to.Ptr(
			"Yes, Terraform supports Azure. Terraform is an open-source infrastructure as code (IaC) tool developed by HashiCorp " +
				"that allows you to define and provision infrastructure using a high-level configuration language. It has extensive " +
				"support for a variety of cloud providers, including Microsoft Azure.",
		)},
		&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(prompt)},
	}

	resp, err := client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: to.Ptr(deployment),
	}, nil)
	if err != nil {
		slog.Error("getting chat completions", "error", err)
		return "", err
	}

	var sb strings.Builder
	for _, c := range resp.Choices {
		sb.WriteString(*c.Message.Content)
		sb.WriteString("\n")
	}
	return sb.String(), nil
}
