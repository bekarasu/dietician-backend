package events

import "context"

// Event topic constants for domain events.
const (
	TopicUserRegistered          = "user.registered"
	TopicProfileUpdated          = "user.profile.updated"
	TopicWeightLogged            = "progress.weight.logged"
	TopicMealLogged              = "tracking.meal.logged"
	TopicHydrationLogged         = "tracking.hydration.logged"
	TopicMedicalFileUploaded     = "medical.file.uploaded"
	TopicRecommendationGenerated = "recommendation.generated"
	TopicWeeklyProgressEvaluated = "progress.weekly.evaluated"
)

// Event represents a domain event with a topic and payload.
type Event struct {
	Topic   string `json:"topic"`
	Payload any    `json:"payload"`
}

// EventPublisher is the interface for publishing domain events.
// Implementations can use Kafka, RabbitMQ, NATS, Redis Streams, etc.
type EventPublisher interface {
	Publish(ctx context.Context, topic string, payload any) error
}

// NoopPublisher is a no-op implementation of EventPublisher.
// It silently discards all events and is used as a placeholder
// until a real event bus (Kafka, NATS, etc.) is integrated.
type NoopPublisher struct{}

func NewNoopPublisher() *NoopPublisher {
	return &NoopPublisher{}
}

func (n *NoopPublisher) Publish(ctx context.Context, topic string, payload any) error {
	return nil
}

// --- Common event payloads ---

type UserRegisteredPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type ProfileUpdatedPayload struct {
	UserID string `json:"user_id"`
}

type WeightLoggedPayload struct {
	UserID   string  `json:"user_id"`
	WeightKg float64 `json:"weight_kg"`
}

type MealLoggedPayload struct {
	UserID   string `json:"user_id"`
	MealType string `json:"meal_type"`
	Calories *int   `json:"calories,omitempty"`
}

type HydrationLoggedPayload struct {
	UserID   string `json:"user_id"`
	AmountMl int    `json:"amount_ml"`
}

type MedicalFileUploadedPayload struct {
	UserID   string `json:"user_id"`
	UploadID string `json:"upload_id"`
}

type RecommendationGeneratedPayload struct {
	UserID    string `json:"user_id"`
	RequestID string `json:"request_id"`
	Type      string `json:"type"`
}

type WeeklyProgressEvaluatedPayload struct {
	UserID    string `json:"user_id"`
	WeekStart string `json:"week_start"`
}
