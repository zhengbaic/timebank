const Fabric_Client = require('fabric-client');
const path = require('path');
const os = require('os')
const util = require('util')


var mongoose = require('mongoose')
let db = mongoose.connect('mongodb://localhost:27017/db38')

let Myuser = mongoose.model("users", {
	uid: String,
	username : String,
	password : String,
	phone: String,
	address: String
})

let Myins = mongoose.model("ins", {
	iid: String,
	name : String,
	password : String,
	duty: String,
	address: String,
	phone: String,
	authority: String
})

//no / yes
let Mytask = mongoose.model("tasks", {
	tid: String,
	publisher: String,
	coin: String,
	type: String,
	content: String,
	owner: String,
	time_create: String,
	time_accept: String,
	time_complete: String,
	cancel: String,
	pconfirm: String,
	aconfirm: String
})

let Mygtask = mongoose.model("gtasks", {
	tid: String,
	publisher: String,
	coin: String,
	type: String,
	content: String,
	owner: String,
	time_create: String,
	time_accept: String,
	time_complete: String,
	cancel: String
})

let Mymessage = mongoose.model("messages", {
	name: String,
	type: String,
	post: String,
	post_type: String,
	time: String,
	mess: String,
	read: String
})

let MyIns_register = mongoose.model("regs", {
	name: String,
	regNum: Number,
	regs: Array
})

let MyTx = mongoose.model("txs", {
	time: String,
	txid: String,
	behavior: String,
	peer: String,
	publisher: String,
	object: String
})

var task_num = 0;
var gtask_num = 0;
var user_num = 0;
var ins_num = 0;
var channel_name = 'mychannel'
//load user1
var fabric_client = new Fabric_Client();
var member_user = null;
var store_path = path.join('/home/zbc/Desktop/timebank/app', 'hfc-key-store');
console.log('Store path:'+store_path);
var tx_id = null;

var channel = fabric_client.newChannel(channel_name);
var peer = fabric_client.newPeer('grpc://localhost:7051');
channel.addPeer(peer);
var order = fabric_client.newOrderer('grpc://localhost:7050')
channel.addOrderer(order);
Fabric_Client.newDefaultKeyValueStore({ path: store_path}).then((state_store) => {
		// assign the store to the fabric client
		fabric_client.setStateStore(state_store);
		var crypto_suite = Fabric_Client.newCryptoSuite();
		// use the same location for the state store (where the users' certificate are kept)
		// and the crypto store (where the users' keys are kept)
		var crypto_store = Fabric_Client.newCryptoKeyStore({path: store_path});
		crypto_suite.setCryptoKeyStore(crypto_store);
		fabric_client.setCryptoSuite(crypto_suite);

		// get the enrolled user from persistence, this user will sign all requests
		return fabric_client.getUserContext('user1', true);
	}).then((user_from_store) => {
		if (user_from_store && user_from_store.isEnrolled()) {
		    console.log('Successfully loaded user1 from persistence');
		    member_user = user_from_store;
		} else {
		    throw new Error('Failed to get user1.... run registerUser.js');
		}
	})


Myuser.find({}, function(err, docs) {
	console.log(docs)
	user_num = docs.length;
	console.log("user_num:" + user_num)
})
Myins.find({}, function(err, docs) {
	console.log(docs)
	ins_num = docs.length;
	console.log("ins_num:" + ins_num)
})
Mytask.find({}, function(err, docs) {
	console.log(docs)
	task_num = docs.length;
	console.log("task_num:" + task_num)
})
Mygtask.find({}, function(err, docs) {
	gtask_num = docs.length;
	console.log("gtask_num:" + gtask_num)
})

async function query_chaincode(request, callback) {
	result = await channel.queryByChaincode(request).then((query_responses) => {
		    console.log("Query has completed, checking results");
		    // query_responses could have more than one  results if there multiple peers were used as targets
		    if (query_responses && query_responses.length == 1) {
		        if (query_responses[0] instanceof Error) {
		            console.error("error from query = ", query_responses[0]);
		            result = "Could not locate tuna"
		        } else {
		            console.log("Response is ", query_responses[0].toString());
		            callback(query_responses[0].toString())
		        }
		    } else {
		        console.log("No payloads were returned from query");
                result = "Could not locate tuna"
		    }
		}).catch((err) => {
            console.error('Failed to query successfully :: ' + err);
            result = 'Failed to query successfully :: ' + err
            
		});
}
function timestampToTime(timestamp) {
        var date = new Date(timestamp);//时间戳为10位需*1000，时间戳为13位的话不需乘1000
        Y = date.getFullYear() + '-';
        M = (date.getMonth()+1 < 10 ? '0'+(date.getMonth()+1) : date.getMonth()+1) + '-';
        D = date.getDate() + ' ';
        h = date.getHours() + ':';
        m = date.getMinutes() + ':';
        s = date.getSeconds();
        return Y+M+D+h+m+s;
    }

async function invoke_chaincode(request, callback) {
	result = await channel.sendTransactionProposal(request).then((results) => {
		    var proposalResponses = results[0];
		    var proposal = results[1];
		    let isProposalGood = false;
		    if (proposalResponses && proposalResponses[0].response &&
		        proposalResponses[0].response.status === 200) {
		            isProposalGood = true;
		            console.log('Transaction proposal was good');
		        } else {
		            console.error('Transaction proposal was bad');
		        }
		    if (isProposalGood) {
		        console.log(util.format(
		            'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s"',
		            proposalResponses[0].response.status, proposalResponses[0].response.message));

		        // build up the request for the orderer to have the transaction committed
		        var request = {
		            proposalResponses: proposalResponses,
		            proposal: proposal
		        };

		        // set the transaction listener and set a timeout of 30 sec
		        // if the transaction did not get committed within the timeout period,
		        // report a TIMEOUT status
		        var transaction_id_string = tx_id.getTransactionID(); //Get the transaction ID string to be used by the event processing
		        var promises = [];

		        var sendPromise = channel.sendTransaction(request);
		        promises.push(sendPromise); //we want the send transaction first, so that we know where to check status

		        // get an eventhub once the fabric client has a user assigned. The user
		        // is required bacause the event registration must be signed
		        let event_hub = fabric_client.newEventHub();
		        event_hub.setPeerAddr('grpc://localhost:7053');

		        // using resolve the promise so that result status may be processed
		        // under the then clause rather than having the catch clause process
		        // the status
		        let txPromise = new Promise((resolve, reject) => {
		            let handle = setTimeout(() => {
		                event_hub.disconnect();
		                resolve({event_status : 'TIMEOUT'}); //we could use reject(new Error('Trnasaction did not complete within 30 seconds'));
		            }, 3000);
		            event_hub.connect();
		            event_hub.registerTxEvent(transaction_id_string, (tx, code) => {
		                // this is the callback for transaction event status
		                // first some clean up of event listener
		                clearTimeout(handle);
		                event_hub.unregisterTxEvent(transaction_id_string);
		                event_hub.disconnect();

		                // now let the application know what happened
		                var return_status = {event_status : code, tx_id : transaction_id_string};
		                if (code !== 'VALID') {
		                    console.error('The transaction was invalid, code = ' + code);
		                    resolve(return_status); // we could use reject(new Error('Problem with the tranaction, event status ::'+code));
		                } else {
		                    console.log('The transaction has been committed on peer ' + event_hub._ep._endpoint.addr);
		                    resolve(return_status);
		                }
		            }, (err) => {
		                //this is the callback if something goes wrong with the event registration or processing
		                reject(new Error('There was a problem with the eventhub ::'+err));
		            });
		        });
		        promises.push(txPromise);

		        return Promise.all(promises);
		    } else {
		        console.error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
		        throw new Error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
		    }
		}).then((results) => {
		    console.log('Send transaction promise and event listener promise have completed');
		    // check the results in the order the promises were added to the promise all list
		    if (results && results[0] && results[0].status === 'SUCCESS') {
		        console.log('Successfully sent transaction to the orderer.');
				//res.send(tx_id.getTransactionID());
				callback(tx_id.getTransactionID());
		    } else {
		        console.error('Failed to order the transaction. Error code: ' + response.status);
		    }

		  //   if(results && results[1] && results[1].event_status === 'VALID') {
		  //       console.log('Successfully committed the change to the ledger by the peer');
				// //res.send(tx_id.getTransactionID());
				// callback(tx_id.getTransactionID());
		  //   } else {
		  //       console.log('Transaction failed to be committed to the ledger due to ::'+results[1].event_status);
		  //   }
		}).catch((err) => {
		    console.error('Failed to invoke successfully :: ' + err);
		});
}

async function queryId(name, type) {
	var id
		if(type == "person") {
			var res = await Myuser.findOne({username: name})
			if(res == null) {
				return res
			}
			console.log(res)
			return res.uid
		}

		else {
			var res = await Myins.findOne({name: name})
			if(res == null) {
				return res
			}
			return res.iid
		}
	}

module.exports = {
	async readMes(req, res) {
		var name = req.query.name
		var type = req.query.type

		var result = await Mymessage.update({
			name: name,
			type: type
		}, {
			$set: {
				read: "1"
			}
		}, {
			multi: true
		})

		var docs = await Mymessage.find({
			name: name,
			type: type
		})

		res.send(docs)
	},
	async addMes(req, res) {
		var name = req.query.name
		var type = req.query.type
		var mess = req.query.mess
		var post = req.query.post
		var ptype = req.query.ptype

		var t = new Date().getTime()

		let new_mess = new Mymessage({
			name: name,
			type: type,
			post: post,
			post_type: ptype,
			mess: mess,
			time: timestampToTime(t),
			read: "0"
		})
		var result = await new_mess.save() 
		res.send("SUCCESS")

	},
	async isMesread(req, res) {
		var name = req.query.name
		var type = req.query.type
		var result = await Mymessage.find({
			name: name,
			type: type,
			read: "0"
		})
		console.log(result.length)

		res.send("" + result.length)
		
	},
	log(req, res) {
		var type = req.query.type
		var people_name = req.query.name
		var password = req.query.password
		if(type == "person") {
			Myuser.findOne({
				username: people_name
			}, function(err, doc) {
				if(err) {
					res.send(err)
					return
				}
				if(doc == null) {
					res.send({
						status: 404,
						message: "Not Exsit"
					})
					return
				}
				console.log(doc)
				if(doc.password == password) {
					res.send({
						status: 200,
						message: ""
					})
				} else {
					res.send({
						status: 404,
						message: "Password Invaild"
					})
				}
			})
		} else {
			Myins.findOne({
				name: people_name
			}, function(err, doc) {
				if(err) {
					res.send(err)
					return
				}
				if(doc == null) {
					res.send({
						status: 404,
						message: "Not Exsit"
					})
					return
				}
				console.log(doc)
				if(doc.password == password) {
					res.send({
						status: 200,
						message: ""
					})
				} else {
					res.send({
						status: 404,
						message: "Password Invaild"
					})
				}
			})
		}
	},
	async query(req, res) {
		var type = req.query.type
		var name = req.query.name
		if(type == "person" || type == "ins") {
			var pro = queryId(name, type)
			pro.then((id) => {
				if(id == null) {
					res.send({
						status: 404,
						message: "Not Exsit"
					})
					return
				}
				console.log("id:" + id)
				const request = {
		    		chaincodeId: 'bank',
		    		txId: tx_id,
		    		fcn: 'query',
		    		args: [id]
				};
				query_chaincode(request, async function(result) {
					if(type == "ins") {
						result = JSON.parse(result)
						var temp = await Myins.findOne({name: name})
						result.duty = temp.duty
						result.address = temp.address
						result.phone = temp.phone
						result.name = temp.name
						console.log(result)
						res.send({
							status: 200,
							message: JSON.stringify(result)
						})
						return
					} else {
						result = JSON.parse(result)
						var temp = await Myuser.findOne({username: name})
						result.address = temp.address
						result.phone = temp.phone
						result.name = temp.username
						console.log(result)
						res.send({
							status: 200,
							message: JSON.stringify(result)
						})
						return
					}
				})
			})
		}
		else {
			id = "task" + req.query.name

			const request = {
			    chaincodeId: 'bank',
			    txId: tx_id,
			    fcn: 'query',
			    args: [id]
			};
			query_chaincode(request, function(result) {
				console.log(result)
				res.send({
					status: 200,
					message: result
				})
			})
		}
	},
	async createInstitution(req, res) {
		var password = req.query.password
		var iid = "ins" + ins_num
		var name = req.query.name
		var duty = req.query.duty
		var address = req.query.address
		var phone = req.query.phone

		var resu = await Myins.findOne({
			name: name
		})
		console.log(resu)
		if(resu != null) {
			console.log("repeated")
			res.send({
				status: 404,
				message: "Duplicated"
			})
			return 
		} else {
			tx_id = fabric_client.newTransactionID();
			console.log("Assigning transaction_id: ", tx_id._transaction_id);
			const request = {
			    chaincodeId: 'bank',
			    txId: tx_id,
			    fcn: 'createInstitution',
			    args: [name],
			    chainId: channel_name
			};
			invoke_chaincode(request, async function(txid) {
				let new_ins = new Myins({
					iid: iid,
					name: name,
					password: password,
					duty: duty,
					address: address,
					phone: phone,
					authority: "0"
				})
				var temp = await new_ins.save()
				ins_num = ins_num + 1
				let newregs = new MyIns_register({
					name: name,
					regNum: 0,
					regs: []
				})
				var temp1 = await newregs.save()
				for(var i = ins_num; i--;) {
					console.log(i)
					var ins = await Myins.findOne({iid: "ins" + i})
					console.log(ins)
					if(ins == null || ins.name == name || ins.authority == "0") continue
						else {
							console.log(ins.name)
							let new_mess = new Mymessage({
								name: ins.name,
								type: "ins",
								post: "system",
								post_type: "system",
								mess: "新机构节点:" + name + "请求加入区块链网络,其信息为:负责人:" + duty + ",地址:" + address + ",电话:" + phone,
								time: timestampToTime(new Date().getTime()),
								read: "0"
							})
							var temp2 = await new_mess.save()
						}
				}
				let new_tx = new MyTx({
					time: timestampToTime(new Date().getTime()),
					txid: txid,
					behavior: "创建机构",
					peer: "",
					publisher: iid,
					object: iid
				})
				var temp3 = await new_tx.save()
				res.send({
					status: 200,
					message: ""
				})
			})
		}
	},
	async registerInstitution(req, res) {
		var name = req.query.name
		var iname = req.query.iname
		var result = await MyIns_register.findOne({name: name})
		if(result == null) {res.send("ERROR: No such ins"); return;}
		for (var i = result.regs.length; i--;) {
			if(result.regs[i] == iname) {
				res.send("ERROR: duplicated ins"); return;
			}
		}
		var temp = result.regs
		var tempnum = result.regNum + 1
		temp.push(iname)
		var docs = await Myins.find({"authority": "1"})
		console.log(docs)
		var authority_num = docs.length
		console.log("tempnum:" + tempnum)
		console.log("authority_num:" + authority_num)
		var r2 = await MyIns_register.update({name: name},{regNum: tempnum, regs: temp})
		if(tempnum == authority_num) {
			var pro = queryId(name, "ins")
			pro.then((id) => {
				tx_id = fabric_client.newTransactionID();
				console.log("Assigning transaction_id: ", tx_id._transaction_id);
				const request = {
				    chaincodeId: 'bank',
				    txId: tx_id,
				    fcn: 'registerInstitution',
				    args: [id],
				    chainId: channel_name
				};
				invoke_chaincode(request, async function(txid) {
					var temp = await Myins.update({name: name}, {authority: "1"})
					let new_tx = new MyTx({
						time: timestampToTime(new Date().getTime()),
						txid: txid,
						behavior: "机构审核通过",
						peer: "",
						publisher: id,
						object: id
					})
					var temp3 = await new_tx.save()
					return res.send("register successfully")
				})
			}) 
		}else {
			res.send("SUCCESS")
			return
		}


	},

	async registerInstitutionbackdoor(req, res) {
		var name = req.query.name
		var pro = queryId(name, "ins")
		pro.then((id) => {
			tx_id = fabric_client.newTransactionID();
			console.log("Assigning transaction_id: ", tx_id._transaction_id);
			const request = {
			    chaincodeId: 'bank',
			    txId: tx_id,
			    fcn: 'registerInstitution',
			    args: [id],
			    chainId: channel_name
			};
			invoke_chaincode(request, async function(txid) {
				var temp = await Myins.update({name: name}, {authority: "1"})
				return res.send("register successfully")
			})
		})
	},
	async giveInstitutionCoin(req, res) {
		tx_id = fabric_client.newTransactionID();
		console.log("Assigning transaction_id: ", tx_id._transaction_id);
		const request = {
		    chaincodeId: 'bank',
		    txId: tx_id,
		    fcn: 'giveInstitutionCoin',
		    args: [],
		    chainId: channel_name
		};
		invoke_chaincode(request, function(txid) {
			return res.send("success")
		})
	},
	async createPeople(req, res) {
		var people_name = req.query.name
		var password = req.query.password
		var phone = req.query.phone
		var address = req.query.address
		var resu = await Myuser.findOne({
			username: people_name
		})
		console.log(resu)
		if(resu) {
			console.log("repeated")
			res.send({
				status: 404,
				message: "Duplicated"
			})
			return 
		} else {
			var uid = "person" + user_num
			user_num = user_num + 1
			tx_id = fabric_client.newTransactionID();
			console.log("Assigning transaction_id: ", tx_id._transaction_id);
			const request = {
			    chaincodeId: 'bank',
			    txId: tx_id,
			    fcn: 'createPeople',
			    args: [people_name],
			    chainId: channel_name
			};
			invoke_chaincode(request, async function(txid) {
				let new_user = new Myuser({
					uid: uid,
					username: people_name,
					password: password,
					phone: phone,
					address: address
				})
				let new_tx = new MyTx({
					time: timestampToTime(new Date().getTime()),
					txid: txid,
					behavior: "创建用户",
					peer: "",
					publisher: uid,
					object: uid
				})
				var temp3 = await new_tx.save()

				new_user.save(function(err) {
					if(err) {
						console.log(err)
						return
					}
					return res.send({
							status: 200,
							message: txid
						})
				})
				
			})
		}
	},
	async createTask(req, res) {
		var index = task_num + gtask_num
		var Tid = 'task' + index;
		task_num = task_num + 1;
		var coin = req.query.coin;
		var publisher = req.query.publisher;
		var type = req.query.type;
		if(type == "ins") type = "institution"
		var content = req.query.content;
		var uid = ""
		var pro = queryId(publisher, type)
		pro.then((id) => {
			uid = id
			console.log(uid)
			console.log(Tid)
			console.log(coin)
			console.log(publisher)
			console.log(type)
			console.log(content)
			tx_id = fabric_client.newTransactionID();
			console.log("Assigning transaction_id: ", tx_id._transaction_id);

			const request = {
			    chaincodeId: 'bank',
			    fcn: 'createTask',
			    args: [Tid, coin, publisher, type, content, uid],
			    chainId: channel_name,
			    txId: tx_id
			};
			invoke_chaincode(request, async function(txid) {
				let new_task = new Mytask({
					publisher: publisher,
					tid: Tid,
					coin: coin,
					type: type,
					content: content,
					owner: "none",
					time_create: new Date().getTime(),
					time_accept: "none",
					time_complete: "none",
					cancel: "0",
					pconfirm: "0",
					aconfirm: "0"
				})
				let new_tx = new MyTx({
					time: timestampToTime(new Date().getTime()),
					txid: txid,
					behavior: "创建任务",
					peer: "",
					publisher: publisher,
					object: Tid
				})
				var temp3 = await new_tx.save()
				console.log(JSON.stringify(new_task))
				new_task.save(function(err){
					if(err) {
						console.log(err)
						return
					}
					res.status(200);
					var temp = {
						"tid" : Tid,
						"txid" : txid
					}
					return res.send(temp)
				})
			})
		})
	},
	async createGroupTask(req, res) {
		var index = task_num + gtask_num
		var Tid = 'task' + index;
		gtask_num = gtask_num + 1;
		var coin = req.query.coin
		var publisher = req.query.publisher;
		var type = req.query.type;
		var content = req.query.content;
		var uid = ""
		uid = queryId(publisher, type)
		var num = req.query.num

		tx_id = fabric_client.newTransactionID();
		console.log("Assigning transaction_id: ", tx_id._transaction_id);

		const request = {
		    chaincodeId: 'bank',
		    fcn: 'createGroupTask',
		    args: [Tid, coin, publisher, type, content, uid, num],
		    chainId: channel_name,
		    txId: tx_id
		};
		invoke_chaincode(request, async function(txid) {
			let new_task = new Mygtask({
				publisher: publisher,
				tid: Tid,
				coin: coin,
				type: type,
				content: content,
				owner: "none",
				time_create: new Date().getTime(),
				time_accept: "none",
				time_complete: "none",
				cancel: "0"
			})
			console.log(JSON.stringify(new_task))
			new_task.save(function(err){
				if(err) {
					console.log(err)
					return
				}
				res.status(200);
				return res.send(txid)
			})
		})
	},
	async queryAllTasks(req, res) {
		let result = ''
		var tx_id = null;
		const request = {
		    chaincodeId: 'bank',
		    txId: tx_id,
		    fcn: 'queryAllTasks',
		    args: [""]
		};
		query_chaincode(request, function(result) {
			console.log(result)
			res.send(result)
		})
	},
	async queryAllIns(req, res) {
		var resluts = await Myins.find({})
		for (var i = resluts.length; i--;) {
			resluts[i].password = undefined
		}
		res.send(resluts)
	},
	async cancelTask(req, res) {
		var tid = "task" + req.query.tid
		var ttype = req.query.ttype
		var ptype = req.query.ptype
		tx_id = fabric_client.newTransactionID();
		console.log("Assigning transaction_id: ", tx_id._transaction_id);
		const request = {
			chaincodeId: 'bank',
			txId: tx_id,
			fcn: 'cancelTask',
			args: [tid, ttype, ptype],
			chainId: channel_name
		};
		invoke_chaincode(request, async function(txid) {
			if(ttype == "single") {
				var doc = await Mytask.findOne({tid: tid})

				let new_tx = new MyTx({
					time: timestampToTime(new Date().getTime()),
					txid: txid,
					behavior: "取消任务",
					peer: "",
					publisher: doc.publisher,
					object: tid
				})
				var temp3 = await new_tx.save()
				Mytask.update({
					tid: tid
				}, {
					cancel: "1"
				}, function(err) {
					res.send(txid)
				})
			} else if(ttype == "group") {
				Mygtask.update({
					tid: tid
				}, {
					cancel: "1"
				}, function(err) {
					res.send(txid)
				})
			}
		})

	},
	async completeSingleTask(req, res) {
		var tid = "task" + req.query.tid
		var name = req.query.name
		var uid = ""
		var pro = queryId(name, "person")
		pro.then((id) => {
			uid = id
			tx_id = fabric_client.newTransactionID();
			console.log("Assigning transaction_id: ", tx_id._transaction_id);
			const request = {
			    chaincodeId: 'bank',
			    txId: tx_id,
			    fcn: 'completeSingleTask',
			    args: [tid, uid],
			    chainId: channel_name
			};
			invoke_chaincode(request, async function(txid) {
				var temp = await Mytask.update({
					tid: tid
				}, {
					time_complete: new Date().getTime()
				})

				var doc = Mytask.findOne({tid: tid})
				var user = Myuser.findOne({uid: uid})
				let new_mess = new Mymessage({
						name: doc.publisher,
						type: doc.type,
						post: "system",
						post_type: "system",
						mess: "用户:" + name + "确认了您发布的任务task" + tid + ",其信息为:地址:" + user.address + ",电话:" + user.phone,
						time: timestampToTime(new Date().getTime()),
						read: "0"
					})
				var temp2 = await new_mess.save()
				let new_tx = new MyTx({
					time: timestampToTime(new Date().getTime()),
					txid: txid,
					behavior: "确认任务",
					peer: "",
					publisher: name,
					object: tid
				})
				var temp3 = await new_tx.save()
				res.send(txid)
			})
		})
	},
	async acceptSingleTask(req, res) {
		var tid = "task" + req.query.tid
		var name = req.query.name
		var uid = ""
		var pro = queryId(name, "person")
		pro.then((id) => {
			uid = id
			tx_id = fabric_client.newTransactionID();
			console.log("Assigning transaction_id: ", tx_id._transaction_id);
			const request = {
			    chaincodeId: 'bank',
			    txId: tx_id,
			    fcn: 'acceptSingleTask',
			    args: [tid, uid, name],
			    chainId: channel_name
			};
			invoke_chaincode(request, async function(txid) {
				var temp = await Mytask.update({
					tid: tid
				}, {
					time_accept: new Date().getTime(),
					owner: name
				})
				var doc = Mytask.findOne({tid: tid})
				var user = Myuser.findOne({uid: uid})
				let new_mess = new Mymessage({
						name: doc.publisher,
						type: doc.type,
						post: "system",
						post_type: "system",
						mess: "用户:" + name + "接受了您发布的任务task" + tid + ",其信息为:地址:" + user.address + ",电话:" + user.phone,
						time: timestampToTime(new Date().getTime()),
						read: "0"
					})
				var temp2 = await new_mess.save()
				let new_tx = new MyTx({
					time: timestampToTime(new Date().getTime()),
					txid: txid,
					behavior: "接受任务",
					peer: "",
					publisher: name,
					object: tid
				})
				var temp3 = await new_tx.save()
				res.send(txid)
			})
		})
	},
	async acceptGroupTask(req, res) {
		var tid = "task" + req.query.tid
		var name = req.query.name
		var uid = queryId(name, "person")
		tx_id = fabric_client.newTransactionID();
		console.log("Assigning transaction_id: ", tx_id._transaction_id);
		const request = {
		    chaincodeId: 'bank',
		    txId: tx_id,
		    fcn: 'acceptSingleTask',
		    args: [tid, uid, name],
		    chainId: channel_name
		};
		invoke_chaincode(request, function(txid) {
			Mygtask.findOne({
				tid: tid
			}, function(err, doc) {
				Mygtask.update({
					tid: tid
				}, {
					time_accept: new Date().getTime(),
					owner: doc.owner + "," + name
				}, function(err) {
					res.send(txid)
				})
			})
		})
	},
	async completeGroupTask(req, res) {
		var tid = "task" + req.query.tid
		tx_id = fabric_client.newTransactionID();
		console.log("Assigning transaction_id: ", tx_id._transaction_id);
		const request = {
		    chaincodeId: 'bank',
		    txId: tx_id,
		    fcn: 'completeGroupTask',
		    args: [tid],
		    chainId: channel_name
		};
		invoke_chaincode(request, function(txid) {
			Mygtask.update({
				tid: tid
			}, {
				time_complete: new Date().getTime()
			}, function(err) {
				res.send(txid)
			})
		})
	},
	async queryConfirm(req, res) {
		var tid = "task" + req.query.tid
		Mytask.findOne({
			tid: tid
		}, function(err, docs) {
			res.send(docs)
		})
	},

	async confirm(req, res) {
		var tid = "task" + req.query.tid
		var type = req.query.type
		var who = req.query.who
		var confirm = req.query.confirm
		var name = req.query.name
		var utype = req.query.utype

		var flag = ""
		if(who == "p") {
			var nothing = await Mytask.update({
				tid: tid
			}, {
				pconfirm: confirm
			})
			var doc = await Mytask.findOne({tid: tid});
			flag = doc.aconfirm
			name = doc.owner

		} else {
			var nothing = await Mytask.update({
				tid: tid
			}, {
				aconfirm: confirm
			})
			var doc = await Mytask.findOne({tid: tid});
			flag = doc.pconfirm
			name = doc.owner
		}
		if(flag == 0) {
			res.send("confirm success")
			return
		} else if(flag == confirm) {
			//complete
			if(type == "single") {
				var uid = ""
				var pro = queryId(name, "person")
				pro.then((id) => {
					uid = id
					tx_id = fabric_client.newTransactionID();
					console.log("!!!!!!!!!!!!")
					console.log("tid:"+tid)
					console.log("uid:"+uid)
					console.log("Assigning transaction_id: ", tx_id._transaction_id);
					const request = {
						chaincodeId: 'bank',
						txId: tx_id,
						fcn: 'completeSingleTask',
						args: [tid, uid],
						chainId: channel_name
					};
					invoke_chaincode(request, function(txid) {
						Mytask.update({
							tid: tid
						}, {
							time_complete: new Date().getTime()
						}, function(err) {
							res.send(txid)
						})
					})
				})
			} else {
				res.send("wrong type")
			}
		} else {
			//recordDisputedTask
			tx_id = fabric_client.newTransactionID();
			console.log("Assigning transaction_id: ", tx_id._transaction_id);
			const request = {
				chaincodeId: 'bank',
				txId: tx_id,
				fcn: 'recordDisputedTask',
				args: [tid],
				chainId: channel_name
			};
			invoke_chaincode(request, function(txid) {
				return res.send("recordDisputedTask successfully")
			})
		}
	},

	async recordDisputedTask(req, res) {
		var tid = "task" + req.query.tid
		tx_id = fabric_client.newTransactionID();
		console.log("Assigning transaction_id: ", tx_id._transaction_id);
		const request = {
		    chaincodeId: 'bank',
		    txId: tx_id,
		    fcn: 'recordDisputedTask',
		    args: [tid],
		    chainId: channel_name
		};
		invoke_chaincode(request, async function(txid) {
			let new_tx = new MyTx({
					time: timestampToTime(new Date().getTime()),
					txid: txid,
					behavior: "记录争议任务",
					peer: "",
					publisher: tid,
					object: tid
				})
				var temp3 = await new_tx.save()
			return res.send("record successfully")
		})
	},

	queryPeoplePublish(req, res) {
		var name = req.query.name
		console.log(name)
		Mytask.find({
			publisher: name
		}, function(err, docs) {
			res.send(JSON.stringify(docs))
		})
	},
	queryPeopleAccept(req, res) {
		var name = req.query.name
		Mytask.find({
			owner: name
		}, function(err, docs) {
			res.send(docs)
		})
	},
	queryPeopleComplete(req, res) {
		var name = req.query.name
		Mytask.find({
			owner: name
		}, function(err, docs) {
			var temp = []
			for (var i = 0; i < docs.length; i++) {
				console.log(JSON.stringify(docs[i]))
				if(docs[i].time_complete != "none") {
					temp.push(docs[i])
				}
			}
			res.send(temp)
		})
	},
	async querypeopleDistued(req, res) {
		var name = req.query.name
		var result = []
		var docs1 = await Mytask.find({
			publisher: name
		})
		for (var i = 0; i < docs1.length; i++) {
			console.log(JSON.stringify(docs1[i]))
			if(docs1[i].aconfirm != 0 && docs1[i].pconfirm != 0 && docs1[i].pconfirm != docs1[i].aconfirm) {
				result.push(docs1[i])
			}
		}
		var docs2 = await Mytask.find({
			owner: name
		})
		for (var i = 0; i < docs2.length; i++) {
			console.log(JSON.stringify(docs2[i]))
			if(docs2[i].aconfirm != 0 && docs2[i].pconfirm != 0 && docs2[i].pconfirm != docs2[i].aconfirm) {
				result.push(docs2[i])
			}
		}
		res.send(result)
	},
	async queryBlockInfo(req, res) {
		var info = await channel.queryInfo()
		var cname = await channel.getName()
		var orderers = await channel.getOrderers()
		var orgs = await channel.getOrganizations()
		var peers = await channel.getPeers()
		var block = await channel.queryBlock(3)
		var response = {}
		response.orderers = []
		response.peers = []
		for (var i = orderers.length; i--;) {
			var temp = {}
			temp.name = orderers[i]._name
			temp.url = orderers[i]._url
			response.orderers.push(temp)
		}
		for (var i = peers.length; i--;) {
			response.peers.push(peers[i]._name)
		}

		response.height = info.height.low
		response.name = cname
		response.orgs = orgs
		res.send(response)
	},
	async queryTransaction(req, res) {
		var tid = req.query.id
		var doc = await MyTx.findOne({txid: tid})
		var result = await channel.queryTransaction(tid)
		doc.peer = result.transactionEnvelope.payload.data.actions[0].header.creator.Mspid
		res.send(doc)
	}

}
