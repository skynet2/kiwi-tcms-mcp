package mcp

import (
	"context"

	"github.com/skynet2/kiwi-tcms-mcp/pkg/kiwi"
)

//go:generate mockgen -source=client.go -destination=mocks/mock_client.go -package=mocks

type KiwiClient interface {
	// TestCase CRUD
	TestCaseFilter(ctx context.Context, q map[string]any) ([]kiwi.TestCase, error)
	TestCaseCreate(ctx context.Context, in kiwi.TestCaseCreateInput) (kiwi.TestCase, error)
	TestCaseUpdate(ctx context.Context, id int64, patch map[string]any) (kiwi.TestCase, error)
	TestCaseRemove(ctx context.Context, q map[string]any) error
	// TestCase relations
	TestCaseAddTag(ctx context.Context, caseID int64, tag string) error
	TestCaseRemoveTag(ctx context.Context, caseID int64, tag string) error
	TestCaseAddComponent(ctx context.Context, caseID int64, component string) error
	TestCaseRemoveComponent(ctx context.Context, caseID int64, componentID int64) error
	TestCaseAddComment(ctx context.Context, caseID int64, comment string) error
	TestCaseAddProperty(ctx context.Context, caseID int64, name, value string) error
	TestCaseRemoveProperty(ctx context.Context, q map[string]any) error
	TestCaseAddLink(ctx context.Context, link map[string]any) (map[string]any, error)
	TestCaseRemoveLink(ctx context.Context, q map[string]any) error
	// Lookups
	CategoryFilter(ctx context.Context, q map[string]any) ([]kiwi.Category, error)
	PriorityFilter(ctx context.Context, q map[string]any) ([]kiwi.Priority, error)
	ProductFilter(ctx context.Context, q map[string]any) ([]kiwi.Product, error)
	ComponentFilter(ctx context.Context, q map[string]any) ([]kiwi.Component, error)
	TagFilter(ctx context.Context, q map[string]any) ([]kiwi.Tag, error)
	TestCaseStatusFilter(ctx context.Context, q map[string]any) ([]kiwi.TestCaseStatus, error)
	UserFilter(ctx context.Context, q map[string]any) ([]kiwi.User, error)
}

// Compile-time check: *kiwi.Client implements KiwiClient
var _ KiwiClient = (*kiwi.Client)(nil)
