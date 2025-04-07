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
	// 2. Create 5 tutors with realistic Polish data.
	// ------------------------------
	tutorData := []struct {
		FirstName   string
		LastName    string
		Email       string
		DiscordID   string
		Description string
		Quote       string
		Rating      float64
		Price       float64
		Subjects    []string
		Levels      []string
	}{
		{
			FirstName:   "Aleksander",
			LastName:    "Nowak",
			Email:       "aleksander.nowak@edu.pl",
			DiscordID:   "AleksanderNowak#4532",
			Description: "Jestem doktorem matematyki stosowanej z 8-letnim doświadczeniem w nauczaniu na poziomie akademickim. Moją specjalnością jest analiza matematyczna i algebra liniowa. Uwielbiam pomagać studentom w odkrywaniu piękna matematyki i jej zastosowań w życiu codziennym. Mój styl nauczania opiera się na zrozumieniu, a nie zapamiętywaniu wzorów.",
			Quote:       "Matematyka to nie tylko liczby, to sposób myślenia o świecie.",
			Rating:      4.8,
			Price:       120.0,
			Subjects:    []string{"Matematyka", "Statystyka"},
			Levels:      []string{"Liceum", "Studia"},
		},
		{
			FirstName:   "Katarzyna",
			LastName:    "Kowalska",
			Email:       "katarzyna.kowalska@edu.pl",
			DiscordID:   "KatarzynaKowalska#7789",
			Description: "Absolwentka fizyki teoretycznej na Uniwersytecie Warszawskim. Od 5 lat prowadzę zajęcia dla uczniów szkół średnich przygotowujących się do olimpiad fizycznych. Moim celem jest pokazanie, że fizyka nie jest trudna, a wręcz fascynująca, gdy zrozumie się podstawowe zasady. Stosuję metody wizualizacji i eksperymenty, które pomagają uczniom lepiej zrozumieć abstrakcyjne pojęcia.",
			Quote:       "Fizyka wyjaśnia świat, a ja pomagam zrozumieć fizykę.",
			Rating:      4.9,
			Price:       110.0,
			Subjects:    []string{"Fizyka", "Astronomia"},
			Levels:      []string{"Szkoła Podstawowa", "Liceum"},
		},
		{
			FirstName:   "Michał",
			LastName:    "Wiśniewski",
			Email:       "michal.wisniewski@edu.pl",
			DiscordID:   "MichalWisniewski#9210",
			Description: "Polonista z pasją, nauczyciel z 12-letnim stażem w liceum. Specjalizuję się w literaturze polskiej XIX i XX wieku. Pomagam uczniom nie tylko zdać maturę, ale także pokochać literaturę. Moje lekcje to połączenie analizy tekstów z ciekawostkami historycznymi i kulturowymi, które sprawiają, że nauka staje się przygodą.",
			Quote:       "Literatura to podróż w czasie i przestrzeni bez wychodzenia z domu.",
			Rating:      4.7,
			Price:       90.0,
			Subjects:    []string{"Język Polski", "Literatura", "Historia"},
			Levels:      []string{"Szkoła Podstawowa", "Liceum"},
		},
		{
			FirstName:   "Agnieszka",
			LastName:    "Lewandowska",
			Email:       "agnieszka.lewandowska@edu.pl",
			DiscordID:   "AgnieszkaLewandowska#6658",
			Description: "Jestem programistką z 10-letnim doświadczeniem w branży IT oraz nauczycielką informatyki. Specjalizuję się w programowaniu w językach Python, Java i C++. Prowadzę kursy zarówno dla początkujących, jak i zaawansowanych. Wierzę, że każdy może nauczyć się programować przy odpowiednim podejściu. Moje lekcje są praktyczne i zorientowane na projekty.",
			Quote:       "Kod to nowy język komunikacji ze światem. Nauczmy się go razem.",
			Rating:      4.9,
			Price:       130.0,
			Subjects:    []string{"Informatyka", "Programowanie"},
			Levels:      []string{"Liceum", "Studia"},
		},
		{
			FirstName:   "Piotr",
			LastName:    "Kowalczyk",
			Email:       "piotr.kowalczyk@edu.pl",
			DiscordID:   "PiotrKowalczyk#2235",
			Description: "Biolog molekularny i pasjonat nauk przyrodniczych. Prowadzę zajęcia dla uczniów wszystkich poziomów, od podstawówki po studia. Moją specjalnością jest genetyka i ekologia. Staram się przekazać nie tylko wiedzę, ale także ciekawość świata. Wykorzystuję multimedialne prezentacje i eksperymenty, które można wykonać w domu.",
			Quote:       "Biologia to nauka o życiu, a zrozumienie jej pozwala lepiej żyć.",
			Rating:      4.6,
			Price:       100.0,
			Subjects:    []string{"Biologia", "Chemia"},
			Levels:      []string{"Szkoła Podstawowa", "Liceum", "Studia"},
		},
	}

	rand.Seed(time.Now().UnixNano())
	var tutors []models.User
	for i, data := range tutorData {
		tutor := models.User{
			Email:       data.Email,
			Username:    fmt.Sprintf("%s%c", data.FirstName, data.LastName[0]),
			Role:        models.UserRoleTutor,
			DiscordID:   data.DiscordID,
			FirstName:   data.FirstName,
			LastName:    data.LastName,
			DateOfBirth: fmt.Sprintf("198%d-0%d-15", i+1, i+1),
			Description: data.Description,
			Rating:      data.Rating,
			Levels:      pq.StringArray(data.Levels),
			Subjects:    pq.StringArray(data.Subjects),
			Quote:       data.Quote,
			Price:       data.Price,
		}
		if err := db.FirstOrCreate(&tutor, models.User{Email: tutor.Email}).Error; err != nil {
			return fmt.Errorf("error creating/finding Tutor %d: %w", i+1, err)
		}
		tutors = append(tutors, tutor)
	}

	// ------------------------------
	// 3. Create weekly availability for each tutor
	// ------------------------------
	weeklySchedules := []struct {
		TutorIndex int
		DayOfWeek  int
		StartTime  string
		EndTime    string
	}{
		// Aleksander Nowak (Monday, Wednesday, Friday afternoons)
		{0, 1, "15:00", "19:00"}, // Monday
		{0, 3, "16:00", "20:00"}, // Wednesday
		{0, 5, "14:00", "18:00"}, // Friday

		// Katarzyna Kowalska (Tuesday, Thursday mornings & evenings)
		{1, 2, "09:00", "12:00"}, // Tuesday morning
		{1, 2, "17:00", "20:00"}, // Tuesday evening
		{1, 4, "09:00", "12:00"}, // Thursday morning
		{1, 4, "17:00", "20:00"}, // Thursday evening

		// Michał Wiśniewski (Monday, Wednesday, Saturday)
		{2, 1, "10:00", "14:00"}, // Monday
		{2, 3, "14:00", "18:00"}, // Wednesday
		{2, 6, "10:00", "16:00"}, // Saturday

		// Agnieszka Lewandowska (Weekday evenings)
		{3, 1, "18:00", "21:00"}, // Monday
		{3, 2, "18:00", "21:00"}, // Tuesday
		{3, 3, "18:00", "21:00"}, // Wednesday
		{3, 4, "18:00", "21:00"}, // Thursday
		{3, 5, "18:00", "21:00"}, // Friday

		// Piotr Kowalczyk (Daily availability)
		{4, 1, "12:00", "15:00"}, // Monday
		{4, 2, "12:00", "15:00"}, // Tuesday
		{4, 3, "12:00", "15:00"}, // Wednesday
		{4, 4, "12:00", "15:00"}, // Thursday
		{4, 5, "12:00", "15:00"}, // Friday
		{4, 6, "10:00", "14:00"}, // Saturday
		{4, 0, "14:00", "18:00"}, // Sunday
	}

	for _, schedule := range weeklySchedules {
		weeklySchedule := models.TutorWeeklySchedule{
			TutorID:   tutors[schedule.TutorIndex].ID,
			DayOfWeek: schedule.DayOfWeek,
			StartTime: schedule.StartTime,
			EndTime:   schedule.EndTime,
		}
		if err := db.Create(&weeklySchedule).Error; err != nil {
			return fmt.Errorf("error creating weekly schedule: %w", err)
		}
	}

	// ------------------------------
	// 4. Create schedule exceptions (both additions and removals)
	// ------------------------------

	// Get dates for next month to create exceptions
	now := time.Now()
	nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())

	// Create holidays (removals)
	holidays := []time.Time{
		// Next month holiday dates
		time.Date(nextMonth.Year(), nextMonth.Month(), 3, 0, 0, 0, 0, nextMonth.Location()),  // Example holiday
		time.Date(nextMonth.Year(), nextMonth.Month(), 15, 0, 0, 0, 0, nextMonth.Location()), // Example holiday
	}

	// Create removals for all tutors on holidays
	for _, tutor := range tutors {
		for _, holiday := range holidays {
			exception := models.TutorScheduleException{
				TutorID:   tutor.ID,
				Date:      holiday,
				IsRemoval: true, // This is a removal (unavailable day)
			}
			if err := db.Create(&exception).Error; err != nil {
				return fmt.Errorf("error creating schedule exception (holiday): %w", err)
			}
		}
	}

	// Create special availability (additions)
	specialAvailability := []struct {
		TutorIndex int
		Date       time.Time
		StartTime  string
		EndTime    string
	}{
		// Aleksander Nowak - special Sunday session next month
		{0, time.Date(nextMonth.Year(), nextMonth.Month(), 5, 0, 0, 0, 0, nextMonth.Location()), "10:00", "14:00"},

		// Katarzyna Kowalska - extra Saturday session
		{1, time.Date(nextMonth.Year(), nextMonth.Month(), 10, 0, 0, 0, 0, nextMonth.Location()), "12:00", "16:00"},

		// Michał Wiśniewski - extra evening slot
		{2, time.Date(nextMonth.Year(), nextMonth.Month(), 8, 0, 0, 0, 0, nextMonth.Location()), "19:00", "21:00"},

		// Agnieszka Lewandowska - morning session
		{3, time.Date(nextMonth.Year(), nextMonth.Month(), 12, 0, 0, 0, 0, nextMonth.Location()), "09:00", "12:00"},

		// Piotr Kowalczyk - special full day session
		{4, time.Date(nextMonth.Year(), nextMonth.Month(), 20, 0, 0, 0, 0, nextMonth.Location()), "09:00", "17:00"},
	}

	for _, special := range specialAvailability {
		exception := models.TutorScheduleException{
			TutorID:   tutors[special.TutorIndex].ID,
			Date:      special.Date,
			StartTime: special.StartTime,
			EndTime:   special.EndTime,
			IsRemoval: false, // This is an addition (special available time)
		}
		if err := db.Create(&exception).Error; err != nil {
			return fmt.Errorf("error creating schedule exception (special): %w", err)
		}
	}

	// ------------------------------
	// 5. Create standalone lessons for each tutor (for the fixed student)
	// with varying statuses. For these lessons, leave CourseID empty (nil).
	// ------------------------------
	lessonSubjects := []struct {
		Title       string
		Description string
	}{
		{
			"Funkcje kwadratowe i ich właściwości",
			"Podczas lekcji omówimy funkcje kwadratowe, ich wykresy (parabole), miejsca zerowe oraz zastosowania praktyczne. Nauczysz się analizować równania kwadratowe i rozwiązywać problemy z życia codziennego.",
		},
		{
			"Mechanika Newtona - podstawy dynamiki",
			"Zajęcia wprowadzające do praw Newtona. Omówimy pojęcie siły, masy i przyspieszenia oraz ich wzajemne relacje. Rozwiążemy przykładowe zadania dotyczące ruchu ciał pod wpływem różnych sił.",
		},
		{
			"Młoda Polska - główne nurty literackie",
			"Przegląd literatury okresu Młodej Polski. Skupimy się na symbolizmie, impresjonizmie i ekspresjonizmie w dziełach Wyspiańskiego, Tetmajera i Kasprowicza. Przeanalizujemy wybrane fragmenty utworów.",
		},
		{
			"Podstawy programowania w Pythonie",
			"Wprowadzenie do języka Python - struktury danych, pętle, funkcje. W trakcie zajęć napiszemy kilka prostych programów demonstrujących możliwości tego języka i jego zastosowania praktyczne.",
		},
		{
			"Genetyka mendlowska i jej zastosowania",
			"Lekcja poświęcona prawom dziedziczenia Mendla. Omówimy krzyżówki genetyczne, dominację cech oraz przykłady dziedziczenia u człowieka. Rozwiążemy zadania z genetyki klasycznej.",
		},
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

			// Select subject and level from tutor's specializations
			tutorSubject := string(tutor.Subjects[0])
			tutorLevel := string(tutor.Levels[0])

			lesson := models.Lesson{
				TutorID:     tutor.ID,
				Students:    []models.User{student},
				Title:       subj.Title,
				Description: subj.Description,
				Subject:     tutorSubject,
				Level:       tutorLevel,
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
	// 6. Create a first real course and fetch it so that its ID is set.
	// ------------------------------
	course1 := models.Course{
		ID:          uuid.New(),
		TutorID:     tutors[0].ID,
		Tutor:       tutors[0],
		Name:        "Matematyka wyższa dla ambitnych",
		Description: "Kompleksowy kurs matematyki wyższej dla uczniów liceum przygotowujących się do olimpiad matematycznych i studiów technicznych. Kurs obejmuje zagadnienia analizy matematycznej, algebry, geometrii analitycznej oraz elementów rachunku prawdopodobieństwa na poziomie zaawansowanym.",
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

	course1LessonDetails := []struct {
		Title       string
		Description string
	}{
		{
			"Granice funkcji i ciągłość",
			"Wprowadzenie do pojęcia granicy funkcji, warunki istnienia granicy, funkcje ciągłe i ich własności.",
		},
		{
			"Rachunek różniczkowy - pochodne",
			"Definicja pochodnej, interpretacja geometryczna, reguły różniczkowania, zastosowania pochodnych.",
		},
		{
			"Ekstrema funkcji i ich zastosowania",
			"Wyznaczanie ekstremów funkcji przy pomocy pochodnych, problemy optymalizacyjne.",
		},
		{
			"Całka nieoznaczona i metody całkowania",
			"Pojęcie całki nieoznaczonej, metody całkowania: podstawianie, przez części, całkowanie funkcji wymiernych.",
		},
		{
			"Całka oznaczona i jej zastosowania",
			"Definicja całki oznaczonej, związek z całką nieoznaczoną, obliczanie pól, objętości brył obrotowych.",
		},
	}

	var course1Lessons []models.Lesson
	for i, lessonDetail := range course1LessonDetails {
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
			Title:       lessonDetail.Title,
			Description: lessonDetail.Description,
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
	// 7. Create a second course (without any student) and add lessons.
	// ------------------------------
	course2 := models.Course{
		ID:          uuid.New(),
		TutorID:     tutors[1].ID,
		Tutor:       tutors[1],
		Name:        "Fizyka kwantowa dla licealistów",
		Description: "Wprowadzenie do fascynującego świata fizyki kwantowej dla uczniów liceum zainteresowanych fizyką teoretyczną. Kurs obejmuje podstawy mechaniki kwantowej, dualizm korpuskularno-falowy, zasadę nieoznaczoności oraz elementy fizyki cząstek elementarnych.",
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

	course2LessonDetails := []struct {
		Title       string
		Description string
	}{
		{
			"Promieniowanie ciała doskonale czarnego",
			"Historia fizyki kwantowej, promieniowanie ciała doskonale czarnego, katastrofa w nadfiolecie, hipoteza Plancka.",
		},
		{
			"Efekt fotoelektryczny i kwanty światła",
			"Zjawisko fotoelektryczne, wyjaśnienie Einsteina, pojęcie fotonu, dualizm korpuskularno-falowy światła.",
		},
		{
			"Model atomu Bohra",
			"Model atomu Rutherforda, postulaty Bohra, poziomy energetyczne, spektroskopia atomowa.",
		},
		{
			"Dualizm korpuskularno-falowy materii",
			"Hipoteza de Broglie'a, fale materii, eksperyment Davissona-Germera, dyfrakcja elektronów.",
		},
		{
			"Zasada nieoznaczoności Heisenberga",
			"Nieoznaczoność pomiaru, relacje nieoznaczoności dla położenia i pędu, energia i czas, interpretacja kopenhaska.",
		},
	}

	var course2Lessons []models.Lesson
	for i, lessonDetail := range course2LessonDetails {
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
			Title:       lessonDetail.Title,
			Description: lessonDetail.Description,
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
