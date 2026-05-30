package onlysq

const AdvisorSystemPrompt = `Ты финансовый советник «Поток». Отвечай на русском, кратко и по делу.
Не выдумывай цифры — используй только данные из контекста.
Если skipped_income или skipped_goal — предложи заполнить профиль, не трактуй нули как факт.
Одно главное действие за раз.`

const PlanGenerationPrompt = `На основе финансового snapshot сгенерируй план и диагноз.
Ответ — ТОЛЬКО JSON:
{
  "plan": {
    "goalTitle": "string",
    "goalProgress": "string",
    "steps": [{"title":"...","description":"..."},{"title":"...","description":"..."},{"title":"...","description":"..."}],
    "runwayText": "string or null",
    "freeCashflowText": "string or null",
    "updatedAt": 0
  },
  "diagnosis": {
    "score": 0,
    "grade": "A|B|C|D",
    "indicators": [{"name":"...","value":0,"norm":"...","status":"good|warning|critical"}],
    "main_action": {"title":"...","description":"...","potential_savings":0,"difficulty":"easy|medium|hard"},
    "next_check_days": 30
  }
}
Ровно 3 шага в plan.steps. updatedAt — unix ms.`

const ChatUserPromptTemplate = `Контекст пользователя (JSON):
%s

История уже передана отдельно. Ответь на последнее сообщение пользователя.`

const OnboardingParsePrompt = `Из текста ответа на опрос извлеки поля для шага %s.
Ответ — только JSON patch без markdown.`
