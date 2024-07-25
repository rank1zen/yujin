package database

import (
	"context"
	"testing"
	"time"
)

// xpzpxnzLQX12ACv3iHZfqgdA8RGZQBLCiqJVa1rfVO8Z3KRiYD7YikD2RZC5mot0YhJNKn1UDxu-Ng
// jpGYaAp070Kd35tVISUaS-tK4ZYoZYfwtGlBUR3sXV1UWVU
func TestProfile(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(testingContext(t), 60*time.Second)
	defer cancel()

	db := setupDB(t)

	db.UpdateInitial(ctx, "xpzpxnzLQX12ACv3iHZfqgdA8RGZQBLCiqJVa1rfVO8Z3KRiYD7YikD2RZC5mot0YhJNKn1UDxu-Ng")
}
