import assert from 'node:assert/strict'
import { describe, it } from 'node:test'
import { repairSplitRussianWords } from './advisorMarkdown.ts'

describe('repairSplitRussianWords', () => {
  it('does not split normal Russian words with embedded na/za/do', () => {
    const words = [
      'необязательных',
      'обязательных',
      'Проанализируйте',
      'проанализируйте',
      'радость',
      'альтернативами',
      'анализ',
      'законодательных',
      'подготовка',
      'сегодня'
    ]
    for (const word of words) {
      assert.equal(repairSplitRussianWords(word), word)
    }
  })

  it('preserves user example paragraph', () => {
    const input =
      'Развлечения и рестораны составляют большую часть ваших необязательных расходов. ' +
      'Сейчас вы тратите 27 660 ₽ на развлечения и 520 ₽ на кафе и рестораны — в сумме 28 180 ₽ в месяц. ' +
      'Сократив эти категории минимум на 50 %, вы сможете освободить около 14 090 ₽.'
    assert.equal(repairSplitRussianWords(input), input)
  })

  it('preserves list items from advisor reply', () => {
    const items = [
      'Проанализируйте, какие развлечения приносят наибольшую радость, и замените их более дешёвыми альтернативами (домашние киносеансы, бесплатные мероприятия).',
      'Ограничьте посещения кафе и ресторанов до одного раза в неделю и выбирайте более доступные заведения.',
      'Установите ежемесячный лимит на развлечения в размере 13 830 ₽ (половина текущих расходов).',
      'Составьте список обязательных развлечений и замените остальные более экономными вариантами уже в этом месяце.'
    ]
    for (const item of items) {
      assert.equal(repairSplitRussianWords(item), item)
    }
  })

  it('fixes known LLM split artifacts', () => {
    assert.equal(repairSplitRussianWords('де по зит'), 'депозит')
    assert.equal(repairSplitRussianWords('го до вых'), 'годовых')
    assert.equal(repairSplitRussianWords('по это му'), 'поэтому')
    assert.equal(repairSplitRussianWords('фи на нсовый'), 'финансовый')
    assert.equal(repairSplitRussianWords('меся -'), 'месяце')
  })

  it('fixes known glued tokens via allowlist', () => {
    assert.equal(repairSplitRussianWords('расходынапродукты'), 'расходы на продукты')
    assert.equal(repairSplitRussianWords('бюджетнаследующий'), 'бюджет на следующий')
    assert.equal(repairSplitRussianWords('откладыватьнацель'), 'откладывать на цель')
    assert.equal(repairSplitRussianWords('Сокративэти'), 'Сократив эти')
  })

  it('does not re-break text on second pass', () => {
    const once = repairSplitRussianWords('необязательных расходов')
    assert.equal(repairSplitRussianWords(once), once)
  })
})
