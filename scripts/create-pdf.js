import PDFDocument from 'pdfkit';
import fs from 'fs';
import fsPromises from 'fs/promises';
import path from 'path';
import { fileURLToPath } from 'url';
import { glob } from 'glob';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const projectRoot = path.join(__dirname, '..');

const OUTPUT_DIR = path.join(projectRoot, 'output', 'print');
const CARDS_DIR = path.join(projectRoot, 'output', 'cards');

// Page dimensions (Letter size: 8.5" x 11" in points, 72 pts per inch)
const PAGE_WIDTH = 8.5 * 72;  // 612 pts
const PAGE_HEIGHT = 11 * 72;  // 792 pts

// Card dimensions in points (2.5" x 3.5" + bleed converted to display size)
// Original card images are 820x1120px, we scale them down for PDF
const CARD_WIDTH = 2.5 * 72;   // 180 pts (without bleed for display)
const CARD_HEIGHT = 3.5 * 72;  // 252 pts (without bleed for display)

// Image dimensions include bleed, so we need the full size
const IMG_SCALE = 0.22;  // Scale factor to fit cards on page
const SCALED_CARD_WIDTH = 820 * IMG_SCALE;   // ~180 pts
const SCALED_CARD_HEIGHT = 1120 * IMG_SCALE; // ~246 pts

// Layout: 3 columns x 3 rows
const CARDS_PER_ROW = 3;
const CARDS_PER_COL = 3;
const CARDS_PER_PAGE = 9;

// Margins
const MARGIN_X = (PAGE_WIDTH - (SCALED_CARD_WIDTH * CARDS_PER_ROW)) / 2;
const MARGIN_Y = (PAGE_HEIGHT - (SCALED_CARD_HEIGHT * CARDS_PER_COL)) / 2;

// Crop mark settings
const CROP_MARK_LENGTH = 15;
const CROP_MARK_OFFSET = 5;

function drawCropMarks(doc, x, y, width, height) {
  doc.save();
  doc.strokeColor('#000000');
  doc.lineWidth(0.5);

  // Top-left corner
  doc.moveTo(x - CROP_MARK_OFFSET, y)
     .lineTo(x - CROP_MARK_OFFSET - CROP_MARK_LENGTH, y)
     .stroke();
  doc.moveTo(x, y - CROP_MARK_OFFSET)
     .lineTo(x, y - CROP_MARK_OFFSET - CROP_MARK_LENGTH)
     .stroke();

  // Top-right corner
  doc.moveTo(x + width + CROP_MARK_OFFSET, y)
     .lineTo(x + width + CROP_MARK_OFFSET + CROP_MARK_LENGTH, y)
     .stroke();
  doc.moveTo(x + width, y - CROP_MARK_OFFSET)
     .lineTo(x + width, y - CROP_MARK_OFFSET - CROP_MARK_LENGTH)
     .stroke();

  // Bottom-left corner
  doc.moveTo(x - CROP_MARK_OFFSET, y + height)
     .lineTo(x - CROP_MARK_OFFSET - CROP_MARK_LENGTH, y + height)
     .stroke();
  doc.moveTo(x, y + height + CROP_MARK_OFFSET)
     .lineTo(x, y + height + CROP_MARK_OFFSET + CROP_MARK_LENGTH)
     .stroke();

  // Bottom-right corner
  doc.moveTo(x + width + CROP_MARK_OFFSET, y + height)
     .lineTo(x + width + CROP_MARK_OFFSET + CROP_MARK_LENGTH, y + height)
     .stroke();
  doc.moveTo(x + width, y + height + CROP_MARK_OFFSET)
     .lineTo(x + width, y + height + CROP_MARK_OFFSET + CROP_MARK_LENGTH)
     .stroke();

  doc.restore();
}

async function createPrintPDF(cardFiles, outputFilename) {
  const doc = new PDFDocument({
    size: 'letter',
    margin: 0,
    autoFirstPage: false,
  });

  const outputPath = path.join(OUTPUT_DIR, outputFilename);
  const writeStream = fs.createWriteStream(outputPath);
  doc.pipe(writeStream);

  for (let i = 0; i < cardFiles.length; i += CARDS_PER_PAGE) {
    doc.addPage();

    const cardsOnPage = cardFiles.slice(i, i + CARDS_PER_PAGE);

    for (let j = 0; j < cardsOnPage.length; j++) {
      const row = Math.floor(j / CARDS_PER_ROW);
      const col = j % CARDS_PER_ROW;

      const x = MARGIN_X + (col * SCALED_CARD_WIDTH);
      const y = MARGIN_Y + (row * SCALED_CARD_HEIGHT);

      // Add card image
      doc.image(cardsOnPage[j], x, y, {
        width: SCALED_CARD_WIDTH,
        height: SCALED_CARD_HEIGHT,
      });

      // Draw crop marks around each card
      drawCropMarks(doc, x, y, SCALED_CARD_WIDTH, SCALED_CARD_HEIGHT);
    }

    // Add page number
    doc.fontSize(8)
       .fillColor('#666666')
       .text(
         `Page ${Math.floor(i / CARDS_PER_PAGE) + 1} of ${Math.ceil(cardFiles.length / CARDS_PER_PAGE)}`,
         0,
         PAGE_HEIGHT - 20,
         { align: 'center', width: PAGE_WIDTH }
       );
  }

  doc.end();

  return new Promise((resolve, reject) => {
    writeStream.on('finish', resolve);
    writeStream.on('error', reject);
  });
}

async function main() {
  console.log('Subway Deal PDF Generator');
  console.log('=========================\n');

  await fsPromises.mkdir(OUTPUT_DIR, { recursive: true });

  // Get all card PNG files
  const cardFiles = await glob(`${CARDS_DIR}/*.png`);
  cardFiles.sort(); // Ensure consistent ordering

  console.log(`Found ${cardFiles.length} cards to layout\n`);

  if (cardFiles.length === 0) {
    console.error('No card files found! Run npm run generate-cards first.');
    process.exit(1);
  }

  // Create main print PDF with all cards
  await createPrintPDF(cardFiles, 'subway-deal-print-all.pdf');
  console.log(`Created: subway-deal-print-all.pdf (${Math.ceil(cardFiles.length / CARDS_PER_PAGE)} pages)`);

  // Create separate PDFs by card type
  const propertyCards = cardFiles.filter(f => path.basename(f).startsWith('prop_'));
  const wildcardCards = cardFiles.filter(f => path.basename(f).startsWith('wild_'));
  const actionCards = cardFiles.filter(f => path.basename(f).startsWith('action_'));
  const rentCards = cardFiles.filter(f => path.basename(f).startsWith('rent_'));
  const moneyCards = cardFiles.filter(f => path.basename(f).startsWith('money_'));

  if (propertyCards.length > 0) {
    await createPrintPDF(propertyCards, 'subway-deal-properties.pdf');
    console.log(`Created: subway-deal-properties.pdf (${Math.ceil(propertyCards.length / CARDS_PER_PAGE)} pages)`);
  }

  if (wildcardCards.length > 0) {
    await createPrintPDF(wildcardCards, 'subway-deal-wildcards.pdf');
    console.log(`Created: subway-deal-wildcards.pdf (${Math.ceil(wildcardCards.length / CARDS_PER_PAGE)} pages)`);
  }

  if (actionCards.length > 0) {
    await createPrintPDF(actionCards, 'subway-deal-actions.pdf');
    console.log(`Created: subway-deal-actions.pdf (${Math.ceil(actionCards.length / CARDS_PER_PAGE)} pages)`);
  }

  if (rentCards.length > 0) {
    await createPrintPDF(rentCards, 'subway-deal-rent.pdf');
    console.log(`Created: subway-deal-rent.pdf (${Math.ceil(rentCards.length / CARDS_PER_PAGE)} pages)`);
  }

  if (moneyCards.length > 0) {
    await createPrintPDF(moneyCards, 'subway-deal-money.pdf');
    console.log(`Created: subway-deal-money.pdf (${Math.ceil(moneyCards.length / CARDS_PER_PAGE)} pages)`);
  }

  console.log('\nPDF generation complete!');
}

main().catch(console.error);
