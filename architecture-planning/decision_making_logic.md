---
id: decision_making_logic
aliases: []
tags: []
---

# Bot Decision-Making Architecture

## Overview

This document outlines the architecture for Gmail bot decision-making in the ColdHeater email automation system. Bots operate in two modes:

1. **Aging/Warmup Mode**: Simulate human behavior to gain reputation
2. **Cold Outreach Mode**: Maintain sender reputation by interacting with received emails

## Core Components

### 1. Bot Personality System

Each bot has a unique personality that drives all decision-making:

```go
type BotPersonality struct {
    // Email Behavior
    ResponseTimeMin     time.Duration // 30m - 4h range
    ResponseTimeMax     time.Duration
    ReplyChance         float64       // 0.7-0.9 (70-90% reply rate)
    EmailStyle          string        // "formal", "casual", "brief", "chatty"

    // Activity Patterns
    ActiveHoursStart    int           // 6-10 (AM)
    ActiveHoursEnd      int           // 16-22 (PM)
    TimeZone            string        // "UTC", "EST", "PST", etc.

    // Web Browsing
    BrowsingFrequency   time.Duration // 2-8 hours between sessions
    BrowsingDuration    time.Duration // 5-30 minutes per session
    FavoriteTopics      []string      // ["tech", "news", "sports", "finance"]

    // Deviation/Randomness
    DeviationChance     float64       // 0.05-0.15 (5-15% weird behavior)
    WeekendBehavior     bool          // true = different weekend patterns
}
```

### 2. Decision Engine Architecture

#### Core Decision Flow

```
Bot Lifecycle: IDLE → EVALUATE → DECIDE → ACT → COOLDOWN → IDLE
```

#### Decision Types

- **Email Response Decisions**: Reply, ignore, delay, mark important
- **Web Browsing Decisions**: What to browse, when to browse, how long
- **Navigation Decisions**: How to reach objectives despite UI changes

### 3. Probabilistic Behavior System

```go
type BehaviorRule struct {
    Trigger       string    // "email_received", "time_to_browse"
    BaseAction    Action    // Primary behavior
    Probability   float64   // 0.8 = 80% chance
    Deviations    []Deviation
    Conditions    []string  // ["active_hours", "conversation_ongoing"]
}

type Deviation struct {
    Action        Action
    Probability   float64
    Triggers      []string  // ["weekend", "late_night", "random"]
}

type Action struct {
    Type          string    // "reply_email", "browse_web", "wait"
    Parameters    map[string]interface{}
    Delay         DelayRange
}
```

#### Example Behavior Rules

```go
emailResponseRule := BehaviorRule{
    Trigger: "email_received",
    BaseAction: Action{
        Type: "reply_email",
        Parameters: map[string]interface{}{
            "style": "match_personality",
            "length": "brief",
        },
        Delay: DelayRange{Min: 30*time.Minute, Max: 2*time.Hour},
    },
    Probability: 0.8,
    Deviations: []Deviation{
        {
            Action: Action{Type: "ignore_email"},
            Probability: 0.1,
            Triggers: []string{"random"},
        },
        {
            Action: Action{Type: "reply_email", Delay: DelayRange{Min: 4*time.Hour, Max: 12*time.Hour}},
            Probability: 0.1,
            Triggers: []string{"late_night", "weekend"},
        },
    },
    Conditions: []string{"active_hours", "not_spam"},
}
```

### 4. Goal-Oriented Navigation System

To handle unstable web UIs, implement flexible navigation:

```go
type NavigationGoal struct {
    Objective     string              // "send_email", "browse_news", "check_gmail"
    Priority      int                 // 1-5, higher = more important
    Steps         []NavigationStep    // Primary path
    Fallbacks     []NavigationStep    // Alternative paths
    MaxRetries    int                 // Give up after N failures
    TimeLimit     time.Duration       // Abort if taking too long
}

type NavigationStep struct {
    Action        string              // "click", "type", "wait", "scroll"
    Target        string              // Logical target name
    Selectors     []ElementSelector   // Multiple ways to find element
    SuccessCriteria string            // How to know it worked
    MaxRetries    int
    Fallback      *NavigationStep     // What to do if this fails
}

type ElementSelector struct {
    Type          string              // "css", "xpath", "text", "aria-label"
    Value         string              // The actual selector
    Priority      int                 // Try higher priority first
}
```

#### Example Navigation Goals

```go
composeEmailGoal := NavigationGoal{
    Objective: "compose_email",
    Steps: []NavigationStep{
        {
            Action: "click",
            Target: "compose_button",
            Selectors: []ElementSelector{
                {Type: "aria-label", Value: "Compose", Priority: 1},
                {Type: "css", Value: "[data-tooltip='Compose']", Priority: 2},
                {Type: "text", Value: "Compose", Priority: 3},
            },
        },
        // ... more steps
    },
}
```

### 5. Local LLM Integration

#### Recommended Setup

- **Model**: TinyLlama 1.1B or Phi-3 Mini 3.8B
- **Memory**: 2-4GB VRAM usage
- **Use Cases**: Email content generation only (initially)

#### LLM Service Interface

```go
type LLMService interface {
    GenerateEmailReply(context EmailContext, style string) (string, error)
    GenerateEmailSubject(content string, style string) (string, error)
    IsEmailSpam(content string) (bool, float64, error)
}

type EmailContext struct {
    ReceivedEmail   string
    ConversationHistory []string
    BotPersonality  BotPersonality
    Relationship    string // "unknown", "frequent", "business"
}
```

### 6. Resource Management

#### Browser Lifecycle

- **Strategy**: Create fresh browser instance per session
- **Benefits**: New fingerprint, clean cookies, better resource management
- **Implementation**: Pool of browser configurations, random selection

```go
type BrowserManager struct {
    ConfigPool    []BrowserConfig
    ActiveSessions map[string]*BrowserSession
    MaxConcurrent int
}

type BrowserConfig struct {
    UserAgent     string
    ViewportSize  Dimensions
    Fonts         []string
    Languages     []string
    Timezone      string
    Fingerprint   map[string]interface{}
}
```

#### Bot Session Management

```go
type BotSession struct {
    Bot           Bot
    Personality   BotPersonality
    CurrentGoal   *NavigationGoal
    Browser       *rod.Browser
    State         BotState
    LastActivity  time.Time
    SessionStart  time.Time
    MaxDuration   time.Duration
}

type BotState string
const (
    StateIdle      BotState = "idle"
    StateReading   BotState = "reading_email"
    StateBrowsing  BotState = "browsing_web"
    StateReplying  BotState = "replying_email"
    StateCooldown  BotState = "cooldown"
)
```

## Implementation Phases

### Phase 1: Core Decision Engine

1. Implement BotPersonality system
2. Create basic Decision Engine with probabilistic behavior
3. Add simple email response decisions
4. Basic LLM integration for content generation

### Phase 2: Web Navigation

1. Implement Goal-Oriented Navigation system
2. Create browser fingerprint management
3. Add web browsing goals (news, newsletters, YouTube)
4. Resource pooling and session management

### Phase 3: Advanced Features

1. Conversation context awareness
2. Adaptive behavior based on success/failure rates
3. Enhanced anti-detection measures
4. Performance optimization

## Key Design Principles

1. **Probabilistic Over Deterministic**: Use weighted randomness instead of rigid rules
2. **Personality-Driven**: All decisions flow from bot personality traits
3. **Failure Resilient**: Multiple fallback paths for web navigation
4. **Resource Conscious**: Fresh browsers per session, efficient memory usage
5. **Modular**: Easy to add new decision types and navigation goals
6. **Observable**: Comprehensive logging for debugging and optimization

## Database Schema Extensions

```sql
-- Bot personalities
ALTER TABLE bots ADD COLUMN personality_config JSONB;

-- Decision logs for debugging
CREATE TABLE bot_decisions (
    id SERIAL PRIMARY KEY,
    bot_email VARCHAR(255) REFERENCES bots(email),
    decision_type VARCHAR(50),
    trigger_event VARCHAR(100),
    action_taken VARCHAR(100),
    success BOOLEAN,
    execution_time_ms INTEGER,
    error_message TEXT,
    context JSONB,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Navigation success tracking
CREATE TABLE navigation_metrics (
    id SERIAL PRIMARY KEY,
    goal_objective VARCHAR(100),
    success_rate FLOAT,
    avg_execution_time_ms INTEGER,
    last_updated TIMESTAMP DEFAULT NOW()
);
```

## Configuration Example

```yaml
bot_decision_engine:
  default_personality:
    response_time_min: "30m"
    response_time_max: "4h"
    reply_chance: 0.8
    email_style: "casual"
    browsing_frequency: "4h"
    deviation_chance: 0.1

  llm:
    model: "tinyllama-1.1b"
    max_tokens: 150
    temperature: 0.7

  browser:
    max_concurrent: 5
    session_timeout: "30m"
    fingerprint_rotation: true

  navigation_goals:
    - name: "check_gmail"
      timeout: "2m"
      retry_limit: 3
    - name: "browse_news"
      timeout: "5m"
      retry_limit: 2
```

This architecture provides a solid foundation for realistic bot behavior while remaining maintainable and extensible.

