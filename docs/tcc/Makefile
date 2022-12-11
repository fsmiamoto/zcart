SHELL = /bin/bash
DATE = $(shell date +"%Y-%m-%d" )
MAIN = main
OUTPUT = zCart-$(DATE)

all: clean $(OUTPUT).pdf

$(OUTPUT).pdf: $(MAIN).tex
	pdflatex $(MAIN)
	bibtex $(MAIN)
	pdflatex $(MAIN)
	pdflatex $(MAIN)
	mv $(MAIN).pdf out/$(OUTPUT).pdf

clean:
	-rm -f $(MAIN).{aux,dvi,lof,log,lot,lsb,lsg,pdf,toc,bbl,blg}

.PHONY: clean all
