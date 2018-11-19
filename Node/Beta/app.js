/*
 * basically what this app  does is
 * it'll call node functions 
 * invoke : to invoke a chaincode passing arguments
 * query : to query a chaincode passing arguments
 */

 var query = require('query');
 var invoke = require('invoke');

//Body-Parser
var bodyParser = require('body-parser');

//Express working
var express = require('express');

//now with 'app' we can use  all of the methods
var app = express();

//var jsonParser = bodyParser.json();

var urlencodeParser = bodyParser.urlencoded({ extended: false});

app.post('/invoke', urlencodeParser, function(req, res) {
    if(!req.body) return res.sendStatus(400)
    console.log(req.body);
    res.render('invoke-process', {data: req.body})
});

//set view-engine
app.set('view engine', 'ejs');

//handling staticFiles (as in assets/) using middleware
app.use('/assets', express.static('assets'));

//sendFile to send html pages
app.get('/', function(req, res) {
    res.render('invoke-form');
});
app.get('', function(req, res) {
    res.render('invoke-form');
});


//render a view using template-engine
app.get('/invoke', function(req, res) {
    var data = {};
    res.render('invoke-form');
});

//express app running in port 3000
app.listen(3000);
