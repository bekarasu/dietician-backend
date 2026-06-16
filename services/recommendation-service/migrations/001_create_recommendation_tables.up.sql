CREATE TABLE IF NOT EXISTS recommendation_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    request_type VARCHAR(50) NOT NULL,
    provider VARCHAR(50) NOT NULL DEFAULT 'mock',
    input_data JSONB,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_recommendation_requests_user_id ON recommendation_requests (user_id);

CREATE TABLE IF NOT EXISTS recommendation_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_id UUID NOT NULL REFERENCES recommendation_requests(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    recommendation_type VARCHAR(50) NOT NULL,
    content JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_recommendation_results_user_id ON recommendation_results (user_id);
CREATE INDEX idx_recommendation_results_request_id ON recommendation_results (request_id);

CREATE TABLE IF NOT EXISTS recommendation_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    result_id UUID NOT NULL REFERENCES recommendation_results(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    rating INT,
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_recommendation_feedback_user_id ON recommendation_feedback (user_id);
