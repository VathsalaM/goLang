var dataObject = {
    "name" : "vathsala",
    "age" : "20"
}

var getDataObject = function(){
    return {
    firstName:$("#firstName").val(),
    lastName:$("#lastName").val(),
    Email:$("#Email").val(),
    gender:$("#gender").val(),
    address1:$("#address1").val(),
    address2:$("#address2").val()
    }
}

var onload = function(){
    console.log("hello........on load..........")
    $("#submit").click(function(){
        console.log("java called........")
        $.post("/register",getDataObject(),function(data){
            console.log("=======================    ",data)
            $("#successPage").append("p").append(data)
        })
    })
}

$.ready(onload)