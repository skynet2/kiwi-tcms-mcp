package kiwi

import "context"

func (c *Client) CategoryFilter(ctx context.Context, q map[string]any) ([]Category, error) {
	var out []Category
	err := c.Do(ctx, "Category.filter", []any{q}, &out)
	return out, err
}

func (c *Client) PriorityFilter(ctx context.Context, q map[string]any) ([]Priority, error) {
	var out []Priority
	err := c.Do(ctx, "Priority.filter", []any{q}, &out)
	return out, err
}

func (c *Client) ProductFilter(ctx context.Context, q map[string]any) ([]Product, error) {
	var out []Product
	err := c.Do(ctx, "Product.filter", []any{q}, &out)
	return out, err
}

func (c *Client) ComponentFilter(ctx context.Context, q map[string]any) ([]Component, error) {
	var out []Component
	err := c.Do(ctx, "Component.filter", []any{q}, &out)
	return out, err
}

func (c *Client) TagFilter(ctx context.Context, q map[string]any) ([]Tag, error) {
	var out []Tag
	err := c.Do(ctx, "Tag.filter", []any{q}, &out)
	return out, err
}

func (c *Client) TestCaseStatusFilter(ctx context.Context, q map[string]any) ([]TestCaseStatus, error) {
	var out []TestCaseStatus
	err := c.Do(ctx, "TestCaseStatus.filter", []any{q}, &out)
	return out, err
}

func (c *Client) UserFilter(ctx context.Context, q map[string]any) ([]User, error) {
	var out []User
	err := c.Do(ctx, "User.filter", []any{q}, &out)
	return out, err
}
