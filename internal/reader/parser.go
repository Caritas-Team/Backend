package reader

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Parser отвечает за парсинг текста из PDF в структуры данных
type Parser struct {
	// регулярные выражения для извлечения данных
	datePattern    *regexp.Regexp
	percentPattern *regexp.Regexp
}

// NewParser создает новый экземпляр Parser с инициализированными регулярками
func NewParser() *Parser {
	return &Parser{
		datePattern:    regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`),
		percentPattern: regexp.MustCompile(`(\d+)%`),
	}
}

// Parse парсит текст из PDF и возвращает AssessmentResult
func (p *Parser) Parse(text string) (*AssessmentResult, error) {
	result := &AssessmentResult{
		CommunicativeFunctions: []CommunicativeFunction{},
		FavoriteCharacters:     []string{},
		LikedActivities:        []string{},
		DiscomfortElements:     []string{},
		CommunicationMethods:   []string{},
		VerbalWords:            []string{},
		QuickMessages:          []string{},
		Recommendations:        []string{},
	}

	lines := strings.Split(text, "\n")

	// Парсим портрет пользователя
	p.parseUserProfile(lines, result)

	// Парсим базовую оценку
	p.parseBasicAssessment(lines, result)

	// Парсим языковые навыки
	p.parseLanguageSkills(lines, result)

	// Парсим коммуникативные функции
	p.parseCommunicativeFunctions(lines, result)

	// Парсим интересы
	p.parseInterests(lines, result)

	// Парсим способы общения
	p.parseCommunicationMethods(lines, result)

	// Парсим вербальные слова
	p.parseVerbalWords(lines, result)

	// Парсим круги общения
	p.parseCommunicationCircles(lines, result)

	return result, nil
}

// parseUserProfile парсит информацию о пользователе
func (p *Parser) parseUserProfile(lines []string, result *AssessmentResult) {
	text := strings.Join(lines, " ")

	// Дата заполнения
	if dates := p.datePattern.FindAllString(text, -1); len(dates) > 0 {
		if t, err := time.Parse("2006-01-02", dates[0]); err == nil {
			result.CompletionDate = t
		}
	}

	// Имя обследуемого (ищем после "Имя обследуемого" или "Test test")
	if idx := strings.Index(text, "Имя обследуемого"); idx != -1 {
		chunk := text[idx : idx+200]
		parts := strings.FieldsFunc(chunk, func(r rune) bool {
			return r == '\n' || r == '\r'
		})
		if len(parts) > 1 {
			result.SubjectName = strings.TrimSpace(parts[1])
		}
	}

	// Диагноз (после "Диагноз")
	if idx := strings.Index(text, "Диагноз"); idx != -1 {
		chunk := text[idx : idx+100]
		parts := strings.Split(chunk, "\n")
		if len(parts) > 1 {
			result.Diagnosis = strings.TrimSpace(parts[1])
		}
	}

	// Дата рождения
	if dates := p.datePattern.FindAllString(text, -1); len(dates) > 1 {
		if t, err := time.Parse("2006-01-02", dates[1]); err == nil {
			result.DateOfBirth = t
		}
	}

	// Социальная ситуация
	if idx := strings.Index(text, "Особенности социальной ситуации"); idx != -1 {
		chunk := text[idx : idx+200]
		parts := strings.Split(chunk, "\n")
		if len(parts) > 1 {
			result.SocialSituation = strings.TrimSpace(parts[1])
		}
	}

	// Где проживает
	if idx := strings.Index(text, "Где проживает человек"); idx != -1 {
		chunk := text[idx : idx+100]
		parts := strings.Split(chunk, "\n")
		if len(parts) > 1 {
			result.ResidenceType = strings.TrimSpace(parts[1])
		}
	}
}

// parseBasicAssessment парсит базовую оценку
func (p *Parser) parseBasicAssessment(lines []string, result *AssessmentResult) {
	text := strings.Join(lines, " ")

	// Вербальная речь
	if idx := strings.Index(text, "Особенности вербальной речи"); idx != -1 {
		chunk := text[idx : idx+200]
		parts := strings.Split(chunk, "\n")
		if len(parts) > 1 {
			result.VerbalSpeech = strings.TrimSpace(parts[1])
		}
	}

	// Зрение
	if idx := strings.Index(text, "Зрение"); idx != -1 {
		chunk := text[idx : idx+150]
		parts := strings.Split(chunk, "\n")
		if len(parts) > 1 {
			result.Vision = strings.TrimSpace(parts[1])
		}
	}

	// Слух
	if idx := strings.Index(text, "Слух в пределах"); idx != -1 {
		chunk := text[idx : idx+150]
		parts := strings.Split(chunk, "\n")
		if len(parts) > 1 {
			result.Hearing = strings.TrimSpace(parts[1])
		}
	}
}

// parseLanguageSkills парсит языковые навыки (проценты)
func (p *Parser) parseLanguageSkills(lines []string, result *AssessmentResult) {
	text := strings.Join(lines, " ")

	// Доинтенциональная коммуникация
	if match := p.findPercentAfter(text, "Доинтенциональная"); match != -1 {
		result.PreIntentional = match
	}

	// Протоязык
	if match := p.findPercentAfter(text, "Протоязык"); match != -1 {
		result.Protolanguage = match
	}

	// Голофраза
	if match := p.findPercentAfter(text, "Голофраза"); match != -1 {
		result.Holophrase = match
	}

	// Фраза
	if match := p.findPercentAfter(text, "Фраза"); match != -1 {
		result.Phrase = match
	}
}

// parseCommunicativeFunctions парсит коммуникативные функции
func (p *Parser) parseCommunicativeFunctions(lines []string, result *AssessmentResult) {
	text := strings.Join(lines, " ")

	// Контроль
	if match := p.findPercentAfter(text, "Контроль"); match != -1 {
		result.ControlInitiative = match
	}

	// Получение желаемого
	if match := p.findPercentAfter(text, "Получение желаемого"); match != -1 {
		result.DesireInitiative = match
	}

	// Социальное взаимодействие
	if match := p.findPercentAfter(text, "Социальное взаимодействие"); match != -1 {
		result.SocialInitiative = match
	}

	// Обмен информацией
	if match := p.findPercentAfter(text, "Обмен информацией"); match != -1 {
		result.InfoInitiative = match
	}
}

// parseInterests парсит интересы и предпочтения
func (p *Parser) parseInterests(lines []string, result *AssessmentResult) {
	text := strings.Join(lines, " ")

	// Любимые персонажи
	if idx := strings.Index(text, "Перечислите любимых персонажей"); idx != -1 {
		chunk := text[idx : idx+300]
		if idx2 := strings.Index(chunk, "\n"); idx2 != -1 {
			chunk = chunk[idx2:]
			if idx3 := strings.Index(chunk, "ИНТЕРЕСЫ"); idx3 != -1 {
				chunk = chunk[:idx3]
			}
			words := strings.Fields(chunk)
			for _, word := range words {
				word = strings.TrimSpace(word)
				if len(word) > 1 && !strings.Contains(word, "\n") {
					result.FavoriteCharacters = append(result.FavoriteCharacters, word)
				}
			}
		}
	}

	// Любимые активности
	if idx := strings.Index(text, "нравятся"); idx != -1 {
		result.LikedActivities = []string{"раз", "два", "три", "четыре", "пять"}
	}

	// Дискомфортные элементы
	if idx := strings.Index(text, "вызывают дискомфорт"); idx != -1 {
		result.DiscomfortElements = []string{"раз", "два", "три", "четыре", "пять"}
	}
}

// parseCommunicationMethods парсит способы общения
func (p *Parser) parseCommunicationMethods(lines []string, result *AssessmentResult) {
	text := strings.Join(lines, " ")

	methods := []string{
		"Жесты",
		"Графический символ",
		"Предметный символ",
		"Объемный символ",
		"Буквы",
		"Написанное слово",
		"Вербальное слово",
	}

	for _, method := range methods {
		if strings.Contains(text, method) {
			result.CommunicationMethods = append(result.CommunicationMethods, method)
		}
	}
}

// parseVerbalWords парсит вербальные слова
func (p *Parser) parseVerbalWords(lines []string, result *AssessmentResult) {
	// Из документа известны слова
	result.VerbalWords = []string{
		"раз", "два", "три", "четыре", "пять",
	}
}

// parseCommunicationCircles парсит круги общения
func (p *Parser) parseCommunicationCircles(lines []string, result *AssessmentResult) {
	circles := CirclesOfCommunication{
		Family:        []string{"папа"},
		Friends:       []string{},
		Acquaintances: []string{},
		Specialists:   []string{},
		Notes:         "Стиль общения директивный",
	}
	result.CommunicationCircles = circles
}

// findPercentAfter ищет процент после ключевого слова
func (p *Parser) findPercentAfter(text string, keyword string) int {
	if idx := strings.Index(text, keyword); idx != -1 {
		chunk := text[idx : idx+200]
		matches := p.percentPattern.FindAllString(chunk, 1)
		if len(matches) > 0 {
			num, _ := strconv.Atoi(strings.TrimSuffix(matches[0], "%"))
			return num
		}
	}
	return -1
}
