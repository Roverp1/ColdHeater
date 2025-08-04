---
id: App-architecture
aliases: []
tags: []
---

## Features

Initial version of the app needs the following set of features [[Basic feature set]]. While being able to scale seamlessly into [[Full feature set]], without major rewrites.

## Development plan

Development should follow [[Version planning]]

## Main modules

1. Sender Module
   - Warmup scheduling
   - Cold email sending
   - Reputation tracking
   - Per-sender delays
2. Bot Module
   - Creation & aging
   - Rotation & limits
   - IP management
   - Reply delays
   - Reputation tracking
3. Campaign Module
   - Sender rotation
   - Campaign IP management
   - Lead management
4. Queue Manager Module:
   - Manages priority queues (reply vs cold outreach)
   - Enforces daily quotas and caps
   - Selects appropriate bots based on interaction history
5. Email Operations Module
   - Actual email sending/receiving (used by all above)
6. Database Module
   - Storage for all entities
7. UI Module
   - Status display

## Entities

1. Senders (real accounts that send cold emails)
   - warmup scheduling, cold email sending, reputation tracking, delay management per sender
2. Bots (fake accounts that build sender reputation)
   - creation, aging, rotation, limit management, ip management, per bot delays for replies, reputation tracking
3. Campaigns (cold outreach sequences)
   - sender rotation and management, ip management,

## Database schema

<!-- senders: id, daily_quota, current_count, reputation_score -->
<!-- bots: id, status, interaction_limits, cooldown_until   -->
<!-- bot_sender_interactions: sender_id, bot_id, interaction_type, last_contact -->
<!-- conversations: id, sender_id, bot_id, thread_status, next_action_time -->
<!-- scheduled_tasks: priority_level, task_type, target_time, payload -->
