package openai

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
	goopenai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
)

//go:embed prompts/diet-recommendations.txt
var dietRecommendationsPrompt string

type Config struct {
	APIKey string
}

type Service interface {
	CreateDietRecommendations(ctx context.Context, req DietRecommendationRequestDto) (*DietRecommendationResponse, error)
}

type service struct {
	client    *goopenai.Client
	validator *validator.Validate
	logger    *logrus.Logger
}

func NewService(cfg Config, logger *logrus.Logger) Service {
	client := goopenai.NewClient(cfg.APIKey)
	return &service{
		client:    client,
		validator: validator.New(),
		logger:    logger,
	}
}

func (s *service) CreateDietRecommendations(ctx context.Context, req DietRecommendationRequestDto) (*DietRecommendationResponse, error) {
	s.logger.Infof("Generating diet recommendations for %d days", req.DurationDays)

	// Validate request
	if err := s.validator.Struct(req); err != nil {
		s.logger.Errorf("Validation error on request: %v", err)
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		s.logger.Errorf("Failed to marshal request payload: %v", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	s.logger.Debugf("Sending request to OpenAI with payload: %s", string(reqBytes))

	resp, err := s.client.CreateChatCompletion(ctx, goopenai.ChatCompletionRequest{
		Model: goopenai.GPT4oMini,
		Messages: []goopenai.ChatCompletionMessage{
			{
				Role:    goopenai.ChatMessageRoleSystem,
				Content: dietRecommendationsPrompt,
			},
			{
				Role:    goopenai.ChatMessageRoleUser,
				Content: string(reqBytes),
			},
		},
		ResponseFormat: &goopenai.ChatCompletionResponseFormat{
			Type: goopenai.ChatCompletionResponseFormatTypeJSONObject,
		},
	})

	if err != nil {
		s.logger.Errorf("OpenAI API error: %v", err)
		return nil, fmt.Errorf("OpenAI request failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		s.logger.Error("OpenAI returned no choices")
		return nil, fmt.Errorf("OpenAI returned no choices")
	}

	respContent := resp.Choices[0].Message.Content

	var dietResp DietRecommendationResponse
	if err := json.Unmarshal([]byte(respContent), &dietResp); err != nil {
		s.logger.Errorf("Failed to unmarshal OpenAI response: %v\nRaw response: %s", err, respContent)
		return nil, fmt.Errorf("invalid JSON response from OpenAI: %w", err)
	}

	// Validate the response structure
	if err := s.validator.Struct(dietResp); err != nil {
		s.logger.Errorf("OpenAI response validation failed: %v", err)
		return nil, fmt.Errorf("OpenAI response validation failed: %w", err)
	}

	s.logger.Info("Successfully generated and validated diet recommendations")

	return &dietResp, nil
}
