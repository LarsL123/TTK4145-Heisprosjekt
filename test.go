package main

import "fmt"
type person struct{
	name string
	age int
	country string
}

func printName(guy *person){
	fmt.Println("Name: " + guy.name)
	fmt.Printf("Age: %d \n", guy.age)
	fmt.Println("Country: " + guy.country)
}

func main(){
	peter := person("peter", 24, "US")
	
}