# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2026-01-19

### ⚠️ Breaking Changes
- **Almanac API**: `NewAlmanac` no longer takes a slice of facts. It now uses functional options (`NewAlmanac(WithAlmanacConditionCaching())`). Facts are added via `almanac.AddFact`.
- **Engine Results**: `Run()` results are now available via the `Results()` and `ReduceResults()` methods on the engine instance for better ergonomics.
- **Rule Compilation**: Rules are now compiled when added to the engine.

### ✨ Added
- **Performance**: Condition result caching (globally or per Almanac).
- **Optimization**: "Smart Skip" dependency tracking (skips rules with missing facts).
- **Observability**: `MetricsCollector` interface for real-time monitoring (Duration, Success, Count).
- **Audit**: Detailed `AuditTrace` capturing every condition evaluation and fact value.
- **Hot-Reload**: Dynamic rule updates from HTTP/JSON sources without engine restart.
- **Operators**: New `regex` operator for pattern matching on strings.
- **API**: Reference REST API implementation in `docs/examples/api`.

### ⚡ Performance Improvements
- **63% reduction** in execution time for rules with shared conditions.
- **8x speedup** for partial fact sets using Smart Skip.
- Concurrent rule evaluation support via Worker Pool.

### ✅ Quality
- Maintained **100% test coverage** across all new features.
- Added comprehensive benchmark suite (`src/engine_bench_test.go`).
- Complete architecture documentation with 8 Mermaid diagrams.

---

## [1.0.0] - 2026-01-10

### Added
- Core engine with `all`/`any`/`none` conditions.
- Basic operator set (equal, greater_than, contains, etc.).
- Event system (Sync/Async).
- JSON support.
- Almanac fact management.
