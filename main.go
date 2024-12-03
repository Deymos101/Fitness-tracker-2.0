package main

import (
	"fmt"
	"math"
	"time"
)

// Общие константы для вычислений.
const (
	MInKm      = 1000 // количество метров в одном километре
	MinInHours = 60   // количество минут в одном часе
	LenStep    = 0.65 // длина одного шага
	CmInM      = 100  // количество сантиметров в одном метре
)

// Training общая структура для всех тренировок
type Training struct {
	

// distance возвращает дистанцию, которую преодолел пользователь.
// Формула расчета:
// количество_повторов * длина_шага / м_в_км
func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / MInKm
}

// meanSpeed возвращает среднюю скорость бега или ходьбы.
func (t Training) meanSpeed() float64 {
	hours := t.Duration.Hours()
	if hours == 0 {
		return 0
	}
	return t.distance() / hours
}

// Calories возвращает количество потраченных килокалорий на тренировке.
// Пока возвращаем 0, так как этот метод будет переопределяться для каждого типа тренировки.
func (t Training) Calories() float64 {
	return 0 
}

// InfoMessage содержит информацию о проведенной тренировке.
type InfoMessage struct {
	TrainingType string  // тип тренировки
	Duration     time.Duration
	Distance     float64 // расстояние в км
	Speed        float64 // средняя скорость в км/ч
	Calories     float64 // потраченные калории
}

// TrainingInfo возвращает труктуру InfoMessage, в которой хранится вся информация о проведенной тренировке.
func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
}

// String возвращает строку с информацией о проведенной тренировке.
func (i InfoMessage) String() string {
	return fmt.Sprintf("Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n",
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

// CaloriesCalculator интерфейс для структур: Running, Walking и Swimming.
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// Константы для расчета потраченных килокалорий при беге.
const (
	CaloriesMeanSpeedMultiplier = 18   // множитель средней скорости бега
	CaloriesMeanSpeedShift      = 1.79 // коэффициент изменения средней скорости
)

// Running структура, описывающая тренировку Бег.
type Running struct {
	Training
}


// Calories возввращает количество потраченных килокалория при беге.
// Формула расчета:
// ((18 * средняя_скорость_в_км/ч + 1.79) * вес_спортсмена_в_кг / м_в_км * время_тренировки_в_часах * мин_в_часе)
// Это переопределенный метод Calories() из Training.
func (r Running) Calories() float64 {
	speed := r.meanSpeed()
	hours := r.Duration.Hours()
	return ((CaloriesMeanSpeedMultiplier*speed + CaloriesMeanSpeedShift) * r.Weight / MInKm) * hours * MinInHours
}

// TrainingInfo возвращает структуру InfoMessage с информацией о проведенной тренировке.
func (r Running) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: r.TrainingType, 
		Duration:     r.Duration,     
		Distance:     r.distance(),   
		Speed:        r.meanSpeed(),  
		Calories:     r.Calories(),   
	}
}

// Константы для расчета потраченных килокалорий при ходьбе.
const (
	CaloriesWeightMultiplier      = 0.035 // коэффициент для веса
	CaloriesSpeedHeightMultiplier = 0.029 // коэффициент для роста
	KmHInMsec                     = 0.278 // коэффициент для перевода км/ч в м/с
)

// Walking структура описывающая тренировку Ходьба
type Walking struct {
	Training
	Height float64 // рост пользователя в сантиметрах
}

// Calories возвращает количество потраченных килокалорий при ходьбе.
// Формула расчета:
// ((0.035 * вес_спортсмена_в_кг + (средняя_скорость_в_метрах_в_секунду**2 / рост_в_метрах)
// * 0.029 * вес_спортсмена_в_кг) * время_тренировки_в_часах * мин_в_ч)
// Это переопределенный метод Calories() из Training.
func (w Walking) Calories() float64 {
	// Перевод средней скорости из км/ч в м/с
	speedMps := w.meanSpeed() * KmHInMsec
	// Конвертация роста в метры
	heightMeters := w.Height / CmInM
	// Длительность тренировки в часах
	hours := w.Duration.Hours()
	// Формула расчета калорий
	return ((CaloriesWeightMultiplier * w.Weight) +
		((math.Pow(speedMps, 2) / heightMeters) * CaloriesSpeedHeightMultiplier * w.Weight)) *
		hours * MinInHours
}

// TrainingInfo возвращает структуру InfoMessage с информацией о проведенной тренировке.
// Это переопределенный метод TrainingInfo() из Training.
func (w Walking) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: w.TrainingType,
		Duration:     w.Duration,
		Distance:     w.distance(),
		Speed:        w.meanSpeed(),
		Calories:     w.Calories(),
	}
}

// Константы для расчета потраченных килокалорий при плавании.
const (
	SwimmingLenStep                  = 1.38 // длина одного гребка
	SwimmingCaloriesMeanSpeedShift   = 1.1  // коэффициент изменения средней скорости
	SwimmingCaloriesWeightMultiplier = 2    // множитель веса пользователя
)

// Swimming структура, описывающая тренировку Плавание
type Swimming struct {
	Training
	LengthPool float64 // Длина бассейна в метрах
	CountPool  int     // Количество пересечений бассейна
}

// meanSpeed возвращает среднюю скорость при плавании.
// Формула расчета:
// длина_бассейна * количество_пересечений / м_в_км / продолжительность_тренировки
func (s Swimming) meanSpeed() float64 {
	totalDistance := s.LengthPool * float64(s.CountPool)
	return (totalDistance / MInKm) / s.Duration.Hours()
}

// Calories возвращает количество калорий, потраченных при плавании.
// Формула расчета:
// (средняя_скорость_в_км/ч + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * вес_спортсмена_в_кг * время_тренировки_в_часах
func (s Swimming) Calories() float64 {
	speed := s.meanSpeed()
	hours := s.Duration.Hours()
	return (speed + SwimmingCaloriesMeanSpeedShift) * SwimmingCaloriesWeightMultiplier * s.Weight * hours
}

// TrainingInfo возвращает структуру InfoMessage с информацией о проведенной тренировке.
func (s Swimming) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: s.TrainingType, // Тип тренировки
		Duration:     s.Duration,     // Длительность тренировки
		Distance:     s.LengthPool * float64(s.CountPool) / MInKm, // Пройденная дистанция
		Speed:        s.meanSpeed(),  // Средняя скорость
		Calories:     s.Calories(),   // Потраченные калории
	}
}
// ReadData возвращает информацию о проведенной тренировке.
func ReadData(training CaloriesCalculator) string {
	// получите количество затраченных калорий
	calories := training.Calories()

	// получите информацию о тренировке
	info := training.TrainingInfo()
	// добавьте полученные калории в структуру с информацией о тренировке
	info.Calories = calories

	return fmt.Sprint(info)
}

func main() {

	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep,
			Duration:     90 * time.Minute,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  5,
	}

	fmt.Println(ReadData(swimming))

	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute,
			Weight:       85,
		},
		Height: 185,
	}

	fmt.Println(ReadData(walking))

	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       85,
		},
	}

	fmt.Println(ReadData(running))

}
