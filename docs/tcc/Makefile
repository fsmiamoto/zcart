SHELL = /bin/bash
MAIN = main

all: clean $(MAIN).pdf

$(MAIN).pdf: $(MAIN).tex
	pdflatex $(MAIN)
	bibtex $(MAIN)
	pdflatex $(MAIN)
	pdflatex $(MAIN)
	mv $(MAIN).pdf out/

clean:
	-rm -f $(MAIN).{aux,dvi,lof,log,lot,lsb,lsg,pdf,toc,bbl,blg}

.PHONY: clean all
