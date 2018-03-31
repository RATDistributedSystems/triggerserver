package main

import (
	"fmt"
	//"bufio"
	//"strings"
	"strconv"
	//"net"
	//"github.com/gocql/gocql"
	"time"
)


func checkBuyTriggers() {

	for{
		//check every 10 milliseconds
		timer1 := time.NewTimer(time.Millisecond * 2000)
		<-timer1.C

		var userid string
		var pendingcash int
		var triggervalue int
		var stock string
		var transactionNum int

		//check if user currently owns any of this stock
		iter := sessionGlobalTR.Query("SELECT userid, pendingcash, triggerValue, stock FROM buyTriggers WHERE pending=TRUE").Iter()
		for iter.Scan(&userid, &pendingcash, &triggervalue, &stock) {
			//delete record
			if err := sessionGlobalTR.Query("DELETE FROM buyTriggers WHERE pending=TRUE AND userid='" + userid + "' AND stock ='" + stock + "'").Exec(); err != nil {
				panic(fmt.Sprintf("Problem DELETING pending buy trigger", err))
			}
			//set record to not pending
			pendingcashstring := strconv.FormatInt(int64(pendingcash), 10)
			triggervaluestring := strconv.FormatInt(int64(triggervalue), 10)
			if err := sessionGlobalTR.Query("INSERT INTO buyTriggers (pending, userid, stock, pendingcash, triggervalue) VALUES (FALSE ,'" + userid + "','" + stock + "'," + pendingcashstring + "," + triggervaluestring + ")").Exec(); err != nil {
				panic(fmt.Sprintf("Problem INSERTING pending buy trigger", err))
			}
			//process the buy trigger
			go processBuyTrigger(userid, stock, triggervalue, pendingcash, transactionNum)

		}
		if err := iter.Close(); err != nil {
			panic(fmt.Sprintf("problem creating session", err))
		}

	}
}

//Execute the buy trigger when the correct condition is met
func processBuyTrigger(userId string, stock string, triggerValue int, pendingCash int, transactionNum int){

	for {
	fmt.Println("4");
	//check every 10 milliseconds
	timer1 := time.NewTimer(time.Millisecond * 500)
	<-timer1.C
	fmt.Println("5");
	var quotePrice = quoteRequest(userId, stock, transactionNum)
	var operation bool = true;

	//Check to ensure the trigger still exists and has not been cancelled
	exists := checkTriggerExists(userId, stock, operation)
	if exists == false {
		return
	}

	fmt.Println("6");
	fmt.Println("quoteprice")
	fmt.Println(quotePrice)
	fmt.Println("triggervalue")
	fmt.Println(triggerValue)
	if(quotePrice <= triggerValue){
		amount, _ := checkStockOwnership(userId, stock)
			fmt.Println("6.5");
			fmt.Println("AMOUNT")
			fmt.Println(amount)
		if(amount != 0){ //-------------------USER ALREADY OWNS SOME OF THIS STOCK ---------------------

			var usableCash int
			var remainingCash int
			//var usid string

			stockamount := amount
			stockValue := quotePrice

			//calculate amount of stocks can be bought
			buyableStocks := pendingCash / stockValue
			buyableStockTotal := buyableStocks + stockamount
			fmt.Println(buyableStocks)
			//remaining money
			fmt.Println(stockValue)
			remainingCash = pendingCash - (buyableStocks * stockValue)
			fmt.Println(remainingCash)

			buyableStocksString := strconv.FormatInt(int64(buyableStockTotal), 10)

			//insert new stock record
			if err := sessionGlobalTS.Query("UPDATE userstocks SET stockamount=" + buyableStocksString + " WHERE userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}

			//check users available cash
			if err := sessionGlobalTS.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}

			//add available cash to leftover cash
			//usableCash = usableCash + remainingCash
			addFunds(userId, remainingCash)

			if err := sessionGlobalTR.Query("DELETE FROM buyTriggers WHERE pending=FALSE AND userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}
			fmt.Println("7");
			return

		}else{ //--------------------------USER DOES NOT OWN ANY OF THIS STOCK----------------------------------

			var usableCash int
			var remainingCash int
			//var usid string

			stockValue := quotePrice

			buyableStocks := pendingCash / stockValue
			remainingCash = pendingCash - (buyableStocks * stockValue)
			fmt.Println("BUYABLESTROCKS")
			fmt.Println(buyableStocks)
			buyableStocksString := strconv.Itoa(buyableStocks)
			fmt.Println("BUYABLESTOCKSSTRING")
			fmt.Println(buyableStocksString)

			//insert new stock record
			if err := sessionGlobalTS.Query("INSERT INTO userstocks (usid, userid, stockamount, stock) VALUES (uuid(), '" + userId + "', " + buyableStocksString + ", '" + stock + "')").Exec(); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}

			//check users available cash
			if err := sessionGlobalTS.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
				panic(fmt.Sprintf("problem creating session", err))				}

			//add available cash to leftover cash

			addFunds(userId, remainingCash)


			if err := sessionGlobalTR.Query("DELETE FROM buyTriggers WHERE pending=FALSE AND userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}
			fmt.Println("7");
			return
		}
	}
	}
}
