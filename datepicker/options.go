package datepicker

import "time"

type Option func(dp *DatePicker)

// StartFromSunday sets the first day of the week to Sunday.
func StartFromSunday() Option {
	return func(dp *DatePicker) {
		dp.startFromSunday = true
	}
}

// CurrentDate sets the current date.
func CurrentDate(current time.Time) Option {
	return func(dp *DatePicker) {
		dp.month = current.Month()
		dp.year = current.Year()
	}
}

// OnCancel sets the callback function for the cancel button.
func OnCancel(f OnCancelHandler) Option {
	return func(dp *DatePicker) {
		dp.onCancel = f
	}
}

// OnError sets the callback function for the error.
func OnError(f OnErrorHandler) Option {
	return func(dp *DatePicker) {
		dp.onError = f
	}
}

// Language sets the language of the datepicker.
func Language(lang string) Option {
	return func(dp *DatePicker) {
		dp.language = lang
	}
}

// Languages sets the languages of the datepicker.
// All supported keys you can see in langs.json file
func Languages(langs LangsData) Option {
	return func(dp *DatePicker) {
		dp.langs = langs
	}
}

// From sets the minimum date.
func From(from time.Time) Option {
	return func(dp *DatePicker) {
		dp.from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	}
}

// To sets the maximum date.
func To(to time.Time) Option {
	return func(dp *DatePicker) {
		dp.to = time.Date(to.Year(), to.Month(), to.Day(), 0, 0, 0, 0, to.Location())
	}
}

// Dates sets the dates. And mode for include/exclude.
func Dates(datesMode DatesMode, dates []time.Time) Option {
	return func(dp *DatePicker) {
		for _, d := range dates {
			dp.dates = append(dp.dates, time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location()))
		}
		dp.datesMode = datesMode
	}
}

// WithPrefix sets a prefix for the widget
func WithPrefix(s string) Option {
	return func(w *DatePicker) {
		w.prefix = s
	}
}

// WithCallbackPrefix sets a specific callback prefix for date picker interactions
func WithCallbackPrefix(prefix string) Option {
	return func(w *DatePicker) {
		w.prefix = prefix
	}
}

// NoDeleteAfterSelect prevents remove keyboard after select
func NoDeleteAfterSelect() Option {
	return func(dp *DatePicker) {
		dp.deleteOnSelect = false
	}
}

// NoDeleteAfterCancel prevents remove keyboard after cancel
func NoDeleteAfterCancel() Option {
	return func(dp *DatePicker) {
		dp.deleteOnCancel = false
	}
}

// ForbidPastDates ensures users cannot pick dates earlier than today.
// It sets the minimum selectable date (from) to the start of today,
// unless a later minimum was already configured via From().
func ForbidPastDates() Option {
	return func(dp *DatePicker) {
		now := time.Now()
		startOfToday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if dp.from.IsZero() || dp.from.Before(startOfToday) {
			dp.from = startOfToday
		}
	}
}
