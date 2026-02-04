# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Subway Deal is a parody card game inspired by Monopoly Deal, themed around the NYC MTA subway system. The project is currently in the design phase with no code implementation yet.

**Goal:** Generate 106 printable playing cards (PNG/PDF) from structured data and templates.

## Planned Architecture

The project will use a data-driven card generation pipeline:

```
data/*.json → templates/*.html + styles.css → scripts/generate-cards.js → output/cards/*.png
                                            → scripts/create-pdf.js → output/print/*.pdf
```

**Recommended tech stack:** Node.js with HTML/CSS templates rendered via Puppeteer.

## Card Types (106 total)

| Type | Count | Data File |
|------|-------|-----------|
| Property cards | 28 | properties.json |
| Property wildcards | 11 | wildcards.json |
| Action cards | 34 | actions.json |
| Rent cards | 13 | rent.json |
| Money cards | 20 | money.json |

## Print Specifications

- Card size: 2.5" × 3.5" (63mm × 88mm)
- Bleed: 3mm all sides
- Resolution: 300 DPI minimum
- Color space: CMYK for print, RGB for digital

## MTA Color Mapping

Colors must match official MTA signage:
- Brown (#996633): J/Z lines
- Blue (#0039A6): A/C/E lines
- Gray (#808183): Shuttles
- Orange (#FF6319): B/D/F lines
- Red (#EE352E): 1/2/3 lines
- Yellow (#FCCC0A): N/Q/R lines
- Green (#00933C): Hub stations
- Black (#000000): Railroads (LIRR, Metro-North, NJ Transit, PATH)
- Light Green (#6CBE45): Utilities (G/L lines)

## Key Design Document

`thoughts/initial_specs.md` contains the complete game design specification including all card data, game rules, visual design concepts, and implementation phases.
