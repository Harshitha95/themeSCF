/** var express = require('express');
var app = express();
var invoke = require('./routes/invoke')
var query = require('./routes/query')
var cors = require('cors')

var http = require('http');
var https = require('https');
var path = require('path');
var fs = require('fs');
var mime = require('mime');
var bodyParser = require('body-parser');
var cfenv = require('cfenv');

var cookieParser = require('cookie-parser');
var session = require('express-session');

var uuid = require('uuid');

//var sessionSecret = env.sessionSecret;
var appEnv = cfenv.getAppEnv();
var vcapServices = require('vcap_services');

var connect = require("connect");
var serveStatic = require('serve-static');
//var app = connect().use(connect.static(__dirname + '/'));
var app = connect().use(connect.static(__dirname + '/'));
app.listen(8180);
var busboy = require('connect-busboy');
app.use(busboy());

app.use(bodyParser.urlencoded({ extended: true }));
app.use(bodyParser.json());
app.set('appName', 'themechainSCF2Golang');
app.set('port', appEnv.port);//app.use(cookieParser(sessionSecret));

/**app.set('views', path.join(__dirname + '/'));
app.engine('html', require('ejs').renderFile);
//app.set('view engine', 'ejs');
app.use(express.static(__dirname + '/'));
//app.use(bodyParser.json());*
////
app.get('/', function(req, res) {
    res.sendFile(path.join(__dirname + '/addBank.html'));
});
app.set('views', path.join(__dirname, 'html'));
app.set('view engine', 'html');
app.engine('html', require('ejs').renderFile);
app.use(express.static(__dirname + '/'))
app.use(bodyParser.json());

var app = express();
app.use(cors());
app.use('/invoke', invoke);
app.use('/query', query);

app.listen(3000);
console.log("app listening at localhost:3000");*/
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

app.use('/invoke', invoke);
app.use('/query', query);
app.use('/postSender',postRequestHandler);
app.use('/getSender',getRequestHandler);
app.listen(3000);
console.log("app listening at localhost:3000");
app.use(express.static(__dirname + '/HTML'));
app.get('/user', function(req, res){
console.log("User ")
//$('#results').html(data);
	res.send("user1");
})
app.get('/', function(req, res){
   // res.sendFile('index.html', { root: __dirname + "/HTML" } );   
});
//app.use('/Images', express.static(__dirname + '/HTML/Images'));
//app.use('/js', express.static(__dirname + '/HTML/js'));

 //adding static functionality for images
/** app.use(express.static('public'));


 // Renders the "image" page (at views/index.ejs) when GETing the URL
 app.get("/image", function(request, response) {
  response.render("image");
 });**/
