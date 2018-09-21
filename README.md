# API doc

- To search coin of a person
GET 'queryPeople?name=[name]' </br>
Return value string

- To create a person
GET 'createPeople?name=[name]' </br>
Return Tx_id

- To create a task
GET 'createTask?coin=[]&publisher=[]&type=[]&content=[]' </br>
Return Tx_id

- To search a Task
GET 'queryTask?tid=[]' </br>
Return JSON </br> '{"accepted":"false","completed":"Yes","id":"Task0","owner":"null","publisher":"hezhiyu","tasktype":"person","timecoin":"50","title":"nothing"}'

- To search all Tasks
GET 'queryAllTasks'</br>
Return JSON Array </br>
[{"Key":"Task0", "Record":{"accepted":"false","completed":"Yes","id":"Task0","owner":"null","publisher":"hezhiyu","tasktype":"person","timecoin":"50","title":"nothing"}},
{"Key":"Task1", "Record":{"accepted":"false","completed":"Yes","id":"Task0","owner":"null","publisher":"hezhiyu","tasktype":"person","timecoin":"50","title":"nothing"}]

- To accept One Task
GET 'changeTaskOwner?tid=[]&name=[]'  </br>
Return Tx_id

- To complete One Task
GET 'changeTaskState?tid=[]&name=[]' </br>
Return Tx_id
