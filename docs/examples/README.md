# 🚀 Go Rules Engine Examples

This directory contains various examples to help you get started with the rules engine. Each example illustrates a specific feature of the engine.

## 📂 List of Examples

| Folder | Description |
| :--- | :--- |
| [**basic**](./basic) | Fundamental usage: simple age check with static facts. |
| [**advanced**](./advanced) | Advanced usage: dynamic facts, priorities, and complex logic. |
| [**json**](./json) | Loading rules and facts from JSON files or strings. |
| [**builder**](./builder) | Using `RuleBuilder` to construct rules elegantly (fluent API). |
| [**audit-trace**](./audit-trace) | **New**: Extraction and JSON serialization of the full evaluation history (Audit Trace). |
| [**hot-reload**](./hot-reload) | **New**: Real-time rule updates via a remote source (HTTP). |
| [**parallel**](./parallel) | **New**: Parallel rule execution to optimize performance with slow facts. |
| [**metrics**](./metrics) | **New**: How to implement monitoring and observability (Prometheus). |
| [**api**](./api) | **New**: Wrapping the engine in a stateless REST API with `net/http`. |
| [**custom-operator**](./custom-operator) | How to extend the engine by adding your own custom operators. |

## 🛠 How to run the examples

You can run any example using the `go run` command from the root of the project or directly within the example's folder.

**From the root:**
```bash
go run docs/examples/basic/main.go
```

**From the example folder:**
```bash
cd docs/examples/builder
go run main.go
```

---

## 💡 Which example to choose?

- **Beginner?** Start with [basic](./basic/main.go). It's the simplest way to understand the "Fact / Condition / Rule" workflow.
- **Need to load configurations?** Check out [json](./json/main.go) or [hot-reload](./hot-reload/main.go) for dynamic loading.
- **Need performance or complex business logic?** See [advanced](./advanced/main.go) to discover dynamic facts (callbacks) and [parallel](./parallel/main.go) for concurrent execution.
- **Observability is a priority?** Look at [metrics](./metrics/main.go) to see how to connect Prometheus.
- **Don't like writing JSON by hand?** The [builder](./builder/main.go) is for you.
