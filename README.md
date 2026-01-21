# axiom-ingest-gateway

## Purpose
The ingest gateway is the external entry point of the Axiom platform.  
It accepts events from clients, validates and authenticates requests, applies admission control, and forwards accepted events to the internal event log.

This service is optimized for low-latency ingress and resilience under bursty, untrusted traffic.

## Responsibilities
- Accept client-facing write requests (HTTP)
- Authenticate and authorize producers
- Validate event schemas and payloads
- Apply rate limiting and backpressure
- Translate external requests into internal event records
- Forward events to the event log engine

## Non-Responsibilities
- Durable storage of events
- Event ordering or replication guarantees
- Stream processing or aggregation
- Serving analytical or read queries
- System configuration or orchestration

## Failure Model
- Must fail fast under overload
- Must never corrupt or mutate event data
- If downstream systems are unavailable, ingestion must degrade gracefully

## Status
Week 0: Minimal HTTP server with `/health` endpoint.
