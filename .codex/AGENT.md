# Project Overview

This project is an embedded analytical database engine written in Go.

The goal is to provide a reusable OLAP engine that can be embedded into applications, similar in spirit to SQLite, but optimized for analytical workloads.

The engine is based on immutable data parts and snapshot-based metadata.

Data is stored as Parquet files.

Parsers (SQL, REST, GraphQL, etc.) are intentionally kept outside of the engine core.

## Architecture

Read Path:

```text
AST
 ↓
Query
 ↓
Planner
 ↓
PhysicalPlan
 ↓
Executor
 ↓
Storage
 ↓
Physical Files
```

Write Path:

```text
Insert
 ↓
Writer
 ↓
Part Builder
 ↓
Storage
 ↓
Part File
 ↓
Catalog Commit
```

Merge Path:

```text
Catalog Snapshot
 ↓
Merger
 ↓
Read Parts
 ↓
Write Merged Part
 ↓
Catalog Replace Commit
```

## Core Concepts

### Catalog

Catalog is the source of truth.

Catalog manages:

* schemas
* manifests
* snapshots
* active parts

Queries must never scan directories directly.

### Query

Query is the normalized semantic representation of a request.

The engine operates on Query objects.

### PhysicalPlan

PhysicalPlan is an execution instruction.

It specifies:

* which parts to read
* which columns to read
* which filters to apply

### Executor

Executor executes a PhysicalPlan.

Executor should not perform path discovery.

### Storage

Storage is responsible for physical I/O.

Storage implementations may support:

* local filesystem
* S3
* other object stores

### Parts

Parts are immutable.

Updates never modify existing parts.

New data is written into new parts.

### Merging

Background merges replace multiple parts with a larger part.

Old parts remain available until garbage collection removes them.

## Current Development Plan

Phase 1:

```text
PhysicalPlan
 ↓
Executor
 ↓
Storage
```

Phase 2:

```text
Query
 ↓
Planner
 ↓
PhysicalPlan
```

Phase 3:

```text
AST
 ↓
Query
 ↓
Planner
```

Future phases:

* manifests
* part pruning
* bloom filters
* write path
* merger
* garbage collection
