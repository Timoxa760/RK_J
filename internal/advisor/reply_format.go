package advisor

import (
	"regexp"
	"strings"
)

var (
	rePipeBlock       = regexp.MustCompile(`(?:\|[^|\n]{1,120}){2,}`)
	reLabelBreak      = regexp.MustCompile(`([.!?])\s*([А-ЯЁ][^\n:]{0,40}:)`)
	reHeaderInline    = regexp.MustCompile(`([^\n])(#{1,6}\s)`)
	reHrInline        = regexp.MustCompile(`([^\n])\s*---\s*([^\n])`)
	reSpaceAfterPunct = regexp.MustCompile(`([,.:;!?])([^\s\d|])`)
	reSpaceRub        = regexp.MustCompile(`₽\s*([А-ЯA-Z])`)
	reCyrBoundary     = regexp.MustCompile(`([а-яё])([А-ЯЁ])`)
	reDashNumber      = regexp.MustCompile(`—\s*(\d)`)
	reNumRub          = regexp.MustCompile(`(\d)\s*₽`)
	rePrepNum         = regexp.MustCompile(`(?i)(на|по|за|от|до)\s*(\d)`)
	reMultiPipe       = regexp.MustCompile(`\|{2,}`)
	rePipeRow         = regexp.MustCompile(`(?m)^[^|\n]*\|[^|\n]+`)
	reOrphanLine      = regexp.MustCompile(`(?m)^[\s|—-]+$`)
	reMultiSpace      = regexp.MustCompile(`[ \t]{2,}`)
)

// Явные склейки от LLM — не разбиваем нормальные слова вроде «расходы», «указать».
var gluedRepairs = []*regexp.Regexp{
	regexp.MustCompile(`(?i)притекущих\s*данных`),
	regexp.MustCompile(`(?i)притекущих`),
	regexp.MustCompile(`(?i)свободных\s*денег`),
	regexp.MustCompile(`(?i)свободныхденег`),
	regexp.MustCompile(`(?i)откладывать\s*на\s*цель`),
	regexp.MustCompile(`(?i)откладыватьнацель`),
	regexp.MustCompile(`(?i)накредит`),
	regexp.MustCompile(`(?i)главный\s*вопрос`),
	regexp.MustCompile(`(?i)главныйвопрос`),
	regexp.MustCompile(`(?i)илиесть`),
	regexp.MustCompile(`(?i)замесяц`),
	regexp.MustCompile(`(?i)занеделю`),
	regexp.MustCompile(`(?i)вмесяц`),
	regexp.MustCompile(`(?i)намесяц`),
	regexp.MustCompile(`(?i)безкатегори`),
	regexp.MustCompile(`(?i)безпробел`),
	regexp.MustCompile(`(?i)([а-яё]{3,})(данных|денег)\b`),
	regexp.MustCompile(`(?i)([а-яё]+)(на)(цель|кредит)\b`),
	regexp.MustCompile(`(\d)(₽)`),
	regexp.MustCompile(`(\d)(%)`),
}

var gluedReplacements = []string{
	"при текущих данных",
	"при текущих",
	"свободных денег",
	"свободных денег",
	"Откладывать на цель",
	"Откладывать на цель",
	"на кредит",
	"Главный вопрос",
	"Главный вопрос",
	"или есть",
	"за месяц",
	"за неделю",
	"в месяц",
	"на месяц",
	"без категори",
	"без пробел",
	"$1 $2",
	"$1 $2 $3",
	"$1 $2",
	"$1 $2",
}

// FormatAdvisorReply приводит ответ LLM к читаемому виду (пробелы, абзацы, без таблиц).
func FormatAdvisorReply(text string) string {
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.TrimSpace(text)
	if text == "" {
		return ""
	}

	text = reMultiPipe.ReplaceAllString(text, "\n")
	text = strings.ReplaceAll(text, "|---|", "\n")
	text = convertPipeBlocks(text)
	text = convertPipeRows(text)
	text = reLabelBreak.ReplaceAllString(text, "$1\n\n$2")
	text = reHeaderInline.ReplaceAllString(text, "$1\n\n$2")
	text = reHrInline.ReplaceAllString(text, "$1\n\n---\n\n$2")
	text = applyGluedRepairs(text)
	text = repairSpacing(text)
	text = reOrphanLine.ReplaceAllString(text, "")

	lines := strings.Split(text, "\n")
	clean := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(repairSpacing(line))
		if line == "" || line == "|" || line == "---" {
			continue
		}
		clean = append(clean, line)
	}
	text = strings.TrimSpace(strings.Join(clean, "\n"))
	text = reMultiSpace.ReplaceAllString(text, " ")
	text = RepairSplitRussianWords(text)
	return text
}

func applyGluedRepairs(s string) string {
	for i, re := range gluedRepairs {
		s = re.ReplaceAllString(s, gluedReplacements[i])
	}
	return s
}

func convertPipeBlocks(text string) string {
	return rePipeBlock.ReplaceAllStringFunc(text, func(block string) string {
		return pipeCellsToList(block)
	})
}

func convertPipeRows(text string) string {
	return rePipeRow.ReplaceAllStringFunc(text, func(row string) string {
		if strings.Count(row, "|") < 1 {
			return row
		}
		return pipeCellsToList(row)
	})
}

func pipeCellsToList(block string) string {
	cells := strings.Split(block, "|")
	var rows []string
	for _, cell := range cells {
		cell = strings.TrimSpace(cell)
		if cell == "" || cell == "---" || strings.HasPrefix(cell, "---") {
			continue
		}
		cell = repairSpacing(cell)
		if cell != "" {
			rows = append(rows, "- "+cell)
		}
	}
	if len(rows) == 0 {
		return repairSpacing(strings.ReplaceAll(block, "|", " "))
	}
	return strings.Join(rows, "\n")
}

func repairSpacing(s string) string {
	s = reSpaceAfterPunct.ReplaceAllString(s, "$1 $2")
	s = strings.ReplaceAll(s, "—", " — ")
	s = reDashNumber.ReplaceAllString(s, "— $1")
	s = reNumRub.ReplaceAllString(s, "$1 ₽")
	s = rePrepNum.ReplaceAllString(s, "$1 $2")
	s = reSpaceRub.ReplaceAllString(s, "₽ $1")
	s = reCyrBoundary.ReplaceAllString(s, "$1 $2")
	s = reMultiSpace.ReplaceAllString(s, " ")
	return strings.TrimSpace(s)
}
