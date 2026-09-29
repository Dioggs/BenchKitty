run:
	go run . -r 100 -d 500 https://jsonplaceholder.typicode.com/todos/1
	
file: 
	go run . -r 10 -d 250 -o ./ https://jsonplaceholder.typicode.com/todos/1
