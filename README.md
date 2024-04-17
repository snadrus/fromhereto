# fromhereto

Go deps getting you down? 

Clone this repo & fix it. How To:
   mkdir ~/scratch
   cd ~/scratch
   git clone github.com/snadrus/fromhereto
   cd <<YOUR_MAIN_FOLDER>>
   go run ~/scratch/fromhereto/fromhereto.go . >deps.json

Open that in your favorite editor. It's all the imports for all the packages your main and all its dependencies import. 

## Wanna determine why you depend on X? 
Find it in the JSON and trace back its parent(s) until you find where that dep shouldn't be. 

## Wanna cut down a big tree?
Find where there is only 1 reference to something (that one's a bit harder)
