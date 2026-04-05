//go:build e2e

package kiwi_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/skynet2/kiwi-tcms-mcp/pkg/kiwi"
)

func e2eConfig(t *testing.T) kiwi.Config {
	t.Helper()
	if testing.Short() {
		t.Skip("skip E2E in -short mode")
	}
	url := os.Getenv("KIWI_E2E_URL")
	user := os.Getenv("KIWI_E2E_USERNAME")
	pass := os.Getenv("KIWI_E2E_PASSWORD")
	if url == "" || user == "" || pass == "" {
		t.Skip("E2E disabled: set KIWI_E2E_URL, KIWI_E2E_USERNAME, KIWI_E2E_PASSWORD")
	}
	return kiwi.Config{BaseURL: url, Username: user, Password: pass}
}

func TestE2E_LoginAndLookups(t *testing.T) {
	client, err := kiwi.New(e2eConfig(t))
	require.NoError(t, err)
	ctx := context.Background()
	require.NoError(t, client.Login(ctx))

	priorities, err := client.PriorityFilter(ctx, map[string]any{})
	require.NoError(t, err)
	assert.NotEmpty(t, priorities)

	statuses, err := client.TestCaseStatusFilter(ctx, map[string]any{})
	require.NoError(t, err)
	assert.NotEmpty(t, statuses)

	users, err := client.UserFilter(ctx, map[string]any{})
	require.NoError(t, err)
	assert.NotEmpty(t, users)

	products, err := client.ProductFilter(ctx, map[string]any{})
	require.NoError(t, err)
	assert.NotEmpty(t, products)

	cats, err := client.CategoryFilter(ctx, map[string]any{"product": products[0].ID})
	require.NoError(t, err)
	assert.NotEmpty(t, cats)

	comps, err := client.ComponentFilter(ctx, map[string]any{"product": products[0].ID})
	require.NoError(t, err)
	assert.NotEmpty(t, comps)

	_, err = client.TagFilter(ctx, map[string]any{})
	require.NoError(t, err)
}

func TestE2E_TestCaseLifecycle(t *testing.T) {
	client, err := kiwi.New(e2eConfig(t))
	require.NoError(t, err)
	ctx := context.Background()
	require.NoError(t, client.Login(ctx))

	products, err := client.ProductFilter(ctx, map[string]any{})
	require.NoError(t, err)
	require.NotEmpty(t, products)
	productID := products[0].ID

	cats, err := client.CategoryFilter(ctx, map[string]any{"product": productID})
	require.NoError(t, err)
	require.NotEmpty(t, cats)

	priorities, err := client.PriorityFilter(ctx, map[string]any{})
	require.NoError(t, err)
	require.NotEmpty(t, priorities)

	confirmedList, err := client.TestCaseStatusFilter(ctx, map[string]any{"name": "CONFIRMED"})
	require.NoError(t, err)
	require.NotEmpty(t, confirmedList)

	summary := "E2E lifecycle " + time.Now().Format(time.RFC3339Nano)
	created, err := client.TestCaseCreate(ctx, kiwi.TestCaseCreateInput{
		Summary:     summary,
		Category:    cats[0].ID,
		Priority:    priorities[0].ID,
		CaseStatus:  confirmedList[0].ID,
		IsAutomated: true,
	})
	require.NoError(t, err)
	require.NotZero(t, created.ID)

	t.Cleanup(func() {
		_ = client.TestCaseRemove(context.Background(), map[string]any{"pk": created.ID})
	})

	t.Run("filter by id", func(t *testing.T) {
		got, err := client.TestCaseFilter(ctx, map[string]any{"pk": created.ID})
		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, summary, got[0].Summary)
	})

	t.Run("update summary", func(t *testing.T) {
		updated, err := client.TestCaseUpdate(ctx, created.ID, map[string]any{"summary": summary + " (upd)"})
		require.NoError(t, err)
		assert.Equal(t, summary+" (upd)", updated.Summary)
	})

	t.Run("add+remove tag", func(t *testing.T) {
		require.NoError(t, client.TestCaseAddTag(ctx, created.ID, "e2e-smoke"))
		require.NoError(t, client.TestCaseRemoveTag(ctx, created.ID, "e2e-smoke"))
	})

	t.Run("add comment", func(t *testing.T) {
		require.NoError(t, client.TestCaseAddComment(ctx, created.ID, "e2e hello"))
	})

	t.Run("add+remove component", func(t *testing.T) {
		comps, err := client.ComponentFilter(ctx, map[string]any{"product": productID})
		require.NoError(t, err)
		require.NotEmpty(t, comps)
		require.NoError(t, client.TestCaseAddComponent(ctx, created.ID, comps[0].Name))
		require.NoError(t, client.TestCaseRemoveComponent(ctx, created.ID, comps[0].ID))
	})

	t.Run("add+remove property", func(t *testing.T) {
		require.NoError(t, client.TestCaseAddProperty(ctx, created.ID, "env", "e2e"))
		require.NoError(t, client.TestCaseRemoveProperty(ctx, map[string]any{"case": created.ID, "name": "env"}))
	})
}

func TestE2E_KnownBroken_AddLink(t *testing.T) {
	client, err := kiwi.New(e2eConfig(t))
	require.NoError(t, err)
	ctx := context.Background()
	require.NoError(t, client.Login(ctx))

	_, err = client.TestCaseAddLink(ctx, map[string]any{"case_id": 1, "url": "http://x"})
	assert.Error(t, err, "expected Method not found: TestCase.add_link not available on this Kiwi TCMS version")
}

func TestE2E_KnownBroken_RemoveLink(t *testing.T) {
	client, err := kiwi.New(e2eConfig(t))
	require.NoError(t, err)
	ctx := context.Background()
	require.NoError(t, client.Login(ctx))

	err = client.TestCaseRemoveLink(ctx, map[string]any{"pk": 999999})
	assert.Error(t, err, "expected Method not found: TestCase.remove_link not available on this Kiwi TCMS version")
}
