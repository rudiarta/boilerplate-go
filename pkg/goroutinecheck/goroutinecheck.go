package goroutinecheck

import (
	"fmt"
	"sync"
)

/* Adding counting go routine function
put this on top of the line on go routine
go func(){
	keyGoRoutine := uuid.New().String() //Generate unique identifier
	goroutinecheck.IncreaseTotalGoRoutineCount(keyGoRoutine)
	defer goroutinecheck.DecreaseTotalGoRoutineCount(keyGoRoutine)

	// Do something

}()
*/

var internalCache sync.Map = sync.Map{}

func IncreaseTotalGoRoutineCount(key string) {
	internalCache.LoadOrStore(key, 1)
}

func DecreaseTotalGoRoutineCount(key string) error {
	_, ok := internalCache.LoadAndDelete(key)
	if !ok {
		return fmt.Errorf("key not exist")
	}

	return nil
}

func GetTotalGoRoutineCount() int {
	totalKey := 0

	internalCache.Range(func(k, v interface{}) bool {
		totalKey++
		return true
	})

	return totalKey
}
