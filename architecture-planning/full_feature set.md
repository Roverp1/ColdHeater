---
id: Full feature set
aliases: []
tags: []
---

1. **Multiple Sender Accounts with Parallel Sending**:

   - Support for **sender accounts** to rotate sending across multiple email accounts, preventing spam flags and scaling outreach volume.
   - Option to have a centralized inbox from where all the replies from the selected campaigns can be accessed.
   - **Smart sender rotation** to balance email load and mimic natural sending patterns.
   - Option to prirotize different days of weeks (e.g. only work days)

2. **Email Warmup with High-Reputation Bot Networks**:

   - A **large bot network** of high-reputation diverse inboxes (bots) across different countries and industries to simulate human-like interactions (opens, replies, marking as important).
   - Automatic/semi-automatic process of creating bots
   - Automatic process of aging bots. Simulate a normal human behavior to gain gmail reputation. It should lead some basic conversation and subscribe to some newspapers, etc...
   - Each bot should have it's own ip, from where it will perform it's functions
   - Each bot should have a unique browser fingeprint to not get detected by gmail
   - If email lends in spam, bot should automatically pull it out
   - smart delay system
     - sender should have a linear delay: send message, wait, send message, wait, send repy wait,
     - bots should have independent delays: 1st bot received message, wait, reply, 1st bot received message, wait, 2nd bot received message, wait, 1st bot reply, 2nd bot reply - each bot should have it's own timer.
   - **AI-driven warmup** that generates emails suitable for warmup.
   - automatic scheduling for warmup campaigns
     - each sender account need to be tracked to know if it was warmed up and for how long
     - gradual increasing of warmup emails on the course of the first 2 weeks of adding new sender email
   - each bot should have a time out, before it can be used again. each bot should have limits on conversation and interactions per day
   - probably single sender should have a linear delay so that cold outreach emails won't get send at the same time as warmup emails (don't know how solve yet, if warmup will be a separate module from cold outreach)
   - **Advanced Quota Management and Bot Interaction System**:
     - Hard daily cap on cold outreach emails per sender (20-25 per day maximum)
     - **Smart email distribution**: configurable percentage split between real leads and internal bots
       - `(100% - botEmailPercent)` sent to real prospects
       - `botEmailPercent` sent to internal bots for reputation building
     - **Bot Interaction Tracking**: comprehensive tracking of sender ↔ bot relationships
       - Track if each bot has had conversation with specific sender
       - Track interaction quality (replied, seen only, never interacted)
     - **Dynamic Bot Pool Management**:
       - **Active Warmup Pool**: Bots that have replied to sender between which configured amount of warmup emails is getting distributed
       - **Probabilistic Targeting**: Bots that only "saw" emails (didn't reply) still receive occasional cold emails but with lower probability
       - **Virgin Bot Pool**: Bots with no prior interaction for initial reputation building
     - **Automated Warmup Distribution**: calculated warmup emails distributed across bots with existing sender relationships
     - **Configurable Bot Reply Behavior**: percentage-based bot reply rates to simulate natural human interaction patterns

3. **AI-Powered Personalization**:
   - **Dynamic content generation** using spintax and AI to create personalized, engaging emails at scale.
     - Allow for customization of ai behavior, with global configuration and campaign specific.

- **Spam word checker** to ensure content avoids spam filters.

4. **Deliverability Optimization**:

   - automatic/help setting up **Domain authentication** (SPF, DKIM, DMARC), and verification before starting the campaign.
   - **Unique IP servers** per campaign to prevent cross-contamination of sender reputation.
   - **Email verification** to ensure valid recipient addresses and reduce bounce rates.
     - integration with email verification tools
   - Blacklist Monitoring: a blacklist check tool to monitor if domains or IPs are blacklisted.
   - Adaptive Sending: adaptive sending algorithms to adjust email volume based on ISP limits or recipient engagement

5. **Campaign Automation and Scheduling**:

   - **Automated workflows** for follow-ups, drip campaigns, and triggered responses based on recipient actions.
   - Each lead needs to be tracked in the database for a specific sender
     - specify a date when the last email was sent and on which category/topic (?)
     - specify categories/topics in which this lead can be interested in (?)
     - condition-based sequences (e.g., if a lead opens but doesn’t reply, send a specific follow-up).
   - Flexible **scheduling** to send emails at optimal times across time zones.
   - **A/B testing** inside a campaign for subject lines, content, and sending patterns to optimize performance.

6. **Comprehensive Analytics**:

- **Reputation monitoring** with alerts for drops in deliverability or inbox placement.
  - Real-time tracking of **open rates, reply rates, bounce rates, and inbox placement**.
  - Detailed **deliverability reports** across major providers (Gmail, Outlook, etc.).
- **Lead progress tracking** to monitor campaign ROI and engagement .

7. **Scalability and Flexibility**:

   - Support for **high-volume sending** (thousands of emails daily) without compromising deliverability.
   - **API integration** for connecting with CRMs, lead generation tools, or custom systems.
   - **Agency-friendly features** like branded client portals for managing multiple campaigns.
   - comprehensive customization of almost all configurable aspects
   - Having a one default configuration for all campaigns, and being able to separate campaigns into clusters which can have their own configuration, and possibility to overwrite all previous configurations with campaign specific configuration.
   - possibility to instantly see configuration changes without needing to restart/relaunch app/campaign
     - possibility to schedule configuration changes.

8. **Lead Generation and Management**:

   - Integration with tools for automatic lead search
   - **CRM capabilities** to manage leads, track interactions, and segment audiences.

9. **UI**:

   - Ability to launch app in CLI
   - potentially pretty interface
   - web interface
   - potentially with live updates

10. **Architecture**

- **Persistent State Management**: full conversation and task persistence across app restarts
- **Priority Queue System**:
  - High-priority Reply Queue for warmup and conversation responses
  - Lower-priority Cold Outreach Queue with daily quotas and caps
  - Queue state persisted to database for restart recovery
- **Centralized server** with possibility to connect via CLI or web interface
- **Dynamic parallelization** for different sender accounts/campaigns
- **Scalable concurrency**: as many concurrent sender accounts/campaigns on a single core → overflowing move onto separate cores
- **Modular monolith design** with well-defined interfaces between components
