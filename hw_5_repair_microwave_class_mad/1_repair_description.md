# Диаграмма последовательности: регистрация и ремонт дефекта

## Обзор
Диаграмма описывает полный цикл работы системы от начала смены до формирования отчётов, включая регистрацию дефекта, назначение ремонта, его выполнение и верификацию.

## Фаза 1: Подготовка к работе (начало смены)

Механик регистрируется в мобильном приложении в начале смены:
1. Mobile отправляет POST /api/shifts/checkin с данными механика и зоны
2. API Gateway передаёт запрос в Resources Service
3. Resources Service обновляет MechanicShift в БД (present=true, availableForJobs=true)
4. Статус механика обновляется в кэше Redis
5. Механик готов принимать задания

## Фаза 2: Регистрация дефекта и диагностика

**Регистрация:**
1. Механик сканирует VIN автомобиля в планшете
2. Mobile отправляет POST /api/defects с параметрами: vehicleId, stage, area, locationOnDiagram, probableCause, severity
3. Defect Service сохраняет дефект в БД со статусом "Registered"
4. API возвращает defectId

**Углубленная диагностика:**
1. Механик проводит детальную диагностику
2. Mobile отправляет POST /api/defects/{id}/diagnosis
3. Defect Service обновляет дефект (status=Diagnosed)
4. Публикует событие DefectDiagnosed в Message Broker для синхронизации других сервисов

## Фаза 3: Направление в ремонтную зону

**Проверка доступности (диспетчер):**
1. Dispatcher запрашивает GET /api/repair-zones/availability с фильтром по механикам
2. Resources Service получает из Redis информацию о свободных местах и доступных механиках
3. Dispatcher видит список зон и свободные места

**Создание заказа на ремонт:**
1. Dispatcher отправляет POST /api/repair-orders с указанием defectId, zoneId, bayId, priority
2. Repair Service создаёт RepairOrder (status=Planned)
3. Occupy Bay - место помечается как OCCUPIED в кэше
4. Публикуется событие RepairOrderCreated

## Фаза 4: Назначение ремонта бригадиром

1. Foreman запрашивает список доступных механиков для зоны
2. Resources Service возвращает механиков, у которых availableForJobs=true в текущую смену
3. Foreman отправляет POST /api/repair-orders/{id}/assign с mechanicId
4. Repair Service обновляет RepairOrder (assignedTo, status=Assigned)
5. Механик помечается как BUSY в кэше
6. Публикуется событие RepairOrderAssigned

## Фаза 5: Выполнение ремонта

**Начало ремонта:**
1. Механик в Mobile нажимает "Начало работы"
2. POST /api/repair-orders/{id}/start
3. Repair Service обновляет status=InProgress, startAt=now()
4. Создаётся запись WorkLog с начальным временем
5. Публикуется событие RepairStarted

**Завершение ремонта:**
1. Механик вводит результаты ремонта (actions, parts used)
2. POST /api/repair-orders/{id}/complete
3. Repair Service обновляет RepairOrder (status=Completed, endAt=now())
4. Место освобождается (status=FREE)
5. Механик помечается как AVAILABLE в кэше
6. WorkLog завершается с расчётом длительности
7. Публикуется событие RepairCompleted

## Фаза 6: Верификация ОТК (инженер контроля качества)

1. QC запрашивает список заказов, ожидающих верификации
2. Repair Service возвращает заказы со статусом Completed

**Если ремонт принят:**
- RepairOrder обновляется (status=Verified)
- Defect закрывается (status=Closed)
- Публикуется событие RepairVerified

**Если ремонт отклонен:**
- RepairOrder обновляется (status=Rejected, rejectionReason)
- Defect переходит в статус Reopened
- Публикуется событие RepairRejected (может потребоваться повторный ремонт)

## Фаза 7: Формирование отчётов (конец смены)

**Отчёт бригадира:**
- GET /api/reports/shift/brigade
- Агрегация: количество ремонтов по каждому механику, время работы, статистика

**Отчёт диспетчера (смены):**
- GET /api/reports/shift/summary
- Общие метрики по зоне + статистика по местам дефектов

**Отчёт ОТК:**
- GET /api/reports/shift/qc
- Статистика верификации, среднее время ремонта

API возвращает отчёты в формате PDF/Excel/CSV.

## Ключевые технологические потоки

- **REST API** - синхронные запросы между клиентами и сервисами
- **Message Broker** - асинхронная публикация событий для развязки сервисов
- **Redis кэш** - быстрые обновления статусов мест и механиков (реал-тайм для диспетчера)
- **БД** - персистентное хранение всех изменений
- **Работа со статусами** - чёткий конечный автомат для RepairOrder и Defect
