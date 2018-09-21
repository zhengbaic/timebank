var express = require('express');
var router = express.Router();
const helper = require('./helper')

/* GET home page. */
router.get('/queryPeople', helper.queryPeople)
	.get('/createPeople', helper.createPeople)
	.get('/createTask', helper.createTask)
	.get('/queryAllTasks', helper.queryAllTasks)
	.get('/changeTaskOwner', helper.changeTaskOwner)
	.get('/changeTaskState', helper.changeTaskState)
	.get('/queryTask', helper.queryTask)

module.exports = router;
