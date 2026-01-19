# Диаграмма состояний микроволновой печи

## Обзор
Диаграмма описывает семь состояний микроволновки (Idle, Power On, Setting Time, Setting Power, Cooking, Paused, Finished) и переходы между ними, управляемые нажатиями кнопок и событиями.

## Семь состояний

**1. Idle (Ожидание)**
- Начальное состояние
- Дисплей показывает текущее время (displayClock)
- Переход в Power On при powerPressed

**2. Power On (Включена)**
- При входе: отображение меню (displayMenu)
- Принимает ввод (acceptInput)
- Переходит в Setting Time при timeButtonPressed
- Переходит в Setting Power при powerButtonPressed
- Переходит в Cooking при startPressed [doorClosed]
- Возврат в Idle при powerPressed

**3. Setting Time (Установка времени)**
- При входе: дисплей мигает (blinkTime)
- Принимает ввод времени (acceptTimeInput)
- Возврат в Power On при confirmPressed

**4. Setting Power (Установка мощности)**
- При входе: мигает уровень мощности (blinkPowerLevel)
- Принимает ввод уровня мощности (acceptPowerLevelInput)
- Возврат в Power On при confirmPressed

**5. Cooking (Приготовление)**
- При входе: включается свет (lightOn), запускается магнетрон (startMagnetron)
- Вращение тарелки (rotatePlate), обратный отсчёт таймера (countdownTimer)
- При выходе: выключается свет (lightOff), магнетрон останавливается (stopMagnetron)
- Переход в Paused при pausePressed
- Переход в Paused при doorOpened
- Переход в Finished при timerExpired
- Переход в Power On при cancelPressed

**6. Paused (Пауза)**
- При входе: звуковой сигнал (doBeep), остановка таймера (stopTimer)
- Дисплей показывает оставшееся время (showTimer)
- Переход в Cooking при startPressed [doorClosed]
- Переход в Power On при cancelPressed

**7. Finished (Завершено)**
- При входе: звуковой сигнал (doBeep)
- Дисплей мигает (blinkDisplay)
- Переход в Power On при pressAnyKey
- Переход в Idle при powerPressed
- Переход в Idle при doorOpened

## События и условия переходов

**Кнопки:**
- powerPressed
- timeButtonPressed
- powerButtonPressed
- startPressed (с условием [doorClosed])
- pausePressed
- cancelPressed
- confirmPressed
- pressAnyKey
- doorOpened

**События таймера:**
- timerExpired

## Основные переходы

- Idle <-> Power On (powerPressed)
- Power On -> Setting Time (timeButtonPressed)
- Power On -> Setting Power (powerButtonPressed)
- Power On -> Cooking (startPressed при doorClosed)
- Cooking -> Paused (pausePressed / startPressed при doorClosed)
- Cooking -> Finished (timerExpired)
- Cooking -> Power On (cancelPressed)
- Cooking -> Paused (doorOpened)
- Paused -> Power On (cancelPressed)
- Finished -> Idle (pressAnyKey)
- Finished -> Power On (powerPressed)
