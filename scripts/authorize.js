var addr = "http://yourAddressHere!!!/";
var inputField = document.getElementById("input");
inputField.onkeydown = 

function authUser(){
    if(event.keyCode==13){
        if(inputField.value !="")
            authorize(addr+"authorize?username="+inputField.value)
    }
}




function authorize(address){
    xmlhttp = new XMLHttpRequest();
    xmlhttp.onreadystatechange = function(){
        if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
            window.location = addr + xmlhttp.responseText
            alert("Auth completed!")
        }
    };
    xmlhttp.open("GET", address, true);
    xmlhttp.send();
}

