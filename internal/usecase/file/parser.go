package file

import "github.com/Caritas-Team/reviewer/internal/model"

type DocumentParser struct{}

func (p *DocumentParser) Parse(rawJSON []byte) (*model.AssessmentDocument, error) {
	// Парсит JSON и маппит в структурированную модель
	// Убирает HTML теги сразу при парсинге
	// Возвращает валидированный документ

	return nil, nil
}

// Stage 2: Validator
type PairValidator struct{}

func (v *PairValidator) Validate(before, after *model.AssessmentDocument) error {
	// Проверяет что документы относятся к одному ученику
	// Проверяет хронологический порядок дат
	// Проверяет консистентность данных

	return nil
}

// Stage 3: Diff Calculator
type DiffCalculator struct{}

func (c *DiffCalculator) Calculate(before, after *model.AssessmentDocument) (*model.AssessmentDiff, error) {
	// Вычисляет все изменения между документами
	// Генерирует рекомендации на основе прогресса
	// Возвращает структурированный результат

	return nil, nil
}

// Stage 4: Report Generator (опционально)
type ReportGenerator struct{}

func (g *ReportGenerator) GenerateReport(diff *model.AssessmentDiff) (string, error) {
	// Генерирует человекочитаемый отчёт в формате Markdown/HTML
	// Для интеграции с фронтендом или печати
	return "", nil
}
