# Bug Bounty Continuous Reconnaissance Platform
## Vision & Technical Design Document

**Version:** 1.0  
**Date:** October 2025  
**Status:** Discovery Phase  
**Author:** Platform Architect  

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Vision & Goals](#vision--goals)
3. [Problem Statement](#problem-statement)
4. [Solution Overview](#solution-overview)
5. [System Architecture](#system-architecture)
6. [Technology Stack](#technology-stack)
7. [Core Components](#core-components)
8. [API Design](#api-design)
9. [Data Model](#data-model)
10. [Deployment Strategy](#deployment-strategy)
11. [Development Roadmap](#development-roadmap)
12. [Cost Analysis](#cost-analysis)
13. [Security & Legal](#security--legal)
14. [Success Metrics](#success-metrics)
15. [Appendix](#appendix)

---

## Executive Summary

This document outlines the design and implementation plan for a **Continuous Reconnaissance and Anomaly Detection Platform** purpose-built for bug bounty hunting. The platform provides 24/7 automated monitoring of target assets, intelligent change detection, temporal pattern analysis, and automated alerting to identify security vulnerabilities before other researchers.

### Key Capabilities

- **Continuous Monitoring**: 24/7 surveillance of bug bounty program assets
- **Behavioral Baselining**: Learn normal deployment patterns per program
- **Anomaly Detection**: Identify suspicious changes using Bayesian scoring
- **Temporal Analysis**: Detect out-of-pattern deployments (weekend deploys, panic fixes)
- **Intelligent Alerting**: Priority-based notifications via Discord/Slack
- **Professional Infrastructure**: Kubernetes-orchestrated, Terraform-provisioned
- **Modern Architecture**: Go microservices with hybrid REST/gRPC API

### Strategic Advantage

While other bug bounty hunters manually check programs during business hours, this platform operates continuously, detecting that 2 AM production deployment or weekend emergency fix that often introduces security vulnerabilities. The system learns each program's normal behavior and alerts when deviations occur, providing a significant competitive advantage in the bug bounty landscape.

### Timeline & Investment

- **MVP Delivery**: 2-3 weeks
- **Initial Cost**: $0 (using cloud free trials)
- **Operational Cost**: $7-20/month (single VPS)
- **Technology Stack**: Go, Kubernetes, TimescaleDB, Redis

---

## Vision & Goals

### Primary Vision

Build an intelligent reconnaissance platform that acts as a **force multiplier** for bug bounty hunting by automating the tedious, time-consuming aspects of asset discovery and change monitoring, while surfacing high-value opportunities through machine learning and behavioral analysis.

### Strategic Goals

1. **Competitive Advantage**: Discover vulnerabilities before other researchers through continuous monitoring
2. **Efficiency**: Reduce manual reconnaissance time from hours daily to minutes of reviewing prioritized alerts
3. **Intelligence**: Learn and adapt to each program's unique patterns and behaviors
4. **Scale**: Monitor 10+ programs simultaneously without additional manual effort
5. **Professionalism**: Deploy using industry-standard DevOps practices (k8s, Terraform, IaC)

### Learning Objectives

As a secondary benefit, this project provides hands-on experience with:

- Go microservices development
- Kubernetes orchestration (k3s)
- Infrastructure as Code (Terraform)
- gRPC and REST API design
- Time-series database optimization (TimescaleDB)
- Real-time streaming architectures
- Container orchestration and deployment

---

## Problem Statement

### Current Bug Bounty Challenges

**Manual Reconnaissance is Time-Consuming**
- Researchers spend 60-80% of time on reconnaissance
- Must manually check each program daily for changes
- Easy to miss new assets or deployments
- No systematic way to detect patterns

**Critical Timing Windows are Missed**
- New deployments often happen outside business hours
- Emergency fixes (highest bug probability) occur on weekends/nights
- By the time researchers check, other hunters have already found issues
- First reporter gets the bounty - timing is critical

**Change Detection is Primitive**
- Most hunters take snapshots and manually diff
- No historical analysis or pattern recognition
- Cannot distinguish "normal" changes from suspicious ones
- Alert fatigue from too many low-value notifications

**No Intelligence Layer**
- Cannot learn which types of changes lead to vulnerabilities
- No risk scoring or prioritization
- Every change treated equally
- Researchers waste time investigating false positives

### Quantified Impact

- **Time Lost**: 10-20 hours/week on manual recon per researcher
- **Opportunities Missed**: Estimated 30-40% of vulnerabilities found by faster researchers
- **Inefficiency**: 80% of investigated changes yield no vulnerabilities
- **Competitive Disadvantage**: First-mover advantage in bug bounties is critical

---

## Solution Overview

### High-Level Approach

Deploy a **distributed, intelligent reconnaissance system** that:

1. **Continuously monitors** all in-scope assets for target programs
2. **Establishes behavioral baselines** for each program (deployment patterns, asset lifecycles)
3. **Detects anomalies** using statistical analysis and machine learning
4. **Scores and prioritizes** changes based on vulnerability probability
5. **Alerts in real-time** via preferred channels (Discord, Slack, CLI)
6. **Learns over time** from feedback to improve accuracy

### Core Innovation: Temporal Anomaly Detection

The platform's key differentiator is **temporal pattern analysis**:

- Learns when each program typically deploys (e.g., Tuesday/Thursday 2-4 AM)
- Flags out-of-pattern activity (weekend deploy = potential emergency fix)
- Correlates multiple signals (new subdomain + weekend + dev environment = HIGH priority)
- Applies Bayesian reasoning to calculate probability of exploitable bugs

### Example Scenario

```
Friday, 11:47 PM - Platform detects new subdomain: dev-api.example.com

Behavioral Analysis:
✓ Outside normal deployment window (Tuesday 2-4 AM)
✓ Weekend deployment (5x risk multiplier)
✓ Subdomain naming pattern: "dev-*" (3x risk multiplier)
✓ No previous deployments this time of day in 6 months
✓ Technology stack: nginx/1.18.0 (new version for this org)

Bayesian Score: 87.3/100 (HIGH PRIORITY)

Alert sent via Discord: 
"🚨 HIGH: New dev subdomain detected during weekend deploy
 Probability of exploitable bugs: 87%
 Investigate immediately: dev-api.example.com"

Researcher investigates within 15 minutes, discovers authentication bypass,
reports vulnerability before other researchers even wake up.

Bounty: $2,500
```

### Competitive Positioning

| Capability | Manual Recon | Basic Automation | This Platform |
|------------|--------------|------------------|---------------|
| Continuous Monitoring | ❌ | ✅ | ✅ |
| Pattern Learning | ❌ | ❌ | ✅ |
| Temporal Analysis | ❌ | ❌ | ✅ |
| Priority Scoring | ❌ | Basic | Advanced (Bayesian) |
| Real-time Alerts | ❌ | ✅ | ✅ (Streaming) |
| Historical Analysis | ❌ | ❌ | ✅ |
| Multi-program Scale | ❌ | Limited | Excellent |

---

## System Architecture

### Architectural Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     Kubernetes Cluster (k3s)                     │
│                                                                   │
│  ┌────────────────────────────────────────────────────────┐    │
│  │                   API Layer                              │    │
│  │  ┌──────────────────┐      ┌─────────────────────┐     │    │
│  │  │   REST API       │      │   gRPC Server        │     │    │
│  │  │   (port 8080)    │      │   (port 9090)        │     │    │
│  │  │                  │      │                      │     │    │
│  │  │ • Programs CRUD  │      │ • StreamAnomalies   │     │    │
│  │  │ • Scans CRUD     │      │ • WatchScan         │     │    │
│  │  │ • Anomalies GET  │      │ • StreamStats       │     │    │
│  │  │ • Auth           │      │ • StreamLogs        │     │    │
│  │  └────────┬─────────┘      └────────┬────────────┘     │    │
│  └───────────┼──────────────────────────┼──────────────────┘    │
│              │                           │                        │
│  ┌───────────┴───────────────────────────┴──────────────────┐   │
│  │                  Business Logic Layer                      │   │
│  │  ┌──────────┐  ┌──────────┐  ┌───────────┐  ┌─────────┐ │   │
│  │  │ Scanner  │  │   Diff    │  │  Scoring  │  │ Alerting│ │   │
│  │  │ Service  │  │  Engine   │  │  Engine   │  │ Service │ │   │
│  │  └──────────┘  └──────────┘  └───────────┘  └─────────┘ │   │
│  └────────────────────────────────┬─────────────────────────┘   │
│                                    │                              │
│  ┌────────────────────────────────┴─────────────────────────┐   │
│  │                    Worker Pool                            │   │
│  │  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐    │   │
│  │  │Worker 1 │  │Worker 2 │  │Worker 3 │  │Worker N │    │   │
│  │  │subfinder│  │  httpx  │  │  nuclei │  │  amass  │    │   │
│  │  └─────────┘  └─────────┘  └─────────┘  └─────────┘    │   │
│  └────────────────────────────────┬─────────────────────────┘   │
│                                    │                              │
│  ┌────────────────────────────────┴─────────────────────────┐   │
│  │               Data & Queue Layer                          │   │
│  │  ┌──────────────────┐         ┌─────────────────────┐   │   │
│  │  │  TimescaleDB     │         │      Redis          │   │   │
│  │  │  (StatefulSet)   │         │   (Deployment)      │   │   │
│  │  │                  │         │                     │   │   │
│  │  │ • Programs       │         │ • Task Queue        │   │   │
│  │  │ • Assets         │         │ • Dedup Cache       │   │   │
│  │  │ • Anomalies      │         │ • Temp Results      │   │   │
│  │  │ • Findings       │         │                     │   │   │
│  │  └──────────────────┘         └─────────────────────┘   │   │
│  └────────────────────────────────────────────────────────┘   │
│                                                                   │
│  ┌────────────────────────────────────────────────────────┐    │
│  │              Scheduled Jobs (CronJobs)                  │    │
│  │  • Passive Recon (hourly)                               │    │
│  │  • Active Scans (daily)                                 │    │
│  │  • Pattern Analysis (daily)                             │    │
│  │  • Database Cleanup (weekly)                            │    │
│  └────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
                               ▲
                               │
                    REST + gRPC over Internet
                               │
┌──────────────────────────────┴──────────────────────────────────┐
│                    CLI Client (Local Machine)                    │
│  ┌────────────────────────────────────────────────────────┐    │
│  │                 recon-cli (Cobra)                       │    │
│  │                                                          │    │
│  │  REST Commands:                                         │    │
│  │  • program add/list/delete                              │    │
│  │  • scan trigger/list                                    │    │
│  │  • anomalies query                                      │    │
│  │                                                          │    │
│  │  gRPC Streaming:                                        │    │
│  │  • dashboard (live view)                                │    │
│  │  • scan watch (progress)                                │    │
│  │  • stats live (metrics)                                 │    │
│  └────────────────────────────────────────────────────────┘    │
└─────────────────────────────────────────────────────────────────┘
```

### Data Flow: Asset Discovery & Alerting

```
1. CronJob triggers → Scan Job created in Redis
2. Worker picks up job → Runs subfinder/httpx
3. Worker stores results → TimescaleDB
4. Diff Engine compares → Current vs. Historical
5. Changes detected → Sent to Scoring Engine
6. Scoring Engine evaluates → Bayesian probability calculation
7. High-priority anomaly → Alert Service notified
8. Alert Service sends → Discord/Slack webhook
9. gRPC stream pushes → Live dashboard updates
10. CLI displays → Real-time notification
```

### Component Interaction: Temporal Pattern Learning

```
Daily at 00:00 UTC:
1. Pattern Analysis Job runs
2. Query all changes from past 24 hours
3. Group by program_id, day_of_week, hour
4. Update deployment_patterns table
5. Calculate statistical norms per program
6. Identify outliers (anomalous patterns)
7. Update Bayesian priors for scoring engine
```

---

## Technology Stack

### Language & Runtime

**Go (Golang) 1.25+**

**Rationale:**
- Excellent concurrency model (goroutines) for parallel scanning
- Single binary deployment (no dependency management)
- Fast compilation and execution
- Strong typing reduces bugs
- Native Kubernetes support
- Industry standard for cloud-native applications
- Great for CLI tools (Cobra ecosystem)
- Many recon tools already written in Go (subfinder, httpx, nuclei)

### API Layer

**Hybrid: REST + gRPC**

**REST API (Chi/Fiber/Echo)**
- CRUD operations (programs, scans, anomalies)
- Authentication and authorization
- Standard HTTP/1.1
- Easy debugging with curl/Postman
- Future web UI compatibility

**gRPC (Protocol Buffers 3)**
- Real-time streaming (dashboard, scan progress)
- Efficient binary protocol (HTTP/2)
- Type-safe contracts
- Bidirectional communication
- Low latency for live updates

**Rationale for Hybrid:**
- REST for simplicity where streaming isn't needed
- gRPC for performance-critical streaming features
- Best of both worlds
- Learn both modern API patterns

### Database

**TimescaleDB (PostgreSQL Extension)**

**Rationale:**
- Built specifically for time-series data (our exact use case)
- 1000x faster queries vs standard PostgreSQL
- 90% data compression (storage savings)
- Automatic data partitioning (hypertables)
- Built-in retention policies
- Continuous aggregates (pre-computed stats)
- SQL interface (no new query language)
- Open source, production-ready

**Alternative Considered:**
- **PostgreSQL**: Good, but lacks time-series optimizations
- **InfluxDB**: Too specialized, less flexible for relational queries
- **MongoDB**: Wrong data model for our use case

### Task Queue

**Asynq or River**

**Asynq (Redis-based)**
- Distributed task queue for Go
- Built on Redis
- Web UI for monitoring
- Retry and timeout handling
- Scheduled and periodic tasks
- Similar to Celery but Go-native

**River (PostgreSQL-based)**
- Job queue stored directly in Postgres
- One less dependency (no Redis needed)
- Simpler architecture
- Good observability

**Recommendation:** Start with **River** for simplicity (one database), migrate to **Asynq** if Redis needed for caching.

### Orchestration

**Kubernetes (k3s)**

**Rationale:**
- Industry-standard container orchestration
- Declarative infrastructure
- Horizontal scaling capability
- Self-healing and health checks
- Production-ready deployment practices
- Excellent learning value

**k3s Specifically:**
- Lightweight Kubernetes (512MB RAM vs 4GB+)
- Perfect for single-node or small clusters
- Runs on VPS efficiently
- Production-ready (used by many companies)
- Same k8s API as full distribution

### Infrastructure as Code

**Terraform**

**Rationale:**
- Industry standard for IaC
- Provider for all major clouds
- Declarative configuration
- State management
- Plan/apply workflow
- Reproducible infrastructure

### CLI Framework

**Cobra + Viper**

**Rationale:**
- Industry standard (kubectl, Hugo, GitHub CLI use it)
- Excellent command structure
- Flag and config management (Viper)
- Auto-generated documentation
- Shell completion support

### Containerization

**Docker**

**Rationale:**
- Standard containerization platform
- Multi-stage builds for small images
- Great Go support
- Easy local development

### Monitoring & Observability (Future)

**Prometheus + Grafana**
- Time-series metrics
- Kubernetes-native
- Industry standard

---

## Core Components

### 1. REST API Server

**Purpose:** Handle CRUD operations and authentication

**Responsibilities:**
- Program management (add, list, update, delete)
- Scan triggering and querying
- Anomaly retrieval and review
- User authentication
- API key management

**Technology:** Go + Chi/Fiber router

**Endpoints:**
```
POST   /api/v1/programs
GET    /api/v1/programs
GET    /api/v1/programs/:id
DELETE /api/v1/programs/:id

POST   /api/v1/scans
GET    /api/v1/scans
GET    /api/v1/scans/:id

GET    /api/v1/anomalies
GET    /api/v1/anomalies/:id
PATCH  /api/v1/anomalies/:id
```

### 2. gRPC Streaming Server

**Purpose:** Provide real-time data streams to CLI

**Responsibilities:**
- Stream new anomalies as detected
- Live scan progress updates
- Real-time statistics
- Log streaming (optional)

**Technology:** Go + gRPC + Protocol Buffers

**Services:**
```protobuf
service StreamService {
  rpc StreamAnomalies(StreamRequest) returns (stream Anomaly);
  rpc WatchScan(WatchScanRequest) returns (stream ScanProgress);
  rpc StreamStats(StatsRequest) returns (stream Stats);
  rpc StreamLogs(LogRequest) returns (stream LogEntry);
}
```

### 3. Worker Pool

**Purpose:** Execute reconnaissance tools and store results

**Responsibilities:**
- Pull jobs from queue
- Execute recon tools (subfinder, httpx, amass, nuclei, etc.)
- Parse and normalize output
- Store results in database
- Handle failures and retries

**Technology:** Go + Asynq/River workers

**Worker Types:**
- **Passive recon workers**: DNS enumeration, certificate transparency
- **Active scan workers**: HTTP probing, port scanning
- **Deep analysis workers**: Technology detection, content discovery

**Scaling:** Horizontal scaling via Kubernetes Deployment (replicas)

### 4. Diff Engine

**Purpose:** Compare current state to historical baseline

**Responsibilities:**
- Query current and previous scan results
- Calculate differences (new assets, changed assets, removed assets)
- Generate change events
- Fingerprint assets (content hashing)
- Detect technology stack changes

**Technology:** Go service, runs periodically or on-demand

**Algorithms:**
- Hash-based comparison (SHA256 of responses)
- Text diff for content changes
- JSON schema diff for API changes
- Time-based queries (TimescaleDB optimized)

### 5. Scoring Engine

**Purpose:** Prioritize anomalies using Bayesian probability

**Responsibilities:**
- Receive detected changes from Diff Engine
- Calculate base probability (prior)
- Apply evidence factors (weekend deploy, dev env, etc.)
- Compute posterior probability
- Assign priority score
- Learn from feedback over time

**Technology:** Go service with statistical library

**Scoring Algorithm:**
```
posterior = prior × ∏(evidence_factors)

Evidence Factors:
- is_weekend_deploy: 5.0x
- is_dev_subdomain: 3.0x
- outside_normal_window: 2.0x
- tech_stack_change: 1.5x
- rapid_deployments: 2.5x
- new_asset: 4.0x
- cert_change: 1.8x

priority_score = posterior × impact_multiplier
```

### 6. Alert Service

**Purpose:** Send notifications to configured channels

**Responsibilities:**
- Receive high-priority anomalies
- Format alert messages
- Send to Discord/Slack/Email
- Manage alert batching (prevent spam)
- Track alert history

**Technology:** Go service with webhook clients

**Channels:**
- Discord webhook
- Slack webhook
- Email (SMTP)
- CLI real-time stream (gRPC)

### 7. Pattern Analysis Service

**Purpose:** Learn deployment behaviors per program

**Responsibilities:**
- Aggregate historical change data
- Calculate deployment windows
- Identify statistical norms
- Detect anomalous patterns
- Update Bayesian priors

**Technology:** Go service, runs as daily CronJob

**Analysis:**
- Group changes by day_of_week, hour_of_day
- Calculate mean, standard deviation, outliers
- Store in deployment_patterns table
- Update scoring engine configuration

### 8. CLI Client

**Purpose:** Primary user interface for platform management

**Responsibilities:**
- Program management commands
- Scan triggering
- Anomaly querying
- Live dashboard (TUI)
- Configuration management

**Technology:** Go + Cobra + Bubble Tea (TUI)

**Commands:**
```bash
recon-cli program add
recon-cli program list
recon-cli scan trigger
recon-cli scan watch
recon-cli anomalies list
recon-cli dashboard
recon-cli stats
```

---

## API Design

### REST API Specification

#### Authentication

**Method:** API Key via Header

```
Authorization: Bearer <api-key>
```

**Key Management:**
- Generated on first deployment
- Stored in Kubernetes Secret
- Configured in CLI config file

#### Program Management

**Add Program**
```http
POST /api/v1/programs
Content-Type: application/json

{
  "name": "Example Corp",
  "platform": "hackerone",
  "scope": ["*.example.com", "*.example.io"],
  "scan_frequency": "1h"
}

Response: 201 Created
{
  "id": 1,
  "name": "Example Corp",
  "platform": "hackerone",
  "scope": ["*.example.com", "*.example.io"],
  "scan_frequency": "1h",
  "created_at": "2025-10-01T10:00:00Z"
}
```

**List Programs**
```http
GET /api/v1/programs

Response: 200 OK
{
  "programs": [
    {
      "id": 1,
      "name": "Example Corp",
      "platform": "hackerone",
      "scope": ["*.example.com"],
      "last_scanned_at": "2025-10-01T09:45:00Z",
      "asset_count": 147,
      "anomaly_count": 3
    }
  ],
  "total": 1
}
```

#### Scan Management

**Trigger Scan**
```http
POST /api/v1/scans
Content-Type: application/json

{
  "program_id": 1,
  "scan_type": "passive"
}

Response: 202 Accepted
{
  "scan_id": 42,
  "status": "queued",
  "created_at": "2025-10-01T10:05:00Z"
}
```

**Get Scan Status**
```http
GET /api/v1/scans/42

Response: 200 OK
{
  "id": 42,
  "program_id": 1,
  "scan_type": "passive",
  "status": "running",
  "progress": 65,
  "assets_found": 23,
  "started_at": "2025-10-01T10:05:30Z"
}
```

#### Anomaly Queries

**List Anomalies**
```http
GET /api/v1/anomalies?program_id=1&min_priority=70&unreviewed=true

Response: 200 OK
{
  "anomalies": [
    {
      "id": 12,
      "program_id": 1,
      "type": "new_subdomain",
      "description": "New subdomain discovered: dev-api.example.com",
      "priority_score": 87.3,
      "detected_at": "2025-10-01T02:47:15Z",
      "is_reviewed": false,
      "metadata": {
        "subdomain": "dev-api.example.com",
        "status_code": 200,
        "tech_stack": ["nginx/1.18.0"]
      }
    }
  ],
  "total": 1
}
```

### gRPC Service Specification

**Protocol Buffers Definition**

```protobuf
syntax = "proto3";
package recon.v1;
option go_package = "github.com/you/recon-platform/gen/go/recon/v1";

// Streaming service for real-time updates
service StreamService {
  // Stream anomalies as they're detected
  rpc StreamAnomalies(StreamAnomaliesRequest) returns (stream Anomaly);
  
  // Watch scan progress in real-time
  rpc WatchScan(WatchScanRequest) returns (stream ScanProgress);
  
  // Stream live statistics
  rpc StreamStats(StreamStatsRequest) returns (stream Stats);
  
  // Stream logs (optional, for debugging)
  rpc StreamLogs(StreamLogsRequest) returns (stream LogEntry);
}

message StreamAnomaliesRequest {
  repeated int64 program_ids = 1;
  float min_priority = 2;
  bool unreviewed_only = 3;
}

message Anomaly {
  int64 id = 1;
  int64 program_id = 2;
  string program_name = 3;
  string type = 4;
  string description = 5;
  float priority_score = 6;
  int64 detected_at = 7;
  map<string, string> metadata = 8;
}

message WatchScanRequest {
  int64 scan_id = 1;
}

message ScanProgress {
  int64 scan_id = 1;
  string status = 2;
  int32 progress = 3;
  string current_step = 4;
  int32 assets_found = 5;
  int64 timestamp = 6;
}

message StreamStatsRequest {
  int32 refresh_interval = 1;
}

message Stats {
  int32 active_programs = 1;
  int32 total_assets = 2;
  int32 unreviewed_anomalies = 3;
  int32 scans_running = 4;
  int64 timestamp = 5;
}

message StreamLogsRequest {
  string level = 1;
  repeated string components = 2;
}

message LogEntry {
  int64 timestamp = 1;
  string level = 2;
  string component = 3;
  string message = 4;
}
```

---

## Data Model

### Database Schema (TimescaleDB)

**Programs Table**
```sql
CREATE TABLE programs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    platform VARCHAR(50),
    scope TEXT[],
    scan_frequency INTERVAL DEFAULT '1 hour',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_scanned_at TIMESTAMPTZ,
    is_active BOOLEAN DEFAULT true,
    metadata JSONB
);

CREATE INDEX idx_programs_active ON programs(is_active) WHERE is_active = true;
```

**Assets Table (Hypertable)**
```sql
CREATE TABLE assets (
    id SERIAL,
    program_id INTEGER REFERENCES programs(id),
    discovered_at TIMESTAMPTZ NOT NULL,
    asset_type VARCHAR(50),
    asset_value TEXT NOT NULL,
    is_live BOOLEAN DEFAULT false,
    status_code INTEGER,
    content_hash TEXT,
    tech_stack JSONB,
    response_headers JSONB,
    cert_info JSONB,
    response_time_ms INTEGER,
    PRIMARY KEY (discovered_at, id)
);

-- Convert to TimescaleDB hypertable
SELECT create_hypertable('assets', 'discovered_at');

-- Indexes
CREATE INDEX idx_assets_program ON assets(program_id);
CREATE INDEX idx_assets_type_value ON assets(asset_type, asset_value);
CREATE INDEX idx_assets_live ON assets(is_live) WHERE is_live = true;
CREATE INDEX idx_assets_hash ON assets(content_hash);

-- Compression (90% storage savings)
ALTER TABLE assets SET (
    timescaledb.compress,
    timescaledb.compress_segmentby = 'program_id,asset_type'
);

-- Compress data older than 7 days
SELECT add_compression_policy('assets', INTERVAL '7 days');

-- Retention policy (auto-delete old data)
SELECT add_retention_policy('assets', INTERVAL '6 months');
```

**Deployment Patterns Table**
```sql
CREATE TABLE deployment_patterns (
    id SERIAL PRIMARY KEY,
    program_id INTEGER REFERENCES programs(id),
    day_of_week INTEGER,
    hour_of_day INTEGER,
    change_count INTEGER DEFAULT 0,
    avg_changes FLOAT,
    stddev_changes FLOAT,
    last_updated_at TIMESTAMPTZ,
    UNIQUE(program_id, day_of_week, hour_of_day)
);

CREATE INDEX idx_patterns_program ON deployment_patterns(program_id);
```

**Anomalies Table (Hypertable)**
```sql
CREATE TABLE anomalies (
    id SERIAL,
    detected_at TIMESTAMPTZ NOT NULL,
    program_id INTEGER REFERENCES programs(id),
    asset_id INTEGER,
    anomaly_type VARCHAR(100),
    description TEXT,
    evidence JSONB,
    base_probability FLOAT,
    posterior_probability FLOAT,
    priority_score FLOAT,
    is_reviewed BOOLEAN DEFAULT false,
    review_notes TEXT,
    reviewed_at TIMESTAMPTZ,
    PRIMARY KEY (detected_at, id)
);

SELECT create_hypertable('anomalies', 'detected_at');

CREATE INDEX idx_anomalies_program ON anomalies(program_id);
CREATE INDEX idx_anomalies_unreviewed ON anomalies(is_reviewed) WHERE is_reviewed = false;
CREATE INDEX idx_anomalies_priority ON anomalies(priority_score DESC);
```

**Findings Table**
```sql
CREATE TABLE findings (
    id SERIAL PRIMARY KEY,
    program_id INTEGER REFERENCES programs(id),
    anomaly_id INTEGER,
    title VARCHAR(500),
    severity VARCHAR(20),
    status VARCHAR(50),
    reported_at TIMESTAMPTZ DEFAULT NOW(),
    resolved_at TIMESTAMPTZ,
    bounty_amount DECIMAL(10,2),
    notes TEXT
);

CREATE INDEX idx_findings_program ON findings(program_id);
CREATE INDEX idx_findings_status ON findings(status);
```

**Scan Jobs Table**
```sql
CREATE TABLE scan_jobs (
    id SERIAL PRIMARY KEY,
    program_id INTEGER REFERENCES programs(id),
    job_type VARCHAR(50),
    status VARCHAR(20),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    results_count INTEGER,
    error_message TEXT,
    metadata JSONB
);

CREATE INDEX idx_scans_program ON scan_jobs(program_id);
CREATE INDEX idx_scans_status ON scan_jobs(status);
```

**Continuous Aggregates (Pre-computed Views)**

```sql
-- Daily asset statistics
CREATE MATERIALIZED VIEW daily_asset_stats
WITH (timescaledb.continuous) AS
SELECT 
    time_bucket('1 day', discovered_at) AS day,
    program_id,
    asset_type,
    COUNT(*) as asset_count,
    COUNT(*) FILTER (WHERE is_live) as live_count,
    AVG(response_time_ms) as avg_response_time
FROM assets
GROUP BY day, program_id, asset_type;

-- Auto-refresh the aggregate
SELECT add_continuous_aggregate_policy('daily_asset_stats',
    start_offset => INTERVAL '1 month',
    end_offset => INTERVAL '1 hour',
    schedule_interval => INTERVAL '1 hour');

-- Anomaly summary by program
CREATE MATERIALIZED VIEW anomaly_summary
WITH (timescaledb.continuous) AS
SELECT 
    time_bucket('1 hour', detected_at) AS hour,
    program_id,
    anomaly_type,
    COUNT(*) as count,
    AVG(priority_score) as avg_priority,
    MAX(priority_score) as max_priority
FROM anomalies
GROUP BY hour, program_id, anomaly_type;
```

### Data Relationships

```
programs (1) ──────< (N) assets
programs (1) ──────< (N) anomalies
programs (1) ──────< (N) findings
programs (1) ──────< (N) deployment_patterns
programs (1) ──────< (N) scan_jobs

assets (1) ────────< (N) anomalies (optional FK)
anomalies (1) ─────< (N) findings (optional FK)
```

---

## Deployment Strategy

### Infrastructure Provisioning (Terraform)

**VPS Setup**

```hcl
# terraform/main.tf
terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2.0"
    }
  }
}

provider "digitalocean" {
  token = var.do_token
}

resource "digitalocean_droplet" "recon_platform" {
  image  = "ubuntu-22-04-x64"
  name   = "recon-platform"
  region = "nyc3"
  size   = "s-2vcpu-4gb"
  
  ssh_keys = [var.ssh_key_fingerprint]
  
  user_data = templatefile("${path.module}/scripts/install-k3s.sh", {
    k3s_token = var.k3s_token
  })
  
  tags = ["recon-platform", "production"]
}

resource "digitalocean_firewall" "recon" {
  name = "recon-platform-firewall"
  
  droplet_ids = [digitalocean_droplet.recon_platform.id]
  
  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = [var.admin_ip]
  }
  
  inbound_rule {
    protocol         = "tcp"
    port_range       = "6443"
    source_addresses = [var.admin_ip]
  }
  
  inbound_rule {
    protocol         = "tcp"
    port_range       = "8080"
    source_addresses = [var.admin_ip]
  }
  
  inbound_rule {
    protocol         = "tcp"
    port_range       = "9090"
    source_addresses = [var.admin_ip]
  }
  
  outbound_rule {
    protocol              = "tcp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
  
  outbound_rule {
    protocol              = "udp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
}

output "droplet_ip" {
  value = digitalocean_droplet.recon_platform.ipv4_address
}
```

**k3s Installation Script**

```bash
#!/bin/bash
# terraform/scripts/install-k3s.sh

set -e

# Update system
apt-get update
apt-get upgrade -y

# Install k3s
curl -sfL https://get.k3s.io | sh -s - \
  --token ${k3s_token} \
  --write-kubeconfig-mode 644

# Wait for k3s to be ready
until kubectl get nodes; do
  echo "Waiting for k3s..."
  sleep 5
done

# Install Helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

echo "k3s installation complete"
```

### Kubernetes Manifests

**Namespace**

```yaml
# k8s/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: recon-platform
```

**ConfigMap**

```yaml
# k8s/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: recon-config
  namespace: recon-platform
data:
  DB_HOST: "timescaledb"
  DB_PORT: "5432"
  DB_NAME: "recon_platform"
  REDIS_HOST: "redis"
  REDIS_PORT: "6379"
  REST_API_PORT: "8080"
  GRPC_PORT: "9090"
  LOG_LEVEL: "info"
```

**Secret**

```yaml
# k8s/secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: recon-secrets
  namespace: recon-platform
type: Opaque
stringData:
  DB_PASSWORD: "your-secure-password"
  API_KEY: "your-api-key"
  DISCORD_WEBHOOK_URL: "https://discord.com/api/webhooks/..."
```

**TimescaleDB StatefulSet**

```yaml
# k8s/timescaledb.yaml
apiVersion: v1
kind: Service
metadata:
  name: timescaledb
  namespace: recon-platform
spec:
  ports:
  - port: 5432
  selector:
    app: timescaledb
  clusterIP: None
---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: timescaledb
  namespace: recon-platform
spec:
  serviceName: timescaledb
  replicas: 1
  selector:
    matchLabels:
      app: timescaledb
  template:
    metadata:
      labels:
        app: timescaledb
    spec:
      containers:
      - name: timescaledb
        image: timescale/timescaledb:latest-pg16
        ports:
        - containerPort: 5432
        env:
        - name: POSTGRES_DB
          valueFrom:
            configMapKeyRef:
              name: recon-config
              key: DB_NAME
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: recon-secrets
              key: DB_PASSWORD
        volumeMounts:
        - name: data
          mountPath: /var/lib/postgresql/data
        resources:
          requests:
            memory: "1Gi"
            cpu: "500m"
          limits:
            memory: "2Gi"
            cpu: "1000m"
  volumeClaimTemplates:
  - metadata:
      name: data
    spec:
      accessModes: ["ReadWriteOnce"]
      resources:
        requests:
          storage: 20Gi
```

**Redis Deployment**

```yaml
# k8s/redis.yaml
apiVersion: v1
kind: Service
metadata:
  name: redis
  namespace: recon-platform
spec:
  ports:
  - port: 6379
  selector:
    app: redis
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis
  namespace: recon-platform
spec:
  replicas: 1
  selector:
    matchLabels:
      app: redis
  template:
    metadata:
      labels:
        app: redis
    spec:
      containers:
      - name: redis
        image: redis:7-alpine
        ports:
        - containerPort: 6379
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

**API Server Deployment**

```yaml
# k8s/api-server.yaml
apiVersion: v1
kind: Service
metadata:
  name: api-server
  namespace: recon-platform
spec:
  type: LoadBalancer
  ports:
  - name: rest
    port: 8080
    targetPort: 8080
  - name: grpc
    port: 9090
    targetPort: 9090
  selector:
    app: api-server
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-server
  namespace: recon-platform
spec:
  replicas: 2
  selector:
    matchLabels:
      app: api-server
  template:
    metadata:
      labels:
        app: api-server
    spec:
      containers:
      - name: api-server
        image: recon-platform/api-server:latest
        ports:
        - containerPort: 8080
        - containerPort: 9090
        envFrom:
        - configMapRef:
            name: recon-config
        env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: recon-secrets
              key: DB_PASSWORD
        - name: API_KEY
          valueFrom:
            secretKeyRef:
              name: recon-secrets
              key: API_KEY
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

**Worker Deployment**

```yaml
# k8s/worker.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: worker
  namespace: recon-platform
spec:
  replicas: 3
  selector:
    matchLabels:
      app: worker
  template:
    metadata:
      labels:
        app: worker
    spec:
      containers:
      - name: worker
        image: recon-platform/worker:latest
        envFrom:
        - configMapRef:
            name: recon-config
        env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: recon-secrets
              key: DB_PASSWORD
        resources:
          requests:
            memory: "512Mi"
            cpu: "500m"
          limits:
            memory: "1Gi"
            cpu: "1000m"
```

**CronJob for Scheduled Scans**

```yaml
# k8s/cronjob-passive-scan.yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: passive-recon
  namespace: recon-platform
spec:
  schedule: "0 * * * *"  # Every hour
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: scheduler
            image: recon-platform/scheduler:latest
            args: ["scan", "--type=passive"]
            envFrom:
            - configMapRef:
                name: recon-config
            env:
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: recon-secrets
                  key: DB_PASSWORD
          restartPolicy: OnFailure
---
apiVersion: batch/v1
kind: CronJob
metadata:
  name: pattern-analysis
  namespace: recon-platform
spec:
  schedule: "0 0 * * *"  # Daily at midnight
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: analyzer
            image: recon-platform/scheduler:latest
            args: ["analyze-patterns"]
            envFrom:
            - configMapRef:
                name: recon-config
            env:
            - name: DB_PASSWORD
              valueFrom:
                secretKeyRef:
                  name: recon-secrets
                  key: DB_PASSWORD
          restartPolicy: OnFailure
```

### Deployment Commands

```bash
# 1. Provision infrastructure
cd terraform
terraform init
terraform plan
terraform apply

# 2. Get kubeconfig from VPS
export KUBECONFIG=~/.kube/recon-platform
scp root@<droplet-ip>:/etc/rancher/k3s/k3s.yaml ~/.kube/recon-platform
# Edit kubeconfig to replace 127.0.0.1 with droplet IP

# 3. Deploy to Kubernetes
kubectl apply -f k8s/namespace.yaml
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secret.yaml
kubectl apply -f k8s/timescaledb.yaml
kubectl apply -f k8s/redis.yaml

# Wait for database to be ready
kubectl wait --for=condition=ready pod -l app=timescaledb -n recon-platform --timeout=300s

# Run database migrations
kubectl run -it --rm migrate --image=recon-platform/migrate:latest \
  --env="DB_HOST=timescaledb" \
  --env="DB_PASSWORD=<password>" \
  -n recon-platform -- migrate up

# Deploy application
kubectl apply -f k8s/api-server.yaml
kubectl apply -f k8s/worker.yaml
kubectl apply -f k8s/cronjob-passive-scan.yaml

# 4. Get LoadBalancer IP
kubectl get svc api-server -n recon-platform

# 5. Configure CLI
recon-cli config set server http://<loadbalancer-ip>:8080
recon-cli config set grpc-server <loadbalancer-ip>:9090
recon-cli config set api-key <your-api-key>
```

---

## Development Roadmap

### MVP Timeline: 2-3 Weeks

#### Week 1: Core Application

**Day 1-2: Project Setup & Database**
- [ ] Initialize Go modules
- [ ] Set up project structure
- [ ] Create database schema
- [ ] Deploy TimescaleDB locally (Docker)
- [ ] Implement pgx connection pool
- [ ] Write migration scripts
- [ ] Test database connectivity

**Day 3-4: Worker Implementation**
- [ ] Implement worker that runs subfinder
- [ ] Implement worker that runs httpx
- [ ] Parse and normalize tool outputs
- [ ] Store results in TimescaleDB
- [ ] Implement basic diff logic (detect new subdomains)
- [ ] Test with real bug bounty program

**Day 5-6: Task Queue & Scheduler**
- [ ] Set up River/Asynq
- [ ] Implement job creation
- [ ] Workers pull from queue
- [ ] Add retry logic
- [ ] Implement scheduled jobs (hourly scans)
- [ ] Test job execution

**Day 7: CLI Tool**
- [ ] Set up Cobra CLI structure
- [ ] Implement `program add` command
- [ ] Implement `scan trigger` command
- [ ] Implement `anomalies list` command
- [ ] REST client library
- [ ] **First end-to-end test: Add program → Trigger scan → View results**

#### Week 2: Deployment & Polish

**Day 8-9: Containerization**
- [ ] Write Dockerfiles (multi-stage builds)
- [ ] Docker Compose for local development
- [ ] Test containers locally
- [ ] Push images to registry (Docker Hub/GHCR)
- [ ] Document container usage

**Day 10-11: Kubernetes**
- [ ] Write k8s manifests
- [ ] Test on local minikube
- [ ] Deploy StatefulSet (database)
- [ ] Deploy Deployments (API, workers)
- [ ] Deploy CronJobs
- [ ] Verify all components communicate

**Day 12-13: VPS + Terraform**
- [ ] Write Terraform configuration
- [ ] Provision DigitalOcean droplet
- [ ] Install k3s via Terraform
- [ ] Deploy k8s manifests to VPS
- [ ] Configure networking/firewall
- [ ] **Production MVP running on VPS**

**Day 14: Alerts & Testing**
- [ ] Implement Discord webhook alerts
- [ ] Test with 2-3 real programs
- [ ] Monitor for 24 hours
- [ ] Fix bugs, tune performance
- [ ] Document usage

### Post-MVP: Weeks 3-8

#### Week 3: Enhanced Scoring
- [ ] Implement Bayesian scoring algorithm
- [ ] Add evidence factors
- [ ] Test scoring accuracy
- [ ] Tune priors based on results

#### Week 4: Pattern Learning
- [ ] Implement deployment pattern analysis
- [ ] Build pattern learning job
- [ ] Store patterns in database
- [ ] Use patterns in scoring

#### Week 5: gRPC Streaming
- [ ] Define Protocol Buffers
- [ ] Implement gRPC server
- [ ] Add streaming endpoints
- [ ] Update CLI with gRPC client
- [ ] Build live dashboard (TUI)

#### Week 6: Additional Tools
- [ ] Integrate amass
- [ ] Integrate nuclei
- [ ] Integrate LinkFinder
- [ ] Add technology detection

#### Week 7: Monitoring & Observability
- [ ] Deploy Prometheus
- [ ] Deploy Grafana
- [ ] Create dashboards
- [ ] Set up alerting rules

#### Week 8: Polish & Optimization
- [ ] Performance tuning
- [ ] Database query optimization
- [ ] Add unit tests
- [ ] Add integration tests
- [ ] Documentation improvements

---

## Cost Analysis

### Development Phase (Weeks 1-2)

**Infrastructure:**
- Local development: $0 (Docker Desktop)
- DigitalOcean free trial: $200 credit (2 months)
- **Total: $0**

**Tools & Services:**
- GitHub (free tier)
- Docker Hub (free tier)
- **Total: $0**

**Time Investment:**
- ~40 hours (2 weeks × 20 hours/week)
- **Total Cost: $0 out of pocket**

### Production (Monthly Costs)

**Recommended Setup: Single VPS**

| Item | Provider | Specs | Cost |
|------|----------|-------|------|
| VPS | DigitalOcean | 2 vCPU, 4GB RAM, 80GB SSD | $24/month |
| OR VPS | Contabo | 4 vCPU, 8GB RAM, 200GB SSD | $7/month |
| Domain (optional) | Namecheap | .com domain | $12/year |

**Recommended: Contabo VPS = $7/month**

**Total Monthly Cost: $7-24/month**

### Scaling Costs (Future)

**If monitoring 20+ programs:**

| Item | Specs | Cost |
|------|-------|------|
| VPS | 4 vCPU, 8GB RAM | $48/month |
| Block Storage | +50GB | $5/month |
| **Total** | | **$53/month** |

**If adding monitoring stack (Prometheus/Grafana):**
- Additional 2GB RAM needed
- Upgrade to 6GB plan: +$12/month

### Cost Comparison

| Solution | Monthly Cost | Capabilities |
|----------|--------------|--------------|
| Manual Recon | $0 | Limited, time-intensive |
| This Platform (MVP) | $7 | Continuous, intelligent, scalable |
| Commercial ASM Platform | $500-2000+ | Enterprise features, overkill |

**ROI:** If platform helps find 1 additional bug per month worth $500+, it pays for itself 70x over.

---

## Security & Legal

### Legal Framework

**Bug Bounty Authorization**
- Only scan programs explicitly enrolled in
- Respect scope limitations (in-scope domains only)
- Follow program rules (rate limits, prohibited actions)
- Maintain records of authorization
- Screenshot program policies for documentation

**Safe Harbor Compliance**
- Bug bounty programs provide safe harbor protection
- Testing is authorized under program terms
- Follow responsible disclosure practices
- Do not exceed scope or authorization

**Best Practices**
- Document all testing activities
- Keep logs of scans and findings
- Respect disclosure timelines
- Never share vulnerabilities publicly before resolution
- Follow CVD (Coordinated Vulnerability Disclosure) principles

### Security Considerations

**API Security**
- API key authentication
- Rate limiting (prevent abuse)
- Input validation
- Parameterized SQL queries (prevent injection)
- TLS/SSL for all communications

**Data Protection**
- Encrypt sensitive data at rest
- Secure Kubernetes Secrets
- No credentials in code or Git
- Regular security updates
- Minimal data retention

**Network Security**
- Firewall rules (restrict access)
- VPN or IP whitelisting for admin access
- Separate production/development environments
- Regular security audits

**Operational Security**
- All scans run from VPS IP (not home IP)
- Rate-limited scanning (respectful, non-aggressive)
- User-Agent headers identify tool
- Maintain scan logs for accountability
- Monitor for abuse/misuse

### Ethical Practices

**Rate Limiting**
- Maximum 10 requests/second per target
- Delays between large scans
- Respect robots.txt
- No resource exhaustion attacks

**Responsible Disclosure**
- Report vulnerabilities promptly
- Provide clear reproduction steps
- Allow reasonable remediation time
- Never weaponize or sell vulnerabilities
- Follow program disclosure policies

**Scope Compliance**
- Automated scope validation
- Reject out-of-scope targets
- Alert on scope violations
- Maintain audit trail

---

## Success Metrics

### Platform Performance Metrics

**Operational Metrics**
- Uptime: Target 99.5% (< 4 hours downtime/month)
- Scan completion rate: > 95%
- Average scan duration: < 10 minutes for passive recon
- Worker utilization: 60-80%
- Database query performance: < 100ms average

**Detection Metrics**
- New assets discovered per day
- Changes detected per day
- Anomalies generated per day
- False positive rate: < 20%
- Time to detection: < 5 minutes

### Bug Bounty Success Metrics

**Efficiency Gains**
- Recon time reduction: 80%+ (from hours to minutes)
- Programs monitored simultaneously: 10+
- Coverage: 100% of in-scope assets
- First-mover advantage: Detect changes within 5 minutes

**Bounty Outcomes**
- Bugs found per month
- Bounty earnings per month
- Hit rate (anomalies → bugs): Target 5-10%
- Average time from detection to report: < 1 hour

**Learning & Improvement**
- Scoring accuracy improvement over time
- Pattern recognition effectiveness
- Alert relevance (reviewed vs ignored)
- Bayesian prior refinement

### Key Performance Indicators (KPIs)

**Week 1 Goals:**
- [ ] Platform deployed and operational
- [ ] 1 program monitored successfully
- [ ] First anomaly detected and alerted

**Month 1 Goals:**
- [ ] 5 programs monitored
- [ ] 50+ anomalies detected
- [ ] 1+ bug reported
- [ ] < 30% false positive rate

**Month 3 Goals:**
- [ ] 10+ programs monitored
- [ ] Bayesian scoring tuned and effective
- [ ] Pattern learning operational
- [ ] 5+ bugs reported
- [ ] Platform pays for itself (bounties > costs)

**Month 6 Goals:**
- [ ] 15+ programs monitored
- [ ] Advanced features deployed
- [ ] Hit rate > 5%
- [ ] ROI > 10x (earnings vs. costs)

---

## Appendix

### A. Technology Decision Matrix

| Decision | Options Considered | Selected | Rationale |
|----------|-------------------|----------|-----------|
| Language | Python, Go, Rust | **Go** | Performance, concurrency, single binary, k8s native |
| Database | PostgreSQL, TimescaleDB, InfluxDB | **TimescaleDB** | Time-series optimized, SQL interface, 1000x faster |
| Task Queue | Celery, RQ, Asynq, River | **River/Asynq** | Go-native, simpler than Celery, production-ready |
| API | REST, gRPC, GraphQL | **REST + gRPC** | REST for simplicity, gRPC for streaming |
| Orchestration | Docker Compose, k8s, Nomad | **Kubernetes (k3s)** | Industry standard, learning value, scalable |
| IaC | Terraform, Pulumi, CloudFormation | **Terraform** | Industry standard, multi-cloud, declarative |
| Hosting | Home server, VPS, Cloud | **VPS** | Cost-effective, dedicated IP, 24/7 uptime |

### B. Recon Tools Integration

**Current Tools (MVP):**
- subfinder: Fast subdomain enumeration
- httpx: HTTP probing and tech detection
- certsh: Certificate transparency logs

**Future Tools:**
- amass: Advanced subdomain enumeration
- nuclei: Vulnerability scanning
- naabu: Port scanning
- katana: Web crawling
- LinkFinder: Endpoint discovery in JavaScript
- wappalyzer: Technology detection
- ffuf: Fuzzing and content discovery

**Tool Execution Pattern:**
```go
type ReconTool interface {
    Run(ctx context.Context, target string) ([]Result, error)
    Name() string
    Type() ToolType
}

// Example: Subfinder
type SubfinderTool struct{}

func (s *SubfinderTool) Run(ctx context.Context, target string) ([]Result, error) {
    cmd := exec.CommandContext(ctx, "subfinder", "-d", target, "-silent", "-json")
    output, err := cmd.Output()
    // Parse and return results
}
```

### C. Bayesian Scoring Deep Dive

**Mathematical Foundation:**

```
P(Bug|Evidence) = P(Evidence|Bug) × P(Bug) / P(Evidence)

Simplified for implementation:
posterior = prior × ∏(evidence_factors)

Where:
- prior = base probability (learned per program)
- evidence_factors = multipliers for each signal
- posterior = final probability estimate
```

**Evidence Factors (Initial Values):**

| Evidence | Multiplier | Rationale |
|----------|-----------|-----------|
| Weekend deployment | 5.0x | Emergency fixes, less testing |
| Outside normal window | 2.0x | Unplanned deployment |
| Dev/staging subdomain | 3.0x | Often less hardened |
| New subdomain | 4.0x | New code, new bugs |
| Tech stack change | 1.5x | Migration issues |
| Multiple rapid changes | 2.5x | Panic mode, shortcuts |
| Certificate change | 1.8x | Infrastructure change |
| Status code change (403→200) | 3.0x | Access control change |
| Version rollback | 2.5x | Reverting due to issues |

**Learning Algorithm:**

```
After each bug report:
1. Record evidence factors present
2. Calculate accuracy: predicted vs actual
3. Adjust factors using gradient descent
4. Update per-program priors
5. Retrain model monthly
```

### D. Deployment Checklist

**Pre-Deployment:**
- [ ] Terraform configuration reviewed
- [ ] Kubernetes manifests validated
- [ ] Docker images built and pushed
- [ ] Secrets generated and stored securely
- [ ] Database migrations tested
- [ ] Backup strategy defined

**Deployment:**
- [ ] Infrastructure provisioned (terraform apply)
- [ ] Kubeconfig obtained and configured
- [ ] Namespace and configs created
- [ ] Database deployed and initialized
- [ ] Migrations run successfully
- [ ] Application deployed
- [ ] Health checks passing
- [ ] CronJobs scheduled

**Post-Deployment:**
- [ ] Add first program via CLI
- [ ] Trigger manual scan
- [ ] Verify results in database
- [ ] Test alert notifications
- [ ] Monitor logs for errors
- [ ] Document any issues
- [ ] Performance baseline established

**Ongoing:**
- [ ] Weekly database backups
- [ ] Monthly security updates
- [ ] Quarterly cost review
- [ ] Semi-annual architecture review

### E. Troubleshooting Guide

**Common Issues:**

**Workers not processing jobs:**
- Check Redis connectivity
- Verify queue name matches
- Check worker logs
- Ensure tools installed in container

**Database connection errors:**
- Verify TimescaleDB StatefulSet is running
- Check password in Secret
- Verify network policies
- Test connection from worker pod

**Anomalies not generating:**
- Verify diff engine running
- Check scan results exist
- Review scoring thresholds
- Inspect diff logic

**Alerts not sending:**
- Verify webhook URLs
- Check network egress rules
- Review alert service logs
- Test webhook manually

### F. Future Enhancements

**Phase 2 (Months 3-6):**
- Web UI for dashboard
- Multi-user support
- Advanced ML models (LSTM for pattern prediction)
- Integration with more recon tools
- Custom wordlist generation
- Automated reporting (generate POC drafts)

**Phase 3 (Months 6-12):**
- Collaboration features (team support)
- Competitive intelligence (track other researchers)
- Historical trending and forecasting
- API marketplace integration
- Mobile app (iOS/Android)
- Browser extension

**Phase 4 (Year 2+):**
- SaaS offering (sell to other researchers)
- Managed service option
- Enterprise features
- Advanced threat intelligence
- Automated exploitation (with permission)
- Integration with SIEM platforms

### G. Resources & References

**Documentation:**
- Go: https://go.dev/doc/
- Kubernetes: https://kubernetes.io/docs/
- Terraform: https://www.terraform.io/docs/
- TimescaleDB: https://docs.timescale.com/
- gRPC: https://grpc.io/docs/
- Cobra: https://cobra.dev/

**Learning Resources:**
- Bug Bounty: HackerOne University, Bugcrowd University
- Kubernetes: "Kubernetes Up & Running" (O'Reilly)
- Go: "The Go Programming Language" (Donovan & Kernighan)
- Terraform: Official HashiCorp tutorials

**Community:**
- r/bugbounty
- Bug bounty Discord servers
- HackerOne community forums
- Kubernetes Slack

### H. Glossary

**ASM**: Attack Surface Management
**CRUD**: Create, Read, Update, Delete
**gRPC**: Google Remote Procedure Call
**IaC**: Infrastructure as Code
**k3s**: Lightweight Kubernetes distribution
**MVP**: Minimum Viable Product
**POC**: Proof of Concept
**ROI**: Return on Investment
**SaaS**: Software as a Service
**TLS**: Transport Layer Security
**TUI**: Terminal User Interface
**VPS**: Virtual Private Server

---

## Conclusion

This platform represents a **modern, intelligent approach to bug bounty reconnaissance** that combines cutting-edge technologies (Go, Kubernetes, TimescaleDB, gRPC) with proven security research methodologies.

**Key Takeaways:**

1. **Competitive Advantage**: Continuous monitoring and intelligent anomaly detection provide a significant edge in the time-sensitive world of bug bounties

2. **Technical Excellence**: Modern architecture using industry-standard tools and practices provides both immediate value and long-term learning

3. **Rapid Deployment**: 2-3 week MVP timeline means you'll be operational quickly while maintaining quality

4. **Cost Effective**: $7-24/month operational cost with potential for 10x+ ROI through increased bug discovery

5. **Scalable Design**: Architecture supports growth from 5 programs to 50+ without major refactoring

**Next Steps:**

1. **Week 1**: Build core application locally
2. **Week 2**: Deploy to production VPS with Kubernetes
3. **Week 3+**: Add intelligence features and scale

**Success Criteria:**

- Platform operational 24/7
- First bug discovered via platform within 30 days
- ROI positive within 90 days
- Professional DevOps skills demonstrated

This document serves as the blueprint for building a production-grade reconnaissance platform that will enhance your bug bounty hunting effectiveness while providing valuable experience with modern cloud-native technologies.

---

**Document Version Control:**

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| 1.0 | 2025-10-01 | Initial vision document | Platform Architect |

**Approval:**

This document represents the technical vision and implementation plan for the Bug Bounty Continuous Reconnaissance Platform. It will be updated as the project evolves and new requirements emerge.

---

**End of Document**



