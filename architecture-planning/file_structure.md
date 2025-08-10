---
id: file_structure
aliases: []
tags: []
---

<!-- prettier-ignore -->
ColdHeater/
├── main.go                 # Main CLI application
├── internal/                 # Private application code
│   ├── bot/                 # Bot Creation Module
│   │   ├── creator.go       # Bot creation logic
│   │   ├── manager.go       # Bot management
│   │   └── models.go        # Bot data structures
│   ├── sender/              # Sender Module
│   ├── campaign/            # Campaign Module
│   ├── queue/               # Queue Manager Module
│   ├── email/               # Email Operations Module
│   ├── database/            # Database Module
│   │   ├── models.go        # All data models
│   │   ├── migrations.go    # Schema setup
│   │   └── postgres.go      # PostgreSQL implementation
│   └── ui/                  # UI Module
│       ├── cli/
│       └── rest/
├── pkg/                     # Reusable libraries (if any)
├── configs/                 # Configuration files
├── scripts/                 # Build/deployment scripts
├── go.mod                   # Dependencies (like package.json)
├── go.sum                   # Lock file (like package-lock.json)
└── README.md

<!-- prettier-ignore-end -->

## Database Setup

- Using Supabase (PostgreSQL)
- Connection via environment variables
- Shared development database for team collaboration
