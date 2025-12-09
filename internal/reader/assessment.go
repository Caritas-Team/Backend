package reader

import "time"

// AssessmentResult представляет результаты оценки коммуникативного развития
type AssessmentResult struct {
	// Портрет пользователя
	CompletionDate  time.Time
	SubjectName     string
	Diagnosis       string
	DateOfBirth     time.Time
	SocialSituation string
	ResidenceType   string

	// Базовая оценка
	VerbalSpeech        string
	WrittenSpeech       string
	Vision              string
	Hearing             string
	SpeechComprehension string
	MotorSkills         string
	PointingMethod      string

	// Языковые навыки (%)
	PreIntentional int
	Protolanguage  int
	Holophrase     int
	Phrase         int

	// Коммуникативные функции (%)
	ControlInitiative int
	DesireInitiative  int
	SocialInitiative  int
	InfoInitiative    int

	// Инициатива и частота
	InitiativeZBR int
	InitiativeBR  int
	Frequency     string

	// Коммуникативные функции - детали
	CommunicativeFunctions []CommunicativeFunction

	// Интересы
	FavoriteCharacters []string
	LikedActivities    []string
	DiscomfortElements []string

	// Способы общения
	CommunicationMethods []string

	// Вербальные слова
	VerbalWords []string

	// Быстрые сообщения
	QuickMessages []string

	// Круги общения
	CommunicationCircles CirclesOfCommunication

	// Рекомендации
	Recommendations []string
}

// CommunicativeFunction описывает коммуникативную функцию
type CommunicativeFunction struct {
	Name         string
	Description  string
	Protolangage string // "сформирован", "уже не используется", "недоступно"
	Holophrase   string
	Phrase       string
	Initiative   int
	Frequency    int
}

// CirclesOfCommunication описывает круги общения
type CirclesOfCommunication struct {
	Family        []string
	Friends       []string
	Acquaintances []string
	Specialists   []string
	Notes         string
}

// VocabularyItem представляет словарный элемент
type VocabularyItem struct {
	Word     string
	Category string
}

// ParseError представляет ошибку парсинга
type ParseError struct {
	Field   string
	Message string
	Line    int
}

func (e ParseError) Error() string {
	return "parse error in field " + e.Field + ": " + e.Message
}
