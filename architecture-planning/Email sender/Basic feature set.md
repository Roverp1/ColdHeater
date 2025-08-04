---
id: Basic feature set
aliases: []
tags: []
---

> [!WARNING]
> SHOULD BE ACVIEVED WITH 0-10$ budget

1. Singel sender account

2. **Email Warmup with High-Reputation Bot Networks** for single sender account:

   - A **small bot network** of high-reputation inboxes (bots) across different countries and industries to simulate human-like interactions (opens, replies, marking as important).
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

3. **Deliverability Optimization**:

   - **Email verification** to ensure valid recipient addresses and reduce bounce rates.
     - integration with email verification tools
   - **Spam word checker** to ensure content avoids spam filters.

4. **Campaign Automation and Scheduling**:

   - Each lead needs to be tracked in the database for a specific sender
     - specify a date when the last email was sent and on which category/topic (?)
     - specify categories/topics in which this lead can be interested in (?)

5. **Comprehensive Analytics**:

- basic account health monitoring

6. **Scalability and Flexibility**:

   - comprehensive customization of almost all configurable aspects

7. **Lead Generation and Management**:

   - Integration with tools for automatic lead search

8. **UI**:

   - Ability to launch app in CLI
   - potentially pretty interface

9. **AI**:
   - spin up local free model to avoid token cost
