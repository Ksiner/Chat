//import {$,jQuery} from 'jQuery'
messageTemplateInit();



var addr = "http://192.168.1.3:8080";
var currentUser = "Me"
var root = document.getElementById("root");
debugger;
var messageContainer = document.getElementsByClassName("messageDiv")[0];
var userContainer = document.getElementsByClassName("userDiv")[0];
var messageInputDiv = document.getElementsByClassName("messageInputDiv")[0];
var messageInputArea = document.getElementsByClassName("messageInputArea")[0];
messageInputArea.onkeyup = sendMessage;

var messageGetUser = addr+"/messages/users";
var messageGetAddr = addr+"/messages/get";
var messagePostAddr = addr+"/messages/send";

function messageTemplateInit(){
    ajax_get(messageGetUser,callbackGet,insertUserDivs)
    ajax_get(messageGetAddr, callbackGet,insetrIntoDivFileds);
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
    div.childNodes[0].innerText = element.user;
    div.childNodes[1].innerText = element.message;
    div.childNodes[2].innerText = unixToDate(element.date);
    messageContainer.appendChild(div);
}

function insertUserDivs(div,element){
    let userDiv = document.createElement("div");
    let userDivChild = document.createElement("p");
    userDivChild.innerText=element;
    userDiv.appendChild(userDivChild);
    userContainer.appendChild(userDiv);
}

var callbackGet = function (data,elemAddFunc) {
    data.forEach(element => {
        if (element.my) {
            debugger;
            let myMessageDv = initDiv("container", "time-right");
            elemAddFunc(myMessageDv,element);
        } else {
            debugger;
            let = messageDiv = initDiv("container darker", "time-left");
            insetrIntoDivFileds(messageDiv,element);
        }
    });
}





function ajax_get(url, callback,elemAddFunc) {
    var xmlhttp = new XMLHttpRequest();
    xmlhttp.onreadystatechange = function () {
        if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
            console.log("Responce text: " + xmlhttp.responseText);
            try {
                var data = JSON.parse(xmlhttp.responseText);
            } catch (err) {
                console.log(err.message + " in " + xmlhttp.responseText);
            }
            if(elemAddFunc !== undefined){
                callback(data,elemAddFunc)
            }else{
                callback(data);
            }
        }
    };
    xmlhttp.open("GET", url, true);
    xmlhttp.send();
}

function ajax_post(Obj,url,callback){
    var xmlhttp = new XMLHttpRequest();
    xmlhttp.open("POST",url,true);
    xmlhttp.setRequestHeader("Content-Type","application/json; charset=utf-8");
    xmlhttp.onreadystatechange = function(){
        if (xmlhttp.readyState === 4 && xmlhttp.status === 200){
            callback(initDiv("container", "time-right"),Obj);
        }
    }
    xmlhttp.send(JSON.stringify(obj));
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
        obj = {
            user: "me",
            message: messageInputArea.value/*.replace('\n','')*/,
            date: Date.now(),
            my:true
        };
        ajax_post(obj,messagePostAddr,insetrIntoDivFileds);
        alert("Message sended!");
        messageInputArea.value = "";
    }
}
