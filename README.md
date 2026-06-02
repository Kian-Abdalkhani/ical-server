1. typetype
2. render
3. db integration
4. handler integration tests

- Use \_test.go in same package
- use t.Helper() for shared assertion helpers
- for db testing, create setupTestDB \*db.Queries for opening :memory: and uses t.Cleanup() to close it
- for handler tests, use httptest.NewServer or just httptest.NewRecorder
