package seeds

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"

	"vibely-backend/src/models"
)

func Seed(db *gorm.DB) error {
	// ------------------------------
	// 1. Retrieve the fixed student.
	// ------------------------------
	studentUUIDStr := "8e686fcf-3849-4efd-8e22-16280d3d310f"
	studentUUID, err := uuid.Parse(studentUUIDStr)
	if err != nil {
		return fmt.Errorf("error parsing student UUID: %w", err)
	}
	var student models.User
	if err := db.First(&student, "id = ?", studentUUID).Error; err != nil {
		return fmt.Errorf("student not found (ID %s): %w", studentUUIDStr, err)
	}

	// ------------------------------
	// 2. Create 10 tutors.
	// ------------------------------
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
	subjectsOptions := [][]string{
		{"Matematyka", "Fizyka"},
		{"Chemia", "Biologia"},
		{"Język Polski", "Język Angielski"},
		{"Historia", "Geografia"},
		{"Informatyka", "Matematyka"},
		{"Sztuka", "Muzyka"},
		{"Ekonomia", "Statystyka"},
		{"Biologia", "Chemia"},
		{"Fizyka", "Astronomia"},
		{"Niemiecki", "Francuski"},
	}
	quotes := []string{
		"Wierzę w odkrywanie pełnego potencjału każdego ucznia.",
		"Nauczanie to moja pasja i praca życia.",
		"Staram się inspirować ciekawość i chęć nauki.",
		"Edukacja jest kluczem do sukcesu!",
		"Każda lekcja to krok w kierunku sukcesu.",
		"Wiedza dzielona jest wiedzą mnożoną.",
		"Inspiracja i nauka idą w parze.",
		"Uczmy się razem i rozwijajmy.",
		"Twoja edukacja to inwestycja w przyszłość.",
		"Dążę do tego, by nauka była przyjemnością.",
	}

	rand.Seed(time.Now().UnixNano())
	var tutors []models.User
	for i := 0; i < 10; i++ {
		tutor := models.User{
			Email:       fmt.Sprintf("tutor%d@example.com", i+1),
			Username:    fmt.Sprintf("%s%c", firstNames[i], lastNames[i][0]),
			Role:        models.UserRoleTutor,
			DiscordID:   discordHandles[i],
			FirstName:   firstNames[i],
			LastName:    lastNames[i],
			DateOfBirth: "1985-05-15",
			Description: fmt.Sprintf("Cześć, jestem %s %s. Mam doświadczenie w %s.",
				firstNames[i], lastNames[i], subjectsOptions[i][0]),
			Rating:   4.0,
			Levels:   pq.StringArray(levelsOptions[i]),
			Subjects: pq.StringArray(subjectsOptions[i]),
			Quote:    quotes[i],
			Price:    50.0 + float64(i)*10.0,
		}
		if err := db.FirstOrCreate(&tutor, models.User{Email: tutor.Email}).Error; err != nil {
			return fmt.Errorf("error creating/finding Tutor %d: %w", i+1, err)
		}
		tutors = append(tutors, tutor)
	}

	// ------------------------------
	// 3. Create standalone lessons for each tutor (for the fixed student)
	// with varying statuses. For these lessons, leave CourseID empty (nil).
	// ------------------------------
	lessonSubjects := []struct {
		Title       string
		Description string
	}{
		{"Wstęp do algebry", "Podstawy algebry, równania liniowe."},
		{"Podstawy fizyki", "Wprowadzenie do praw Newtona."},
		{"Chemia: Układ okresowy", "Budowa atomu i pierwiastków."},
		{"Historia", "Omówienie głównych wydarzeń."},
		{"Biologia", "Wprowadzenie do dziedziczenia cech."},
	}
	// We'll create 3 lessons per tutor.
	lessonsPerTutor := 3
	nowTime := time.Now()
	var standaloneLessons []models.Lesson
	for i, tutor := range tutors {
		for j := 0; j < lessonsPerTutor; j++ {
			statusIndex := (i*lessonsPerTutor + j) % 6
			var status string
			var start, end time.Time
			switch statusIndex {
			case 0:
				status = models.LessonStatusScheduled // future date
				start = nowTime.Add(2 * time.Hour)
				end = start.Add(time.Hour)
			case 1:
				status = models.LessonStatusConfirmed // future date
				start = nowTime.Add(1 * time.Hour)
				end = start.Add(time.Hour)
			case 2:
				status = models.LessonStatusInProgress // around now
				start = nowTime.Add(-15 * time.Minute)
				end = nowTime.Add(45 * time.Minute)
			case 3:
				status = models.LessonStatusDone // past
				start = nowTime.Add(-2 * time.Hour)
				end = start.Add(time.Hour)
			case 4:
				status = models.LessonStatusFailed // past
				start = nowTime.Add(-3 * time.Hour)
				end = start.Add(time.Hour)
			case 5:
				status = models.LessonStatusCancelled // past
				start = nowTime.Add(-4 * time.Hour)
				end = start.Add(time.Hour)
			}
			idx := (i*lessonsPerTutor + j) % len(lessonSubjects)
			subj := lessonSubjects[idx]
			lesson := models.Lesson{
				TutorID:     tutor.ID,
				Students:    []models.User{student},
				Title:       subj.Title,
				Description: subj.Description,
				Subject:     subjectsOptions[i%len(subjectsOptions)][0],
				Level:       levelsOptions[i%len(levelsOptions)][0],
				StartTime:   start,
				EndTime:     end,
				Status:      status,
				// Leave CourseID empty (nil) for standalone lessons.
				CourseID: nil,
			}
			standaloneLessons = append(standaloneLessons, lesson)
		}
	}
	if err := db.Create(&standaloneLessons).Error; err != nil {
		return fmt.Errorf("error creating standalone lessons: %w", err)
	}

	// ------------------------------
	// 4. Create a first real course and fetch it so that its ID is set.
	// ------------------------------
	course1 := models.Course{
		ID:          uuid.New(),
		TutorID:     tutors[0].ID,
		Tutor:       tutors[0],
		Name:        "Mathematics Mastery",
		Description: "A comprehensive course on mathematics fundamentals.",
		Subject:     "Matematyka",
		Level:       "Liceum",
		Students:    []models.User{student},
		CreatedAt:   nowTime,
	}
	if err := db.Create(&course1).Error; err != nil {
		return fmt.Errorf("error creating course1: %w", err)
	}
	if err := db.First(&course1, "id = ?", course1.ID).Error; err != nil {
		return fmt.Errorf("error fetching course1 after creation: %w", err)
	}

	var course1Lessons []models.Lesson
	for i := 0; i < 5; i++ {
		statusIndex := i % 6
		var status string
		var start, end time.Time
		switch statusIndex {
		case 0:
			status = models.LessonStatusScheduled
			start = nowTime.Add(2 * time.Hour)
			end = start.Add(time.Hour)
		case 1:
			status = models.LessonStatusConfirmed
			start = nowTime.Add(1 * time.Hour)
			end = start.Add(time.Hour)
		case 2:
			status = models.LessonStatusInProgress
			start = nowTime.Add(-15 * time.Minute)
			end = nowTime.Add(45 * time.Minute)
		case 3:
			status = models.LessonStatusDone
			start = nowTime.Add(-2 * time.Hour)
			end = start.Add(time.Hour)
		case 4:
			status = models.LessonStatusFailed
			start = nowTime.Add(-3 * time.Hour)
			end = start.Add(time.Hour)
		case 5:
			status = models.LessonStatusCancelled
			start = nowTime.Add(-4 * time.Hour)
			end = start.Add(time.Hour)
		}
		lesson := models.Lesson{
			TutorID:     tutors[0].ID,
			Title:       fmt.Sprintf("Mathematics Lesson %d", i+1),
			Description: fmt.Sprintf("Description for Mathematics Lesson %d", i+1),
			Subject:     "Matematyka",
			Level:       "Liceum",
			StartTime:   start,
			EndTime:     end,
			Status:      status,
			// Set CourseID to course1's ID.
			CourseID: &course1.ID,
			Students: []models.User{student},
		}
		course1Lessons = append(course1Lessons, lesson)
	}
	if err := db.Create(&course1Lessons).Error; err != nil {
		return fmt.Errorf("error creating lessons for course1: %w", err)
	}

	// ------------------------------
	// 5. Create a second course (without any student) and add lessons.
	// ------------------------------
	course2 := models.Course{
		ID:          uuid.New(),
		TutorID:     tutors[1].ID,
		Tutor:       tutors[1],
		Name:        "Physics Fundamentals",
		Description: "An introductory course on physics covering essential concepts.",
		Subject:     "Fizyka",
		Level:       "Liceum",
		CreatedAt:   nowTime,
	}
	if err := db.Create(&course2).Error; err != nil {
		return fmt.Errorf("error creating course2: %w", err)
	}
	if err := db.First(&course2, "id = ?", course2.ID).Error; err != nil {
		return fmt.Errorf("error fetching course2 after creation: %w", err)
	}
	var course2Lessons []models.Lesson
	for i := 0; i < 5; i++ {
		statusIndex := i % 6
		var status string
		var start, end time.Time
		switch statusIndex {
		case 0:
			status = models.LessonStatusScheduled
			start = nowTime.Add(2 * time.Hour)
			end = start.Add(time.Hour)
		case 1:
			status = models.LessonStatusConfirmed
			start = nowTime.Add(1 * time.Hour)
			end = start.Add(time.Hour)
		case 2:
			status = models.LessonStatusInProgress
			start = nowTime.Add(-15 * time.Minute)
			end = nowTime.Add(45 * time.Minute)
		case 3:
			status = models.LessonStatusDone
			start = nowTime.Add(-2 * time.Hour)
			end = start.Add(time.Hour)
		case 4:
			status = models.LessonStatusFailed
			start = nowTime.Add(-3 * time.Hour)
			end = start.Add(time.Hour)
		case 5:
			status = models.LessonStatusCancelled
			start = nowTime.Add(-4 * time.Hour)
			end = start.Add(time.Hour)
		}
		lesson := models.Lesson{
			TutorID:     tutors[1].ID,
			Title:       fmt.Sprintf("Physics Lesson %d", i+1),
			Description: fmt.Sprintf("Description for Physics Lesson %d", i+1),
			Subject:     "Fizyka",
			Level:       "Liceum",
			StartTime:   start,
			EndTime:     end,
			Status:      status,
			// Set CourseID to course2's ID.
			CourseID: &course2.ID,
			// No students assigned.
		}
		course2Lessons = append(course2Lessons, lesson)
	}
	if err := db.Create(&course2Lessons).Error; err != nil {
		return fmt.Errorf("error creating lessons for course2: %w", err)
	}

	fmt.Println("==== SEEDING SUCCESS ====")
	return nil
}
