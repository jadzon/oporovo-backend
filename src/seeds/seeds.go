package seeds

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq" // make sure this is imported
	"gorm.io/gorm"

	"vibely-backend/src/models"
)

func Seed(db *gorm.DB) error {
	// 1) Retrieve an existing Student with the fixed ID.
	studentUUIDStr := "8e686fcf-3849-4efd-8e22-16280d3d310f"
	studentUUID, err := uuid.Parse(studentUUIDStr)
	if err != nil {
		return fmt.Errorf("failed to parse student UUID: %w", err)
	}

	var student models.User
	if err := db.First(&student, "id = ?", studentUUID).Error; err != nil {
		return fmt.Errorf("failed to find student with ID %s: %w", studentUUIDStr, err)
	}

	// 2) Create multiple Tutors (10 in this example) with richer, realistic details.
	var tutors []models.User
	// Realistic sample data arrays.
	firstNames := []string{"Michał", "Anna", "Paweł", "Katarzyna", "Tomasz", "Magdalena", "Grzegorz", "Joanna", "Piotr", "Ewa"}
	lastNames := []string{"Kowalski", "Nowak", "Wiśniewski", "Wójcik", "Kowalczyk", "Kamińska", "Lewandowski", "Zielińska", "Szymański", "Woźniak"}
	discordHandles := []string{
		"MichałKowalski#9821",
		"AnnaNowak#4532",
		"PawełWiśniewski#9210",
		"KatarzynaWójcik#3345",
		"TomaszKowalczyk#6658",
		"MagdalenaKamińska#7789",
		"GrzegorzLewandowski#2235",
		"JoannaZielińska#5567",
		"PiotrSzymański#9901",
		"EwaWoźniak#1123",
	}
	levelsOptions := [][]string{
		{"podstawówka", "liceum"},
		{"studia"},
		{"podstawówka", "studia"},
		{"liceum", "studia"},
		{"podstawówka"},
		{"liceum"},
		{"studia"},
		{"podstawówka", "liceum", "studia"},
		{"liceum", "studia"},
		{"studia"},
	}

	for i := 0; i < 10; i++ {
		// Generate a realistic username: first name + first letter of last name.
		username := fmt.Sprintf("%s%c", firstNames[i], lastNames[i][0])
		tutor := models.User{
			Email:       fmt.Sprintf("tutor%d@example.com", i+1),
			Username:    username,
			Role:        models.UserRoleTutor,
			DiscordID:   discordHandles[i],
			FirstName:   firstNames[i],
			LastName:    lastNames[i],
			DateOfBirth: "1985-05-15",
			Description: fmt.Sprintf("Cześć, jestem %s %s. Mam wieloletnie doświadczenie w nauczaniu oraz pasję do rozwijania umiejętności moich uczniów.", firstNames[i], lastNames[i]),
			Rating:      3.5 + float64(i%3)*0.5,
			Levels:      pq.StringArray(levelsOptions[i]), // Use pq.StringArray to ensure proper conversion.
		}

		if err := db.FirstOrCreate(&tutor, models.User{Email: tutor.Email}).Error; err != nil {
			return fmt.Errorf("failed to create/find tutor %d: %w", i+1, err)
		}
		tutors = append(tutors, tutor)
	}

	// 3) Define statuses and detailed lesson topics.
	statuses := []string{
		models.LessonStatusScheduled,
		models.LessonStatusConfirmed,
		models.LessonStatusInProgress,
		models.LessonStatusDone,
		models.LessonStatusCancelled,
		models.LessonStatusFailed,
	}

	subjects := []struct {
		Title       string
		Description string
	}{
		{"Algebra Introduction", "This lesson covers the basics of algebra including linear equations, variables, and problem-solving techniques."},
		{"Physics Basics", "An in-depth introduction to Newton’s laws, kinematics, and the fundamentals of mechanics with practical experiments."},
		{"Chemistry: The Periodic Table", "Explore atomic structures, chemical bonds, and the organization of elements in the periodic table."},
		{"World History Overview", "Discuss major global events, cultural revolutions, and historical milestones that have shaped modern society."},
		{"Biology: Genetics", "Learn about DNA, genetic inheritance, and molecular biology through interactive examples."},
		{"Computer Science Fundamentals", "An introduction to algorithms, data structures, and computational thinking with hands-on coding exercises."},
		{"Philosophy 101", "Examine major philosophical theories, influential thinkers, and the evolution of critical reasoning."},
		{"Music Theory", "Discover scales, chord progressions, and rhythm structures that form the foundation of musical composition."},
		{"Art & Painting", "Learn about color theory, composition, and creative techniques for artistic expression."},
		{"English Literature", "Analyze important literary works, writing styles, and methods for literary criticism in depth."},
		{"Advanced Algebra", "Dive deeper into algebra with complex equations, polynomial functions, and advanced problem-solving strategies."},
		{"Calculus I", "Understand limits, derivatives, and introductory integrals with practical, real-world examples."},
	}

	now := time.Now()
	var allLessons []models.Lesson
	lessonsPerTutor := 3 // Total 30 lessons

	for i, tutor := range tutors {
		for j := 0; j < lessonsPerTutor; j++ {
			statusIndex := (i*lessonsPerTutor + j) % len(statuses)
			chosenStatus := statuses[statusIndex]

			subjectIndex := (i*lessonsPerTutor + j) % len(subjects)
			chosenSubject := subjects[subjectIndex]

			// Offset time so lessons don't overlap.
			offsetHours := (i * 10) + (j * 2)
			start := now.Add(time.Duration(offsetHours) * time.Hour)
			end := start.Add(time.Hour)

			lesson := models.Lesson{
				TutorID:     tutor.ID,
				Students:    []models.User{student},
				Title:       chosenSubject.Title,       // Only the lesson topic.
				Description: chosenSubject.Description, // Detailed description.
				StartTime:   start,
				EndTime:     end,
				Status:      chosenStatus,
			}

			allLessons = append(allLessons, lesson)
		}
	}

	// 4) Insert all lessons in a batch.
	if err := db.Create(&allLessons).Error; err != nil {
		return fmt.Errorf("failed to create lessons: %w", err)
	}

	fmt.Println("==== SEEDING SUCCESS ====")
	fmt.Printf("Created %d tutors and %d total lessons for student (ID: %s).\n",
		len(tutors), len(allLessons), studentUUIDStr)
	return nil
}
