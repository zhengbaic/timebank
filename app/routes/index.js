var express = require('express');
var router = express.Router();
const helper = require('./helper')
const multer = require('multer')
const fs = require('fs')

var uploadimage = multer({dest: 'tmp/'})

// router.all('*', function(req, res, next) {
//     res.header("Access-Control-Allow-Origin", "*");
//     res.header("Access-Control-Allow-Headers", "X-Requested-With");
//     res.header("Access-Control-Allow-Methods","PUT,POST,GET,DELETE,OPTIONS");
//     res.header("X-Powered-By",' 3.2.1')
//     res.header("Content-Type", "application/json;charset=utf-8");
//     next();
// });

/* GET home page. */
router.get('/query', helper.query)
	.get('/createPeople', helper.createPeople)
	.get('/createInstitution', helper.createInstitution)
	.get('/registerInstitution', helper.registerInstitution)
	.get('/giveInstitutionCoin', helper.giveInstitutionCoin)

	.get('/createTask', helper.createTask)
	.get('/createGroupTask', helper.createGroupTask)
	.get('/queryAllTasks', helper.queryAllTasks)
	.get('/queryAllIns', helper.queryAllIns)

	.get('/acceptSingleTask', helper.acceptSingleTask)
	.get('/completeSingleTask', helper.completeSingleTask)
	.get('/acceptGroupTask', helper.acceptGroupTask)
	.get('/completeGroupTask', helper.completeGroupTask)
	.get('/cancelTask', helper.cancelTask)

	.get('/queryConfirm', helper.queryConfirm)
	.get('/confirm', helper.confirm)
	.get('/recordDisputedTask', helper.recordDisputedTask)
	.get('/log', helper.log)
	.get('/readMes', helper.readMes)
	.get('/addMes', helper.addMes)
	.get('/isMesread', helper.isMesread)

	.get('/queryPeopleAccept', helper.queryPeopleAccept)
	.get('/queryPeoplePublish', helper.queryPeoplePublish)
	.get('/queryPeopleComplete', helper.queryPeopleComplete)
	.get('/queryPeopleDisputed', helper.querypeopleDistued)

	.get('/queryBlockInfo', helper.queryBlockInfo)
	.get('/queryTransaction', helper.queryTransaction)

	.get('/registerInstitutionbackdoor', helper.registerInstitutionbackdoor)

	.post('/uploadImage', uploadimage.single('imageFile'), function(req, res) {
		var tmp_path = req.file.path
		var target_path = 'upload/' + req.file.originalname
		console.log(req.body.name)
		console.log(req.file.path)
		console.log(req.file.originalname)
		var src = fs.createReadStream(tmp_path);
  		var dest = fs.createWriteStream(target_path);
  		src.pipe(dest);
  		src.on('end', function() { res.send('complete'); });
  		src.on('error', function(err) { res.send('error'); });
	})

	.get('/showImage', function(req, res) {
		fs.readFile('./upload/' + req.query.name + ".jpg", function(err, file) {
			if(err) {
				res.send(err)
			} else {
				res.writeHead(200, {"Content-Type": "image/png"})
				res.write(file, "binary")
				res.end()
			}
		})
	})

	.get('/index', function(req, res) {
		res.render("index")
	})
	


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
