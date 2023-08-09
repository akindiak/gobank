package main

func main() {
	s := NewApiServer(":3000")

	s.Run()
}
