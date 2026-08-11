package example

import (
	"fmt"
	"math"
)

// Struct must use all defined functions in order to implement it
type Shape interface {
	Area() float64
}

type Rectangle struct {
	width, height float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

func calculateArea(s Shape) float64 {
	return s.Area()
}

func main() {
	rect := Rectangle{width: 5, height: 5}
	circle := Circle{radius: 2}

	fmt.Println("Rectangle area:", calculateArea(rect))
	fmt.Println("Rectangle area:", calculateArea(circle))
}
