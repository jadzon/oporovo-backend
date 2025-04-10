package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
	"vibely-backend/src/models"
	"vibely-backend/src/repositories"
)

type TutorAvailabilityService interface {
	// Weekly schedule methods
	CreateWeeklySchedule(tutorID uuid.UUID, dayOfWeek int, startTime, endTime string) (models.TutorWeeklySchedule, error)
	GetWeeklySchedulesByTutorID(tutorID uuid.UUID) ([]models.TutorWeeklySchedule, error)
	UpdateWeeklySchedule(scheduleID uuid.UUID, dayOfWeek int, startTime, endTime string) (models.TutorWeeklySchedule, error)
	DeleteWeeklySchedule(scheduleID uuid.UUID) error

	// Exception methods
	AddException(tutorID uuid.UUID, date time.Time, startTime, endTime string, isRemoval bool) (models.TutorScheduleException, error)
	GetExceptionsByTutorID(tutorID uuid.UUID, startDate, endDate time.Time) ([]models.TutorScheduleException, error)
	UpdateException(exceptionID uuid.UUID, date time.Time, startTime, endTime string, isRemoval bool) (models.TutorScheduleException, error)
	DeleteException(exceptionID uuid.UUID) error

	// Availability calculation methods
	GetAvailabilityForDateRange(tutorID uuid.UUID, startDate, endDate time.Time) ([]models.AvailabilitySlot, error)
}

type tutorAvailabilityService struct {
	availabilityRepo repositories.TutorAvailabilityRepository
	userRepo         repositories.UserRepository
}

func NewTutorAvailabilityService(availabilityRepo repositories.TutorAvailabilityRepository, userRepo repositories.UserRepository) TutorAvailabilityService {
	return &tutorAvailabilityService{
		availabilityRepo: availabilityRepo,
		userRepo:         userRepo,
	}
}

// CreateWeeklySchedule adds a new recurring time slot to a tutor's weekly schedule
func (s *tutorAvailabilityService) CreateWeeklySchedule(tutorID uuid.UUID, dayOfWeek int, startTime, endTime string) (models.TutorWeeklySchedule, error) {
	// Validate day of week (0-6)
	if dayOfWeek < 0 || dayOfWeek > 6 {
		return models.TutorWeeklySchedule{}, errors.New("day of week must be between 0 (Sunday) and 6 (Saturday)")
	}

	// Validate time format (24-hour)
	_, err := time.Parse("15:04", startTime)
	if err != nil {
		return models.TutorWeeklySchedule{}, errors.New("invalid start time format, use 24-hour format (e.g., 14:30)")
	}

	_, err = time.Parse("15:04", endTime)
	if err != nil {
		return models.TutorWeeklySchedule{}, errors.New("invalid end time format, use 24-hour format (e.g., 16:45)")
	}

	// Ensure startTime is before endTime
	if startTime >= endTime {
		return models.TutorWeeklySchedule{}, errors.New("start time must be before end time")
	}

	// Verify that the user exists and is a tutor
	user, err := s.userRepo.GetUserByID(tutorID)
	if err != nil {
		return models.TutorWeeklySchedule{}, errors.New("tutor not found")
	}

	if user.Role != models.UserRoleTutor {
		return models.TutorWeeklySchedule{}, errors.New("user is not a tutor")
	}

	// Create schedule in repository
	schedule := models.TutorWeeklySchedule{
		TutorID:   tutorID,
		DayOfWeek: dayOfWeek,
		StartTime: startTime,
		EndTime:   endTime,
	}

	return s.availabilityRepo.CreateWeeklySchedule(schedule)
}

// GetWeeklySchedulesByTutorID retrieves all weekly schedule slots for a tutor
func (s *tutorAvailabilityService) GetWeeklySchedulesByTutorID(tutorID uuid.UUID) ([]models.TutorWeeklySchedule, error) {
	// Verify that the user exists and is a tutor
	user, err := s.userRepo.GetUserByID(tutorID)
	if err != nil {
		return nil, errors.New("tutor not found")
	}

	if user.Role != models.UserRoleTutor {
		return nil, errors.New("user is not a tutor")
	}

	return s.availabilityRepo.GetWeeklySchedulesByTutorID(tutorID)
}

// UpdateWeeklySchedule updates an existing weekly schedule slot
func (s *tutorAvailabilityService) UpdateWeeklySchedule(scheduleID uuid.UUID, dayOfWeek int, startTime, endTime string) (models.TutorWeeklySchedule, error) {
	// Validate day of week (0-6)
	if dayOfWeek < 0 || dayOfWeek > 6 {
		return models.TutorWeeklySchedule{}, errors.New("day of week must be between 0 (Sunday) and 6 (Saturday)")
	}

	// Validate time format (24-hour)
	_, err := time.Parse("15:04", startTime)
	if err != nil {
		return models.TutorWeeklySchedule{}, errors.New("invalid start time format, use 24-hour format (e.g., 14:30)")
	}

	_, err = time.Parse("15:04", endTime)
	if err != nil {
		return models.TutorWeeklySchedule{}, errors.New("invalid end time format, use 24-hour format (e.g., 16:45)")
	}

	// Ensure startTime is before endTime
	if startTime >= endTime {
		return models.TutorWeeklySchedule{}, errors.New("start time must be before end time")
	}

	// Get existing schedule
	schedule, err := s.availabilityRepo.GetWeeklyScheduleByID(scheduleID)
	if err != nil {
		return models.TutorWeeklySchedule{}, errors.New("schedule not found")
	}

	// Update fields
	schedule.DayOfWeek = dayOfWeek
	schedule.StartTime = startTime
	schedule.EndTime = endTime

	// Save to repository
	err = s.availabilityRepo.UpdateWeeklySchedule(schedule)
	if err != nil {
		return models.TutorWeeklySchedule{}, err
	}

	return schedule, nil
}

// DeleteWeeklySchedule deletes a weekly schedule slot
func (s *tutorAvailabilityService) DeleteWeeklySchedule(scheduleID uuid.UUID) error {
	// Verify schedule exists
	_, err := s.availabilityRepo.GetWeeklyScheduleByID(scheduleID)
	if err != nil {
		return errors.New("schedule not found")
	}

	return s.availabilityRepo.DeleteWeeklySchedule(scheduleID)
}

// AddException adds a new exception to the tutor's schedule
func (s *tutorAvailabilityService) AddException(tutorID uuid.UUID, date time.Time, startTime, endTime string, isRemoval bool) (models.TutorScheduleException, error) {
	// Verify that the user exists and is a tutor
	user, err := s.userRepo.GetUserByID(tutorID)
	if err != nil {
		return models.TutorScheduleException{}, errors.New("tutor not found")
	}

	if user.Role != models.UserRoleTutor {
		return models.TutorScheduleException{}, errors.New("user is not a tutor")
	}

	// If it's not a removal, validate the time range
	if !isRemoval {
		// Validate time format (24-hour)
		_, err := time.Parse("15:04", startTime)
		if err != nil {
			return models.TutorScheduleException{}, errors.New("invalid start time format, use 24-hour format (e.g., 14:30)")
		}

		_, err = time.Parse("15:04", endTime)
		if err != nil {
			return models.TutorScheduleException{}, errors.New("invalid end time format, use 24-hour format (e.g., 16:45)")
		}

		// Ensure startTime is before endTime
		if startTime >= endTime {
			return models.TutorScheduleException{}, errors.New("start time must be before end time")
		}
	}

	// Create exception in repository
	exception := models.TutorScheduleException{
		TutorID:   tutorID,
		Date:      date,
		StartTime: startTime,
		EndTime:   endTime,
		IsRemoval: isRemoval,
	}

	return s.availabilityRepo.CreateException(exception)
}

// GetExceptionsByTutorID retrieves all exceptions for a tutor in a date range
func (s *tutorAvailabilityService) GetExceptionsByTutorID(tutorID uuid.UUID, startDate, endDate time.Time) ([]models.TutorScheduleException, error) {
	// Verify that the user exists and is a tutor
	user, err := s.userRepo.GetUserByID(tutorID)
	if err != nil {
		return nil, errors.New("tutor not found")
	}

	if user.Role != models.UserRoleTutor {
		return nil, errors.New("user is not a tutor")
	}

	return s.availabilityRepo.GetExceptionsByTutorID(tutorID, startDate, endDate)
}

// UpdateException updates an existing schedule exception
func (s *tutorAvailabilityService) UpdateException(exceptionID uuid.UUID, date time.Time, startTime, endTime string, isRemoval bool) (models.TutorScheduleException, error) {
	// Get existing exception
	exception, err := s.availabilityRepo.GetExceptionByID(exceptionID)
	if err != nil {
		return models.TutorScheduleException{}, errors.New("exception not found")
	}

	// If it's not a removal, validate the time range
	if !isRemoval {
		// Validate time format (24-hour)
		_, err := time.Parse("15:04", startTime)
		if err != nil {
			return models.TutorScheduleException{}, errors.New("invalid start time format, use 24-hour format (e.g., 14:30)")
		}

		_, err = time.Parse("15:04", endTime)
		if err != nil {
			return models.TutorScheduleException{}, errors.New("invalid end time format, use 24-hour format (e.g., 16:45)")
		}

		// Ensure startTime is before endTime
		if startTime >= endTime {
			return models.TutorScheduleException{}, errors.New("start time must be before end time")
		}
	}

	// Update fields
	exception.Date = date
	exception.StartTime = startTime
	exception.EndTime = endTime
	exception.IsRemoval = isRemoval

	// Save to repository
	err = s.availabilityRepo.UpdateException(exception)
	if err != nil {
		return models.TutorScheduleException{}, err
	}

	return exception, nil
}

// DeleteException deletes a schedule exception
func (s *tutorAvailabilityService) DeleteException(exceptionID uuid.UUID) error {
	// Verify exception exists
	_, err := s.availabilityRepo.GetExceptionByID(exceptionID)
	if err != nil {
		return errors.New("exception not found")
	}

	return s.availabilityRepo.DeleteException(exceptionID)
}

// GetAvailabilityForDateRange calculates a tutor's availability for a date range
func (s *tutorAvailabilityService) GetAvailabilityForDateRange(tutorID uuid.UUID, startDate, endDate time.Time) ([]models.AvailabilitySlot, error) {
	// Verify that the user exists and is a tutor
	user, err := s.userRepo.GetUserByID(tutorID)
	if err != nil {
		return nil, errors.New("tutor not found")
	}

	if user.Role != models.UserRoleTutor {
		return nil, errors.New("user is not a tutor")
	}

	// Get the tutor's weekly schedule
	weeklySchedules, err := s.availabilityRepo.GetWeeklySchedulesByTutorID(tutorID)
	if err != nil {
		return nil, err
	}

	// Get exceptions for the date range
	exceptions, err := s.availabilityRepo.GetExceptionsByTutorID(tutorID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	fmt.Printf("\n==== AVAILABILITY CALCULATION ====\n")
	fmt.Printf("Tutor ID: %s\n", tutorID)
	fmt.Printf("Date range: %s to %s\n", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	fmt.Printf("Found %d weekly schedules and %d exceptions\n", len(weeklySchedules), len(exceptions))

	// Calculate available slots for each day in the range
	var availableSlots []models.AvailabilitySlot

	// Create a map of exceptions by date for faster lookup
	exceptionsByDate := make(map[time.Time][]models.TutorScheduleException)
	for _, exception := range exceptions {
		dateKey := time.Date(exception.Date.Year(), exception.Date.Month(), exception.Date.Day(), 0, 0, 0, 0, exception.Date.Location())
		exceptionsByDate[dateKey] = append(exceptionsByDate[dateKey], exception)
		fmt.Printf("Exception for %s: IsRemoval=%v, Time=%s-%s\n",
			dateKey.Format("2006-01-02"), exception.IsRemoval, exception.StartTime, exception.EndTime)
	}

	// Iterate through each day in the date range
	for date := startDate; date.Before(endDate) || date.Equal(endDate); date = date.AddDate(0, 0, 1) {
		dayOfWeek := int(date.Weekday())
		dateKey := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

		fmt.Printf("\nProcessing day: %s (day of week: %d)\n", dateKey.Format("2006-01-02"), dayOfWeek)

		// Check if there are any full-day removals
		hasFullDayRemoval := false
		for _, exception := range exceptionsByDate[dateKey] {
			if exception.IsRemoval && exception.StartTime == "" && exception.EndTime == "" {
				hasFullDayRemoval = true
				fmt.Printf("  Full day removal exception found - skipping this day\n")
				break
			}
		}

		if hasFullDayRemoval {
			continue // Skip this day entirely
		}

		// Get regular schedule for this day of the week
		var regularSlots []struct {
			StartTime string
			EndTime   string
		}

		for _, schedule := range weeklySchedules {
			if schedule.DayOfWeek == dayOfWeek {
				regularSlots = append(regularSlots, struct {
					StartTime string
					EndTime   string
				}{
					StartTime: schedule.StartTime,
					EndTime:   schedule.EndTime,
				})
				fmt.Printf("  Regular slot for this day: %s-%s\n", schedule.StartTime, schedule.EndTime)
			}
		}

		// Apply exceptions to the regular slots
		finalSlots := make(map[string]string) // map of start time to end time

		// Add all regular slots to the map
		for _, slot := range regularSlots {
			finalSlots[slot.StartTime] = slot.EndTime
			fmt.Printf("  Added regular slot: %s-%s\n", slot.StartTime, slot.EndTime)
		}

		// Process exceptions for this day
		if exceptions, ok := exceptionsByDate[dateKey]; ok {
			fmt.Printf("  Found %d exceptions for this day\n", len(exceptions))
			for _, exception := range exceptions {
				if exception.IsRemoval {
					// Remove specific time slot
					if exception.StartTime != "" && exception.EndTime != "" {
						if _, exists := finalSlots[exception.StartTime]; exists {
							fmt.Printf("  Removing slot due to exception: %s-%s\n",
								exception.StartTime, exception.EndTime)
							delete(finalSlots, exception.StartTime)
						} else {
							fmt.Printf("  Tried to remove slot %s-%s but it didn't exist\n",
								exception.StartTime, exception.EndTime)
						}
					}
				} else {
					// Add additional time slot
					finalSlots[exception.StartTime] = exception.EndTime
					fmt.Printf("  Added slot from exception: %s-%s\n",
						exception.StartTime, exception.EndTime)
				}
			}
		} else {
			fmt.Printf("  No exceptions for this day\n")
		}

		// Convert the map to AvailabilitySlot array
		if len(finalSlots) > 0 {
			fmt.Printf("  Final slots for %s:\n", dateKey.Format("2006-01-02"))
			for startTime, endTime := range finalSlots {
				availableSlots = append(availableSlots, models.AvailabilitySlot{
					Date:      date,
					StartTime: startTime,
					EndTime:   endTime,
				})
				fmt.Printf("    %s-%s\n", startTime, endTime)
			}
		} else {
			fmt.Printf("  No available slots for this day\n")
		}
	}

	fmt.Printf("\nTotal available slots: %d\n", len(availableSlots))
	fmt.Printf("==== END AVAILABILITY CALCULATION ====\n\n")

	return availableSlots, nil
}
