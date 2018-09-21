#API doc

GET 'queryPeople?name=[name]'
Return value string

GET 'createPeople?name=[name]'
Return Tx_id

GET 'createTask?coin=[]&publisher=[]&type=[]&content=[]'
Return Tx_id

GET 'queryTask?tid=[]'
Return JSON '{"accepted":"false","completed":"Yes","id":"Task0","owner":"null","publisher":"hezhiyu","tasktype":"person","timecoin":"50","title":"nothing"}'

GET 'queryAllTasks'
Return JSON Array 
[{"Key":"Task0", "Record":{"accepted":"false","completed":"Yes","id":"Task0","owner":"null","publisher":"hezhiyu","tasktype":"person","timecoin":"50","title":"nothing"}},
{"Key":"Task1", "Record":{"accepted":"false","completed":"Yes","id":"Task0","owner":"null","publisher":"hezhiyu","tasktype":"person","timecoin":"50","title":"nothing"}]

GET 'changeTaskOwner?tid=[]&name=[]' (To accept One Task)
Return Tx_id

GET 'changeTaskState?tid=[]&name=[]' (To complete One Task)
Return Tx_id
