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
		timer1 := time.NewTimer(time.Millisecond * 10)
		<-timer1.C

		var userid string
		var pendingcash int
		var triggervalue int
		var stock string
		var transactionNum int

		//check if user currently owns any of this stock
		iter := sessionGlobalTR.Query("SELECT userid, pendingcash, triggerValue, stock FROM buyTriggers WHERE pending=TRUE").Iter()
		for iter.Scan(&userid, &pendingcash, &triggervalue, &stock) {

			//set record to "not pending"
			if err := sessionGlobalTR.Query("UPDATE buyTriggers SET pending=FALSE WHERE userid='" + userid + "' AND stock ='" + stock + "'").Exec(); err != nil {
				panic(fmt.Sprintf("Problem UPDATING pending buy trigger", err))
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


	//check every 10 milliseconds
	timer1 := time.NewTimer(time.Millisecond * 10)
	<-timer1.C

	var quotePrice = quoteRequest(userId, stock, transactionNum)
	var operation bool = true;

	//Check to ensure the trigger still exists and has not been cancelled
	exists := checkTriggerExists(userId, stock, operation)
	if exists == false {
		return
	}


	if(quotePrice >= triggerValue){
		amount, _ := checkStockOwnership(userId, stock)
		
		if(amount != 0){ //-------------------USER ALREADY OWNS SOME OF THIS STOCK ---------------------

			var usableCash int
			var remainingCash int
			//var usid string

			stockamount := amount
			stockValue := quotePrice

			//calculate amount of stocks can be bought
			buyableStocks := pendingCash / stockValue
			buyableStocks = buyableStocks + stockamount
			//remaining money
			remainingCash = pendingCash - (buyableStocks * stockValue)

			buyableStocksString := strconv.FormatInt(int64(buyableStocks), 10)

			//insert new stock record
			if err := sessionGlobalTS.Query("UPDATE userstocks SET stockamount=" + buyableStocksString + " WHERE userid=" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}

			//check users available cash
			if err := sessionGlobalTS.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}

			//add available cash to leftover cash
			usableCash = usableCash + remainingCash
			addFunds(userId, usableCash)

			if err := sessionGlobalTR.Query("DELETE FROM buyTriggers WHERE pending=FALSE AND userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}

			return

		}else{ //--------------------------USER DOES NOT OWN ANY OF THIS STOCK----------------------------------

			var usableCash int
			var remainingCash int
			//var usid string

			stockValue := quotePrice

			buyableStocks := pendingCash / stockValue
			remainingCash = pendingCash - (buyableStocks * stockValue)
			buyableStocksString := strconv.FormatInt(int64(buyableStocks), 10)


			//insert new stock record
			if err := sessionGlobalTS.Query("INSERT INTO userstocks (usid, userid, stockamount, stock) VALUES (uuid(), '" + userId + "', " + buyableStocksString + ", '" + stock + "')").Exec(); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}

			//check users available cash
			if err := sessionGlobalTS.Query("select usableCash from users where userid='" + userId + "'").Scan(&usableCash); err != nil {
				panic(fmt.Sprintf("problem creating session", err))				}

			//add available cash to leftover cash
			usableCash = usableCash + remainingCash
			usableCashString := strconv.FormatInt(int64(usableCash), 10)

			//re input the new cash value in to the user db
			if err := sessionGlobalTS.Query("UPDATE users SET usableCash =" + usableCashString + " WHERE userid='" + userId + "'").Exec(); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}

			if err := sessionGlobalTR.Query("DELETE FROM buyTriggers WHERE pending=FALSE AND userid='" + userId + "' AND stock='" + stock + "'").Exec(); err != nil {
				panic(fmt.Sprintf("problem creating session", err))
			}

			return
		}
	}
}
