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
   - Warmup threading
   - Cold email sending
   - Reputation tracking
2. Bot Module
   <!-- - Rotation & limits -->
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
5. Bot Creation Module
   - Bot creation
   - Bot aging
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

## Database schema (postgresql)

- bots: id, email, ip, status, created_at, last_used, aging_end_date
- campaigns: id, ip
- senders: id, email, daily_quota (?), warmup_end_date, campaign_id
- sender_bot_interactions: id, sender_id, bot_id, interaction_type, created_at
- queue_tasks: id, campaign_id, queue_type, task_data, priority, scheduled_time, status
- conversations: id, sender_id, bot_id, thread_id, message_id, subject, next_step_time

## Complete Flow (cold emails)

_Before adding a sender, to a list of available senders. One need to select a proxy for a specific sender and configure SPF with selected proxy's ip_

1. Campaign Manager: User initializes new campaign (sets limits, uploads leads)
2. Campaign Manager: Checks available sender accounts after rotation, assigns them to campaign
3. Campaign Manager: Spawns Queue Module instances for each assigned sender (parallelization)
4. Queue Module: Checks if assigned sender needs warmup (sender.warmup_end_date > current_date)
5. Queue Module: If warmup needed → modifies parameters:

- daily_quota = reduced warmup limit
- percentageOfColdEmailsToBots = 100%
- If normal mode → uses original sender parameters

6. Queue Module: Initializes cold outreach queue with data loaded from database

7. Queue Module: Populates cold outreach queue with targets:

- Warmup mode: only bots (calculated amount based on daily_quota)
- Normal mode: leads + bots (based on percentageOfColdEmailsToBots)

7. Queue Module: Processes queue, with calculated delay
8. Sender Module: Sends email from sender to selected target
9. Database: Records interaction (sender_id, target_id, interaction_type, timestamp)
10. Bot Module: If target was bot → decides to reply based on configured rates → schedules reply in high-priority Reply Queue

## File structure

Initial file structure is described in [[file_structure]]
