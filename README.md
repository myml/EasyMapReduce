# EasyMapReduce
一个简单的MapReduce实现
## 使用
简单的例子
```
	//生成数据源
	var data []int
	for i := 0; i < 1000; i++ {
		data = append(data, i)
	}
	Map := func(i int) int {
		return i * i
	}
	var result []int
	Reduce := func(r int) error {
		result = append(result, r)
		return nil
	}
	err := mapreduce(data, goroutineLimit, Map, Reduce)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
```			
复杂一点的例子
```

	source := make(chan int)
	go func() {
		for i := 0; i < 1000; i++ {
			source <- i
		}
		//！！！通道类型的源要自己关闭
		close(source)
	}()
	Map := func(i int) *string {
		if i%2 == 0 {
			return nil
		}
		s := strconv.Itoa(i)
		return &s
	}
	var result []string
	Reduce := func(r *string) error {
		if r != nil {
			result = append(result, *r)
		}
		return nil
	}
	err := mapreduce(source, goroutineLimit, Map, Reduce)
	if err != nil {
		panic(err)
	}
	fmt.Println(result)
```
## 说明
只有一个mapreduce函数
```

func mapreduce(source interface{}, goroutineLimit int, Map interface{}, Reduce interface{})
```
以下说明中的T为任意类型
* source:chan T or []T		
	一个slice或chan类型
* goroutineLimit:int
	并发限制，最多同时执行goroutineLimit个Map函数
* Map:func(T)... 		
	只能有一个参数的函数，参数类型要为source中的元素类型,不限制返回值个数和类型，返回值会传给Reduce
* Reduce:func(...)error		
	参数类型和个数要和Map返回值相同		
	Reduce的返回值如不为空，则等待正在运行的Map结束后，退出mapreduce，返回值为Reduce的返回值		
	在多线程执行Map函数后返回的结果传递给Reduce,Reduce时单线程执行，线程安全的