package main

import "fmt"

type Animal struct {
	Name  string
	Heart int
}

type Cat struct {
	Animal
	Claws int
}

type Dolphin struct {
	Animal
	Fin int
}

func main() {
	var myCat Cat
	var myDolphin Dolphin

	myCat.Name = "Miou miou"
	myCat.Claws = 20

	fmt.Println("inheritance PoC")
	fmt.Println(myCat)
	fmt.Printf("%+v\n", myCat)
	fmt.Println(myDolphin)

	printName(myCat.Animal)
}

func printName(a Animal) {
	fmt.Println("The animal name is", a.Name)
}
