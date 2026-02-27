.PHONY: help gource gource-30s gource-15s gource-60s gource-last-year-30s gource-last-year-30s-gif gource-render gource-render-range

help:
	@$(MAKE) --no-print-directory -f Makefile.gource help

gource:
	@$(MAKE) --no-print-directory -f Makefile.gource gource

gource-30s:
	@$(MAKE) --no-print-directory -f Makefile.gource gource-30s

gource-15s:
	@$(MAKE) --no-print-directory -f Makefile.gource gource-15s

gource-60s:
	@$(MAKE) --no-print-directory -f Makefile.gource gource-60s

gource-last-year-30s:
	@$(MAKE) --no-print-directory -f Makefile.gource gource-last-year-30s

gource-last-year-30s-gif:
	@$(MAKE) --no-print-directory -f Makefile.gource gource-last-year-30s-gif

gource-render:
	@$(MAKE) --no-print-directory -f Makefile.gource gource-render

gource-render-range:
	@$(MAKE) --no-print-directory -f Makefile.gource gource-render-range
