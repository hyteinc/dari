package dari

import (
	"context"
	"errors"
	"fmt"
)

// GetPut reads an item, applies mutator to the freshly-read copy, and writes
// it back. If the put fails because of a version constraint (another writer
// updated the item between the Get and Put), GetPut re-reads the item and
// retries, up to maxRetries attempts.
//
// Warning: mutator may be called multiple times — once per attempt — so it
// must be safe to re-run against a freshly-read copy each time. If mutator
// returns a non-nil error, GetPut aborts immediately without retrying.
func GetPut[T Keys](ctx context.Context, t *Table, item *T, mutator func(*T) error, maxRetries int) error {
	if item == nil {
		return errors.New("nil item")
	}
	if mutator == nil {
		return errors.New("nil mutator")
	}
	if maxRetries <= 0 {
		return errors.New("maxRetries must be > 0")
	}

	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		current, err := Get[T](ctx, t, any(item).(Keys))
		if err != nil {
			return err
		}

		if err := mutator(current); err != nil {
			return err
		}

		err = Put[T](ctx, t, current)
		if err == nil {
			return nil
		}
		if errors.Is(err, ErrAlreadyExists) {
			lastErr = err
			continue
		}
		return err
	}

	return fmt.Errorf("GetPut exhausted %d retries: %w", maxRetries, lastErr)
}
