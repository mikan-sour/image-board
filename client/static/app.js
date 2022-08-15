const img = document.getElementById("good-boy-image");

const btn = document.getElementById("good-boy-button").addEventListener("click", function(){
    const currentImage = img.src.split(".")[0].slice(-1)
    let num = currentImage
    do {
        num = Math.floor(Math.random() * 5) + 1;
    } while (num == currentImage)

    const newURL = `http://localhost/s3/saint-bernards/images/sb${num}.jpg`;
    if(img.hasAttribute("src")){
        img.src = newURL;
    }
})