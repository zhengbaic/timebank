# API doc

- To search coin of a person
GET 'queryPeople?name=[name]'
Return value string

- To create a person
GET 'createPeople?name=[name]'
Return Tx_id

- To create a task
GET 'createTask?coin=[]&publisher=[]&type=[]&content=[]'
Return Tx_id

- To search a Task
GET 'queryTask?tid=[]'
Return JSON '{"accepted":"false","completed":"Yes","id":"Task0","owner":"null","publisher":"hezhiyu","tasktype":"person","timecoin":"50","title":"nothing"}'

- To search all Tasks
GET 'queryAllTasks'
Return JSON Array
[{"Key":"Task0", "Record":{"accepted":"false","completed":"Yes","id":"Task0","owner":"null","publisher":"hezhiyu","tasktype":"person","timecoin":"50","title":"nothing"}},
{"Key":"Task1", "Record":{"accepted":"false","completed":"Yes","id":"Task0","owner":"null","publisher":"hezhiyu","tasktype":"person","timecoin":"50","title":"nothing"}]

- To accept One Task
GET 'changeTaskOwner?tid=[]&name=[]' 
Return Tx_id

- To complete One Task
GET 'changeTaskState?tid=[]&name=[]'
Return Tx_id
