package models

import (
	//"github.com/astaxie/beego"
	"container/list"
	"math/rand"
	"time"
)

func humanCreatedAt() string {
	return time.Now().Format(TIME_FORMAT)
}

//洗牌算法
func random(nums []int) []int {
	rand.Seed(time.Now().UnixNano())
	n := len(nums)
	var index int
	for i := 0; i < n; i++ {
		index = rand.Intn(n - 1)
		if index != i {
			nums[i], nums[index] = nums[index], nums[i]

		}

	}
	return nums
}

func toArray(l *list.List) []string {
	var array []string
	for e := l.Front(); e != nil; e = e.Next() {
		array = append(array, e.Value.(string))
	}
	return array
}

//判断list是否包含对应字符串
func Contains(l *list.List, value string) (bool, *list.Element) {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == value {
			return true, e
		}
	}
	return false, nil
}

//返回当前回合做任务的人数
//* 5人：2-3-2-3-3（均为出现一个任务失败就判定为任务失败）
//* 6人：2-3-4-3-4（均为出现一个任务失败就判定为任务失败）
//* 7人：2-3-3-4-4（第一个4人任务需要出现两个任务失败才判定为失败，其余只需要一个）
//* 8-10人：3-4-4-5-5（第一个5人任务需要出现两个任务失败才判定为失败，其余只需要一个）
func getMissionerNumber(usernumber , round int) int{
	ret:=0
    switch usernumber {
        case 5:
            switch round{
                case 1:
                    ret = 2
                    
                case 2:
                ret = 3
                   
                case 3:
                ret = 2
                    
                case 4:
                ret = 3
                    
                case 5:
                ret = 3
                    
                default:
                    break
            }
        case 6:
            switch round{
                case 1:
                    ret = 2
                    
                case 2:
                ret = 3
                    
                case 3:
                ret = 4
                    
                case 4:
                ret = 3
                    
                case 5:
                ret = 4
                    
                default:
                    break
            }
        case 7:
            switch round{
                case 1:
                    ret = 2
                    
                case 2:
                ret = 3
                    
                case 3:
                ret = 3
                    
                case 4:
                ret = 4
                    
                case 5:
                ret = 4
                    
                default:
                    break
            }
        default:    //8-10
            switch round{
                case 1:
                    ret = 3
                    
                case 2:
                ret = 4
                    
                case 3:
                ret = 4
                    
                case 4:
                ret = 5
                    
                case 5:
                ret = 5
                    
                default:
                    break
            }
    }
    return ret
}