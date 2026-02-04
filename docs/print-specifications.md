# Print Service Specifications

## Card Specifications

- **Size:** 2.5" x 3.5" (63mm x 88mm) - Poker size
- **Stock:** 330gsm Black Core or 300gsm Premium
- **Finish:** Linen texture recommended for shuffle feel
- **Coating:** Gloss or UV coating for durability
- **Bleed:** 3mm all sides
- **Resolution:** 300 DPI
- **Color:** CMYK (converted from RGB source)

## Generated Files

| File | Contents | Pages |
|------|----------|-------|
| `subway-deal-print-all.pdf` | All 106 cards | 12 |
| `subway-deal-properties.pdf` | 28 property cards | 4 |
| `subway-deal-wildcards.pdf` | 11 wildcard cards | 2 |
| `subway-deal-actions.pdf` | 34 action cards | 4 |
| `subway-deal-rent.pdf` | 13 rent cards | 2 |
| `subway-deal-money.pdf` | 20 money cards | 3 |

## Recommended Print Services

### The Game Crafter
- **URL:** thegamecrafter.com
- **MOQ:** 1 deck
- **Card stock:** 12pt (270gsm)
- **Turnaround:** 5-7 business days
- **Cost estimate:** ~$15-20 per deck
- **Notes:** Great for prototyping, US-based

### DriveThruCards
- **URL:** drivethrucards.com
- **MOQ:** 1 deck
- **Card stock:** 300gsm with linen finish
- **Turnaround:** 2-3 weeks
- **Cost estimate:** ~$12-18 per deck
- **Notes:** Part of DriveThruRPG network

### MakePlayingCards
- **URL:** makeplayingcards.com
- **MOQ:** 1 deck
- **Card stock:** 330gsm Black Core (recommended)
- **Finish:** Multiple options including holographic
- **Turnaround:** 7-10 business days
- **Cost estimate:** ~$18-25 per deck
- **Notes:** Best quality for final production

### PrinterStudio
- **URL:** printerstudio.com
- **MOQ:** 1 deck
- **Card stock:** Various options
- **Turnaround:** 5-10 business days
- **Cost estimate:** ~$15-22 per deck
- **Notes:** Good for custom card sizes

## Upload Instructions

1. Export PDFs from `output/print/`
2. Most services prefer individual card images rather than PDFs
   - Use files from `output/cards/` for individual uploads
3. Verify color mode in upload preview
4. Confirm card dimensions match service requirements
5. Order single proof deck before bulk production
6. Review proof for color accuracy and cut alignment

## Individual Card Upload

For services that require individual card images:

```bash
# Cards are in output/cards/*.png
# 820x1120 pixels at 300 DPI equivalent
# Includes 3mm bleed on all sides
```

## Quality Checklist for Proof Deck

When you receive your proof deck, verify:

- [ ] Card quality and finish meet expectations
- [ ] Colors match design intent (accounting for CMYK conversion)
- [ ] Text is legible and sharp
- [ ] Cards are cut accurately with minimal misalignment
- [ ] Bleed areas are trimmed correctly
- [ ] Cards shuffle and handle well
- [ ] No visible printing artifacts or banding
- [ ] Corners are properly rounded

## Color Notes

The cards use MTA official colors which are designed for signage (RGB/sRGB). When converting to CMYK for print:

| Color | RGB Hex | Notes |
|-------|---------|-------|
| Brown (J/Z) | #996633 | Converts well |
| Blue (A/C/E) | #0039A6 | May appear slightly different |
| Gray (Shuttles) | #808183 | Converts accurately |
| Orange (B/D/F) | #FF6319 | Vibrant, converts well |
| Red (1/2/3) | #EE352E | Converts accurately |
| Yellow (N/Q/R) | #FCCC0A | May appear slightly muted |
| Green (Hubs) | #00933C | Converts well |
| Dark Blue (Stadiums) | #2A344D | Converts accurately |
| Black (Railroads) | #000000 | Pure black, no conversion needed |
| Light Green (Utilities) | #6CBE45 | Converts well |

## Cost Estimation

For a typical print run:

| Quantity | Per-Deck Cost | Total |
|----------|---------------|-------|
| 1 (proof) | $20-25 | $20-25 |
| 10 | $12-15 | $120-150 |
| 50 | $8-10 | $400-500 |
| 100+ | $6-8 | $600-800 |

*Prices are estimates and vary by service, card stock, and finish options.*
