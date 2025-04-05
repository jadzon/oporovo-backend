package seeds

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq" // Ensure this is imported
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
		{"Szkoła Podstawowa", "Liceum"},
		{"Studia"},
		{"Szkoła Podstawowa", "Studia"},
		{"Liceum", "Studia"},
		{"Szkoła Podstawowa"},
		{"Liceum"},
		{"Studia"},
		{"Szkoła Podstawowa", "Liceum", "Studia"},
		{"Liceum", "Studia"},
		{"Studia"},
	}
	// New subjects options for tutors.
	subjectsOptions := [][]string{
		{"Matematyka", "Fizyka"},
		{"Chemia", "Biologia"},
		{"Język Polski", "Język Angielski"},
		{"Historia", "Geografia"},
		{"Informatyka", "Matematyka"},
		{"Chemia", "Fizyka"},
		{"Biologia", "Język Polski"},
		{"Historia", "Informatyka"},
		{"Matematyka", "Chemia"},
		{"Fizyka", "Biologia"},
	}
	// New inspirational quotes for tutors.
	quotes := []string{
		"I believe in unlocking every student's potential.",
		"Teaching is my passion and life's work.",
		"I strive to inspire curiosity and learning.",
		"Education is the key to success!",
		"Let's learn together and grow.",
		"Every lesson is a step toward success.",
		"Knowledge shared is knowledge multiplied.",
		"I'm here to guide you on your learning journey.",
		"Passionate about making learning fun.",
		"Empowering students one lesson at a time.",
	}

	// Set prices creatively (e.g., 50 + (i%5)*15 gives values between 50 and 110 zł).
	var tutorPrice float64

	for i := 0; i < 10; i++ {
		// Generate a realistic username: first name + first letter of last name.
		username := fmt.Sprintf("%s%c", firstNames[i], lastNames[i][0])
		tutorPrice = 50.0 + float64(i%5)*15.0

		// Create a tutor with enriched data.
		tutor := models.User{
			Email:       fmt.Sprintf("tutor%d@example.com", i+1),
			Username:    username,
			Role:        models.UserRoleTutor,
			DiscordID:   discordHandles[i],
			FirstName:   firstNames[i],
			LastName:    lastNames[i],
			DateOfBirth: "1985-05-15",
			// Description remains as a longer bio.
			Description: fmt.Sprintf("Cześć, jestem %s %s. Mam wieloletnie doświadczenie w nauczaniu i specjalizuję się w %s. Moim celem jest inspirowanie uczniów do osiągania sukcesów.",
				firstNames[i], lastNames[i], subjectsOptions[i][0]),
			Rating:   3.5 + float64(i%3)*0.5,
			Levels:   pq.StringArray(levelsOptions[i]),
			Subjects: pq.StringArray(subjectsOptions[i]),
			// New fields for creative data.
			Quote: quotes[i],
			Price: tutorPrice,
		}

		// Create or retrieve the tutor based on unique email.
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

	lessonSubjects := []struct {
		Title       string
		Description string
	}{
		{"Algebra Introduction", "This lesson covers the basics of algebra including linear equations, variables, and problem-solving techniques."},
		{"Physics Basics", "An in-depth introduction to Newton’s laws, kinematics, and fundamentals of mechanics with practical experiments."},
		{"Chemistry: The Periodic Table", "Explore atomic structures, chemical bonds, and the organization of elements."},
		{"World History Overview", "Discuss major global events, cultural revolutions, and historical milestones."},
		{"Biology: Genetics", "Learn about DNA, genetic inheritance, and molecular biology through interactive examples."},
		{"Computer Science Fundamentals", "An introduction to algorithms, data structures, and coding with hands-on exercises."},
		{"Philosophy 101", "Examine major philosophical theories and influential thinkers in a critical discussion."},
		{"Music Theory", "Discover scales, chord progressions, and rhythm structures foundational to musical composition."},
		{"Art & Painting", "Learn about color theory, composition, and creative techniques for artistic expression."},
		{"English Literature", "Analyze important literary works and methods for literary criticism in depth."},
		{"Advanced Algebra", "Dive deeper into algebra with complex equations, polynomial functions, and advanced strategies."},
		{"Calculus I", "Understand limits, derivatives, and introductory integrals with practical examples."},
	}

	now := time.Now()
	var allLessons []models.Lesson
	lessonsPerTutor := 3 // Total 30 lessons

	for i, tutor := range tutors {
		for j := 0; j < lessonsPerTutor; j++ {
			statusIndex := (i*lessonsPerTutor + j) % len(statuses)
			chosenStatus := statuses[statusIndex]

			subjectIndex := (i*lessonsPerTutor + j) % len(lessonSubjects)
			chosenSubject := lessonSubjects[subjectIndex]

			// Offset time so lessons don't overlap.
			offsetHours := (i * 10) + (j * 2)
			start := now.Add(time.Duration(offsetHours) * time.Hour)
			end := start.Add(time.Hour)

			lesson := models.Lesson{
				TutorID:     tutor.ID,
				Students:    []models.User{student},
				Title:       chosenSubject.Title,
				Description: chosenSubject.Description,
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
