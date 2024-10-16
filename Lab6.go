package main

import (
	t1 "PP_LAB6/Pack1"
	t2 "PP_LAB6/Pack2"
	"bufio"
	"math/rand"
	"os"
	"sync"

	"fmt"
	"time"
)

func main() {
	//Запуск и создание горутин
	fmt.Println("Параллельное выполнение")
	a := t1.CreateSeries(100)
	go t1.Fact(2)
	go t1.Rand()
	go t1.Sum(a)
	time.Sleep(1 * time.Second)

	//2 10 чисел Фибоначчи
	fmt.Println("Первые 10 чисел Фибоначчи")
	mes := make(chan int, 10)
	go func() { t2.Fib(10, mes) }()
	go func() {
		for i := range mes {
			fmt.Println(i)
		}
	}()
	time.Sleep(1 * time.Second)

	//Две горутины случайные числа и их четность
	fmt.Println("Две горутины ")
	msgRand := make(chan int)
	msgEven := make(chan string)

	go func() {
		for i := 0; i < 10; i++ {
			num := rand.Intn(100)
			msgRand <- num
			time.Sleep(100 * time.Millisecond)
		}
		close(msgRand)
	}()

	go func() {
		for num := range msgRand {
			if num%2 == 0 {
				msgEven <- fmt.Sprintf("Число %d чётное", num)
			} else {
				msgEven <- fmt.Sprintf("Число %d нечётное", num)
			}
		}
		close(msgEven)
	}()
	go func() {
		for i := 0; i < 10; i++ {
			select {
			case msg1, ok := <-msgRand:
				if ok {
					fmt.Println("Случайное число:", msg1)
				}
			case msg2, ok := <-msgEven:
				if ok {
					fmt.Println(msg2)
				} else {

				}
			}
		}
	}()

	//Исользование мьютоксов для увеличения общей переменной
	fmt.Println("Мьютексы")
	var total int
	var mutex = &sync.Mutex{}
	for i := 0; i < 10; i++ {
		go func() {
			mutex.Lock()
			total += 1
			fmt.Printf("total=%v, Выполняет горутиная №%v\n", total, i)
			mutex.Unlock()
		}()
	}
	time.Sleep(1 * time.Second)
	mutex.Lock()
	fmt.Println("Общая переменная ", total)
	mutex.Unlock()
	time.Sleep(1 * time.Second)

	//многопоточный калькулятор
	var wg sync.WaitGroup
	requests := make(chan CalcRequest)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go calculator(requests, &wg)
	}
	sendReq := func(op1, op2 float32, operation string) {
		resultChan := make(chan float32)
		req := CalcRequest{Op1: op1, Op2: op2, Operation: operation, Res: resultChan}
		requests <- req
		result := <-resultChan
		fmt.Printf("Результат %v %v %v = %v\n", op1, operation, op2, result)
	}

	go sendReq(5, 4, "+")
	go sendReq(10, 8, "-")
	go sendReq(5, 0, "/")

	go func() {
		wg.Wait()
		close(requests)
	}()
	time.Sleep(time.Second)

	//пул воркеров
	if len(os.Args) < 2 {
		fmt.Println("Введите при запуске go run main.go <input-file>")
		return
	}
	inputFileName := os.Args[1]
	file, err := os.Open(inputFileName)
	if err != nil {
		fmt.Println("Ошибка открытия: ", err)
		return
	}
	defer file.Close()

	var tasks []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		tasks = append(tasks, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка:", err)
		return
	}
	var numWorkers int
	fmt.Print("Введите количество воркеров: ")
	fmt.Scanf("%d", &numWorkers)
	taskChan := make(chan string, len(tasks))
	resultChan := make(chan string, len(tasks))

	var wg1 sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg1.Add(1)
		go worker(i, taskChan, resultChan, &wg1)
	}

	for _, task := range tasks {
		taskChan <- task
	}
	close(taskChan)
	go func() {
		wg1.Wait()
		close(resultChan)
	}()
	outputFile, err := os.Create("output.txt")
	if err != nil {
		fmt.Println("Ошибка создания:", err)
		return
	}
	defer outputFile.Close()
	for result := range resultChan {
		outputFile.WriteString(result + "\n")
	}

	fmt.Println("Все готово, результаты в output.txt")
}

type CalcRequest struct {
	Op1       float32
	Op2       float32
	Operation string
	Res       chan float32
}

func calculator(requests chan CalcRequest, wg *sync.WaitGroup) {
	for req := range requests {
		var res float32
		switch req.Operation {
		case "+":
			res = req.Op1 + req.Op2
		case "-":
			res = req.Op1 - req.Op2
		case "*":
			res = req.Op1 * req.Op2
		case "/":
			if req.Op2 != 0 {
				res = req.Op1 / req.Op2
			} else {

				req.Res <- 0
			}
		}
		req.Res <- res
	}
	wg.Done()
}

func worker(i int, tasks <-chan string, results chan<- string, wg *sync.WaitGroup) {
	for task := range tasks {
		reversed := reverse(task)
		fmt.Printf("Worker %d обработал %s\n", i, reversed)
		results <- reversed
	}
	wg.Done()
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
