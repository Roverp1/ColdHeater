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

### 1. Time-Based Activity System

Each bot operates with configurable activity windows that define behavior patterns throughout the day:

```go
type ActivityWindow struct {
    StartHour         int           // 8
    EndHour           int           // 12
    EvaluationDelay   DelayRange // How often to "think" during this window
    ActionProbability float64       // 0.7 = 70% chance to act when evaluating
    MaxActionsPerHour int           // Rate limiting
    ActionTypes       map[string]float64 // "reply_email": 0.8, "browse_web": 0.3
}

type BotPersonality struct {
    // Time-Based Activity Management
    TimeZone            string            // "UTC", "EST", "PST", etc.
    ActiveWindows       []ActivityWindow  // High activity periods
    OffHoursWindow      *ActivityWindow   // Low activity for off-hours
    WeekendWindows      []ActivityWindow  // Different weekend patterns (optional)

    // Email Behavior
    ReplyChance         float64           // Base reply chance, modified by window probability
    EmailStyle          string            // "formal", "casual", "brief", "chatty"

    // Web Browsing
    BrowsingDuration    time.Duration     // 5-30 minutes per session
    FavoriteTopics      []string          // ["tech", "news", "sports", "finance"]

    // Deviation/Randomness
    DeviationChance     float64           // 0.05-0.15 (5-15% weird behavior)
}
```

### 2. Time-Window Based Decision Engine

#### Core Decision Flow

```
Bot Lifecycle: SLEEP → WAKE (if in activity window) → EVALUATE → DECIDE → ACT → MICRO-COOLDOWN → EVALUATE...
```

The bot continuously evaluates during active windows, rather than using fixed delays after actions.

#### Activity Window Logic

```go
func (b *Bot) ShouldEvaluate(currentTime time.Time) (bool, time.Duration) {
    window := b.GetCurrentActivityWindow(currentTime)

    if window == nil {
        // Outside active hours - use off-hours window if available
        if b.Personality.OffHoursWindow != nil {
            return true, b.Personality.OffHoursWindow.EvaluationDelay
        }
        return false, time.Hour // Sleep until next potential window
    }

    return true, window.EvaluationDelay
}

func (b *Bot) DecideAction(triggers []string) *Action {
    window := b.GetCurrentActivityWindow(time.Now())
    baseChance := window.ActionProbability

    // Apply randomness and personality factors
    if rand.Float64() > baseChance {
        return &Action{Type: "wait"}
    }

    // Select action type based on window configuration and triggers
    return b.selectActionFromWindow(window, triggers)
}
```

#### Decision Types

- **Email Response Decisions**: Reply, ignore, delay, mark important
- **Web Browsing Decisions**: What to browse, when to browse, how long
- **Navigation Decisions**: How to reach objectives despite UI changes
- **Wait Decisions**: Continue evaluation cycle without taking action

### 3. Window-Based Probabilistic Behavior System

The behavior system is now integrated with activity windows, removing fixed delays in favor of evaluation-based timing:

```go
type BehaviorRule struct {
    Trigger       string              // "email_received", "evaluation_cycle"
    ActionType    string              // "reply_email", "browse_web", "wait"
    Parameters    map[string]interface{}
    WindowModifiers map[string]float64 // Probability modifiers per window type (?)
    Conditions    []string            // ["not_spam", "business_hours"]
    Deviations    []Deviation
}

type Deviation struct {
    Action        Action
    Probability   float64
    Triggers      []string            // ["weekend", "off_hours", "random"]
}

type Action struct {
    Type          string              // "reply_email", "browse_web", "wait"
    Parameters    map[string]interface{}
    MicroCooldown time.Duration       // Brief pause after action (30s-5m) (?)
}
```

#### Example Behavior Rules for Time-Window System

```go
// Email response behavior adapts to current activity window
emailResponseRule := BehaviorRule{
    Trigger: "email_received",
    ActionType: "reply_email",
    Parameters: map[string]interface{}{
        "style": "match_personality",
        "length": "brief",
    },

    // its probably redundant
    // Different probabilities based on current window
    WindowModifiers: map[string]float64{
        "active_morning":   1.0,    // Full probability during active hours
        "active_afternoon": 0.8,    // Slightly lower in afternoon
        "off_hours":        0.2,    // Much lower during off hours
        "weekend":          0.6,    // Moderate weekend activity
    },

    Conditions: []string{"not_spam", "within_rate_limit"},

    Deviations: []Deviation{
        {
            Action: Action{
                Type: "ignore_email",
                MicroCooldown: 1*time.Minute,
            },
            Probability: 0.1,
            Triggers: []string{"random_deviation"},
        },
    },
}

// Browsing behavior triggered during evaluation cycles
browsingRule := BehaviorRule{
    Trigger: "evaluation_cycle",
    ActionType: "browse_web",
    Parameters: map[string]interface{}{
        "topic": "select_from_favorites",
        "duration": "5-15m",
    },

    WindowModifiers: map[string]float64{
        "active_morning":   0.3,    // Lower browsing during focused morning hours
        "active_afternoon": 0.6,    // Higher browsing in afternoon
        "off_hours":        0.1,    // Minimal off-hours browsing
        "weekend":          0.8,    // High weekend browsing
    },

    Conditions: []string{"no_pending_emails", "within_rate_limit"},
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

1. **Time-Window Based Activity**: Bots operate within configurable activity windows rather than fixed delays
2. **Fully Configurable**: All timing, probabilities, and behaviors configurable via database/config
3. **Probabilistic Over Deterministic**: Use weighted randomness instead of rigid rules
4. **Personality-Driven**: All decisions flow from bot personality traits and current activity window
5. **Realistic Human Patterns**: Off-hours activity possible but much less frequent
6. **Evaluation-Driven**: Continuous evaluation cycles during active periods replace action-delay patterns
7. **Failure Resilient**: Multiple fallback paths for web navigation
8. **Resource Conscious**: Fresh browsers per session, efficient memory usage
9. **Modular**: Easy to add new decision types and navigation goals
10. **Observable**: Comprehensive logging for debugging and optimization

## Benefits of Time-Window Architecture

- **Natural Response Times**: Response delays emerge from evaluation frequency + probability rather than forced waits
- **Realistic Behavior**: Humans check email periodically during work hours, not on rigid schedules
- **Flexible Off-Hours Activity**: Bots can still act during off-hours but with appropriate lower frequency
- **No Scheduling Conflicts**: Actions never spill outside appropriate time boundaries
- **Configurable Per Bot**: Each bot can have unique activity patterns and personalities
- **Rate Limiting**: Built-in action limits prevent unrealistic activity spikes

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
    timezone: "EST"
    reply_chance: 0.8
    email_style: "casual"
    deviation_chance: 0.1
    browsing_duration: "15m"

    # High activity periods
    active_windows:
      - start_hour: 8
        end_hour: 12
        evaluation_delay: "2m" # Check every 2 minutes
        action_probability: 0.8 # 80% chance to act when evaluating
        max_actions_per_hour: 15
        action_types:
          reply_email: 0.9 # High email priority during morning
          browse_web: 0.3
          check_gmail: 0.7

      - start_hour: 13
        end_hour: 17
        evaluation_delay: "3m" # Slightly less frequent afternoon checks
        action_probability: 0.6 # More selective in afternoon
        max_actions_per_hour: 10
        action_types:
          reply_email: 0.7
          browse_web: 0.5 # More browsing in afternoon
          check_gmail: 0.4

    # Low activity for off-hours
    off_hours_window:
      evaluation_delay: "30m" # Check every 30 minutes
      action_probability: 0.1 # Only 10% chance to act
      max_actions_per_hour: 2 # Very limited actions
      action_types:
        reply_email: 0.3 # Rare off-hours replies
        browse_web: 0.1
        check_gmail: 0.2

    # Different weekend behavior (optional)
    weekend_windows:
      - start_hour: 10
        end_hour: 14
        evaluation_delay: "10m" # Lazy weekend checking
        action_probability: 0.4
        max_actions_per_hour: 5
        action_types:
          reply_email: 0.5
          browse_web: 0.7 # More casual browsing on weekends
          check_gmail: 0.3

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
