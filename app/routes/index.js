var express = require('express');
var router = express.Router();
const helper = require('./helper')


// router.all('*', function(req, res, next) {
//     res.header("Access-Control-Allow-Origin", "*");
//     res.header("Access-Control-Allow-Headers", "X-Requested-With");
//     res.header("Access-Control-Allow-Methods","PUT,POST,GET,DELETE,OPTIONS");
//     res.header("X-Powered-By",' 3.2.1')
//     res.header("Content-Type", "application/json;charset=utf-8");
//     next();
// });

/* GET home page. */
router.get('/queryPeople', helper.queryPeople)
	.get('/createPeople', helper.createPeople)
	.get('/createTask', helper.createTask)
	.get('/queryAllTasks', helper.queryAllTasks)
	.get('/changeTaskOwner', helper.changeTaskOwner)
	.get('/changeTaskState', helper.changeTaskState)
	.get('/queryTask', helper.queryTask)
	.get('/queryPeopleAccept', helper.queryPeopleAccept)
	.get('/queryPeoplePublish', helper.queryPeoplePublish)
	.get('/queryPeopleComplete', helper.queryPeopleComplete)
	.get('/log', helper.log)


// API doc

// GET 'queryPeople?name=[name]'
// Return value string

// GET 'createPeople?name=[name]'
// Return Tx_id

// GET 'createTask?coin=[]&publisher=[]&type=[]&content=[]'
// Return Tx_id

// GET 'queryTask?tid=[]'
// Return JSON '{"accepted":"false","completed":"Yes","id":"Task0","owner":"null","publisher":"hezhiyu","tasktype":"person","timecoin":"50","title":"nothing"}'

// GET 'queryAllTasks'
// Return JSON Array 
// [{"Key":"Task0", "Record":{"accepted":"false","completed":"Yes","id":"Task0","owner":"null","publisher":"hezhiyu","tasktype":"person","timecoin":"50","title":"nothing"}},
// {"Key":"Task1", "Record":{"accepted":"false","completed":"Yes","id":"Task0","owner":"null","publisher":"hezhiyu","tasktype":"person","timecoin":"50","title":"nothing"}]

// GET 'changeTaskOwner?tid=[]&name=[]' (To accept One Task)
// Return Tx_id

// GET 'changeTaskState?tid=[]&name=[]' (To complete One Task)
// Return Tx_id

module.exports = router;
