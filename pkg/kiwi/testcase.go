package kiwi

import "context"

func (c *Client) TestCaseFilter(ctx context.Context, query map[string]any) ([]TestCase, error) {
	var out []TestCase
	err := c.Do(ctx, "TestCase.filter", []any{query}, &out)
	return out, err
}

func (c *Client) TestCaseCreate(ctx context.Context, in TestCaseCreateInput) (TestCase, error) {
	var out TestCase
	err := c.Do(ctx, "TestCase.create", []any{in}, &out)
	return out, err
}

func (c *Client) TestCaseUpdate(ctx context.Context, id int64, patch map[string]any) (TestCase, error) {
	var out TestCase
	err := c.Do(ctx, "TestCase.update", []any{id, patch}, &out)
	return out, err
}

func (c *Client) TestCaseRemove(ctx context.Context, query map[string]any) error {
	return c.Do(ctx, "TestCase.remove", []any{query}, nil)
}

func (c *Client) TestCaseAddTag(ctx context.Context, caseID int64, tag string) error {
	return c.Do(ctx, "TestCase.add_tag", []any{caseID, tag}, nil)
}

func (c *Client) TestCaseRemoveTag(ctx context.Context, caseID int64, tag string) error {
	return c.Do(ctx, "TestCase.remove_tag", []any{caseID, tag}, nil)
}

func (c *Client) TestCaseAddComponent(ctx context.Context, caseID int64, component string) error {
	return c.Do(ctx, "TestCase.add_component", []any{caseID, component}, nil)
}

func (c *Client) TestCaseRemoveComponent(ctx context.Context, caseID int64, componentID int64) error {
	return c.Do(ctx, "TestCase.remove_component", []any{caseID, componentID}, nil)
}

func (c *Client) TestCaseAddComment(ctx context.Context, caseID int64, comment string) error {
	return c.Do(ctx, "TestCase.add_comment", []any{caseID, comment}, nil)
}

func (c *Client) TestCaseAddProperty(ctx context.Context, caseID int64, name, value string) error {
	return c.Do(ctx, "TestCase.add_property", []any{caseID, name, value}, nil)
}

func (c *Client) TestCaseRemoveProperty(ctx context.Context, query map[string]any) error {
	return c.Do(ctx, "TestCase.remove_property", []any{query}, nil)
}

func (c *Client) TestCaseAddLink(ctx context.Context, link map[string]any) (map[string]any, error) {
	out := map[string]any{}
	err := c.Do(ctx, "TestCase.add_link", []any{link}, &out)
	return out, err
}

func (c *Client) TestCaseRemoveLink(ctx context.Context, query map[string]any) error {
	return c.Do(ctx, "TestCase.remove_link", []any{query}, nil)
}
