# pipeline

Generic stage-based data pipeline framework in Go, with a blockchain block processing demo.

## Architecture

The pipeline connects three stages via buffered channels:

```
Fetch → Transform → Store
```

Each stage runs in its own goroutine, reads from an inbound channel, processes data through a worker function, and sends results to an outbound channel. Errors from any stage cancel the entire pipeline via context.

## Usage

```bash
go run main.go start
```

The demo simulates blockchain block fetching (~4s/block), transformation, and storage with configurable failure injection.

## Built with

- Go 1.23
- [Cobra](https://github.com/spf13/cobra) CLI framework
