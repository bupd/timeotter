package calendar

import (
	"testing"
)

func TestConvertTimeToCron(t *testing.T) {
	tests := []struct {
		name                 string
		timeStr              string
		triggerBeforeMinutes int
		expected             string
	}{
		{
			name:                 "RFC3339 with positive offset +05:30",
			timeStr:              "2025-02-02T20:30:00+05:30",
			triggerBeforeMinutes: 5,
			expected:             "25 20 2 2 0", // Sun Feb 2, 2025 - 5 mins before 20:30 = 20:25
		},
		{
			name:                 "RFC3339 with positive offset +00:00",
			timeStr:              "2025-01-15T10:00:00+00:00",
			triggerBeforeMinutes: 5,
			expected:             "55 9 15 1 3", // Wed Jan 15, 2025 - 5 mins before 10:00 = 9:55
		},
		{
			name:                 "RFC3339 with negative offset -08:00",
			timeStr:              "2025-03-10T14:30:00-08:00",
			triggerBeforeMinutes: 5,
			expected:             "25 14 10 3 1", // Mon Mar 10, 2025 - 5 mins before 14:30 = 14:25
		},
		{
			name:                 "zero trigger offset",
			timeStr:              "2025-06-20T12:00:00+05:30",
			triggerBeforeMinutes: 0,
			expected:             "0 12 20 6 5", // Fri Jun 20, 2025 - no offset = 12:00
		},
		{
			name:                 "10 minute trigger offset",
			timeStr:              "2025-04-15T09:15:00+05:30",
			triggerBeforeMinutes: 10,
			expected:             "5 9 15 4 2", // Tue Apr 15, 2025 - 10 mins before 9:15 = 9:05
		},
		{
			name:                 "date with slashes format",
			timeStr:              "2025/07/04T18:00:00+05:30",
			triggerBeforeMinutes: 5,
			expected:             "55 17 4 7 5", // Fri Jul 4, 2025 - 5 mins before 18:00 = 17:55
		},
		{
			name:                 "midnight crossing - subtract from 00:05",
			timeStr:              "2025-08-10T00:05:00+05:30",
			triggerBeforeMinutes: 10,
			expected:             "55 23 9 8 6", // subtracting 10 mins from 00:05 crosses to previous day
		},
		{
			name:                 "exact midnight",
			timeStr:              "2025-09-01T00:00:00+05:30",
			triggerBeforeMinutes: 5,
			expected:             "55 23 31 8 0", // 5 mins before midnight = 23:55 previous day
		},
		{
			name:                 "end of month crossing",
			timeStr:              "2025-02-01T00:03:00+05:30",
			triggerBeforeMinutes: 5,
			expected:             "58 23 31 1 5", // crosses to Jan 31, 2025
		},
		{
			name:                 "large trigger offset 30 minutes",
			timeStr:              "2025-05-15T10:20:00+05:30",
			triggerBeforeMinutes: 30,
			expected:             "50 9 15 5 4", // Thu May 15, 2025 - 30 mins before 10:20 = 9:50
		},
		{
			name:                 "end of year crossing",
			timeStr:              "2025-01-01T00:02:00+05:30",
			triggerBeforeMinutes: 5,
			expected:             "57 23 31 12 2", // crosses to Dec 31, 2024
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertTimeToCron(tt.timeStr, tt.triggerBeforeMinutes)
			if result != tt.expected {
				t.Errorf("ConvertTimeToCron(%q, %d) = %q, want %q",
					tt.timeStr, tt.triggerBeforeMinutes, result, tt.expected)
			}
		})
	}
}

func TestConvertTimeToCron_CronFormat(t *testing.T) {
	// Test that output is valid cron format: MIN HOUR DAY MONTH WEEKDAY
	result := ConvertTimeToCron("2025-03-15T14:30:00+05:30", 5)

	// Should have 5 space-separated values
	parts := splitCronParts(result)
	if len(parts) != 5 {
		t.Errorf("expected 5 cron parts, got %d: %q", len(parts), result)
	}
}

func TestConvertTimeToCron_MinuteRange(t *testing.T) {
	// Test various times to ensure minute is in valid range 0-59
	times := []string{
		"2025-01-01T00:00:00+05:30",
		"2025-01-01T00:30:00+05:30",
		"2025-01-01T00:59:00+05:30",
	}

	for _, timeStr := range times {
		result := ConvertTimeToCron(timeStr, 0)
		parts := splitCronParts(result)
		if len(parts) < 1 {
			t.Errorf("invalid cron output for %s", timeStr)
			continue
		}
		min := parseIntOrNeg1(parts[0])
		if min < 0 || min > 59 {
			t.Errorf("minute %d out of range for time %s", min, timeStr)
		}
	}
}

func TestConvertTimeToCron_HourRange(t *testing.T) {
	// Test various times to ensure hour is in valid range 0-23
	times := []string{
		"2025-01-01T00:30:00+05:30",
		"2025-01-01T12:30:00+05:30",
		"2025-01-01T23:30:00+05:30",
	}

	for _, timeStr := range times {
		result := ConvertTimeToCron(timeStr, 0)
		parts := splitCronParts(result)
		if len(parts) < 2 {
			t.Errorf("invalid cron output for %s", timeStr)
			continue
		}
		hour := parseIntOrNeg1(parts[1])
		if hour < 0 || hour > 23 {
			t.Errorf("hour %d out of range for time %s", hour, timeStr)
		}
	}
}

func TestConvertTimeToCron_WeekdayRange(t *testing.T) {
	// Test that weekday is in valid range 0-6
	times := []string{
		"2025-01-05T12:00:00+05:30", // Sunday
		"2025-01-06T12:00:00+05:30", // Monday
		"2025-01-11T12:00:00+05:30", // Saturday
	}

	for _, timeStr := range times {
		result := ConvertTimeToCron(timeStr, 0)
		parts := splitCronParts(result)
		if len(parts) < 5 {
			t.Errorf("invalid cron output for %s", timeStr)
			continue
		}
		weekday := parseIntOrNeg1(parts[4])
		if weekday < 0 || weekday > 6 {
			t.Errorf("weekday %d out of range for time %s", weekday, timeStr)
		}
	}
}

// Helper function to split cron parts
func splitCronParts(cron string) []string {
	var parts []string
	current := ""
	for _, c := range cron {
		if c == ' ' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// Helper to parse int or return -1
func parseIntOrNeg1(s string) int {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		} else {
			return -1
		}
	}
	return result
}

// Mock-based tests using testutil

func TestConvertTimeToCron_WithMockEvents(t *testing.T) {
	// Test conversion for typical calendar event times
	testCases := []struct {
		name     string
		dateTime string
		offset   int
		wantMin  int
		wantHour int
	}{
		{
			name:     "morning meeting",
			dateTime: "2025-03-15T09:00:00+05:30",
			offset:   5,
			wantMin:  55,
			wantHour: 8,
		},
		{
			name:     "afternoon standup",
			dateTime: "2025-03-15T14:30:00+05:30",
			offset:   5,
			wantMin:  25,
			wantHour: 14,
		},
		{
			name:     "evening reminder",
			dateTime: "2025-03-15T18:00:00+05:30",
			offset:   10,
			wantMin:  50,
			wantHour: 17,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ConvertTimeToCron(tc.dateTime, tc.offset)
			parts := splitCronParts(result)

			if len(parts) != 5 {
				t.Fatalf("expected 5 parts, got %d", len(parts))
			}

			gotMin := parseIntOrNeg1(parts[0])
			gotHour := parseIntOrNeg1(parts[1])

			if gotMin != tc.wantMin {
				t.Errorf("minute: got %d, want %d", gotMin, tc.wantMin)
			}
			if gotHour != tc.wantHour {
				t.Errorf("hour: got %d, want %d", gotHour, tc.wantHour)
			}
		})
	}
}

func TestConvertTimeToCron_AllDayEvent(t *testing.T) {
	// All-day events have Date instead of DateTime
	// The current implementation would fail to parse date-only strings
	// This test documents the expected behavior

	// Note: The current implementation only handles DateTime format
	// Date-only events (e.g., "2025-03-15") would not parse correctly
	// This is a known limitation that should be documented

	dateOnly := "2025-03-15"
	result := ConvertTimeToCron(dateOnly, 0)

	// The current implementation will produce unexpected results for date-only
	// This test verifies the actual behavior (not necessarily correct)
	parts := splitCronParts(result)
	if len(parts) != 5 {
		t.Logf("date-only parsing produced %d parts: %q", len(parts), result)
	}
}

func TestConvertTimeToCron_DifferentTimezones(t *testing.T) {
	// Test that timezone is preserved (not converted)
	testCases := []struct {
		name     string
		dateTime string
		offset   int
	}{
		{
			name:     "UTC timezone",
			dateTime: "2025-03-15T12:00:00+00:00",
			offset:   5,
		},
		{
			name:     "IST timezone",
			dateTime: "2025-03-15T12:00:00+05:30",
			offset:   5,
		},
		{
			name:     "PST timezone",
			dateTime: "2025-03-15T12:00:00-08:00",
			offset:   5,
		},
		{
			name:     "JST timezone",
			dateTime: "2025-03-15T12:00:00+09:00",
			offset:   5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ConvertTimeToCron(tc.dateTime, tc.offset)
			parts := splitCronParts(result)

			if len(parts) != 5 {
				t.Fatalf("expected 5 parts, got %d", len(parts))
			}

			// All should produce minute 55 and hour 11 (12:00 - 5 min)
			gotMin := parseIntOrNeg1(parts[0])
			gotHour := parseIntOrNeg1(parts[1])

			if gotMin != 55 {
				t.Errorf("minute: got %d, want 55", gotMin)
			}
			if gotHour != 11 {
				t.Errorf("hour: got %d, want 11", gotHour)
			}
		})
	}
}
