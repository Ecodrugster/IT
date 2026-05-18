package handlers

import "testing"

func TestNormalizeRole(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "admin", in: "admin", want: "admin"},
		{name: "teacher uppercase", in: "TEACHER", want: "teacher"},
		{name: "student spaced", in: " student ", want: "student"},
		{name: "invalid fallback", in: "moderator", want: "student"},
		{name: "empty fallback", in: "", want: "student"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeRole(tt.in)
			if got != tt.want {
				t.Fatalf("normalizeRole(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestNormalizeAttendanceStatus(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{name: "present", in: "present", want: "present"},
		{name: "late uppercase", in: "LATE", want: "late"},
		{name: "excused spaced", in: " excused ", want: "excused"},
		{name: "invalid empty", in: "holiday", want: ""},
		{name: "empty", in: "", want: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeAttendanceStatus(tt.in)
			if got != tt.want {
				t.Fatalf("normalizeAttendanceStatus(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestValidateScheduleInput(t *testing.T) {
	valid := ScheduleInput{
		Subject:    "Math",
		GroupName:  "P-21",
		TeacherID:  "teacher-1",
		DayOfWeek:  2,
		PairNumber: 3,
		StartsAt:   "09:00",
		EndsAt:     "10:20",
		Room:       "302",
	}
	if !validateScheduleInput(&valid) {
		t.Fatalf("expected valid schedule input to pass validation")
	}

	badDay := valid
	badDay.DayOfWeek = 8
	if validateScheduleInput(&badDay) {
		t.Fatalf("expected invalid day_of_week to fail validation")
	}

	badPair := valid
	badPair.PairNumber = 0
	if validateScheduleInput(&badPair) {
		t.Fatalf("expected invalid pair_number to fail validation")
	}

	badStart := valid
	badStart.StartsAt = "99:00"
	if validateScheduleInput(&badStart) {
		t.Fatalf("expected invalid starts_at format to fail validation")
	}

	badEnd := valid
	badEnd.EndsAt = "25:10"
	if validateScheduleInput(&badEnd) {
		t.Fatalf("expected invalid ends_at format to fail validation")
	}

	missing := valid
	missing.Subject = " "
	if validateScheduleInput(&missing) {
		t.Fatalf("expected missing subject to fail validation")
	}
}

func TestNormalizeClubInput(t *testing.T) {
	input := ClubInput{
		Name:        "  Chess Club  ",
		Description: "  Learn and play  ",
		Icon:        "",
		Color:       "",
	}

	if !normalizeClubInput(&input) {
		t.Fatalf("expected normalizeClubInput to pass for valid name/description")
	}
	if input.Name != "Chess Club" {
		t.Fatalf("expected trimmed club name, got %q", input.Name)
	}
	if input.Description != "Learn and play" {
		t.Fatalf("expected trimmed club description, got %q", input.Description)
	}
	if input.Icon == "" {
		t.Fatalf("expected default icon to be populated")
	}
	if input.Color == "" {
		t.Fatalf("expected default color to be populated")
	}

	invalid := ClubInput{Name: "", Description: "desc"}
	if normalizeClubInput(&invalid) {
		t.Fatalf("expected invalid club input to fail")
	}
}
