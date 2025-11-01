package service

import (
	"context"
	"strconv"
	"testing"
	"time"

	"quizizz.com/internal/domain"
	"quizizz.com/internal/repository"
)

// Benchmark tests for UserService

func BenchmarkUserService_GetByID(b *testing.B) {
	// Disable logging for benchmarks
	DisableLoggingForBenchmark(b)

	// Setup
	ctx := context.Background()
	repo := repository.NewMockUserRepository()
	service := NewUserService(repo)

	// Create test user
	user := &domain.User{
		ID:        "bench-user",
		Name:      "Benchmark User",
		Email:     "bench@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, user)
	if err != nil {
		b.Fatalf("Failed to create test user: %v", err)
	}

	// Run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GetByID(ctx, "bench-user")
		if err != nil {
			b.Fatalf("GetByID failed: %v", err)
		}
	}
}

func BenchmarkUserService_List(b *testing.B) {
	// Disable logging for benchmarks
	DisableLoggingForBenchmark(b)

	// Setup
	ctx := context.Background()
	repo := repository.NewMockUserRepository()
	service := NewUserService(repo)

	// Create test users
	for i := 0; i < 100; i++ {
		user := &domain.User{
			ID:        "bench-user-" + strconv.Itoa(i),
			Name:      "Benchmark User " + strconv.Itoa(i),
			Email:     "bench" + strconv.Itoa(i) + "@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		err := repo.Create(ctx, user)
		if err != nil {
			b.Fatalf("Failed to create test user: %v", err)
		}
	}

	// Run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.List(ctx)
		if err != nil {
			b.Fatalf("List failed: %v", err)
		}
	}
}

func BenchmarkUserService_Create(b *testing.B) {
	// Disable logging for benchmarks
	DisableLoggingForBenchmark(b)

	// Setup
	ctx := context.Background()
	repo := repository.NewMockUserRepository()
	service := NewUserService(repo)

	// Run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		user := &domain.User{
			ID:        "bench-create-" + strconv.Itoa(i),
			Name:      "Benchmark Create User",
			Email:     "benchcreate" + strconv.Itoa(i) + "@example.com",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		err := service.Create(ctx, user)
		if err != nil {
			b.Fatalf("Create failed: %v", err)
		}
	}
}

func BenchmarkUserService_Update(b *testing.B) {
	// Disable logging for benchmarks
	DisableLoggingForBenchmark(b)

	// Setup
	ctx := context.Background()
	repo := repository.NewMockUserRepository()
	service := NewUserService(repo)

	// Create test user
	user := &domain.User{
		ID:        "bench-update",
		Name:      "Benchmark Update User",
		Email:     "benchupdate@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := repo.Create(ctx, user)
	if err != nil {
		b.Fatalf("Failed to create test user: %v", err)
	}

	// Run benchmark
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a new user with updated fields
		updatedUser := &domain.User{
			ID:        "bench-update",
			Name:      "Updated Name " + strconv.Itoa(i),
			Email:     "updated" + strconv.Itoa(i) + "@example.com",
			CreatedAt: user.CreatedAt,
			UpdatedAt: time.Now(),
		}

		err := service.Update(ctx, updatedUser)
		if err != nil {
			b.Fatalf("Update failed: %v", err)
		}
	}
}

func BenchmarkUserService_Delete(b *testing.B) {
	b.Run("Single user", func(b *testing.B) {
		// Disable logging for benchmarks
		DisableLoggingForBenchmark(b)

		// This is a simpler benchmark that recreates and deletes a single user repeatedly
		ctx := context.Background()
		repo := repository.NewMockUserRepository()
		service := NewUserService(repo)

		// Run benchmark
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// First create a user
			user := &domain.User{
				ID:        "bench-delete",
				Name:      "Benchmark Delete User",
				Email:     "benchdelete@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			err := repo.Create(ctx, user)
			if err != nil {
				b.Fatalf("Failed to create test user: %v", err)
			}

			// Then delete it
			err = service.Delete(ctx, "bench-delete")
			if err != nil {
				b.Fatalf("Delete failed: %v", err)
			}
		}
	})

	b.Run("Multiple users", func(b *testing.B) {
		// Disable logging for benchmarks
		DisableLoggingForBenchmark(b)

		// Setup
		ctx := context.Background()
		repo := repository.NewMockUserRepository()
		service := NewUserService(repo)

		// Create many users before starting the benchmark
		for i := 0; i < b.N; i++ {
			userId := "bench-delete-multi-" + strconv.Itoa(i)
			user := &domain.User{
				ID:        userId,
				Name:      "Benchmark Delete User " + strconv.Itoa(i),
				Email:     "benchdelete" + strconv.Itoa(i) + "@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			err := repo.Create(ctx, user)
			if err != nil {
				b.Fatalf("Failed to create test user: %v", err)
			}
		}

		// Run benchmark
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			userId := "bench-delete-multi-" + strconv.Itoa(i)
			err := service.Delete(ctx, userId)
			if err != nil {
				b.Fatalf("Delete failed: %v", err)
			}
		}
	})
}
