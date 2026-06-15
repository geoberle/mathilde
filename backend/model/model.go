package model

import "time"

// Profile lives at users/{uid}/profile/main. Backend-owned.
type Profile struct {
	Mission     string            `firestore:"mission,omitempty"`
	Preferences map[string]string `firestore:"preferences,omitempty"`
}

// Progress lives at users/{uid}/progress/main. Shell-owned.
type Progress struct {
	XP    int `firestore:"xp"`
	Level int `firestore:"level"`
}

// Concept is an atomic learning unit within a Topic.
type Concept struct {
	ID            string        `firestore:"id" json:"id"`
	Name          string        `firestore:"name" json:"name"`
	Prerequisites []string      `firestore:"prerequisites,omitempty" json:"prerequisites"`
	Status        ConceptStatus `firestore:"status,omitempty" json:"status,omitempty"`
}

// Topic lives at users/{uid}/topics/{id}.
type Topic struct {
	Name     string    `firestore:"name"`
	Concepts []Concept `firestore:"concepts"`
	AddedAt  time.Time `firestore:"addedAt"`
}

// SessionStatus tracks the lifecycle of a session.
type SessionStatus string

const (
	SessionPendingReview SessionStatus = "pending_review"
	SessionApproved      SessionStatus = "approved"
	SessionActive        SessionStatus = "active"
	SessionCompleted     SessionStatus = "completed"
	SessionRejected      SessionStatus = "rejected"
)

// SessionType defines what kind of session the planner generated.
type SessionType string

const (
	SessionLearn    SessionType = "learn"
	SessionDrill    SessionType = "drill"
	SessionQuiz     SessionType = "quiz"
	SessionReview   SessionType = "review"
	SessionExamPrep SessionType = "exam-prep"
)

// Session lives at users/{uid}/sessions/{id}.
type Session struct {
	TopicID    string        `firestore:"topicId"`
	ConceptID  string        `firestore:"conceptId,omitempty"`
	Type       SessionType   `firestore:"type"`
	Status     SessionStatus `firestore:"status"`
	HTML       string        `firestore:"html"`
	CreatedAt  time.Time     `firestore:"createdAt"`
	ApprovedAt time.Time     `firestore:"approvedAt,omitempty"`
	Order      int           `firestore:"order"`
}

// ExerciseResult captures the outcome of a single exercise.
type ExerciseResult struct {
	ExerciseID string `firestore:"exerciseId"`
	Correct    bool   `firestore:"correct"`
	HintsUsed  int    `firestore:"hintsUsed"`
	TimeSpent  int    `firestore:"timeSpent"`
}

// Result lives at users/{uid}/results/{id}.
type Result struct {
	SessionID   string           `firestore:"sessionId"`
	Exercises   []ExerciseResult `firestore:"exercises"`
	TotalXP     int              `firestore:"totalXp"`
	CompletedAt time.Time        `firestore:"completedAt"`
}

// ConceptStatus tracks mastery of a concept.
type ConceptStatus string

const (
	StatusNotStarted ConceptStatus = "not-started"
	StatusIntroduced ConceptStatus = "introduced"
	StatusStruggling ConceptStatus = "struggling"
	StatusMastered   ConceptStatus = "mastered"
)

// LearningRecord lives at users/{uid}/learning-records/{id}.
// Shape matches PRD exactly.
type LearningRecord struct {
	ConceptID     string        `firestore:"conceptId"`
	Status        ConceptStatus `firestore:"status"`
	Narrative     string        `firestore:"narrative"`
	LastPracticed time.Time     `firestore:"lastPracticed,omitempty"`
	NextReview    time.Time     `firestore:"nextReview,omitempty"`
	Interval      int           `firestore:"interval"`
	ErrorPatterns []string      `firestore:"errorPatterns,omitempty"`
}

// EvaluationStatus tracks the handwriting evaluation lifecycle.
type EvaluationStatus string

const (
	EvalPending   EvaluationStatus = "pending"
	EvalCompleted EvaluationStatus = "completed"
)

// Evaluation lives at users/{uid}/evaluations/{id}.
type Evaluation struct {
	ExerciseID string           `firestore:"exerciseId"`
	SessionID  string           `firestore:"sessionId"`
	ImageData  string           `firestore:"imageData,omitempty"`
	Status     EvaluationStatus `firestore:"status"`
	Result     string           `firestore:"result,omitempty"`
	CreatedAt  time.Time        `firestore:"createdAt"`
}

// Reference lives at users/{uid}/reference/{id}.
type Reference struct {
	Type    string `firestore:"type"`
	Title   string `firestore:"title"`
	Content string `firestore:"content"`
}
