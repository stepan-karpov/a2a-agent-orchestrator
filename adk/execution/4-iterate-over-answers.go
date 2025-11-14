package execution

import (
	"adk/a2a/server"
	"adk/agents"
	"adk/providers"
	"context"
	"fmt"
	"strings"

	"google.golang.org/protobuf/types/known/structpb"
)

func IterateOverAnswers(ctx context.Context, provider providers.Provider, task *server.Task, message *server.Message, history []providers.Message, availableAgents []agents.Agent) (*server.Task, error) {
	fmt.Printf("IterateOverAnswers: received message with content='%s' (len=%d)\n", message.Content, len(message.Content))

	if message.Content == "" {
		return nil, fmt.Errorf("message content is empty")
	}

	history = append(history,
		providers.Message{
			Role:    "user",
			Content: message.Content,
		},
	)

	calledAgents := make([]string, 0)
	maxIterations := 10

	for i := 0; i < maxIterations; i++ {
		response, err := provider.ChatCompletion(history, availableAgents)
		if err != nil {
			return nil, err
		}

		if response.AgentQuery != nil {
			fmt.Printf("Query to agent '%s': message='%s'\n", response.AgentQuery.Agent.Name, response.AgentQuery.Message)
			var shouldContinue bool
			history, shouldContinue, err = handleAgentResponse(ctx, response.AgentQuery, history, &calledAgents, task.ContextId)
			if err != nil {
				return nil, err
			}
			if shouldContinue {
				continue
			}
		}

		if response.Content != "" {
			return handleRegularResponse(task, response.Content, calledAgents)
		}

		return nil, fmt.Errorf("response has no content and no agent query")
	}

	return nil, fmt.Errorf("max iterations reached (%d)", maxIterations)
}

func handleAgentResponse(ctx context.Context, agentQuery *providers.AgentQuery, history []providers.Message, calledAgents *[]string, contextId string) ([]providers.Message, bool, error) {
	*calledAgents = append(*calledAgents, agentQuery.Agent.Name)

	agentResponse, err := agents.CallAgent(ctx, agentQuery.Agent, agentQuery.Message, contextId)
	if err != nil {
		return nil, false, fmt.Errorf("failed to call agent %s: %w", agentQuery.Agent.Name, err)
	}

	history = append(history,
		providers.Message{
			Role:    "assistant",
			Content: fmt.Sprintf("Tool call to %s: %s", agentQuery.Agent.Name, agentQuery.Message),
		},
		providers.Message{
			Role:    "user",
			Content: fmt.Sprintf("Response from %s: %s", agentQuery.Agent.Name, agentResponse),
		},
	)

	return history, true, nil
}

func handleRegularResponse(task *server.Task, content string, calledAgents []string) (*server.Task, error) {
	finalContent := content

	if len(calledAgents) > 0 {
		agentsList := make([]interface{}, len(calledAgents))
		for i, name := range calledAgents {
			agentsList[i] = name
		}

		metadata := map[string]interface{}{
			"called_agents": agentsList,
		}

		metadataStruct, err := structpb.NewStruct(metadata)
		if err == nil {
			task.Metadata = metadataStruct
		}

		agentsText := fmt.Sprintf("\n\n[Agents used: %s]", strings.Join(calledAgents, ", "))
		finalContent = content + agentsText
	}

	task.Artifacts = append(task.Artifacts, &server.Artifact{
		Type:    "text",
		Content: finalContent,
	})

	task.Status = server.TaskState_TASK_STATE_COMPLETED
	return task, nil
}
