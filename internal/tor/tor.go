package tor

import (
	"context"
	"time"

	"github.com/cretz/bine/tor"
)

// StartTor initializes and starts a new Tor client.
func StartTor() (*tor.Tor, error) {
	t, err := tor.Start(nil, &tor.StartConf{ProcessTimeout: 3 * time.Minute})
	if err != nil {
		return nil, err
	}
	// Wait until Tor is fully bootstrapped.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	if err := t.Wait(ctx); err != nil {
		return nil, err
	}
	return t, nil
}
