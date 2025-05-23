Budget Tracker
Ashton Carter
ascarter@chapman.edu
2428763

REQUIREMENTS:
Must have Go Installed and Live server on VSCODE to run front end

RUN INSTRUCTIONS:
To start backend(must be in folder backend):
go run ./cmd
To start frontend click on index.html and then Open With Live Server (VSCODE extension)

NOTES:
Included is my .env file with my API key and SQL connection
If you upload from CSV, it takes a bit to run becuase of the AI response, look at backend terminal for completion notification and then refresh frontend


LINK TO SCREEN CAPTURE
https://www.loom.com/share/d208dcef4fc945f3829b535ef0c61bab?sid=34349877-9600-4a39-8eb6-fe5241834ea8

RUBRIC LOCATIONS:

1. Print/display records from your database/tables.
Shown in the video on dashboard and in the budgets/categories page.

2. Query for data/results with various parameters/filters
Shown in the video on the dashboard and in the budgets/categories page.  Filter by month
backend request is in router_commands. go at line 66

3. Create a new record
Can upload .csv (Test file is in /backend/csv_parser/test.CSV)
backend request is in router_commands. go at line 672

4. Delete records (soft delete function would be ideal)
It is not a soft delete, but all goals/budgets/transactions are able to be deleted.  Manage trasactions has the delete feature for transactions and budgets/goals has delete buttons on the display page
backend request is in router_commands. go at line 503, 533, 641

5. Update records
Can edit goals/transactions
backend request is in router_commands. go at line 599

6. Make use of transactions (commit & rollback)
Bit confused on this part, all the backend commands use commit

7. Generate reports that can be exported (excel or csv format)
This is done in the manage transactions tab, backend request is in router_commands. go at line 702

8. One query must perform an aggregation/group-by clause
backend request is in router_commands. go at line 450

9. One query must contain a subquery.
Done throughout, backend request is in router_commands. go at line 89

10. Two queries must involve joins across at least 3 tables
backend request is in router_commands. go at line 450 and 175

11. Enforce referential integrality (PK/FK Constraints)
In SQL dump

12. Include Database Views, Indexes
View is called transactions_with_cats and is in backend request is in router_commands. go at line 712

13. Use at least 5 entities
In SQL dump(budgets, goals, transactions, categories, users)

