package benchmark; import "testing"; func BenchmarkExample(b *testing.B) { for i := 0; i < b.N; i++ { _ = i } }
