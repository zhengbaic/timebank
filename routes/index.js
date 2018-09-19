var express = require('express');
var router = express.Router();
const helper = require('helper')

/* GET home page. */
router.get('/', IndexController.test);

module.exports = router;
