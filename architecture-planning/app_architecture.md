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

## Complete Flow (cold emails)

1. Campaign Manager: User initializes new campaign (sets limits, uploads leads)
2. Campaign Manager: Checks available sender accounts after rotation, assigns them to campaign
3. Campaign Manager: Spawns Queue Module instances for each assigned sender (parallelization)
4. Queue Module: Checks if assigned sender needs warmup (sender.warmup_end_date > current_date)
5. Queue Module: If warmup needed → modifies parameters:

- daily_quota = reduced warmup limit
- percentageOfColdEmailsToBots = 100%
- If normal mode → uses original sender parameters

6. Queue Module: Populates cold outreach queue with targets:

- Warmup mode: only bots (calculated amount based on daily_quota)
- Normal mode: leads + bots (based on percentageOfColdEmailsToBots)

7. Queue Manager: Processes queue, selects specific bot using selection logic (virgin → probabilistic → active pool)
8. Email Operations: Sends email from sender to selected target
9. Database: Records interaction (sender_id, target_id, interaction_type, timestamp)
10. Bot Module: If target was bot → decides to reply based on configured rates → schedules reply in high-priority Reply Queue
