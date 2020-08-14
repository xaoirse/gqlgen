package graphql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestGetOperationContext(t *testing.T) {
	rc := &OperationContext{}

	t.Run("with operation context", func(t *testing.T) {
		ctx := WithOperationContext(context.Background(), rc)

		require.True(t, HasOperationContext(ctx))
		require.Equal(t, rc, GetOperationContext(ctx))
	})

	t.Run("without operation context", func(t *testing.T) {
		ctx := context.Background()

		require.False(t, HasOperationContext(ctx))
		require.Panics(t, func() {
			GetOperationContext(ctx)
		})
	})
}

func TestCollectAllFields(t *testing.T) {
	t.Run("collect fields", func(t *testing.T) {
		ctx := testContext(ast.SelectionSet{
			&ast.Field{
				Name: "field",
			},
		})
		s := CollectAllFields(ctx)
		require.Equal(t, []string{"field"}, s)
	})

	t.Run("unique field names", func(t *testing.T) {
		ctx := testContext(ast.SelectionSet{
			&ast.Field{
				Name: "field",
			},
			&ast.Field{
				Name:  "field",
				Alias: "field alias",
			},
		})
		s := CollectAllFields(ctx)
		require.Equal(t, []string{"field"}, s)
	})

	t.Run("collect fragments", func(t *testing.T) {
		ctx := testContext(ast.SelectionSet{
			&ast.Field{
				Name: "fieldA",
			},
			&ast.InlineFragment{
				TypeCondition: "ExampleTypeA",
				SelectionSet: ast.SelectionSet{
					&ast.Field{
						Name: "fieldA",
					},
				},
			},
			&ast.InlineFragment{
				TypeCondition: "ExampleTypeB",
				SelectionSet: ast.SelectionSet{
					&ast.Field{
						Name: "fieldB",
					},
				},
			},
		})
		s := CollectAllFields(ctx)
		require.Equal(t, []string{"fieldA", "fieldB"}, s)
	})

	t.Run("collect fragments with same field name", func(t *testing.T) {
		ctx := testContext(ast.SelectionSet{
			&ast.InlineFragment{
				TypeCondition: "ExampleTypeA",
				SelectionSet: ast.SelectionSet{
					&ast.Field{
						Name: "fieldA",
						ObjectDefinition:  &ast.Definition{Name:"ExampleTypeA"},
					},
				},
			},
			&ast.InlineFragment{
				TypeCondition: "ExampleTypeB",
				SelectionSet: ast.SelectionSet{
					&ast.Field{
						Name: "fieldA",
						ObjectDefinition:  &ast.Definition{Name:"ExampleTypeB"},
					},
				},
			},
		})

		resctx := GetFieldContext(ctx)
		collected := CollectFields(GetOperationContext(ctx), resctx.Field.Selections, nil)
		s := make([]string, 0, len(collected))
		for _, f := range collected {
			s = append(s, f.Name)
		}
		require.Equal(t, []string{"fieldA", "fieldA"}, s)
	})
}
