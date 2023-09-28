package main

import "fmt"
import "time"
import "runtime"
import "sync"
import "os"
import "bufio"
import "strconv"
import "math"

/* A brute force recursive implementation of 0-1 Knapsack problem 
modified from: https://www.geeksforgeeks.org/0-1-knapsack-problem-dp-10 */
var wg sync.WaitGroup
var k[]int
var n[]string

func Max(x int,nx string ,y int,ny string) (int, string){
	if x < y {
		return y, ny
	}
	return x, nx
}

//function made for reading file
func ReadFile(name string) ([]string,[]int,[]int,int) {
	data, err := os.Open(name)
	if err != nil {
		fmt.Println("File reading error", err)
		return nil, nil, nil, 1
	}

	fileScanner := bufio.NewScanner(data)
	fileScanner .Split(bufio.ScanWords)
	fileScanner.Scan()
	num, _:= strconv.Atoi(fileScanner.Text())

	weights := make([]int, num)
	values  := make([]int, num)
	names := make([]string, num)

	for i := 0; i < num; i++ {
		fileScanner.Scan()
		names[i] = fileScanner.Text()
		fileScanner.Scan()
		temp, _:= strconv.Atoi(fileScanner.Text())
		values[i] = temp;
		fileScanner.Scan()
		temp, _= strconv.Atoi(fileScanner.Text())
		weights[i] = temp;
	}
	
	fileScanner.Scan()
	W, _:= strconv.Atoi(fileScanner.Text())
	return names, values, weights, W
}

// Returns the maximum value that 
// can be put in a knapsack of capacity W 
func KnapSack(W int, wt []int, val []int, n []string)(int, string) { 

	// Base Case 
	if (len(wt) == 0 || W == 0) {
		return 0 , ""
		
	}
	last := len(wt)-1

	// If weight of the nth item is more 
	// than Knapsack capacity W, then 
	// this item cannot be included 
	// in the optimal solution 
	if wt[last] > W { 
		return KnapSack(W, wt[:last], val[:last],n[:last])

	// Return the maximum of two cases: 
	// (1) nth item included 
	// (2) item not included 
	} else {
		
		k1,n1 := KnapSack(W - wt[last], wt[:last], val[:last],n[:last])	 
		k2,n2 :=KnapSack(W, wt[:last], val[:last],n[:last])

		return Max(k1 + val[last],n1 + n[last]+", ",k2,n2)
	}
	
} 

//split solve function solves the equation by splitting it splits # of times
func SplitSolve(W int, weights []int, values []int, names []string, splits int) (int, string) { 
	pow := int(math.Pow(float64(2), float64(splits)))
	var tempW int
	var bin string
	last := len(weights)-splits
	
	k := make([]int, pow)
	n := make([]string,pow)
	
	for i := 0; i < pow ; i++ {
		tempW = W
		bin = strconv.FormatInt(int64(i), 2)
		for len(bin) < splits{
			bin = "0"+ bin
		}
		
		for j := 0; j < splits ; j++ {

			if(bin[j] == '1' ){
				tempW -=  weights [last+j];
				k[i] += values [last+j]
				n[i] += names [last+j] + ", "
			}
		}
		if tempW > 0 {
			wg.Add(1)
			go f(tempW, weights [:last], values [:last],names [:last],&k[i],&n[i])
		}else{
			k[i] = 0
			n[i] = ""
		}
	}
	wg.Wait()
	maxIndex := 0
	for i := 0; i < len(k) ; i++ {
		
    		if (k[i]> k[maxIndex]){
			maxIndex = i;
		}
	}
	
	return k[maxIndex], n[maxIndex]

}
func f(tempW int, weights []int, values []int, names []string,k *int, n *string) {
	var x int
	var y string
	x,y = KnapSack(tempW, weights, values, names)	

	*k += x
	*n += y
	wg.Done()
}

// Driver code 
func main()  { 
	

	names, values, weights, W := ReadFile("14.txt")
	runtime.GOMAXPROCS(2)
	fmt.Println("Number of cores: ",runtime.NumCPU())

	start := time.Now();
	
	fmt.Println(SplitSolve(W, weights, values, names, 1))		//splits into 2 solutions
	end := time.Now();
	fmt.Printf("Total runtime: %s\n", end.Sub(start))	

} 
