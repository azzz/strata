# Merge Algorithm

## Assumptions

Every part is immutable.

Every part is sorted by table primary/sort key.

Example:

```text
part_1
-------
key1
key1
key2
key3

part_2
-------
key1
key4
key5
```

Merge operations never modify input parts.

A merge always creates a new output part.

---

# Goal

Given multiple sorted parts:

```text
part_1
part_2
...
part_N
```

produce a new sorted part:

```text
part_M
```

while applying table merge strategy:

* AppendOnly
* Replacing
* Summing
* Aggregating
* etc.

---

# High-Level Algorithm

Merge is implemented as a streaming k-way merge.

Input parts are never fully loaded into memory.

The merger maintains:

```text
one iterator per input part
```

Each iterator exposes:

```go
type Iterator interface {
    Current() Row
    Next() bool
}
```

The merger keeps only the current row of every iterator in memory.

---

# Processing

At every step:

1. Find the smallest key among all iterators.
2. Call it currentKey.
3. Collect all rows having currentKey from all iterators.
4. Apply merge strategy.
5. Write resulting row into output part.
6. Continue with next key.

Example:

Input:

```text
part_1:
  key1
  key1
  key2
  key3

part_2:
  key1
  key4
  key5
```

Step 1:

```text
currentKey = key1
```

Collect:

```text
part_1:
  key1
  key1

part_2:
  key1
```

Result:

```text
merge(key1 rows)
```

Write:

```text
key1
```

Iterators now point to:

```text
part_1 -> key2
part_2 -> key4
```

Next:

```text
currentKey = key2
```

and so on.

---

# Why This Works

Because every part is sorted.

Once an iterator reaches:

```text
key2
```

it is guaranteed that:

```text
key1
```

will never appear again in that iterator.

Therefore rows can be processed in a streaming fashion.

---

# Complexity

For two parts:

```text
O(rows_part_1 + rows_part_2)
```

For N parts:

```text
O(total_rows × log(N))
```

when a min-heap is used.

Memory usage:

```text
O(number_of_parts)
```

plus accumulator state for the currently processed key.

Memory consumption is independent of total dataset size.

---

# Merge Strategies

## AppendOnly

Rows are copied as-is.

Duplicate primary keys are allowed.

Example:

```text
key1 value=5
key1 value=17
```

Output:

```text
key1 value=5
key1 value=17
```

---

## Replacing

Keep one row for a key.

Example:

```text
key1 value=5
key1 value=17
```

Output:

```text
key1 value=17
```

---

## Summing

Aggregate numeric columns.

Example:

```text
key1 requests=5
key1 requests=17
key1 requests=3
```

Output:

```text
key1 requests=25
```

---

# Important Invariant

All parts must be sorted by table primary/sort key.

Without sorted parts:

```text
streaming merge becomes impossible
```

and the system would require:

* full hash aggregation
* external sorting
* significantly higher memory usage

Therefore sorted immutable parts are a fundamental requirement of the storage engine.
