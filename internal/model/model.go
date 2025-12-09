package model

import "time"

// AssessmentDocument - весь документ оценки
type AssessmentDocument struct {
	Metadata       AssessmentMetadata   `json:"metadata"`
	BasicInfo      BasicInfo            `json:"basic_info"`
	Activities     []ActivityAssessment `json:"activities"`      // newAct*
	ActivityBlocks []ActivityBlock      `json:"activity_blocks"` // actBlock*
	DiagramSummary DiagramSummary       `json:"diagram_summary"` // diagramBlock
	Dictionaries   Dictionaries         `json:"dictionaries"`
	Rounds         []Round              `json:"rounds"`
	LevelZone      string               `json:"level_zone"` // levelZoneZbr
}

// AssessmentMetadata - метаинформация о документе
type AssessmentMetadata struct {
	StudentID      string    `json:"student_id"` // из por02 или имени файла
	AssessmentType string    `json:"type"`       // "before" или "after"
	Date           time.Time `json:"date"`       // из por01
	NextDate       time.Time `json:"next_date"`  // из por03
	Location       string    `json:"location"`   // из por05
	Specialist     string    `json:"specialist"` // из por06
}

// BasicInfo - базовые данные об ученике
type BasicInfo struct {
	Speech        string            `json:"speech"`  // Baz01
	Hearing       string            `json:"hearing"` // Baz04
	Vision        string            `json:"vision"`  // Baz03
	Communication CommunicationInfo `json:"communication"`
}

type CommunicationInfo struct {
	Method       string `json:"method"`        // ADK01
	Access       string `json:"access"`        // ADK02
	GetMethod    string `json:"get_method"`    // ADK03
	Activities   string `json:"activities"`    // ADK04
	Gaze         string `json:"gaze"`          // Baz07
	ResponseType string `json:"response_type"` // Baz05
}

// ActivityAssessment - оценка конкретного навыка
type ActivityAssessment struct {
	ID       string `json:"id"`        // newAct05
	Title    string `json:"title"`     // Название навыка
	Answer   string `json:"answer"`    // Текущий статус
	Passed   bool   `json:"passed"`    // Навык освоен?
	TextNote string `json:"text_note"` // Дополнительные заметки
	Progress int    `json:"progress"`  // % прогресса (из newAct01-04)
}

// ActivityBlock - блок активности с метриками
type ActivityBlock struct {
	ID                 string               `json:"id"`                 // actBlock01
	IsAvailable        bool                 `json:"is_available"`       // bodyBlockElem == "shown"
	Protocommunication CommunicationMetrics `json:"protocommunication"` // prot*
	Voice              CommunicationMetrics `json:"voice"`              // gol*
	Phrase             CommunicationMetrics `json:"phrase"`             // fra*
	Targets            ActivityTargets      `json:"targets"`            // *Target поля
}

// CommunicationMetrics - метрики для каждого типа коммуникации
type CommunicationMetrics struct {
	FormedPercent      int    `json:"formed_percent"`    // *SforProcElem
	WithSupportPercent int    `json:"support_percent"`   // *InitProcElem
	FrequencyText      string `json:"frequency_text"`    // *ChastTextElem
	FrequencyPercent   int    `json:"frequency_percent"` // *ChastProcElem
	Description        string `json:"description"`       // *OpisanieElem
	IsUnavailable      bool   `json:"is_unavailable"`    // *UnElem
	Zone               string `json:"zone"`              // ЗБР/ЗАР
}

// ActivityTargets - целевые навыки
type ActivityTargets struct {
	ProtRequest   string `json:"prot_request"`   // protZorTarget
	ProtRefuse    string `json:"prot_refuse"`    // protZbrTarget
	VoiceRequest  string `json:"voice_request"`  // golZorTarget
	VoiceRefuse   string `json:"voice_refuse"`   // golZbrTarget
	PhraseRequest string `json:"phrase_request"` // fraZorTarget
	PhraseRefuse  string `json:"phrase_refuse"`  // fraZbrTarget
}

// DiagramSummary - итоговые показатели
type DiagramSummary struct {
	Predictive         MetricsPair `json:"predictive"`         // pred*
	Protocommunication MetricsPair `json:"protocommunication"` // prot*
	Voice              MetricsPair `json:"voice"`              // gol*
	Phrase             MetricsPair `json:"phrase"`             // fra*
}

type MetricsPair struct {
	ActivePercent  int    `json:"active_percent"`  // *ActProcNumElem
	SupportPercent int    `json:"support_percent"` // *InitProcNumElem
	Zone           string `json:"zone"`            // *ZoneElem (ЗБР/ЗАР)
}

// Dictionaries - словари
type Dictionaries struct {
	Basic         []DictionaryItem `json:"basic"`         // basicDictionary
	BasicMore     []string         `json:"basic_more"`    // dictBasicMore
	Communication []string         `json:"communication"` // dictSposObsh
	Verbal        []string         `json:"verbal"`        // dictWerbSlov
	FastMessages  []string         `json:"fast_messages"` // dictBystrSoobsh
}

type DictionaryItem struct {
	Content    string `json:"content"`
	ColorStyle string `json:"color_style"`
	IsActive   bool   `json:"is_active"` // itemOffStyle == ""
}

// Round - информация о человеке из социального окружения ученика
type Round struct {
	Relation     string   `json:"relation"`      // roundRelation: сестра, мама, папа
	Name         string   `json:"name"`          // roundName: имя
	StyleText    string   `json:"style_text"`    // roundStyleText: стиль общения
	StyleAlert   int      `json:"style_alert"`   // roundStyleAlert: 0-2 (уровень внимания)
	ReactText    string   `json:"react_text"`    // roundReactText: как реагирует
	ReactAlert   int      `json:"react_alert"`   // roundReactAlert: 0-2
	SymbText     string   `json:"symb_text"`     // roundSymbText: использование символов
	SymbAlert    int      `json:"symb_alert"`    // roundSymbAlert: 0-2
	HowInteracts []string `json:"how_interacts"` // roundHowTextArr: способы взаимодействия
}
