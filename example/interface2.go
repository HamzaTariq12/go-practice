package example

import "fmt"

type Animal interface {
	Speak()
	Legs()
}

type Dog struct{}

func (d Dog) Speak() {
	fmt.Println("I say, Woof Woof")
}

func (d Dog) Legs() {
	fmt.Println("I have 4 legs")
}

type Cat struct{}

func (c Cat) Speak() {
	fmt.Println("I say, Meow Meow")
}

func (c Cat) Legs() {
	fmt.Println("I have 4 legs")
}

func guessAnimal(a Animal) {
	a.Speak() // Capitalized
	a.Legs()  // Capitalized
}

func mainExample() {
	cat := Cat{}
	dog := Dog{}

	animals := []Animal{cat, dog}

	for _, animal := range animals {
		animal.Speak()
		animal.Legs()
	}

	// First way - direct calls
	cat.Speak() // Capitalized
	cat.Legs()  // Capitalized
	dog.Speak() // Capitalized
	dog.Legs()  // Capitalized

	// Second way - pass to function
	guessAnimal(cat)
	guessAnimal(dog)

	// Third way - using the interface variable
	var animalService Animal
	animalService = cat // Assign a concrete type to the interface
	animalService.Speak()
	animalService.Legs()

	animalService = dog
	animalService.Speak()
	animalService.Legs()
}
