import type { ReceiptListItem } from '~/types/api'

/** Локальный каталог чеков ФНС (~50 000 ₽, 13 покупок). */
export const fnsReceiptCatalog: ReceiptListItem[] = [
  {
    id: 'fns-20260529-pyaterochka',
    store: 'Пятёрочка',
    amount: 2_847,
    date: '2026-05-29',
    category: 'Продукты',
    source: 'fns',
    items: [
      { name: 'Молоко 3,2%', price: 89, quantity: 2 },
      { name: 'Хлеб бородинский', price: 54, quantity: 2 },
      { name: 'Яйца С1', price: 119, quantity: 2 },
      { name: 'Сыр российский', price: 249, quantity: 3 },
      { name: 'Курица охлаждённая', price: 389, quantity: 1 }
    ]
  },
  {
    id: 'fns-20260528-apteka',
    store: 'Аптека 36,6',
    amount: 1_890,
    date: '2026-05-28',
    category: 'Здоровье',
    source: 'fns',
    items: [
      { name: 'Витамин D3', price: 620, quantity: 1 },
      { name: 'Омега-3', price: 890, quantity: 1 },
      { name: 'Пластырь', price: 380, quantity: 1 }
    ]
  },
  {
    id: 'fns-20260527-coffee',
    store: 'Surf Coffee',
    amount: 620,
    date: '2026-05-27',
    category: 'Кафе',
    source: 'fns',
    impulse_count: 1,
    items: [{ name: 'Капучино 400 мл', price: 620, quantity: 1, impulse: true }]
  },
  {
    id: 'fns-20260526-ozon',
    store: 'Ozon',
    amount: 4_290,
    date: '2026-05-26',
    category: 'Развлечения',
    source: 'fns',
    items: [
      { name: 'Наушники беспроводные', price: 3_490, quantity: 1 },
      { name: 'Чехол для телефона', price: 800, quantity: 1 }
    ]
  },
  {
    id: 'fns-20260525-vkusvill',
    store: 'ВкусВилл',
    amount: 2_193,
    date: '2026-05-25',
    category: 'Продукты',
    source: 'fns',
    items: [
      { name: 'Салат Цезарь', price: 349, quantity: 2 },
      { name: 'Смузи манго', price: 199, quantity: 3 },
      { name: 'Хумус', price: 249, quantity: 2 }
    ]
  },
  {
    id: 'fns-20260524-metro',
    store: 'Метро',
    amount: 5_670,
    date: '2026-05-24',
    category: 'Продукты',
    source: 'fns',
    items: [
      { name: 'Куриное филе', price: 389, quantity: 4 },
      { name: 'Рис', price: 149, quantity: 2 },
      { name: 'Овощи микс', price: 199, quantity: 5 },
      { name: 'Сыр Гауда', price: 449, quantity: 2 }
    ]
  },
  {
    id: 'fns-20260523-dns',
    store: 'DNS',
    amount: 10_990,
    date: '2026-05-23',
    category: 'Развлечения',
    source: 'fns',
    items: [
      { name: 'SSD 1 ТБ', price: 8_990, quantity: 1 },
      { name: 'Кабель USB-C', price: 990, quantity: 1 },
      { name: 'Коврик для мыши', price: 1_010, quantity: 1 }
    ]
  },
  {
    id: 'fns-20260522-transport',
    store: 'Яндекс Go',
    amount: 890,
    date: '2026-05-22',
    category: 'Транспорт',
    source: 'fns',
    items: [{ name: 'Поездка', price: 890, quantity: 1 }]
  },
  {
    id: 'fns-20260521-lenta',
    store: 'Лента',
    amount: 3_450,
    date: '2026-05-21',
    category: 'Продукты',
    source: 'fns',
    items: [
      { name: 'Говядина тушёная', price: 520, quantity: 2 },
      { name: 'Макароны', price: 89, quantity: 4 },
      { name: 'Томаты черри', price: 199, quantity: 3 },
      { name: 'Сок апельсиновый', price: 149, quantity: 4 }
    ]
  },
  {
    id: 'fns-20260520-sportmaster',
    store: 'Спортмастер',
    amount: 4_590,
    date: '2026-05-20',
    category: 'Развлечения',
    source: 'fns',
    items: [
      { name: 'Кроссовки беговые', price: 3_990, quantity: 1 },
      { name: 'Носки спортивные', price: 600, quantity: 1 }
    ]
  },
  {
    id: 'fns-20260519-shell',
    store: 'Shell',
    amount: 3_200,
    date: '2026-05-19',
    category: 'Транспорт',
    source: 'fns',
    items: [
      { name: 'АИ-95', price: 2_800, quantity: 1 },
      { name: 'Вода 0,5 л', price: 400, quantity: 1 }
    ]
  },
  {
    id: 'fns-20260518-wb',
    store: 'Wildberries',
    amount: 2_780,
    date: '2026-05-18',
    category: 'Другое',
    source: 'fns',
    items: [
      { name: 'Футболка базовая', price: 990, quantity: 2 },
      { name: 'Ремень кожаный', price: 800, quantity: 1 }
    ]
  },
  {
    id: 'fns-20260517-mvideo',
    store: 'М.Видео',
    amount: 8_590,
    date: '2026-05-17',
    category: 'Развлечения',
    source: 'fns',
    items: [
      { name: 'Пылесос вертикальный', price: 7_990, quantity: 1 },
      { name: 'Мешки для пылесоса', price: 600, quantity: 1 }
    ]
  }
]
