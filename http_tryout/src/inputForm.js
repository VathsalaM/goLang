var onload = function(){
    $("button").click(function(){
        console.log('button clicked')
        $.post('/view')
    })
}
$.ready(onload)