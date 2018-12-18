var express = require('express');
var app = express();
var invoke = require('./routes/invoke')
var query = require('./routes/query')
var postRequestHandler =require('./routes/postRequestHandle')
var getRequestHandler =require('./routes/getRequestHandle')
var http = require('http');
var https = require('https');
var path = require('path');
var fs = require('fs');
var mime = require('mime');
//var findme= require('find-me');
var bodyParser = require('body-parser');
var cfenv = require('cfenv');
var tryRequire = require('try-require');
var registerUserImport = tryRequire.resolve('./registerUser.js');
//console.log("user in import "+ user)
/* if( user  != nil) {
    var user = tryRequire('./registerUser.js');    
} */
//var user = tryRequire.('./registerUser.js');
//var registerUser1 =require('./registerUser').member_user;
var cookieParser = require('cookie-parser');
var session = require('express-session');

var vcapServices = require('vcap_services');
var uuid = require('uuid');


//var sessionSecret = env.sessionSecret;
var appEnv = cfenv.getAppEnv();

var busboy = require('connect-busboy');
app.use(busboy());

app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());
app.set('appName', 'themechainSCF2Golang');
app.set('port', appEnv.port);//app.use(cookieParser(sessionSecret));

app.set('views', path.join(__dirname + '/'));
app.engine('html', require('ejs').renderFile);
app.set('view engine', 'ejs');
app.use(express.static(__dirname + '/'));
app.use(bodyParser.json());
var globals = require('node-global-storage');
app.use('/invoke', invoke);
app.use('/query', query);
app.use('/postSender',postRequestHandler);
app.use('/getSender',getRequestHandler);
app.listen(3000);
console.log("app listening at localhost:3000");
app.use(express.static(__dirname + '/HTML'));
/* app.get('/user', function(req, res){
/* console.log("registerUser1 ", app.locals.name ) 

var userName1 = globals.get('userName1'); 
console.log("userName1 ", app.locals.name ) 
    var uNAme= registerUserImport.getUserName;
    res.send(uNAme);
   //  res.send()
}) */
app.get('/', function(req, res){
   // res.sendFile('index.html', { root: __dirname + "/HTML" } );   
});

