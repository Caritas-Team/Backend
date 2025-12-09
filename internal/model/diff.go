package model

import "time"

// AssessmentDiff - результат сравнения "до/после"
type AssessmentDiff struct {
	StudentID          string             `json:"student_id"`
	PeriodStart        time.Time          `json:"period_start"`
	PeriodEnd          time.Time          `json:"period_end"`
	OverallProgress    OverallProgress    `json:"overall_progress"`
	ActivitiesProgress []ActivityProgress `json:"activities_progress"`
	BlocksProgress     []BlockProgress    `json:"blocks_progress"`
	DiagramDiff        DiagramDiff        `json:"diagram_diff"`
	DictionaryGrowth   DictionaryGrowth   `json:"dictionary_growth"`
	LevelChange        LevelChange        `json:"level_change"`
}

// OverallProgress - общий прогресс
type OverallProgress struct {
	NewSkillsAcquired    int      `json:"new_skills_acquired"`    // сколько навыков освоено
	SkillsImproved       int      `json:"skills_improved"`        // сколько улучшилось
	AverageProgressDelta float64  `json:"average_progress_delta"` // средний прирост %
	Recommendations      []string `json:"recommendations"`
}

// ActivityProgress - прогресс по навыку
type ActivityProgress struct {
	ActivityID    string `json:"activity_id"`
	Title         string `json:"title"`
	StatusBefore  string `json:"status_before"`
	StatusAfter   string `json:"status_after"`
	PassedBefore  bool   `json:"passed_before"`
	PassedAfter   bool   `json:"passed_after"`
	IsImproved    bool   `json:"is_improved"`
	ProgressDelta int    `json:"progress_delta"` // изменение в %
}

// BlockProgress - прогресс по блоку активности
type BlockProgress struct {
	BlockID         string               `json:"block_id"`
	ProtProgress    MetricsProgressDelta `json:"prot_progress"`
	VoiceProgress   MetricsProgressDelta `json:"voice_progress"`
	PhraseProgress  MetricsProgressDelta `json:"phrase_progress"`
	BecameAvailable bool                 `json:"became_available"`
}

type MetricsProgressDelta struct {
	FormedPercentDelta  int    `json:"formed_percent_delta"` // разница в %
	SupportPercentDelta int    `json:"support_percent_delta"`
	FrequencyImproved   bool   `json:"frequency_improved"`
	ZoneChanged         bool   `json:"zone_changed"`
	ZoneBefore          string `json:"zone_before"`
	ZoneAfter           string `json:"zone_after"`
}

// DiagramDiff - изменения итоговых показателей
type DiagramDiff struct {
	PredictiveDelta         int `json:"predictive_delta"`
	ProtocommunicationDelta int `json:"protocommunication_delta"`
	VoiceDelta              int `json:"voice_delta"`
	PhraseDelta             int `json:"phrase_delta"`
}

// DictionaryGrowth - рост словаря
type DictionaryGrowth struct {
	BasicWordsAdded   int      `json:"basic_words_added"`
	BasicWordsRemoved int      `json:"basic_words_removed"`
	TotalWordsGrowth  int      `json:"total_words_growth"`
	NewActiveWords    []string `json:"new_active_words"` // слова, ставшие активными
}

// LevelChange - изменение уровня развития
type LevelChange struct {
	LevelBefore string `json:"level_before"`
	LevelAfter  string `json:"level_after"`
	HasImproved bool   `json:"has_improved"`
}
