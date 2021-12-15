package common

import (
	"github.com/rs/zerolog"
	"github.com/schollz/progressbar/v3"
)

// ProgressBar is a wrapper around the progressbar.ProgressBar type
// instead of returning errors from functions, it only logs errors
type ProgressBar struct {
	ProgressBar *progressbar.ProgressBar
	Log         zerolog.Logger
}

func NewProgressBar(log zerolog.Logger, max int64, description ...string) *ProgressBar {
	return &ProgressBar{
		ProgressBar: progressbar.Default(max, description...),
		Log:         log.With().Str("component", "progressbar").Logger(),
	}
}

func (p *ProgressBar) Add(n int) {
	err := p.ProgressBar.Add(n)
	if err != nil {
		p.Log.Error().Err(err).Msg("progress bar error")
	}
}

func (p *ProgressBar) Increment() {
	p.Add(1)
}

func (p *ProgressBar) Finish() {
	err := p.ProgressBar.Finish()
	if err != nil {
		p.Log.Error().Err(err).Msg("progress bar error")
	}
}
