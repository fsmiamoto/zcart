MAIN=main
rm -f $MAIN.{aux,dvi,lof,log,lot,lsb,lsg,pdf,toc,bbl,blg}
rm -f ./out/*
pdflatex $MAIN
bibtex $MAIN
pdflatex $MAIN
pdflatex $MAIN
mv $MAIN.pdf out/