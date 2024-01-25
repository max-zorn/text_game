package main

import (
	"strings"
)

type Location struct {
	name              string
	condition         string
	info              string
	quest             Quest
	nearby            []*Location
	nearbyPermissions map[*string]int
	furniture         []*Furniture
	*Quest
}

type Quest struct {
	name   string
	stages []*QuestStep
}

type QuestStep struct {
	goalName string
	check    func() bool
}

type Furniture struct {
	name  string
	items []*Item
	bag   int
}

type Backpack struct {
	name  string
	items map[*string]*Item
}

type Item struct {
	name  string
	items map[string]*Item
}

type Person struct {
	condition       string
	inventoryExists int
	bag             Item
	position        Location
	*Location
	*Item
}

func checkItems(l *Location) (result bool) {
	for _, furniture := range l.furniture {
		if len(furniture.items) != 0 {
			return true
		}
	}
	return false
}

func (p *Person) lookAround() (message string) {
	var text string
	if len(p.Location.furniture) != 0 {
		if checkItems(p.Location) == true {
			indexFuniture := 0
			for _, furniture := range p.Location.furniture {
				if len(furniture.items) != 0 && indexFuniture == 0 {
					indexFuniture += 1
					text = text + "на " + furniture.name + "е: "
					indexItems := 0
					for _, item := range furniture.items {
						if indexItems == 0 {
							text = text + item.name
							indexItems += 1
						} else {
							text = text + ", " + item.name
						}
					}
				} else if len(furniture.items) != 0 {
					text = text + ", на " + furniture.name + "е: "
					indexItems := 0
					for _, item := range furniture.items {
						if indexItems == 0 {
							text = text + item.name
							indexItems += 1
						} else {
							text = text + ", " + item.name
						}
					}
				}
			}
		} else {
			text = "пустая комната"
		}
	}

	message = message + p.Location.info + text

	if p.Location.Quest != nil {
		indexQuest := 0
		for _, questStep := range p.Location.Quest.stages {
			if questStep.check() != true && indexQuest == 0 {
				message = message + ", надо " + questStep.goalName
				indexQuest += 1
			} else if questStep.check() != true {
				message = message + " и " + questStep.goalName
			}
		}
	}
	var toGo string
	index := 0
	for _, nearLocation := range p.Location.nearby {
		if index == 0 {
			toGo = nearLocation.name
			index += 1
		} else {
			toGo = toGo + ", " + nearLocation.name
		}
	}
	message = message + ". можно пройти - " + toGo
	return
}

func (p *Person) move(direction string) (message string) {
	for _, v := range p.Location.nearby {
		if direction == v.name {
			if p.Location.nearbyPermissions[&v.name] == 0 {
				p.Location = v
				var toGo string
				idx := 0
				for _, value := range p.Location.nearby {
					if idx == 0 {
						toGo = value.name
						idx += 1
					} else {
						toGo = toGo + ", " + value.name
					}
				}
				return p.Location.condition + " можно пройти - " + toGo
			} else {
				return "дверь закрыта"
			}
		}
	}
	return "нет пути в " + direction
}

func (p *Person) putBagOn() (message string) {
	for _, value := range p.Location.furniture {
		for i, _ := range value.items {
			if value.bag == 1 {
				value.bag = 0
				value.items = append(value.items[:i], value.items[i+1:]...)

				p.inventoryExists = 1
				p.bag = Bag
				return "вы надели: рюкзак"
			}
		}
	}
	return "error"
}

func (p *Person) takeItem(item string) (message string) {
	if p.bag.name == Bag.name {
		for _, v := range p.Location.furniture {
			for key, value := range v.items {
				if item == value.name {
					p.bag.items[value.name] = value

					copy(v.items[key:], v.items[key+1:])
					v.items[len(v.items)-1] = nil
					v.items = v.items[:len(v.items)-1]

					return "предмет добавлен в инвентарь: " + item
				}
			}
		}
		return "нет такого"
	}
	return "некуда класть"
}

func (p *Person) use(sourceItem, targetItem string) (message string) {
	for _, v := range p.bag.items {
		if sourceItem == v.name {
			switch targetItem {
			case "дверь":
				if p.Location.nearbyPermissions[&Street.name] == 1 {
					p.Location.nearbyPermissions[&Street.name] = 0
					return "дверь открыта"
				} else {
					return "дверь уже открыта"
				}
			default:
				return "не к чему применить"
			}
		}
	}
	return "нет предмета в инвентаре - " + sourceItem
}

func (p Person) unknownAction() (message string) {
	return "неизвестная команда"
}

var Street, Hall, Kitchen, Room, Home Location
var TableInKitchen, TableInRoom, ChairInRoom Furniture
var Tea, Keys, Summary, Bag Item
var Player Person
var KitchenQuest Quest
var FirstStep, SecondStep QuestStep

func initGame() {
	Kitchen = Location{
		name:              "кухня",
		condition:         "кухня, ничего интересного.",
		info:              "ты находишься на кухне, ",
		Quest:             &KitchenQuest,
		nearby:            []*Location{&Hall},
		nearbyPermissions: map[*string]int{&Hall.name: 0},
		furniture:         []*Furniture{&TableInKitchen},
	}

	Hall = Location{
		name:              "коридор",
		condition:         "ничего интересного.",
		info:              "",
		nearby:            []*Location{&Kitchen, &Room, &Street},
		nearbyPermissions: map[*string]int{&Kitchen.name: 0, &Room.name: 0, &Street.name: 1},
		furniture:         []*Furniture{},
	}

	Room = Location{
		name:              "комната",
		condition:         "ты в своей комнате.",
		info:              "",
		nearby:            []*Location{&Hall},
		nearbyPermissions: map[*string]int{&Hall.name: 0},
		furniture:         []*Furniture{&TableInRoom, &ChairInRoom},
	}

	Street = Location{
		name:              "улица",
		condition:         "на улице весна.",
		info:              "",
		nearby:            []*Location{&Home},
		nearbyPermissions: map[*string]int{&Home.name: 0},
		furniture:         []*Furniture{},
	}

	Home = Location{
		name:              "домой",
		condition:         " ",
		info:              "",
		nearby:            []*Location{},
		nearbyPermissions: map[*string]int{},
		furniture:         []*Furniture{},
	}

	TableInKitchen = Furniture{
		name:  "стол",
		items: []*Item{&Tea},
		bag:   0,
	}

	TableInRoom = Furniture{
		name:  "стол",
		items: []*Item{&Keys, &Summary},
		bag:   0,
	}

	ChairInRoom = Furniture{
		name:  "стул",
		items: []*Item{&Bag},
		bag:   1,
	}

	Tea = Item{
		name: "чай",
	}

	Keys = Item{
		name: "ключи",
	}

	Summary = Item{
		name: "конспекты",
	}

	Bag = Item{
		name:  "рюкзак",
		items: map[string]*Item{},
	}

	Player = Person{
		condition:       "какое-то состояние игрока",
		inventoryExists: 0,
		Location:        &Kitchen,
		// Item: 				&Bag,
	}

	FirstStep = QuestStep{
		goalName: "собрать рюкзак",
	}

	FirstStep.check = func() bool {
		checkSummary := func(item *Item) bool {
			for _, value := range item.items {
				if value == &Summary {
					return true
				}
			}
			return false
		}

		if Player.inventoryExists == 1 && checkSummary(&Bag) == true {
			return true
		}
		return false
	}

	SecondStep = QuestStep{
		goalName: "идти в универ",
	}

	SecondStep.check = func() bool {
		return false
	}

	KitchenQuest = Quest{
		name:   "квест для кухни",
		stages: []*QuestStep{&FirstStep, &SecondStep},
	}

	return
}

func handleCommand(command string) (message string) {
	if command != "" {
		action := strings.Fields(command)
		switch action[0] {
		case "осмотреться":
			message = Player.lookAround()
		case "идти":
			if len(action) < 2 {
				message = "Неправильная команда"
			} else {
				message = Player.move(action[1])
			}
		case "надеть":
			message = Player.putBagOn()
		case "взять":
			if len(action) < 2 {
				message = "Неправильная команда"
			} else {
				message = Player.takeItem(action[1])
			}
		case "применить":
			if len(action) < 3 {
				message = "Неправильная команда"
			} else {
				message = Player.use(action[1], action[2])
			}
		default:
			message = Player.unknownAction()
		}
	} else {
		message = "Неправильная команда"
	}
	return message
}

func main() {
}
