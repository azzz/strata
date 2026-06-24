# Read Path:

Read Path:

```
Parser-Specific AST
 ↓
Query
 ↓
Planner
 ├─ ask catalog: table schema
 ├─ ask catalog: current manifest
 ├─ check columns and types
 ├─ pruning files accordint to manifest 
 ↓
PhysicalPlan
 ↓
Executor
 ↓
Storage
 ↓
{Physical File}
```

## Project structure
```
docs/ -- you are here
engine/
- catalog   // table schema, manifest etc
- query     // Query, Expressions
- planner   // Creates physical plan from a query
- plan      // Physical Plan
- executor  // Runs a physical plan
- storage   // Storage adapters: disks, s3 etc
```

## Implementation phases

### Phase 1

Basic MVP. Goal: read part files.

```
Physical Plan
 ↓
Executor
 ↓
Storage
```

### Phase 2

Goal: execute queries

```
Query
 ↓
Planner
 ↓
{ Phase 1 }
```

### Phase 3

Goal: prepare for parsers

```
AST
 ↓
{ Phase 2 }
```

# Storage Model

Data is stored as immutable parts.

Example:

```
events/
  metadata/
    CURRENT
    manifest_000001.json

  data/
    dt=2026-06-22/
      part_1.parquet
      part_2.parquet

    dt=2026-06-23/
      part_3.parquet
```

A part is an immutable data file (initially Parquet).

Updates never modify existing parts.

New data is written into new parts.

Background compaction merges multiple parts into a larger one.
