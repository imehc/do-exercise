// Package integration 承载租户隔离回归测试（真实 PostgreSQL）。
// 用例文件都带 `//go:build integration`，日常 go test ./... 不会运行；
// 需要真实库时用 `make test-integration`（先创建独立测试库 do_exercise_test）。
package integration
