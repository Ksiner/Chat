//import {$,jQuery} from 'jQuery'




var addr = "http://90.188.151.177/";
var messageAddr = "";
var usersAddr = "";
var currentUser = "";
var targetUser = "";
var root = document.getElementById("root");
var dataArray = new Array();
var userArray = new Array();

var messageContainer = document.getElementsByClassName("messageDiv")[0];
var formCurrentTargetUser = document.getElementById("currentTargetUser")
var userContainer = document.getElementsByClassName("userDiv")[0];
var messageInputDiv = document.getElementsByClassName("messageInputDiv")[0];
var messageInputArea = document.getElementsByClassName("messageInputArea")[0];
messageInputArea.onkeyup = sendMessage;


ParseURLArgs(addr);



function ParseURLArgs(addres){
    let url = new URL(window.location.href)
    currentUser = url.searchParams.get("username")
    if (currentUser == null){
        window.location = addres
        return
    }
    targetUser = url.searchParams.get("targetuser")
    if (targetUser==null){
        targetUser = ""
    }
    messageAddr = addr+"messages/request?username="+currentUser+"&targetuser="+targetUser
    usersAddr = addr+"users/request?username="+currentUser
    userTemplateInit();
    messageTemplateInit();
}

//var messageGetUser = addr+"/messages/users";


// var messageGetAddr = addr+"/messages/get";
// var messagePostAddr = addr+"/messages/send";

function updateMessageAddr(){
    messageAddr = addr+"messages/request?username="+currentUser+"&targetuser="+targetUser;
}


function changeTargetUser(userDiv){
    return function(){ 
        if(targetUser == userDiv.childNodes[0].innerText){
            return;
        }
        targetUser = userDiv.childNodes[0].innerText;
        dataArray = new Array();
        while(messageContainer.firstChild){
            messageContainer.removeChild(messageContainer.firstChild);
        }
        formCurrentTargetUser.innerText = targetUser;
        updateMessageAddr();
    }
}

function userTemplateInit(){
    ajax_get_users(usersAddr,userArray,callbackGet,addUserElemetsToItsDivs,insertUserDivs);
    setTimeout(userTemplateInit,2000);
}

function messageTemplateInit(){
    ajax_get_messages(messageAddr,dataArray, callbackGet,addMessageElemetsToItsDivs,insetrIntoDivFileds);
    setTimeout(messageTemplateInit,1000);
}


function unixToDate(unix_timestamp){
    let date = new Date(unix_timestamp);
    let days = date.getDate();
    let months = date.getMonth()+1;
    let hours = date.getHours();
    let minutes = date.getMinutes();
    let seconds = date.getSeconds();
    if(seconds<10)
        seconds = "0"+seconds;
    if(minutes<10)
        minutes = "0"+minutes
    if(hours<10)
        hours = "0"+hours
    if(months<10)
        months = "0"+months
    if(days<10)
        days = "0"+ days
    return days+"."+months+"."+date.getFullYear()+" "+hours+":"+minutes+":"+seconds;
}

function insetrIntoDivFileds(div,element){
    div.childNodes[0].innerText = element.userfrom;
    div.childNodes[1].innerText = element.message;
    div.childNodes[2].innerText = unixToDate(element.date);
    messageContainer.appendChild(div);
}

function insertUserDivs(element){
    let userDiv = document.createElement("div");
    userDiv.className = "container user"
    let userDivChild = document.createElement("p");
    userDivChild.innerText=element.Name;
    userDivChild.onclick = changeTargetUser(userDiv);
    userDiv.appendChild(userDivChild);
    userContainer.appendChild(userDiv);
}

function addMessageElemetsToItsDivs(element,elemAddFunc){
    if (element.userfrom === currentUser) {
        debugger;
        let myMessageDv = initDiv("container", "time-right");
        elemAddFunc(myMessageDv,element);
    } else {
        debugger;
        let = messageDiv = initDiv("container darker", "time-left");
        elemAddFunc(messageDiv,element);
    }
}
function addUserElemetsToItsDivs(element,elemAddFunc){
        elemAddFunc(element);
}

function callbackGet(data,cacheArray,elemTypeFunc,elemAddFunc) {
    data.forEach(element => {
        if(cacheArray.find(cacheElem=>JSON.stringify(cacheElem)==JSON.stringify(element))==undefined){
            cacheArray.push(element);
            elemTypeFunc(element,elemAddFunc)
        } 
    });
}





function ajax_get_messages(url,cacheArray,callback,elemTypeFunc,elemAddFunc) {
    let xmlhttp = new XMLHttpRequest();
    xmlhttp.onreadystatechange = function () {
        if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
            try {
                var data = JSON.parse(xmlhttp.responseText);
            } catch (err) {
                console.log(err.message + " in " + xmlhttp.responseText);
            }
            if(elemAddFunc !== undefined){
                callback(data,cacheArray,elemTypeFunc,elemAddFunc)
            }else{
                callback(data);
            }
        }
    };
    xmlhttp.open("GET", url, true);
    xmlhttp.send();
}

function ajax_get_users(url,cacheArray,callback,elemTypeFunc,elemAddFunc) {
    let xmlhttp = new XMLHttpRequest();
    xmlhttp.onreadystatechange = function () {
        if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
            try {
                var data = JSON.parse(xmlhttp.responseText);
            } catch (err) {
                console.log(err.message + " in " + xmlhttp.responseText);
            }
            if(elemAddFunc !== undefined){
                callback(data,cacheArray,elemTypeFunc,elemAddFunc)
            }else{
                callback(data);
            }
        }
    };
    xmlhttp.open("GET", url, true);
    xmlhttp.send();
}

// function ajax_get_users(url,callback){
//     let xmlhttp = new XMLHttpRequest();
//     xmlhttp.open("GET",url,true);
//     xmlhttp.onreadystatechange = function(){
//         if (xmlhttp.readyState === 4 && xmlhttp.status === 200){
//             callback(JSON.parse(xmlhttp.responseText),insertUserDivs);
//         }
//     }
//     xmlhttp.send();
// }

function ajax_post_message(Obj,url){
    let xmlhttp = new XMLHttpRequest();
    xmlhttp.open("POST",url,true);
    xmlhttp.setRequestHeader("Content-Type","application/json; charset=utf-8");
    xmlhttp.onreadystatechange = function(){
        if (xmlhttp.readyState === 4 && xmlhttp.status === 200){
            //alert("Message sended!");
        }
    }
    xmlhttp.send(JSON.stringify(Obj));
}

function init_elem(elemTag, className) {
    let elem = document.createElement(elemTag);
    elem.className = className;
    return elem;
}

function initDiv(divClassName, spanClassName) {
    debugger;
    let div = init_elem("div", divClassName);
    div.appendChild(init_elem("p", "userName"));
    div.appendChild(init_elem("p", "messageText"));
    div.appendChild(init_elem("span", spanClassName));
    return div;
}

function sendMessage(event){
    if(event.keyCode==13){
        if (targetUser=="")
            return
        obj = {
            userfrom: currentUser,
            userto: targetUser,
            message: messageInputArea.value/*.replace('\n','')*/,
            date: Date.now()
        };
        ajax_post_message(obj,messageAddr);
        messageInputArea.value = "";
    }
}
