# creating bank:
echo "creating bank"
node invoke bankcc writeBankInfo 1bank kvb chennai 40A 1000 1000 1000 1000 1000

# creating business:
echo "creating business"
node invoke businesscc putNewBusinessInfo 1bus tata 12348901 4000000 1000 1000 1000 12 8 1000 1000
node invoke businesscc putNewBusinessInfo 2bus mrf 12348902 4000000 1000 1000 1000 12 8 1000 1000

# creating Program
echo "creating Program"
node invoke programcc writeProgram 1prg program1 1bus Accounts_Payable 10/04/2019 10000 6 buyer 4 100 pragadeesh 123452

# creating PPR
echo "creating PPR"
node invoke pprcc createPPR 1ppr 1prg 2bus seller 12000 3 100 5 40 34tf2

# creating instrument
echo "creating instrument"
node invoke instrumentcc enterInstrument 1ins 23/10/2018 2bus 1bus 1000 23/07/2019 1prg 1ppr 34 04/01/2018:12:43:59

# creating loan
echo "creating loan"
node invoke loancc newLoanInfo 1loan 1ins 1bus 1prg 900 pragadeesh 5 23/10/2018 25/09/2018:20:45:01 0 0 0 1bus 2bus

# invoking disbursement transaction
echo "disbursement transaction"
node invoke txncc newTxnInfo 1txn disbursement 23/04/2018 1loan 1inst 300 1bank 1bus pragadeesh

# invoking accrual transaction
# node invoke txncc newTxnInfo 2txn accrual 23/04/2018 1loan 1inst 800 1bus 1bank pragadeesh

# invoking charges transaction
# node invoke txncc newTxnInfo 2txn charges 23/04/2018 1loan 1inst 800 1bank 1bus pragadeesh

# invoking interest_accrued_charge transaction

# invoking interest_in_advance transaction
# node invoke txncc newTxnInfo 2txn interest_in_advance 23/04/2018 1loan 1inst 800 1bank 1bus pragadeesh

# invoking interest_refund transaction
# node invoke txncc newTxnInfo 2txn interest_refund 23/04/2018 1loan 1inst 800 1bank 1bus pragadeesh

# invoking margin_refund transaction
# node invoke txncc newTxnInfo 2txn margin_refund 23/04/2018 1loan 1inst 800 1bank 1bus pragadeesh

# invoking penal_charges transaction
# node invoke txncc newTxnInfo 2txn penal_charges 23/04/2018 1loan 1inst 800 1bank 1bus pragadeesh

# invoking penal_interest_collection transaction
# node invoke txncc newTxnInfo 2txn penal_interest_collection 23/04/2018 1loan 1inst 800 1bank 1bus pragadeesh

# invoking repayment transaction
echo "repayment transaction"
node invoke txncc newTxnInfo 2txn repayment 23/04/2018 1loan 1inst 800 1bus 1bank pragadeesh 

# invoking tds transaction
#node invoke txncc newTxnInfo 2txn TDS 23/04/2018 1loan 1inst 800 1bank 1bus pragadeesh 
